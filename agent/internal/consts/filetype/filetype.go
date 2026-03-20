package filetype

type FileType int

const (
    INVALID FileType = iota
    FILE
    SYMLINK
)

func (ft FileType) ToString() string {
    return [...]string{"INVALID", "FILE", "SYMLINK"}[ft]
}


func FromString(s string) FileType {
    switch s {
    case "FILE":
        return FILE
    case "SYMLINK":
        return SYMLINK
    default:
        return INVALID
    }
}
