package rest

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/tmrts/flamingo/pkg/datasrc/metadata"
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

type Response *http.Response

type RESTResponse struct {
}

func JSON(r Response) (*metadata.GCEv1, error) {
	var f metadata.GCEv1

	//buf, _ := ioutil.ReadAll(r.Body)

	buf := []byte("{\"attributes\":{},\"cpuPlatform\":\"Intel Sandy Bridge\",\"description\":\"\",\"disks\":[{\"deviceName\":\"persistent-disk-0\",\"index\":0,\"mode\":\"READ_WRITE\",\"type\":\"PERSISTENT\"}],\"hostname\":\"centos.c.total-tooling-96110.internal\",\"id\":533952008838200936,\"image\":\"\",\"machineType\":\"projects/939642202004/machineTypes/n1-standard-1\",\"maintenanceEvent\":\"NONE\",\"networkInterfaces\":[{\"accessConfigs\":[{\"externalIp\":\"104.155.29.138\",\"type\":\"ONE_TO_ONE_NAT\"}],\"forwardedIps\":[],\"ip\":\"10.240.41.122\",\"network\":\"projects/939642202004/networks/default\"}],\"scheduling\":{\"automaticRestart\":\"TRUE\",\"onHostMaintenance\":\"MIGRATE\"},\"serviceAccounts\":{\"939642202004-compute@developer.gserviceaccount.com\":{\"aliases\":[\"default\"],\"email\":\"939642202004-compute@developer.gserviceaccount.com\",\"scopes\":[\"https://www.googleapis.com/auth/computeaccounts.readonly\",\"https://www.googleapis.com/auth/devstorage.read_only\",\"https://www.googleapis.com/auth/logging.write\"]},\"default\":{\"aliases\":[\"default\"],\"email\":\"939642202004-compute@developer.gserviceaccount.com\",\"scopes\":[\"https://www.googleapis.com/auth/computeaccounts.readonly\",\"https://www.googleapis.com/auth/devstorage.read_only\",\"https://www.googleapis.com/auth/logging.write\"]}},\"tags\":[],\"virtualClock\":{\"driftToken\":\"11239462191112056598\"},\"zone\":\"projects/939642202004/zones/europe-west1-b\"}")
	err := json.Unmarshal(buf, &f)
	if err != nil {
		return nil, err
	}

	fmt.Printf("%v", f.Digest())

	//m := f.(map[string]interface{})
	//for k, v := range m {
	//switch v.(type) {
	//default:
	//fmt.Printf("%T %v: %v\n", v, k, v)
	//}
	//}
	return &f, nil
}

type Request struct {
	Method, URL string

	Headers http.Header
}

func (r *Request) Normalize() *http.Request {
	req, err := http.NewRequest(r.Method, r.URL, nil)
	if err != nil {
		panic(err)
	}

	req.Header = r.Headers

	return req
}

type Parameter func(*Request)

//request.Header("Metadata-Flavour", "Google")
func Header(key, value string) Parameter {
	return func(r *Request) {
		r.Headers.Add(key, value)
	}
}

type Client interface {
	Perform(string, ...Parameter) (Response, error)
}

type ClientImplementation struct {
	*http.Client
}

var DefaultClient = ClientImplementation{http.DefaultClient}

func (c *ClientImplementation) Perform(r *Request) (Response, error) {
	req := r.Normalize()

	resp, err := c.Do(req)
	if err != nil {
		return nil, err
	}

	return Response(resp), nil
}

func Get(url string, params ...Parameter) (Response, error) {
	r := &Request{
		URL:     url,
		Method:  "GET",
		Headers: http.Header{},
	}

	for _, parametrize := range params {
		parametrize(r)
	}

	return DefaultClient.Perform(r)
}

//func Put(string, ...Parameter) (*Response, error)     {}
//func Post(string, ...Parameter) (*Response, error)    {}
//func Head(string, ...Parameter) (*Response, error)    {}
//func Delete(string, ...Parameter) (*Response, error)  {}
//func Options(string, ...Parameter) (*Response, error) {}
