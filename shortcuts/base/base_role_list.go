// Copyright (c) 2026 Lark Technologies Pte. Ltd.
// SPDX-License-Identifier: MIT

package base

import (
	"context"
	"github.com/larksuite/cli/shortcuts/common"
	"strings"
)

var BaseRoleList = common.Shortcut{
	Service:     "base",
	Command:     "+role-list",
	Description: "List all roles in a Base",
	Risk:        "read",
	Scopes:      []string{"base:role:read"},
	AuthTypes:   []string{"user", "bot"},
	HasFormat:   true,
	Flags: []common.Flag{
		{Name: "base-token", Desc: "base token", Required: true},
	},
	Validate: func(ctx context.Context, runtime *common.RuntimeContext) error {
		if strings.TrimSpace(runtime.Str("base-token")) == "" {
			return common.FlagErrorf("--base-token must not be blank")
		}
		return nil
	},
	DryRun: func(ctx context.Context, runtime *common.RuntimeContext) *common.DryRunAPI {
		return common.NewDryRunAPI().
			GET("/open-apis/bitable/v1/apps/:base_token/roles").
			Set("base_token", runtime.Str("base-token"))
	},
	Execute: func(ctx context.Context, runtime *common.RuntimeContext) error {
		data, err := baseV3Call(runtime, "GET", baseRolePath(runtime.Str("base-token")), nil, nil)
		if err != nil {
			return err
		}
		items, _ := data["items"].([]interface{})
		if len(items) == 0 {
			items, _ = data["roles"].([]interface{})
		}
		if len(items) == 0 {
			if _, ok := data["role_id"]; ok {
				items = []interface{}{data}
			}
		}
		runtime.Out(map[string]interface{}{"items": items, "total": len(items)}, nil)
		return nil
	},
}
