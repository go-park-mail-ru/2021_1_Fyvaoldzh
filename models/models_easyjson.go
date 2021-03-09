// Code generated by easyjson for marshaling/unmarshaling. DO NOT EDIT.

package models

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

func easyjsonD2b7633eDecodeGithubComGoParkMailRu20211FyvaoldzhModels(in *jlexer.Lexer, out *UserProfile) {
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
		case "Uid":
			out.Uid = int(in.Int())
		case "name":
			out.Name = string(in.String())
		case "age":
			out.Age = int(in.Int())
		case "city":
			out.City = string(in.String())
		case "followers":
			out.Followers = int(in.Int())
		case "about":
			out.About = string(in.String())
		case "avatar":
			out.Avatar = string(in.String())
		case "events":
			(out.Event).UnmarshalEasyJSON(in)
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
func easyjsonD2b7633eEncodeGithubComGoParkMailRu20211FyvaoldzhModels(out *jwriter.Writer, in UserProfile) {
	out.RawByte('{')
	first := true
	_ = first
	{
		const prefix string = ",\"Uid\":"
		out.RawString(prefix[1:])
		out.Int(int(in.Uid))
	}
	{
		const prefix string = ",\"name\":"
		out.RawString(prefix)
		out.String(string(in.Name))
	}
	{
		const prefix string = ",\"age\":"
		out.RawString(prefix)
		out.Int(int(in.Age))
	}
	{
		const prefix string = ",\"city\":"
		out.RawString(prefix)
		out.String(string(in.City))
	}
	{
		const prefix string = ",\"followers\":"
		out.RawString(prefix)
		out.Int(int(in.Followers))
	}
	{
		const prefix string = ",\"about\":"
		out.RawString(prefix)
		out.String(string(in.About))
	}
	{
		const prefix string = ",\"avatar\":"
		out.RawString(prefix)
		out.String(string(in.Avatar))
	}
	{
		const prefix string = ",\"events\":"
		out.RawString(prefix)
		(in.Event).MarshalEasyJSON(out)
	}
	out.RawByte('}')
}

// MarshalJSON supports json.Marshaler interface
func (v UserProfile) MarshalJSON() ([]byte, error) {
	w := jwriter.Writer{}
	easyjsonD2b7633eEncodeGithubComGoParkMailRu20211FyvaoldzhModels(&w, v)
	return w.Buffer.BuildBytes(), w.Error
}

// MarshalEasyJSON supports easyjson.Marshaler interface
func (v UserProfile) MarshalEasyJSON(w *jwriter.Writer) {
	easyjsonD2b7633eEncodeGithubComGoParkMailRu20211FyvaoldzhModels(w, v)
}

// UnmarshalJSON supports json.Unmarshaler interface
func (v *UserProfile) UnmarshalJSON(data []byte) error {
	r := jlexer.Lexer{Data: data}
	easyjsonD2b7633eDecodeGithubComGoParkMailRu20211FyvaoldzhModels(&r, v)
	return r.Error()
}

// UnmarshalEasyJSON supports easyjson.Unmarshaler interface
func (v *UserProfile) UnmarshalEasyJSON(l *jlexer.Lexer) {
	easyjsonD2b7633eDecodeGithubComGoParkMailRu20211FyvaoldzhModels(l, v)
}
func easyjsonD2b7633eDecodeGithubComGoParkMailRu20211FyvaoldzhModels1(in *jlexer.Lexer, out *UserOwnProfile) {
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
		case "Uid":
			out.Uid = int(in.Int())
		case "name":
			out.Name = string(in.String())
		case "birthday":
			out.Birthday = string(in.String())
		case "city":
			out.City = string(in.String())
		case "email":
			out.Email = string(in.String())
		case "visited":
			out.Visited = int(in.Int())
		case "planning":
			out.Planning = int(in.Int())
		case "followers":
			out.Followers = int(in.Int())
		case "about":
			out.About = string(in.String())
		case "avatar":
			out.Avatar = string(in.String())
		case "events":
			(out.Event).UnmarshalEasyJSON(in)
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
func easyjsonD2b7633eEncodeGithubComGoParkMailRu20211FyvaoldzhModels1(out *jwriter.Writer, in UserOwnProfile) {
	out.RawByte('{')
	first := true
	_ = first
	{
		const prefix string = ",\"Uid\":"
		out.RawString(prefix[1:])
		out.Int(int(in.Uid))
	}
	{
		const prefix string = ",\"name\":"
		out.RawString(prefix)
		out.String(string(in.Name))
	}
	{
		const prefix string = ",\"birthday\":"
		out.RawString(prefix)
		out.String(string(in.Birthday))
	}
	{
		const prefix string = ",\"city\":"
		out.RawString(prefix)
		out.String(string(in.City))
	}
	{
		const prefix string = ",\"email\":"
		out.RawString(prefix)
		out.String(string(in.Email))
	}
	{
		const prefix string = ",\"visited\":"
		out.RawString(prefix)
		out.Int(int(in.Visited))
	}
	{
		const prefix string = ",\"planning\":"
		out.RawString(prefix)
		out.Int(int(in.Planning))
	}
	{
		const prefix string = ",\"followers\":"
		out.RawString(prefix)
		out.Int(int(in.Followers))
	}
	{
		const prefix string = ",\"about\":"
		out.RawString(prefix)
		out.String(string(in.About))
	}
	{
		const prefix string = ",\"avatar\":"
		out.RawString(prefix)
		out.String(string(in.Avatar))
	}
	{
		const prefix string = ",\"events\":"
		out.RawString(prefix)
		(in.Event).MarshalEasyJSON(out)
	}
	out.RawByte('}')
}

// MarshalJSON supports json.Marshaler interface
func (v UserOwnProfile) MarshalJSON() ([]byte, error) {
	w := jwriter.Writer{}
	easyjsonD2b7633eEncodeGithubComGoParkMailRu20211FyvaoldzhModels1(&w, v)
	return w.Buffer.BuildBytes(), w.Error
}

// MarshalEasyJSON supports easyjson.Marshaler interface
func (v UserOwnProfile) MarshalEasyJSON(w *jwriter.Writer) {
	easyjsonD2b7633eEncodeGithubComGoParkMailRu20211FyvaoldzhModels1(w, v)
}

// UnmarshalJSON supports json.Unmarshaler interface
func (v *UserOwnProfile) UnmarshalJSON(data []byte) error {
	r := jlexer.Lexer{Data: data}
	easyjsonD2b7633eDecodeGithubComGoParkMailRu20211FyvaoldzhModels1(&r, v)
	return r.Error()
}

// UnmarshalEasyJSON supports easyjson.Unmarshaler interface
func (v *UserOwnProfile) UnmarshalEasyJSON(l *jlexer.Lexer) {
	easyjsonD2b7633eDecodeGithubComGoParkMailRu20211FyvaoldzhModels1(l, v)
}
func easyjsonD2b7633eDecodeGithubComGoParkMailRu20211FyvaoldzhModels2(in *jlexer.Lexer, out *UserData) {
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
		case "name":
			out.Name = string(in.String())
		case "birthday":
			out.Birthday = string(in.String())
		case "city":
			out.City = string(in.String())
		case "email":
			out.Email = string(in.String())
		case "about":
			out.About = string(in.String())
		case "password":
			out.Password = string(in.String())
		case "avatar":
			out.Avatar = string(in.String())
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
func easyjsonD2b7633eEncodeGithubComGoParkMailRu20211FyvaoldzhModels2(out *jwriter.Writer, in UserData) {
	out.RawByte('{')
	first := true
	_ = first
	{
		const prefix string = ",\"name\":"
		out.RawString(prefix[1:])
		out.String(string(in.Name))
	}
	{
		const prefix string = ",\"birthday\":"
		out.RawString(prefix)
		out.String(string(in.Birthday))
	}
	{
		const prefix string = ",\"city\":"
		out.RawString(prefix)
		out.String(string(in.City))
	}
	{
		const prefix string = ",\"email\":"
		out.RawString(prefix)
		out.String(string(in.Email))
	}
	{
		const prefix string = ",\"about\":"
		out.RawString(prefix)
		out.String(string(in.About))
	}
	{
		const prefix string = ",\"password\":"
		out.RawString(prefix)
		out.String(string(in.Password))
	}
	{
		const prefix string = ",\"avatar\":"
		out.RawString(prefix)
		out.String(string(in.Avatar))
	}
	out.RawByte('}')
}

// MarshalJSON supports json.Marshaler interface
func (v UserData) MarshalJSON() ([]byte, error) {
	w := jwriter.Writer{}
	easyjsonD2b7633eEncodeGithubComGoParkMailRu20211FyvaoldzhModels2(&w, v)
	return w.Buffer.BuildBytes(), w.Error
}

// MarshalEasyJSON supports easyjson.Marshaler interface
func (v UserData) MarshalEasyJSON(w *jwriter.Writer) {
	easyjsonD2b7633eEncodeGithubComGoParkMailRu20211FyvaoldzhModels2(w, v)
}

// UnmarshalJSON supports json.Unmarshaler interface
func (v *UserData) UnmarshalJSON(data []byte) error {
	r := jlexer.Lexer{Data: data}
	easyjsonD2b7633eDecodeGithubComGoParkMailRu20211FyvaoldzhModels2(&r, v)
	return r.Error()
}

// UnmarshalEasyJSON supports easyjson.Unmarshaler interface
func (v *UserData) UnmarshalEasyJSON(l *jlexer.Lexer) {
	easyjsonD2b7633eDecodeGithubComGoParkMailRu20211FyvaoldzhModels2(l, v)
}
func easyjsonD2b7633eDecodeGithubComGoParkMailRu20211FyvaoldzhModels3(in *jlexer.Lexer, out *User) {
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
		case "Id":
			out.Id = int(in.Int())
		case "login":
			out.Login = string(in.String())
		case "password":
			out.Password = string(in.String())
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
func easyjsonD2b7633eEncodeGithubComGoParkMailRu20211FyvaoldzhModels3(out *jwriter.Writer, in User) {
	out.RawByte('{')
	first := true
	_ = first
	{
		const prefix string = ",\"Id\":"
		out.RawString(prefix[1:])
		out.Int(int(in.Id))
	}
	{
		const prefix string = ",\"login\":"
		out.RawString(prefix)
		out.String(string(in.Login))
	}
	{
		const prefix string = ",\"password\":"
		out.RawString(prefix)
		out.String(string(in.Password))
	}
	out.RawByte('}')
}

// MarshalJSON supports json.Marshaler interface
func (v User) MarshalJSON() ([]byte, error) {
	w := jwriter.Writer{}
	easyjsonD2b7633eEncodeGithubComGoParkMailRu20211FyvaoldzhModels3(&w, v)
	return w.Buffer.BuildBytes(), w.Error
}

// MarshalEasyJSON supports easyjson.Marshaler interface
func (v User) MarshalEasyJSON(w *jwriter.Writer) {
	easyjsonD2b7633eEncodeGithubComGoParkMailRu20211FyvaoldzhModels3(w, v)
}

// UnmarshalJSON supports json.Unmarshaler interface
func (v *User) UnmarshalJSON(data []byte) error {
	r := jlexer.Lexer{Data: data}
	easyjsonD2b7633eDecodeGithubComGoParkMailRu20211FyvaoldzhModels3(&r, v)
	return r.Error()
}

// UnmarshalEasyJSON supports easyjson.Unmarshaler interface
func (v *User) UnmarshalEasyJSON(l *jlexer.Lexer) {
	easyjsonD2b7633eDecodeGithubComGoParkMailRu20211FyvaoldzhModels3(l, v)
}
func easyjsonD2b7633eDecodeGithubComGoParkMailRu20211FyvaoldzhModels4(in *jlexer.Lexer, out *RegData) {
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
		case "Id":
			out.Id = int(in.Int())
		case "name":
			out.Name = string(in.String())
		case "login":
			out.Login = string(in.String())
		case "password":
			out.Password = string(in.String())
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
func easyjsonD2b7633eEncodeGithubComGoParkMailRu20211FyvaoldzhModels4(out *jwriter.Writer, in RegData) {
	out.RawByte('{')
	first := true
	_ = first
	{
		const prefix string = ",\"Id\":"
		out.RawString(prefix[1:])
		out.Int(int(in.Id))
	}
	{
		const prefix string = ",\"name\":"
		out.RawString(prefix)
		out.String(string(in.Name))
	}
	{
		const prefix string = ",\"login\":"
		out.RawString(prefix)
		out.String(string(in.Login))
	}
	{
		const prefix string = ",\"password\":"
		out.RawString(prefix)
		out.String(string(in.Password))
	}
	out.RawByte('}')
}

// MarshalJSON supports json.Marshaler interface
func (v RegData) MarshalJSON() ([]byte, error) {
	w := jwriter.Writer{}
	easyjsonD2b7633eEncodeGithubComGoParkMailRu20211FyvaoldzhModels4(&w, v)
	return w.Buffer.BuildBytes(), w.Error
}

// MarshalEasyJSON supports easyjson.Marshaler interface
func (v RegData) MarshalEasyJSON(w *jwriter.Writer) {
	easyjsonD2b7633eEncodeGithubComGoParkMailRu20211FyvaoldzhModels4(w, v)
}

// UnmarshalJSON supports json.Unmarshaler interface
func (v *RegData) UnmarshalJSON(data []byte) error {
	r := jlexer.Lexer{Data: data}
	easyjsonD2b7633eDecodeGithubComGoParkMailRu20211FyvaoldzhModels4(&r, v)
	return r.Error()
}

// UnmarshalEasyJSON supports easyjson.Unmarshaler interface
func (v *RegData) UnmarshalEasyJSON(l *jlexer.Lexer) {
	easyjsonD2b7633eDecodeGithubComGoParkMailRu20211FyvaoldzhModels4(l, v)
}
func easyjsonD2b7633eDecodeGithubComGoParkMailRu20211FyvaoldzhModels5(in *jlexer.Lexer, out *PlanningEvents) {
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
		case "Uid":
			out.Uid = int(in.Int())
		case "Eid":
			out.Eid = uint64(in.Uint64())
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
func easyjsonD2b7633eEncodeGithubComGoParkMailRu20211FyvaoldzhModels5(out *jwriter.Writer, in PlanningEvents) {
	out.RawByte('{')
	first := true
	_ = first
	{
		const prefix string = ",\"Uid\":"
		out.RawString(prefix[1:])
		out.Int(int(in.Uid))
	}
	{
		const prefix string = ",\"Eid\":"
		out.RawString(prefix)
		out.Uint64(uint64(in.Eid))
	}
	out.RawByte('}')
}

// MarshalJSON supports json.Marshaler interface
func (v PlanningEvents) MarshalJSON() ([]byte, error) {
	w := jwriter.Writer{}
	easyjsonD2b7633eEncodeGithubComGoParkMailRu20211FyvaoldzhModels5(&w, v)
	return w.Buffer.BuildBytes(), w.Error
}

// MarshalEasyJSON supports easyjson.Marshaler interface
func (v PlanningEvents) MarshalEasyJSON(w *jwriter.Writer) {
	easyjsonD2b7633eEncodeGithubComGoParkMailRu20211FyvaoldzhModels5(w, v)
}

// UnmarshalJSON supports json.Unmarshaler interface
func (v *PlanningEvents) UnmarshalJSON(data []byte) error {
	r := jlexer.Lexer{Data: data}
	easyjsonD2b7633eDecodeGithubComGoParkMailRu20211FyvaoldzhModels5(&r, v)
	return r.Error()
}

// UnmarshalEasyJSON supports easyjson.Unmarshaler interface
func (v *PlanningEvents) UnmarshalEasyJSON(l *jlexer.Lexer) {
	easyjsonD2b7633eDecodeGithubComGoParkMailRu20211FyvaoldzhModels5(l, v)
}
func easyjsonD2b7633eDecodeGithubComGoParkMailRu20211FyvaoldzhModels6(in *jlexer.Lexer, out *Events) {
	isTopLevel := in.IsStart()
	if in.IsNull() {
		in.Skip()
		*out = nil
	} else {
		in.Delim('[')
		if *out == nil {
			if !in.IsDelim(']') {
				*out = make(Events, 0, 0)
			} else {
				*out = Events{}
			}
		} else {
			*out = (*out)[:0]
		}
		for !in.IsDelim(']') {
			var v1 Event
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
func easyjsonD2b7633eEncodeGithubComGoParkMailRu20211FyvaoldzhModels6(out *jwriter.Writer, in Events) {
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
func (v Events) MarshalJSON() ([]byte, error) {
	w := jwriter.Writer{}
	easyjsonD2b7633eEncodeGithubComGoParkMailRu20211FyvaoldzhModels6(&w, v)
	return w.Buffer.BuildBytes(), w.Error
}

// MarshalEasyJSON supports easyjson.Marshaler interface
func (v Events) MarshalEasyJSON(w *jwriter.Writer) {
	easyjsonD2b7633eEncodeGithubComGoParkMailRu20211FyvaoldzhModels6(w, v)
}

// UnmarshalJSON supports json.Unmarshaler interface
func (v *Events) UnmarshalJSON(data []byte) error {
	r := jlexer.Lexer{Data: data}
	easyjsonD2b7633eDecodeGithubComGoParkMailRu20211FyvaoldzhModels6(&r, v)
	return r.Error()
}

// UnmarshalEasyJSON supports easyjson.Unmarshaler interface
func (v *Events) UnmarshalEasyJSON(l *jlexer.Lexer) {
	easyjsonD2b7633eDecodeGithubComGoParkMailRu20211FyvaoldzhModels6(l, v)
}
func easyjsonD2b7633eDecodeGithubComGoParkMailRu20211FyvaoldzhModels7(in *jlexer.Lexer, out *Event) {
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
			out.ID = uint64(in.Uint64())
		case "title":
			out.Title = string(in.String())
		case "place":
			out.Place = string(in.String())
		case "description":
			out.Description = string(in.String())
		case "time":
			out.Date = string(in.String())
		case "subway":
			out.Subway = string(in.String())
		case "street":
			out.Street = string(in.String())
		case "typeEvent":
			out.TypeEvent = string(in.String())
		case "image":
			out.Image = string(in.String())
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
func easyjsonD2b7633eEncodeGithubComGoParkMailRu20211FyvaoldzhModels7(out *jwriter.Writer, in Event) {
	out.RawByte('{')
	first := true
	_ = first
	{
		const prefix string = ",\"id\":"
		out.RawString(prefix[1:])
		out.Uint64(uint64(in.ID))
	}
	{
		const prefix string = ",\"title\":"
		out.RawString(prefix)
		out.String(string(in.Title))
	}
	{
		const prefix string = ",\"place\":"
		out.RawString(prefix)
		out.String(string(in.Place))
	}
	{
		const prefix string = ",\"description\":"
		out.RawString(prefix)
		out.String(string(in.Description))
	}
	{
		const prefix string = ",\"time\":"
		out.RawString(prefix)
		out.String(string(in.Date))
	}
	{
		const prefix string = ",\"subway\":"
		out.RawString(prefix)
		out.String(string(in.Subway))
	}
	{
		const prefix string = ",\"street\":"
		out.RawString(prefix)
		out.String(string(in.Street))
	}
	{
		const prefix string = ",\"typeEvent\":"
		out.RawString(prefix)
		out.String(string(in.TypeEvent))
	}
	{
		const prefix string = ",\"image\":"
		out.RawString(prefix)
		out.String(string(in.Image))
	}
	out.RawByte('}')
}

// MarshalJSON supports json.Marshaler interface
func (v Event) MarshalJSON() ([]byte, error) {
	w := jwriter.Writer{}
	easyjsonD2b7633eEncodeGithubComGoParkMailRu20211FyvaoldzhModels7(&w, v)
	return w.Buffer.BuildBytes(), w.Error
}

// MarshalEasyJSON supports easyjson.Marshaler interface
func (v Event) MarshalEasyJSON(w *jwriter.Writer) {
	easyjsonD2b7633eEncodeGithubComGoParkMailRu20211FyvaoldzhModels7(w, v)
}

// UnmarshalJSON supports json.Unmarshaler interface
func (v *Event) UnmarshalJSON(data []byte) error {
	r := jlexer.Lexer{Data: data}
	easyjsonD2b7633eDecodeGithubComGoParkMailRu20211FyvaoldzhModels7(&r, v)
	return r.Error()
}

// UnmarshalEasyJSON supports easyjson.Unmarshaler interface
func (v *Event) UnmarshalEasyJSON(l *jlexer.Lexer) {
	easyjsonD2b7633eDecodeGithubComGoParkMailRu20211FyvaoldzhModels7(l, v)
}
