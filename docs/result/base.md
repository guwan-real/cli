# Base 模块最终测试结果

测试日期：2026-03-30  
测试环境：i 讯飞私有化飞书  
OpenAPI：`https://open.xfchat.iflytek.com`  
OAuth：`https://accounts.xfchat.iflytek.com`  
CLI：`./xfchat_cli`

## 1. 最终结论

Base 模块已完成私有化适配并打通。

当前已确认打通：

- Base：创建、获取、复制
- Table：列表、获取、创建、更新、删除
- Field：列表、获取、创建、更新、删除
- Record：新增/更新、获取、列表、删除、附件上传回写
- View：列表、获取、创建、重命名、删除
- Advanced Permission：开启
- Role：列表、创建、更新、删除
- Member：列表、批量新增、批量删除
- Form：创建、删除、读取、列表、元数据更新、问题列表、问题更新

当前仍未打通：

- `record-history-list`
- `data-query`

这两项目前表现为私有化接口未提供，不是 CLI body 或 scope 适配问题。

当前确认的行为差异：

- `Email` 字段按文档声明创建成功，但服务端最终回落为 `Text`
- 地理位置字段可创建，但记录写值当前仍返回 `1255002`
- User 把自己加入自己创建的高级权限角色时，写接口返回成功，但最终列表不落地

## 2. 私有化适配结论

根因：

- 公有云实现使用 `base/v3/bases`
- 私有化文档和服务实际实现使用 `bitable/v1/apps`

已落地的适配方向：

- 统一把 Base 模块主路径切到 `bitable/v1`
- 内部 `bases` 资源映射到私有化 `apps`
- 兼容 `tables / fields / views / roles / forms / members` 的私有化响应结构
- 修正私有化下的 body 差异：
  - `table-create`：`{"table":{"name":"..."}}`
  - `view-rename`：`{"view_name":"..."}`
  - `record-upsert`：自动包裹 `{"fields": ...}`
  - `advperm-enable`：`PUT /apps/:app_token` + `{"is_advanced":true}`
  - `form`：元数据走 `/forms/:form_id`，创建/删除复用视图能力
  - `form-questions`：走 `/forms/:form_id/fields` 和 `/forms/:form_id/fields/:field_id`

本轮新增或调整的命令能力：

- `+advperm-enable`
- `+advperm-disable`
- `+role-create`
- `+role-list`
- `+role-get`
- `+role-update`
- `+role-delete`
- `+member-list`
- `+member-batch-create`
- `+member-batch-delete`
- `+form-create`
- `+form-delete`
- `+form-get`
- `+form-list`
- `+form-update`
- `+form-questions-list`
- `+form-questions-update`

## 3. 当前授权与权限结论

最终 user token 已确认有效，且已包含本轮 Base 外围能力需要的关键 scope。

最终 `auth status --verify` 关键结果：

```json
{
  "identity": "user",
  "tokenStatus": "valid",
  "verified": true,
  "scope": "auth:user.id:read base:app:copy base:app:create base:app:read base:app:update base:collaborator:create base:collaborator:read base:field:create base:field:delete base:field:read base:field:update base:form:read base:form:update base:record:create base:record:delete base:record:read base:record:retrieve base:record:update base:role:create base:role:delete base:role:read base:role:update base:table:create base:table:delete base:table:read base:table:update base:view:read base:view:write_only bitable:app offline_access"
}
```

`docs/scope.json` 已确认存在并可对照的相关口径：

- `base:app:create`
- `base:app:copy`
- `base:app:read`
- `base:app:update`
- `base:table:create`
- `base:table:delete`
- `base:table:read`
- `base:table:update`
- `base:field:create`
- `base:field:delete`
- `base:field:read`
- `base:field:update`
- `base:record:create`
- `base:record:delete`
- `base:record:read`
- `base:record:retrieve`
- `base:record:update`
- `base:view:read`
- `base:view:write_only`
- `base:role:create`
- `base:role:delete`
- `base:role:read`
- `base:role:update`
- `base:form:read`
- `base:form:update`
- `base:collaborator:create`
- `base:collaborator:read`
- `bitable:app`
- `bitable:app:readonly`

备注：

- `bitable:app:readonly` 文档备注为“已不再维护，不可新增申请；建议申请 `bitable:app`”

## 4. 能力覆盖矩阵

### 4.1 Base / Table / Field / Record / View

