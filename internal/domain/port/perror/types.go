package perror

type Type int

const (
	TypeUnspecified Type = iota
	TypeNotFound
	TypeAlreadyExists
	TypeInvalidParam
	TypeNoRowsUpdated
)
