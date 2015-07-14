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

func (r *Response) JSON(f interface{}) error {
	buf, _ := ioutil.ReadAll(r.Body)
	defer r.Body.Close()

	return json.Unmarshal(buf, &f)
}
