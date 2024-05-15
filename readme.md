# go-llm-proxy

## Introduction

When using third-party tools like LangChain to call large OpenAI models, these tools process the prompts in ways that can obscure what specific information the application is sending to the models. This tool provides a proxy service that captures and prints messages sent to GPT models and their responses to the console, facilitating developers in identifying and solving issues.

## Build

Compile this project using Go:

```sh
go build
```

## Run

Run the compiled executable:

```sh
./go-llm-proxy
```

## Test

Use curl to test if the proxy service can correctly forward requests and display return data:

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

The models currently supported include:
[models](https://platform.openai.com/docs/models)

Please configure the specific model the proxy should connect to as needed.

## Contribution

Community developers are welcome to contribute code to add new features or improve existing functionalities. Please submit Pull Requests to our GitHub repository.

## License

This project is licensed under the MIT License, see the `LICENSE` file for details.