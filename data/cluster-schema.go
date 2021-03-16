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

type CreateClusterRequestArray struct {
	RequestArray []CreateClusterRequest	`json: "RequestArray"`
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

type ClusterArrayByIDResponse struct{

	// Array of current cluster IDs
	//
	ClusterIDArray []string 

	// Array of deliveries assigned to respective clusters
	//
	AssignedDeliveryArray []PendingDeliveryBulk
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

				//Name of the customer placing delivery
				//
				CustomerName    string  `json:"CustomerName"`

				//ClusterID of the cluster this delivery falls into
				//
				ClusterID       string  `json:"clusterID"`

				//Delivery Observed Distance (in meters)
				//
				DistanceObserved  float64  `json:"distanceObserved"`

				//AgentID of the agent associated with the delivery
				//
				DeliveryAgentID string  `json:"deliveryAgentID"`

				//Longitude of delivery location
				//
				Longitude       float64 `json:"longitude"`
			} `json:"_source"`
		} `json:"hits"`
	} `json:"hits"`
}

//get all deliveries Response struct
type SingleClusterResponseBulk struct {
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
				//Pincode of delivery location
				//
				Pincode         string  `json:"pincode"`

				//API Key used in the delivery
				//
				APIKey          string  `json:"apiKey"`

				//Latitude of delivery location
				//
				Latitude        float64 `json:"latitude"`

				//ClusterID of the cluster this delivery falls into
				//
				ClusterID       string  `json:"clusterID"`

				//AgentID of the agent associated with the delivery
				//
				DeliveryAgentID string  `json:"deliveryAgentID"`

				//Phone number of the customer placing delivery
				//
				Phone           string  `json:"phone"`

				//Name of the customer placing delivery
				//
				CustomerName    string  `json:"CustomerName"`

				//Business ID associated with the delivery
				//
				BybID           string  `json:"BybID"`

				//Weight of Item delivered
				//
				ItemWeight      float64     `json:"itemWeight"`

				//Is payment done or not
				//
				PaymentStatus   bool    `json:"paymentStatus"`

				//Status of Delivery
				//
				DeliveryStatus  string  `json:"deliveryStatus"`

				//Delivery Observed Distance (in meters)
				//
				DistanceObserved  string  `json:"distanceObserved"`

				//Address of delivery
				//
				CustomerAddress string  `json:"CustomerAddress"`

				//Longitude of delivery location
				//
				Longitude       float64 `json:"longitude"`

				//Delivery Ranking Time (It will be set using an internal algo)
				//
				RankingTime int64 `json:"rankingTime"`
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
	ItemWeight float64 `json:"itemWeight"`
	ClusterID string
	AgentID string
}

type AgentIDArrayStruct struct {
	BybID string `json:"bybid"`
}

type ClusterIDArray struct {
	ClusterID []string `json:"clusterid"`
}

type ClusterArrayObject struct{
	CurrentClusterArr ClusterIDArray `json:"currentClusterArr"`
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

type PythonClusterAPIResponse struct {
	ClusterIDArray []string `json:"clusterIDArray"`
	Status         string   `json:"status"`
}

func (d *CreateClusterRequest) ValidateCreateClusters() error {
	validate := validator.New()
	return validate.Struct(d)
}

/*GET _search
{
  "query": {
    "match_all": {}
  }
}

GET /_all/_search
{
  "query": {
			"ids": {
				"values": "o4axZ3cBMbaQ18HIDI95"
			}
		}
	}
	
#DELETE /01-29-2021/_doc/BYYSTncBMbaQ18HIMI2q

POST /_all/_search?size=500
{
		"query": {
		  "bool": {
			"filter": [
			  {"term": {"deliveryAgentID": "601a9e411b66a92ffc0f1cf2"}}
			]
		  }
		}
	  }
	  
POST /_all/_search?size=500
{
		"query": {
		  "bool": {
			"filter": [
			  {"term": {"BybID": "601a4c31565e622c7f6816f5"}}
			]
		  }
		}
	  }	  
	  
POST /_all/_update_by_query?conflicts=proceed
 {
		"script" : {
			"source": "ctx._source.clusterID='q123eiropncjknjk';ctx._source.deliveryAgentID='q1evds23eiropncjknjk';",
			"lang": "painless"  
		  },
		  "query": {
			  "ids" : {
		    	"values" : "GIZ0T3cBMbaQ18HIq41k"
			    }
		  }
	}

POST /_all/_update_by_query?conflicts=proceed
 {
		"script" : {
			"source": "ctx._source.clusterID='';ctx._source.deliveryAgentID='';",
			"lang": "painless"  
		  },
		  "query": {
		  "bool": {
			"filter": [
			  {"term": {"BybID": "6015789575219a98b573e497"}}
			]
		  }
		}
	}	
	
POST /_all/_search?size=500
{
		"query": {
		  "bool": {
			"filter": [
			  {"term": {"clusterID.keyword": "1oCVufAyKeaHHIxWKyfoXOlQags"}}
			]
		  }
		}
}	

POST /_all/_update_by_query?conflicts=proceed
{
		"script" : {
			"source": "ctx._source.distanceObserved=1234;",
			"lang": "painless"  
		  },
		  "query": {
			  "ids" : {
			"values" : "GIZ0T3cBMbaQ18HIq41k"
			}
		  }
	  }*/