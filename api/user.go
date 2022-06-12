package api

import (
	"context"
	"errors"
	"fmt"
	"github.com/NNKulickov/forum/forms"
	"github.com/NNKulickov/forum/response"
	"github.com/jackc/pgx/v4"
	"github.com/valyala/fasthttp"
	"net/http"
)

func CreateUser(fastCtx *fasthttp.RequestCtx) {
	ctx := context.Background()
	nickname, err := getSlugUsername(fastCtx)
	if err != nil {
		fmt.Println("CreateUser (1):", err)
		return
	}
	user := new(forms.User)
	if err := user.UnmarshalJSON(fastCtx.Request.Body()); err != nil {
		fmt.Println("CreateUser (2):", err)
		return
	}
	user.Nickname = nickname
	_, err = DBS.Exec(ctx, `
		insert into 
		    actor (nickname,fullname,about,email) 
		values ($1,$2,$3,$4)
		`,
		user.Nickname, user.Fullname, user.About, user.Email)
	if err == nil {
		response.Send(http.StatusCreated, user, fastCtx)
		return
	}
	rows, err := DBS.Query(ctx, `
		select nickname,fullname,about,email 
		from actor 
		where  lower(nickname) = lower($1) or lower(email) = lower($2)
		`, user.Nickname, user.Email)
	defer rows.Close()
	users := new(forms.Users)
	for rows.Next() {
		rowUser := forms.User{}
		if err = rows.Scan(
			&rowUser.Nickname,
			&rowUser.Fullname,
			&rowUser.About,
			&rowUser.Email,
		); err != nil {
			response.Send(http.StatusInternalServerError, forms.Error{Message: "Cannot get user" + err.Error()}, fastCtx)
			return
		}
		*users = append(*users, rowUser)
	}
	response.Send(http.StatusConflict, users, fastCtx)
}

func GetUserProfile(fastCtx *fasthttp.RequestCtx) {
	nickname, err := getSlugUsername(fastCtx)
	ctx := context.Background()
	user, err := getUserByNicknam(ctx, nickname)
	if err != nil {
		response.Send(http.StatusNotFound, forms.Error{Message: "Not found"}, fastCtx)
		return
	}
	fmt.Println("user:", user)
	response.Send(http.StatusOK, user, fastCtx)
}

func getUserByNicknam(ctx context.Context, nickname string) (forms.User, error) {
	user := forms.User{}
	err := DBS.QueryRow(ctx, ` 
			select nickname,fullname,about,email 
				from actor
			where  lower(nickname) = lower($1)`, nickname).
		Scan(&user.Nickname, &user.Fullname, &user.About, &user.Email)

	return user, err
}

func UpdateUserProfile(fastCtx *fasthttp.RequestCtx) {
	nickname, err := getSlugUsername(fastCtx)
	if err != nil {
		response.Send(http.StatusInternalServerError, forms.Error{
			Message: " smth wrong",
		}, fastCtx)
		return
	}
	user := new(forms.User)
	if err := user.UnmarshalJSON(fastCtx.Request.Body()); err != nil {
		fmt.Println("UpdateUserProfile:", err)
		return
	}
	user.Nickname = nickname
	userModel := forms.User{}
	ctx := context.Background()
	if err = DBS.QueryRow(ctx,
		`select nickname,fullname,about,email from actor 
                where lower(nickname) = lower($1)
		`, user.Nickname).Scan(
		&userModel.Nickname,
		&userModel.Fullname,
		&userModel.About,
		&userModel.Email,
	); err == pgx.ErrNoRows {
		response.Send(http.StatusNotFound, forms.Error{Message: "none such user"}, fastCtx)
		return
	}
	if user.About == "" {
		user.About = userModel.About
	}
	if user.Email == "" {
		user.Email = userModel.Email
	}
	if user.Fullname == "" {
		user.Fullname = userModel.Fullname
	}
	if err = DBS.QueryRow(ctx, `
		update actor 
		set fullname = $2,
		    about = $3,
		    email = $4
		where lower(nickname) = lower($1)
		returning nickname,fullname,about,email
		`,
		user.Nickname, user.Fullname, user.About, user.Email).
		Scan(
			&user.Nickname,
			&user.Fullname,
			&user.About,
			&user.Email,
		); err != nil {

		response.Send(http.StatusConflict, forms.Error{Message: "new params don't suit"}, fastCtx)
		return
	}
	response.Send(200, user, fastCtx)
	return
}

func getSlugUsername(fastCtx *fasthttp.RequestCtx) (string, error) {
	username := fastCtx.UserValue(usernameSlug).(string)
	if username == "" {
		fmt.Println("cannot get username")
		return username, errors.New("None user")
	}
	return username, nil
}
