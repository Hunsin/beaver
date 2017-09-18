package beaver

import (
	"encoding/json"
	"net/http"
	"os"
)

// WriteFile writes bytes in given pathname.
// It overwrites the file if it already exists.
// The caller should close the file if necessary.
func WriteFile(path string, body []byte) error {
	f, err := os.Create(path)
	if err != nil {
		return err
	}
	defer f.Close()

	if _, err = f.Write(body); err != nil {
		return err
	}

	return f.Sync()
}

// OpenJSON reads JSON file from given path and decode.
func OpenJSON(path string, body interface{}) error {
	f, err := os.Open(path)
	if err != nil {
		return err
	}
	defer f.Close()

	return json.NewDecoder(f).Decode(body)
}

// WriteJSON saves an object to JSON file by given path.
// It appends ".json" to the path if it doesn't end with the string.
func WriteJSON(path string, body interface{}) error {
	if path[len(path)-5:] != ".json" {
		path += ".json"
	}

	b, err := json.MarshalIndent(body, "", "\t")
	if err != nil {
		return err
	}

	return WriteFile(path, b)
}

// ServeJSON sends response to client with given status code and object.
// The Content-Type header is set to "application/json; charset=utf-8".
// Additional headers must be set before calling ServeJSON.
func ServeJSON(w http.ResponseWriter, status int, body interface{}) error {
	js, err := json.Marshal(body)
	if err != nil {
		return err
	}

	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(status)
	_, err = w.Write(js)
	return err
}
