package model

// Post 메시지 구조체
type Post struct {
	Channel string
	Message string
	Param   interface{}
}

// NewPost 메시지 생성
func NewPost(channel string, message string, param interface{}) *Post {
	return &Post{Channel: channel, Message: message, Param: param}
}
