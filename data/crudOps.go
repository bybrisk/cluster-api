package data

import ( //"log"
		"fmt"
		"github.com/segmentio/ksuid"
		//"io/ioutil"
		//"net/http"
		//"net/url"
		//"encoding/json"
		//"github.com/muesli/clusters"
		//"github.com/muesli/kmeans"
	)

func CreateClusterByID (d *CreateClusterRequest) *CreateClusterResponse{

	var res CreateClusterResponse
	//lat-long of the business
	//geoCodes := GetGeocodes(d.BybID)
	
	//all the pending deliveries from es queue where cluster is null
	queueArr := GetPendingDeliveries(d.BybID)

	var arr []LatLongAndID

	for _, hit := range queueArr.Hits.Hits {
		singleArr := LatLongAndID{
			BybID : hit.ID,
			Latitude: hit.Source.Latitude,
			Longitude: hit.Source.Longitude,
		}
		arr = append(arr, singleArr)
	 }
	NOC := int(d.NumberOfCluster)

	//Get k agentIDs in an array else return error
	agentIDArray,err := GetAgentIDArray(d.NumberOfCluster,d.BybID)
	if err!=nil {
		return &res	
	}
	
	// Feed the array and the no. of clusters to the algo
	clusterArr := ClusteringAlgoFunc(arr,NOC,agentIDArray,d.BybID)

	res=CreateClusterResponse{
		ClusterIDArray:clusterArr ,
		Message: "Clusters Created",
	}

	return &res
}

func ClusteringAlgoFunc (arrLatLong []LatLongAndID, k int,agentIDArray []AgentIDArrayStruct,accountID string) ClusterIDArray {
	//geoCodeArr []LatLongStruct, k int64
	var geoCodeArr []LatLongAndID
	var clusterIDArr ClusterIDArray
	var d Observations
	d = arrLatLong

	fmt.Printf("%d data points\n", len(d))

	// Partition the data points into 7 clusters
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
	SaveClusterID(geoCodeArr)
	SaveClusterIDToMongo(clusterIDArr,accountID)
	//fmt.Println(clusterIDArr)
	return clusterIDArr
}