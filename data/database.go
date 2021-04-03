package data

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"encoding/json"
	//"bytes"
	//"github.com/bybrisk/structs"
	"go.mongodb.org/mongo-driver/bson"
	"github.com/shashank404error/shashankMongo"
	"go.mongodb.org/mongo-driver/bson/primitive"
	log "github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var resultID string 

var apiKey string = "AIzaSyAZDoWPn-emuLvzohH3v-cS_En-u9NSA1A"

func GetGeocodes (docID string) LatLongOfBusiness {
	collectionName := shashankMongo.DatabaseName.Collection("businessAccounts")
	id, _ := primitive.ObjectIDFromHex(docID)
	filter := bson.M{"_id": id}

	var document LatLongOfBusiness

	err:= collectionName.FindOne(shashankMongo.CtxForDB, filter).Decode(&document)
	if err != nil {
		log.Error("GetGeocodes ERROR:")
		log.Error(err)
	}
	return document
}	

func GetDistanceFromBusiness(queueArr PendingDeliveryBulk ,geoCodes LatLongOfBusiness) []DeliveryHitsArr {
	//Fetch distance between business and deliveries using google matrix api

	var middleArrs []DeliveryHitsArr
	for _, hit := range queueArr.Hits.Hits {

		destinationLat := fmt.Sprintf("%f", hit.Source.Latitude)
		destinationLng := fmt.Sprintf("%f", hit.Source.Longitude)

		originLat := fmt.Sprintf("%f", geoCodes.Latitude)
		originLng := fmt.Sprintf("%f", geoCodes.Longitude)

		url := "https://maps.googleapis.com/maps/api/distancematrix/json?origins="+originLat+","+originLng+"&destinations="+destinationLat+","+destinationLng+"&key="+apiKey
		
		//get distance and time using latlng
		response, err := http.Get(url)

		if err != nil {
			fmt.Print(err.Error())
		}

		responseData, err := ioutil.ReadAll(response.Body)
		if err != nil {
			log.Fatal(err)
		}
		var responseObject GoogleDistanceMatrix

		json.Unmarshal(responseData, &responseObject)
		var middleArr *DeliveryHitsArr
		middleArr = &DeliveryHitsArr{
			DeliveryID : hit.ID,
			Weight : hit.Source.ItemWeight,
			DistanceValue :responseObject.Rows[0].Elements[0].Distance.Text,
			ETA:responseObject.Rows[0].Elements[0].Duration.Text,	
		} 
		
		fmt.Println(middleArr)
		middleArrs = append(middleArrs,*middleArr)

	}
	
	return middleArrs
}

func GetAgentIDArray(k int64, docID string) ([]AgentIDArrayStruct){
	var agentIDs []AgentIDArrayStruct

	options := options.Find()
	options.SetLimit(k)

	collectionName := shashankMongo.DatabaseName.Collection("agents")

	cursor, err := collectionName.Find(shashankMongo.CtxForDB, bson.M{"businessid":docID},options)
	if err != nil {
		log.Error("AllAgentsByBusinessID Read ERROR : ")
		log.Error(err)
		return agentIDs
	}
	if err = cursor.All(shashankMongo.CtxForDB, &agentIDs); err != nil {
		log.Error("AllAgentsByBusinessID Write ERROR : ")
		log.Error(err)
		return agentIDs
	}

	return agentIDs
}

func SaveClusterIDToMongo(clusterIDArr ClusterIDArray, docID string) int64 {
	collectionName := shashankMongo.DatabaseName.Collection("cluster")
	//id, _ := primitive.ObjectIDFromHex(account.BybID)
	update := bson.M{"$set":bson.M{"currentClusterArr": clusterIDArr}}
	filter := bson.M{"bybid": docID}
	res,err := collectionName.UpdateOne(shashankMongo.CtxForDB,filter, update)
	if err!=nil{
		log.Error("SaveClusterIDToMongo ERROR:")
		log.Error(err)
		}	
	
	return res.ModifiedCount
}

