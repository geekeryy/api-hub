# API Hub

## 📦 项目说明

* 本项目基于go-zero微服务框架，旨在打造一个开箱即用的微服务开发脚手架。
* 通过框架与业务逻辑的解耦设计，让开发者能够专注于核心业务开发。只需简单配置，即可快速构建完整的微服务应用。

## 🏗️ 项目架构概览

![excalidraw.com](doc/assets/architecture.svg)

### 外部依赖层

* **go-zero** - 微服务框架，提供API服务基础能力
* **gin-gonic/gin** - Web框架，用于HTTP服务
* **gRPC** - 远程过程调用框架
* **go-i18n** - 国际化支持库
* **go-playground/validator** - 数据验证库

### 目录结构

```plain
.
├── api
│   └── gateway // 网关（BFF层）
│       ├── doc
│       └── internal
│           ├── config
│           ├── handler
│           ├── logic
│           │   ├── auth                    // 认证中心（授权域）
│           │   │   ├── admin B端授权
│           │   │   ├── jwks  jwt密钥开放接口
│           │   │   └── member C端授权
│           │   ├── healthz 健康检查
│           │   ├── oms                     // 运维管理服务（运维域）
│           │   │   └── jwks jwt密钥管理
│           │   └── user                    // 业务服务（用户域）
│           │       ├── admin B端用户管理
│           │       └── member C端用户管理
│           ├── middleware
│           ├── svc
│           └── types
├── core          核心库
│   ├── consts    常量
│   ├── email     邮件
│   ├── error     错误处理
│   ├── facebook  facebook
│   ├── google    google
│   ├── handler   响应包装
│   ├── jwks      jwks
│   ├── language  国际化
│   │   └── i18n  内置翻译字典
│   ├── limiter   限流器
│   ├── pgcache   pg缓存
│   ├── tracing   链路追踪
│   ├── validate  参数校验
│   ├── xcontext  上下文包装
│   ├── xgorm     sql日志
│   └── xstrings  字符串处理
├── deploy        部署
│   ├── local
│   └── prod
├── doc           文档
│   ├── assets    静态资源
│   ├── sql       sql
│   └── swagger   swagger
├── library       业务库
│   ├── consts    业务常量
│   ├── localization 国际化
│   │   └── i18n     国际化字典
│   ├── validator  自定义参数校验
│   └── xerror     自定义错误码
├── rpc
│   ├── model     gorm
│   └── user      用户服务（业务）
├── test          测试
└── tpl           项目模板
```

## 💡 设计特点

1. **分层架构** - 核心功能与业务逻辑分离
2. **统一错误处理** - 所有错误都通过核心错误处理器
3. **国际化支持** - 完整的多语言错误消息系统
4. **可扩展性** - 模块化设计便于功能扩展
5. **代码生成** - 基于模板的代码生成能力

## 🚀 使用方式

### 自定义参数校验

```go
// 实现自定义参数校验函数
// library/validator/example_validator.go
func ExampleValidator(v *validator.Validate, trans ut.Translator) error {
 if err := v.RegisterValidation("example_validator", chineseName); err != nil {
  return err
 }
 if err := v.RegisterTranslation("example_validator", trans, func(ut ut.Translator) error {
  ...
 }, func(ut ut.Translator, fe validator.FieldError) string {
  ...
 }
 return nil
}
```

```go
// 注册自定义参数校验函数
// api/gateway/svc/servicecontext.go
func NewServiceContext(c config.Config) *ServiceContext {
  ...
  validate.New([]validate.ValidatorFn{validator.ExampleValidator}, []string{"zh", "en"})
  ...
}
```

```go
// 在API定义中使用自定义参数校验函数
// api/gateway/gateway.api
type ExampleRequest {
  Name string `json:"name" comment:"FIELD_USERNAME" validate:"example_validator" `
}
```

### 自定义错误码

```go
// 定义错误码
// library/xerror/codes.go
var EmailFormatError = newerr(10047, "EMAIL_FORMAT_ERROR")
```

