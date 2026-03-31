// Copyright (c) 2026 Lark Technologies Pte. Ltd.
// SPDX-License-Identifier: MIT

package doc

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/larksuite/cli/internal/output"
	"github.com/larksuite/cli/internal/validate"
	"github.com/larksuite/cli/shortcuts/common"
)

type documentRef struct {
	Kind  string
	Token string
}

func parseDocumentRef(input string) (documentRef, error) {
	raw := strings.TrimSpace(input)
	if raw == "" {
		return documentRef{}, output.ErrValidation("--doc cannot be empty")
	}

	if token, ok := extractDocumentToken(raw, "/wiki/"); ok {
		return documentRef{Kind: "wiki", Token: token}, nil
	}
	if token, ok := extractDocumentToken(raw, "/docx/"); ok {
		return documentRef{Kind: "docx", Token: token}, nil
	}
	if token, ok := extractDocumentToken(raw, "/doc/"); ok {
		return documentRef{Kind: "doc", Token: token}, nil
	}
	if strings.Contains(raw, "://") {
		return documentRef{}, output.ErrValidation("unsupported --doc input %q: use a docx URL/token or a wiki URL that resolves to docx", raw)
	}
	if strings.ContainsAny(raw, "/?#") {
		return documentRef{}, output.ErrValidation("unsupported --doc input %q: use a docx token or a wiki URL", raw)
	}

	return documentRef{Kind: "docx", Token: raw}, nil
}

func extractDocumentToken(raw, marker string) (string, bool) {
	idx := strings.Index(raw, marker)
	if idx < 0 {
		return "", false
	}
	token := raw[idx+len(marker):]
	if end := strings.IndexAny(token, "/?#"); end >= 0 {
		token = token[:end]
	}
	token = strings.TrimSpace(token)
	if token == "" {
		return "", false
	}
	return token, true
}

func buildDriveRouteExtra(docID string) (string, error) {
	extra, err := json.Marshal(map[string]string{"drive_route_token": docID})
	if err != nil {
		return "", output.Errorf(output.ExitInternal, "internal_error", "failed to marshal upload extra data: %v", err)
	}
	return string(extra), nil
}

func shouldFallbackToDocxOpenAPI(result map[string]interface{}) bool {
	if len(result) == 0 {
		return true
	}
	if len(result) == 1 {
		if v, ok := result["result"]; ok && v == nil {
			return true
		}
		if v, ok := result["message"]; ok {
			if s, ok := v.(string); ok && strings.TrimSpace(s) == "" {
				return true
			}
		}
	}
	return false
}

func createDocxViaOpenAPI(runtime *common.RuntimeContext, title, markdown string) (map[string]interface{}, error) {
	body := map[string]interface{}{}
	if strings.TrimSpace(title) != "" {
		body["title"] = title
	}
	data, err := runtime.CallAPI("POST", "/open-apis/docx/v1/documents", nil, body)
	if err != nil {
		return nil, err
	}
	document, _ := data["document"].(map[string]interface{})
	docID, _ := document["document_id"].(string)
	if docID == "" {
		return nil, output.Errorf(output.ExitAPI, "api_error", "docx create returned no document_id")
	}
	if strings.TrimSpace(markdown) != "" {
		if err := appendMarkdownAsPlainText(runtime, docID, markdown); err != nil {
			return nil, err
		}
	}
	return map[string]interface{}{
		"document_id": docID,
		"title":       document["title"],
		"revision_id": document["revision_id"],
		"backend":     "openapi_fallback",
	}, nil
}

func fetchDocxViaOpenAPI(runtime *common.RuntimeContext, docInput string) (map[string]interface{}, error) {
	docID, err := resolveDocxDocumentID(runtime, docInput)
	if err != nil {
		return nil, err
	}
	meta, err := runtime.CallAPI("GET", fmt.Sprintf("/open-apis/docx/v1/documents/%s", validate.EncodePathSegment(docID)), nil, nil)
	if err != nil {
		return nil, err
	}
	document, _ := meta["document"].(map[string]interface{})
	blocks, err := runtime.CallAPI("GET", fmt.Sprintf("/open-apis/docx/v1/documents/%s/blocks", validate.EncodePathSegment(docID)), nil, nil)
	if err != nil {
		return nil, err
	}
	items, _ := blocks["items"].([]interface{})
	lines := make([]string, 0, len(items))
	for _, item := range items {
		block, _ := item.(map[string]interface{})
		if block == nil {
			continue
		}
		if text := extractBlockText(block); strings.TrimSpace(text) != "" {
			lines = append(lines, text)
		}
	}
	title, _ := document["title"].(string)
	return map[string]interface{}{
		"document_id": docID,
		"title":       title,
		"markdown":    strings.Join(lines, "\n\n"),
		"backend":     "openapi_fallback",
	}, nil
}

