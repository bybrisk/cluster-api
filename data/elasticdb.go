package data

import (
	"encoding/json"
	//"fmt"
	//"time"
	log "github.com/sirupsen/logrus"
	"bytes"
	//"context"
	"net/http"
	"io/ioutil"
	//"github.com/elastic/go-elasticsearch/v8"
	//"github.com/elastic/go-elasticsearch/v8/esutil"
	//"github.com/mitchellh/mapstructure"
	
)

var Elasticurl string = "https://67b69d7c039a4780b5aa5abbecb09c1a.ap-south-1.aws.elastic-cloud.com:9243"
var UsernameElastic string = "elastic"
var Elasticpassword string = "OmantonZe8l1MaugY6ypelSE"
var urlAuthenticate string = "https://elastic:OmantonZe8l1MaugY6ypelSE@67b69d7c039a4780b5aa5abbecb09c1a.ap-south-1.aws.elastic-cloud.com:9243"

var (
	clusterURLs = []string{Elasticurl}
	username    = UsernameElastic
	password    = Elasticpassword
  )

func GetPendingDeliveries(docID string) PendingDeliveryBulk {
	var deliveries PendingDeliveryBulk

	postBody:=`{
		"query": {
				  "bool": {
					"filter": [
							 { "term" : {"BybID" : "`+docID+`" }},
					   {"term" : {"deliveryStatus.keyword" : "Pending" }},
					   {"term" : {"clusterID.keyword" : "" }}
					]
				  }
				}
		}`

	 responseBody := bytes.NewBufferString(postBody)
  	//Leverage Go's HTTP Post function to make request
	 resp, err := http.Post(urlAuthenticate+"/_all/_search?size=5000", "application/json", responseBody)
  
	 //Handle Error
	 if err != nil {
		log.Fatalf("An Error Occured %v", err)
	 }
	 defer resp.Body.Close()

	 body, err := ioutil.ReadAll(resp.Body)
	 if err != nil {
		log.Error("ReadAll ERROR : ")
		log.Error(err)
	 }
	 
	 err = json.Unmarshal(body, &deliveries)
	 if err != nil {
		log.Error("json.Unmarshal ERROR : ")
		log.Error(err)
    	} 
	return deliveries
}

func SaveTimeNDistanceES(arr []DeliveryHitsArr) {
	
	for _, d := range arr {
	
	//Encode the data
	postBody:=`{
		"script" : {
			"source": "ctx._source.ETA='`+d.ETA+`';ctx._source.distanceValue='`+d.DistanceValue+`';",
			"lang": "painless"  
		  },
		  "query": {
			  "ids" : {
			"values" : "`+d.DeliveryID+`"
			}
		  }
	  }`

	 responseBody := bytes.NewBufferString(postBody)
  	//Leverage Go's HTTP Post function to make request
	 resp, err := http.Post(urlAuthenticate+"/_all/_update_by_query?conflicts=proceed", "application/json", responseBody)
  
	 //Handle Error
	 if err != nil {
		log.Fatalf("An Error Occured %v", err)
	 }
	 defer resp.Body.Close()

	}
}

func SaveClusterID(arr []LatLongAndID) {
	
	for _, d := range arr {
	
	//Encode the data
	postBody:=`{
		"script" : {
			"source": "ctx._source.clusterID='`+d.ClusterID+`';ctx._source.deliveryAgentID='`+d.AgentID+`';",
			"lang": "painless"  
		  },
		  "query": {
			  "ids" : {
		    	"values" : "`+d.BybID+`"
			    }
		  }
	}`

	 responseBody := bytes.NewBufferString(postBody)
  	//Leverage Go's HTTP Post function to make request
	 resp, err := http.Post(urlAuthenticate+"/_all/_update_by_query?conflicts=proceed", "application/json", responseBody)
  
	 //Handle Error
	 if err != nil {
		log.Fatalf("An Error Occured %v", err)
	 }
	 defer resp.Body.Close()

	}
}

func AllClusterDetailByID(arr []string) []PendingDeliveryBulk{
	var clusterMetaArray []PendingDeliveryBulk
	for _,c := range arr{
		var singleCluster PendingDeliveryBulk
		
		postBody:=`{
			"query": {
			  "bool": {
				"filter": [
				  {"term": {"clusterID.keyword": "`+c+`"}}
				]
			  }
			}
	}`
	
		 responseBody := bytes.NewBufferString(postBody)
		  //Leverage Go's HTTP Post function to make request
		 resp, err := http.Post(urlAuthenticate+"/_all/_search?size=5000", "application/json", responseBody)
	  
		 //Handle Error
		 if err != nil {
			log.Fatalf("An Error Occured %v", err)
		 }
		 defer resp.Body.Close()
	
		 body, err := ioutil.ReadAll(resp.Body)
		 if err != nil {
			log.Error("ReadAll ERROR : ")
			log.Error(err)
		 }
		 
		 err = json.Unmarshal(body, &singleCluster)
		 if err != nil {
			log.Error("json.Unmarshal ERROR : ")
			log.Error(err)
			} 

		clusterMetaArray = append(clusterMetaArray,singleCluster)
	}
	return clusterMetaArray
}

func FetchDeliveryByClusterID(docID string) SingleClusterResponseBulk{
	var singleCluster SingleClusterResponseBulk
		
		postBody:=`{
			"query": {
			  "bool": {
				"filter": [
				  {"term": {"clusterID.keyword": "`+docID+`"}}
				]
			  }
			}
	}`
	
		 responseBody := bytes.NewBufferString(postBody)
		  //Leverage Go's HTTP Post function to make request
		 resp, err := http.Post(urlAuthenticate+"/_all/_search?size=5000", "application/json", responseBody)
	  
		 //Handle Error
		 if err != nil {
			log.Fatalf("An Error Occured %v", err)
		 }
		 defer resp.Body.Close()
	
		 body, err := ioutil.ReadAll(resp.Body)
		 if err != nil {
			log.Error("ReadAll ERROR : ")
			log.Error(err)
		 }
		 
		 err = json.Unmarshal(body, &singleCluster)
		 if err != nil {
			log.Error("json.Unmarshal ERROR : ")
			log.Error(err)
			}

		return	singleCluster
}