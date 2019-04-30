package sorter

import (
	"log"
	"testing"
	"sort"
	"math/rand"
)

type Sorter struct {}

// INT SORTING

func (s *Sorter) SortIntSlice(slice []int, algorithm string) (sliceSorted []int){

	switch algorithm {
	case "insertion":
		sliceSorted = SortUsingInsertionAlgo(slice)
	case "quicksort":
		sliceSorted = SortUsingQuicksortAlgo(slice)
	}
	return
}

func (s *Sorter) InverseSortIntSlice(slice []int, algorithm string) (sliceSorted []int){

	switch algorithm {
	case "insertion":
		sliceSorted = SortUsingInverseInsertionAlgo(slice)
	case "quicksort":
		// sliceSorted = SortUsingInverseQuicksortAlgo(slice)
	}
	return

}

func SortUsingInverseInsertionAlgo(slice []int) (sliceSorted []int){

	originalLength := len(slice)
	for i := 0; i < originalLength; i++ {
		maximumValue, positionOfMaximum := GetMaximum(slice)
		sliceSorted = append(sliceSorted, maximumValue)
		slice = DeleteIntOfSlice(slice, positionOfMaximum)
	}
	return
}

func SortUsingInverseQuicksortAlgo(a []int) []int {
	// TODO quicksort inverse algo undone, exists?
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

func SortUsingInsertionAlgo(slice []int) (sliceSorted []int){

	originalLength := len(slice)
	for i := 0; i < originalLength; i++ {
		minimumValue, positionOfMinimum := GetMinimum(slice)
		sliceSorted = append(sliceSorted, minimumValue)
		slice = DeleteIntOfSlice(slice, positionOfMinimum)
	}
	return
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

// STRING SORTING

func (s *Sorter) SortStringSlice(slice []string) (sliceSorted []string){

	originalLength := len(slice)
	_ = originalLength

	for i := 0; i < originalLength; i++ {
		firstAlphabetString, position := GetFirstAlphabetically(slice)
		sliceSorted = append(sliceSorted, firstAlphabetString)
		slice = DeleteStringOfSlice(slice, position)
	}

	return
}

// HELP FUNC

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

func GetMaximum(slice []int) (maximumValue int, positionOfMaximum int){

	maximumValue = slice[len(slice)-1]

	for i := len(slice)-1; i >= 0; i-- {
		if maximumValue <= slice[i]{
			maximumValue = slice[i]
			positionOfMaximum = i
		}
	}
	return
}

func GetFirstAlphabetically(slice []string) (firstAlphabetString string, position int){

	firstAlphabetString = slice[len(slice)-1]

	for i := 0; i < len(slice); i++ {
		if firstAlphabetString >= slice[i] {
			firstAlphabetString = slice[i]
			position = i
		}
	}
	return
}

func AddValueToSlice(slice []int, value int){
	slice = append(slice, value)
}

func DeleteIntOfSlice(slice []int, positionToDelete int)(sliceCleaned []int){
	slice[positionToDelete] = slice[len(slice)-1] // Save last value into the must be deleted value position.
	sliceCleaned = slice[:len(slice)-1] // Truncate slice forgetting last position
	return
}

func DeleteStringOfSlice(slice []string, positionToDelete int)(sliceCleaned []string){
	slice[positionToDelete] = slice[len(slice)-1]
	sliceCleaned = slice[:len(slice)-1]
	return 
}

// BENCHMARKS

func (s *Sorter)RunBenchmarks(){
	br := testing.Benchmark(BenchmarkSortInsertionAlgo)
	log.Print(br)

	br2 := testing.Benchmark(BenchmarkSortNativeGO)
	log.Print(br2)

	br3 := testing.Benchmark(BenchmarkQuickSorted)
	log.Print(br3)
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

