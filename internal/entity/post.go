package entity

// Post is returned to the client with voting and comments metadata included.
type Post struct {
	PostID        uint      `json:"post_id"`
	UserID        uint      `json:"user_id"`
	NickName      string    `json:"nickname"`
	Title         string    `json:"title"`
	Data          string    `json:"data"`
	Likes         uint      `json:"likes"`
	Dislikes      uint      `json:"dislikes"`
	VoteStatus    uint      `json:"vote_status"` // 0: no vote, 1: like, 2: dislike
	Comments      []Comment `json:"comments"`
	CommentsCount uint      `json:"comments_count"`
	Categorys     []string  `json:"categories"`
}

// Category describes a community topic a post can be tagged with.
type Category struct {
	CategoryID  uint   `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
}

// CategoryAndPost joins a post to one of its categories.
type CategoryAndPost struct {
	CategoryID uint
	PostID     uint
}

// PostVote captures an upvote/downvote action on a post.
type PostVote struct {
	UserID uint `json:"user_id"`
	PostID uint `json:"post_id"`
	Vote   int  `json:"vote"`
}
