# 管理页面功能开发流程模板（后端 + RBAC 对应版）

这份文档用于指导你在 `go-vue-admin-server` 中新增一个后台管理功能，并且和前端页面开发步骤一一对应。

适用场景：

- 新增一个管理页面对应的后端接口
- 新增一个业务模块的增删改查
- 需要同时处理菜单、角色、按钮权限、Casbin 策略
- 需要让前后端联调时不漏步骤

对应前端文档：

- 前端文档位置：`go-vue-admin-ui/docs/admin-page-development-template.md`

建议使用方式：

1. 先看这份后端文档
2. 再对照前端文档一起推进
3. 每新增一个功能都复制“开发记录模板”使用

---

## 1. 先理解当前项目的权限模型

你这个项目当前已经具备完整的 RBAC 基础能力：

- 用户绑定角色
- 角色绑定菜单
- 菜单包含目录、页面菜单、按钮权限
- Casbin 根据角色菜单关系生成 API 访问策略
- 前端根据路由和按钮权限做页面展示控制

当前关键链路是：

1. 用户登录
2. 用户拿到角色
3. 角色关联菜单和按钮
4. 前端根据菜单生成路由
5. 后端根据角色菜单同步 Casbin 策略
6. 请求接口时由 `middleware.CasbinAuth()` 校验 API 权限

一句话理解：

**前端控制“看不看得到按钮”，后端控制“接口能不能真正调用成功”。**

---

## 2. 当前项目里和新增功能相关的关键目录

以后新增一个管理功能，通常会涉及这些目录：

```text
api/v1/            HTTP 接口层
services/v1/       业务逻辑层
models/            数据模型
router/v1/         路由注册
docs/              Swagger 文档输出
core/              Casbin 初始化逻辑
middleware/        JWT / Casbin / 日志 / 限流中间件
```

常规新增一个模块时，至少会动到：

```text
models/xxx.go
services/v1/xxx.go
api/v1/xxx.go
router/v1/system.go
```

如果新增功能要接入 RBAC，还要额外确认：

- 菜单数据是否要初始化
- 按钮权限标识是否要约定
- Casbin 映射规则是否要扩展

---

## 3. 后端新增一个管理功能的标准顺序

以后每新增一个管理功能，建议严格按下面顺序执行：

1. 明确功能范围
2. 明确数据库表和字段
3. 明确前端页面路由和权限前缀
4. 设计接口
5. 写 `model`
6. 写 `service`
7. 写 `api`
8. 注册 `router`
9. 接入 RBAC 菜单和按钮权限
10. 补 Casbin 映射
11. 本地接口自测
12. 和前端联调
13. 回归权限测试

不要跳过第 9 和第 10 步，否则页面能看到但接口不一定能调通。

---

## 4. 开发前先确认的内容

在开始写代码之前，先写清楚下面这些信息。

### 4.1 功能基础信息

- 功能名称：
- 所属模块：
- 路由地址：
- 前端页面组件：
- 权限前缀：
- 数据表名：

示例：

- 功能名称：设备管理
- 所属模块：系统管理
- 路由地址：`/system/device`
- 前端页面组件：`system/device/index`
- 权限前缀：`system:device`
- 数据表名：`system_device`

### 4.2 功能清单

明确这个页面需要哪些能力：

- 列表查询
- 分页
- 详情
- 新增
- 编辑
- 删除
- 批量删除
- 导出
- 启用/禁用

### 4.3 权限点清单

建议统一命名为：

`模块:资源:动作`

例如：

- `system:device:list`
- `system:device:add`
- `system:device:edit`
- `system:device:delete`
- `system:device:view`
- `system:device:export`

### 4.4 前后端对齐信息

这一步很重要，要和前端提前对齐：

- 前端路由 path 是什么
- 前端页面里会出现哪些按钮
- 前端接口地址打算怎么命名
- 前端查询字段、表单字段有哪些
- 后端返回结构是什么

