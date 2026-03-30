// Copyright (c) 2026 Lark Technologies Pte. Ltd.
// SPDX-License-Identifier: MIT

package base

import (
	"context"
	"io"
	"strings"

	"github.com/larksuite/cli/internal/output"
	"github.com/larksuite/cli/shortcuts/common"
)

var BaseFormUpdate = common.Shortcut{
	Service:     "base",
	Command:     "+form-update",
	Description: "Update a form in a Base table",
	Risk:        "write",
	Scopes:      []string{"base:form:update"},
	AuthTypes:   []string{"user", "bot"},
	HasFormat:   true,
	Flags: []common.Flag{
		{Name: "base-token", Desc: "Base token (base_token)", Required: true},
		{Name: "table-id", Desc: "table ID", Required: true},
		{Name: "form-id", Desc: "form ID", Required: true},
		{Name: "name", Desc: "new form name"},
		{Name: "description", Desc: "new form description (plain text or markdown link like [text](https://example.com))"},
		{Name: "shared", Desc: "whether form is shared: true/false"},
		{Name: "shared-limit", Desc: "share scope: off|tenant_editable|anyone_editable"},
		{Name: "submit-limit-once", Desc: "limit submit once: true/false"},
	},
	DryRun: func(ctx context.Context, runtime *common.RuntimeContext) *common.DryRunAPI {
		return common.NewDryRunAPI().
			PATCH("/open-apis/bitable/v1/apps/:base_token/tables/:table_id/forms/:form_id").
			Set("base_token", runtime.Str("base-token")).
			Set("table_id", runtime.Str("table-id")).
			Set("form_id", runtime.Str("form-id"))
	},
	Execute: func(ctx context.Context, runtime *common.RuntimeContext) error {
		baseToken := runtime.Str("base-token")
		tableId := runtime.Str("table-id")
		formId := runtime.Str("form-id")
		name := runtime.Str("name")
		description := runtime.Str("description")
		shared := strings.TrimSpace(runtime.Str("shared"))
		sharedLimit := strings.TrimSpace(runtime.Str("shared-limit"))
		submitLimitOnce := strings.TrimSpace(runtime.Str("submit-limit-once"))

		body := map[string]interface{}{}
		if name != "" {
			body["name"] = name
		}
		if description != "" {
			body["description"] = description
		}
		if shared != "" {
			body["shared"] = strings.EqualFold(shared, "true")
		}
		if sharedLimit != "" {
			body["shared_limit"] = sharedLimit
		}
		if submitLimitOnce != "" {
			body["submit_limit_once"] = strings.EqualFold(submitLimitOnce, "true")
		}

		data, err := baseV3Call(runtime, "PATCH", baseFormPath(baseToken, tableId, formId), nil, body)
		if err != nil {
			return err
		}

		runtime.OutFormat(data, nil, func(w io.Writer) {
			output.PrintTable(w, []map[string]interface{}{
				{
					"id":                data["id"],
					"name":              data["name"],
					"description":       data["description"],
					"shared":            data["shared"],
					"shared_limit":      data["shared_limit"],
					"submit_limit_once": data["submit_limit_once"],
				},
			})
		})
		return nil
	},
}
