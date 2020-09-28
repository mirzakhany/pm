package config

type Int interface {
	Int() int
	Int64() int64
}

type String interface {
	String() string
}

type Float interface {
	Float32() float32
	Float64() float64
}

type Bool interface {
	Bool() bool
}

type intHolder struct {
	value *int64
}

type stringHolder struct {
	value *string
}

type floatHolder struct {
	value *float64
}

type boolHolder struct {
	value *bool
}

func (sh stringHolder) String() string {
	return *sh.value
}

func (ih intHolder) Int() int {
	return int(*ih.value)
}

func (ih intHolder) Int64() int64 {
	return *ih.value
}

func (fh floatHolder) Float32() float32 {
	return float32(*fh.value)
}

func (fh floatHolder) Float64() float64 {
	return *fh.value
}

func (bh boolHolder) Bool() bool {
	return *bh.value
}
