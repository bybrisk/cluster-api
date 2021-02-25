package data

import (
	"fmt"
	"math"
	"github.com/segmentio/ksuid"
	"sort")

var MKM_C []int64
type MKM_intArr []int64
var MKM_NAvg int
var MKM_Omega []MKM_intArr
var MKM_Index []int64
var MKM_Index0 []int64
var MKM_DeliveryID []string
var MKM_ClusterID []string
var MKM_Weights []float64
var MKM_Weights_Sum float64
type MKM_floatArr []float64
var MKM_LatLong []float64
var MKM_DataPoints []MKM_floatArr

type kv struct {
	Key   int64
	Value float64
}


func ModifiedCluster (geoCodeArr []LatLongAndID,k int) ([]string, []string, []MKM_intArr, []MKM_floatArr ) {
	//emptying the variables
	MKM_C = nil
	MKM_NAvg=0
	MKM_Omega = nil
	MKM_Index = nil
	MKM_Index0 = nil
	MKM_DeliveryID = nil
	MKM_ClusterID = nil
	MKM_Weights = nil
	MKM_Weights_Sum = 0
	MKM_LatLong = nil
	MKM_DataPoints = nil

	C1,I0,I,Navg := InitialiseFirstCluster(geoCodeArr,k)
	MKM_Omega = append(MKM_Omega,C1)

	MKM_OmegaNew:= InitiliseRestClusters(I,I0,k)

	MKM_OmegaNew1 := AssignPoints(I,MKM_OmegaNew)
	//fmt.Println("Assign Points")
	//fmt.Println(MKM_OmegaNew1)

	//final_Omega:=AdmissibleMove(MKM_OmegaNew1,k,Navg)
	final_Omega:=AdmissibleMoveImproved(&MKM_OmegaNew1,k,Navg)
	//fmt.Println()
	//fmt.Println("After Admissible move")
	//fmt.Println(final_Omega)
	return MKM_ClusterID,MKM_DeliveryID,*final_Omega,MKM_DataPoints
}

func InitialiseFirstCluster(geoCodeArr []LatLongAndID,k int) ([]int64,[]int64,[]int64,int){
	
	MKM_C = nil
	MKM_Index0 = nil
	// filled all variables
	for i,v := range geoCodeArr{
		MKM_LatLong = nil
		MKM_Index = append(MKM_Index,int64(i))
		MKM_DeliveryID = append(MKM_DeliveryID,v.BybID)
		MKM_Weights = append(MKM_Weights,v.ItemWeight)
		MKM_Weights_Sum = MKM_Weights_Sum+v.ItemWeight

		MKM_LatLong = append(MKM_LatLong,v.Latitude)
		MKM_LatLong = append(MKM_LatLong,v.Longitude)
		
		MKM_DataPoints = append(MKM_DataPoints,MKM_LatLong)
	}
	MKM_NAvg = (len(MKM_Index)/k)
	//fmt.Println(MKM_Index)
	//fmt.Println(MKM_DeliveryID)
	//fmt.Println(MKM_DataPoints)
	if (len(MKM_DeliveryID)!=0) {
		clusterId := ksuid.New()
		MKM_ClusterID = append(MKM_ClusterID,clusterId.String())
	}
	fmt.Printf("Total Pending Delivery %d\n",len(MKM_DeliveryID))
	fmt.Printf("Avg. Delivery %d\n",MKM_NAvg)

	firstCentroidIndex := FindCentriodIndex(MKM_DataPoints)
	MKM_C = append(MKM_C,int64(firstCentroidIndex))
	MKM_Index0 = append(MKM_Index0,int64(firstCentroidIndex))
	return MKM_C,MKM_Index0,MKM_Index,MKM_NAvg
}

func InitiliseRestClusters(I []int64,I0 []int64,k int) []MKM_intArr {
	for h:=2; h<k+1;h++{
		MKM_C = nil
		m := make(map[int]float64)
		var point int = 0
		var maxDistance float64 = 0
		for _,v:= range I{
			var minDistance float64 = 999999999999
			var distance float64
			for _,u:=range I0{
				
				distance = getDistance(MKM_DataPoints[v][0],MKM_DataPoints[v][1],MKM_DataPoints[u][0],MKM_DataPoints[u][1])
				if distance < minDistance {
					minDistance = distance
				}
			}
			m[int(v)]=minDistance
		}
		for key,element:=range m{
			if element>maxDistance {
				maxDistance=element
				point = key
			}
		}
	
		//get the point and add it to the cluster
		I0 = append(I0,int64(point))
		MKM_C = append(MKM_C,int64(point))
		//Add the cluster to omega
		MKM_Omega = append(MKM_Omega,MKM_C)
		clusterId := ksuid.New()
		MKM_ClusterID = append(MKM_ClusterID,clusterId.String())
	}
	return MKM_Omega
}

