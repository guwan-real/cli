// Copyright (c) 2026 Lark Technologies Pte. Ltd.
// SPDX-License-Identifier: MIT

package base

import (
	"context"
	"github.com/larksuite/cli/shortcuts/common"
	"strings"
)

var BaseAdvpermEnable = common.Shortcut{
	Service:     "base",
	Command:     "+advperm-enable",
	Description: "Enable advanced permissions for a Base",
	Risk:        "write",
	Scopes:      []string{"base:app:update"},
	AuthTypes:   []string{"user", "bot"},
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
			PUT("/open-apis/bitable/v1/apps/:base_token").
			Body(map[string]interface{}{"is_advanced": true}).
			Set("base_token", runtime.Str("base-token"))
	},
	Execute: func(ctx context.Context, runtime *common.RuntimeContext) error {
		data, err := baseV3Call(runtime, "PUT", baseV3Path("bases", runtime.Str("base-token")), nil, map[string]interface{}{"is_advanced": true})
		if err != nil {
			return err
		}
		runtime.Out(data, nil)
		return nil
	},
}
