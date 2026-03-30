// Copyright (c) 2026 Lark Technologies Pte. Ltd.
// SPDX-License-Identifier: MIT

package base

import (
	"context"
	"io"

	"github.com/larksuite/cli/internal/output"
	"github.com/larksuite/cli/shortcuts/common"
)

var BaseFormCreate = common.Shortcut{
	Service:     "base",
	Command:     "+form-create",
	Description: "Create a form in a Base table",
	Risk:        "write",
	Scopes:      []string{"base:view:write_only"},
	AuthTypes:   []string{"user", "bot"},
	HasFormat:   true,
	Flags: []common.Flag{
		{Name: "base-token", Desc: "Base token (base_token)", Required: true},
		{Name: "table-id", Desc: "table ID", Required: true},
		{Name: "name", Desc: "form name", Required: true},
		{Name: "description", Desc: `form description (plain text or markdown link like [text](https://example.com))`},
	},
	DryRun: func(ctx context.Context, runtime *common.RuntimeContext) *common.DryRunAPI {
		return common.NewDryRunAPI().
			POST("/open-apis/bitable/v1/apps/:base_token/tables/:table_id/views").
			Body(map[string]interface{}{"view_name": runtime.Str("name"), "view_type": "form"}).
			Set("base_token", runtime.Str("base-token")).
			Set("table_id", runtime.Str("table-id"))
	},
	Execute: func(ctx context.Context, runtime *common.RuntimeContext) error {
		baseToken := runtime.Str("base-token")
		tableId := runtime.Str("table-id")
		name := runtime.Str("name")
		description := runtime.Str("description")

		data, err := baseV3Call(runtime, "POST",
			baseV3Path("bases", baseToken, "tables", tableId, "views"), nil,
			map[string]interface{}{"view_name": name, "view_type": "form"})
		if err != nil {
			return err
		}
		formID := viewID(data)
		if description != "" && formID != "" {
			updated, updateErr := baseV3Call(runtime, "PATCH", baseFormPath(baseToken, tableId, formID), nil, map[string]interface{}{"description": description})
			if updateErr == nil {
				data = updated
			} else {
				data["description_update_error"] = updateErr.Error()
			}
		}

		runtime.OutFormat(data, nil, func(w io.Writer) {
			output.PrintTable(w, []map[string]interface{}{
				{
					"id":          viewID(data),
					"name":        viewName(data),
					"description": data["description"],
				},
			})
		})
		return nil
	},
}
