// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.25.0-devel
// 	protoc        v3.14.0
// source: chat.proto

package proto

import (
	protoreflect "google.golang.org/protobuf/reflect/protoreflect"
	protoimpl "google.golang.org/protobuf/runtime/protoimpl"
	reflect "reflect"
	sync "sync"
)

const (
	// Verify that this generated code is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(20 - protoimpl.MinVersion)
	// Verify that runtime/protoimpl is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(protoimpl.MaxVersion - 20)
)

type Ids struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	List []uint64 `protobuf:"varint,1,rep,packed,name=list,proto3" json:"list,omitempty"`
}

func (x *Ids) Reset() {
	*x = Ids{}
	if protoimpl.UnsafeEnabled {
		mi := &file_chat_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Ids) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Ids) ProtoMessage() {}

func (x *Ids) ProtoReflect() protoreflect.Message {
	mi := &file_chat_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Ids.ProtoReflect.Descriptor instead.
func (*Ids) Descriptor() ([]byte, []int) {
	return file_chat_proto_rawDescGZIP(), []int{0}
}

func (x *Ids) GetList() []uint64 {
	if x != nil {
		return x.List
	}
	return nil
}

type MailingIn struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	UserId  uint64 `protobuf:"varint,1,opt,name=userId,proto3" json:"userId,omitempty"`
	EventId uint64 `protobuf:"varint,2,opt,name=eventId,proto3" json:"eventId,omitempty"`
	To      *Ids   `protobuf:"bytes,3,opt,name=to,proto3" json:"to,omitempty"`
}

func (x *MailingIn) Reset() {
	*x = MailingIn{}
	if protoimpl.UnsafeEnabled {
		mi := &file_chat_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *MailingIn) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*MailingIn) ProtoMessage() {}

func (x *MailingIn) ProtoReflect() protoreflect.Message {
	mi := &file_chat_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use MailingIn.ProtoReflect.Descriptor instead.
func (*MailingIn) Descriptor() ([]byte, []int) {
	return file_chat_proto_rawDescGZIP(), []int{1}
}

func (x *MailingIn) GetUserId() uint64 {
	if x != nil {
		return x.UserId
	}
	return 0
}

func (x *MailingIn) GetEventId() uint64 {
	if x != nil {
		return x.EventId
	}
	return 0
}

func (x *MailingIn) GetTo() *Ids {
	if x != nil {
		return x.To
	}
	return nil
}

type SearchIn struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Uid  uint64 `protobuf:"varint,1,opt,name=uid,proto3" json:"uid,omitempty"`
	Id   int32  `protobuf:"varint,2,opt,name=id,proto3" json:"id,omitempty"`
	Str  string `protobuf:"bytes,3,opt,name=str,proto3" json:"str,omitempty"`
	Page int32  `protobuf:"varint,4,opt,name=page,proto3" json:"page,omitempty"`
}

func (x *SearchIn) Reset() {
	*x = SearchIn{}
	if protoimpl.UnsafeEnabled {
		mi := &file_chat_proto_msgTypes[2]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *SearchIn) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*SearchIn) ProtoMessage() {}

func (x *SearchIn) ProtoReflect() protoreflect.Message {
	mi := &file_chat_proto_msgTypes[2]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use SearchIn.ProtoReflect.Descriptor instead.
func (*SearchIn) Descriptor() ([]byte, []int) {
	return file_chat_proto_rawDescGZIP(), []int{2}
}

func (x *SearchIn) GetUid() uint64 {
	if x != nil {
		return x.Uid
	}
	return 0
}

func (x *SearchIn) GetId() int32 {
	if x != nil {
		return x.Id
	}
	return 0
}

func (x *SearchIn) GetStr() string {
	if x != nil {
		return x.Str
	}
	return ""
}

func (x *SearchIn) GetPage() int32 {
	if x != nil {
		return x.Page
	}
	return 0
}

type IdPage struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Id   uint64 `protobuf:"varint,1,opt,name=id,proto3" json:"id,omitempty"`
	Page int32  `protobuf:"varint,2,opt,name=page,proto3" json:"page,omitempty"`
}

func (x *IdPage) Reset() {
	*x = IdPage{}
	if protoimpl.UnsafeEnabled {
		mi := &file_chat_proto_msgTypes[3]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *IdPage) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*IdPage) ProtoMessage() {}

func (x *IdPage) ProtoReflect() protoreflect.Message {
	mi := &file_chat_proto_msgTypes[3]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use IdPage.ProtoReflect.Descriptor instead.
func (*IdPage) Descriptor() ([]byte, []int) {
	return file_chat_proto_rawDescGZIP(), []int{3}
}

func (x *IdPage) GetId() uint64 {
	if x != nil {
		return x.Id
	}
	return 0
}

func (x *IdPage) GetPage() int32 {
	if x != nil {
		return x.Page
	}
	return 0
}

type SendEditMessage struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Id1  uint64 `protobuf:"varint,1,opt,name=id1,proto3" json:"id1,omitempty"`
	Text string `protobuf:"bytes,2,opt,name=text,proto3" json:"text,omitempty"`
	Id2  uint64 `protobuf:"varint,3,opt,name=id2,proto3" json:"id2,omitempty"`
}