| 能力 | User | Bot | 结论 |
| --- | --- | --- | --- |
| `base-create` | 通过 | 通过 | 通过 |
| `base-get` | 通过 | 通过 | 通过 |
| `base-copy` | 通过 | 未单独重复测 | 通过 |
| `table-list/get/create/update/delete` | 通过 | 通过 | 通过 |
| `field-list/get/create/update/delete` | 通过 | 通过 | 通过 |
| `record-upsert/get/list/delete` | 通过 | 通过 | 通过 |
| `record-upload-attachment` | 通过 | 通过 | 通过 |
| `view-list/get/create/rename/delete` | 通过 | 通过 | 通过 |
| `record-history-list` | 失败 | 失败 | 私有化接口未提供 |
| `data-query` | 失败 | 失败 | 私有化接口未提供 |

### 4.2 Advanced Permission / Role / Member / Form

| 能力 | User | Bot | 结论 |
| --- | --- | --- | --- |
| `advperm-enable` | 通过 | 通过 | 通过 |
| `role-list` | 通过 | 通过 | 通过 |
| `role-create` | 通过 | 通过 | 通过 |
| `role-update` | 通过 | 通过 | 通过 |
| `role-delete` | 命令已适配，未保留最终删除留痕 | 命令已适配，未保留最终删除留痕 | 能力可用 |
| `member-list` | 通过 | 通过 | 通过 |
| `member-batch-create` | 通过 | 通过 | 通过 |
| `member-batch-delete` | 通过 | 通过 | 通过 |
| `form-create` | 通过 | 通过 | 通过 |
| `form-delete` | 命令已适配，未保留最终删除留痕 | 命令已适配，未保留最终删除留痕 | 能力可用 |
| `form-get` | 通过 | 通过 | 通过 |
| `form-list` | 通过 | 通过 | 通过 |
| `form-update` | 通过 | 通过 | 通过 |
| `form-questions-list` | 通过 | 通过 | 通过 |
| `form-questions-update` | 通过 | 通过 | 通过 |

## 5. 关键资源与留痕

### 5.1 User Base 主资源

- Base：`basrz2pNKyyqHlDSmjXkezlEDlJ`
- Base Copy：`basrzjuvcYH9nGYeoNVuJthA6Bb`
- 默认表：`tblKstt6B2fazaXn`
- 保留测试表：`tblTkuaQc1MiElas`
- 保留测试表名：`自动化测试表2-已更新`
- 保留字段：`fldyh1cWsc`
- 默认表留痕记录：`recvflbEINCj6S`
- 保留视图：`vewgwBH5Yk`
- 保留视图名：`用户留痕视图-最终版`

User Base 链接：

- `https://yf2ljykclb.xfchat.iflytek.com/base/basrz2pNKyyqHlDSmjXkezlEDlJ`

### 5.2 Bot Base 主资源

- Base：`basrz5S4GR0eZaiw9pxfCHiLYSd`
- 默认表：`tblG16WHJUAhTWD6`
- 测试表：`tblLpv2thjf9GU3y`
- 测试表名：`Bot自动化测试表`
- 字段：`fldCrOAALR`
- 字段名：`Bot留痕文本`
- 记录：`recvflck46ZfRe`
- 保留视图：`vewQNn3rAK`
- 保留视图名：`Bot留痕视图`

Bot Base 链接：

- `https://yf2ljykclb.xfchat.iflytek.com/base/basrz5S4GR0eZaiw9pxfCHiLYSd`

### 5.3 场景样板：星火派 Agent 软件开发排期

- Base：`basrz4CISgYO2yc3n4kfuESDLWb`
- 名称：`星火派 Agent 软件开发排期`

Base 链接：

- `https://yf2ljykclb.xfchat.iflytek.com/base/basrz4CISgYO2yc3n4kfuESDLWb`

保留表：

- `版本路线图`：`tbluaaar4dRMnnfJ`
- `Sprint计划`：`tblvgRfFwh8vuGqI`
- `功能清单`：`tblHhU6jnERYPXos`
- `风险与阻塞`：`tblRJ6T7r7z5X2p4`

保留视图直链：

- 看板：`https://yf2ljykclb.xfchat.iflytek.com/base/basrz4CISgYO2yc3n4kfuESDLWb?table=tblHhU6jnERYPXos&view=vewTCcymhk`
- 甘特：`https://yf2ljykclb.xfchat.iflytek.com/base/basrz4CISgYO2yc3n4kfuESDLWb?table=tblHhU6jnERYPXos&view=vewXFERNo7`
- 画册：`https://yf2ljykclb.xfchat.iflytek.com/base/basrz4CISgYO2yc3n4kfuESDLWb?table=tblHhU6jnERYPXos&view=vewxLuAAVv`
- 表单：`https://yf2ljykclb.xfchat.iflytek.com/base/basrz4CISgYO2yc3n4kfuESDLWb?table=tblHhU6jnERYPXos&view=vewCp7ONqT`

