# go-llm-proxy

## buid
```
go build
```
## run 
```
v
```
## test
```
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
    "stream":true
}'

available models

- gpt-3.5-turbo
- gpt-4o