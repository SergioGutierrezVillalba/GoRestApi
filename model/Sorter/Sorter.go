package sorter

import (
	// "time"
	"log"
	"testing"
	"sort"
	"math/rand"
)

type Sorter struct {}

func (s *Sorter) SortIntSlice(slice []int, algorithm string) (sliceSorted []int){

	switch algorithm {
	case "insertion":
		sliceSorted = SortUsingInsertionAlgo(slice)
	case "quicksort":
		sliceSorted = SortUsingQuicksortAlgo(slice)
	}
	return
}

func (s *Sorter)RunBenchmarks(){
	br := testing.Benchmark(BenchmarkSortInsertionAlgo)
	log.Print(br)

	br2 := testing.Benchmark(BenchmarkSortNativeGO)
	log.Print(br2)

	br3 := testing.Benchmark(BenchmarkQuickSorted)
	log.Print(br3)
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
		SortUsingQuicksortAlgo(slice)
	}
}

func SortUsingQuicksortAlgo(a []int) []int {
    if len(a) < 2 {
        return a
    }
      
    left, right := 0, len(a)-1
      
    pivot := rand.Int() % len(a)
      
    a[pivot], a[right] = a[right], a[pivot]
      
    for i, _ := range a {
        if a[i] < a[right] {
            a[left], a[i] = a[i], a[left]
            left++
        }
    }
      
    a[left], a[right] = a[right], a[left]
      
    SortUsingQuicksortAlgo(a[:left])
    SortUsingQuicksortAlgo(a[left+1:])
      
    return a
}