func (x *SendEditMessage) Reset() {
	*x = SendEditMessage{}
	if protoimpl.UnsafeEnabled {
		mi := &file_chat_proto_msgTypes[4]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *SendEditMessage) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*SendEditMessage) ProtoMessage() {}

func (x *SendEditMessage) ProtoReflect() protoreflect.Message {
	mi := &file_chat_proto_msgTypes[4]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use SendEditMessage.ProtoReflect.Descriptor instead.
func (*SendEditMessage) Descriptor() ([]byte, []int) {
	return file_chat_proto_rawDescGZIP(), []int{4}
}

func (x *SendEditMessage) GetId1() uint64 {
	if x != nil {
		return x.Id1
	}
	return 0
}

func (x *SendEditMessage) GetText() string {
	if x != nil {
		return x.Text
	}
	return ""
}

func (x *SendEditMessage) GetId2() uint64 {
	if x != nil {
		return x.Id2
	}
	return 0
}

type IdId struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Id1 uint64 `protobuf:"varint,1,opt,name=id1,proto3" json:"id1,omitempty"`
	Id2 uint64 `protobuf:"varint,2,opt,name=id2,proto3" json:"id2,omitempty"`
}

func (x *IdId) Reset() {
	*x = IdId{}
	if protoimpl.UnsafeEnabled {
		mi := &file_chat_proto_msgTypes[5]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *IdId) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*IdId) ProtoMessage() {}

func (x *IdId) ProtoReflect() protoreflect.Message {
	mi := &file_chat_proto_msgTypes[5]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use IdId.ProtoReflect.Descriptor instead.
func (*IdId) Descriptor() ([]byte, []int) {
	return file_chat_proto_rawDescGZIP(), []int{5}
}

func (x *IdId) GetId1() uint64 {
	if x != nil {
		return x.Id1
	}
	return 0
}

func (x *IdId) GetId2() uint64 {
	if x != nil {
		return x.Id2
	}
	return 0
}

type IdIdPage struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Id1  uint64 `protobuf:"varint,1,opt,name=id1,proto3" json:"id1,omitempty"`
	Id2  uint64 `protobuf:"varint,2,opt,name=id2,proto3" json:"id2,omitempty"`
	Page int32  `protobuf:"varint,3,opt,name=page,proto3" json:"page,omitempty"`
}

func (x *IdIdPage) Reset() {
	*x = IdIdPage{}
	if protoimpl.UnsafeEnabled {
		mi := &file_chat_proto_msgTypes[6]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *IdIdPage) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*IdIdPage) ProtoMessage() {}

func (x *IdIdPage) ProtoReflect() protoreflect.Message {
	mi := &file_chat_proto_msgTypes[6]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use IdIdPage.ProtoReflect.Descriptor instead.
func (*IdIdPage) Descriptor() ([]byte, []int) {
	return file_chat_proto_rawDescGZIP(), []int{6}
}

func (x *IdIdPage) GetId1() uint64 {
	if x != nil {
		return x.Id1
	}
	return 0
}

func (x *IdIdPage) GetId2() uint64 {
	if x != nil {
		return x.Id2
	}
	return 0
}

func (x *IdIdPage) GetPage() int32 {
	if x != nil {
		return x.Page
	}
	return 0
}

type DialogueCard struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	ID           uint64       `protobuf:"varint,1,opt,name=ID,proto3" json:"ID,omitempty"`
	Interlocutor *UserOnEvent `protobuf:"bytes,2,opt,name=interlocutor,proto3" json:"interlocutor,omitempty"`
	LastMessage  *Message     `protobuf:"bytes,3,opt,name=lastMessage,proto3" json:"lastMessage,omitempty"`
}

func (x *DialogueCard) Reset() {
	*x = DialogueCard{}
	if protoimpl.UnsafeEnabled {
		mi := &file_chat_proto_msgTypes[7]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *DialogueCard) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*DialogueCard) ProtoMessage() {}

func (x *DialogueCard) ProtoReflect() protoreflect.Message {
	mi := &file_chat_proto_msgTypes[7]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use DialogueCard.ProtoReflect.Descriptor instead.
func (*DialogueCard) Descriptor() ([]byte, []int) {
	return file_chat_proto_rawDescGZIP(), []int{7}
}

func (x *DialogueCard) GetID() uint64 {
	if x != nil {
		return x.ID
	}
	return 0
}

func (x *DialogueCard) GetInterlocutor() *UserOnEvent {
	if x != nil {
		return x.Interlocutor
	}
	return nil
}

func (x *DialogueCard) GetLastMessage() *Message {
	if x != nil {
		return x.LastMessage
	}
	return nil
}

type DialogueCards struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	List []*DialogueCard `protobuf:"bytes,1,rep,name=list,proto3" json:"list,omitempty"`
}

