package entity

type Post struct {
	PostID        uint      `json:"post_id"`
	UserID        uint      `json:"user_id"`
	UserName      string    `json:"username"`
	Title         string    `json:"title"`
	Data          string    `json:"data"`
	Likes         uint      `json:"likes"`
	Dislikes      uint      `json:"dislikes"`
	VoteStatus    uint      `json:"vote_status"` // 0: no vote, 1: like, 2: dislike
	Comments      []Comment `json:"comments"`
	CommentsCount uint      `json:"comments_count"`
	Categorys     []string  `json:"categories"`
}

type Category struct {
	CategoryID  uint   `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
}

type CategoryAndPost struct {
	CategoryID uint
	PostID     uint
}

type PostVote struct {
	UserID uint `json:"user_id"`
	PostID uint `json:"post_id"`
	Vote   int  `json:"vote"`
}
