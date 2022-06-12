// Code generated by easyjson for marshaling/unmarshaling. DO NOT EDIT.

package forms

import (
	json "encoding/json"
	easyjson "github.com/mailru/easyjson"
	jlexer "github.com/mailru/easyjson/jlexer"
	jwriter "github.com/mailru/easyjson/jwriter"
)

// suppress unused package warning
var (
	_ *json.RawMessage
	_ *jlexer.Lexer
	_ *jwriter.Writer
	_ easyjson.Marshaler
)

func easyjson5a72dc82DecodeGithubComNNKulickovForumForms(in *jlexer.Lexer, out *ThreadPosts) {
	isTopLevel := in.IsStart()
	if in.IsNull() {
		if isTopLevel {
			in.Consumed()
		}
		in.Skip()
		return
	}
	in.Delim('{')
	for !in.IsDelim('}') {
		key := in.UnsafeFieldName(false)
		in.WantColon()
		if in.IsNull() {
			in.Skip()
			in.WantComma()
			continue
		}
		switch key {
		case "limit":
			out.Limit = int(in.Int())
		case "since":
			out.Since = int(in.Int())
		case "sort":
			out.Sort = string(in.String())
		case "desc":
			out.Desc = bool(in.Bool())
		default:
			in.SkipRecursive()
		}
		in.WantComma()
	}
	in.Delim('}')
	if isTopLevel {
		in.Consumed()
	}
}
func easyjson5a72dc82EncodeGithubComNNKulickovForumForms(out *jwriter.Writer, in ThreadPosts) {
	out.RawByte('{')
	first := true
	_ = first
	if in.Limit != 0 {
		const prefix string = ",\"limit\":"
		first = false
		out.RawString(prefix[1:])
		out.Int(int(in.Limit))
	}
	if in.Since != 0 {
		const prefix string = ",\"since\":"
		if first {
			first = false
			out.RawString(prefix[1:])
		} else {
			out.RawString(prefix)
		}
		out.Int(int(in.Since))
	}
	if in.Sort != "" {
		const prefix string = ",\"sort\":"
		if first {
			first = false
			out.RawString(prefix[1:])
		} else {
			out.RawString(prefix)
		}
		out.String(string(in.Sort))
	}
	{
		const prefix string = ",\"desc\":"
		if first {
			first = false
			out.RawString(prefix[1:])
		} else {
			out.RawString(prefix)
		}
		out.Bool(bool(in.Desc))
	}
	out.RawByte('}')
}

// MarshalJSON supports json.Marshaler interface
func (v ThreadPosts) MarshalJSON() ([]byte, error) {
	w := jwriter.Writer{}
	easyjson5a72dc82EncodeGithubComNNKulickovForumForms(&w, v)
	return w.Buffer.BuildBytes(), w.Error
}

// MarshalEasyJSON supports easyjson.Marshaler interface
func (v ThreadPosts) MarshalEasyJSON(w *jwriter.Writer) {
	easyjson5a72dc82EncodeGithubComNNKulickovForumForms(w, v)
}

// UnmarshalJSON supports json.Unmarshaler interface
func (v *ThreadPosts) UnmarshalJSON(data []byte) error {
	r := jlexer.Lexer{Data: data}
	easyjson5a72dc82DecodeGithubComNNKulickovForumForms(&r, v)
	return r.Error()
}

// UnmarshalEasyJSON supports easyjson.Unmarshaler interface
func (v *ThreadPosts) UnmarshalEasyJSON(l *jlexer.Lexer) {
	easyjson5a72dc82DecodeGithubComNNKulickovForumForms(l, v)
}
func easyjson5a72dc82DecodeGithubComNNKulickovForumForms1(in *jlexer.Lexer, out *Posts) {
	isTopLevel := in.IsStart()
	if in.IsNull() {
		in.Skip()
		*out = nil
	} else {
		in.Delim('[')
		if *out == nil {
			if !in.IsDelim(']') {
				*out = make(Posts, 0, 0)
			} else {
				*out = Posts{}
			}
		} else {
			*out = (*out)[:0]
		}
		for !in.IsDelim(']') {
			var v1 Post
			(v1).UnmarshalEasyJSON(in)
			*out = append(*out, v1)
			in.WantComma()
		}
		in.Delim(']')
	}
	if isTopLevel {
		in.Consumed()
	}
}
func easyjson5a72dc82EncodeGithubComNNKulickovForumForms1(out *jwriter.Writer, in Posts) {
	if in == nil && (out.Flags&jwriter.NilSliceAsEmpty) == 0 {
		out.RawString("null")
	} else {
		out.RawByte('[')
		for v2, v3 := range in {
			if v2 > 0 {
				out.RawByte(',')
			}
			(v3).MarshalEasyJSON(out)
		}
		out.RawByte(']')
	}
}