---

## 5. 数据库设计步骤

### Step 1：确认是否需要新表

先判断：

- 是在老表上扩展
- 还是单独建一个新表

如果是标准管理页，建议字段尽量统一：

- `id`
- `name`
- `code`
- `status`
- `sort`
- `remark`
- `created_at`
- `updated_at`

如果是业务表，按业务需要扩展。

### Step 2：明确字段约束

至少确认：

- 哪些字段必填
- 哪些字段唯一
- 哪些字段可为空
- 哪些字段需要索引

比如：

- `device_code` 唯一
- `status` 只能是 `0/1`
- `device_name` 必填

### Step 3：确认和 RBAC 的关系

注意：

- 业务数据表本身通常不直接和角色表关联
- 权限控制由“菜单 + 角色菜单关系 + Casbin”完成

也就是说：

- 你的业务表存业务数据
- 菜单表 `system_menu` 存页面菜单和按钮权限
- 角色菜单表 `system_role_menu` 决定谁有权限

---

## 6. API 设计步骤

建议优先采用你项目当前的 REST 风格。

从现有路由看，项目更接近这种风格：

```text
GET    /api/v1/system/users
GET    /api/v1/system/users/:id
POST   /api/v1/system/users
PUT    /api/v1/system/users/:id
DELETE /api/v1/system/users/:id
```

新增模块时建议保持一致。

### 标准接口清单

以“设备管理”为例：

```text
GET    /api/v1/system/devices
GET    /api/v1/system/devices/:id
POST   /api/v1/system/devices
PUT    /api/v1/system/devices/:id
DELETE /api/v1/system/devices/:id
DELETE /api/v1/system/devices?ids=1,2,3   或单独 batch 接口
```

建议你先确定好下面这些内容：

- 列表接口参数名
- 分页参数名
- 返回结构
- 详情路径参数
- 批量删除方案

### 分页返回结构建议

建议和现有响应方式保持统一，例如：

```json
{
  "code": 0,
  "msg": "success",
  "data": {
    "list": [],
    "total": 0,
    "page": 1,
    "pageSize": 10
  }
}
```

这样前端最容易复用。

---

## 7. Model 开发步骤

### Step 1：创建模型文件

建议在 `models/` 下新增对应模型文件，例如：

```text
models/system_device.go
```

### Step 2：定义结构体

一个基础模型大致包含：

- 数据库字段
- JSON 标签
- GORM 标签

示例思路：

```go
type SystemDevice struct {
    ID         uint      `json:"id" gorm:"primaryKey"`
    DeviceName string    `json:"deviceName" gorm:"size:100;not null"`
    DeviceCode string    `json:"deviceCode" gorm:"size:100;uniqueIndex;not null"`
    Status     int       `json:"status" gorm:"default:1"`
    Remark     string    `json:"remark" gorm:"size:255"`
    CreatedAt  time.Time `json:"createdAt"`
    UpdatedAt  time.Time `json:"updatedAt"`
}
```

### Step 3：补请求 DTO 和查询 DTO

如果项目里 DTO 直接写在 API 层，也可以放 API 层；如果想更清晰，也可以在 model 层补结构。

建议至少有：

- 创建请求
- 更新请求
- 查询请求

---

## 8. Service 开发步骤

### Step 1：创建 service 文件

例如：

```text
services/v1/system_device.go
```

### Step 2：先写最基础的 5 个方法

标准 CRUD 最少包括：

- 列表查询
- 根据 ID 查详情
- 新增
- 更新
- 删除

可选：

- 批量删除
- 修改状态

### Step 3：列表查询要处理的内容

列表接口建议统一处理：

- 分页
- 关键字查询
- 状态筛选
- 排序

### Step 4：新增时要做的校验

新增前建议至少校验：

- 名称是否为空
- 编码是否重复
- 状态是否合法

### Step 5：更新时要做的校验

更新前建议至少校验：

