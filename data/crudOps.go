package data

import ( //"log"
		"fmt"
		"sync"
		//"io/ioutil"
		//"net/http"
		//"net/url"
		//"encoding/json"
		//"github.com/muesli/clusters"
		//"github.com/muesli/kmeans"
	)

var wg sync.WaitGroup

func CreateClusterByID () *[]CreateClusterResponse{

	var resArr []CreateClusterResponse
	
	//get array object from clusterQ
	//loop over the array
	clusterArrObj:=GetClusteQArrayObj()
	for _,CurrClusterObj:=range clusterArrObj.RequestArray{
	
		//Clear all variables

		var res CreateClusterResponse
		
		//all the pending deliveries from es queue where cluster is null
		queueArr := GetPendingDeliveries(CurrClusterObj.BybID)

		var arr []LatLongAndID

		//Get k agentIDs in an array else return error
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

	//clear queue array from the clusterQ
	/*if resArr!=nil {
		ClearClusterQArray()
	}*/


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

