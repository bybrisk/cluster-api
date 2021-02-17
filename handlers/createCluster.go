
package handlers

import (
	"net/http"
	"fmt"
	"github.com/bybrisk/cluster-api/data"
)

// swagger:route POST /cluster/create cluster createCluster
// Create clusters of all the pending delivery of business account
//
// responses:
//	200: clusterCreateResponse
//  422: errorValidation
//  501: errorResponse

func (p *Cluster) CreateCluster (w http.ResponseWriter, r *http.Request){
	p.l.Println("Handle POST request -> cluster-api Module")
	cluster := &data.CreateClusterRequest{}

	err:=cluster.FromJSONToCreateClusterStruct(r.Body)
	if err!=nil {
		http.Error(w,"Data failed to unmarshel", http.StatusBadRequest)
	}

	//validate the data
	err = cluster.ValidateCreateClusters()
	if err!=nil {
		p.l.Println("Validation error in POST request -> cluster-api Module \n",err)
		http.Error(w,fmt.Sprintf("Error in data validation : %s",err), http.StatusBadRequest)
		return
	} 

	//add request to cluster queue
	res := data.AddToClusterQueue(cluster)

	//writing to the io.Writer
	err = res.CreateClusterToJSON(w)
	if err!=nil {
		http.Error(w,"Data with ID failed to marshel",http.StatusInternalServerError)		
	}
}