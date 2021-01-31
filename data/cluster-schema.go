package data

import (
	"github.com/go-playground/validator/v10"
	//"github.com/bybrisk/structs"
)

//post request for creating clusters
type CreateClusterRequest struct{

	// BybID of the business which needs its deliveries converted to clusters
	//
	// required: true
	// max length: 1000
	BybID string `json;"bybID" validate:"required"` 

	// Number of clusters you want to divide the deliveries to
	//
	// required: true
	NumberOfCluster int64 `json: "numberOfCluster" validate:"required"`
}

//Response for creating clusters
type CreateClusterResponse struct{

	// Array of newly created cluster IDs
	//
	ClusterIDArray ClusterIDArray 

	// Message to be shown to the user
	//
	Message string
}

//struct for getting the pending deliveries to convert into clusters
type PendingDeliveryBulk struct {
	Hits struct {
		Hits []struct {
			//Date of delivery
			//
			Index  string `json:"_index"`

			//ID of delivery
			//
			ID     string `json:"_id"`

			//Delivery details
			//
			Source struct {
				//Latitude of delivery location
				//
				Latitude        float64 `json:"latitude"`

				//Weight of Item delivered
				//
				ItemWeight      float64     `json:"itemWeight"`

				//Longitude of delivery location
				//
				Longitude       float64 `json:"longitude"`
			} `json:"_source"`
		} `json:"hits"`
	} `json:"hits"`
}

type DeliveryHitsArr struct{
	DeliveryID string
	Weight float64
	DistanceValue string
	ETA string
	ClusterID string
	AgentID string 
}


//helper struct
type LatLongOfBusiness struct {
	Latitude float64 `json:"latitude"`
	Longitude float64 `json:"longitude"`
}

type LatLongAndID struct {
	BybID string
	Latitude float64 `json:"latitude"`
	Longitude float64 `json:"longitude"`
	ClusterID string
	AgentID string
}

type AgentIDArrayStruct struct {
	BybID string `json:"bybid"`
}

type ClusterIDArray struct {
	ClusterID []string `json:"clusterID"`
}

//struct to store google distance matrix data
type GoogleDistanceMatrix struct {
	Rows []struct {
		Elements []struct {
			Distance struct {
				Text  string `json:"text"`
				Value int    `json:"value"`
			} `json:"distance"`
			Duration struct {
				Text  string `json:"text"`
				Value int    `json:"value"`
			} `json:"duration"`
			Status string `json:"status"`
		} `json:"elements"`
	} `json:"rows"`
	Status string `json:"status"`
}

func (d *CreateClusterRequest) ValidateCreateClusters() error {
	validate := validator.New()
	return validate.Struct(d)
}