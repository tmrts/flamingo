package rest

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

	//buf := []byte("{\"attributes\":{},\"cpuPlatform\":\"Intel Sandy Bridge\",\"description\":\"\",\"disks\":[{\"deviceName\":\"persistent-disk-0\",\"index\":0,\"mode\":\"READ_WRITE\",\"type\":\"PERSISTENT\"}],\"hostname\":\"centos.c.total-tooling-96110.internal\",\"id\":533952008838200936,\"image\":\"\",\"machineType\":\"projects/939642202004/machineTypes/n1-standard-1\",\"maintenanceEvent\":\"NONE\",\"networkInterfaces\":[{\"accessConfigs\":[{\"externalIp\":\"104.155.29.138\",\"type\":\"ONE_TO_ONE_NAT\"}],\"forwardedIps\":[],\"ip\":\"10.240.41.122\",\"network\":\"projects/939642202004/networks/default\"}],\"scheduling\":{\"automaticRestart\":\"TRUE\",\"onHostMaintenance\":\"MIGRATE\"},\"serviceAccounts\":{\"939642202004-compute@developer.gserviceaccount.com\":{\"aliases\":[\"default\"],\"email\":\"939642202004-compute@developer.gserviceaccount.com\",\"scopes\":[\"https://www.googleapis.com/auth/computeaccounts.readonly\",\"https://www.googleapis.com/auth/devstorage.read_only\",\"https://www.googleapis.com/auth/logging.write\"]},\"default\":{\"aliases\":[\"default\"],\"email\":\"939642202004-compute@developer.gserviceaccount.com\",\"scopes\":[\"https://www.googleapis.com/auth/computeaccounts.readonly\",\"https://www.googleapis.com/auth/devstorage.read_only\",\"https://www.googleapis.com/auth/logging.write\"]}},\"tags\":[],\"virtualClock\":{\"driftToken\":\"11239462191112056598\"},\"zone\":\"projects/939642202004/zones/europe-west1-b\"}")
	return json.Unmarshal(buf, &f)
}
