package entity

type Comment struct {
	CommentID  uint   `json:"comment_id"`
	UserID     uint   `json:"user_id"`
	NickName   string `json:"nickname"`
	PostID     uint   `json:"post_id"`
	Data       string `json:"data"`
	Likes      uint   `json:"likes"`
	Dislikes   uint   `json:"dislikes"`
	VoteStatus uint   `json:"vote_status"`
}

type CommentVote struct {
	UserID    uint `json:"user_id"`
	CommentID uint `json:"comment_id"`
	Vote      int  `json:"vote"`
}
