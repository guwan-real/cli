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

var BaseFormsList = common.Shortcut{
	Service:     "base",
	Command:     "+form-list",
	Description: "List all forms in a Base table (auto-paginated)",
	Risk:        "read",
	Scopes:      []string{"base:form:read"},
	AuthTypes:   []string{"user", "bot"},
	HasFormat:   true,
	Flags: []common.Flag{
		{Name: "base-token", Desc: "Base token (base_token)", Required: true},
		{Name: "table-id", Desc: "table ID", Required: true},
		{Name: "page-size", Type: "int", Default: "100", Desc: "page size per request (max 100)"},
	},
	DryRun: func(ctx context.Context, runtime *common.RuntimeContext) *common.DryRunAPI {
		return common.NewDryRunAPI().
			GET("/open-apis/bitable/v1/apps/:base_token/tables/:table_id/views").
			Set("base_token", runtime.Str("base-token")).
			Set("table_id", runtime.Str("table-id"))
	},
	Execute: func(ctx context.Context, runtime *common.RuntimeContext) error {
		views, _, err := listAllViews(runtime, runtime.Str("base-token"), runtime.Str("table-id"), 0, runtime.Int("page-size"))
		if err != nil {
			return err
		}
		allForms := make([]interface{}, 0)
		for _, view := range views {
			if viewType(view) == "form" {
				allForms = append(allForms, map[string]interface{}{
					"id":          viewID(view),
					"name":        viewName(view),
					"description": view["description"],
					"type":        viewType(view),
				})
			}
		}

		outData := map[string]interface{}{
			"forms": allForms,
			"total": len(allForms),
		}
		runtime.OutFormat(outData, nil, func(w io.Writer) {
			if len(allForms) == 0 {
				fmt.Fprintln(w, "No forms found.")
				return
			}
			var rows []map[string]interface{}
			for _, item := range allForms {
				m, _ := item.(map[string]interface{})
				rows = append(rows, map[string]interface{}{
					"id":          m["id"],
					"name":        m["name"],
					"description": m["description"],
					"type":        m["type"],
				})
			}
			output.PrintTable(w, rows)
			fmt.Fprintf(w, "\n%d form(s) total\n", len(allForms))
		})
		return nil
	},
}
