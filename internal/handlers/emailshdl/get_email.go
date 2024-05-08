package emailshdl

type getEmail struct {
	ID string `param:"id" validate:"empty=false"`
}
