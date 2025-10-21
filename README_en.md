[中文](README.md) | [English](README_en.md)

# ctx

A simple tool for managing kubectl contexts.

## Install

```bash
git clone https://github.com/uniquejava/ctx.git
cd ctx
make build
```

Add `bin/ctx` to your PATH.

## Usage

```bash
ctx                          # List contexts
ctx use <context>             # Switch context
ctx use <context> <namespace> # Switch context and set namespace
ctx rm <context>              # Remove context
```

## License

MIT