package protocol

const (
	C2GS_HELLO          = 0x01
	C2GS_IDENTIFY       = 0x02
	C2GS_LOGIN          = 0x03
	C2GS_ROLE           = 0x04
	C2GS_LOADROLEINFO   = 0x05
	C2GS_LOGIN_FINISHED = 0x06
)

const (
	GS2C_HELLO          = 0x01
	GS2C_IDENTIFY       = 0x02
	GS2C_LOGIN          = 0x03
	GS2C_ROLE           = 0x04
	GS2C_LOADROLEINFO   = 0x05
	GS2C_LOGIN_FINISHED = 0x06
)

var TEST_ECHO int = 0xff
