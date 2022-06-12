package main

import (
	"context"
	"fmt"
	"github.com/NNKulickov/forum/api"
	_ "github.com/NNKulickov/forum/docs"
	"github.com/fasthttp/router"
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/valyala/fasthttp"
	"io/ioutil"
	"log"
	"os"
	"time"
)

const initialScriptPath = "./db/db.sql"

func main() {
	fastRouter := router.New()
	api.DBS = initDB(context.Background(), initialScriptPath)
	api.InitRoutes(fastRouter.Group("/api"))
	log.Fatal(fasthttp.ListenAndServe("0.0.0.0:5000", wrapperHeader(fastRouter.Handler)).Error())
}

func wrapperHeader(handler fasthttp.RequestHandler) fasthttp.RequestHandler {
	return func(ctx *fasthttp.RequestCtx) {
		ctx.Response.Header.Set("Content-Type", "application/json")
		handler(ctx)
	}
}
func initDB(defaultCtx context.Context, initDBPath string) *pgxpool.Pool {
	connectString := fmt.Sprintf(
		"user=%s password=%s dbname=%s host=%s port=5432 sslmode=disable TimeZone=Europe/Moscow",
		os.Getenv("POSTGRES_USER"),
		os.Getenv("POSTGRES_PASSWORD"),
		os.Getenv("POSTGRES_DB"),
		os.Getenv("DB_HOST"),
	)
	fmt.Println(connectString)
	connectConf, err := pgxpool.ParseConfig(connectString)
	connectConf.MaxConns = 100
	connectConf.MaxConnLifetime = time.Minute
	connectConf.MaxConnIdleTime = time.Second * 5
	if err != nil {
		log.Fatal("cannot parse connect conf ", err)
	}
	fmt.Println("connectConf:", connectConf)

	pool, err := pgxpool.ConnectConfig(defaultCtx, connectConf)
	if err != nil {
		log.Fatal("cannot parse ", err)
	}
	pool.Ping(defaultCtx)
	if err != nil {
		log.Fatal("cannot connect ", err)
	}

	sql, err := ioutil.ReadFile(initDBPath)
	if err != nil {
		log.Fatal(err)
	}

	_, err = pool.Exec(defaultCtx, string(sql))
	if err != nil {
		log.Fatal(err)
	}

	return pool
}
