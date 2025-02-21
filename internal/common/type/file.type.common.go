package types

import "onx-outgoing-go/internal/common/enum"

type BufferedFile struct {
	MediaType    string `json:"mediaType" validate:"required"`
	OriginalName string `json:"originalName" validate:"required"`
	Encoding     string `json:"encoding" validate:"required"`
	MimeType     string `json:"mimetype" validate:"required"`
	Size         int    `json:"size" validate:"required"`
	Buffer       []byte `json:"buffer" validate:"required"`
}

type BufferedFiles map[string][]BufferedFile

type FileType enum.FileTypeEnum

func (e FileType) IsValidImage(file *BufferedFile) bool {
	if enum.FileTypeEnum(e) == enum.IMAGE {
		switch file.MimeType {
		case "image/jpeg", "image/png", "image/gif", "image/bmp", "image/webp", "image/tiff", "image/svg+xml", "image/x-icon", "image/heic", "image/heif":
			return true
		}
	}

	return false
}

func (e FileType) IsValidVideo(file *BufferedFile) bool {
	if enum.FileTypeEnum(e) == enum.VIDEO {
		switch file.MimeType {
		case "video/mp4", "video/webm", "video/ogg", "video/avi", "video/mkv", "video/quicktime", "video/x-flv", "video/x-msvideo":
			return true
		}
	}

	return false
}
