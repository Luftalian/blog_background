package model

// isAllowedContentType は指定されたMIMEタイプが許可されているかを確認します
func IsAllowedContentType(contentType string, allowedTypes []string) bool {
	for _, t := range allowedTypes {
		if contentType == t {
			return true
		}
	}
	return false
}

// mimeExtension はMIMEタイプからファイル拡張子を返します
func MimeExtension(contentType string) string {
	switch contentType {
	case "image/jpeg":
		return ".jpg"
	case "image/png":
		return ".png"
	case "image/gif":
		return ".gif"
	default:
		return ""
	}
}
