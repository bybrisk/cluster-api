package data

import ( //"log"
		"fmt"
		"sync"
		"io/ioutil"
		log "github.com/sirupsen/logrus"
		"net/http"
		//"net/url"
		"encoding/json"
		//"github.com/muesli/clusters"
		//"github.com/muesli/kmeans"
	)

var wg sync.WaitGroup

/*func CreateClusterByID () *[]CreateClusterResponse{

	var resArr []CreateClusterResponse
	respresp
	//get array object from clusterQ
	//loop over the array
	clusterArrObj:=GetClusteQArrayObj()
	for _,CurrClusterObj:=range clusterArrObj.RequestArray{
	
		//Clear all variables

		var res CreateClusterResponse
		
		//all the pending deliveries from es queue where cluster is null
		queueArr := GetPendingDeliveries(CurrClusterObj.BybID)

		var arr []LatLongArespn an array else return error
		//agentIDArr := GetAgentIDArray(CurrClusterObj.NumberOfCluster,CurrClusterObj.BybID)

		for _, hit := range queueArr.Hits.Hits {
			singleArr := LatLongAndID{
				BybID : hit.ID,
				Latitude: hit.Source.Latitude,
				Longitude: hit.Source.Longitude,
				ItemWeight: hit.Source.ItemWeight,
			}
			arr = append(arr, singleArr)
		}

		clusterIdArr,deliveryIdArr,omega,latlng := ModifiedCluster(arr,int(CurrClusterObj.NumberOfCluster))
        
		_=prepareAndSaveClusterID(clusterIdArr,deliveryIdArr,omega,GetAgentIDArray(CurrClusterObj.NumberOfCluster,CurrClusterObj.BybID),latlng,CurrClusterObj.NumberOfCluster)


		we := make([]float64,len(omega))

		for i,v:=range omega{
			for _,u:=range v{
				we[i] = we[i] + MKM_Weights[u]
			}
		}

		fmt.Printf("Weight Array %v",we)

		var clusterIDArrObj ClusterIDArray
		clusterIDArrObj.ClusterID = clusterIdArr

		SaveClusterIDToMongo(clusterIDArrObj,CurrClusterObj.BybID)
		

		if len(deliveryIdArr)==0{
			clusterIDArrObj.ClusterID = nil
			res=CreateClusterResponse{
				ClusterIDArray:clusterIDArrObj ,
				Message: "No Pending Deliveries",
			}
		} else{
			res=CreateClusterResponse{
				ClusterIDArray:clusterIDArrObj ,
				Message: "Clusters Created",
			}
		}

		resArr = append(resArr,res)
		//delete this element from array
		ClearClusterQArray(CurrClusterObj.BybID)

	}


	return &resArr
}*/

func CreateClusterByID () *[]CreateClusterResponse{

	var resArr []CreateClusterResponse
	
	//get array object frorespm clusterQ
	//loop over the array
	clusterArrObj:=GetClusteQArrayObj()
	for _,CurrClusterObj:=range clusterArrObj.RequestArray{
	
		//Clear all variables

		var res CreateClusterResponse
		
		//fetch python API
		pythonResponse, err := http.Get("http://ec2-18-218-54-244.us-east-2.compute.amazonaws.com/api/bybrisk/cluster/create?bybid="+CurrClusterObj.BybID)
		if err != nil {
			log.Error(err)
		}
		//We Read the response body on the line below.
		body, err := ioutil.ReadAll(pythonResponse.Body)
		if err != nil {
			log.Error(err)
		}

		var betaResponseCrossOrigin PythonClusterAPIResponse
		err = json.Unmarshal(body, &betaResponseCrossOrigin)
		
		var clusterIDArrObj ClusterIDArray
		clusterIDArrObj.ClusterID = betaResponseCrossOrigin.ClusterIDArray //response.clusterArr attach
		
		res=CreateClusterResponse{
			ClusterIDArray:clusterIDArrObj ,
			Message: "Clusters Created",
		}

		/*if len(deliveryIdArr)==0{
			clusterIDArrObj.ClusterID = nil
			res=CreateClusterResponse{
				ClusterIDArray:clusterIDArrObj ,
				Message: "No Pending Deliveries",
			}
		} else{
			res=CreateClusterResponse{
				ClusterIDArray:clusterIDArrObj ,
				Message: "Clusters Created",
			}
		}*/

		resArr = append(resArr,res)
		//delete this element from array
		ClearClusterQArray(CurrClusterObj.BybID)

	}


	return &resArr
}

func prepareAndSaveClusterID (clusterIdArr []string,deliveryIdArr []string,omega []MKM_intArr,agentIDArr []AgentIDArrayStruct, latlng []MKM_floatArr,n int64) int64 {
	if (int64(len(agentIDArr))<n) {
		return 0
	} else {
    var geoCodeArr []LatLongAndID

	for i,v:=range omega{
		for _,u:=range v{
			geoCodeArrSingle := LatLongAndID{
				AgentID : agentIDArr[i].BybID,
				BybID:deliveryIdArr[u],
				Latitude:latlng[u][0],
				Longitude:latlng[u][1],
				ClusterID:clusterIdArr[i],
			}
			geoCodeArr = append(geoCodeArr,geoCodeArrSingle)
		}
	}
		SaveClusterID(geoCodeArr)
	}
	return 1
}

