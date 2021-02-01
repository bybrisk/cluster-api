
package handlers

import (
	"net/http"
	"github.com/gorilla/mux"
	"github.com/bybrisk/cluster-api/data"
)

// swagger:route GET /cluster/all/{id} cluster getAllCluster
// Get clusterIds and details of all the deliveries assigned to that Id using businessID.
//
// responses:
//	200: getAllClusterDetails
//  422: errorValidation
//  501: errorResponse

func (p *Cluster) GetAllCluster(w http.ResponseWriter, r *http.Request) {
	p.l.Println("Handle GET request -> cluster-api Module")
	
	vars := mux.Vars(r)
	id := vars["id"]

	lp := data.GetAllClusterByID(id)

	err := lp.AllClusterResponseToJSON(w)
	if err!=nil {
		http.Error(w,"Data failed to marshel",http.StatusInternalServerError)		
	}
}