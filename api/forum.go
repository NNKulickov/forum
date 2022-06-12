package api

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/NNKulickov/forum/forms"
	"github.com/NNKulickov/forum/response"
	"github.com/go-openapi/strfmt"
	"github.com/valyala/fasthttp"
	"net/http"
	"strconv"
	"strings"
)

func CreateForum(fastCtx *fasthttp.RequestCtx) {
	forum := forms.PostForum{}
	ctx := context.Background()

	if err := forum.UnmarshalJSON(fastCtx.Request.Body()); err != nil {
		fmt.Println("CreateForum (1):", err)
		response.Send(http.StatusInternalServerError, forms.Error{
			Message: " smth wrong",
		}, fastCtx)
		return
	}
	_, ok := checkForumSlug(ctx, forum.Slug)
	if ok {
		if res, err := getForum(ctx, forum.Slug); err == nil {
			response.Send(http.StatusConflict, res, fastCtx)
			return
		}
	}
	user, err := getUserByNicknam(ctx, forum.User)
	if err != nil {
		response.Send(http.StatusNotFound, forms.Error{
			Message: "none such user",
		}, fastCtx)

		return
	}
	_, err = DBS.Exec(ctx, `
		insert into forum (slug,title,host)
		values ($1,$2,$3)
	`,
		forum.Slug,
		forum.Title,
		user.Nickname,
	)
	if err == nil {
		response.Send(http.StatusCreated, forms.ForumResult{
			Title:   forum.Title,
			User:    user.Nickname,
			Slug:    forum.Slug,
			Posts:   0,
			Threads: 0,
		}, fastCtx)
		return
	}
	fmt.Println("CreateForum (2):", err)
	response.Send(http.StatusNotFound, forms.Error{
		Message: "none such user",
	}, fastCtx)
}

func GetForumDetails(fastCtx *fasthttp.RequestCtx) {
	forumParam := fastCtx.UserValue(forumSlug).(string)
	ctx := context.Background()
	fmt.Println("details:")
	forum, err := getForum(ctx, forumParam)
	if err != nil {
		fmt.Println("GetForumDetails (1):", err)
		response.Send(http.StatusNotFound, forms.Error{
			Message: "none such forum",
		}, fastCtx)
		return
	}
	response.Send(http.StatusOK, forum, fastCtx)
}

func CreateForumThread(fastCtx *fasthttp.RequestCtx) {
	slug := fastCtx.Value(forumSlug).(string)
	thread := forms.ThreadForm{}
	ctx := context.Background()
	if err := thread.UnmarshalJSON(fastCtx.Request.Body()); err != nil {
		fmt.Println("CreateForumThread (1):", err)
		return
	}
	var err error
	slug, ok := checkForumSlug(ctx, slug)
	if !ok {
		fmt.Println("CreateForumThread not found (3)", err)
		response.Send(http.StatusNotFound, forms.Error{
			Message: "none such forum"}, fastCtx)
		return
	}
	if err != nil {
		fmt.Println("CreateForumThread (2):", err)
		response.Send(http.StatusNotFound, forms.Error{
			Message: "none such user or forum",
		}, fastCtx)
		return
	}
	threadModel := forms.ThreadModel{
		Title:   thread.Title,
		Author:  thread.Author,
		Forum:   slug,
		Message: thread.Message,
		Slug:    sql.NullString{String: thread.Slug, Valid: true},
	}

	// try insert
	builder := strings.Builder{}
	builder.WriteString("insert into thread (title,author,forum,message,slug")
	if thread.Created != "" {
		builder.WriteString(",created")
	}
	builder.WriteString(") values ($1,$2,$3,$4,nullif($5,'')")
	if thread.Created != "" {

		builder.WriteString(fmt.Sprintf(",'%s'", thread.Created))
	}
	builder.WriteString(") returning id,created")
	if err = DBS.QueryRow(ctx, builder.String(),
		threadModel.Title,
		threadModel.Author,
		threadModel.Forum,
		threadModel.Message,
		threadModel.Slug,
	).
		Scan(
			&threadModel.Id,
			&threadModel.Created,
		); err == nil {
		response.Send(http.StatusCreated, forms.ThreadForm{
			Id:      threadModel.Id,
			Title:   threadModel.Title,
			Author:  threadModel.Author,
			Forum:   threadModel.Forum,
			Message: threadModel.Message,
			Slug:    threadModel.Slug.String,
			Votes:   threadModel.Votes,
			Created: strfmt.DateTime(threadModel.Created.UTC()).String(),
		}, fastCtx)

		return
	}
	// select if exists
	if err = DBS.QueryRow(ctx, `
		select
		id,
		title,
		author,
		forum,
		message,
		slug,
		created
		from thread 
		where lower(slug) = lower($1)`,
		threadModel.Slug.String,
	).
		Scan(
			&threadModel.Id,
			&threadModel.Title,
			&threadModel.Author,
			&threadModel.Forum,
			&threadModel.Message,
			&threadModel.Slug,
			&threadModel.Created,
		); err == nil {
		response.Send(http.StatusConflict, forms.ThreadForm{
			Id:      threadModel.Id,
			Title:   threadModel.Title,
			Author:  threadModel.Author,
			Forum:   threadModel.Forum,
			Message: threadModel.Message,
			Slug:    threadModel.Slug.String,
			Created: strfmt.DateTime(threadModel.Created.UTC()).String(),
		}, fastCtx)
		return
	}
	fmt.Println("CreateForumThread not found (3)", err)
	response.Send(http.StatusNotFound, forms.Error{
		Message: "none such user or forum"}, fastCtx)
}

