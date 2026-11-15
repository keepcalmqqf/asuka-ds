# 基于 Asuka 做类 shadcn/ui 的 Design System：调研与方案

调研对象：<https://asuka.egoist.dev/zh>（`github.com/egoist/asuka`，v0.1.1，MIT，Go 1.26+）

结论先行：**可行，但不能照搬 shadcn 的"纯自绘"路线**。Asuka 的核心卖点是"每个平台用原生控件"，而 shadcn 的核心卖点是"像素完全可控"。两者在按钮、开关、分段控件这类小控件上直接冲突，必须做取舍。建议做一个**双后端的设计系统**：一套 token + 一套 API，组件可以选"原生皮肤"或"自绘皮肤"，并且默认跟随平台。

本文所有结论都在本仓库的原型上跑过：Web（Chromium，明/暗/中性/hover/键盘焦点/390pt）、macOS（AppKit 原生窗口），`asuka build`、`asuka vet` 均通过。M1→M2→M3 的实现与验收记录见 `rfcs/2026-11-14-asuka-ds-m1-m3.md`。

> 本文是调研当时的记录，文中出现的 `demo/`、`demo-form/`、`site/`、`docs/screenshots/`、`docs/superpowers/` 均为旧路径；当前布局分别是 `examples/showcase/`、`examples/form/`、`website/`、`artwork/screenshots/`、`DESIGN.md` 与 `rfcs/`。

---

## 1. Asuka 是什么

用 Go 写应用，界面在 iOS、Android、macOS、Windows、Linux、Web 上分别用**平台自己的控件**渲染：UIKit、Android Views + Material 3、AppKit、Win32、GTK 4 + libadwaita、DOM。

关键设计（决定了组件库能做什么）：

| 层 | 实现 | 对组件库的含义 |
|---|---|---|
| 语言 | 纯 Go，app 就是一个返回 `ui.Node` 的函数 | 组件是普通 Go 函数，可直接复制、改、内联 |
| 更新 | 细粒度 signal（Solid 式），组件只跑一次，只重发变化的属性 | 组件不需要"重渲染"概念，动态值传 `signal.Accessor` |
| 布局 | 自己实现的 flexbox/grid，单位 pt，跨平台一致 | 尺寸、间距、圆角可以做成真正的 token |
| 渲染 | 宿主创建/更新平台控件；Web 是 WASM + DOM | 原生控件的行为（IME、选择、无障碍、滚动）免费获得 |
| 样式 | `ui.Style`：背景、圆角、边框、阴影、渐变、材质、变换、裁剪、字体、颜色；`StyleFunc` 响应 signal | 容器层可完全换肤；**控件本身只在平台支持的地方生效** |
| 主题 | 语义色（`colors.Label` 等）+ 固定色；`asuka.json` 的 `accent`；`ui.SetColorScheme` | 深色模式/对比度/Android 动态色免费，但品牌色需要自带明暗两套 |

平台支持：iOS 16、Android 13 (API 33)、macOS 11、Windows 10、Linux（GTK 4.14 + libadwaita 1.5）、主流浏览器。桌面宿主不需要 cgo。

---

## 2. shadcn/ui 的四个支柱，在 Asuka 里各自对应什么

| shadcn/ui | Asuka 的对应物 | 可行性 |
|---|---|---|
| **代码归你所有**：`shadcn add button` 把源码拷进项目 | Go 包 + 可复制的 `.go` 文件；可自建 registry + CLI | ✅ 天然适合，Go 源码比 TSX 更好复制（无构建期类名） |
| **tokens 即 CSS 变量**：`--background`、`--primary`… | `Theme` 结构体 + signal；没有 CSS 变量层 | ✅ 可自建，且能做**运行时切换主题** |
| **Radix 提供无障碍原语** | 平台原生控件 + 可按压 `View`（`Role`/`Label`/`Hint`/`Actions`/`OnAdjust`） | ⚠️ 部分：原生控件满分；自绘控件有真实缺口（见 §4） |
| **CVA 变体 + `className` 覆盖** | 变体枚举 + props 结构体 + `Style` 组合 / `ds.Themed` 逃生口 | ✅ 可行，API 更"Go" |