func AssignPoints(I []int64,o1 []MKM_intArr)  []MKM_intArr {
	var cluster_index int 
	for _,v:=range I{
		var minDistance float64 = 999999999999999
		centroidArr:=FindCentriodVal(o1)
		for j,u:=range centroidArr{
			distance := getDistance(MKM_DataPoints[v][0],MKM_DataPoints[v][1],u[0],u[1])
			//get cluster_index having minimum distance
			if (distance<minDistance) {
				minDistance = distance
				cluster_index = j
			}
		}
		//add that data_index to the cluster_index
			o1[cluster_index] = append(o1[cluster_index],int64(v))
	}
	for a,b:=range o1{
		b=removeIndex(b, 0)
		o1[a]=b
	}
	return o1
}

func FindCentriodIndex(MKM_DataPoints []MKM_floatArr) int{
	
	cc := make([]float64, 2)

	for _,v := range MKM_DataPoints{
		cc[0] = cc[0] + v[0]
		cc[1] = cc[1] + v[1]
	}
	
	cc[0] = cc[0]/float64(len(MKM_DataPoints))
	cc[1] = cc[1]/float64(len(MKM_DataPoints))
	
	var minDistance float64 = 999999999999 
	var minDistanceIndex int
	for i,v := range MKM_DataPoints{
		distance := getDistance(cc[0],cc[1],v[0],v[1])
		if distance < minDistance {
			minDistance = distance
			minDistanceIndex = i
		}
	}
	return minDistanceIndex
}

func FindCentriodVal(MKM_OmegaNew []MKM_intArr) []MKM_floatArr{
	
	MKM_CentroidArray := make([]MKM_floatArr, len(MKM_OmegaNew))
    for i,v:= range MKM_OmegaNew{
		cc := make([]float64, 2)
		for _,u:= range v{
			cc[0] = cc[0] + MKM_DataPoints[u][0]
			cc[1] = cc[1] + MKM_DataPoints[u][1]
		}
		cc[0] = cc[0]/float64(len(v))
		cc[1] = cc[1]/float64(len(v))
		// add it to centroid array
		MKM_CentroidArray[i] = cc
	}

	return MKM_CentroidArray
}

func FindCentriodSingle(singleOmega []int64) []float64{

	cc := make([]float64, 2)
    for _,v:= range singleOmega{
		cc[0] = cc[0] + MKM_DataPoints[v][0]
		cc[1] = cc[1] + MKM_DataPoints[v][1]
	}
	cc[0] = cc[0]/float64(len(singleOmega))
	cc[1] = cc[1]/float64(len(singleOmega))

	return cc
}

func getDistance(lat1 float64,long1 float64,lat2 float64,long2 float64) float64{
	var r float64
	r = math.Pow(lat1-lat2, 2) + math.Pow(long1-long2, 2)
	return r
}

func balenceValue(Cr []int64,Cs []int64, NAvg int) int64{
	if (len(Cr)>NAvg && len(Cs)<NAvg){
		return 2
	} 
	if (len(Cr)<NAvg && len(Cs)>NAvg && (len(Cr)>len(Cs)+1)) {
		return 1
	} 
	if (len(Cr)==len(Cs)+1) {
		return 0
	} 
	if (len(Cr)<len(Cs) || len(Cr)==len(Cs)) {
		return -1
	}
	return -2
}

func weightValue(Cr []int64,Cs []int64,CrPi []int64,CsPi []int64, k int) float64 {
	Wavg := MKM_Weights_Sum/float64(k)
	var WCr float64
	var WCs float64
	var WCrPi float64
	var WCsPi float64
	for _,b:= range Cr{
		WCr = WCr + MKM_Weights[b]  
	}
	for _,b:= range Cs{
		WCs = WCs + MKM_Weights[b]  
	}
	for _,b:= range CrPi{
		WCrPi = WCrPi + MKM_Weights[b]  
	}
	for _,b:= range CsPi{
		WCsPi = WCsPi + MKM_Weights[b]  
	}
	WeightValueRes := (WCr - Wavg) + (WCs - Wavg) - (WCrPi - Wavg) - (WCsPi - Wavg)
	return WeightValueRes
}

func removeIndex(a []int64, i int) []int64 {
	copy(a[i:], a[i+1:]) // Shift a[i+1:] left one index.
    a[len(a)-1] = 0     // Erase last element (write zero value).
    a = a[:len(a)-1]
	return a
}

func removeItem(items []int64, item int) []int64 {
	var newitems []int64

    for _, i := range items {
        if i != int64(item) {
            newitems = append(newitems, i)
        }
    }

    return newitems
}

