package enum

type FileTypeEnum string

const (
	IMAGE FileTypeEnum = "image"
	VIDEO FileTypeEnum = "video"
	FILE  FileTypeEnum = "file"
)

func (e FileTypeEnum) ToString() string {
	switch e {
	case IMAGE:
		return "image"
	case FILE:
		return "file"
	case VIDEO:
		return "video"
	default:
		return ""
	}
}

func (e FileTypeEnum) IsValid() bool {
	switch e {
	case IMAGE, FILE, VIDEO:
		return true
	}

	return false
}
