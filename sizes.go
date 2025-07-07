package msgpack

const (
	mask4 = 0x0f
	mask5 = 0x1f
	mask6 = 0x3f

	mask8 = 0xff
	size8 = mask8 + 1
	mask7 = mask8 >> 1

	mask16 = 0xffff
	size16 = mask16 + 1
	mask15 = mask16 >> 1

	mask32 = 0xffffffff
	size32 = mask32 + 1
	mask31 = mask32 >> 1

	mask34 = 0x03ffffffff
	size34 = mask34 + 1
)

const (
	intFixMax = mask7
	intFixMin = -32

	int8Max = mask7
	int8Min = -int8Max - 1

	int16Max = mask15
	int16Min = -int16Max - 1

	int32Max = mask31
	int32Min = -int32Max - 1
)