- 记录是否存在
- 唯一字段是否冲突
- 是否允许修改某些系统字段

### Step 6：删除时要考虑关联关系

删除前先判断：

- 是否有子数据关联
- 是否允许物理删除
- 是否需要软删除

如果存在关联数据，尽量返回明确错误，不要直接删。

---

## 9. API 层开发步骤

### Step 1：创建 API 文件

例如：

```text
api/v1/system_device.go
```

### Step 2：每个接口做这几件事

API 层职责要轻：

1. 接收参数
2. 参数绑定
3. 参数校验
4. 调用 service
5. 返回统一响应

不要把复杂 SQL 或复杂业务规则堆进 API 层。

### Step 3：补 Swagger 注释

你这个项目已经在使用 Swagger，所以新增接口后建议同步补注释。

开发完成后执行：

```bash
make swagger
```

或者：

```bash
swag init -g main.go -o ./docs
```

---

## 10. Router 注册步骤

### Step 1：在 `router/v1/system.go` 注册接口

你当前系统管理模块已经统一挂在：

```text
/api/v1/system/*
```

并且受下面这些中间件保护：

- `JWTAuth()`
- `TokenBlacklistMiddleware()`
- `CasbinAuth()`
- `OperationLog()`
- `APIRateLimit()`

新增接口时，正常情况下应该注册到受保护的 `authRouter` 中。

### Step 2：接口风格保持一致

建议使用已有风格：

```go
authRouter.GET("/devices", systemDeviceApi.GetDeviceList)
authRouter.GET("/devices/:id", systemDeviceApi.GetDeviceDetail)
authRouter.POST("/devices", systemDeviceApi.CreateDevice)
authRouter.PUT("/devices/:id", systemDeviceApi.UpdateDevice)
authRouter.DELETE("/devices/:id", systemDeviceApi.DeleteDevice)
```

---

## 11. RBAC 接入步骤

这是最关键的部分。

新增一个管理页面，后端不只是加接口，还要让菜单、角色、按钮权限、Casbin 对上。

### 11.1 先理解当前菜单类型

从你项目现有逻辑看，菜单分三类：

- `1`：目录
- `2`：菜单
- `3`：按钮

它们的作用分别是：

- 目录：用于前端左侧导航分组
- 菜单：用于前端页面路由
- 按钮：用于前端按钮权限和后端 Casbin API 权限映射

### 11.2 新增一个页面通常需要补哪些菜单数据

以“设备管理”为例，通常至少需要这些菜单记录：

1. 目录或上级菜单
2. 页面菜单：设备管理
3. 按钮：新增
4. 按钮：编辑
5. 按钮：删除
6. 按钮：查看

示例概念：

- 菜单：`/system/device`
- 组件：`system/device/index`
- 权限前缀：`system:device`

按钮权限：

- `system:device:add`
- `system:device:edit`
- `system:device:delete`
- `system:device:view`

### 11.3 菜单和前端如何对应

前端关注这些字段：

- `path`
- `component`
- `menuName`
- `icon`
- `menuType`
- `perm`

对应关系建议这样理解：

- 页面菜单的 `path` 对应前端路由地址
- 页面菜单的 `component` 对应前端组件路径
- 按钮菜单的 `perm` 对应前端按钮权限字符串

### 11.4 角色如何获得权限

在你项目里，角色不是直接配 API 权限，而是：

1. 角色绑定菜单
2. 角色绑定按钮
3. `SetRoleMenus` 时同步更新 Casbin 策略

所以新增功能后，一定要做两件事：

1. 在菜单表中创建页面和按钮权限
2. 给测试角色或管理员角色分配这些菜单

---

## 12. Casbin 对接步骤

这一步非常关键。

你当前项目里的 Casbin 权限并不是完全自动推导所有业务模块，而是依赖映射规则。

现有逻辑主要在：

- `services/v1/system_role.go`
- `core/casbin.go`

### 12.1 当前 Casbin 的工作方式

