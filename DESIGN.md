# Asuka DS Design System Design

## Intent

为 Asuka 做一个类 shadcn/ui 的组件库：不是 npm 包式的黑盒控件，而是应用可以复制、拥有和修改的一组 Go 源码；同时保留 Asuka 的核心价值——用平台原生控件获得输入法、无障碍、焦点、滚动和手感。

成功标准是：一套 token 与组件 API 能覆盖常见产品界面；品牌皮肤在六个平台上有一致的视觉结构；关键表单控件保留平台原生行为；组件代码在用户工程里可读、可改。

## Goals

- 提供 shadcn 风格的语义 token、类型化主题、组合式组件和零值可用的 props API。
- 提供自绘与原生两套按钮后端；输入框、开关、滑块、选择器、弹层优先复用平台能力。
- 主题切换通过 signal 驱动，不重建页面。
- 组件的所有交互语义都可在 `uitest` 中验证，包括 disabled、selected、expanded、field label/error。
- 文档与 showcase 直接由 Asuka 构建，展示 Web 与 macOS 两种宿主。

## Non-goals

- 不追求把 shadcn 的 CSS 变量、Tailwind、Radix、className 覆盖机制逐字移植到 Go。
- 不做运行时主题包或动态库插件。
- 不把原生控件完全替换成自绘控件；不承诺 iOS/Android 的像素级一致。
- 不在 v0.1 覆盖虚拟列表、图表、复杂数据网格或完整设计令牌编辑器。

## Anti-goals

- 组件看起来“统一”，但破坏原生输入、读屏状态或键盘导航。
- 用户改组件必须 fork 整个框架，或者必须通过全局 monkey patch 才能覆盖样式。
- 组件内部散落十六进制颜色、硬编码字号、圆角或动画时长。

## Constraints

- Asuka `v0.1.1`，Go `1.26+`；组件接口遵循 `ui.Node` 与 `signal.Accessor`。
- 组件只跑一次；所有会变的值必须放在 accessor、`ui.StyleFunc` 或 `ui.Dynamic` 内读取。
- 原生按钮只接受 `Font`/`Color`；`TextField{Plain:true}` 用于自定义输入框外壳。
- `View` 没有原生 `Disabled`、`Selected`、`Expanded` 语义；`Pressable` 用动态 `Hint` 与事件门控表达这些状态，并在静态禁用时不挂载交互事件。
- `ui.Sheet` 的 children 不参与布局；`Overlay` 的内容必须通过 `Content` 构建，children 只用于 `Popover` 锚点。
- 主题是 app 级 signal；同一进程内不支持多个相互独立的设计系统主题。

## Approaches

### A. Native-first

只提供 token、排版和表面组件；交互控件全部包装原生控件。

- 优点：平台行为、无障碍和升级稳定性最好。
- 缺点：品牌表达弱，和 shadcn 的视觉目标差距大。

### B. Skin-first

所有控件尽量自绘，只保留输入框等平台能力。

- 优点：跨平台视觉最一致，最接近 shadcn 的“脸”。
- 缺点：要自行补齐 disabled/selected/expanded、焦点、动效、平台手感；合规风险高。

### C. Hybrid（采用）

一套 token、一套 props API；容器和品牌控件自绘，交互与系统能力优先原生；Button 同时给出 `Button` 和 `NativeButton`。

- 原因：这是 Asuka 的价值观与 shadcn 的“代码归你所有”之间可长期维护的交集。
- 代价：文档必须明确哪些组件是品牌皮肤、哪些组件是平台原生包装，不能承诺所有平台像素一致。

## Architecture

- `ds/theme.go`：`Palette`、`Radii`、`TypeScale`、`Motion`、`Elevation`、`Theme.Validate()`、`SetTheme()`、`Themed()`。
- `pressable.go`：所有自绘交互的统一语义入口；负责 disabled/selected/expanded hint、事件门控与禁用态透明度。
- `field.go`：label、description、error 的通用表单语义；错误与描述同时显示，且都进入 accessible name。
- 组件按职责分文件：`button.go`、`card.go`、`input.go`、`controls.go`、`overlay.go`、`navigation.go`、`display.go`、`command.go`。
- `examples/showcase/`：全组件、双主题、双外观的 showcase；`examples/form/`：表单专项。
- `docs/research.md`：一手调研、硬约束、里程碑与风险；本文件：产品边界与已定方案。

