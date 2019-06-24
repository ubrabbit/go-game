package protocol

import (
	. "server/leaf/network"
)

type C2GSHello struct {
	Seed int
}

type GS2CHello struct {
	Seed int
}

type C2GSIdentity struct {
	Identity string
}

type GS2CIdentity struct {
}

type C2GSLogin struct {
	User     string
	Password string
}

type GS2CLogin struct {
	Type int
	Pid  int
}

type C2GSRoleID struct {
}

type GS2CRoleID struct {
}

type C2GSLoadRole struct {
	Pid int
}

type GS2CLoadRole struct {
	NLen int
	Name string
}

type C2GSLoginFinished struct {
}

type GS2CLoginFinished struct {
}

func (p *C2GSHello) Protocol() uint8 {
	return C2GS_HELLO
}

func (p *C2GSHello) PacketData() (uint8, []byte) {
	data := make([]byte, 0)
	data = PacketInt(data, p.Seed, 4)
	return p.Protocol(), data
}

func (p *C2GSHello) UnpackData(from []byte) {
	seed, _ := UnpackInt(from, 4)
	p.Seed = int(seed)
}

func (p *GS2CHello) Protocol() uint8 {
	return GS2C_HELLO
}

func (p *GS2CHello) PacketData() (uint8, []byte) {
	data := make([]byte, 0)
	data = PacketInt(data, p.Seed, 4)
	return p.Protocol(), data
}

func (p *GS2CHello) UnpackData(from []byte) {
	seed, _ := UnpackInt(from, 4)
	p.Seed = int(seed)
}

func (p *C2GSIdentity) Protocol() uint8 {
	return C2GS_IDENTIFY
}

func (p *C2GSIdentity) PacketData() (uint8, []byte) {
	data := make([]byte, 0)
	data = PacketString(data, p.Identity, 24)
	return p.Protocol(), data
}

func (p *C2GSIdentity) UnpackData(from []byte) {
	ident, _ := UnpackString(from, 24)
	p.Identity = ident
}

func (p *GS2CIdentity) Protocol() uint8 {
	return GS2C_IDENTIFY
}

func (p *GS2CIdentity) PacketData() (uint8, []byte) {
	return p.Protocol(), nil
}

func (p *GS2CIdentity) UnpackData(from []byte) {
}

func (p *C2GSLogin) Protocol() uint8 {
	return C2GS_LOGIN
}

func (p *C2GSLogin) PacketData() (uint8, []byte) {
	data := make([]byte, 0)
	nLen, pLen := len(p.User), len(p.Password)
	data = PacketInt(data, nLen, 1)
	data = PacketString(data, p.User, nLen)
	data = PacketInt(data, pLen, 1)
	data = PacketString(data, p.Password, pLen)
	return p.Protocol(), data
}

func (p *C2GSLogin) UnpackData(from []byte) {
	nLen, from := UnpackInt(from, 1)
	user, from := UnpackString(from, nLen)
	pLen, from := UnpackInt(from, 1)
	password, from := UnpackString(from, pLen)
	p.User = user
	p.Password = password
}

func (p *GS2CLogin) Protocol() uint8 {
	return GS2C_LOGIN
}

func (p *GS2CLogin) PacketData() (uint8, []byte) {
	data := make([]byte, 0)
	data = PacketInt(data, p.Type, 1)
	data = PacketInt(data, p.Pid, 4)
	return p.Protocol(), data
}

func (p *GS2CLogin) UnpackData(from []byte) {
	t, from := UnpackInt(from, 1)
	pid, _ := UnpackInt(from, 4)
	p.Type = int(t)
	p.Pid = pid
}

func (p *C2GSRoleID) Protocol() uint8 {
	return C2GS_ROLE
}

func (p *C2GSRoleID) PacketData() (uint8, []byte) {
	data := make([]byte, 0)
	return p.Protocol(), data
}

func (p *C2GSRoleID) UnpackData(from []byte) {
}

func (p *GS2CRoleID) Protocol() uint8 {
	return GS2C_ROLE
}

func (p *GS2CRoleID) PacketData() (uint8, []byte) {
	data := make([]byte, 0)
	return p.Protocol(), data
}

func (p *GS2CRoleID) UnpackData(from []byte) {
}

func (p *C2GSLoadRole) Protocol() uint8 {
	return C2GS_LOADROLEINFO
}

func (p *C2GSLoadRole) PacketData() (uint8, []byte) {
	data := make([]byte, 0)
	data = PacketInt(data, p.Pid, 4)
	return p.Protocol(), data
}

func (p *C2GSLoadRole) UnpackData(from []byte) {
	pid, _ := UnpackInt(from, 4)
	p.Pid = pid
}

func (p *GS2CLoadRole) Protocol() uint8 {
	return GS2C_LOADROLEINFO
}

func (p *GS2CLoadRole) PacketData() (uint8, []byte) {
	data := make([]byte, 0)
	data = PacketInt(data, p.NLen, 1)
	data = PacketString(data, p.Name, p.NLen)
	return p.Protocol(), data
}

func (p *GS2CLoadRole) UnpackData(from []byte) {
	nLen, from := UnpackInt(from, 1)
	name, from := UnpackString(from, nLen)
	p.NLen = nLen
	p.Name = name
}

func (p *C2GSLoginFinished) Protocol() uint8 {
	return C2GS_LOGIN_FINISHED
}

func (p *C2GSLoginFinished) PacketData() (uint8, []byte) {
	data := make([]byte, 0)
	return p.Protocol(), data
}

func (p *C2GSLoginFinished) UnpackData(from []byte) {
}

func (p *GS2CLoginFinished) Protocol() uint8 {
	return GS2C_LOGIN_FINISHED
}

func (p *GS2CLoginFinished) PacketData() (uint8, []byte) {
	data := make([]byte, 0)
	return p.Protocol(), data
}

func (p *GS2CLoginFinished) UnpackData(from []byte) {
}
