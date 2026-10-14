## SwiftUI Design Skill v1.0.0

> 一句话，一个能交付的 SwiftUI 界面。

首个正式版本。跨 agent 通用的 SwiftUI 前端视觉设计 skill——不 AI slop。

### ✨ 内置能力

| 能力 | 说明 |
|------|------|
| 反 AI Slop 六条铁律 | 6 条硬规则 + 详细参考文档 |
| 设计方向顾问 | 5 大设计流派 × N 种哲学，推荐 3 方向 |
| 品牌资产协议 | 5 步硬流程，从问到写 brand-spec.md |
| 五维评审 | 哲学一致性 / 视觉层级 / 细节执行 / 功能性 / 原创性 |
| 布局模式库 | 9 种常见布局，可复制 SwiftUI 代码 |
| 字体与色彩系统 | 层级字体 + 色彩 token + 设计系统模板 |
| Swift 扩展 | Color(hex:)、动画模式、SF Symbols 指南 |

### 📦 安装

```bash
npx skills add wholiver/swiftui-design-skill -g -y
```

### 📁 文件结构

```
swiftui-design-skill/
├── SKILL.md                      ← 核心定义（278 行）
├── references/
│   ├── anti-ai-slop.md           ← 反 AI Slop 详细规则（268 行）
│   ├── layout-patterns.md        ← 9 种布局模式（265 行）
│   ├── typography-color.md       ← 字体层级 + 色彩系统（172 行）
│   ├── design-review.md          ← 五维评审流程（151 行）
│   └── swift-extensions.md       ← Color/Font/Animation 扩展（373 行）
└── templates/
    └── brand-spec.md             ← 品牌规范模板（77 行）
```

### 🎯 与 swiftui-expert-skill 的区别

| | SwiftUI Design Skill | swiftui-expert-skill |
|---|---|---|
| 焦点 | 视觉设计、美学、品牌感 | 代码质量、性能、正确性 |
| 问题 | 「看起来怎么样？」 | 「代码写得对不对？」 |

### 📋 Changelog

- ✅ 反 AI Slop 六条铁律 + 详细参考
- ✅ 设计方向顾问（5 流派）
- ✅ 品牌资产协议（5 步流程）
- ✅ 五维评审系统
- ✅ 9 种布局模式 + 可复制代码
- ✅ 字体层级与色彩系统
- ✅ Swift 扩展（Color/Font/Animation）
- ✅ 品牌规范模板
- ✅ MIT License