---

## 3. 硬约束：哪些能自绘，哪些必须原生

### 3.1 原生控件不可深度换肤（实测源码 + 文档）

- `ui.Button` 的 `Style` 只认 `Font` 和 `Color`（tint）；圆角、背景、边框由平台决定（iOS 26 的 Liquid Glass、Android 的 Material 按钮）。
- `ui.TextField` 有一个关键开关 **`Plain: true`**：去掉原生边框、背景和内边距，只留文字。这让"DS 画外壳 + 原生输入框在内"成为可能（原型已验证）。
- `ui.Toggle`、`ui.Slider`、`ui.Segmented`、`ui.DatePicker`、`ui.Wheel` 基本只能给 tint / font，不能改形状。

### 3.2 用 `ui.View` 自绘时，能力是够的

- 视觉：背景、圆角、边框（内描边）、阴影、线性/径向渐变、`Material` 毛玻璃、`Clip`、`Transform`，以及 `ui.Shape` / `ui.Canvas` / `ui.Shader`。
- 交互：`OnPress`、`OnLongPress`、`OnPan`/`OnPinch`/`OnRotate`/`OnSwipe`、`Keys`（键盘）、`Drag`/`Drop`。
- 状态样式：`Style.Hovered` / `Pressed` / `Dragging` / `DragOver`（CSS 的 `:hover`/`:active`），`Style.When` / `Unless` 跟随任意 signal。
- 无障碍：带 `OnPress` 的 View 在六个平台都会被当成按钮；`Label`/`Hint`/`Role`/`Actions`/`OnAdjust` 可用。
- 键盘：Windows / Linux / Web 上可按压 View 是 Tab 停靠点（原型在 Web 上实测到 `role=button`、`tabindex=0` 和系统焦点环）。

### 3.3 自绘的真实缺口（做组件库必须正面处理）

1. **没有 `Disabled`**：`ViewProps` 里没有这个字段。禁用只能"不挂 `OnPress` + 降透明度 + 换 a11y 文案"，屏幕阅读器不会说"已禁用"。
2. **没有 selected / checked / expanded 状态语义**：Tabs、Accordion、Switch 无法让读屏软件读出"当前选中"。原生 `Segmented` 有，自绘没有。
3. **`Label` 只在创建时设置一次**：会变的标签必须重建节点（`Show` / `Dynamic`），不能原地更新。
4. **hover/pressed 变化没有过渡**：`Style.Transition` 只在 `ui.Animate` 内用于插入/移除，指针引起的状态变化是硬切。要动效只能自己做（例如 `ui.Animate` 包住业务状态，而不是 hover）。
5. **自绘控件不参与平台的"高级手感"**：Android ripple、Windows Fluent hover 填充、iOS 26 Liquid Glass、iPad 指针效果都拿不到（设置 `Style.Pressed` 后平台高亮会被主动关掉）。
6. **没有应用级 context**：`ui.UseEnvironment` 只给环境（深色、字号、平台…），不能注入自定义的 Theme。多主题并存只能靠全局切换。

> 这些不是"实现不好"，而是框架的价值观：它优先保证"像这个平台的应用"。做 shadcn 式设计系统，本质是在它的价值观上开一个"品牌皮肤"的口子。

---

## 4. 三种产品形态（建议选 C）

### A. Native-first（原生优先）
只做 token、排版、`Card`/`Badge`/`Separator`/`Skeleton` 等表面组件，交互控件全部包装原生。

- 优点：无障碍、IME、焦点、平台手感全部免费；升级 Asuka 最稳。
- 缺点：不是 shadcn 的脸；品牌表达力弱。
- 适合：面向 iOS/Android 的产品团队。

### B. Skin-first（shadcn 克隆）
Button / Tabs / Accordion / Dialog 触发器等全部 `View` 自绘，只有输入框、开关、选择器保留原生。跨平台几乎像素一致。