// MarshalJSON supports json.Marshaler interface
func (v Posts) MarshalJSON() ([]byte, error) {
	w := jwriter.Writer{}
	easyjson5a72dc82EncodeGithubComNNKulickovForumForms1(&w, v)
	return w.Buffer.BuildBytes(), w.Error
}

// MarshalEasyJSON supports easyjson.Marshaler interface
func (v Posts) MarshalEasyJSON(w *jwriter.Writer) {
	easyjson5a72dc82EncodeGithubComNNKulickovForumForms1(w, v)
}

// UnmarshalJSON supports json.Unmarshaler interface
func (v *Posts) UnmarshalJSON(data []byte) error {
	r := jlexer.Lexer{Data: data}
	easyjson5a72dc82DecodeGithubComNNKulickovForumForms1(&r, v)
	return r.Error()
}

// UnmarshalEasyJSON supports easyjson.Unmarshaler interface
func (v *Posts) UnmarshalEasyJSON(l *jlexer.Lexer) {
	easyjson5a72dc82DecodeGithubComNNKulickovForumForms1(l, v)
}
func easyjson5a72dc82DecodeGithubComNNKulickovForumForms2(in *jlexer.Lexer, out *PostUpdate) {
	isTopLevel := in.IsStart()
	if in.IsNull() {
		if isTopLevel {
			in.Consumed()
		}
		in.Skip()
		return
	}
	in.Delim('{')
	for !in.IsDelim('}') {
		key := in.UnsafeFieldName(false)
		in.WantColon()
		if in.IsNull() {
			in.Skip()
			in.WantComma()
			continue
		}
		switch key {
		case "message":
			out.Message = string(in.String())
		default:
			in.SkipRecursive()
		}
		in.WantComma()
	}
	in.Delim('}')
	if isTopLevel {
		in.Consumed()
	}
}
func easyjson5a72dc82EncodeGithubComNNKulickovForumForms2(out *jwriter.Writer, in PostUpdate) {
	out.RawByte('{')
	first := true
	_ = first
	{
		const prefix string = ",\"message\":"
		out.RawString(prefix[1:])
		out.String(string(in.Message))
	}
	out.RawByte('}')
}

// MarshalJSON supports json.Marshaler interface
func (v PostUpdate) MarshalJSON() ([]byte, error) {
	w := jwriter.Writer{}
	easyjson5a72dc82EncodeGithubComNNKulickovForumForms2(&w, v)
	return w.Buffer.BuildBytes(), w.Error
}

// MarshalEasyJSON supports easyjson.Marshaler interface
func (v PostUpdate) MarshalEasyJSON(w *jwriter.Writer) {
	easyjson5a72dc82EncodeGithubComNNKulickovForumForms2(w, v)
}

// UnmarshalJSON supports json.Unmarshaler interface
func (v *PostUpdate) UnmarshalJSON(data []byte) error {
	r := jlexer.Lexer{Data: data}
	easyjson5a72dc82DecodeGithubComNNKulickovForumForms2(&r, v)
	return r.Error()
}

// UnmarshalEasyJSON supports easyjson.Unmarshaler interface
func (v *PostUpdate) UnmarshalEasyJSON(l *jlexer.Lexer) {
	easyjson5a72dc82DecodeGithubComNNKulickovForumForms2(l, v)
}
func easyjson5a72dc82DecodeGithubComNNKulickovForumForms3(in *jlexer.Lexer, out *PostFull) {
	isTopLevel := in.IsStart()
	if in.IsNull() {
		if isTopLevel {
			in.Consumed()
		}
		in.Skip()
		return
	}
	in.Delim('{')
	for !in.IsDelim('}') {
		key := in.UnsafeFieldName(false)
		in.WantColon()
		if in.IsNull() {
			in.Skip()
			in.WantComma()
			continue
		}
		switch key {
		case "post":
			(out.Post).UnmarshalEasyJSON(in)
		case "author":
			(out.Author).UnmarshalEasyJSON(in)
		case "thread":
			(out.Thread).UnmarshalEasyJSON(in)
		default:
			in.SkipRecursive()
		}
		in.WantComma()
	}
	in.Delim('}')
	if isTopLevel {
		in.Consumed()
	}
}
func easyjson5a72dc82EncodeGithubComNNKulickovForumForms3(out *jwriter.Writer, in PostFull) {
	out.RawByte('{')
	first := true
	_ = first
	{
		const prefix string = ",\"post\":"
		out.RawString(prefix[1:])
		(in.Post).MarshalEasyJSON(out)
	}
	{
		const prefix string = ",\"author\":"
		out.RawString(prefix)
		(in.Author).MarshalEasyJSON(out)
	}
	{
		const prefix string = ",\"thread\":"
		out.RawString(prefix)
		(in.Thread).MarshalEasyJSON(out)
	}
	out.RawByte('}')
}

