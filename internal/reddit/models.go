package reddit

import (
	"time"
)

type Post struct {
	Subreddit    string    `json:"subreddit"`
	Title        string    `json:"title"`
	Body         string    `json:"selftext"`
	Ups          int       `json:"ups"`
	Downs        int       `json:"downs"`
	Score        int       `json:"score"`
	UpvoteRatio  float64   `json:"upvote_ratio"`
	Created      float64   `json:"created_utc"`
	Author       string    `json:"author"`
	Permalink    string    `json:"permalink"`
	NumComments  int       `json:"num_comments"`
	Comments     []Comment `json:"comments"`
	DownloadTime time.Time `json:"-"`
}

type Comment struct {
	Id      string    `json:"id"`
	Body    string    `json:"body"`
	Ups     int       `json:"ups"`
	Downs   int       `json:"downs"`
	Score   int       `json:"score"`
	Created time.Time `json:"created_utc"`
	Depth   int       `json:"depth"`
	Author  string    `json:"author"`

	// Parent and children needs to be filled in after struct has been unmarshalled
	//Parent   *Comment   `json:"-"`
	Replies []Comment `json:"replies"`
}
