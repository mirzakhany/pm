package config

type stringHolderMock struct {
	v string
}

type intHolderMock struct {
	v int
}

func (i intHolderMock) Int() int {
	return i.v
}

func (i intHolderMock) Int64() int64 {
	return int64(i.v)
}

func (h stringHolderMock) String() string {
	return h.v
}

func RegisterStringMock(key, defValue string) String {
	return stringHolderMock{v: defValue}
}

func RegisterIntMock(key string, defValue int) Int {
	return intHolderMock{v: defValue}
}
