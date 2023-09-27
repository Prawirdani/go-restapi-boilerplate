package httputil

import (
	"bytes"
	"encoding/json"
	"net/http"
)

func SendJson(w http.ResponseWriter, status_code int, data interface{}) {
	buf := &bytes.Buffer{}
	enc := json.NewEncoder(buf)
	enc.SetEscapeHTML(true)

	response := NewResponse(status_code, data)
	
	if err := enc.Encode(response); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status_code)
	w.Write(buf.Bytes())
}

func BindJson(r *http.Request, request interface{}) error {
	defer r.Body.Close()
	dec := json.NewDecoder(r.Body)
	if err := dec.Decode(request); err != nil {
		return err
	}
	return nil
}

func SendError(w http.ResponseWriter, Err error) {
	buf := &bytes.Buffer{}
	enc := json.NewEncoder(buf)
	enc.SetEscapeHTML(true)

	response := NewErrorResponse(Err)

	if err := enc.Encode(response); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(response.Code)
	w.Write(buf.Bytes())
}
