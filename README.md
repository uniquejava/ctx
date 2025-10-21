[中文](README.md) | [English](README_en.md)

# ctx

管理kubectl contexts的简单工具。

## 安装

```bash
git clone https://github.com/uniquejava/ctx.git
cd ctx
make build
```

将 `bin/ctx` 加入PATH。

## 使用

```bash
ctx                          # 列出contexts
ctx use <context>             # 切换context
ctx use <context> <namespace> # 切换context并设置namespace
ctx rm <context>              # 删除context
```

## 许可证

MIT