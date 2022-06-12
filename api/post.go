package api

import (
	"context"
	"errors"
	"fmt"
	"github.com/NNKulickov/forum/forms"
	"github.com/NNKulickov/forum/response"
	"github.com/go-openapi/strfmt"
	"github.com/jackc/pgtype"
	"github.com/valyala/fasthttp"
	"net/http"
	"strconv"
	"strings"
)

func GetPostDetails(fastCtx *fasthttp.RequestCtx) {
	idStr := fastCtx.UserValue(postSlug).(string)
	id := 0
	var err error
	if id, err = strconv.Atoi(idStr); err != nil {
		fmt.Println("GetPostDetails (1):", err)
		return
	}
	related := string(fastCtx.QueryArgs().Peek("related"))
	isUser := strings.Contains(related, "user")
	isForum := strings.Contains(related, "forum")
	isThread := strings.Contains(related, "thread")

	ctx := context.Background()
	details := forms.PostDetails{}

	post, err := getSinglePost(ctx, id)
	if err != nil {
		fmt.Println("GetPostDetails (2):", err)
		response.Send(http.StatusNotFound, forms.Error{
			Message: err.Error(),
		}, fastCtx)
		return
	}
	details.Post = post
	if isUser {
		author, err := getUserByNicknam(ctx, post.Author)
		if err != nil {
			fmt.Println("GetPostDetails (3):", err)
			response.Send(http.StatusNotFound, forms.Error{
				Message: "Not found user",
			}, fastCtx)
			return
		}
		details.Author = &author
	}
	if isThread {
		threadModel := forms.ThreadModel{}
		if threadModel, err = getThreadById(ctx, post.Thread); err != nil {
			fmt.Println("GetPostDetails (4):", err)
			response.Send(http.StatusNotFound, forms.Error{
				Message: "Not found thread",
			}, fastCtx)
			return
		}
		details.Thread = &forms.ThreadForm{
			Id:      threadModel.Id,
			Title:   threadModel.Title,
			Author:  threadModel.Author,
			Forum:   threadModel.Forum,
			Message: threadModel.Message,
			Slug:    threadModel.Slug.String,
			Votes:   threadModel.Votes,
			Created: strfmt.DateTime(threadModel.Created.UTC()).String(),
		}
	}
	if isForum {
		forum, err := getForum(ctx, post.Forum)
		if err != nil {
			fmt.Println("GetPostDetails (5):", err)
			response.Send(http.StatusNotFound, forms.Error{
				Message: "Not found forum",
			}, fastCtx)
			return
		}
		details.Forum = &forum
	}
	response.Send(http.StatusOK, details, fastCtx)
}

func getSinglePost(ctx context.Context, id int) (forms.Post, error) {
	post := forms.Post{}
	created := pgtype.Timestamp{}
	if err := DBS.QueryRow(ctx, `
		select id,parent,author,message,isedited,forum,threadid,created from post
		where id = $1`, id).
		Scan(
			&post.Id,
			&post.Parent,
			&post.Author,
			&post.Message,
			&post.IsEdited,
			&post.Forum,
			&post.Thread,
			&created,
		); err != nil {
		fmt.Println("getSinglePost:", err)
		return post, errors.New("Not found post")
	}
	post.Created = strfmt.DateTime(created.Time.UTC()).String()
	return post, nil
}

func UpdatePostDetails(fastCtx *fasthttp.RequestCtx) {
	idStr := fastCtx.UserValue(postSlug).(string)
	id := 0
	var err error
	ctx := context.Background()
	postUpdate := new(forms.PostUpdate)
	if err = postUpdate.UnmarshalJSON(fastCtx.Request.Body()); err != nil {
		fmt.Println("UpdatePostDetails (1):", err)
		return
	}

	if id, err = strconv.Atoi(idStr); err != nil {
		fmt.Println("UpdatePostDetails (2):", err)
		return
	}
	post := forms.Post{}
	post, err = getSinglePost(ctx, id)
	if err != nil {
		fmt.Println("UpdatePostDetails (3):", err)
		response.Send(http.StatusNotFound, forms.Error{
			Message: err.Error(),
		}, fastCtx)
		return
	}

	if postUpdate.Message == "" || postUpdate.Message == post.Message {
		response.Send(http.StatusOK, post, fastCtx)
		return
	}
	created := pgtype.Timestamp{}
	if err = DBS.QueryRow(ctx, `
		update post set message = $1,isedited = true  where id = $2
		returning id,parent,author,message,isedited,forum,threadid,created
		`, postUpdate.Message, id).
		Scan(
			&post.Id,
			&post.Parent,
			&post.Author,
			&post.Message,
			&post.IsEdited,
			&post.Forum,
			&post.Thread,
			&created,
		); err != nil {
		fmt.Println("UpdatePostDetails (4):", err)
		response.Send(http.StatusNotFound, forms.Error{
			Message: "Not found post",
		}, fastCtx)
		return
	}
	post.Created = strfmt.DateTime(created.Time.UTC()).String()
	response.Send(http.StatusOK, post, fastCtx)
}

