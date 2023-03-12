package example

type MyStruct struct {
	Age     int     `visibility:"readonly"`
	Ptr     *string `visibility:"readonly"` // want `cannot define visibility of pointer types`
	Invalid string  `visibility:"invalid"`  // want `invalid visibility tag: "invalid"`
}

func (s *MyStruct) SetAge(age int) {
	s.Age = age
}

func _() {
	m := MyStruct{}

	m.Age = 1  // want `misuse of readonly field: cannot assign`
	m.Age += 1 // want `misuse of readonly field: cannot assign`
	m.Age++    // want `misuse of readonly field: cannot increment`
	_ = &m.Age // want `misuse of readonly field: cannot take address`

	x := struct {
		m MyStruct
	}{}

	x.m.Age = 1 // want `misuse of readonly field: cannot assign`
}
