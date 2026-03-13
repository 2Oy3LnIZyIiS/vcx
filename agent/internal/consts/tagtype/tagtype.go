package tagtype

type TagType int

const (
	INVALID TagType = iota
	SYSTEM_PROJECT
	SYSTEM_BRANCH
	SYSTEM_FILE
	USER
	USER_PROJECT
	USER_BRANCH
	USER_FILE
)

func (tt TagType) ToString() string {
	return [...]string{"INVALID", "SYSTEM_PROJECT", "SYSTEM_BRANCH", "SYSTEM_FILE", "USER", "USER_PROJECT", "USER_BRANCH", "USER_FILE"}[tt]
}


func FromString(s string) TagType {
	switch s {
	case "SYSTEM_PROJECT":
		return SYSTEM_PROJECT
	case "SYSTEM_BRANCH":
		return SYSTEM_BRANCH
	case "SYSTEM_FILE":
		return SYSTEM_FILE
	case "USER":
		return USER
	case "USER_PROJECT":
		return USER_PROJECT
	case "USER_BRANCH":
		return USER_BRANCH
	case "USER_FILE":
		return USER_FILE
	default:
		return INVALID
	}
}