func (x *DialogueCards) Reset() {
	*x = DialogueCards{}
	if protoimpl.UnsafeEnabled {
		mi := &file_chat_proto_msgTypes[8]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *DialogueCards) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*DialogueCards) ProtoMessage() {}

func (x *DialogueCards) ProtoReflect() protoreflect.Message {
	mi := &file_chat_proto_msgTypes[8]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use DialogueCards.ProtoReflect.Descriptor instead.
func (*DialogueCards) Descriptor() ([]byte, []int) {
	return file_chat_proto_rawDescGZIP(), []int{8}
}

func (x *DialogueCards) GetList() []*DialogueCard {
	if x != nil {
		return x.List
	}
	return nil
}

type UserOnEvent struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Id     uint64 `protobuf:"varint,1,opt,name=id,proto3" json:"id,omitempty"`
	Name   string `protobuf:"bytes,2,opt,name=name,proto3" json:"name,omitempty"`
	Avatar string `protobuf:"bytes,3,opt,name=avatar,proto3" json:"avatar,omitempty"`
}

func (x *UserOnEvent) Reset() {
	*x = UserOnEvent{}
	if protoimpl.UnsafeEnabled {
		mi := &file_chat_proto_msgTypes[9]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *UserOnEvent) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*UserOnEvent) ProtoMessage() {}

func (x *UserOnEvent) ProtoReflect() protoreflect.Message {
	mi := &file_chat_proto_msgTypes[9]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use UserOnEvent.ProtoReflect.Descriptor instead.
func (*UserOnEvent) Descriptor() ([]byte, []int) {
	return file_chat_proto_rawDescGZIP(), []int{9}
}

func (x *UserOnEvent) GetId() uint64 {
	if x != nil {
		return x.Id
	}
	return 0
}

func (x *UserOnEvent) GetName() string {
	if x != nil {
		return x.Name
	}
	return ""
}

func (x *UserOnEvent) GetAvatar() string {
	if x != nil {
		return x.Avatar
	}
	return ""
}

type Message struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	ID     uint64 `protobuf:"varint,1,opt,name=ID,proto3" json:"ID,omitempty"`
	FromMe bool   `protobuf:"varint,2,opt,name=FromMe,proto3" json:"FromMe,omitempty"`
	Text   string `protobuf:"bytes,3,opt,name=Text,proto3" json:"Text,omitempty"`
	Date   string `protobuf:"bytes,4,opt,name=Date,proto3" json:"Date,omitempty"`
	Redact bool   `protobuf:"varint,5,opt,name=Redact,proto3" json:"Redact,omitempty"`
	Read   bool   `protobuf:"varint,6,opt,name=Read,proto3" json:"Read,omitempty"`
}

func (x *Message) Reset() {
	*x = Message{}
	if protoimpl.UnsafeEnabled {
		mi := &file_chat_proto_msgTypes[10]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Message) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Message) ProtoMessage() {}

func (x *Message) ProtoReflect() protoreflect.Message {
	mi := &file_chat_proto_msgTypes[10]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Message.ProtoReflect.Descriptor instead.
func (*Message) Descriptor() ([]byte, []int) {
	return file_chat_proto_rawDescGZIP(), []int{10}
}

func (x *Message) GetID() uint64 {
	if x != nil {
		return x.ID
	}
	return 0
}

func (x *Message) GetFromMe() bool {
	if x != nil {
		return x.FromMe
	}
	return false
}

func (x *Message) GetText() string {
	if x != nil {
		return x.Text
	}
	return ""
}

func (x *Message) GetDate() string {
	if x != nil {
		return x.Date
	}
	return ""
}

func (x *Message) GetRedact() bool {
	if x != nil {
		return x.Redact
	}
	return false
}

func (x *Message) GetRead() bool {
	if x != nil {
		return x.Read
	}
	return false
}

type Messages struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	List []*Message `protobuf:"bytes,1,rep,name=list,proto3" json:"list,omitempty"`
}

func (x *Messages) Reset() {
	*x = Messages{}
	if protoimpl.UnsafeEnabled {
		mi := &file_chat_proto_msgTypes[11]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Messages) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Messages) ProtoMessage() {}