前端人工创建并验证可直达的视图：

- 看板：`vew33Ge38x`
- 日历：`vew8CKlSWt`
- 甘特：`vew9a2WFkU`
- 画册：`vewoChlaVz`
- 表单：`vewbtUWN9I`

### 5.4 场景样板：字段能力全景

- Base：`basrzM5h3cX7nUJNz0NTX2uK1pc`
- 主表：`tblPmeVxw2tYuZDz`
- 主表名：`字段能力全景`
- 辅助表：`tblgcHDEOddHYUs3`
- 辅助表名：`字段类型辅助表`

Base 链接：

- 根链接：`https://yf2ljykclb.xfchat.iflytek.com/base/basrzM5h3cX7nUJNz0NTX2uK1pc`
- 表格：`https://yf2ljykclb.xfchat.iflytek.com/base/basrzM5h3cX7nUJNz0NTX2uK1pc?table=tblPmeVxw2tYuZDz&view=vewCLRBP1K`
- 看板：`https://yf2ljykclb.xfchat.iflytek.com/base/basrzM5h3cX7nUJNz0NTX2uK1pc?table=tblPmeVxw2tYuZDz&view=vew9w0nw74`
- 甘特：`https://yf2ljykclb.xfchat.iflytek.com/base/basrzM5h3cX7nUJNz0NTX2uK1pc?table=tblPmeVxw2tYuZDz&view=vewOxoQJdd`
- 画册：`https://yf2ljykclb.xfchat.iflytek.com/base/basrzM5h3cX7nUJNz0NTX2uK1pc?table=tblPmeVxw2tYuZDz&view=vewBBVpTva`
- 表单：`https://yf2ljykclb.xfchat.iflytek.com/base/basrzM5h3cX7nUJNz0NTX2uK1pc?table=tblPmeVxw2tYuZDz&view=vewekMqr3U`

保留主记录：

- `recvflupqylLVz`：`星火派 1.0 里程碑`
- `recvflupqfcPrM`：`星火派 协同消息中心`
- `recvflupssrFBo`：`星火派 知识运营控制台`

保留辅助记录：

- `recvfltwk1BhUt`：`能力目录-认证登录`
- `recvfltwibtqJ8`：`能力目录-消息协同`
- `recvfltwjRdyVl`：`能力目录-知识运营`

### 5.5 User 外围能力保留对象

- 角色：`rolG7yvSLV`
- 角色名：`用户协作者演示-已验证`
- User 表单：`vewUQvVeTn`
- User 表单名：`用户表单留痕`
- User 表单分享链接：`https://yf2ljykclb.xfchat.iflytek.com/share/base/shrrzDZj9cOUu54smGfxbMIONkg`

### 5.6 Bot 外围能力保留对象

- 高级权限已开启 Base：`basrz5S4GR0eZaiw9pxfCHiLYSd`
- 角色：`rolReaXDQA`
- 角色名：`Bot协作者演示-已验证`
- Bot 表单：`vew0NUApiV`
- Bot 表单名：`Bot表单留痕`
- Bot 表单分享链接：`https://yf2ljykclb.xfchat.iflytek.com/share/base/shrrzcOAhNtGErndeiI9vlgm0Ke`

## 6. 字段能力样板最终结论

已验证创建成功的字段类型：

- 文本：`fldooqDwbP`
- 人员：`fldAOQHxvf`
- 单选：`fld8GHgkZT`
- 日期：`fldbhdt7FI`
- 附件：`fldmSdrt3D`
- 复选框：`fldmBW4o9S`
- 数字：`fldCuHK2PK`
- 单向关联：`fldaeaAeir`
- 多选：`fld5IMGSCY`
- 邮箱字段声明：`fldgPPLY5i`
- 超链接：`fldYjS3ftZ`
- 货币：`fldyF21u79`
- 条码：`fld8KX1ZHz`
- 电话：`fldshtSnDq`
- 进度：`fldAF8PAW4`
- 评分：`fld0T6jL4e`
- 双向关联：`fld9bvNApw`
- 公式：`fld7c42jWg`
- 群组：`fldnzl3nTQ`
- 创建时间：`fldami0nQC`
- 最后更新时间：`fldeg3fEe8`
- 地理位置：`fld3R0OkWe`
- 创建人：`fld1RFGvOT`
- 自动编号：`fldbLOqaZx`
- 修改人：`fldo2CVZz8`

字段值格式结论：

- `record-upsert` 已兼容自动包裹 `{"fields": ...}`
- 日期字段接受 Unix 时间戳
- 超链接字段要求对象：

