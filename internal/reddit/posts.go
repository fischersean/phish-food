package reddit

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"strconv"
)

type postResponse struct {
	Data struct {
		Children []struct {
			Data Post `json:"data"`
		} `json:"children"`
	} `json:"data"`
}

func parsePostResponse(b []byte) (posts []Post, err error) {

	rawResponse := postResponse{}
	err = json.Unmarshal(b, &rawResponse)
	if err != nil {
		return posts, err
	}

	for _, child := range rawResponse.Data.Children {
		// Turning off auto fetch of comments
		//tmpPost := child.Data
		//tmpPost.Comments, err = FetchPostComments(tmpPost)
		//if err != nil {
		//return posts, err
		//}
		posts = append(posts, child.Data)
	}

	return posts, err
}

func GetPosts(subreddit string, limit int, sort string, authorization AuthToken) (posts []Post, err error) {

	validSorts := map[string]int{
		"top":    0,
		"hot":    0,
		"new":    0,
		"rising": 0,
	}

	if _, ok := validSorts[sort]; !ok {
		return posts, fmt.Errorf(fmt.Sprintf("Sort type not supported: %s", sort))
	}

	res, err := redditGet(fmt.Sprintf("r/%s/%s.json", subreddit, sort), map[string]string{
		"limit": strconv.FormatInt(int64(limit), 10),
	}, authorization)
	if err != nil {
		return posts, err
	}
	defer res.Body.Close()

	b, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return posts, err
	}

	return parsePostResponse(b)
}

func GetHot(subreddit string, limit int, authorization AuthToken) ([]Post, error) {
	return GetPosts(subreddit, limit, "hot", authorization)
}

func GetTop(subreddit string, limit int, authorization AuthToken) ([]Post, error) {
	return GetPosts(subreddit, limit, "top", authorization)
}

func GetNew(subreddit string, limit int, authorization AuthToken) ([]Post, error) {
	return GetPosts(subreddit, limit, "new", authorization)
}

func GetRising(subreddit string, limit int, authorization AuthToken) ([]Post, error) {
	return GetPosts(subreddit, limit, "rising", authorization)
}
