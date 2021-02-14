package data

import ( //"log"
		"fmt"
		"github.com/segmentio/ksuid"
		"sync"
		//"io/ioutil"
		//"net/http"
		//"net/url"
		//"encoding/json"
		//"github.com/muesli/clusters"
		//"github.com/muesli/kmeans"
	)

var wg sync.WaitGroup

func CreateClusterByID (d *CreateClusterRequest) *CreateClusterResponse{

	var res CreateClusterResponse
	
	//all the pending deliveries from es queue where cluster is null
	queueArr := GetPendingDeliveries(d.BybID)

	var arr []LatLongAndID

	//Get k agentIDs in an array else return error
	agentIDArr,err := GetAgentIDArray(d.NumberOfCluster,d.BybID)
	if err!=nil {
		return &res	
	}

	for _, hit := range queueArr.Hits.Hits {
		singleArr := LatLongAndID{
			BybID : hit.ID,
			Latitude: hit.Source.Latitude,
			Longitude: hit.Source.Longitude,
			ItemWeight: hit.Source.ItemWeight,
		}
		arr = append(arr, singleArr)
	 }

	
	 clusterIdArr,deliveryIdArr,omega,latlng := ModifiedCluster(arr,int(d.NumberOfCluster))

	 prepareAndSaveClusterID(clusterIdArr,deliveryIdArr,omega,agentIDArr,latlng)


	 we := make([]float64,len(omega))

	 for i,v:=range omega{
		 for _,u:=range v{
			we[i] = we[i] + MKM_Weights[u]
		 }
	 }

	 fmt.Printf("Weight Array %v",we)

	 var clusterIDArrObj ClusterIDArray
	 clusterIDArrObj.ClusterID = clusterIdArr

	 SaveClusterIDToMongo(clusterIDArrObj,d.BybID)
	

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


	return &res
}

func prepareAndSaveClusterID (clusterIdArr []string,deliveryIdArr []string,omega []MKM_intArr,agentIDArr []AgentIDArrayStruct, latlng []MKM_floatArr) {
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

func ClusteringAlgoFunc (arrLatLong []LatLongAndID, k int,agentIDArray []AgentIDArrayStruct,accountID string) ClusterIDArray {
	
	var geoCodeArr []LatLongAndID
	var clusterIDArr ClusterIDArray
	var d Observations
	d = arrLatLong

	fmt.Printf("%d data points\n", len(d))

	m := Kmeans{
		deltaThreshold:     0.01,
		iterationThreshold: 96,
	}
	clusters, _ := m.Partition(d, k)

	for i, c := range clusters {
		clusterId := ksuid.New()
		clusterIDArr.ClusterID=append(clusterIDArr.ClusterID,clusterId.String())
		agentIdContext := agentIDArray[i].BybID
		fmt.Printf("Cluster %d : %s\n",i, clusterId.String())
		for _, dataPoints :=range c.Observations {
			geoCodeArrSingle := LatLongAndID{
				AgentID : agentIdContext,
				BybID:dataPoints.BybID,
				Latitude:dataPoints.Latitude,
				Longitude:dataPoints.Longitude,
				ClusterID:clusterId.String(),
			}
			geoCodeArr = append(geoCodeArr,geoCodeArrSingle)
		}
	}

	//run this function in background
	SaveClusterID(geoCodeArr)
	SaveClusterIDToMongo(clusterIDArr,accountID)
	return clusterIDArr
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


