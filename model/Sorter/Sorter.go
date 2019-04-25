package sorter

import (
	// "time"
	"log"
	"testing"
	"sort"
)

type Sorter struct {}

func (s *Sorter) SortIntSlice(slice []int, algorithm string) (sliceSorted []int){

	br := testing.Benchmark(BenchmarkSortInsertionAlgo)
	log.Print(br)

	br2 := testing.Benchmark(BenchmarkSortNativeGO)
	log.Print(br2)

	br3 := testing.Benchmark(BenchmarkQuickSorted)
	log.Print(br3)
	
	// switch algorithm {
	// case "insertion":
	// 	return SortUsingInsertionAlgo(slice)
	// case "quicksort":
	// 	return SortUsingQuicksortAlgo(slice) 
	// }
	return slice
}

func SortUsingInsertionAlgo(slice []int) (sliceSorted []int){

	originalLength := len(slice)
	for i := 0; i < originalLength; i++ {
		minimumValue, positionOfMinimum := GetMinimum(slice)
		sliceSorted = append(sliceSorted, minimumValue)
		slice = DeleteOfSlice(slice, positionOfMinimum)
		// AddValueToSlice(sliceSorted, minimumValue)
	}
	return
}

func GetMinimum(slice []int) (minimumValue int, positionOfMinimum int){

	minimumValue = slice[0]

	for i := 0; i < len(slice); i++ {
		if minimumValue > slice[i]{
			minimumValue = slice[i]
			positionOfMinimum = i
		}
	}
	return
}

func AddValueToSlice(slice []int, value int){
	slice = append(slice, value)
}

func DeleteOfSlice(slice []int, positionToDelete int)(sliceCleaned []int){
	slice[positionToDelete] = slice[len(slice)-1] // Save last value into the must be deleted value position.
	sliceCleaned = slice[:len(slice)-1] // Truncate slice forgetting last position
	return
}

func BenchmarkSortInsertionAlgo(b *testing.B){
	slice := []int{2,5,9,2,3,7,9,10,34}
	for n := 0; n < b.N; n++ {
		SortUsingInsertionAlgo(slice)
	}
}

func BenchmarkSortNativeGO(b *testing.B){
	slice := []int{2,5,9,2,3,7,9,10,34}
	for n := 0; n < b.N; n++ {
		sort.Ints(slice)
	}
}

func BenchmarkQuickSorted(b *testing.B){
	slice := []int{2,5,9,2,3,7,9,10,34}
	for n := 0; n < b.N; n++ {
		Quicksort(slice)
	}
}

func Quicksort(slice []int){
	if len(slice) > 2 {
		pivote := (slice[0] + slice[len(slice)-1] + slice[len(slice)-1/2])/3  // 1/2??
	}
}

// func TraverseAllValuesAfterOnePosition(sliceToAdapt []int, positionToFree int)(sliceAdapted []int){
// 	while len(sliceToAdapt) {

// 	}
// }

func SortUsingQuicksortAlgo(slice []int) (sliceSorted []int){
	return slice
}