```toml
# 定义错误消息
# library/localization/error.zh.toml
EMAIL_FORMAT_ERROR =  "邮箱格式错误"

# library/localization/error.en.toml
EMAIL_FORMAT_ERRO = "Email format error"
```

## 功能列表

### 登录类型

* 账号+密码
* 邮箱+密码
* 手机号+密码
* 邮箱+验证码
* 手机号+验证码
* 第三方登录：wechat、google、facebook

## 🤔 QA

* 使用模板生成代码出现`<no value>`，如何解决？
  * 由于目前go-zero官方还未合并我的pr，暂时不支持projectPkg模板变量，可以选择使用我fork的goctl版本，或者等待官方合并。

* 为什么配置项的类型不能使用`int64`？
  * 因为go-zero使用环境变量覆盖配置项时，使用类型为int64的配置项会被误认为time.Duration，我已经提交pr修复，等待官方合并 [#4979](https://github.com/zeromicro/go-zero/pull/4979)

## 项目管理

### 开发流程

需求评审 -> 持续开发、测试用例评审 -> 开发自测 -> 产品走查、冒烟测试 -> 多轮测试 -> 预发布测试、开启下个迭代 -> 上线并测试、Checklist -> 项目复盘

### Commit规范

> 每次push需要压缩提交，commit描述清楚，让每次提交的功能尽可能清晰和单一；分支合并使用rebase，多次少量原则

* feat： 新增 feature
* fix: 修复 bug
* docs: 仅仅修改了文档，比如 README, CHANGELOG, CONTRIBUTE等等
* style: 仅仅修改了空格、格式缩进、逗号等等，不改变代码逻辑
* refactor: 代码重构，没有加新功能或者修复 bug
* perf: 优化相关，比如提升性能、体验
* test: 测试用例，包括单元测试、集成测试等
* chore: 改变构建流程、或者增加依赖库、工具等
* revert: 回滚到上一个版本

### Git分支管理

* master：稳定分支
* hotfix：紧急修复分支 线上紧急bug修复，修改后需合入dev
* release：发布分支 用于提测 格式：release/20250630（下个迭代开始后，release不再合并dev分支，而是直接在分支上修复问题，修改后需合入dev）
* dev：开发分支
* feature：功能分支

### 版本管理

[语义化版本 2.0.0](https://semver.org/lang/zh-CN/)

版本格式：主版本号.次版本号.修订号，版本号递增规则如下：

* 主版本号：当你做了不兼容的 API 修改
* 次版本号：当你做了向下兼容的功能性新增
* 修订号：当你做了向下兼容的问题修正

先行版本号及版本编译信息可以加到“主版本号.次版本号.修订号”的后面，作为延伸

定制版本：主版本号.次版本号.修订号-<客户> 意思是某客户基于某个版本的定制

## 代码规范

### 命名规范

[Effective Go](https://go.dev/doc/effective_go)

* go文件名：下划线
* 包名：小写的单个单词
* 变量名、函数名：驼峰
* 单方法接口以方法名加上-er后缀或类似的修饰来构造代理名词来命名：Reader

### 配置管理

* 原则：配置与程序隔离，使用环境变量覆盖（防止配置泄露）
* 开发环境：推荐使用`.env`文件，使用`dotenv`库加载

### 错误处理

* 业务代码在发生错误的地方打印错误日志并上报error
* 共享库warp详细信息并上报error，不打印日志

## 资源

* [百度翻译API](https://api.fanyi.baidu.com) 100w字符/月免费
* [谷歌翻译API](https://cloud.google.com/translate?hl=zh-cn) 50w字符/月免费

## 📝 TODO

* 框架：熔断、限流、降级、排队、监控、告警、日志、链路追踪
* 框架：不用鉴权的api利用签名机制防止盗刷
* 框架：CICD
* 测试：单元测试、集成测试、性能测试、暴力测试
* 功能：AI单轮聊天
* 功能：翻译池，聚合多家翻译API，通过账号池实现无限量翻译
