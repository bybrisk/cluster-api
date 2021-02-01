package data

import (
	"encoding/json"
	"io"
)	
func (d *CreateClusterRequest) FromJSONToCreateClusterStruct (r io.Reader) error {
	e := json.NewDecoder(r)
	return e.Decode(d)
}

func (d *CreateClusterResponse) CreateClusterToJSON (w io.Writer) error {
	e := json.NewEncoder(w)
	return e.Encode(d)
}

func (d *ClusterArrayByIDResponse) AllClusterResponseToJSON (w io.Writer) error {
	e := json.NewEncoder(w)
	return e.Encode(d)
}

func (d *SingleClusterResponseBulk) SingleClusterResponseToJSON (w io.Writer) error {
	e := json.NewEncoder(w)
	return e.Encode(d)
}