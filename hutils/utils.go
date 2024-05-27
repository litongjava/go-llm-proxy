package hutils

import (
  "github.com/cloudwego/hertz/pkg/app"
  "github.com/cloudwego/hertz/pkg/common/utils"
  "net/http"
)

func Fail400(reqCtx *app.RequestContext, err error) {
  reqCtx.JSON(http.StatusBadRequest, utils.H{
    "msg":  "Invalid request:" + err.Error(),
    "ok":   false,
    "code": 0,
  })
}

func Fail500(reqCtx *app.RequestContext, err error) {
  reqCtx.JSON(http.StatusInternalServerError, utils.H{
    "msg":  "Server Error:" + err.Error(),
    "ok":   false,
    "code": 0,
  })
}

func Ok(reqCtx *app.RequestContext, data interface{}) {
  reqCtx.JSON(http.StatusOK, utils.H{
    "code": 1,
    "data": data,
    "ok":   true,
  })
}
