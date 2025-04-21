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
    "train_no": "G1234",
    "from_station": "北京",
    "to_station": "上海",
    "depart_time": "2024-04-15T10:00:00+08:00",
    "arrive_time": "2024-04-15T13:00:00+08:00",
    "duration": "3小时",
    "seats": [
      {
        "type": "商务座",
        "count": 10,
        "price": 553.5
      },
      {
        "type": "一等座",
        "count": 20,
        "price": 333.5
      },
      {
        "type": "二等座",
        "count": 100,
        "price": 208.5
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