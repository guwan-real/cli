// Copyright (c) 2026 Lark Technologies Pte. Ltd.
// SPDX-License-Identifier: MIT

package base

import (
	"context"
	"fmt"
	"io"

	"github.com/larksuite/cli/internal/output"
	"github.com/larksuite/cli/shortcuts/common"
)

var BaseFormQuestionsList = common.Shortcut{
	Service:     "base",
	Command:     "+form-questions-list",
	Description: "List questions of a form in a Base table",
	Risk:        "read",
	Scopes:      []string{"base:form:read"},
	AuthTypes:   []string{"user", "bot"},
	HasFormat:   true,
	Flags: []common.Flag{
		{Name: "base-token", Desc: "Base app token (base_token)", Required: true},
		{Name: "table-id", Desc: "table ID", Required: true},
		{Name: "form-id", Desc: "form ID", Required: true},
	},
	DryRun: func(ctx context.Context, runtime *common.RuntimeContext) *common.DryRunAPI {
		return common.NewDryRunAPI().
			GET("/open-apis/bitable/v1/apps/:base_token/tables/:table_id/forms/:form_id/fields").
			Set("base_token", runtime.Str("base-token")).
			Set("table_id", runtime.Str("table-id")).
			Set("form_id", runtime.Str("form-id"))
	},
	Execute: func(ctx context.Context, runtime *common.RuntimeContext) error {
		baseToken := runtime.Str("base-token")
		tableId := runtime.Str("table-id")
		formId := runtime.Str("form-id")

		data, err := baseV3Call(runtime, "GET", baseFormPath(baseToken, tableId, formId, "fields"), nil, nil)
		if err != nil {
			return err
		}

		items, _ := data["items"].([]interface{})
		if len(items) == 0 {
			items, _ = data["fields"].([]interface{})
		}
		total := toInt(data["total"])
		if total == 0 {
			total = len(items)
		}
		outData := map[string]interface{}{
			"questions": items,
			"total":     total,
		}

		runtime.OutFormat(outData, nil, func(w io.Writer) {
			if len(items) == 0 {
				fmt.Fprintln(w, "No questions found.")
				return
			}
			var rows []map[string]interface{}
			for _, item := range items {
				m, _ := item.(map[string]interface{})
				rows = append(rows, map[string]interface{}{
					"id":          fieldID(m),
					"title":       m["title"],
					"description": m["description"],
					"required":    m["required"],
					"visible":     m["visible"],
				})
			}
			output.PrintTable(w, rows)
			fmt.Fprintf(w, "\n%d question(s) total\n", len(items))
		})
		return nil
	},
}
