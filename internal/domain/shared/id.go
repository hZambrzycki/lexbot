package shared

type ID string

func NewID(value string) ID {
	return ID(value)
}

func (id ID) String() string {
	return string(id)
}