# 12306-MCP

基于 Model Context Protocol (MCP) 的12306购票搜索服务器。提供了简单的API接口，允许大模型利用接口搜索12306购票信息。

## 功能特点

* 查询12306购票信息

## 待办事项

* 完成12306其余接口（中转，过站查询等）

## 安装

```bash
git clone <your-repo-url>
cd 12306-mcp
go mod init 12306-mcp
go mod tidy
```

## 快速开始


1. 运行服务器

```bash
go run main.go
```


## API 示例

查询车票信息：

```json
{
    "from_station": "杭州",
    "to_station": "武汉",
    "date": "2025-04-30"
}
```

响应示例：

```json
[
  {
	"车次": "G1442",
	"出发站": "杭州西",
	"到达站": "武汉",
	"出发时间": "2025-04-26 06:07:00",
	"到达时间": "2025-04-26 10:23:00",
	"历时": "04:16",
	"余票信息": [{
			"座位类型": "商务座",
			"余票数量": 3,
			"价格": 1018
		},
		{
			"座位类型": "一等座",
			"余票数量": 10,
			"价格": 528
		},
		{
			"座位类型": "二等座",
			"余票数量": 21,
			"价格": 329
		}
	]
}
]
```
## 在Cursor中使用

要在Cursor IDE中使用此MCP服务：

1. 确保12306-MCP服务已在本地运行
2. 在Cursor设置中启用MCP功能
3. 本地mcp.json添加：
```
    "12306-mcp": {
      "url": "http://127.0.0.1:8080/sse"
    }
```

4. 连接后即可通过对话框使用以下查询示例：
   ```
   请帮我查询杭州到武汉明天的高铁票
   ```
5. Cursor将通过MCP协议调用本服务获取实时车票信息





## 参考

* [Model Context Protocol](https://github.com/modelcontextprotocol/modelcontextprotocol)
* [MCP Go SDK](https://github.com/mark3labs/mcp-go)

## 许可证

MIT License 