func GetAllClusterByID(docID string) *ClusterArrayByIDResponse{
	var response ClusterArrayByIDResponse
	clusterArray := AllClusterIDsByBusinessID(docID)

	clusterDetailArray := AllClusterDetailByID(clusterArray.CurrentClusterArr.ClusterID)

	response = ClusterArrayByIDResponse{
		ClusterIDArray: clusterArray.CurrentClusterArr.ClusterID,
		AssignedDeliveryArray: clusterDetailArray,
	}
	
	return &response
}

func GetSingleClusterByID(docID string) *SingleClusterResponseBulk{
	res := FetchDeliveryByClusterID(docID)
	return &res 
}

func AddToClusterQueue(d *CreateClusterRequest) *CreateClusterResponse{

	var res CreateClusterResponse
	var clusterIDArrObj ClusterIDArray
	//save request to cluster Array in mongodb
	_=SaveToClusterArray(d)
	res=CreateClusterResponse{
		ClusterIDArray: clusterIDArrObj,
		Message: "Clusters Request Queued! Clusters will be created within 5 minutes.",
	}
	return &res
}

func SendRequestToPythonAPI(d *CreateClusterRequest) *CreateClusterResponse{

	var res CreateClusterResponse
		
	//fetch python API
	pythonResponse, err := http.Get("http://ec2-18-218-54-244.us-east-2.compute.amazonaws.com/api/bybrisk/cluster/create?bybid="+d.BybID)
	if err != nil {
		log.Error(err)
	}
	//We Read the response body on the line below.
	body, err := ioutil.ReadAll(pythonResponse.Body)
	if err != nil {
		log.Error(err)
	}

	var betaResponseCrossOrigin PythonClusterAPIResponse
	err = json.Unmarshal(body, &betaResponseCrossOrigin)
		
	var clusterIDArrObj ClusterIDArray
	clusterIDArrObj.ClusterID = betaResponseCrossOrigin.ClusterIDArray //response.clusterArr attach
		
	res=CreateClusterResponse{
		ClusterIDArray:clusterIDArrObj ,
		Message: "Clusters Created",
	}
	
	return &res
}

