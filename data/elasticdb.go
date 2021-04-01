package data

import (
	"encoding/json"
	"fmt"
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

var Elasticurl string = "https://search-krayfin-ewgnw3zevhpuznh2kvo5hbbmtq.us-east-2.es.amazonaws.com"
var UsernameElastic string = "elastic"
var Elasticpassword string = "K8txmFf6hwnGvaNs7HxNcg2w$"
var urlAuthenticate string = "https://elastic:K8txmFf6hwnGvaNs7HxNcg2w$@search-krayfin-ewgnw3zevhpuznh2kvo5hbbmtq.us-east-2.es.amazonaws.com"
var awsYearIndex string = "/*2021"

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
					   {"term" : {"deliveryStatus.keyword" : "Pending" }}
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
	 resp, err := http.Post(urlAuthenticate+awsYearIndex+"/_update_by_query?conflicts=proceed", "application/json", responseBody)
  
	 //Handle Error
	 if err != nil {
		log.Fatalf("An Error Occured %v", err)
	 }
	 defer resp.Body.Close()

	}
}

func SaveClusterID(arr []LatLongAndID){
	
	for i, d := range arr {
		fmt.Printf("\r %.0f %c", float64(i*100/len(arr)),'%')
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
		 resp, err := http.Post(urlAuthenticate+awsYearIndex+"/_update_by_query?conflicts=proceed", "application/json", responseBody)
	  
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