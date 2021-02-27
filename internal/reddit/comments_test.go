package reddit

import (
	"testing"
)

func TestFetchComments(t *testing.T) {

	res := `[
  {},
  {
    "data": {
      "children": [
        {
          "kind": "t1",
          "data": {
            "id": "this is the id",
            "score": 1234,
            "replies": {
              "kind": "t1",
              "data": {
                "children": [
                  {
                    "kind": "t1",
                    "data": {
                      "id": "this is the 2nd id",
                      "score": 1234,
                      "replies": "",
                      "noop": ""
                    }
                  }
                ]
              }
            }
          }
        }
      ]
    }
  }
]`

	c, err := parseCommentResponse([]byte(res))
	if err != nil {
		t.Fatal(err.Error())
	}

	if len(c) == 0 {
		t.Fatalf("Empty slice returned")
	}

	if c[0].Id != "this is the id" {
		t.Errorf("Incorrect id found: %s != %s", c[0].Id, "this is the id")
	}

	if c[0].Replies[0].Id != "this is the 2nd id" {
		t.Errorf("Incorrect id found: %s != %s", c[0].Replies[0].Id, "this is the 2nd id")
	}

	_, err = parseCommentResponse([]byte(`this is bad json`))
	if err == nil {
		t.Errorf("Json unmarshal error not thrown with bad json")
	}

}
