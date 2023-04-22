package main

import (
	"fmt"
	"log"
	"sort"
	"sync"
)

func main() {

	fmt.Printf("Enter size of your array: ")
	var size int
	fmt.Scanln(&size)
	if size < 4 {
		log.Fatal("Size can not be less then 4!")
	}
	var arr = make([]int, size)
	var wg sync.WaitGroup
	var subsize int
	subsize = size / 4

	for i := 0; i < size; i++ {
		fmt.Printf("Enter %dth element: ", i)
		fmt.Scanln(&arr[i])
	}
	fmt.Println("Your array is: ", arr)

	//split array to 4 parts
	subarr1 := arr[0 : subsize*1]
	subarr2 := arr[subsize*1 : subsize*2]
	subarr3 := arr[subsize*2 : subsize*3]
	subarr4 := arr[subsize*3:]
	//sort each of 4 arrays
	wg.Add(4)
	go sorting(&wg, subarr1, 1)
	go sorting(&wg, subarr2, 2)
	go sorting(&wg, subarr3, 3)
	go sorting(&wg, subarr4, 4)
	wg.Wait()
	//merge into main array
	arr = subarr1
	arr = append(arr, subarr2...)
	arr = append(arr, subarr3...)
	arr = append(arr, subarr4...)
	sort.Ints(arr)
	//printout main sorted array
	fmt.Println("Sorted main array is:", arr)
}

//sort array of integers
func sorting(wg *sync.WaitGroup, arr []int, id int) {
	defer wg.Done()
	fmt.Println("subarray", id, "is:", arr)
	sort.Ints(arr)
}
