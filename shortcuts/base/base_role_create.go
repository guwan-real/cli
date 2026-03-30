// Copyright (c) 2026 Lark Technologies Pte. Ltd.
// SPDX-License-Identifier: MIT

package base

import (
	"context"
	"encoding/json"
	"github.com/larksuite/cli/shortcuts/common"
	"strings"
)

var BaseRoleCreate = common.Shortcut{
	Service:     "base",
	Command:     "+role-create",
	Description: "Create a custom role in a Base",
	Risk:        "write",
	Scopes:      []string{"base:role:create"},
	AuthTypes:   []string{"user", "bot"},
	Flags: []common.Flag{
		{Name: "base-token", Desc: "base token", Required: true},
		{Name: "json", Desc: `body JSON (AdvPermBaseRoleConfig), e.g. {"role_name":"Reviewer","role_type":"custom_role","table_rule_map":{...}}`, Required: true},
	},
	Validate: func(ctx context.Context, runtime *common.RuntimeContext) error {
		if strings.TrimSpace(runtime.Str("base-token")) == "" {
			return common.FlagErrorf("--base-token must not be blank")
		}
		var body map[string]any
		if err := json.Unmarshal([]byte(runtime.Str("json")), &body); err != nil {
			return common.FlagErrorf("--json must be valid JSON: %v", err)
		}
		return nil
	},
	DryRun: func(ctx context.Context, runtime *common.RuntimeContext) *common.DryRunAPI {
		var body map[string]any
		json.Unmarshal([]byte(runtime.Str("json")), &body)
		return common.NewDryRunAPI().
			POST("/open-apis/bitable/v1/apps/:base_token/roles").
			Body(body).
			Set("base_token", runtime.Str("base-token"))
	},
	Execute: func(ctx context.Context, runtime *common.RuntimeContext) error {
		var body map[string]any
		json.Unmarshal([]byte(runtime.Str("json")), &body)

		data, err := baseV3Call(runtime, "POST", baseRolePath(runtime.Str("base-token")), nil, body)
		if err != nil {
			return err
		}
		runtime.Out(data, nil)
		return nil
	},
}
