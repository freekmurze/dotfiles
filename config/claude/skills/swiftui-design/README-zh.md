<div align="center">

# SwiftUI Design Skill

> *「一句话，一个能交付的 SwiftUI 界面。」*
> *"One prompt. A shippable SwiftUI interface."*

[![License](https://img.shields.io/badge/License-MIT-blue.svg)](LICENSE)
[![Agent-Agnostic](https://img.shields.io/badge/Agent-Agnostic-blueviolet)](https://skills.sh)
[![Skills](https://img.shields.io/badge/skills.sh-Compatible-green)](https://skills.sh)

</div>

<p align="center">
  <a href="README.md"><img src="https://img.shields.io/badge/English-blue?style=flat-square" alt="English"></a>
  <a href="README-zh.md"><img src="https://img.shields.io/badge/中文-red?style=flat-square" alt="中文"></a>
</p>

**在你的 agent 里打一句话，拿回一个不 AI slop 的 SwiftUI 界面。**

3 到 30 分钟，你能设计出一个**有签名感的 iOS 界面**——不是「AI 做的还行」那种水平，是看起来像有品味的设计师做的。

反 AI slop 六条铁律、设计方向顾问、品牌资产协议、五维评审——全部内置。给 skill 你的品牌色板，它会读懂你的气质；什么都不给，内置的 5 种设计语汇也能兜底到不出丑。

跨 agent 通用——Claude Code、Cursor、Codex、OpenCode、Hermes 都能装。

## 装上就能用

```bash
npx skills add wholiver/swiftui-design-skill -g -y
```

## 能做什么

| 能力 | 交付物 | 典型耗时 |
|------|--------|----------|
| SwiftUI 界面设计 | 可编译的 SwiftUI 代码 · 设计系统 token · 布局模式 | 10–15 min |
| 设计方向顾问 | 5 流派 × N 种设计哲学 · 推荐 3 方向 · 视觉锚点描述 | 5 min |
| 品牌资产集成 | brand-spec.md · 色板 · 字体 · 间距系统 | 5–10 min |
| 反 AI Slop 审查 | 逐条检查 · 具体修复代码 · 替代方案 | 3–5 min |
| 五维评审 | 雷达图评分 · Keep/Fix/Quick Wins · 可操作修复清单 | 3 min |
| 布局模式库 | 9 种常见布局 · 可复制的 SwiftUI 代码 | 即时 |
| 动画设计 | spring/parallax/pull-to-refresh · 可编译代码 | 5–10 min |

## 核心机制

### 反 AI Slop 六条铁律

这是 skill 里最硬的一段规则。所有设计都必须通过这 6 条检查：

| # | 铁律 | ❌ 禁止 | ✅ 替代方案 |
|---|------|---------|------------|
| 1 | 起点从上下文开始 | 空白画布凭空捏造 | 先问设计系统/UI kit/品牌资产 |
| 2 | 初级设计师模式 | 等完美方案 | 灰色占位符 > 烂 SVG |
| 3 | 给变体不给终稿 | 一个「最终答案」 | 2–3 个差异化方向 |
| 4 | 占位符 > 烂实现 | AI 生成的剪贴画 | 干净的灰色占位 + 文字标注 |
| 5 | 系统优先，别填充 | 塞满每一个像素 | 每个元素必须证明自己存在的理由 |
| 6 | 禁止 AI slop 模式 | 紫蓝渐变、emoji 图标、圆角卡片+左边框 | 单一暖色强调色、SF Symbols、衬线标题 |

详见 `references/anti-ai-slop.md`。

### 设计方向顾问

当需求模糊时，skill 会从 5 大设计流派中推荐 3 个差异化方向：

| 流派 | 特征 | 代表风格 |
|------|------|----------|
| **信息型** | 数据优先、图表密集 | Bloomberg Terminal |
| **编辑型** | 杂志排版、衬线字体、大量留白 | NYT、Medium |
| **表现型** | 大胆配色、不对称布局、动效前置 | Stripe、Linear |
| **功能型** | 密集工具感、等宽字体点缀 | Things、OmniFocus |
| **温暖极简** | 柔和中性色、圆角、微妙质感 | Notion、Craft |

### 品牌资产协议

涉及具体品牌时强制执行 5 步硬流程：

| 步骤 | 动作 | 目的 |
|------|------|------|
| 1 · 问 | 用户有 brand guidelines 吗？ | 尊重已有资源 |
| 2 · 搜 | 搜索品牌官方页面 | 获取真实素材 |
| 3 · 下 | 下载实际资产文件 | PNG/SVG logo、字体 |
| 4 · 验 | 验证颜色匹配官方来源 | 核对 hex 值 |
| 5 · 写 | 生成 brand-spec.md | 记录完整设计系统 |

质量门槛：5 个真实品牌色、10 个设计 token、2 个字体家族、8pt 间距网格。

### 五维评审

每个设计交付前必须通过 5 个维度的评审：

| 维度 | 评分标准 | 最低分 |
|------|----------|--------|
| 🎯 哲学一致性 | 设计是否贯彻了选定的设计哲学 | ≥ 7/10 |
| 📐 视觉层级 | 信息优先级是否清晰 | ≥ 7/10 |
| 🔧 细节执行 | 间距、字体、颜色是否精确 | ≥ 7/10 |
| ⚡ 功能性 | 布局是否服务用户目标 | ≥ 7/10 |
| ✨ 原创性 | 是否有 1 个签名细节 | ≥ 7/10 |

## 文件结构

```
swiftui-design-skill/
├── SKILL.md                           ← 核心定义（278 行）
├── references/
│   ├── anti-ai-slop.md                ← 反 AI Slop 详细规则（268 行）
│   ├── layout-patterns.md             ← 9 种布局模式 + 可复制代码（265 行）
│   ├── typography-color.md            ← 字体层级 + 色彩系统（172 行）
│   ├── design-review.md               ← 五维评审流程（151 行）
│   └── swift-extensions.md            ← Color/Font/Animation 扩展（373 行，必备）
└── templates/
    └── brand-spec.md                  ← 品牌规范模板（77 行）
```

## 与 swiftui-expert-skill 的区别

| | SwiftUI Design Skill | swiftui-expert-skill |
|---|---|---|
| **焦点** | 视觉设计、美学、品牌感 | 代码质量、性能、正确性 |
| **问题** | 「看起来怎么样？」 | 「代码写得对不对？」 |
| **输出** | 设计方向、色板、布局、评审 | 代码审查、Instruments 分析、API 现代化 |
| **互补** | 先用 design 确定方向 | 再用 expert 确保代码质量 |

两个 skill 可以配合使用：Design 确定「做什么」，Expert 确保「怎么做对」。

## License

MIT — 随便用，但请保留原作者署名。