// MarshalJSON supports json.Marshaler interface
func (v PostFull) MarshalJSON() ([]byte, error) {
	w := jwriter.Writer{}
	easyjson5a72dc82EncodeGithubComNNKulickovForumForms3(&w, v)
	return w.Buffer.BuildBytes(), w.Error
}

// MarshalEasyJSON supports easyjson.Marshaler interface
func (v PostFull) MarshalEasyJSON(w *jwriter.Writer) {
	easyjson5a72dc82EncodeGithubComNNKulickovForumForms3(w, v)
}

// UnmarshalJSON supports json.Unmarshaler interface
func (v *PostFull) UnmarshalJSON(data []byte) error {
	r := jlexer.Lexer{Data: data}
	easyjson5a72dc82DecodeGithubComNNKulickovForumForms3(&r, v)
	return r.Error()
}

// UnmarshalEasyJSON supports easyjson.Unmarshaler interface
func (v *PostFull) UnmarshalEasyJSON(l *jlexer.Lexer) {
	easyjson5a72dc82DecodeGithubComNNKulickovForumForms3(l, v)
}
func easyjson5a72dc82DecodeGithubComNNKulickovForumForms4(in *jlexer.Lexer, out *PostDetails) {
	isTopLevel := in.IsStart()
	if in.IsNull() {
		if isTopLevel {
			in.Consumed()
		}
		in.Skip()
		return
	}
	in.Delim('{')
	for !in.IsDelim('}') {
		key := in.UnsafeFieldName(false)
		in.WantColon()
		if in.IsNull() {
			in.Skip()
			in.WantComma()
			continue
		}
		switch key {
		case "post":
			(out.Post).UnmarshalEasyJSON(in)
		case "author":
			if in.IsNull() {
				in.Skip()
				out.Author = nil
			} else {
				if out.Author == nil {
					out.Author = new(User)
				}
				(*out.Author).UnmarshalEasyJSON(in)
			}
		case "thread":
			if in.IsNull() {
				in.Skip()
				out.Thread = nil
			} else {
				if out.Thread == nil {
					out.Thread = new(ThreadForm)
				}
				(*out.Thread).UnmarshalEasyJSON(in)
			}
		case "forum":
			if in.IsNull() {
				in.Skip()
				out.Forum = nil
			} else {
				if out.Forum == nil {
					out.Forum = new(ForumResult)
				}
				(*out.Forum).UnmarshalEasyJSON(in)
			}
		default:
			in.SkipRecursive()
		}
		in.WantComma()
	}
	in.Delim('}')
	if isTopLevel {
		in.Consumed()
	}
}
func easyjson5a72dc82EncodeGithubComNNKulickovForumForms4(out *jwriter.Writer, in PostDetails) {
	out.RawByte('{')
	first := true
	_ = first
	{
		const prefix string = ",\"post\":"
		out.RawString(prefix[1:])
		(in.Post).MarshalEasyJSON(out)
	}
	if in.Author != nil {
		const prefix string = ",\"author\":"
		out.RawString(prefix)
		(*in.Author).MarshalEasyJSON(out)
	}
	if in.Thread != nil {
		const prefix string = ",\"thread\":"
		out.RawString(prefix)
		(*in.Thread).MarshalEasyJSON(out)
	}
	if in.Forum != nil {
		const prefix string = ",\"forum\":"
		out.RawString(prefix)
		(*in.Forum).MarshalEasyJSON(out)
	}
	out.RawByte('}')
}

// MarshalJSON supports json.Marshaler interface
func (v PostDetails) MarshalJSON() ([]byte, error) {
	w := jwriter.Writer{}
	easyjson5a72dc82EncodeGithubComNNKulickovForumForms4(&w, v)
	return w.Buffer.BuildBytes(), w.Error
}

