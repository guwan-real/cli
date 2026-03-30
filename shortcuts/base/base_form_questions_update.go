// Copyright (c) 2026 Lark Technologies Pte. Ltd.
// SPDX-License-Identifier: MIT

package base

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"strings"

	"github.com/larksuite/cli/internal/output"
	"github.com/larksuite/cli/shortcuts/common"
)

var BaseFormQuestionsUpdate = common.Shortcut{
	Service:     "base",
	Command:     "+form-questions-update",
	Description: "Update questions of a form in a Base table",
	Risk:        "write",
	Scopes:      []string{"base:form:update"},
	AuthTypes:   []string{"user", "bot"},
	HasFormat:   true,
	Flags: []common.Flag{
		{Name: "base-token", Desc: "Base token (base_token)", Required: true},
		{Name: "table-id", Desc: "table ID", Required: true},
		{Name: "form-id", Desc: "form ID", Required: true},
		{Name: "questions", Desc: `questions JSON array, each item must include "id". Supported fields: "id"(required),"title","description","required","visible","pre_field_id". E.g. '[{"id":"fld_xxx","title":"Updated?","required":true}]'`, Required: true},
	},
	DryRun: func(ctx context.Context, runtime *common.RuntimeContext) *common.DryRunAPI {
		return common.NewDryRunAPI().
			PATCH("/open-apis/bitable/v1/apps/:base_token/tables/:table_id/forms/:form_id/fields/:field_id").
			Set("base_token", runtime.Str("base-token")).
			Set("table_id", runtime.Str("table-id")).
			Set("form_id", runtime.Str("form-id"))
	},
	Execute: func(ctx context.Context, runtime *common.RuntimeContext) error {
		baseToken := runtime.Str("base-token")
		tableId := runtime.Str("table-id")
		formId := runtime.Str("form-id")
		questionsJSON := runtime.Str("questions")

		var questions []interface{}
		if err := json.Unmarshal([]byte(questionsJSON), &questions); err != nil {
			return output.Errorf(output.ExitValidation, "invalid_json", "--questions must be a valid JSON array: %s", err)
		}

		items := make([]interface{}, 0, len(questions))
		for idx, question := range questions {
			m, _ := question.(map[string]interface{})
			fieldIDValue, _ := m["id"].(string)
			fieldIDValue = strings.TrimSpace(fieldIDValue)
			if fieldIDValue == "" {
				return output.Errorf(output.ExitValidation, "invalid_json", "--questions item %d must include non-empty string field \"id\"", idx+1)
			}
			delete(m, "id")
			data, err := baseV3Call(runtime, "PATCH", baseFormPath(baseToken, tableId, formId, "fields", fieldIDValue), nil, m)
			if err != nil {
				return err
			}
			items = append(items, data)
		}
		outData := map[string]interface{}{"questions": items}

		runtime.OutFormat(outData, nil, func(w io.Writer) {
			var rows []map[string]interface{}
			for _, item := range items {
				m, _ := item.(map[string]interface{})
				rows = append(rows, map[string]interface{}{
					"id":       fieldID(m),
					"title":    m["title"],
					"required": m["required"],
					"visible":  m["visible"],
				})
			}
			output.PrintTable(w, rows)
			fmt.Fprintf(w, "\n%d question(s) updated\n", len(items))
		})
		return nil
	},
}
