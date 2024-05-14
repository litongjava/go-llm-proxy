package handlers

import (
  "context"
  "encoding/json"
  "github.com/cloudwego/hertz/pkg/app"
  "github.com/cloudwego/hertz/pkg/common/hlog"
  "github.com/hertz-contrib/sse"
  "strings"
)

func OpenaiV1ChatCompletions(ctx context.Context, c *app.RequestContext) {
  // Read request body and headers
  body, _ := c.Body()
  headers := make(map[string]string)
  c.Request.Header.VisitAll(func(key, value []byte) {
    headers[string(key)] = string(value)
  })
  //headers.put("host", "api.openai.com");
  headers["Host"] = "api.openai.com"
  hlog.Info("body:", string(body))
  hlog.Info("headers:", headers)

  // Setup client to connect to the remote server
  client := sse.NewClient("https://api.openai.com/v1/chat/completions")
  client.SetMethod("POST")
  client.SetHeaders(headers)
  client.SetBody(body)

  // Setup the stream for the original client
  s := sse.NewStream(c)
  errChan := make(chan error)
  var completeContent strings.Builder
  go func() {
    err := client.Subscribe(func(msg *sse.Event) {
      if msg.Data != nil {
        // Forwarding the received event back to the original client
        event := &sse.Event{
          Data: msg.Data,
        }
        err := s.Publish(event)
        if err != nil {
          hlog.CtxErrorf(ctx, "failed to send event to client: %s", err)
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
      hlog.CtxErrorf(ctx, "error from remote server: %s", err)
      return
    }
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
        hlog.Info("content:", choice.Delta.Content)
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
