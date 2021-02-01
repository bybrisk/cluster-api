
package handlers

import (
	"net/http"
	"github.com/gorilla/mux"
	"github.com/bybrisk/cluster-api/data"
)

// swagger:route GET /cluster/one/{clusterID} cluster getSingleCluster
// Get details of all the deliveries assigned to a cluster.
//
// responses:
//	200: getSingleClusterDetails
//  422: errorValidation
//  501: errorResponse

func (p *Cluster) GetSingleCluster(w http.ResponseWriter, r *http.Request) {
	p.l.Println("Handle GET request -> cluster-api Module")
	
	vars := mux.Vars(r)
	id := vars["clusterID"]

	lp := data.GetSingleClusterByID(id)

	err := lp.SingleClusterResponseToJSON(w)
	if err!=nil {
		http.Error(w,"Data failed to marshel",http.StatusInternalServerError)		
	}
}