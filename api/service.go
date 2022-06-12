package api

import (
	"context"
	"fmt"
	"github.com/NNKulickov/forum/forms"
	"github.com/NNKulickov/forum/response"
	"github.com/valyala/fasthttp"
	"io/ioutil"
	"log"
	"net/http"
)

const initialScriptPath = "./db/db.sql"

func GetServiceStatus(fastCtx *fasthttp.RequestCtx) {
	ctx := context.Background()
	status := forms.Status{}
	if err := DBS.QueryRow(ctx, `select count(*) from forum`).
		Scan(&status.Forum); err != nil {
		fmt.Println("GetServiceStatus (1) :", err)
		return
	}

	if err := DBS.QueryRow(ctx, `select count(*) from post`).
		Scan(&status.Post); err != nil {
		fmt.Println("GetServiceStatus (2) :", err)
		return
	}

	if err := DBS.QueryRow(ctx, `select count(*) from thread`).
		Scan(&status.Thread); err != nil {
		fmt.Println("GetServiceStatus (3) :", err)
		return
	}

	if err := DBS.QueryRow(ctx, `select count(*) from actor`).
		Scan(&status.User); err != nil {
		fmt.Println("GetServiceStatus (4) :", err)
		return
	}
	response.Send(http.StatusOK, status, fastCtx)
}

func ClearServiceData(fastCtx *fasthttp.RequestCtx) {
	ctx := context.Background()

	_, err := DBS.Exec(ctx, `
		truncate actor,forum,post,thread,vote,forum_actors
		`)
	sql, err := ioutil.ReadFile(initialScriptPath)
	if err != nil {
		log.Fatal(err)
	}

	_, err = DBS.Exec(ctx, string(sql))
	if err != nil {
		log.Fatal(err)
	}
	response.Send(http.StatusOK, forms.Error{
		Message: " smth wrong",
	}, fastCtx)
}