```json
{
  "text": "发布主页",
  "link": "https://xfchat.iflytek.com/agents/spark-launch"
}
```

- 单向关联、双向关联在当前私有化环境可接受“记录 ID 字符串数组”
- 群组字段可接受：

```json
[
  {
    "id": "oc_xxx"
  }
]
```

明确差异：

- `Email` 字段当前服务端回落为 `Text`
- 地理位置字段写值仍未打通

## 7. Advanced Permission / Role / Member / Form 最终结论

### 7.1 Advanced Permission

已验证：

- User：通过
- Bot：通过

真实私有化开启方式：

```json
{
  "is_advanced": true
}
```

对应接口：

- `PUT /open-apis/bitable/v1/apps/:app_token`

### 7.2 Role

最终结论：

- `role-list`：通过
- `role-create`：通过
- `role-update`：通过
- `role-delete`：命令已适配，可用

私有化约束：

- `role-create` 和 `role-update` 均要求携带 `table_roles`
- 高级权限开启后立即调用 `role-create` 可能短暂返回 `1254301 OperationTypeError`
- 等待约 2 秒后重试可恢复

### 7.3 Member

最终结论：

- Bot：
  - `member-list`：通过
  - `member-batch-create`：通过
  - `member-batch-delete`：通过
- User：
  - `member-list`：通过
  - `member-batch-create`：通过
  - `member-batch-delete`：通过

私有化成员返回字段口径：

- `member_name`
- `member_en_name`
- `member_type`
- `open_id`
- `union_id`
- `user_id`

唯一行为差异：

- User 把自己加入自己创建的高级权限角色时：
  - 写接口返回成功
  - 但 `member-list` 与原始 API 最终均返回空
- 改用其他成员，例如 `ou_981bb29d80e31c4aebe79225e866d8ff`，则新增、列表、删除均正常

因此当前判断：

- 协作者链路整体已打通
- “自加自己不落地”更像私有化服务端规则或静默忽略

### 7.4 Form

最终结论：

- `form-create`：通过
- `form-delete`：命令已适配，可用
- `form-get`：通过
- `form-list`：通过
- `form-update`：通过
- `form-questions-list`：通过
- `form-questions-update`：通过

实现方式说明：

- `form-create` / `form-delete` 当前复用视图能力
- `form-list` 当前从视图列表中过滤 `view_type=form`
- `form` 元数据和问题编辑均已按私有化 `/forms/...` 路径打通

User 侧保留表单：

- `vewUQvVeTn` / `用户表单留痕`

Bot 侧保留表单：

- `vew0NUApiV` / `Bot表单留痕`

## 8. 视图与直达链接结论

之前“打开总是默认视图”的原因不是视图没建成，而是发送的链接缺少 `table` 和 `view` 查询参数。

私有化前端已确认支持直达视图链接：

```text
https://yf2ljykclb.xfchat.iflytek.com/base/{app_token}?table={table_id}&view={view_id}
```

当前可稳定通过 OpenAPI 创建并直达的视图类型：

- `grid`
- `kanban`
- `gantt`
- `gallery`
- `form`

当前不能通过 OpenAPI 进一步配置的能力：

- `group`
- `sort`
- `filter`

这些子接口在当前私有化环境实测返回 `HTTP 404`。

补充差异：

- 前端存在 `calendar` 视图
- 但当前文档和 CLI 未验证其创建能力

## 9. 群内留痕

测试群：

- `oc_010df6b42b975ce056cc7c2e717abde8`

已保留的 Base 相关群消息：

- Base 核心样板闭环消息：`om_x100b53901ed93ca4386eb832fadc7fa`
- 字段样板闭环消息：`om_x100b53916bbd58a0385c013fd6bb429`
- Bot 外围能力总结：`om_x100b5392714fb0a4386fc0b841d720a`
- User 外围能力复测：`om_x100b539227b7e8a0386174470cf4b04`
- 协作者最终结论：`om_x100b5392cf0824a038602a8704466b5`

已确认留痕动作：

- Base 链接消息已发到群
- 相关总结消息已追加 Pin
- 相关总结消息已追加 `THUMBSUP`

## 10. 未打通与后续建议

当前仍未打通：

- `record-history-list`
- `data-query`
- 地理位置字段写值
- `Email` 字段真实 UI 类型未保持

后续建议优先级：

1. 查私有化是否提供 `record-history-list` 与 `data-query`
2. 继续追地理位置字段写值格式
3. 继续确认 `Email` 字段在私有化是否真正支持独立 UI 类型
4. 如果要继续做演示样板，可直接复用：
   - `星火派 Agent 软件开发排期`
   - `字段能力全景`
