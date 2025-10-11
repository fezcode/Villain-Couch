package str

type Str string

func (s Str) Empty() bool {
	return len(s) == 0
}

func (s Str) String() string {
	return string(s)
}

func (s Str) StrPtr() *string {
	x := string(s)
	return &x
}

func (s Str) Len() int {
	return len(s)
}
