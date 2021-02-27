package reddit

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"mime/multipart"
	"net/http"
	"net/url"
)

type AuthToken struct {
	AccessToken string `json:"access_token"`
	TokenType   string `json:"token_type"`
	DeviceId    string `json:"device_id"`
	ExpiresIn   int    `json:"expires_in"`
	Scope       string `json:"scope"`
}

const (
	baseUrl   = "oauth.reddit.com"
	userAgent = "script:Subreddit Sentiment:v0.1.0 (by /u/hapins)"
)

func GetAuthToken(appId, appSecret string) (token AuthToken, err error) {

	url := "https://www.reddit.com/api/v1/access_token?scope=read"
	method := "POST"

	payload := &bytes.Buffer{}
	writer := multipart.NewWriter(payload)
	_ = writer.WriteField("grant_type", "https://oauth.reddit.com/grants/installed_client")
	_ = writer.WriteField("device_id", "DO_NOT_TRACK_THIS_DEVICE")
	err = writer.Close()
	if err != nil {
		return token, err
	}

	client := &http.Client{}
	req, err := http.NewRequest(method, url, payload)

	if err != nil {
		return token, err
	}

	req.SetBasicAuth(appId, appSecret)
	req.Header.Set("User-agent", userAgent)
	req.Header.Set("Content-Type", writer.FormDataContentType())

	res, err := client.Do(req)
	if err != nil {
		return token, err
	}
	defer res.Body.Close()

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return token, err
	}

	err = json.Unmarshal(body, &token)
	if err != nil {
		return token, err
	}

	if token.AccessToken == "" {
		return token, fmt.Errorf("could not retrieve token")
	}

	return token, nil

}

func oauthGet(u *url.URL, authorizatoin AuthToken) (resp *http.Response, err error) {

	// TODO: handle when the token is expired
	//if Authorization.AccessToken == "" {
	//Authorization, err = getAuthToken()
	//if err != nil {
	//return nil, err
	//}
	//}

	client := &http.Client{}

	req, err := http.NewRequest("GET", u.String(), nil)
	if err != nil {
		return nil, err
	}
	req.Header.Add("Authorization", fmt.Sprintf("%s %s", authorizatoin.TokenType, authorizatoin.AccessToken))
	req.Header.Add("User-Agent", userAgent)

	resp, err = client.Do(req)

	// **** Commenting this out since reddit doesn't appear to be providing these headers anymore ****
	//rateLimitUsed, _ = strconv.Atoi(resp.Header.Get("X-Ratelimit-Used"))
	//rateLimitRemaining, _ = strconv.Atoi(resp.Header.Get("X-Ratelimit-Remaining"))
	//rateLimitReset, _ = strconv.Atoi(resp.Header.Get("X-Ratelimit-Reset"))

	//fmt.Printf("Used %d requests. %d remaining", rateLimitUsed, rateLimitRemaining)

	//if rateLimitRemaining == 0 {
	//time.Sleep(time.Duration(rateLimitReset) * time.Second)
	//}

	return resp, err

}

func redditGet(path string, headers map[string]string, authorization AuthToken) (resp *http.Response, err error) {

	u := &url.URL{
		Scheme: "https",
		Host:   baseUrl,
		Path:   path,
	}

	v := url.Values{}
	for k, h := range headers {
		v.Add(k, h)
	}
	u.RawQuery = v.Encode()

	return oauthGet(u, authorization)

}
