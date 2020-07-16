package web

import (
	"encoding/json"

	"github.com/astaxie/beego"
	"github.com/astaxie/beego/context"
)

type ResponseResult struct {
	State   int         `json:"code"`              //0代表业务没有执行成功,但没有报错,error说明;-1代表通用错误;-500代表没有权限
	Message *string     `json:"message,omitempty"` //存放需要给用户进行提示的信息
	Error   *string     `json:"error,omitempty"`   //存放报错信息,一般不放入
	Data    interface{} `json:"data,omitempty"`    //需要的时候会放入数据
}

func _NewResponseResult(state int, msg *string, e error, data interface{}) *ResponseResult {
	rr := &ResponseResult{
		State:   state,
		Message: msg,
		Data:    data,
	}
	if e != nil {
		emsg := e.Error()
		rr.Error = &emsg
	}
	return rr
}

//输出错误
func ResponseError(state int, msg string, e error, c *beego.Controller) {
	rr := _NewResponseResult(state, &msg, e, nil)
	ResponseJson(rr, c)
}

func ResponseState(state int, msg string, c *beego.Controller) {
	rr := _NewResponseResult(state, &msg, nil, nil)
	ResponseJson(rr, c)
}

//输出data
func ResponseData(data interface{}, msg string, c *beego.Controller) {
	rr := _NewResponseResult(1, &msg, nil, data)
	ResponseJson(rr, c)
}

func ResponseToContext(rr *ResponseResult, ctx *context.Context) {
	body, e := json.Marshal(rr)
	if e != nil {
		rr = _NewResponseResult(-1, rr.Message, e, nil)
		ResponseToContext(rr, ctx)
		return
	}
	ctx.Output.Header("Access-Control-Allow-Origin", "*")
	ctx.Output.Header("content-type", "application/javascript")
	ctx.Output.Body(body)
}

//输出ResponseResult结构
func ResponseJson(rr *ResponseResult, c *beego.Controller) {
	ResponseToContext(rr, c.Ctx)
	c.Finish()

}