func AllClusterIDsByBusinessID(docID string) ClusterArrayObject{
	collectionName := shashankMongo.DatabaseName.Collection("cluster")
	filter := bson.M{"bybid": docID}

	var document ClusterArrayObject

	err:= collectionName.FindOne(shashankMongo.CtxForDB, filter).Decode(&document)
	if err != nil {
		log.Error("GetGeocodes ERROR:")
		log.Error(err)
	}
	return document
}

func SaveToClusterArray(d *CreateClusterRequest) int64 {
	collectionName := shashankMongo.DatabaseName.Collection("clusterQ")
	id, _ := primitive.ObjectIDFromHex("602cc9ef0a7d47656411f63a")
	filter := bson.M{"_id": id}
	updateResult, err := collectionName.UpdateOne(shashankMongo.CtxForDB, filter, bson.D{{Key: "$push", Value: bson.M{"RequestArray": d}}})
	if err != nil {
		log.Error("SaveToClusterArray ERROR:")
		log.Error(err)
	}
	fmt.Println(updateResult.ModifiedCount)
	return updateResult.ModifiedCount
}

func GetClusteQArrayObj() CreateClusterRequestArray{
	collectionName := shashankMongo.DatabaseName.Collection("clusterQ")
	id, _ := primitive.ObjectIDFromHex("602cc9ef0a7d47656411f63a")
	filter := bson.M{"_id": id}

	var document CreateClusterRequestArray

	err:= collectionName.FindOne(shashankMongo.CtxForDB, filter).Decode(&document)
	if err != nil {
		log.Error("GetClusteQArrayObj ERROR:")
		log.Error(err)
	}
	return document
}

func ClearClusterQArray(bybid string) {
	collectionName := shashankMongo.DatabaseName.Collection("clusterQ")
	id, _ := primitive.ObjectIDFromHex("602cc9ef0a7d47656411f63a")
	filter := bson.M{"_id": id}
	_, err := collectionName.UpdateOne(shashankMongo.CtxForDB, filter, bson.M{"$pull":bson.M{"RequestArray": bson.M{"bybid":bybid}}})
	if err != nil {
		log.Error("ClearClusterQArray ERROR:")
		log.Error(err)
	}
}

func SavePathAndDetailToMongo(BybID string,docID string,mongoObj MongoStructForTimeAndDistance) int64{
	collectionName := shashankMongo.DatabaseName.Collection("cluster")
	//id, _ := primitive.ObjectIDFromHex("602cc9ef0a7d47656411f63a")
	filter := bson.M{"bybid": BybID}
	docIDPresent := docID+"Present"
	updateResult, err := collectionName.UpdateOne(shashankMongo.CtxForDB, filter, bson.D{{Key: "$push", Value: bson.M{"DeliveryDetailObj": mongoObj}}})
	if err != nil {
		log.Error("SavePathAndDetailToMongo1 ERROR:")
		log.Error(err)
	}

	_, err = collectionName.UpdateOne(shashankMongo.CtxForDB, filter, bson.D{{Key: "$set", Value: bson.M{docIDPresent:"True"}}})
	if err != nil {
		log.Error("SavePathAndDetailToMongo2 ERROR:")
		log.Error(err)
	}
	fmt.Println(updateResult.ModifiedCount)
	return updateResult.ModifiedCount
}

func GetSavedPathAndDetailFromMongo(BybID string,docID string) (ExtractTimeAndDistanceFromMongo,error){
	collectionName := shashankMongo.DatabaseName.Collection("cluster")
	docIDPresent := docID+"Present"
	filter := bson.M{docIDPresent: "True"}

	var document ExtractTimeAndDistanceFromMongo

	err:= collectionName.FindOne(shashankMongo.CtxForDB, filter).Decode(&document)
	if err != nil {
		log.Error("GetGeocodes ERROR:")
		log.Error(err)
	}
	return document,err
}