func getPostsFlat(ctx context.Context, threadid, limit, since int, desc bool) ([]forms.Post, error) {
	posts := []forms.Post{}
	builder := strings.Builder{}

	builder.WriteString(`
		select id,parent,author,message,isedited,forum,threadid,created
			from post where threadid = $1`)
	if since > 0 {
		if desc {
			builder.WriteString(fmt.Sprintf(" and id < %d", since))

		} else {
			builder.WriteString(fmt.Sprintf(" and id > %d", since))

		}
	}

	builder.WriteString(" order by id")

	if desc {
		builder.WriteString(" desc")
	}
	builder.WriteString(",created limit nullif($2,0)")
	rows, err := DBS.Query(ctx, builder.String(),
		threadid, limit)
	if err != nil {
		fmt.Println("getPostsFlat(1): ", err)
		return nil, err
	}
	for rows.Next() {
		post := forms.Post{}
		created := pgtype.Timestamp{}
		err = rows.Scan(
			&post.Id,
			&post.Parent,
			&post.Author,
			&post.Message,
			&post.IsEdited,
			&post.Forum,
			&post.Thread,
			&created,
		)
		if err != nil {
			fmt.Println("getPostsFlat(2): ", err, post)
			return nil, err
		}
		post.Created = strfmt.DateTime(created.Time.UTC()).String()
		posts = append(posts, post)
	}
	return posts, nil
}
func getPostsParentTree(ctx context.Context, threadid, limit, since int, desc bool) ([]forms.Post, error) {
	posts := []forms.Post{}
	builder := strings.Builder{}
	subBuilder := strings.Builder{}
	subBuilder.WriteString(`
		select id from post where threadid = $1 and parent = 0`)
	if since > 0 {
		subBuilder.WriteString(" and pathtree[1]")
		if desc {
			subBuilder.WriteString(" <")
		} else {
			subBuilder.WriteString(" >")
		}
		subBuilder.WriteString(fmt.Sprintf(" (select pathtree[1] from post where id = %d)", since))

	}

	subBuilder.WriteString(" order by id")

	if desc {
		subBuilder.WriteString(" desc")
	}
	subBuilder.WriteString(" limit nullif($2,0)")
	builder.WriteString(fmt.Sprintf(`
				select id,parent,author,message,isedited,forum,threadid,created
			from post where pathtree[1] in (%s)
		 `, subBuilder.String()))
	if desc {
		builder.WriteString(" order by pathtree[1] desc,pathtree")
	} else {
		builder.WriteString(" order by pathtree")

	}
	fmt.Println("getPostsParentTree:", builder.String())
	rows, err := DBS.Query(ctx, builder.String(),
		threadid, limit)
	if err != nil {
		fmt.Println("getPostsParentTree(1): ", err)
		return nil, err
	}
	for rows.Next() {
		post := forms.Post{}
		created := pgtype.Timestamp{}
		err = rows.Scan(
			&post.Id,
			&post.Parent,
			&post.Author,
			&post.Message,
			&post.IsEdited,
			&post.Forum,
			&post.Thread,
			&created,
		)
		if err != nil {
			fmt.Println("getPostsParentTree(2): ", err, post)
			return nil, err
		}
		post.Created = strfmt.DateTime(created.Time.UTC()).String()
		posts = append(posts, post)
	}
	return posts, nil
}

func getPostsTree(ctx context.Context, threadid, limit, since int, desc bool) ([]forms.Post, error) {
	posts := []forms.Post{}
	builder := strings.Builder{}

	builder.WriteString(`
		select id,parent,author,message,isedited,forum,threadid,created
			from post where threadid = $1`)
	if since > 0 {
		builder.WriteString(" and pathtree")
		if desc {
			builder.WriteString(" <")
		} else {
			builder.WriteString(" >")
		}
		builder.WriteString(fmt.Sprintf(" (select pathtree from post where id = %d)", since))
	}

	builder.WriteString(" order by pathtree ")

	if desc {
		builder.WriteString("desc")
	}
	builder.WriteString(" limit nullif($2,0)")
	rows, err := DBS.Query(ctx, builder.String(),
		threadid, limit)
	if err != nil {
		fmt.Println("getPostsTree(1): ", err)
		return nil, err
	}
	for rows.Next() {
		post := forms.Post{}
		created := pgtype.Timestamp{}
		err = rows.Scan(
			&post.Id,
			&post.Parent,
			&post.Author,
			&post.Message,
			&post.IsEdited,
			&post.Forum,
			&post.Thread,
			&created,
		)
		if err != nil {
			fmt.Println("getPostsTree(2): ", err, post)
			return nil, err
		}
		post.Created = strfmt.DateTime(created.Time.UTC()).String()
		posts = append(posts, post)
	}
	return posts, nil
}
