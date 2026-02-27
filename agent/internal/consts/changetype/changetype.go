package changetype

type ChangeType int

const (
    INVALID ChangeType = iota
	ACCOUNT
	PROJECT
	BRANCH
	FILE
	TAG
)

func (ct ChangeType) ToString() string {
	return [...]string{"INVALID", "ACCOUNT", "PROJECT", "BRANCH", "FILE", "TAG"}[ct]
}

func FromString(s string) ChangeType {
	switch s {
	case "ACCOUNT":
		return ACCOUNT
	case "PROJECT":
		return PROJECT
	case "BRANCH":
		return BRANCH
	case "FILE":
		return FILE
	case "TAG":
		return TAG
	default:
		return INVALID
	}
}
