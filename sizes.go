package msgpack

const (
	mask4  = 0x0f
	mask5  = 0x1f
	mask6  = 0x3f
	mask7  = 0x7f
	mask8  = 0xff
	mask16 = 0xffff
	mask32 = 0xffffffff
	mask34 = 0x03ffffffff
)

const (
	intFixMin = -(mask6 >> 1) - 1
	intFixMax = mask8 >> 1
	int8Min   = -(mask8 >> 1) - 1
	int8Max   = mask8 >> 1
	int16Min  = -(mask16 >> 1) - 1
	int16Max  = mask16 >> 1
	int32Min  = -(mask32 >> 1) - 1
	int32Max  = mask32 >> 1
)
