package tests

import (
  "context"
  "github.com/cloudwego/hertz/pkg/common/hlog"
  "github.com/hertz-contrib/sse"
  "sync"
  "testing"
)

var wg sync.WaitGroup

func TestOpenaiStreamClient(t *testing.T) {
  wg.Add(1) // If you're only connecting one client, change this to 1
  go func() {
    c := sse.NewClient("https://llm-proxy.fly.dev/openai/v1/chat/completions")

    // Set the HTTP method to POST
    c.SetMethod("POST")
    // Set the authorization header
    headers := map[string]string{
      "Authorization": "",
      "Content-Type":  "application/json",
    }

    c.SetHeaders(headers)

    // Set the request body
    c.SetBody([]byte(`{"messages": [{"role": "system", "content": "Just say hi"}], "model": "gpt-3.5-turbo", "stream": true}`))

    // touch off when connected to the server
    c.SetOnConnectCallback(func(ctx context.Context, client *sse.Client) {
      hlog.Infof("client connect to server %s success with %s method", c.GetURL(), c.GetMethod())
    })

    // touch off when the connection is shutdown
    c.SetDisconnectCallback(func(ctx context.Context, client *sse.Client) {
      hlog.Infof("client disconnect to server %s success with %s method", c.GetURL(), c.GetMethod())
    })

    events := make(chan *sse.Event)
    errChan := make(chan error)
    go func() {
      cErr := c.Subscribe(func(msg *sse.Event) {
        if msg.Data != nil {
          events <- msg
          return
        }
      })
      errChan <- cErr
    }()
    for {
      select {
      case e := <-events:
        hlog.Info(e.Event)
        hlog.Info(string(e.Data))
      case err := <-errChan:
        hlog.CtxErrorf(context.Background(), "err = %s", err.Error())
        wg.Done()
        return
      }
    }
  }()
  wg.Wait()
}
