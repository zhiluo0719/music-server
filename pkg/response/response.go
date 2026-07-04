package response

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type Response struct {
	Code int         `json:"code"`
	Msg  string      `json:"msg"`
	Data interface{} `json:"data,omitempty"`
}

type PageData struct {
	List     interface{} `json:"list"`
	Total    int64       `json:"total"`
	Page     int         `json:"page"`
	PageSize int         `json:"page_size"`
}

func Success(c *gin.Context, data interface{}) {
	c.JSON(http.StatusOK, Response{Code: 0, Msg: "ok", Data: data})
}

func SuccessPage(c *gin.Context, list interface{}, total int64, page, pageSize int) {
	c.JSON(http.StatusOK, Response{
		Code: 0,
		Msg:  "ok",
		Data: PageData{List: list, Total: total, Page: page, PageSize: pageSize},
	})
}

func Error(c *gin.Context, code int, msg string) {
	c.JSON(http.StatusOK, Response{Code: code, Msg: msg})
}

func BadRequest(c *gin.Context, msg string) {
	Error(c, 400, msg)
}

func NotFound(c *gin.Context, msg string) {
	Error(c, 404, msg)
}

func InternalError(c *gin.Context, msg string) {
	Error(c, 500, msg)
}
