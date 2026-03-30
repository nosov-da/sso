package main

import "fmt"

// BinarySearch performs binary search on a sorted slice
func BinarySearch(arr []int, target int) int {
	left, right := 0, len(arr)-1

	for left <= right {
		mid := left + (right-left)/2

		if arr[mid] == target {
			return mid
		} else if arr[mid] < target {
			left = mid + 1
		} else {
			right = mid - 1
		}
	}

	return -1 // target not found
}

func main() {
	// Example usage
	arr := []int{2, 5, 8, 12, 16, 23, 38, 45, 56, 67, 78}
	target := 23

	result := BinarySearch(arr, target)

	if result != -1 {
		fmt.Printf("Element found at index: %d\n", result)
	} else {
		fmt.Println("Element not found")
	}
}