## API Decisions

- `XxxProps` 结构体；零值给出合理默认，变体使用枚举。
- 静态文本用 `string`，动态文本用 `signal.Accessor[string]`；会变的布尔语义用 `signal.Accessor[bool]`。
- `Overlay` 支持 `Content`；未提供 `Content` 时，在 `Popover` 模式下允许 children 同时作为锚点和内容，但 `Sheet` 不接收 children。
- `Input` 复用 `Field`，因此 label/description/error、disabled、焦点环与 a11y 名称保持一致。
- `Select` 在 disabled 时不挂 `ui.Menu`，避免读屏/键盘仍能打开菜单。
- `Command` 在过滤结果变少时把选中项夹回有效范围，避免选中态消失。
- `Tabs` 使用 `Pressable`，让选中态进入动态 `Hint`。

## Distribution

- 主线是 Go module：用户 `go get` 后可直接使用；需要深度定制时复制或 vendor 单个 `.go` 文件。
- registry 与 `asuka-ds add` 留到 M4：先把 API 和组件质量冻结，再引入源码拷贝工具。
- 文档站复用 Asuka web 构建产物，避免为了组件库再引入第二套前端框架。

## Milestone Acceptance

### M1 Infrastructure

- `Theme`/`Palette`/`TypeScale`/`Radii`/`Motion`/`Elevation` 已冻结为强类型 API。
- `Native()` 使用平台语义色，`Neutral()` 自带 shadcn zinc 明暗两套；二者都通过 `Validate()`。
- `Validate()` 检查缺失 token 与 WCAG AA 文字对比度；`MotionFor()` 支持 Reduce Motion。
- `SetTheme()` 通过 signal 实时换肤；`Themed()` 提供逃生口。

### M2 Primitives and Accessibility

- `Pressable` 统一 disabled/selected/expanded 语义；静态禁用不注册 press/long-press，动态禁用阻止事件并更新 hint。
- `Field` 统一 label、description、error、required 与 accessible name；`Input` 复用它。
- `Overlay` 包装 `Sheet`/`Popover`；`Menu` 包装平台菜单并支持 disabled/checked/role；`Toast` 包装 Snackbar。
- 29 个 `uitest` 覆盖按钮、字段、Tabs、Radio、Accordion、Select、Table、Avatar、Progress、Dialog、Toast、主题对比度与 motion。

### M3 Component Coverage

- Surface：`Card`、`Badge`、`Separator`、`Skeleton`、`Alert`、`Avatar`、`Progress`、`Table`。
- Form：`Input`、`Switch`/Checkbox、`RadioGroup`、`Select`、`Slider`、`Field`、`Command`。
- Navigation：`Tabs`、`Accordion`。
- Overlay/feedback：`Dialog`、`ConfirmDialog`、`Overlay`、`Menu`、`Toast`、`Tooltip`。
- 两个 demo 在 Web 与 macOS 构建中展示组件与主题切换。

## Remaining Work

- M4：registry 索引、`asuka-ds add`、文档站、版本兼容策略。
- M5：`uitest` 截图快照、六平台基线、CI 对比度检查、iOS/Android 真机走查。
- 上游建议：为 `ViewProps` 增加真正的 `Disabled`、`Selected`、`Expanded` trait，避免继续用 `Hint` 兜底。

## Builder's Calls

- v0.1 以 Go module 为主，不先做 CLI registry。
- 组件文件保持 root package `ds`，不提前拆成多个 Go package。
- `Overlay` 的 `Sheet` 模式要求调用方提供 `Content`；children 仅在 `Popover` 模式下作为锚点。
- 错误信息与描述同时显示，因为表单用户需要同时知道“哪里错了”和“这个字段是什么”。
- 自绘控件的状态语义用 `Hint` 表达，是 Asuka `v0.1.1` 限制下的短期方案；M2 文档明确其边界。
