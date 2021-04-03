
package handlers

import (
	"net/http"
	"github.com/gorilla/mux"
	"github.com/bybrisk/cluster-api/data"
)

// swagger:route GET /cluster/timeNdistance/{id} cluster getClusterTimeNDistance
// Get cluster calculated time and distance to complete a cluster using clusterID.
//
// responses:
//	200: getClusterTimeNDistanceResp
//  422: errorValidation
//  501: errorResponse

func (p *Cluster) GetClusterTimeNDistance(w http.ResponseWriter, r *http.Request) {
	p.l.Println("Handle GET request -> cluster-api Module")
	
	vars := mux.Vars(r)
	id := vars["id"]

	lp := data.GetClusterTNDCRUDOPS(id)

	err := lp.ClusterTimeNDistanceResponseToJSON(w)
	if err!=nil {
		http.Error(w,"Data failed to marshel",http.StatusInternalServerError)		
	}
}