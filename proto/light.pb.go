// Code generated by protoc-gen-go.
// source: light.proto
// DO NOT EDIT!

/*
Package light is a generated protocol buffer package.

It is generated from these files:
	light.proto

It has these top-level messages:
	ShareMessage
	PeerMessage
	FileMessage
*/
package light

import proto "code.google.com/p/goprotobuf/proto"
import math "math"

// Reference imports to suppress errors if they are not otherwise used.
var _ = proto.Marshal
var _ = math.Inf

type ShareAction int32

const (
	ShareAction_ENTERING ShareAction = 0
	ShareAction_LEAVING  ShareAction = 1
)

var ShareAction_name = map[int32]string{
	0: "ENTERING",
	1: "LEAVING",
}
var ShareAction_value = map[string]int32{
	"ENTERING": 0,
	"LEAVING":  1,
}

func (x ShareAction) Enum() *ShareAction {
	p := new(ShareAction)
	*p = x
	return p
}
func (x ShareAction) String() string {
	return proto.EnumName(ShareAction_name, int32(x))
}
func (x *ShareAction) UnmarshalJSON(data []byte) error {
	value, err := proto.UnmarshalJSONEnum(ShareAction_value, data, "ShareAction")
	if err != nil {
		return err
	}
	*x = ShareAction(value)
	return nil
}

type FileAction int32

const (
	FileAction_CREATED FileAction = 0
	FileAction_UPDATED FileAction = 1
	FileAction_REMOVED FileAction = 2
)

var FileAction_name = map[int32]string{
	0: "CREATED",
	1: "UPDATED",
	2: "REMOVED",
}
var FileAction_value = map[string]int32{
	"CREATED": 0,
	"UPDATED": 1,
	"REMOVED": 2,
}

func (x FileAction) Enum() *FileAction {
	p := new(FileAction)
	*p = x
	return p
}
func (x FileAction) String() string {
	return proto.EnumName(FileAction_name, int32(x))
}
func (x *FileAction) UnmarshalJSON(data []byte) error {
	value, err := proto.UnmarshalJSONEnum(FileAction_value, data, "FileAction")
	if err != nil {
		return err
	}
	*x = FileAction(value)
	return nil
}

type ShareMessage struct {
	ShareName        *string      `protobuf:"bytes,1,req,name=share_name" json:"share_name,omitempty"`
	Action           *ShareAction `protobuf:"varint,2,req,name=action,enum=light.ShareAction" json:"action,omitempty"`
	XXX_unrecognized []byte       `json:"-"`
}

func (m *ShareMessage) Reset()         { *m = ShareMessage{} }
func (m *ShareMessage) String() string { return proto.CompactTextString(m) }
func (*ShareMessage) ProtoMessage()    {}

func (m *ShareMessage) GetShareName() string {
	if m != nil && m.ShareName != nil {
		return *m.ShareName
	}
	return ""
}

func (m *ShareMessage) GetAction() ShareAction {
	if m != nil && m.Action != nil {
		return *m.Action
	}
	return ShareAction_ENTERING
}

type PeerMessage struct {
	PeerName         *string  `protobuf:"bytes,1,req,name=peer_name" json:"peer_name,omitempty"`
	Address          *string  `protobuf:"bytes,2,req,name=address" json:"address,omitempty"`
	Port             *string  `protobuf:"bytes,3,req,name=port" json:"port,omitempty"`
	Shares           []string `protobuf:"bytes,4,rep,name=shares" json:"shares,omitempty"`
	XXX_unrecognized []byte   `json:"-"`
}

func (m *PeerMessage) Reset()         { *m = PeerMessage{} }
func (m *PeerMessage) String() string { return proto.CompactTextString(m) }
func (*PeerMessage) ProtoMessage()    {}

func (m *PeerMessage) GetPeerName() string {
	if m != nil && m.PeerName != nil {
		return *m.PeerName
	}
	return ""
}

func (m *PeerMessage) GetAddress() string {
	if m != nil && m.Address != nil {
		return *m.Address
	}
	return ""
}

func (m *PeerMessage) GetPort() string {
	if m != nil && m.Port != nil {
		return *m.Port
	}
	return ""
}

func (m *PeerMessage) GetShares() []string {
	if m != nil {
		return m.Shares
	}
	return nil
}

type FileMessage struct {
	Filename         *string     `protobuf:"bytes,1,req,name=filename" json:"filename,omitempty"`
	Folder           *bool       `protobuf:"varint,2,req,name=folder" json:"folder,omitempty"`
	Action           *FileAction `protobuf:"varint,3,req,name=action,enum=light.FileAction" json:"action,omitempty"`
	Hash             []byte      `protobuf:"bytes,4,opt,name=hash" json:"hash,omitempty"`
	XXX_unrecognized []byte      `json:"-"`
}

func (m *FileMessage) Reset()         { *m = FileMessage{} }
func (m *FileMessage) String() string { return proto.CompactTextString(m) }
func (*FileMessage) ProtoMessage()    {}

func (m *FileMessage) GetFilename() string {
	if m != nil && m.Filename != nil {
		return *m.Filename
	}
	return ""
}

func (m *FileMessage) GetFolder() bool {
	if m != nil && m.Folder != nil {
		return *m.Folder
	}
	return false
}

func (m *FileMessage) GetAction() FileAction {
	if m != nil && m.Action != nil {
		return *m.Action
	}
	return FileAction_CREATED
}

func (m *FileMessage) GetHash() []byte {
	if m != nil {
		return m.Hash
	}
	return nil
}

func init() {
	proto.RegisterEnum("light.ShareAction", ShareAction_name, ShareAction_value)
	proto.RegisterEnum("light.FileAction", FileAction_name, FileAction_value)
}
