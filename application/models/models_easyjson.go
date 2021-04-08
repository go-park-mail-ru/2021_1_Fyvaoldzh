// Code generated by easyjson for marshaling/unmarshaling. DO NOT EDIT.

package models

import (
	sql "database/sql"
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

func easyjsonD2b7633eDecodeKudagoApplicationModels(in *jlexer.Lexer, out *UserOwnProfile) {
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
		case "login":
			out.Login = string(in.String())
		case "birthday":
			out.Birthday = string(in.String())
		case "city":
			out.City = string(in.String())
		case "email":
			out.Email = string(in.String())
		case "visited":
			(out.Visited).UnmarshalEasyJSON(in)
		case "planning":
			(out.Planning).UnmarshalEasyJSON(in)
		case "followers":
			if in.IsNull() {
				in.Skip()
				out.Followers = nil
			} else {
				in.Delim('[')
				if out.Followers == nil {
					if !in.IsDelim(']') {
						out.Followers = make([]uint64, 0, 8)
					} else {
						out.Followers = []uint64{}
					}
				} else {
					out.Followers = (out.Followers)[:0]
				}
				for !in.IsDelim(']') {
					var v1 uint64
					v1 = uint64(in.Uint64())
					out.Followers = append(out.Followers, v1)
					in.WantComma()
				}
				in.Delim(']')
			}
		case "about":
			out.About = string(in.String())
		case "avatar":
			out.Avatar = string(in.String())
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
func easyjsonD2b7633eEncodeKudagoApplicationModels(out *jwriter.Writer, in UserOwnProfile) {
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
		const prefix string = ",\"login\":"
		out.RawString(prefix)
		out.String(string(in.Login))
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
		(in.Visited).MarshalEasyJSON(out)
	}
	{
		const prefix string = ",\"planning\":"
		out.RawString(prefix)
		(in.Planning).MarshalEasyJSON(out)
	}
	{
		const prefix string = ",\"followers\":"
		out.RawString(prefix)
		if in.Followers == nil && (out.Flags&jwriter.NilSliceAsEmpty) == 0 {
			out.RawString("null")
		} else {
			out.RawByte('[')
			for v2, v3 := range in.Followers {
				if v2 > 0 {
					out.RawByte(',')
				}
				out.Uint64(uint64(v3))
			}
			out.RawByte(']')
		}
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
		const prefix string = ",\"password\":"
		out.RawString(prefix)
		out.String(string(in.Password))
	}
	out.RawByte('}')
}

// MarshalJSON supports json.Marshaler interface
func (v UserOwnProfile) MarshalJSON() ([]byte, error) {
	w := jwriter.Writer{}
	easyjsonD2b7633eEncodeKudagoApplicationModels(&w, v)
	return w.Buffer.BuildBytes(), w.Error
}

// MarshalEasyJSON supports easyjson.Marshaler interface
func (v UserOwnProfile) MarshalEasyJSON(w *jwriter.Writer) {
	easyjsonD2b7633eEncodeKudagoApplicationModels(w, v)
}

// UnmarshalJSON supports json.Unmarshaler interface
func (v *UserOwnProfile) UnmarshalJSON(data []byte) error {
	r := jlexer.Lexer{Data: data}
	easyjsonD2b7633eDecodeKudagoApplicationModels(&r, v)
	return r.Error()
}

// UnmarshalEasyJSON supports easyjson.Unmarshaler interface
func (v *UserOwnProfile) UnmarshalEasyJSON(l *jlexer.Lexer) {
	easyjsonD2b7633eDecodeKudagoApplicationModels(l, v)
}
func easyjsonD2b7633eDecodeKudagoApplicationModels1(in *jlexer.Lexer, out *UserEvents) {
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
func easyjsonD2b7633eEncodeKudagoApplicationModels1(out *jwriter.Writer, in UserEvents) {
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
	easyjsonD2b7633eEncodeKudagoApplicationModels1(&w, v)
	return w.Buffer.BuildBytes(), w.Error
}

// MarshalEasyJSON supports easyjson.Marshaler interface
func (v UserEvents) MarshalEasyJSON(w *jwriter.Writer) {
	easyjsonD2b7633eEncodeKudagoApplicationModels1(w, v)
}

// UnmarshalJSON supports json.Unmarshaler interface
func (v *UserEvents) UnmarshalJSON(data []byte) error {
	r := jlexer.Lexer{Data: data}
	easyjsonD2b7633eDecodeKudagoApplicationModels1(&r, v)
	return r.Error()
}

// UnmarshalEasyJSON supports easyjson.Unmarshaler interface
func (v *UserEvents) UnmarshalEasyJSON(l *jlexer.Lexer) {
	easyjsonD2b7633eDecodeKudagoApplicationModels1(l, v)
}
func easyjsonD2b7633eDecodeKudagoApplicationModels2(in *jlexer.Lexer, out *UserData) {
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
		case "Name":
			easyjsonD2b7633eDecodeDatabaseSql(in, &out.Name)
		case "Login":
			out.Login = string(in.String())
		case "Birthday":
			easyjsonD2b7633eDecodeDatabaseSql1(in, &out.Birthday)
		case "City":
			easyjsonD2b7633eDecodeDatabaseSql(in, &out.City)
		case "Email":
			easyjsonD2b7633eDecodeDatabaseSql(in, &out.Email)
		case "About":
			easyjsonD2b7633eDecodeDatabaseSql(in, &out.About)
		case "Password":
			easyjsonD2b7633eDecodeDatabaseSql(in, &out.Password)
		case "Avatar":
			easyjsonD2b7633eDecodeDatabaseSql(in, &out.Avatar)
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
func easyjsonD2b7633eEncodeKudagoApplicationModels2(out *jwriter.Writer, in UserData) {
	out.RawByte('{')
	first := true
	_ = first
	{
		const prefix string = ",\"Id\":"
		out.RawString(prefix[1:])
		out.Uint64(uint64(in.Id))
	}
	{
		const prefix string = ",\"Name\":"
		out.RawString(prefix)
		easyjsonD2b7633eEncodeDatabaseSql(out, in.Name)
	}
	{
		const prefix string = ",\"Login\":"
		out.RawString(prefix)
		out.String(string(in.Login))
	}
	{
		const prefix string = ",\"Birthday\":"
		out.RawString(prefix)
		easyjsonD2b7633eEncodeDatabaseSql1(out, in.Birthday)
	}
	{
		const prefix string = ",\"City\":"
		out.RawString(prefix)
		easyjsonD2b7633eEncodeDatabaseSql(out, in.City)
	}
	{
		const prefix string = ",\"Email\":"
		out.RawString(prefix)
		easyjsonD2b7633eEncodeDatabaseSql(out, in.Email)
	}
	{
		const prefix string = ",\"About\":"
		out.RawString(prefix)
		easyjsonD2b7633eEncodeDatabaseSql(out, in.About)
	}
	{
		const prefix string = ",\"Password\":"
		out.RawString(prefix)
		easyjsonD2b7633eEncodeDatabaseSql(out, in.Password)
	}
	{
		const prefix string = ",\"Avatar\":"
		out.RawString(prefix)
		easyjsonD2b7633eEncodeDatabaseSql(out, in.Avatar)
	}
	out.RawByte('}')
}

// MarshalJSON supports json.Marshaler interface
func (v UserData) MarshalJSON() ([]byte, error) {
	w := jwriter.Writer{}
	easyjsonD2b7633eEncodeKudagoApplicationModels2(&w, v)
	return w.Buffer.BuildBytes(), w.Error
}

// MarshalEasyJSON supports easyjson.Marshaler interface
func (v UserData) MarshalEasyJSON(w *jwriter.Writer) {
	easyjsonD2b7633eEncodeKudagoApplicationModels2(w, v)
}

// UnmarshalJSON supports json.Unmarshaler interface
func (v *UserData) UnmarshalJSON(data []byte) error {
	r := jlexer.Lexer{Data: data}
	easyjsonD2b7633eDecodeKudagoApplicationModels2(&r, v)
	return r.Error()
}

// UnmarshalEasyJSON supports easyjson.Unmarshaler interface
func (v *UserData) UnmarshalEasyJSON(l *jlexer.Lexer) {
	easyjsonD2b7633eDecodeKudagoApplicationModels2(l, v)
}
func easyjsonD2b7633eDecodeDatabaseSql1(in *jlexer.Lexer, out *sql.NullTime) {
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
		case "Time":
			if data := in.Raw(); in.Ok() {
				in.AddError((out.Time).UnmarshalJSON(data))
			}
		case "Valid":
			out.Valid = bool(in.Bool())
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
func easyjsonD2b7633eEncodeDatabaseSql1(out *jwriter.Writer, in sql.NullTime) {
	out.RawByte('{')
	first := true
	_ = first
	{
		const prefix string = ",\"Time\":"
		out.RawString(prefix[1:])
		out.Raw((in.Time).MarshalJSON())
	}
	{
		const prefix string = ",\"Valid\":"
		out.RawString(prefix)
		out.Bool(bool(in.Valid))
	}
	out.RawByte('}')
}
func easyjsonD2b7633eDecodeDatabaseSql(in *jlexer.Lexer, out *sql.NullString) {
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
		case "String":
			out.String = string(in.String())
		case "Valid":
			out.Valid = bool(in.Bool())
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
func easyjsonD2b7633eEncodeDatabaseSql(out *jwriter.Writer, in sql.NullString) {
	out.RawByte('{')
	first := true
	_ = first
	{
		const prefix string = ",\"String\":"
		out.RawString(prefix[1:])
		out.String(string(in.String))
	}
	{
		const prefix string = ",\"Valid\":"
		out.RawString(prefix)
		out.Bool(bool(in.Valid))
	}
	out.RawByte('}')
}
func easyjsonD2b7633eDecodeKudagoApplicationModels3(in *jlexer.Lexer, out *User) {
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
func easyjsonD2b7633eEncodeKudagoApplicationModels3(out *jwriter.Writer, in User) {
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
	easyjsonD2b7633eEncodeKudagoApplicationModels3(&w, v)
	return w.Buffer.BuildBytes(), w.Error
}

// MarshalEasyJSON supports easyjson.Marshaler interface
func (v User) MarshalEasyJSON(w *jwriter.Writer) {
	easyjsonD2b7633eEncodeKudagoApplicationModels3(w, v)
}

// UnmarshalJSON supports json.Unmarshaler interface
func (v *User) UnmarshalJSON(data []byte) error {
	r := jlexer.Lexer{Data: data}
	easyjsonD2b7633eDecodeKudagoApplicationModels3(&r, v)
	return r.Error()
}

// UnmarshalEasyJSON supports easyjson.Unmarshaler interface
func (v *User) UnmarshalEasyJSON(l *jlexer.Lexer) {
	easyjsonD2b7633eDecodeKudagoApplicationModels3(l, v)
}
func easyjsonD2b7633eDecodeKudagoApplicationModels4(in *jlexer.Lexer, out *Tags) {
	isTopLevel := in.IsStart()
	if in.IsNull() {
		in.Skip()
		*out = nil
	} else {
		in.Delim('[')
		if *out == nil {
			if !in.IsDelim(']') {
				*out = make(Tags, 0, 2)
			} else {
				*out = Tags{}
			}
		} else {
			*out = (*out)[:0]
		}
		for !in.IsDelim(']') {
			var v4 Tag
			(v4).UnmarshalEasyJSON(in)
			*out = append(*out, v4)
			in.WantComma()
		}
		in.Delim(']')
	}
	if isTopLevel {
		in.Consumed()
	}
}
func easyjsonD2b7633eEncodeKudagoApplicationModels4(out *jwriter.Writer, in Tags) {
	if in == nil && (out.Flags&jwriter.NilSliceAsEmpty) == 0 {
		out.RawString("null")
	} else {
		out.RawByte('[')
		for v5, v6 := range in {
			if v5 > 0 {
				out.RawByte(',')
			}
			(v6).MarshalEasyJSON(out)
		}
		out.RawByte(']')
	}
}

// MarshalJSON supports json.Marshaler interface
func (v Tags) MarshalJSON() ([]byte, error) {
	w := jwriter.Writer{}
	easyjsonD2b7633eEncodeKudagoApplicationModels4(&w, v)
	return w.Buffer.BuildBytes(), w.Error
}

// MarshalEasyJSON supports easyjson.Marshaler interface
func (v Tags) MarshalEasyJSON(w *jwriter.Writer) {
	easyjsonD2b7633eEncodeKudagoApplicationModels4(w, v)
}

// UnmarshalJSON supports json.Unmarshaler interface
func (v *Tags) UnmarshalJSON(data []byte) error {
	r := jlexer.Lexer{Data: data}
	easyjsonD2b7633eDecodeKudagoApplicationModels4(&r, v)
	return r.Error()
}

// UnmarshalEasyJSON supports easyjson.Unmarshaler interface
func (v *Tags) UnmarshalEasyJSON(l *jlexer.Lexer) {
	easyjsonD2b7633eDecodeKudagoApplicationModels4(l, v)
}
func easyjsonD2b7633eDecodeKudagoApplicationModels5(in *jlexer.Lexer, out *Tag) {
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
		case "name":
			out.Name = string(in.String())
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
func easyjsonD2b7633eEncodeKudagoApplicationModels5(out *jwriter.Writer, in Tag) {
	out.RawByte('{')
	first := true
	_ = first
	{
		const prefix string = ",\"id\":"
		out.RawString(prefix[1:])
		out.Uint64(uint64(in.ID))
	}
	{
		const prefix string = ",\"name\":"
		out.RawString(prefix)
		out.String(string(in.Name))
	}
	out.RawByte('}')
}

// MarshalJSON supports json.Marshaler interface
func (v Tag) MarshalJSON() ([]byte, error) {
	w := jwriter.Writer{}
	easyjsonD2b7633eEncodeKudagoApplicationModels5(&w, v)
	return w.Buffer.BuildBytes(), w.Error
}

// MarshalEasyJSON supports easyjson.Marshaler interface
func (v Tag) MarshalEasyJSON(w *jwriter.Writer) {
	easyjsonD2b7633eEncodeKudagoApplicationModels5(w, v)
}

// UnmarshalJSON supports json.Unmarshaler interface
func (v *Tag) UnmarshalJSON(data []byte) error {
	r := jlexer.Lexer{Data: data}
	easyjsonD2b7633eDecodeKudagoApplicationModels5(&r, v)
	return r.Error()
}

// UnmarshalEasyJSON supports easyjson.Unmarshaler interface
func (v *Tag) UnmarshalEasyJSON(l *jlexer.Lexer) {
	easyjsonD2b7633eDecodeKudagoApplicationModels5(l, v)
}
func easyjsonD2b7633eDecodeKudagoApplicationModels6(in *jlexer.Lexer, out *RegData) {
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
func easyjsonD2b7633eEncodeKudagoApplicationModels6(out *jwriter.Writer, in RegData) {
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
	easyjsonD2b7633eEncodeKudagoApplicationModels6(&w, v)
	return w.Buffer.BuildBytes(), w.Error
}

// MarshalEasyJSON supports easyjson.Marshaler interface
func (v RegData) MarshalEasyJSON(w *jwriter.Writer) {
	easyjsonD2b7633eEncodeKudagoApplicationModels6(w, v)
}

// UnmarshalJSON supports json.Unmarshaler interface
func (v *RegData) UnmarshalJSON(data []byte) error {
	r := jlexer.Lexer{Data: data}
	easyjsonD2b7633eDecodeKudagoApplicationModels6(&r, v)
	return r.Error()
}

// UnmarshalEasyJSON supports easyjson.Unmarshaler interface
func (v *RegData) UnmarshalEasyJSON(l *jlexer.Lexer) {
	easyjsonD2b7633eDecodeKudagoApplicationModels6(l, v)
}
func easyjsonD2b7633eDecodeKudagoApplicationModels7(in *jlexer.Lexer, out *OtherUserProfile) {
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
		case "about":
			out.About = string(in.String())
		case "avatar":
			out.Avatar = string(in.String())
		case "visited":
			(out.Visited).UnmarshalEasyJSON(in)
		case "planning":
			(out.Planning).UnmarshalEasyJSON(in)
		case "followers":
			if in.IsNull() {
				in.Skip()
				out.Followers = nil
			} else {
				in.Delim('[')
				if out.Followers == nil {
					if !in.IsDelim(']') {
						out.Followers = make([]uint64, 0, 8)
					} else {
						out.Followers = []uint64{}
					}
				} else {
					out.Followers = (out.Followers)[:0]
				}
				for !in.IsDelim(']') {
					var v7 uint64
					v7 = uint64(in.Uint64())
					out.Followers = append(out.Followers, v7)
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
func easyjsonD2b7633eEncodeKudagoApplicationModels7(out *jwriter.Writer, in OtherUserProfile) {
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
		const prefix string = ",\"visited\":"
		out.RawString(prefix)
		(in.Visited).MarshalEasyJSON(out)
	}
	{
		const prefix string = ",\"planning\":"
		out.RawString(prefix)
		(in.Planning).MarshalEasyJSON(out)
	}
	{
		const prefix string = ",\"followers\":"
		out.RawString(prefix)
		if in.Followers == nil && (out.Flags&jwriter.NilSliceAsEmpty) == 0 {
			out.RawString("null")
		} else {
			out.RawByte('[')
			for v8, v9 := range in.Followers {
				if v8 > 0 {
					out.RawByte(',')
				}
				out.Uint64(uint64(v9))
			}
			out.RawByte(']')
		}
	}
	out.RawByte('}')
}

// MarshalJSON supports json.Marshaler interface
func (v OtherUserProfile) MarshalJSON() ([]byte, error) {
	w := jwriter.Writer{}
	easyjsonD2b7633eEncodeKudagoApplicationModels7(&w, v)
	return w.Buffer.BuildBytes(), w.Error
}

// MarshalEasyJSON supports easyjson.Marshaler interface
func (v OtherUserProfile) MarshalEasyJSON(w *jwriter.Writer) {
	easyjsonD2b7633eEncodeKudagoApplicationModels7(w, v)
}

// UnmarshalJSON supports json.Unmarshaler interface
func (v *OtherUserProfile) UnmarshalJSON(data []byte) error {
	r := jlexer.Lexer{Data: data}
	easyjsonD2b7633eDecodeKudagoApplicationModels7(&r, v)
	return r.Error()
}

// UnmarshalEasyJSON supports easyjson.Unmarshaler interface
func (v *OtherUserProfile) UnmarshalEasyJSON(l *jlexer.Lexer) {
	easyjsonD2b7633eDecodeKudagoApplicationModels7(l, v)
}
func easyjsonD2b7633eDecodeKudagoApplicationModels8(in *jlexer.Lexer, out *Events) {
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
			var v10 Event
			(v10).UnmarshalEasyJSON(in)
			*out = append(*out, v10)
			in.WantComma()
		}
		in.Delim(']')
	}
	if isTopLevel {
		in.Consumed()
	}
}
func easyjsonD2b7633eEncodeKudagoApplicationModels8(out *jwriter.Writer, in Events) {
	if in == nil && (out.Flags&jwriter.NilSliceAsEmpty) == 0 {
		out.RawString("null")
	} else {
		out.RawByte('[')
		for v11, v12 := range in {
			if v11 > 0 {
				out.RawByte(',')
			}
			(v12).MarshalEasyJSON(out)
		}
		out.RawByte(']')
	}
}

// MarshalJSON supports json.Marshaler interface
func (v Events) MarshalJSON() ([]byte, error) {
	w := jwriter.Writer{}
	easyjsonD2b7633eEncodeKudagoApplicationModels8(&w, v)
	return w.Buffer.BuildBytes(), w.Error
}

// MarshalEasyJSON supports easyjson.Marshaler interface
func (v Events) MarshalEasyJSON(w *jwriter.Writer) {
	easyjsonD2b7633eEncodeKudagoApplicationModels8(w, v)
}

// UnmarshalJSON supports json.Unmarshaler interface
func (v *Events) UnmarshalJSON(data []byte) error {
	r := jlexer.Lexer{Data: data}
	easyjsonD2b7633eDecodeKudagoApplicationModels8(&r, v)
	return r.Error()
}

// UnmarshalEasyJSON supports easyjson.Unmarshaler interface
func (v *Events) UnmarshalEasyJSON(l *jlexer.Lexer) {
	easyjsonD2b7633eDecodeKudagoApplicationModels8(l, v)
}
func easyjsonD2b7633eDecodeKudagoApplicationModels9(in *jlexer.Lexer, out *EventSQL) {
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
		case "ID":
			out.ID = uint64(in.Uint64())
		case "Title":
			out.Title = string(in.String())
		case "Place":
			out.Place = string(in.String())
		case "Description":
			out.Description = string(in.String())
		case "StartDate":
			easyjsonD2b7633eDecodeDatabaseSql1(in, &out.StartDate)
		case "EndDate":
			easyjsonD2b7633eDecodeDatabaseSql1(in, &out.EndDate)
		case "Subway":
			easyjsonD2b7633eDecodeDatabaseSql(in, &out.Subway)
		case "Street":
			easyjsonD2b7633eDecodeDatabaseSql(in, &out.Street)
		case "Category":
			out.Category = string(in.String())
		case "Image":
			easyjsonD2b7633eDecodeDatabaseSql(in, &out.Image)
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
func easyjsonD2b7633eEncodeKudagoApplicationModels9(out *jwriter.Writer, in EventSQL) {
	out.RawByte('{')
	first := true
	_ = first
	{
		const prefix string = ",\"ID\":"
		out.RawString(prefix[1:])
		out.Uint64(uint64(in.ID))
	}
	{
		const prefix string = ",\"Title\":"
		out.RawString(prefix)
		out.String(string(in.Title))
	}
	{
		const prefix string = ",\"Place\":"
		out.RawString(prefix)
		out.String(string(in.Place))
	}
	{
		const prefix string = ",\"Description\":"
		out.RawString(prefix)
		out.String(string(in.Description))
	}
	{
		const prefix string = ",\"StartDate\":"
		out.RawString(prefix)
		easyjsonD2b7633eEncodeDatabaseSql1(out, in.StartDate)
	}
	{
		const prefix string = ",\"EndDate\":"
		out.RawString(prefix)
		easyjsonD2b7633eEncodeDatabaseSql1(out, in.EndDate)
	}
	{
		const prefix string = ",\"Subway\":"
		out.RawString(prefix)
		easyjsonD2b7633eEncodeDatabaseSql(out, in.Subway)
	}
	{
		const prefix string = ",\"Street\":"
		out.RawString(prefix)
		easyjsonD2b7633eEncodeDatabaseSql(out, in.Street)
	}
	{
		const prefix string = ",\"Category\":"
		out.RawString(prefix)
		out.String(string(in.Category))
	}
	{
		const prefix string = ",\"Image\":"
		out.RawString(prefix)
		easyjsonD2b7633eEncodeDatabaseSql(out, in.Image)
	}
	out.RawByte('}')
}

// MarshalJSON supports json.Marshaler interface
func (v EventSQL) MarshalJSON() ([]byte, error) {
	w := jwriter.Writer{}
	easyjsonD2b7633eEncodeKudagoApplicationModels9(&w, v)
	return w.Buffer.BuildBytes(), w.Error
}

// MarshalEasyJSON supports easyjson.Marshaler interface
func (v EventSQL) MarshalEasyJSON(w *jwriter.Writer) {
	easyjsonD2b7633eEncodeKudagoApplicationModels9(w, v)
}

// UnmarshalJSON supports json.Unmarshaler interface
func (v *EventSQL) UnmarshalJSON(data []byte) error {
	r := jlexer.Lexer{Data: data}
	easyjsonD2b7633eDecodeKudagoApplicationModels9(&r, v)
	return r.Error()
}

// UnmarshalEasyJSON supports easyjson.Unmarshaler interface
func (v *EventSQL) UnmarshalEasyJSON(l *jlexer.Lexer) {
	easyjsonD2b7633eDecodeKudagoApplicationModels9(l, v)
}
func easyjsonD2b7633eDecodeKudagoApplicationModels10(in *jlexer.Lexer, out *EventCards) {
	isTopLevel := in.IsStart()
	if in.IsNull() {
		in.Skip()
		*out = nil
	} else {
		in.Delim('[')
		if *out == nil {
			if !in.IsDelim(']') {
				*out = make(EventCards, 0, 0)
			} else {
				*out = EventCards{}
			}
		} else {
			*out = (*out)[:0]
		}
		for !in.IsDelim(']') {
			var v13 EventCard
			(v13).UnmarshalEasyJSON(in)
			*out = append(*out, v13)
			in.WantComma()
		}
		in.Delim(']')
	}
	if isTopLevel {
		in.Consumed()
	}
}
func easyjsonD2b7633eEncodeKudagoApplicationModels10(out *jwriter.Writer, in EventCards) {
	if in == nil && (out.Flags&jwriter.NilSliceAsEmpty) == 0 {
		out.RawString("null")
	} else {
		out.RawByte('[')
		for v14, v15 := range in {
			if v14 > 0 {
				out.RawByte(',')
			}
			(v15).MarshalEasyJSON(out)
		}
		out.RawByte(']')
	}
}

// MarshalJSON supports json.Marshaler interface
func (v EventCards) MarshalJSON() ([]byte, error) {
	w := jwriter.Writer{}
	easyjsonD2b7633eEncodeKudagoApplicationModels10(&w, v)
	return w.Buffer.BuildBytes(), w.Error
}

// MarshalEasyJSON supports easyjson.Marshaler interface
func (v EventCards) MarshalEasyJSON(w *jwriter.Writer) {
	easyjsonD2b7633eEncodeKudagoApplicationModels10(w, v)
}

// UnmarshalJSON supports json.Unmarshaler interface
func (v *EventCards) UnmarshalJSON(data []byte) error {
	r := jlexer.Lexer{Data: data}
	easyjsonD2b7633eDecodeKudagoApplicationModels10(&r, v)
	return r.Error()
}

// UnmarshalEasyJSON supports easyjson.Unmarshaler interface
func (v *EventCards) UnmarshalEasyJSON(l *jlexer.Lexer) {
	easyjsonD2b7633eDecodeKudagoApplicationModels10(l, v)
}
func easyjsonD2b7633eDecodeKudagoApplicationModels11(in *jlexer.Lexer, out *EventCardWithDateSQL) {
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
		case "ID":
			out.ID = uint64(in.Uint64())
		case "Title":
			out.Title = string(in.String())
		case "Description":
			out.Description = string(in.String())
		case "Image":
			easyjsonD2b7633eDecodeDatabaseSql(in, &out.Image)
		case "StartDate":
			if data := in.Raw(); in.Ok() {
				in.AddError((out.StartDate).UnmarshalJSON(data))
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
func easyjsonD2b7633eEncodeKudagoApplicationModels11(out *jwriter.Writer, in EventCardWithDateSQL) {
	out.RawByte('{')
	first := true
	_ = first
	{
		const prefix string = ",\"ID\":"
		out.RawString(prefix[1:])
		out.Uint64(uint64(in.ID))
	}
	{
		const prefix string = ",\"Title\":"
		out.RawString(prefix)
		out.String(string(in.Title))
	}
	{
		const prefix string = ",\"Description\":"
		out.RawString(prefix)
		out.String(string(in.Description))
	}
	{
		const prefix string = ",\"Image\":"
		out.RawString(prefix)
		easyjsonD2b7633eEncodeDatabaseSql(out, in.Image)
	}
	{
		const prefix string = ",\"StartDate\":"
		out.RawString(prefix)
		out.Raw((in.StartDate).MarshalJSON())
	}
	out.RawByte('}')
}

// MarshalJSON supports json.Marshaler interface
func (v EventCardWithDateSQL) MarshalJSON() ([]byte, error) {
	w := jwriter.Writer{}
	easyjsonD2b7633eEncodeKudagoApplicationModels11(&w, v)
	return w.Buffer.BuildBytes(), w.Error
}

// MarshalEasyJSON supports easyjson.Marshaler interface
func (v EventCardWithDateSQL) MarshalEasyJSON(w *jwriter.Writer) {
	easyjsonD2b7633eEncodeKudagoApplicationModels11(w, v)
}

// UnmarshalJSON supports json.Unmarshaler interface
func (v *EventCardWithDateSQL) UnmarshalJSON(data []byte) error {
	r := jlexer.Lexer{Data: data}
	easyjsonD2b7633eDecodeKudagoApplicationModels11(&r, v)
	return r.Error()
}

// UnmarshalEasyJSON supports easyjson.Unmarshaler interface
func (v *EventCardWithDateSQL) UnmarshalEasyJSON(l *jlexer.Lexer) {
	easyjsonD2b7633eDecodeKudagoApplicationModels11(l, v)
}
func easyjsonD2b7633eDecodeKudagoApplicationModels12(in *jlexer.Lexer, out *EventCard) {
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
		case "description":
			out.Description = string(in.String())
		case "image":
			out.Image = string(in.String())
		case "startDate":
			out.StartDate = string(in.String())
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
func easyjsonD2b7633eEncodeKudagoApplicationModels12(out *jwriter.Writer, in EventCard) {
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
		const prefix string = ",\"description\":"
		out.RawString(prefix)
		out.String(string(in.Description))
	}
	{
		const prefix string = ",\"image\":"
		out.RawString(prefix)
		out.String(string(in.Image))
	}
	{
		const prefix string = ",\"startDate\":"
		out.RawString(prefix)
		out.String(string(in.StartDate))
	}
	out.RawByte('}')
}

// MarshalJSON supports json.Marshaler interface
func (v EventCard) MarshalJSON() ([]byte, error) {
	w := jwriter.Writer{}
	easyjsonD2b7633eEncodeKudagoApplicationModels12(&w, v)
	return w.Buffer.BuildBytes(), w.Error
}

// MarshalEasyJSON supports easyjson.Marshaler interface
func (v EventCard) MarshalEasyJSON(w *jwriter.Writer) {
	easyjsonD2b7633eEncodeKudagoApplicationModels12(w, v)
}

// UnmarshalJSON supports json.Unmarshaler interface
func (v *EventCard) UnmarshalJSON(data []byte) error {
	r := jlexer.Lexer{Data: data}
	easyjsonD2b7633eDecodeKudagoApplicationModels12(&r, v)
	return r.Error()
}

// UnmarshalEasyJSON supports easyjson.Unmarshaler interface
func (v *EventCard) UnmarshalEasyJSON(l *jlexer.Lexer) {
	easyjsonD2b7633eDecodeKudagoApplicationModels12(l, v)
}
func easyjsonD2b7633eDecodeKudagoApplicationModels13(in *jlexer.Lexer, out *Event) {
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
		case "startDate":
			out.StartDate = string(in.String())
		case "endDate":
			out.EndDate = string(in.String())
		case "subway":
			out.Subway = string(in.String())
		case "street":
			out.Street = string(in.String())
		case "tags":
			(out.Tags).UnmarshalEasyJSON(in)
		case "category":
			out.Category = string(in.String())
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
func easyjsonD2b7633eEncodeKudagoApplicationModels13(out *jwriter.Writer, in Event) {
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
		const prefix string = ",\"startDate\":"
		out.RawString(prefix)
		out.String(string(in.StartDate))
	}
	{
		const prefix string = ",\"endDate\":"
		out.RawString(prefix)
		out.String(string(in.EndDate))
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
		const prefix string = ",\"tags\":"
		out.RawString(prefix)
		(in.Tags).MarshalEasyJSON(out)
	}
	{
		const prefix string = ",\"category\":"
		out.RawString(prefix)
		out.String(string(in.Category))
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
	easyjsonD2b7633eEncodeKudagoApplicationModels13(&w, v)
	return w.Buffer.BuildBytes(), w.Error
}

// MarshalEasyJSON supports easyjson.Marshaler interface
func (v Event) MarshalEasyJSON(w *jwriter.Writer) {
	easyjsonD2b7633eEncodeKudagoApplicationModels13(w, v)
}

// UnmarshalJSON supports json.Unmarshaler interface
func (v *Event) UnmarshalJSON(data []byte) error {
	r := jlexer.Lexer{Data: data}
	easyjsonD2b7633eDecodeKudagoApplicationModels13(&r, v)
	return r.Error()
}

// UnmarshalEasyJSON supports easyjson.Unmarshaler interface
func (v *Event) UnmarshalEasyJSON(l *jlexer.Lexer) {
	easyjsonD2b7633eDecodeKudagoApplicationModels13(l, v)
}