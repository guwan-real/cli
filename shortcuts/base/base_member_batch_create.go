// Copyright (c) 2026 Lark Technologies Pte. Ltd.
// SPDX-License-Identifier: MIT

package base

import (
	"context"
	"strings"

	"github.com/larksuite/cli/shortcuts/common"
)

var BaseMemberBatchCreate = common.Shortcut{
	Service:     "base",
	Command:     "+member-batch-create",
	Description: "Batch add collaborators to a custom role",
	Risk:        "write",
	Scopes:      []string{"base:role:update"},
	AuthTypes:   []string{"user", "bot"},
	HasFormat:   true,
	Flags: []common.Flag{
		{Name: "base-token", Desc: "base token", Required: true},
		{Name: "role-id", Desc: "role ID", Required: true},
		{Name: "json", Desc: `member JSON array, e.g. '[{"type":"open_id","id":"ou_xxx"}]'`, Required: true},
	},
	Validate: func(ctx context.Context, runtime *common.RuntimeContext) error {
		if strings.TrimSpace(runtime.Str("base-token")) == "" {
			return common.FlagErrorf("--base-token must not be blank")
		}
		if strings.TrimSpace(runtime.Str("role-id")) == "" {
			return common.FlagErrorf("--role-id must not be blank")
		}
		if _, err := parseJSONArray(runtime.Str("json"), "json"); err != nil {
			return err
		}
		return nil
	},
	DryRun: func(ctx context.Context, runtime *common.RuntimeContext) *common.DryRunAPI {
		members, _ := parseJSONArray(runtime.Str("json"), "json")
		return common.NewDryRunAPI().
			POST("/open-apis/bitable/v1/apps/:base_token/roles/:role_id/members/batch_create").
			Body(map[string]interface{}{"member_list": members}).
			Set("base_token", runtime.Str("base-token")).
			Set("role_id", runtime.Str("role-id"))
	},
	Execute: func(ctx context.Context, runtime *common.RuntimeContext) error {
		members, err := parseJSONArray(runtime.Str("json"), "json")
		if err != nil {
			return err
		}
		data, err := baseV3Call(runtime, "POST", baseRolePath(runtime.Str("base-token"), runtime.Str("role-id"), "members", "batch_create"), nil, map[string]interface{}{"member_list": members})
		if err != nil {
			return err
		}
		runtime.Out(data, nil)
		return nil
	},
}
