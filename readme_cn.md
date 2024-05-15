# go-llm-proxy

## Introduction

在使用如 LangChain 等第三方工具调用 OpenAI 大型模型时，这些工具会对 prompt
进行处理，开发者难以看清楚应用程序向大模型发送了哪些具体信息。本工具提供了一个代理服务，可以捕获并打印向 GPT
模型发送的消息和模型的回复到控制台，方便开发者定位和解决问题。

## Build

使用 Go 语言编译本项目：

```sh
go build
```

## Run

编译后运行生成的可执行文件：

```sh
./go-llm-proxy
```

## Test

使用 curl 测试代理服务是否能正确转发请求并显示返回数据：

```sh
curl --location --request POST 'http://127.0.0.1:8888/openai/v1/chat/completions' \
--header 'Authorization: Bearer <token>' \
--header 'Content-Type: application/json' \
--data-raw '{
    "messages": [
        {
            "role": "system",
            "content": "Just say hi"
        }
    ],
    "model": "gpt-3.5-turbo",
    "stream": true
}'
```

## Available Models

目前支持的模型包括：
[models](https://platform.openai.com/docs/models)

请根据需要配置代理连接的具体模型。

## Contribution

欢迎社区开发者贡献代码，增加新功能或改进现有功能。请提交 Pull Requests 到我们的 GitHub 仓库。

## License

本项目采用 MIT 许可证，详见 `LICENSE` 文件。