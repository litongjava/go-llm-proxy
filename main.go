package main

import (
  "github.com/cloudwego/hertz/pkg/app/server"
  "go-llm-proxy/handlers"
)

func main() {
  h := server.Default()
  h.GET("/", handlers.PingCompletions)
  h.GET("/ping", handlers.PingCompletions)
  h.POST("/test", handlers.Test)
  h.POST("/v1/chat/completions", handlers.OpenaiV1ChatCompletions)
  h.POST("/openai/v1/chat/completions", handlers.OpenaiV1ChatCompletions)

  h.Spin()
}