func updateDocxViaOpenAPI(runtime *common.RuntimeContext, docInput, mode, markdown string) (map[string]interface{}, error) {
	docID, err := resolveDocxDocumentID(runtime, docInput)
	if err != nil {
		return nil, err
	}
	if mode != "append" && mode != "overwrite" {
		return nil, output.ErrValidation("private deployment fallback currently supports --mode append|overwrite")
	}
	if mode == "overwrite" {
		rootData, err := runtime.CallAPI("GET",
			fmt.Sprintf("/open-apis/docx/v1/documents/%s/blocks/%s/children", validate.EncodePathSegment(docID), validate.EncodePathSegment(docID)),
			nil, nil)
		if err != nil {
			return nil, err
		}
		children, _ := rootData["items"].([]interface{})
		if len(children) > 0 {
			if _, err := runtime.CallAPI("DELETE",
				fmt.Sprintf("/open-apis/docx/v1/documents/%s/blocks/%s/children/batch_delete", validate.EncodePathSegment(docID), validate.EncodePathSegment(docID)),
				nil,
				map[string]interface{}{"start_index": 0, "end_index": len(children)}); err != nil {
				return nil, err
			}
		}
	}
	if err := appendMarkdownAsPlainText(runtime, docID, markdown); err != nil {
		return nil, err
	}
	return map[string]interface{}{
		"document_id": docID,
		"mode":        mode,
		"backend":     "openapi_fallback",
	}, nil
}

func appendMarkdownAsPlainText(runtime *common.RuntimeContext, docID, markdown string) error {
	lines := markdownLines(markdown)
	for idx, line := range lines {
		if _, err := runtime.CallAPI("POST",
			fmt.Sprintf("/open-apis/docx/v1/documents/%s/blocks/%s/children", validate.EncodePathSegment(docID), validate.EncodePathSegment(docID)),
			nil,
			map[string]interface{}{
				"children": []interface{}{
					map[string]interface{}{
						"block_type": 2,
						"text": map[string]interface{}{
							"elements": []interface{}{
								map[string]interface{}{
									"text_run": map[string]interface{}{"content": line},
								},
							},
						},
					},
				},
				"index": idx,
			}); err != nil {
			return err
		}
	}
	return nil
}

func markdownLines(markdown string) []string {
	rawLines := strings.Split(strings.ReplaceAll(markdown, "\r\n", "\n"), "\n")
	lines := make([]string, 0, len(rawLines))
	for _, line := range rawLines {
		line = strings.TrimSpace(line)
		line = strings.TrimLeft(line, "#")
		line = strings.TrimSpace(strings.TrimPrefix(line, "- "))
		line = strings.TrimSpace(strings.TrimPrefix(line, "* "))
		if line == "" {
			continue
		}
		lines = append(lines, line)
	}
	return lines
}

func extractBlockText(block map[string]interface{}) string {
	for _, key := range []string{"page", "text", "heading1", "heading2", "heading3", "bullet", "ordered"} {
		if node, ok := block[key].(map[string]interface{}); ok {
			if text := extractElementsText(node["elements"]); text != "" {
				return text
			}
		}
	}
	return ""
}

func extractElementsText(raw interface{}) string {
	elements, _ := raw.([]interface{})
	parts := make([]string, 0, len(elements))
	for _, item := range elements {
		elem, _ := item.(map[string]interface{})
		if elem == nil {
			continue
		}
		if run, ok := elem["text_run"].(map[string]interface{}); ok {
			if content, ok := run["content"].(string); ok && content != "" {
				parts = append(parts, content)
			}
		}
	}
	return strings.TrimSpace(strings.Join(parts, ""))
}