大致流程是：

1. 根据角色拿到菜单 ID
2. 查出菜单信息
3. 如果是页面菜单，通过 `path` 映射 API
4. 如果是按钮权限，通过 `perm` 映射 API 和 HTTP 方法
5. 生成 Casbin policy

### 12.2 新增模块时你必须确认两类映射

#### 第一类：菜单路径映射

例如：

- 前端菜单 path：`/system/device`
- 后端 API 资源：`/v1/system/devices`

如果当前 `mapMenuPathToAPI` 没有这个映射，就需要补上。

#### 第二类：权限标识映射

例如：

- `system:device:add`
- `system:device:edit`
- `system:device:delete`

如果当前 `mapPermToAPI` 不认识 `system:device`，就需要补上：

```go
case "system:device":
    return "/v1/system/devices"
```

### 12.3 HTTP 方法映射也要确认

当前逻辑大致是：

- `:add` -> `POST`
- `:edit` -> `PUT`
- `:delete` -> `DELETE`
- `:export` -> `GET`
- 默认 -> `GET`

如果你新增了特殊动作，例如：

- `system:device:enable`
- `system:device:disable`
- `system:device:reset`

你要确认：

1. 它们走什么 HTTP 方法
2. `mapPermToMethod` 是否需要扩展

### 12.4 什么时候要重新同步 Casbin

当你做了这些操作后，要注意重新同步：

- 新增菜单
- 修改菜单 path
- 修改按钮 perm
- 给角色重新分配菜单

否则角色可能已经看得到页面，但接口调用还是 403。

---

## 13. 前后端对应关系表

以后新增一个页面，可以直接按这个表对应。

| 项目 | 前端 | 后端 |
|---|---|---|
| 页面路由 | `/system/device` | 菜单 `path=/system/device` |
| 页面组件 | `system/device/index` | 菜单 `component=system/device/index` |
| 查询列表 | `GET /system/devices` | `GET /api/v1/system/devices` |
| 新增按钮 | `system:device:add` | 按钮菜单 `perm=system:device:add` |
| 编辑按钮 | `system:device:edit` | 按钮菜单 `perm=system:device:edit` |
| 删除按钮 | `system:device:delete` | 按钮菜单 `perm=system:device:delete` |
| 按钮展示 | 前端 `hasPerm(...)` | 菜单权限 + 角色菜单关系 |
| 接口放行 | 无 | Casbin policy |

---

## 14. 一次完整新增“设备管理”功能的执行顺序

下面给你一个完整的实操顺序，以后可以照着走。

### 第一阶段：设计

1. 确认页面名称：设备管理
2. 确认前端路由：`/system/device`
3. 确认权限前缀：`system:device`
4. 确认数据表：`system_device`
5. 确认字段：`id/device_name/device_code/status/remark`
6. 确认接口：列表、详情、新增、编辑、删除

### 第二阶段：后端 CRUD

1. 新建 `models/system_device.go`
2. 新建 `services/v1/system_device.go`
3. 新建 `api/v1/system_device.go`
4. 在 `router/v1/system.go` 注册 `/devices` 路由
5. 补 Swagger 注释
6. 本地测通 CRUD 接口

### 第三阶段：RBAC

1. 在菜单表中新增“设备管理”页面菜单
2. 在菜单表中新增按钮权限：
   - `system:device:add`
   - `system:device:edit`
   - `system:device:delete`
   - `system:device:view`
3. 给管理员角色分配这些菜单
4. 给测试角色分配部分菜单
5. 在 `mapMenuPathToAPI` 中补 `/system/device -> /v1/system/devices`
6. 在 `mapPermToAPI` 中补 `system:device -> /v1/system/devices`
7. 如果有特殊动作，再补 `mapPermToMethod`
8. 重新分配一次角色菜单，触发 Casbin 同步

### 第四阶段：前端联调

