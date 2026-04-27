# dong-labs

咚咚家族 (DongDong Family) - AI 原生的个人数据管理 CLI 工具集。

## 工具列表

| CLI | 命令 | 说明 | 安装命令 |
|-----|------|------|---------|
| **dong-think** | `dong-think` | 思咚咚 - 记录灵感和想法 | `brew install dong-labs/tap/dong-think` |
| **dong-log** | `dong-log` | 记咚咚 - 日常日志记录 | `brew install dong-labs/tap/dong-log` |
| **dong-read** | `dong-read` | 读咚咚 - 个人知识管理 | `brew install dong-labs/tap/dong-read` |
| **dong-dida** | `dong-dida` | 事咚咚 - 待办事项管理 | `brew install dong-labs/tap/dong-dida` |
| **dong-cang** | `dong-cang` | 仓咚咚 - 个人财务管理 | `brew install dong-labs/tap/dong-cang` |
| **dong-expire** | `dong-expire` | 到期咚 - 订阅到期管理 | `brew install dong-labs/tap/dong-expire` |
| **dong-pass** | `dong-pass` | 密码咚 - 账号密码管理 | `brew install dong-labs/tap/dong-pass` |
| **dong-timeline** | `dong-timeline` | 时间咚 - 里程碑记录 | `brew install dong-labs/tap/dong-timeline` |
| **dong-member** | `dong-member` | 会员咚 - 会员信息管理 | `brew install dong-labs/tap/dong-member` |

## 安装

### macOS (推荐)

使用 Homebrew 安装单个工具：

```bash
# 安装思咚咚
brew install dong-labs/tap/dong-think

# 安装记咚咚
brew install dong-labs/tap/dong-log

# 或安装全部工具
brew install dong-labs/tap/dong-think \
                 dong-labs/tap/dong-log \
                 dong-labs/tap/dong-read \
                 dong-labs/tap/dong-dida \
                 dong-labs/tap/dong-cang \
                 dong-labs/tap/dong-expire \
                 dong-labs/tap/dong-pass \
                 dong-labs/tap/dong-timeline \
                 dong-labs/tap/dong-member
```

### Go 开发者

```bash
# 安装单个工具
go install github.com/dong-labs/think/cmd/dong-think@latest

# 确保 $GOPATH/bin 在 $PATH 中
export PATH=$PATH:$(go env GOPATH)/bin
```

### 从源码编译

```bash
git clone https://github.com/dong-labs/dong-labs-go.git
cd dong-labs-go
go build -o dong-think ./cmd/dong-think
sudo mv dong-think /usr/local/bin/
```

## 快速开始

### dong-think (思咚咚)

```bash
# 初始化
dong-think init

# 添加想法
dong-think add "今天有个好主意"

# 列出所有想法
dong-think list

# 搜索
dong-think search "主意"
```

### dong-log (记咚咚)

```bash
# 初始化
dong-log init

# 记录日志
dong-log add "今天完成了项目迁移"

# 查看统计
dong-log stats
```

### dong-dida (事咚咚)

```bash
# 初始化
dong-dida init

# 添加待办
dong-dida add "完成代码重构" --priority high

# 列出待办
dong-dida list
```

## 数据存储

所有数据存储在本地 `~/.dong/` 目录：

```
~/.dong/
├── config.json       # 统一配置文件
├── think.db          # 思咚咚数据
├── log.db            # 记咚咚数据
├── read.db           # 读咚咚数据
├── dida.db           # 事咚咚数据
├── cang.db           # 仓咚咚数据
├── expire.db         # 到期咚数据
├── pass.db           # 密码咚数据
├── timeline.db       # 时间咚数据
└── member.db         # 会员咚数据
```

**隐私承诺：** 所有数据本地存储，不上云、不同步、不追踪。

## 设计原则

1. **AI First, Human Second** - 所有命令优先为 AI 调用设计
2. **JSON Native** - 每个命令返回结构化 JSON，方便程序处理
3. **Local & Private** - 数据完全本地化，保护隐私
4. **Minimal Core** - 每个工具只做一件事，做好一件事

## 输出格式

所有命令统一返回 JSON：

```json
{
  "success": true,
  "data": {
    // 响应数据
  }
}
```

```json
{
  "success": false,
  "error": {
    "code": "ERROR_CODE",
    "message": "错误信息"
  }
}
```

## 开发

```bash
# 克隆仓库
git clone https://github.com/dong-labs/dong-labs-go.git
cd dong-labs-go

# 构建
go build -o dong-think ./cmd/dong-think

# 运行
./dong-think init
```

## License

MIT License

## 相关链接

- [Homebrew Tap](https://github.com/dong-labs/homebrew-tap)
- [GitHub](https://github.com/dong-labs)