- 优点：真正的"品牌设计系统"，一套视觉走六个平台。
- 缺点：§3.3 的坑要自己填；无障碍和键盘要逐平台验证；手感不如原生。
- 适合：品牌 > 平台一致性，Web/桌面为主的产品。

### C. Hybrid（推荐）
一套 props、两个实现：

```go
ds.Button(ds.ButtonProps{Text: "Save", Variant: ds.Primary, OnPress: save})   // 自绘
ds.NativeButton(ds.ButtonProps{Text: "Save", Variant: ds.Primary, OnPress: save}) // 平台原生
```

默认自绘保证品牌一致；在 iOS/Android 的关键路径（支付、登录、表单）换 `NativeButton`，拿到 Liquid Glass / Material 按钮 / ripple。`Input` 永远是"原生输入 + DS 外壳"，这是两个世界的最佳结合点。

原型的 `demo` 里就并排展示了三种：原生填充、原生描边、自绘。

---

## 5. 原型：已经跑起来的东西

```
asuka-ds/
  theme.go      主题与 token：Palette / Radii / TypeScale；Native() 与 Neutral() 两套
  button.go     Button + NativeButton，6 变体 × 4 尺寸，hover/pressed/disabled
  card.go       Card / CardHeader / CardTitle / CardDescription / CardContent / CardFooter
  badge.go      Badge / Separator / Skeleton
  input.go      Input：原生 Plain TextField + label/description/error/focus ring
  tabs.go       Tabs：自绘分段控件，signal 驱动
  text.go       Heading / Title / Text / Muted / Screen（页面外壳）
  demo/         全组件展示 + 主题/外观切换
  demo-form/    表单专项：验证原生输入框在自绘外壳里的表现
```

设计要点：

- **零值可用**：`ds.Button(ButtonProps{Text: "Save"})` 就是主按钮中号。
- **token 是结构体，不是字符串**：`Palette` 的每个字段是一个 `ui.Color`，类型安全，能直接做对比度计算。
- **主题是 signal**：`ds.SetTheme(ds.Neutral())` 会让所有 `StyleFunc` 重跑，实现不重建页面的实时换肤（原型里两个 `Tabs` 就是切换器）。
- **两套调色板**：
  - `Native()`：映射到 `colors.*` 语义色，深色模式、增强对比度、Android 动态色全自动；代价是桌面端语义背景会"塌陷"（见 §6）。
  - `Neutral()`：shadcn zinc 的明/暗两套十六进制值，跨平台一致。
- **`ds.Themed(func(t, c) ui.Style{...})`**：给库没覆盖的部分留的逃生口，用户写的样式同样参与换肤。

验证结果：

| 检查 | 结果 |
|---|---|
| `go build ./...` | ✅ |
| `asuka build --platform web demo` | ✅ 7.4 MB wasm（gzip 2.0 MB） |
| `asuka build demo`（macOS .app） | ✅ ad-hoc 签名 |
| `asuka vet ./...` | ✅ 无告警 |
| Web：浅色 / 深色 / Native / Neutral / hover / Tab 焦点 / 390pt | ✅ 见 `docs/screenshots/` |
| macOS：AppKit 窗口、Plain TextField 焦点环、错误态、原生 Toggle | ✅ 见 `docs/screenshots/macos-*.png` |
| iOS / Android | 未验证（本机无模拟器环境），组件只用了跨平台 API |

截图证据：

- `docs/screenshots/web-neutral-light.png` / `web-neutral-dark.png`
- `docs/screenshots/web-focus-ring.png`（键盘焦点由平台绘制）
- `docs/screenshots/web-phone.png`（390pt 自动换行）
- `docs/screenshots/macos-native-dark.png`（原生宿主 + 系统强调色）
- `docs/screenshots/macos-form.png`（原生输入框、焦点、错误、原生开关）

---

## 6. 实测发现的坑（组件库必须处理）

