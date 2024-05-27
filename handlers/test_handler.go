package handlers

import (
  "context"
  "github.com/bytedance/sonic"
  "github.com/cloudwego/hertz/pkg/app"
  "github.com/cloudwego/hertz/pkg/common/utils"
)

func Test(ctx context.Context, c *app.RequestContext) {
  body, err := c.Body()
  if err != nil {
    c.JSON(500, utils.H{"error": err})
    return
  }
  var jsonObj = utils.H{}
  err = sonic.ConfigDefault.Unmarshal(body, &jsonObj)
  if err != nil {
    c.JSON(500, utils.H{"error": err})
    return
  }
  c.JSON(200, jsonObj)
}
