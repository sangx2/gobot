package model

// Post 메시지 구조체
type Post struct {
	Message string
	RootID  interface{}
}

// NewPost 메시지 생성
func NewPost(message string, rootID interface{}) *Post {
	return &Post{Message: message, RootID: rootID}
}
