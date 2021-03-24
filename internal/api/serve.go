package api

import (
	"encoding/json"
	"net/http"
)

func setHeaders(w http.ResponseWriter) {

	w.Header().Set("Access-Control-Allow-Origin", "*")

}

// HttpServeJSON sets common headers for serving JSON data
func HttpServeJSON(w http.ResponseWriter, b []byte) (int, error) {

	setHeaders(w)
	w.Header().Set("Content-Type", "application/json")

	wlen, err := w.Write(b)
	if err != nil {
		return wlen, err
	}

	return wlen, nil

}

// ServeMarshallableData serves to w the interface. Will throw an error if data cannot be marhalled
func HttpServeMarahallableData(w http.ResponseWriter, v interface{}) (int, error) {

	b, err := json.Marshal(v)
	if err != nil {
		return 0, err
	}

	return HttpServeJSON(w, b)
}