1. **桌面端语义背景色塌陷**：macOS 上 `SecondaryBackground` 与窗口背景几乎一致，"Secondary" 按钮在深色下看起来像 Ghost。需要引入*elevation token*（例如 `Surface2 = Foreground 8%`）而不是复用语义背景。
2. **`ui.Style` 的零值不是"透明"而是"不设置"**：`Color` 是 `uint64`，零值表示用平台默认。要"无色"必须显式 `colors.Clear`。写通用组件时容易踩。
3. **`View` 的 `Style` 会复制给孩子**：把带背景/圆角的样式同时传给容器和文字，文字会盖住圆角。原型里已拆成 `itemStyle` 和 `labelStyle`。
4. **`app.RunWindow` 的尺寸在这台机器上没生效**：请求 900×900，窗口仍是 800×632（800×600 内容 + 标题栏）。桌面端尺寸策略需要单独验证。
5. **自绘按钮的 `Disabled` 只能"看起来禁用"**：语义上仍是按钮。若产品有无障碍合规要求，禁用态建议改用原生 `Button`（`Disabled` 是真语义）。
6. **Tabs 的选中态对读屏不可见**：自绘 Tabs 没有 selected trait；对无障碍敏感的页面用 `ui.Segmented`。

---

## 7. 建议的库结构

```
asuka-ds/
  ds/                 # 一个包，按组件分文件（同 shadcn 的 components/ui）
    theme.go
    button.go
    ...
  registry/           # 可选：给 copy-in 模式用的组件清单
    index.json
  cmd/asuka-ds/       # 可选：`asuka-ds add button` 把源码拷进 app
  docs/               # 用 Asuka 自己构建的文档站（web 平台）
  examples/showcase/
```

**分发方式建议**：

- **主线：Go module**（`go get github.com/you/asuka-ds`）。Go 的版本语义、`go doc`、依赖管理都比 copy-paste 成熟；用户想改某个组件时，`go mod vendor` 或直接 fork 单文件即可。
- **补充：registry CLI**。shadcn 的灵魂是"你拥有这段代码"。Go 里可以做得更彻底：`asuka-ds add button` 把 `button.go` + 它依赖的 token 片段写进 `app/ui/`，之后完全脱离库版本演进。建议 M4 再做。
- **不要做**：运行时主题包 / 插件。Asuka 的插件机制是给 C/Zig/Rust 动态库用的，组件库不需要。

**API 约定**：

- `XxxProps` 结构体；零值 = 合理默认（`Variant: 0` 是 Primary，`Size: 0` 是 Medium）。
- 动态文本用 `signal.Accessor[string]`，静态用 `ui.Val("...")`。
- 变体用枚举，不用字符串（Go 里字符串变体没有编译期检查）。
- 所有颜色/圆角/字号只从 `Theme` 来；组件内部不出现十六进制。
- 每个组件都能被 `ds.Themed` 覆盖，不让用户和库打架。

---

## 8. 路线图

| 阶段 | 内容 | 产出 |
|---|---|---|
| **M0 原型**（已完成） | token、双主题、Button/Card/Badge/Input/Tabs、双 demo、Web+macOS 验证 | 本仓库 |
| **M1 基础设施**（已完成） | 冻结 `Theme`/`Palette`/`TypeScale` API；补 elevation、motion token；WCAG AA 对比度单测；`Theme.Validate()` / `MotionFor()` | `theme.go`、`theme_test.go` |
| **M2 原语与无障碍**（已完成） | `Pressable`（动态禁用语义、焦点、hover）、`Field`（label/desc/error + 关联）、`Overlay`（Sheet/Popover 包装）、`Menu`/`Toast` 包装；29 个 `uitest` | `pressable.go`、`field.go`、`overlay.go`、`audit_test.go` |
| **M3 组件扩充**（已完成） | Switch/Checkbox/Radio/Select/Slider/Dialog/Dropdown/Toast/Tooltip/Avatar/Progress/Table/Accordion/Command | 组件源码 + 双 demo |
| **M4 分发与文档** | registry + `asuka-ds add`；文档站（Asuka web 构建，双主题并排 + 实时改 token）；版本兼容策略 | v1.x |
| **M5 质量** | `uitest` 快照、`asuka vet`、iOS/Android 真机走查、CI 对比度检查 | 持续 |