// MarshalEasyJSON supports easyjson.Marshaler interface
func (v PostDetails) MarshalEasyJSON(w *jwriter.Writer) {
	easyjson5a72dc82EncodeGithubComNNKulickovForumForms4(w, v)
}

// UnmarshalJSON supports json.Unmarshaler interface
func (v *PostDetails) UnmarshalJSON(data []byte) error {
	r := jlexer.Lexer{Data: data}
	easyjson5a72dc82DecodeGithubComNNKulickovForumForms4(&r, v)
	return r.Error()
}

// UnmarshalEasyJSON supports easyjson.Unmarshaler interface
func (v *PostDetails) UnmarshalEasyJSON(l *jlexer.Lexer) {
	easyjson5a72dc82DecodeGithubComNNKulickovForumForms4(l, v)
}
func easyjson5a72dc82DecodeGithubComNNKulickovForumForms5(in *jlexer.Lexer, out *Post) {
	isTopLevel := in.IsStart()
	if in.IsNull() {
		if isTopLevel {
			in.Consumed()
		}
		in.Skip()
		return
	}
	in.Delim('{')
	for !in.IsDelim('}') {
		key := in.UnsafeFieldName(false)
		in.WantColon()
		if in.IsNull() {
			in.Skip()
			in.WantComma()
			continue
		}
		switch key {
		case "id":
			out.Id = int(in.Int())
		case "parent":
			out.Parent = int(in.Int())
		case "author":
			out.Author = string(in.String())
		case "message":
			out.Message = string(in.String())
		case "isEdited":
			out.IsEdited = bool(in.Bool())
		case "forum":
			out.Forum = string(in.String())
		case "thread":
			out.Thread = int(in.Int())
		case "created":
			out.Created = string(in.String())
		default:
			in.SkipRecursive()
		}
		in.WantComma()
	}
	in.Delim('}')
	if isTopLevel {
		in.Consumed()
	}
}
func easyjson5a72dc82EncodeGithubComNNKulickovForumForms5(out *jwriter.Writer, in Post) {
	out.RawByte('{')
	first := true
	_ = first
	if in.Id != 0 {
		const prefix string = ",\"id\":"
		first = false
		out.RawString(prefix[1:])
		out.Int(int(in.Id))
	}
	if in.Parent != 0 {
		const prefix string = ",\"parent\":"
		if first {
			first = false
			out.RawString(prefix[1:])
		} else {
			out.RawString(prefix)
		}
		out.Int(int(in.Parent))
	}
	{
		const prefix string = ",\"author\":"
		if first {
			first = false
			out.RawString(prefix[1:])
		} else {
			out.RawString(prefix)
		}
		out.String(string(in.Author))
	}
	{
		const prefix string = ",\"message\":"
		out.RawString(prefix)
		out.String(string(in.Message))
	}
	if in.IsEdited {
		const prefix string = ",\"isEdited\":"
		out.RawString(prefix)
		out.Bool(bool(in.IsEdited))
	}
	if in.Forum != "" {
		const prefix string = ",\"forum\":"
		out.RawString(prefix)
		out.String(string(in.Forum))
	}
	{
		const prefix string = ",\"thread\":"
		out.RawString(prefix)
		out.Int(int(in.Thread))
	}
	if in.Created != "" {
		const prefix string = ",\"created\":"
		out.RawString(prefix)
		out.String(string(in.Created))
	}
	out.RawByte('}')
}

// MarshalJSON supports json.Marshaler interface
func (v Post) MarshalJSON() ([]byte, error) {
	w := jwriter.Writer{}
	easyjson5a72dc82EncodeGithubComNNKulickovForumForms5(&w, v)
	return w.Buffer.BuildBytes(), w.Error
}

// MarshalEasyJSON supports easyjson.Marshaler interface
func (v Post) MarshalEasyJSON(w *jwriter.Writer) {
	easyjson5a72dc82EncodeGithubComNNKulickovForumForms5(w, v)
}

// UnmarshalJSON supports json.Unmarshaler interface
func (v *Post) UnmarshalJSON(data []byte) error {
	r := jlexer.Lexer{Data: data}
	easyjson5a72dc82DecodeGithubComNNKulickovForumForms5(&r, v)
	return r.Error()
}

// UnmarshalEasyJSON supports easyjson.Unmarshaler interface
func (v *Post) UnmarshalEasyJSON(l *jlexer.Lexer) {
	easyjson5a72dc82DecodeGithubComNNKulickovForumForms5(l, v)
}