func GetForumUsers(fastCtx *fasthttp.RequestCtx) {
	slug := fastCtx.UserValue(forumSlug).(string)
	ctx := context.Background()
	var err error
	slug, ok := checkForumSlug(ctx, slug)
	if !ok {
		fmt.Println("CreateForumThread not found (3)", err)
		response.Send(http.StatusNotFound, forms.Error{
			Message: "none such forum"}, fastCtx)
		return
	}
	limit := 0
	desc := false
	limitString := string(fastCtx.QueryArgs().Peek("limit"))
	descString := string(fastCtx.QueryArgs().Peek("desc"))
	if limit, err = strconv.Atoi(limitString); err != nil {
		fmt.Println("GetForumUsers (1):", err)
	}
	if desc, err = strconv.ParseBool(descString); err != nil {
		fmt.Println("GetForumUsers (2):", err)
	}
	users := forms.UserFilter{
		Limit: limit,
		Desc:  desc,
		Since: string(fastCtx.QueryArgs().Peek("since")),
	}
	build := strings.Builder{}

	build.WriteString(`
		select fa.nickname,a.fullname,a.about,a.email from forum_actors fa
		join actor a on lower(fa.nickname) = lower(a.nickname)
		where lower(fa.forum) = lower($1)`)
	if users.Since != "" {
		if users.Desc {
			build.WriteString(fmt.Sprintf(` and lower(fa.nickname) collate "C" <  lower('%s') collate "C"`, users.Since))

		} else {
			build.WriteString(fmt.Sprintf(` and lower(fa.nickname) collate "C" >  lower('%s') collate "C"`, users.Since))

		}
	}
	build.WriteString(` order by lower(fa.nickname) collate "C"`)
	if users.Desc {
		build.WriteString(" desc")
	}
	build.WriteString(" limit nullif($2,0)")
	if rows, err := DBS.Query(ctx, build.String(),
		slug, users.Limit); err == nil {
		usersResponse := new(forms.Users)
		for rows.Next() {
			user := forms.User{}
			if err = rows.
				Scan(
					&user.Nickname,
					&user.Fullname,
					&user.About,
					&user.Email,
				); err != nil {
				fmt.Println("GetForumUsers (2):", err)
				response.Send(http.StatusInternalServerError, forms.Error{
					Message: " smth wrong",
				}, fastCtx)
			}
			*usersResponse = append(*usersResponse, user)
		}
		if len(*usersResponse) == 0 {
			empty := forms.EmptyArray{}
			response.Send(http.StatusOK, empty, fastCtx)
			return
		}
		response.Send(http.StatusOK, usersResponse, fastCtx)
		return
	}
	response.Send(http.StatusNotFound, forms.Error{
		Message: "none such forum"}, fastCtx)
	return
}