func (x *Messages) ProtoReflect() protoreflect.Message {
	mi := &file_chat_proto_msgTypes[11]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Messages.ProtoReflect.Descriptor instead.
func (*Messages) Descriptor() ([]byte, []int) {
	return file_chat_proto_rawDescGZIP(), []int{11}
}

func (x *Messages) GetList() []*Message {
	if x != nil {
		return x.List
	}
	return nil
}

type Dialogue struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	ID             uint64       `protobuf:"varint,1,opt,name=ID,proto3" json:"ID,omitempty"`
	Interlocutor   *UserOnEvent `protobuf:"bytes,2,opt,name=interlocutor,proto3" json:"interlocutor,omitempty"`
	DialogMessages *Messages    `protobuf:"bytes,3,opt,name=dialogMessages,proto3" json:"dialogMessages,omitempty"`
}

func (x *Dialogue) Reset() {
	*x = Dialogue{}
	if protoimpl.UnsafeEnabled {
		mi := &file_chat_proto_msgTypes[12]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Dialogue) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Dialogue) ProtoMessage() {}

func (x *Dialogue) ProtoReflect() protoreflect.Message {
	mi := &file_chat_proto_msgTypes[12]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Dialogue.ProtoReflect.Descriptor instead.
func (*Dialogue) Descriptor() ([]byte, []int) {
	return file_chat_proto_rawDescGZIP(), []int{12}
}

func (x *Dialogue) GetID() uint64 {
	if x != nil {
		return x.ID
	}
	return 0
}

func (x *Dialogue) GetInterlocutor() *UserOnEvent {
	if x != nil {
		return x.Interlocutor
	}
	return nil
}

func (x *Dialogue) GetDialogMessages() *Messages {
	if x != nil {
		return x.DialogMessages
	}
	return nil
}

type Answer struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Flag bool   `protobuf:"varint,1,opt,name=flag,proto3" json:"flag,omitempty"`
	Msg  string `protobuf:"bytes,2,opt,name=msg,proto3" json:"msg,omitempty"`
}

func (x *Answer) Reset() {
	*x = Answer{}
	if protoimpl.UnsafeEnabled {
		mi := &file_chat_proto_msgTypes[13]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Answer) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Answer) ProtoMessage() {}

