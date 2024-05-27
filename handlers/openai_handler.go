package handlers

import (
  "context"
  "encoding/json"
  "github.com/bytedance/sonic"
  "github.com/cloudwego/hertz/pkg/app"
  "github.com/cloudwego/hertz/pkg/app/client"
  "github.com/cloudwego/hertz/pkg/common/hlog"
  "github.com/cloudwego/hertz/pkg/protocol"
  "github.com/hertz-contrib/sse"
  "go-llm-proxy/hutils"
  "strings"
)

var OpenAiChatCompletionUrl = "https://api.openai.com/v1/chat/completions"

func OpenaiV1ChatCompletions(ctx context.Context, reqCtx *app.RequestContext) {
  // Read request body and headers
  body, _ := reqCtx.Body()
  // Decode the body into a map
  var requestMap map[string]interface{}
  sonic.ConfigDefault.Unmarshal(body, &requestMap)

  headers := make(map[string]string)
  reqCtx.Request.Header.VisitAll(func(key, value []byte) {
    headers[string(key)] = string(value)
  })
  //headers.put("host", "api.openai.com");
  headers["Host"] = "api.openai.com"
  hlog.Info("body:", string(body))
  hlog.Info("headers:", headers)
  if requestMap["stream"] == true {
    // Setup client to connect to the remote server
    client := sse.NewClient(OpenAiChatCompletionUrl)
    client.SetMethod("POST")
    client.SetHeaders(headers)
    client.SetBody(body)

    // Setup the stream for the original client
    var sEvent *sse.Stream = sse.NewStream(reqCtx)
    errChan := make(chan error)
    var completeContent strings.Builder
    go func() {
      err := client.Subscribe(func(msg *sse.Event) {
        if msg.Data != nil {
          // Forwarding the received event back to the original client
          event := &sse.Event{
            Data: msg.Data,
          }
          err := sEvent.Publish(event)
          if err != nil {
            hlog.CtxErrorf(ctx, "failed to send event to client: %sEvent", err)
            return
          }
          printResponseContent(msg, &completeContent)
        }
      })
      errChan <- err
    }()

    select {
    case err := <-errChan:
      if err != nil {
        hlog.CtxErrorf(ctx, "error from remote server: %sEvent", err)
        hutils.Fail500(reqCtx, err)
        return
      }
    }
  } else {
    httpClient, _ := client.NewClient()
    request := &protocol.Request{}
    response := &protocol.Response{}
    request.SetRequestURI(OpenAiChatCompletionUrl)
    request.SetMethod("POST")
    request.SetHeaders(headers)
    request.SetBody(body)

    // 执行请求
    var err = httpClient.Do(context.Background(), request, response)
    defer response.CloseBodyStream()
    if err != nil {
      hutils.Fail500(reqCtx, err)
      return
    }
    // 设置响应头和状态码
    response.Header.VisitAll(func(key, value []byte) {
      reqCtx.Response.Header.Set(string(key), string(value))
    })
    // 设置响应体流
    reqCtx.Response.SetBody(response.Body())
  }
}

func printResponseContent(msg *sse.Event, completeContent *strings.Builder) {
  var responseData struct {
    Choices []struct {
      Delta struct {
        Content string `json:"content"`
      } `json:"delta"`
    } `json:"choices"`
  }
  if len(msg.Data) > 6 {
    if err := json.Unmarshal(msg.Data, &responseData); err == nil {
      for _, choice := range responseData.Choices {
        //hlog.Info("content:", choice.Delta.Content)
        completeContent.WriteString(choice.Delta.Content)
      }
    } else {
      hlog.Error("json parsing error: %s", err)
    }
  }
  // Check if it's the end of the data
  if string(msg.Data) == "[DONE]" {
    hlog.Info("Complete content:", completeContent.String())
    completeContent.Reset() // Clear the builder for the next message
  }
}