1. 前端新建 `device` 页面
2. 对接列表接口
3. 对接新增编辑接口
4. 对接删除接口
5. 对接按钮权限
6. 用不同角色测试按钮展示和接口权限

### 第五阶段：回归测试

1. 管理员可以正常访问和操作
2. 只读角色能看到页面但不能新增删除
3. 无权限角色无法访问接口
4. 菜单、路由、按钮权限都生效

---

## 15. 联调时最容易出错的点

你后面可以重点检查这些地方：

1. 前端路由 path 和后端菜单 path 不一致
2. 前端按钮权限字符串和后端 `perm` 不一致
3. 新增了接口，但没在 `router` 注册
4. 注册了路由，但没进 `authRouter`
5. 菜单配好了，但 `mapPermToAPI` 没补
6. 菜单配好了，但角色没重新分配，Casbin 没刷新
7. 前端能看到按钮，但后端接口还是 403
8. 返回结构和前端预期不一致
9. `PUT /:id`、`DELETE /:id` 和按钮权限映射不一致

---

## 16. 提交前自测清单

### 16.1 后端接口

- [ ] 列表接口正常
- [ ] 详情接口正常
- [ ] 新增接口正常
- [ ] 编辑接口正常
- [ ] 删除接口正常
- [ ] 错误参数会返回明确错误

### 16.2 路由和文档

- [ ] 路由已注册
- [ ] 路由注册在正确分组下
- [ ] Swagger 注释已补
- [ ] Swagger 已重新生成

### 16.3 RBAC

- [ ] 菜单已创建
- [ ] 按钮权限已创建
- [ ] 角色已分配新菜单
- [ ] Casbin 映射已补
- [ ] 不同角色权限符合预期

### 16.4 联调

- [ ] 前端列表能查到数据
- [ ] 前端新增编辑删除都正常
- [ ] 前端按钮显示符合角色权限
- [ ] 无权限时后端接口正确拒绝

---

## 17. 开发记录模板

以后每次新增功能时，把下面这段复制出来填写。

```md
# 后端管理功能开发记录

## 1. 基础信息
- 功能名称：
- 所属模块：
- 前端路由：
- 前端组件：
- 权限前缀：
- 数据表名：

## 2. 功能范围
- [ ] 列表
- [ ] 分页
- [ ] 详情
- [ ] 新增
- [ ] 编辑
- [ ] 删除
- [ ] 批量删除
- [ ] 导出
- [ ] 启用/禁用

## 3. 数据库
- [ ] 表结构确认
- [ ] 字段确认
- [ ] 索引确认
- [ ] 唯一性确认

## 4. 代码开发
- [ ] 新建 model
- [ ] 新建 service
- [ ] 新建 api
- [ ] 注册 router
- [ ] 补 swagger

## 5. RBAC
- [ ] 新增页面菜单
- [ ] 新增按钮权限
- [ ] 角色分配菜单
- [ ] 补 mapMenuPathToAPI
- [ ] 补 mapPermToAPI
- [ ] 补 mapPermToMethod（如需要）
- [ ] 重新同步 Casbin

## 6. 自测
- [ ] CRUD 正常
- [ ] 权限正常
- [ ] 403 生效
- [ ] 联调通过

## 7. 备注
- 参考模块：
- 接口前缀：
- 特殊权限：
- 特殊说明：
```

---

## 18. 推荐原则

最后记住这几条，能帮你减少很多返工。

1. 先确定权限前缀，再写接口
2. 先完成 CRUD，再接 RBAC
3. 菜单、按钮、Casbin 映射必须一起做
4. 前端展示权限和后端接口权限必须同时验证
5. 每新增一个模块，都检查 `mapMenuPathToAPI` 和 `mapPermToAPI`
6. 每次联调都至少测试两种角色

如果你后面需要，我还可以继续帮你补一份更细的文档：

- “菜单表初始化 SQL 模板”
- “新增一个模块时的后端代码骨架模板”
- “RBAC 权限点命名规范模板”
