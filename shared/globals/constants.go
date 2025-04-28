package globals

const (
	// all units in words
	DataSectionLength = 16
	StackLength       = 64
	FlagLength        = 1
	TextSectionLength = 175
)

const (
	// Byte masks to check against Mem(79) / Process Flag
	HaltFlag  = 0x80 // 0b 1000 0000
	SleepFlag = 0x40 // 0b 0100 0000
)