func (x *Answer) ProtoReflect() protoreflect.Message {
	mi := &file_chat_proto_msgTypes[13]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Answer.ProtoReflect.Descriptor instead.
func (*Answer) Descriptor() ([]byte, []int) {
	return file_chat_proto_rawDescGZIP(), []int{13}
}

func (x *Answer) GetFlag() bool {
	if x != nil {
		return x.Flag
	}
	return false
}

func (x *Answer) GetMsg() string {
	if x != nil {
		return x.Msg
	}
	return ""
}

var File_chat_proto protoreflect.FileDescriptor

var file_chat_proto_rawDesc = []byte{
	0x0a, 0x0a, 0x63, 0x68, 0x61, 0x74, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x22, 0x19, 0x0a, 0x03,
	0x49, 0x64, 0x73, 0x12, 0x12, 0x0a, 0x04, 0x6c, 0x69, 0x73, 0x74, 0x18, 0x01, 0x20, 0x03, 0x28,
	0x04, 0x52, 0x04, 0x6c, 0x69, 0x73, 0x74, 0x22, 0x53, 0x0a, 0x09, 0x4d, 0x61, 0x69, 0x6c, 0x69,
	0x6e, 0x67, 0x49, 0x6e, 0x12, 0x16, 0x0a, 0x06, 0x75, 0x73, 0x65, 0x72, 0x49, 0x64, 0x18, 0x01,
	0x20, 0x01, 0x28, 0x04, 0x52, 0x06, 0x75, 0x73, 0x65, 0x72, 0x49, 0x64, 0x12, 0x18, 0x0a, 0x07,
	0x65, 0x76, 0x65, 0x6e, 0x74, 0x49, 0x64, 0x18, 0x02, 0x20, 0x01, 0x28, 0x04, 0x52, 0x07, 0x65,
	0x76, 0x65, 0x6e, 0x74, 0x49, 0x64, 0x12, 0x14, 0x0a, 0x02, 0x74, 0x6f, 0x18, 0x03, 0x20, 0x01,
	0x28, 0x0b, 0x32, 0x04, 0x2e, 0x49, 0x64, 0x73, 0x52, 0x02, 0x74, 0x6f, 0x22, 0x52, 0x0a, 0x08,
	0x53, 0x65, 0x61, 0x72, 0x63, 0x68, 0x49, 0x6e, 0x12, 0x10, 0x0a, 0x03, 0x75, 0x69, 0x64, 0x18,
	0x01, 0x20, 0x01, 0x28, 0x04, 0x52, 0x03, 0x75, 0x69, 0x64, 0x12, 0x0e, 0x0a, 0x02, 0x69, 0x64,
	0x18, 0x02, 0x20, 0x01, 0x28, 0x05, 0x52, 0x02, 0x69, 0x64, 0x12, 0x10, 0x0a, 0x03, 0x73, 0x74,
	0x72, 0x18, 0x03, 0x20, 0x01, 0x28, 0x09, 0x52, 0x03, 0x73, 0x74, 0x72, 0x12, 0x12, 0x0a, 0x04,
	0x70, 0x61, 0x67, 0x65, 0x18, 0x04, 0x20, 0x01, 0x28, 0x05, 0x52, 0x04, 0x70, 0x61, 0x67, 0x65,
	0x22, 0x2c, 0x0a, 0x06, 0x49, 0x64, 0x50, 0x61, 0x67, 0x65, 0x12, 0x0e, 0x0a, 0x02, 0x69, 0x64,
	0x18, 0x01, 0x20, 0x01, 0x28, 0x04, 0x52, 0x02, 0x69, 0x64, 0x12, 0x12, 0x0a, 0x04, 0x70, 0x61,
	0x67, 0x65, 0x18, 0x02, 0x20, 0x01, 0x28, 0x05, 0x52, 0x04, 0x70, 0x61, 0x67, 0x65, 0x22, 0x49,
	0x0a, 0x0f, 0x53, 0x65, 0x6e, 0x64, 0x45, 0x64, 0x69, 0x74, 0x4d, 0x65, 0x73, 0x73, 0x61, 0x67,
	0x65, 0x12, 0x10, 0x0a, 0x03, 0x69, 0x64, 0x31, 0x18, 0x01, 0x20, 0x01, 0x28, 0x04, 0x52, 0x03,
	0x69, 0x64, 0x31, 0x12, 0x12, 0x0a, 0x04, 0x74, 0x65, 0x78, 0x74, 0x18, 0x02, 0x20, 0x01, 0x28,
	0x09, 0x52, 0x04, 0x74, 0x65, 0x78, 0x74, 0x12, 0x10, 0x0a, 0x03, 0x69, 0x64, 0x32, 0x18, 0x03,
	0x20, 0x01, 0x28, 0x04, 0x52, 0x03, 0x69, 0x64, 0x32, 0x22, 0x2a, 0x0a, 0x04, 0x49, 0x64, 0x49,
	0x64, 0x12, 0x10, 0x0a, 0x03, 0x69, 0x64, 0x31, 0x18, 0x01, 0x20, 0x01, 0x28, 0x04, 0x52, 0x03,
	0x69, 0x64, 0x31, 0x12, 0x10, 0x0a, 0x03, 0x69, 0x64, 0x32, 0x18, 0x02, 0x20, 0x01, 0x28, 0x04,
	0x52, 0x03, 0x69, 0x64, 0x32, 0x22, 0x42, 0x0a, 0x08, 0x49, 0x64, 0x49, 0x64, 0x50, 0x61, 0x67,
	0x65, 0x12, 0x10, 0x0a, 0x03, 0x69, 0x64, 0x31, 0x18, 0x01, 0x20, 0x01, 0x28, 0x04, 0x52, 0x03,
	0x69, 0x64, 0x31, 0x12, 0x10, 0x0a, 0x03, 0x69, 0x64, 0x32, 0x18, 0x02, 0x20, 0x01, 0x28, 0x04,
	0x52, 0x03, 0x69, 0x64, 0x32, 0x12, 0x12, 0x0a, 0x04, 0x70, 0x61, 0x67, 0x65, 0x18, 0x03, 0x20,
	0x01, 0x28, 0x05, 0x52, 0x04, 0x70, 0x61, 0x67, 0x65, 0x22, 0x7c, 0x0a, 0x0c, 0x44, 0x69, 0x61,
	0x6c, 0x6f, 0x67, 0x75, 0x65, 0x43, 0x61, 0x72, 0x64, 0x12, 0x0e, 0x0a, 0x02, 0x49, 0x44, 0x18,
	0x01, 0x20, 0x01, 0x28, 0x04, 0x52, 0x02, 0x49, 0x44, 0x12, 0x30, 0x0a, 0x0c, 0x69, 0x6e, 0x74,
	0x65, 0x72, 0x6c, 0x6f, 0x63, 0x75, 0x74, 0x6f, 0x72, 0x18, 0x02, 0x20, 0x01, 0x28, 0x0b, 0x32,
	0x0c, 0x2e, 0x55, 0x73, 0x65, 0x72, 0x4f, 0x6e, 0x45, 0x76, 0x65, 0x6e, 0x74, 0x52, 0x0c, 0x69,
	0x6e, 0x74, 0x65, 0x72, 0x6c, 0x6f, 0x63, 0x75, 0x74, 0x6f, 0x72, 0x12, 0x2a, 0x0a, 0x0b, 0x6c,
	0x61, 0x73, 0x74, 0x4d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x18, 0x03, 0x20, 0x01, 0x28, 0x0b,
	0x32, 0x08, 0x2e, 0x4d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x52, 0x0b, 0x6c, 0x61, 0x73, 0x74,
	0x4d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x22, 0x32, 0x0a, 0x0d, 0x44, 0x69, 0x61, 0x6c, 0x6f,
	0x67, 0x75, 0x65, 0x43, 0x61, 0x72, 0x64, 0x73, 0x12, 0x21, 0x0a, 0x04, 0x6c, 0x69, 0x73, 0x74,
	0x18, 0x01, 0x20, 0x03, 0x28, 0x0b, 0x32, 0x0d, 0x2e, 0x44, 0x69, 0x61, 0x6c, 0x6f, 0x67, 0x75,
	0x65, 0x43, 0x61, 0x72, 0x64, 0x52, 0x04, 0x6c, 0x69, 0x73, 0x74, 0x22, 0x49, 0x0a, 0x0b, 0x55,
	0x73, 0x65, 0x72, 0x4f, 0x6e, 0x45, 0x76, 0x65, 0x6e, 0x74, 0x12, 0x0e, 0x0a, 0x02, 0x69, 0x64,
	0x18, 0x01, 0x20, 0x01, 0x28, 0x04, 0x52, 0x02, 0x69, 0x64, 0x12, 0x12, 0x0a, 0x04, 0x6e, 0x61,
	0x6d, 0x65, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x04, 0x6e, 0x61, 0x6d, 0x65, 0x12, 0x16,
	0x0a, 0x06, 0x61, 0x76, 0x61, 0x74, 0x61, 0x72, 0x18, 0x03, 0x20, 0x01, 0x28, 0x09, 0x52, 0x06,
	0x61, 0x76, 0x61, 0x74, 0x61, 0x72, 0x22, 0x85, 0x01, 0x0a, 0x07, 0x4d, 0x65, 0x73, 0x73, 0x61,
	0x67, 0x65, 0x12, 0x0e, 0x0a, 0x02, 0x49, 0x44, 0x18, 0x01, 0x20, 0x01, 0x28, 0x04, 0x52, 0x02,
	0x49, 0x44, 0x12, 0x16, 0x0a, 0x06, 0x46, 0x72, 0x6f, 0x6d, 0x4d, 0x65, 0x18, 0x02, 0x20, 0x01,
	0x28, 0x08, 0x52, 0x06, 0x46, 0x72, 0x6f, 0x6d, 0x4d, 0x65, 0x12, 0x12, 0x0a, 0x04, 0x54, 0x65,
	0x78, 0x74, 0x18, 0x03, 0x20, 0x01, 0x28, 0x09, 0x52, 0x04, 0x54, 0x65, 0x78, 0x74, 0x12, 0x12,
	0x0a, 0x04, 0x44, 0x61, 0x74, 0x65, 0x18, 0x04, 0x20, 0x01, 0x28, 0x09, 0x52, 0x04, 0x44, 0x61,
	0x74, 0x65, 0x12, 0x16, 0x0a, 0x06, 0x52, 0x65, 0x64, 0x61, 0x63, 0x74, 0x18, 0x05, 0x20, 0x01,
	0x28, 0x08, 0x52, 0x06, 0x52, 0x65, 0x64, 0x61, 0x63, 0x74, 0x12, 0x12, 0x0a, 0x04, 0x52, 0x65,
	0x61, 0x64, 0x18, 0x06, 0x20, 0x01, 0x28, 0x08, 0x52, 0x04, 0x52, 0x65, 0x61, 0x64, 0x22, 0x28,
	0x0a, 0x08, 0x4d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x73, 0x12, 0x1c, 0x0a, 0x04, 0x6c, 0x69,
	0x73, 0x74, 0x18, 0x01, 0x20, 0x03, 0x28, 0x0b, 0x32, 0x08, 0x2e, 0x4d, 0x65, 0x73, 0x73, 0x61,
	0x67, 0x65, 0x52, 0x04, 0x6c, 0x69, 0x73, 0x74, 0x22, 0x7f, 0x0a, 0x08, 0x44, 0x69, 0x61, 0x6c,
	0x6f, 0x67, 0x75, 0x65, 0x12, 0x0e, 0x0a, 0x02, 0x49, 0x44, 0x18, 0x01, 0x20, 0x01, 0x28, 0x04,
	0x52, 0x02, 0x49, 0x44, 0x12, 0x30, 0x0a, 0x0c, 0x69, 0x6e, 0x74, 0x65, 0x72, 0x6c, 0x6f, 0x63,
	0x75, 0x74, 0x6f, 0x72, 0x18, 0x02, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x0c, 0x2e, 0x55, 0x73, 0x65,
	0x72, 0x4f, 0x6e, 0x45, 0x76, 0x65, 0x6e, 0x74, 0x52, 0x0c, 0x69, 0x6e, 0x74, 0x65, 0x72, 0x6c,
	0x6f, 0x63, 0x75, 0x74, 0x6f, 0x72, 0x12, 0x31, 0x0a, 0x0e, 0x64, 0x69, 0x61, 0x6c, 0x6f, 0x67,
	0x4d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x73, 0x18, 0x03, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x09,
	0x2e, 0x4d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x73, 0x52, 0x0e, 0x64, 0x69, 0x61, 0x6c, 0x6f,
	0x67, 0x4d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x73, 0x22, 0x2e, 0x0a, 0x06, 0x41, 0x6e, 0x73,
	0x77, 0x65, 0x72, 0x12, 0x12, 0x0a, 0x04, 0x66, 0x6c, 0x61, 0x67, 0x18, 0x01, 0x20, 0x01, 0x28,
	0x08, 0x52, 0x04, 0x66, 0x6c, 0x61, 0x67, 0x12, 0x10, 0x0a, 0x03, 0x6d, 0x73, 0x67, 0x18, 0x02,
	0x20, 0x01, 0x28, 0x09, 0x52, 0x03, 0x6d, 0x73, 0x67, 0x32, 0xc6, 0x02, 0x0a, 0x04, 0x43, 0x68,
	0x61, 0x74, 0x12, 0x2c, 0x0a, 0x0f, 0x47, 0x65, 0x74, 0x41, 0x6c, 0x6c, 0x44, 0x69, 0x61, 0x6c,
	0x6f, 0x67, 0x75, 0x65, 0x73, 0x12, 0x07, 0x2e, 0x49, 0x64, 0x50, 0x61, 0x67, 0x65, 0x1a, 0x0e,
	0x2e, 0x44, 0x69, 0x61, 0x6c, 0x6f, 0x67, 0x75, 0x65, 0x43, 0x61, 0x72, 0x64, 0x73, 0x22, 0x00,
	0x12, 0x28, 0x0a, 0x0e, 0x47, 0x65, 0x74, 0x4f, 0x6e, 0x65, 0x44, 0x69, 0x61, 0x6c, 0x6f, 0x67,
	0x75, 0x65, 0x12, 0x09, 0x2e, 0x49, 0x64, 0x49, 0x64, 0x50, 0x61, 0x67, 0x65, 0x1a, 0x09, 0x2e,
	0x44, 0x69, 0x61, 0x6c, 0x6f, 0x67, 0x75, 0x65, 0x22, 0x00, 0x12, 0x22, 0x0a, 0x0e, 0x44, 0x65,
	0x6c, 0x65, 0x74, 0x65, 0x44, 0x69, 0x61, 0x6c, 0x6f, 0x67, 0x75, 0x65, 0x12, 0x05, 0x2e, 0x49,
	0x64, 0x49, 0x64, 0x1a, 0x07, 0x2e, 0x41, 0x6e, 0x73, 0x77, 0x65, 0x72, 0x22, 0x00, 0x12, 0x2a,
	0x0a, 0x0b, 0x53, 0x65, 0x6e, 0x64, 0x4d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x12, 0x10, 0x2e,
	0x53, 0x65, 0x6e, 0x64, 0x45, 0x64, 0x69, 0x74, 0x4d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x1a,
	0x07, 0x2e, 0x41, 0x6e, 0x73, 0x77, 0x65, 0x72, 0x22, 0x00, 0x12, 0x2a, 0x0a, 0x0b, 0x45, 0x64,
	0x69, 0x74, 0x4d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x12, 0x10, 0x2e, 0x53, 0x65, 0x6e, 0x64,
	0x45, 0x64, 0x69, 0x74, 0x4d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x1a, 0x07, 0x2e, 0x41, 0x6e,
	0x73, 0x77, 0x65, 0x72, 0x22, 0x00, 0x12, 0x21, 0x0a, 0x0d, 0x44, 0x65, 0x6c, 0x65, 0x74, 0x65,
	0x4d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x12, 0x05, 0x2e, 0x49, 0x64, 0x49, 0x64, 0x1a, 0x07,
	0x2e, 0x41, 0x6e, 0x73, 0x77, 0x65, 0x72, 0x22, 0x00, 0x12, 0x20, 0x0a, 0x07, 0x4d, 0x61, 0x69,
	0x6c, 0x69, 0x6e, 0x67, 0x12, 0x0a, 0x2e, 0x4d, 0x61, 0x69, 0x6c, 0x69, 0x6e, 0x67, 0x49, 0x6e,
	0x1a, 0x07, 0x2e, 0x41, 0x6e, 0x73, 0x77, 0x65, 0x72, 0x22, 0x00, 0x12, 0x25, 0x0a, 0x06, 0x53,
	0x65, 0x61, 0x72, 0x63, 0x68, 0x12, 0x09, 0x2e, 0x53, 0x65, 0x61, 0x72, 0x63, 0x68, 0x49, 0x6e,
	0x1a, 0x0e, 0x2e, 0x44, 0x69, 0x61, 0x6c, 0x6f, 0x67, 0x75, 0x65, 0x43, 0x61, 0x72, 0x64, 0x73,
	0x22, 0x00, 0x42, 0x0a, 0x5a, 0x08, 0x2e, 0x2e, 0x2f, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x06,
	0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_chat_proto_rawDescOnce sync.Once
	file_chat_proto_rawDescData = file_chat_proto_rawDesc
)

func file_chat_proto_rawDescGZIP() []byte {
	file_chat_proto_rawDescOnce.Do(func() {
		file_chat_proto_rawDescData = protoimpl.X.CompressGZIP(file_chat_proto_rawDescData)
	})
	return file_chat_proto_rawDescData
}

var file_chat_proto_msgTypes = make([]protoimpl.MessageInfo, 14)
var file_chat_proto_goTypes = []interface{}{
	(*Ids)(nil),             // 0: Ids
	(*MailingIn)(nil),       // 1: MailingIn
	(*SearchIn)(nil),        // 2: SearchIn
	(*IdPage)(nil),          // 3: IdPage
	(*SendEditMessage)(nil), // 4: SendEditMessage
	(*IdId)(nil),            // 5: IdId
	(*IdIdPage)(nil),        // 6: IdIdPage
	(*DialogueCard)(nil),    // 7: DialogueCard
	(*DialogueCards)(nil),   // 8: DialogueCards
	(*UserOnEvent)(nil),     // 9: UserOnEvent
	(*Message)(nil),         // 10: Message
	(*Messages)(nil),        // 11: Messages
	(*Dialogue)(nil),        // 12: Dialogue
	(*Answer)(nil),          // 13: Answer
}
var file_chat_proto_depIdxs = []int32{
	0,  // 0: MailingIn.to:type_name -> Ids
	9,  // 1: DialogueCard.interlocutor:type_name -> UserOnEvent
	10, // 2: DialogueCard.lastMessage:type_name -> Message
	7,  // 3: DialogueCards.list:type_name -> DialogueCard
	10, // 4: Messages.list:type_name -> Message
	9,  // 5: Dialogue.interlocutor:type_name -> UserOnEvent
	11, // 6: Dialogue.dialogMessages:type_name -> Messages
	3,  // 7: Chat.GetAllDialogues:input_type -> IdPage
	6,  // 8: Chat.GetOneDialogue:input_type -> IdIdPage
	5,  // 9: Chat.DeleteDialogue:input_type -> IdId
	4,  // 10: Chat.SendMessage:input_type -> SendEditMessage
	4,  // 11: Chat.EditMessage:input_type -> SendEditMessage
	5,  // 12: Chat.DeleteMessage:input_type -> IdId
	1,  // 13: Chat.Mailing:input_type -> MailingIn
	2,  // 14: Chat.Search:input_type -> SearchIn
	8,  // 15: Chat.GetAllDialogues:output_type -> DialogueCards
	12, // 16: Chat.GetOneDialogue:output_type -> Dialogue
	13, // 17: Chat.DeleteDialogue:output_type -> Answer
	13, // 18: Chat.SendMessage:output_type -> Answer
	13, // 19: Chat.EditMessage:output_type -> Answer
	13, // 20: Chat.DeleteMessage:output_type -> Answer
	13, // 21: Chat.Mailing:output_type -> Answer
	8,  // 22: Chat.Search:output_type -> DialogueCards
	15, // [15:23] is the sub-list for method output_type
	7,  // [7:15] is the sub-list for method input_type
	7,  // [7:7] is the sub-list for extension type_name
	7,  // [7:7] is the sub-list for extension extendee
	0,  // [0:7] is the sub-list for field type_name
}

func init() { file_chat_proto_init() }
func file_chat_proto_init() {
	if File_chat_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_chat_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*Ids); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_chat_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*MailingIn); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_chat_proto_msgTypes[2].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*SearchIn); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_chat_proto_msgTypes[3].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*IdPage); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_chat_proto_msgTypes[4].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*SendEditMessage); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_chat_proto_msgTypes[5].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*IdId); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_chat_proto_msgTypes[6].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*IdIdPage); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_chat_proto_msgTypes[7].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*DialogueCard); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_chat_proto_msgTypes[8].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*DialogueCards); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_chat_proto_msgTypes[9].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*UserOnEvent); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_chat_proto_msgTypes[10].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*Message); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_chat_proto_msgTypes[11].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*Messages); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_chat_proto_msgTypes[12].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*Dialogue); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_chat_proto_msgTypes[13].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*Answer); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
	}
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: file_chat_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   14,
			NumExtensions: 0,
			NumServices:   1,
		},
		GoTypes:           file_chat_proto_goTypes,
		DependencyIndexes: file_chat_proto_depIdxs,
		MessageInfos:      file_chat_proto_msgTypes,
	}.Build()
	File_chat_proto = out.File
	file_chat_proto_rawDesc = nil
	file_chat_proto_goTypes = nil
	file_chat_proto_depIdxs = nil
}