func GetClusterTNDCRUDOPS(docID string) *ClusterTimeNDistanceResponse{
	var res ClusterTimeNDistanceResponse
	
	//check if the distance and time in mongo against the cluster is not zero.
	
		//Get the sorted delivery from internal API
		bybriskResponse, err := http.Get("https://developers.bybrisk.com/delivery/all/pending/"+docID)
		if err != nil {
				log.Error(err)
			}
		//We Read the response body on the line below.
		body, err := ioutil.ReadAll(bybriskResponse.Body)
		if err != nil {
	    		log.Error(err)
			}


		totalTimeInSec:=0
		totalDistanceInMeters:=0	

		var allDeliveryArr DeliveryResponseBulk
		var arrayOfDeliveryObj []DeliveryWithTimeAndDistance
		err = json.Unmarshal(body, &allDeliveryArr)

		//Check if the process has been cached
		mongoDeliveryObj,err:=GetSavedPathAndDetailFromMongo(allDeliveryArr.Hits.Hits[0].Source.BybID,allDeliveryArr.Hits.Hits[0].Source.ClusterID)
		businessGeocodes := GetGeocodes(allDeliveryArr.Hits.Hits[0].Source.BybID)
		totalStandByDuration := businessGeocodes.Standbyduration*int64(len(allDeliveryArr.Hits.Hits))
		//fmt.Println(allDeliveryArr)

		if err!=nil{	

		originLat :=  fmt.Sprintf("%f", businessGeocodes.Latitude)
		originLng := fmt.Sprintf("%f",businessGeocodes.Longitude)
		FinalDestinationLat :=  originLat
		FinalDestinationLng := originLng

		for _,v:=range allDeliveryArr.SortedIdString{
			for index2,m:= range allDeliveryArr.Hits.Hits{
				if v==m.ID {

					//distance and time from depot to first delivery

					destinationLat1 := fmt.Sprintf("%f",allDeliveryArr.Hits.Hits[index2].Source.Latitude)
					destinationLng1 := fmt.Sprintf("%f",allDeliveryArr.Hits.Hits[index2].Source.Longitude)		
					
					googleResponse, err := http.Get("https://maps.googleapis.com/maps/api/directions/json?origin="+originLat+","+originLng+"&destination="+destinationLat1+","+destinationLng1+"&key=AIzaSyAZDoWPn-emuLvzohH3v-cS_En-u9NSA1A")
					if err != nil {
							log.Error(err)
						}
					//We Read the response body on the line below.
					body, err := ioutil.ReadAll(googleResponse.Body)
					if err != nil {
							log.Error(err)
						}

					var googleResponseCrossOrigin GoogleDirectionAPIStruct
					err = json.Unmarshal(body, &googleResponseCrossOrigin)
					//fmt.Println(googleResponseCrossOrigin)
					if (googleResponseCrossOrigin.Status != "OK") {
						//fmt.Println("Couldn't Fetch this ")
						continue
					}	
					totalTimeInSec = totalTimeInSec + googleResponseCrossOrigin.Routes[0].Legs[0].Duration.Value
					totalDistanceInMeters = totalDistanceInMeters + googleResponseCrossOrigin.Routes[0].Legs[0].Distance.Value	
					//fmt.Printf("origin: %s , %s ; Destination: %s , %s ; name: %s\n",originLat,originLng,destinationLat1,destinationLng1,allDeliveryArr.Hits.Hits[index2].Source.CustomerName)
					originLat =  destinationLat1
					originLng = destinationLng1

					var NewDeliveryObj DeliveryWithTimeAndDistance
					NewDeliveryObj.DeliveryID = m.ID
					NewDeliveryObj.Distance = int64(googleResponseCrossOrigin.Routes[0].Legs[0].Distance.Value)
					NewDeliveryObj.Time = int64(googleResponseCrossOrigin.Routes[0].Legs[0].Duration.Value)
					arrayOfDeliveryObj = append(arrayOfDeliveryObj,NewDeliveryObj)
				}
			} 
		
		}

		//final leg calculation
		googleResponse, err := http.Get("https://maps.googleapis.com/maps/api/directions/json?origin="+originLat+","+originLng+"&destination="+FinalDestinationLat+","+FinalDestinationLng+"&key=AIzaSyAZDoWPn-emuLvzohH3v-cS_En-u9NSA1A")
		if err != nil {
				log.Error(err)
			}
		//We Read the response body on the line below.
		body, err = ioutil.ReadAll(googleResponse.Body)
		if err != nil {
				log.Error(err)
			}

		var googleResponseCrossOrigin GoogleDirectionAPIStruct
		err = json.Unmarshal(body, &googleResponseCrossOrigin)
		totalTimeInSec = totalTimeInSec + googleResponseCrossOrigin.Routes[0].Legs[0].Duration.Value
		totalDistanceInMeters = totalDistanceInMeters + googleResponseCrossOrigin.Routes[0].Legs[0].Distance.Value	
		//fmt.Printf("origin: %s , %s ; Destination: %s , %s ; name: Depot\n",originLat,originLng,FinalDestinationLat,FinalDestinationLng)	
		
		var NewDeliveryObjOuter DeliveryWithTimeAndDistance
		NewDeliveryObjOuter.DeliveryID = "Depot"
		NewDeliveryObjOuter.Distance = int64(googleResponseCrossOrigin.Routes[0].Legs[0].Distance.Value)
		NewDeliveryObjOuter.Time = int64(googleResponseCrossOrigin.Routes[0].Legs[0].Duration.Value)
		arrayOfDeliveryObj = append(arrayOfDeliveryObj,NewDeliveryObjOuter)

		document := MongoStructForTimeAndDistance{
			ArrayOfDeliveryDetail:arrayOfDeliveryObj,
			AgentID:allDeliveryArr.Hits.Hits[0].Source.ClusterID,
		}

		//fmt.Println(arrayOfDeliveryObj)
		//fmt.Println(totalTimeInSec)
		//fmt.Println(totalDistanceInMeters)
	
		//save data to mongo against this agentID
		SavePathAndDetailToMongo(allDeliveryArr.Hits.Hits[0].Source.BybID,allDeliveryArr.Hits.Hits[0].Source.ClusterID,document)

		res=ClusterTimeNDistanceResponse{
			AgentID:allDeliveryArr.Hits.Hits[0].Source.DeliveryAgentID ,
			ClusterTime: int64(totalTimeInSec)+totalStandByDuration,
			ClusterDistance:int64(totalDistanceInMeters),
			Message:"Calculated Time and Distance",

		}
	} else {
		time,distance := getTimeAndDistanceFromCachedData(mongoDeliveryObj,allDeliveryArr.Hits.Hits[0].Source.ClusterID)
		res=ClusterTimeNDistanceResponse{
			AgentID:allDeliveryArr.Hits.Hits[0].Source.DeliveryAgentID ,
			ClusterTime: time+totalStandByDuration,
			ClusterDistance:distance,
			Message:"Cached Time and Distance",

		}
	}
		
	return &res
}

func getTimeAndDistanceFromCachedData(mongoDeliveryObj ExtractTimeAndDistanceFromMongo,agentID string) (int64,int64){
	var time int64 = 0
	var distance int64 = 0
	for _,v:= range mongoDeliveryObj.DeliveryDetailObj{
		if (v.AgentID==agentID){
			for _,val:= range v.ArrayOfDeliveryDetail{
				distance = distance + val.Distance
				time = time + val.Time
			} 
		}
	}
	return time,distance
}