组件优先级建议：先做"高频且好做"的（Card、Badge、Input、Tabs、Button），再做"必须原生"的（Select、DatePicker、Sheet），最后做"纯自绘成本高"的（Command palette、Data table）。

## 8.1 M1→M2→M3 验收记录（2026-11-14）

本轮把调研文档中的路线图落到了 M3，并按 `docs/superpowers/plans/2026-11-14-asuka-ds-m1-m3.md` 验收：

- `Theme` / `Palette` / `Radii` / `TypeScale` / `Motion` / `Elevation` 已冻结为强类型 API；`Native()` 与 `Neutral()` 均通过 `Validate()`。
- `Pressable` 统一 disabled / selected / expanded 提示；动态禁用后事件不再触发，hint 更新为 `Unavailable`，样式同步变淡。
- `Field` 统一 label、description、error 与 accessible name；`Input` 已改为复用 `Field`，错误与描述可同时显示。
- `Overlay` 明确区分 `Content` 与 `Popover` 锚点 children；`Sheet` 不接收 children。
- `Select` disabled 时不再挂载 `ui.Menu`；`Menu` 补齐 action role；`Tabs` 改用 `Pressable` 输出 selected 提示。
- `Command` 在过滤结果变少时把选中项夹回有效范围，保证始终有且仅有一个选中行。
- `go test ./...`：29 项通过；`go build ./...`、`asuka vet ./...` 通过；Web 与 macOS 构建产物已生成。
- iOS / Android 仍未做真机走查，组件只用了跨平台 API；这是 M5 的范围。

新增回归测试均先写失败用例、观察红测，再做最小实现并观察绿测：`audit_test.go` 覆盖动态禁用、Tabs 选中、overlay 内容、禁用 Select、Input 字段语义与 Command 过滤边界。

---

## 9. 风险与建议

| 风险 | 影响 | 应对 |
|---|---|---|
| Asuka 还在 0.1.x，API 可能变 | 组件库跟着 breaking | 锁定 `v0.1.1`；把所有 `ui.*` 调用收敛到 `ds/internal/ask` 薄封装；CI 跑 `go test` + `asuka vet` |
| 原生与自绘混用的视觉割裂 | 用户体验不一致 | 以"组件级"选择后端，不要在同一个控件组里混；文档明确默认值 |
| 自绘控件的无障碍缺口 | 合规风险 | 关键路径用原生；给 Asuka 提 issue（`ViewProps.Disabled`、selected state、可更新 Label）；短期用 `Show` 重建节点 + `Hint` 兜底 |
| 主题全局单例 | 一个页面不能同时展示两套品牌 | 文档站用两个 app 实例 / 切换器；未来可给 Asuka 提"自定义 context"需求 |
| 桌面端语义色塌陷 | 深色下层次丢失 | 自带 elevation token，不依赖 `SecondaryBackground`；每个平台出基线截图 |
| 自绘控件没有平台手感 | "不像原生应用" | 明确产品定位：品牌一致性优先；提供 `NativeButton` 等原生变体 |

**给上游的两个具体建议**（可以直接开 issue）：

1. `ViewProps.Disabled`（或 `Interaction.Disabled`）：让可按压 View 真正被读屏软件读成"已禁用"。
2. 可按压 View 的 `Selected` / `Checked` / `Expanded` 状态：Tabs、Accordion、Switch、Menu item 都需要。

---

## 10. 怎么跑

```sh
go install github.com/egoist/asuka/cmd/asuka@latest

cd asuka-ds
asuka dev demo                  # 本机（macOS 上是 AppKit 窗口）
asuka dev --platform web demo   # 浏览器，localhost:8790
asuka build --platform web demo # 静态站点到 demo/dist/web
asuka build demo                # macOS .app
asuka vet ./...
```

文档站建议直接复用 `demo/`：Asuka 的 Web 构建产物是纯静态站点，可以部署到任意 CDN（`_headers`/`_redirects` 已生成）。
