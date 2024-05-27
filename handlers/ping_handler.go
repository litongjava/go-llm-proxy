package handlers

import (
  "context"
  "github.com/cloudwego/hertz/pkg/app"
  "github.com/cloudwego/hertz/pkg/common/utils"
)

func PingCompletions(ctx context.Context, c *app.RequestContext) {
  c.JSON(200, utils.H{
    "status": "running",
  })
}
