package data_test

import (
	"testing"
	"fmt"
	//"math/rand"
	"github.com/bybrisk/cluster-api/data"
)

/*func TestCreateCluster(t *testing.T) {

	res:=data.CreateClusterByID()
	fmt.Println(res)
}*/

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

/*func TestAddToClusterQ (t *testing.T) {
	payload:= &data.CreateClusterRequest {
		BybID : "6038bd0fc35e3b8e8bd9f81a",
		NumberOfCluster : 2,
	}
	res:=data.SendRequestToPythonAPI(payload)
	fmt.Println(res)
}*/

/*func TestGetClusterQ(t *testing.T) {
	res:=data.GetClusteQArrayObj()
	fmt.Println(res)
}*/

/*func TestClearArr(t *testing.T) {
	data.ClearClusterQArray()
}*/

/*func TestClusterDetailByID(t *testing.T){
	res:=data.FetchDeliveryByClusterID("bcdd0fc3-7e2a-4158-bcd4-c65bffc6b53f")
	fmt.Println(res)
}*/

func TestGetClusterTNDCRUDOPS(t *testing.T){
	res:= data.GetClusterTNDCRUDOPS("60780951ec2eb2585f154498")
	fmt.Println(res)
}