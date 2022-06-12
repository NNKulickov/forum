package response

import (
	"fmt"
	"github.com/valyala/fasthttp"
	"net/http"
)

type Data interface {
	MarshalJSON() ([]byte, error)
}

const (
	errMsg = "smthing wrong"
	okMsg  = "ok"
)

func Send[T Data](code int, payload T, ctx *fasthttp.RequestCtx) {

	ctx.SetStatusCode(code)

	if code == http.StatusInternalServerError {
		ctx.SetBody([]byte(errMsg))
		return
	}
	data, err := payload.MarshalJSON()
	if err != nil {
		fmt.Println("Can't send data")
		ctx.SetBody([]byte(errMsg))
		return
	}
	ctx.SetBody(data)
}
