package data_test

import (
	"testing"
	"fmt"
	//"math/rand"
	"github.com/bybrisk/cluster-api/data"
)

func TestCreateCluster(t *testing.T) {

	payload:= &data.CreateClusterRequest {
		BybID : "60250f9d4063fe8843356b82",
		NumberOfCluster : 3,
	}

	res:=data.CreateClusterByID(payload)
	fmt.Println(res)
}

/*func TestGetPendingDeliveries(t *testing.T) {
	res := data.GetPendingDeliveries("600d95c5d72ee5dd5896dd75")
	fmt.Println(res)
}*/

/*func TestClusteringAlgo(t *testing.T) {
	 var arr []data.LatLongOfBusiness

	 for i:=0; i<100;i++{
		singleArr := data.LatLongOfBusiness{
			Latitude:rand.Float64()*5 ,
			Longitude: rand.Float64()*5,
		}
		arr = append(arr, singleArr)
	 }

	 data.ClusteringAlgoFunc(arr,3)
}*/

/*func TestGetAgentIDArray(t *testing.T) {
	arr,err:=data.GetAgentIDArray(2,"6013bc1aeef443c14c31f250")
	if err!=nil{
		t.Fail()
	}
	fmt.Println(arr)
}*/

/*func TestSaveClusterIDToMongo(t *testing.T){
	arr:=[]string{"abc","bcd","123"}
	clusterIDArr:=data.ClusterIDArray{
		ClusterID:arr,
	}
	_=data.SaveClusterIDToMongo(clusterIDArr, "6016e907d3d30bd3d65a565d")
}*/

/*func TestGetAllClusterByID(t *testing.T){
	res:=data.GetAllClusterByID("6016ee473ae10bd996052f15")
	fmt.Println(res)
}*/

/*func TestSingleClusterByID(t *testing.T) {
	res:=data.GetSingleClusterByID("1nrr1kxlRUf42iChV7WMrqZs9L9")
	fmt.Println(res)
}*/