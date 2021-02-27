package reddit

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"time"
)

type commentData struct {
	Kind string `json:"kind"`
	Data struct {
		Id      string           `json:"id"`
		Body    string           `json:"body"`
		Ups     int              `json:"ups"`
		Downs   int              `json:"downs"`
		Score   int              `json:"score"`
		Created float64          `json:"created_utc"`
		Depth   int              `json:"depth"`
		Author  string           `json:"author"`
		Replies *commentResponse `json:"replies"`
	} `json:"data"`
}

type commentResponse struct {
	Data struct {
		Comments []commentData `json:"children"`
	} `json:"data"`
}

func parseCommentResponse(b []byte) (comments []Comment, err error) {

	rawResponse := []commentResponse{}

	// Fix reddit's busted replies return type
	b = bytes.ReplaceAll(b, []byte("\"replies\": \"\","), []byte{})

	err = json.Unmarshal(b, &rawResponse)
	if err != nil {
		return comments, err
	}

	for _, v := range rawResponse[1].Data.Comments {
		comments = append(comments, commentTreeTraverse(v))
	}

	return comments, err
}

func commentTreeTraverse(comment commentData) Comment {

	data := comment.Data
	c := Comment{
		Id:      data.Id,
		Body:    data.Body,
		Ups:     data.Ups,
		Downs:   data.Downs,
		Score:   data.Score,
		Created: time.Unix(int64(data.Created), 0),
		Depth:   data.Depth,
		Author:  data.Author,
		Replies: []Comment{},
	}

	if data.Replies == nil {
		return c
	}
	replies := data.Replies.Data.Comments
	for _, cmt := range replies {
		c.Replies = append(c.Replies, commentTreeTraverse(cmt))
	}

	return c
}

func FetchPostComments(p Post, authorization AuthToken) (c []Comment, err error) {

	res, err := redditGet(fmt.Sprintf("%s.json", p.Permalink), nil, authorization)
	if err != nil {
		return c, err
	}
	defer res.Body.Close()

	b, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return c, err
	}

	return parseCommentResponse(b)
}