/*func AdmissibleMoveImproved(oNew *[]MKM_intArr, k int,Navg int) *[]MKM_intArr {
	var counter int64

	for i:=0;i<len(*oNew);i++{
		//fmt.Printf("Value of ONew After one cluster [%d] %v\n",i,*oNew)
		
		for j:=0;j<len(*oNew);j++{
			if j!=i {

				CrPi := (*oNew)[i]
				//iterate over Cr
				for l:=0;l<len((*oNew)[i]);l++{
					//fmt.Printf("Value in Cr [%d] %d\n",k,(*oNew)[i][l])
					//fmt.Printf("Value in Cr [%d] %d\n",k,CrPi)
					valShifted := (*oNew)[i][l]
					CrPi=removeItem((*oNew)[i],int(valShifted))
					//fmt.Printf("After Deleting ------------->%v\n",CrPi)
					var CsPi []int64

					//iterate through Cs
					for m:=0;m<len((*oNew)[j]);m++{
						CsPi=append((*oNew)[j],valShifted)	
						//fmt.Printf("After Adding ------------->%v\n",CsPi)

						//Balencing condition
						if (balenceValue((*oNew)[i],(*oNew)[j], Navg) > 0 || (balenceValue((*oNew)[i],(*oNew)[j], Navg)>0 || balenceValue((*oNew)[i],(*oNew)[j], Navg)==0 && weightValue((*oNew)[i],(*oNew)[j],CrPi,CsPi,k)>0)) {
							counter = counter+1
							(*oNew)[i]=CrPi
							(*oNew)[j]=CsPi
							break
						}
						CsPi=nil
					}

				}
			}
		}
	}
	return oNew
}*/

func AdmissibleMoveImproved(oNew *[]MKM_intArr, k int,Navg int) *[]MKM_intArr {
	var counter int64

	for i:=0;i<len(*oNew);i++{
		//fmt.Printf("Value of ONew After one cluster [%d] %v\n",i,*oNew)
		
		for j:=0;j<len(*oNew);j++{
			if j!=i {

				// centroid of (*oNew)[j]
				ojCentroid:=FindCentriodSingle((*oNew)[j])
				//fmt.Printf("centroid [%d] : %v",i+j,ojCentroid)

				//sort (*oNew)[i] on the basis of distance from centroid
				distanceFromCentroidMap := make(map[int64]float64)
				for fIndex:=0;fIndex<len((*oNew)[i]);fIndex++{
					distanceFromCentroid := getDistance(MKM_DataPoints[(*oNew)[i][fIndex]][0],MKM_DataPoints[(*oNew)[i][fIndex]][1],ojCentroid[0],ojCentroid[1])
					distanceFromCentroidMap[(*oNew)[i][fIndex]] = distanceFromCentroid
				}
				//fmt.Printf("distanceFromCentroidMap [%d] : %v",i+j,distanceFromCentroidMap)
				//sorting here
				var sortedSlicePoint []kv
				for indexPointOfMap, distanceValueOfMap := range distanceFromCentroidMap {
					sortedSlicePoint = append(sortedSlicePoint, kv{indexPointOfMap, distanceValueOfMap})
				}
				sort.Slice(sortedSlicePoint, func(x, y int) bool {
					return sortedSlicePoint[x].Value < sortedSlicePoint[y].Value
				})
				for indexOfOriginalArraySlected, indexPointOfMapInPlace := range sortedSlicePoint {
					(*oNew)[i][indexOfOriginalArraySlected]=indexPointOfMapInPlace.Key
				}
				sortedSlicePoint = nil
				//fmt.Printf("Sorted array [%d] : %v",i+j,(*oNew)[i])

				CrPi := (*oNew)[i]

				//iterate over Cr
				for l:=0;l<len((*oNew)[i]);l++{
					//fmt.Printf("Value in Cr [%d] %d\n",k,(*oNew)[i][l])
					//fmt.Printf("Value in Cr [%d] %d\n",k,CrPi)
					valShifted := (*oNew)[i][l]
					CrPi=removeItem((*oNew)[i],int(valShifted))
					//fmt.Printf("After Deleting ------------->%v\n",CrPi)
					var CsPi []int64

					CsPi=append((*oNew)[j],valShifted)	
					//fmt.Printf("After Adding ------------->%v\n",CsPi)

					//Balencing condition
					if (balenceValue((*oNew)[i],(*oNew)[j], Navg) > 0 || (balenceValue((*oNew)[i],(*oNew)[j], Navg)>0 || balenceValue((*oNew)[i],(*oNew)[j], Navg)==0 && weightValue((*oNew)[i],(*oNew)[j],CrPi,CsPi,k)>0)) {
						counter = counter+1
						(*oNew)[i]=CrPi
						(*oNew)[j]=CsPi
						l=l-1
					}
					CsPi=nil
					CrPi=nil	
				}
			}
		}
	}
	return oNew
}