func GetForumThreads(fastCtx *fasthttp.RequestCtx) {
	slug := fastCtx.UserValue(forumSlug).(string)
	ctx := context.Background()
	forum := ""
	if err := DBS.QueryRow(ctx,
		`select slug from forum where lower(slug) = lower($1)`,
		slug).Scan(&forum); err != nil {
		response.Send(http.StatusNotFound, forms.Error{
			Message: "none forum"}, fastCtx)
		return
	}
	slug = forum
	limit, err := strconv.Atoi(string(fastCtx.QueryArgs().Peek("limit")))
	if err != nil {
		fmt.Println("err:", err)
		limit = 0
	}
	desc, err := strconv.ParseBool(string(fastCtx.QueryArgs().Peek("desc")))
	if err != nil {
		fmt.Println("err:", err)
		desc = false
	}
	threadsFilter := forms.ThreadFilter{
		Limit: limit,
		Since: string(fastCtx.QueryArgs().Peek("since")),
		Desc:  desc,
	}
	build := strings.Builder{}
	build.WriteString(`
		select id, title, author, forum, message,
			slug, created,votes from thread
			where lower(forum) = lower($1)`)
	if threadsFilter.Since != "" {
		if desc {
			build.WriteString(fmt.
				Sprintf(
					` and created <= '%s'`,
					threadsFilter.Since,
				),
			)
		} else {
			build.WriteString(fmt.
				Sprintf(
					` and created >= '%s'`,
					threadsFilter.Since,
				),
			)
		}
	}
	build.WriteString(" order by created")
	if threadsFilter.Desc {
		build.WriteString(" desc")
	}
	build.WriteString(" limit nullif($2,0)")
	if rows, err := DBS.
		Query(
			ctx,
			build.String(),
			slug,
			threadsFilter.Limit); err == nil {
		threadsResponse := new(forms.ThreadsFrom)
		for rows.Next() {
			thread := forms.ThreadModel{}
			if err = rows.
				Scan(
					&thread.Id,
					&thread.Title,
					&thread.Author,
					&thread.Forum,
					&thread.Message,
					&thread.Slug,
					&thread.Created,
					&thread.Votes,
				); err != nil {
				fmt.Println("GetForumThreads (2):", err)
				response.Send(http.StatusInternalServerError, forms.Error{
					Message: " smth wrong",
				}, fastCtx)
				return
			}
			*threadsResponse = append(*threadsResponse, forms.ThreadForm{
				Id:      thread.Id,
				Title:   thread.Title,
				Author:  thread.Author,
				Forum:   thread.Forum,
				Message: thread.Message,
				Slug:    thread.Slug.String,
				Votes:   thread.Votes,
				Created: strfmt.DateTime(thread.Created.UTC()).String(),
			})
		}
		if len(*threadsResponse) == 0 {
			empty := forms.EmptyArray{}
			response.Send(http.StatusOK, empty, fastCtx)
			return
		}
		response.Send(http.StatusOK, threadsResponse, fastCtx)
		return
	}
	response.Send(http.StatusNotFound, forms.Error{
		Message: "none such forum"}, fastCtx)

	return
}
func checkForumSlug(ctx context.Context, slug string) (string, bool) {
	findSlug := ""
	err := DBS.QueryRow(ctx, `
	select slug from forum where lower(slug) = lower($1);
	`, slug).
		Scan(
			&findSlug,
		)
	if err != nil {
		return "", false
	}
	return findSlug, true
}
func getForum(ctx context.Context, slug string) (forms.ForumResult, error) {
	forum := forms.ForumResult{}

	err := DBS.QueryRow(ctx, `
	select slug,title,host,threads,posts from forum where lower(slug) = lower($1);
	`, slug).
		Scan(
			&forum.Slug,
			&forum.Title,
			&forum.User,
			&forum.Threads,
			&forum.Posts,
		)
	return forum, err
}
