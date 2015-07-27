package request

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
)

/*
 *type Response interface {
 *    StatusCode() int
 *
 *    Headers() http.Header
 *    Content() chan []byte
 *
 *    Request() *Request
 *}
 */

type Response struct {
	*http.Response
}

func (r *Response) Text() ([]byte, error) {
	buf, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return nil, err
	}
	defer r.Body.Close()

	return buf, nil
}

func (r *Response) JSON(f interface{}) error {
	buf, err := r.Text()
	if err != nil {
		return err
	}

	return json.Unmarshal(buf, &f)
}
