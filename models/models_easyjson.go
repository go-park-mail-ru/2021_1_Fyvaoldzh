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

func easyjsonD2b7633eDecodeKudagoModels(in *jlexer.Lexer, out *UserProfile) {
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
			out.Uid = uint64(in.Uint64())
		case "name":
			out.Name = string(in.String())
		case "age":
			out.Age = uint8(in.Uint8())
		case "city":
			out.City = string(in.String())
		case "followers":
			out.Followers = uint64(in.Uint64())
		case "about":
			out.About = string(in.String())
		case "avatar":
			out.Avatar = string(in.String())
		case "events":
			if in.IsNull() {
				in.Skip()
				out.Event = nil
			} else {
				in.Delim('[')
				if out.Event == nil {
					if !in.IsDelim(']') {
						out.Event = make([]uint64, 0, 8)
					} else {
						out.Event = []uint64{}
					}
				} else {
					out.Event = (out.Event)[:0]
				}
				for !in.IsDelim(']') {
					var v1 uint64
					v1 = uint64(in.Uint64())
					out.Event = append(out.Event, v1)
					in.WantComma()
				}
				in.Delim(']')
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
func easyjsonD2b7633eEncodeKudagoModels(out *jwriter.Writer, in UserProfile) {
	out.RawByte('{')
	first := true
	_ = first
	{
		const prefix string = ",\"Uid\":"
		out.RawString(prefix[1:])
		out.Uint64(uint64(in.Uid))
	}
	{
		const prefix string = ",\"name\":"
		out.RawString(prefix)
		out.String(string(in.Name))
	}
	{
		const prefix string = ",\"age\":"
		out.RawString(prefix)
		out.Uint8(uint8(in.Age))
	}
	{
		const prefix string = ",\"city\":"
		out.RawString(prefix)
		out.String(string(in.City))
	}
	{
		const prefix string = ",\"followers\":"
		out.RawString(prefix)
		out.Uint64(uint64(in.Followers))
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
		if in.Event == nil && (out.Flags&jwriter.NilSliceAsEmpty) == 0 {
			out.RawString("null")
		} else {
			out.RawByte('[')
			for v2, v3 := range in.Event {
				if v2 > 0 {
					out.RawByte(',')
				}
				out.Uint64(uint64(v3))
			}
			out.RawByte(']')
		}
	}
	out.RawByte('}')
}

// MarshalJSON supports json.Marshaler interface
func (v UserProfile) MarshalJSON() ([]byte, error) {
	w := jwriter.Writer{}
	easyjsonD2b7633eEncodeKudagoModels(&w, v)
	return w.Buffer.BuildBytes(), w.Error
}

// MarshalEasyJSON supports easyjson.Marshaler interface
func (v UserProfile) MarshalEasyJSON(w *jwriter.Writer) {
	easyjsonD2b7633eEncodeKudagoModels(w, v)
}

// UnmarshalJSON supports json.Unmarshaler interface
func (v *UserProfile) UnmarshalJSON(data []byte) error {
	r := jlexer.Lexer{Data: data}
	easyjsonD2b7633eDecodeKudagoModels(&r, v)
	return r.Error()
}

// UnmarshalEasyJSON supports easyjson.Unmarshaler interface
func (v *UserProfile) UnmarshalEasyJSON(l *jlexer.Lexer) {
	easyjsonD2b7633eDecodeKudagoModels(l, v)
}
func easyjsonD2b7633eDecodeKudagoModels1(in *jlexer.Lexer, out *UserOwnProfile) {
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
			out.Uid = uint64(in.Uint64())
		case "name":
			out.Name = string(in.String())
		case "birthday":
			out.Birthday = string(in.String())
		case "city":
			out.City = string(in.String())
		case "email":
			out.Email = string(in.String())
		case "visited":
			out.Visited = uint64(in.Uint64())
		case "planning":
			out.Planning = uint64(in.Uint64())
		case "followers":
			out.Followers = uint64(in.Uint64())
		case "about":
			out.About = string(in.String())
		case "avatar":
			out.Avatar = string(in.String())
		case "events":
			if in.IsNull() {
				in.Skip()
				out.Event = nil
			} else {
				in.Delim('[')
				if out.Event == nil {
					if !in.IsDelim(']') {
						out.Event = make([]uint64, 0, 8)
					} else {
						out.Event = []uint64{}
					}
				} else {
					out.Event = (out.Event)[:0]
				}
				for !in.IsDelim(']') {
					var v4 uint64
					v4 = uint64(in.Uint64())
					out.Event = append(out.Event, v4)
					in.WantComma()
				}
				in.Delim(']')
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
func easyjsonD2b7633eEncodeKudagoModels1(out *jwriter.Writer, in UserOwnProfile) {
	out.RawByte('{')
	first := true
	_ = first
	{
		const prefix string = ",\"Uid\":"
		out.RawString(prefix[1:])
		out.Uint64(uint64(in.Uid))
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
		out.Uint64(uint64(in.Visited))
	}
	{
		const prefix string = ",\"planning\":"
		out.RawString(prefix)
		out.Uint64(uint64(in.Planning))
	}
	{
		const prefix string = ",\"followers\":"
		out.RawString(prefix)
		out.Uint64(uint64(in.Followers))
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
		if in.Event == nil && (out.Flags&jwriter.NilSliceAsEmpty) == 0 {
			out.RawString("null")
		} else {
			out.RawByte('[')
			for v5, v6 := range in.Event {
				if v5 > 0 {
					out.RawByte(',')
				}
				out.Uint64(uint64(v6))
			}
			out.RawByte(']')
		}
	}
	out.RawByte('}')
}

// MarshalJSON supports json.Marshaler interface
func (v UserOwnProfile) MarshalJSON() ([]byte, error) {
	w := jwriter.Writer{}
	easyjsonD2b7633eEncodeKudagoModels1(&w, v)
	return w.Buffer.BuildBytes(), w.Error
}

// MarshalEasyJSON supports easyjson.Marshaler interface
func (v UserOwnProfile) MarshalEasyJSON(w *jwriter.Writer) {
	easyjsonD2b7633eEncodeKudagoModels1(w, v)
}

// UnmarshalJSON supports json.Unmarshaler interface
func (v *UserOwnProfile) UnmarshalJSON(data []byte) error {
	r := jlexer.Lexer{Data: data}
	easyjsonD2b7633eDecodeKudagoModels1(&r, v)
	return r.Error()
}

// UnmarshalEasyJSON supports easyjson.Unmarshaler interface
func (v *UserOwnProfile) UnmarshalEasyJSON(l *jlexer.Lexer) {
	easyjsonD2b7633eDecodeKudagoModels1(l, v)
}
func easyjsonD2b7633eDecodeKudagoModels2(in *jlexer.Lexer, out *UserEvents) {
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
			out.Uid = uint64(in.Uint64())
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
func easyjsonD2b7633eEncodeKudagoModels2(out *jwriter.Writer, in UserEvents) {
	out.RawByte('{')
	first := true
	_ = first
	{
		const prefix string = ",\"Uid\":"
		out.RawString(prefix[1:])
		out.Uint64(uint64(in.Uid))
	}
	{
		const prefix string = ",\"Eid\":"
		out.RawString(prefix)
		out.Uint64(uint64(in.Eid))
	}
	out.RawByte('}')
}

// MarshalJSON supports json.Marshaler interface
func (v UserEvents) MarshalJSON() ([]byte, error) {
	w := jwriter.Writer{}
	easyjsonD2b7633eEncodeKudagoModels2(&w, v)
	return w.Buffer.BuildBytes(), w.Error
}

// MarshalEasyJSON supports easyjson.Marshaler interface
func (v UserEvents) MarshalEasyJSON(w *jwriter.Writer) {
	easyjsonD2b7633eEncodeKudagoModels2(w, v)
}

// UnmarshalJSON supports json.Unmarshaler interface
func (v *UserEvents) UnmarshalJSON(data []byte) error {
	r := jlexer.Lexer{Data: data}
	easyjsonD2b7633eDecodeKudagoModels2(&r, v)
	return r.Error()
}

// UnmarshalEasyJSON supports easyjson.Unmarshaler interface
func (v *UserEvents) UnmarshalEasyJSON(l *jlexer.Lexer) {
	easyjsonD2b7633eDecodeKudagoModels2(l, v)
}
func easyjsonD2b7633eDecodeKudagoModels3(in *jlexer.Lexer, out *UserData) {
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
func easyjsonD2b7633eEncodeKudagoModels3(out *jwriter.Writer, in UserData) {
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
	easyjsonD2b7633eEncodeKudagoModels3(&w, v)
	return w.Buffer.BuildBytes(), w.Error
}

// MarshalEasyJSON supports easyjson.Marshaler interface
func (v UserData) MarshalEasyJSON(w *jwriter.Writer) {
	easyjsonD2b7633eEncodeKudagoModels3(w, v)
}

// UnmarshalJSON supports json.Unmarshaler interface
func (v *UserData) UnmarshalJSON(data []byte) error {
	r := jlexer.Lexer{Data: data}
	easyjsonD2b7633eDecodeKudagoModels3(&r, v)
	return r.Error()
}

// UnmarshalEasyJSON supports easyjson.Unmarshaler interface
func (v *UserData) UnmarshalEasyJSON(l *jlexer.Lexer) {
	easyjsonD2b7633eDecodeKudagoModels3(l, v)
}
func easyjsonD2b7633eDecodeKudagoModels4(in *jlexer.Lexer, out *User) {
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
			out.Id = uint64(in.Uint64())
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
func easyjsonD2b7633eEncodeKudagoModels4(out *jwriter.Writer, in User) {
	out.RawByte('{')
	first := true
	_ = first
	{
		const prefix string = ",\"Id\":"
		out.RawString(prefix[1:])
		out.Uint64(uint64(in.Id))
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
	easyjsonD2b7633eEncodeKudagoModels4(&w, v)
	return w.Buffer.BuildBytes(), w.Error
}

// MarshalEasyJSON supports easyjson.Marshaler interface
func (v User) MarshalEasyJSON(w *jwriter.Writer) {
	easyjsonD2b7633eEncodeKudagoModels4(w, v)
}

// UnmarshalJSON supports json.Unmarshaler interface
func (v *User) UnmarshalJSON(data []byte) error {
	r := jlexer.Lexer{Data: data}
	easyjsonD2b7633eDecodeKudagoModels4(&r, v)
	return r.Error()
}

// UnmarshalEasyJSON supports easyjson.Unmarshaler interface
func (v *User) UnmarshalEasyJSON(l *jlexer.Lexer) {
	easyjsonD2b7633eDecodeKudagoModels4(l, v)
}
func easyjsonD2b7633eDecodeKudagoModels5(in *jlexer.Lexer, out *RegData) {
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
			out.Id = uint64(in.Uint64())
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
func easyjsonD2b7633eEncodeKudagoModels5(out *jwriter.Writer, in RegData) {
	out.RawByte('{')
	first := true
	_ = first
	{
		const prefix string = ",\"Id\":"
		out.RawString(prefix[1:])
		out.Uint64(uint64(in.Id))
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
	easyjsonD2b7633eEncodeKudagoModels5(&w, v)
	return w.Buffer.BuildBytes(), w.Error
}

// MarshalEasyJSON supports easyjson.Marshaler interface
func (v RegData) MarshalEasyJSON(w *jwriter.Writer) {
	easyjsonD2b7633eEncodeKudagoModels5(w, v)
}

// UnmarshalJSON supports json.Unmarshaler interface
func (v *RegData) UnmarshalJSON(data []byte) error {
	r := jlexer.Lexer{Data: data}
	easyjsonD2b7633eDecodeKudagoModels5(&r, v)
	return r.Error()
}

// UnmarshalEasyJSON supports easyjson.Unmarshaler interface
func (v *RegData) UnmarshalEasyJSON(l *jlexer.Lexer) {
	easyjsonD2b7633eDecodeKudagoModels5(l, v)
}
func easyjsonD2b7633eDecodeKudagoModels6(in *jlexer.Lexer, out *Events) {
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
			var v7 Event
			(v7).UnmarshalEasyJSON(in)
			*out = append(*out, v7)
			in.WantComma()
		}
		in.Delim(']')
	}
	if isTopLevel {
		in.Consumed()
	}
}
func easyjsonD2b7633eEncodeKudagoModels6(out *jwriter.Writer, in Events) {
	if in == nil && (out.Flags&jwriter.NilSliceAsEmpty) == 0 {
		out.RawString("null")
	} else {
		out.RawByte('[')
		for v8, v9 := range in {
			if v8 > 0 {
				out.RawByte(',')
			}
			(v9).MarshalEasyJSON(out)
		}
		out.RawByte(']')
	}
}

// MarshalJSON supports json.Marshaler interface
func (v Events) MarshalJSON() ([]byte, error) {
	w := jwriter.Writer{}
	easyjsonD2b7633eEncodeKudagoModels6(&w, v)
	return w.Buffer.BuildBytes(), w.Error
}

// MarshalEasyJSON supports easyjson.Marshaler interface
func (v Events) MarshalEasyJSON(w *jwriter.Writer) {
	easyjsonD2b7633eEncodeKudagoModels6(w, v)
}

// UnmarshalJSON supports json.Unmarshaler interface
func (v *Events) UnmarshalJSON(data []byte) error {
	r := jlexer.Lexer{Data: data}
	easyjsonD2b7633eDecodeKudagoModels6(&r, v)
	return r.Error()
}

// UnmarshalEasyJSON supports easyjson.Unmarshaler interface
func (v *Events) UnmarshalEasyJSON(l *jlexer.Lexer) {
	easyjsonD2b7633eDecodeKudagoModels6(l, v)
}
func easyjsonD2b7633eDecodeKudagoModels7(in *jlexer.Lexer, out *Event) {
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
		case "date":
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
func easyjsonD2b7633eEncodeKudagoModels7(out *jwriter.Writer, in Event) {
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
		const prefix string = ",\"date\":"
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
	easyjsonD2b7633eEncodeKudagoModels7(&w, v)
	return w.Buffer.BuildBytes(), w.Error
}

// MarshalEasyJSON supports easyjson.Marshaler interface
func (v Event) MarshalEasyJSON(w *jwriter.Writer) {
	easyjsonD2b7633eEncodeKudagoModels7(w, v)
}

// UnmarshalJSON supports json.Unmarshaler interface
func (v *Event) UnmarshalJSON(data []byte) error {
	r := jlexer.Lexer{Data: data}
	easyjsonD2b7633eDecodeKudagoModels7(&r, v)
	return r.Error()
}

// UnmarshalEasyJSON supports easyjson.Unmarshaler interface
func (v *Event) UnmarshalEasyJSON(l *jlexer.Lexer) {
	easyjsonD2b7633eDecodeKudagoModels7(l, v)
}
