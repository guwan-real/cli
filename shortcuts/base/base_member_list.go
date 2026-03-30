// Copyright (c) 2026 Lark Technologies Pte. Ltd.
// SPDX-License-Identifier: MIT

package base

import (
	"context"
	"strings"

	"github.com/larksuite/cli/shortcuts/common"
)

var BaseMemberList = common.Shortcut{
	Service:     "base",
	Command:     "+member-list",
	Description: "List collaborators in a custom role",
	Risk:        "read",
	Scopes:      []string{"base:role:read"},
	AuthTypes:   []string{"user", "bot"},
	HasFormat:   true,
	Flags: []common.Flag{
		{Name: "base-token", Desc: "base token", Required: true},
		{Name: "role-id", Desc: "role ID", Required: true},
		{Name: "page-size", Type: "int", Default: "100", Desc: "page size per request (max 100)"},
	},
	Validate: func(ctx context.Context, runtime *common.RuntimeContext) error {
		if strings.TrimSpace(runtime.Str("base-token")) == "" {
			return common.FlagErrorf("--base-token must not be blank")
		}
		if strings.TrimSpace(runtime.Str("role-id")) == "" {
			return common.FlagErrorf("--role-id must not be blank")
		}
		return nil
	},
	DryRun: func(ctx context.Context, runtime *common.RuntimeContext) *common.DryRunAPI {
		return common.NewDryRunAPI().
			GET("/open-apis/bitable/v1/apps/:base_token/roles/:role_id/members").
			Params(map[string]interface{}{"page_size": runtime.Int("page-size")}).
			Set("base_token", runtime.Str("base-token")).
			Set("role_id", runtime.Str("role-id"))
	},
	Execute: func(ctx context.Context, runtime *common.RuntimeContext) error {
		items, total, err := listAllRoleMembers(runtime, runtime.Str("base-token"), runtime.Str("role-id"), runtime.Int("page-size"))
		if err != nil {
			return err
		}
		runtime.Out(map[string]interface{}{
			"items": simplifyRoleMembers(items),
			"total": total,
		}, nil)
		return nil
	},
}
