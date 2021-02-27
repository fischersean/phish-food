package reddit

import (
	//"encoding/json"
	//"io/ioutil"
	"testing"
)

func TestParsePostResponse(t *testing.T) {

	res := `{
	"data": {
		"children": [
			{
				"data": {
					"subreddit": "sub",
					"selftext": "",
					"title": "This is the title",
					"downs": 0,
					"upvote_ratio": 0.99,
					"ups": 12345,
					"score": 54321,
					"author": "Author Name",
					"num_comments": 98766,
					"permalink": "This is the permalink",
					"created_utc": 1592410647
				}
			}
		]
	}
}`

	p, err := parsePostResponse([]byte(res))
	if err != nil {
		t.Fatalf(err.Error())
	}

	if p[0].Subreddit != "sub" {
		t.Errorf("Incorrect sub detected: %s != %s", p[0].Subreddit, "sub")
	}

	if p[0].Score != 54321 {
		t.Errorf("Incorrect score detected: %d != %d", p[0].Score, 54321)
	}

	_, err = parsePostResponse([]byte(`this is bad json`))
	if err == nil {
		t.Errorf("Json unmarshal error not thrown with bad json")
	}
}
