package main

func LongestMountain(A []int) int {
	/*
		Find the length of the longest mountain in A.

		Inputs: A []int - a mountain array.
		Outputs: int - the length of the longest mountain in A.

		Example:
		A = [0, 1, 0]
		LongestMountain(A) = 3

		A = [0, 2, 1, 0]
		LongestMountain(A) = 4

		Steps:
		1. Initialize a variable to store the length of the longest mountain.
		2. Iterate through the array.
		3. Check if the current element is a peak.
		4. If it is a peak, expand to the left and right to find the length of the mountain.
		5. Update the length of the longest mountain if the current mountain is longer.
		6. Return the length of the longest mountain.
	*/

	// Initialize the length of the longest mountain.
	longestMountain := 0

	// Iterate through the array.
	for index, value := range A {
		// check that we are not on the edges of the array
		if index > 0 && index < len(A)-1 {
			// check if the current element is a peak
			if A[index-1] < value && value > A[index+1] {
				// expand to the left and right to find the length of the mountain
				left := index - 1
				right := index + 1
				for left > 0 && A[left-1] < A[left] {
					left--
				}
				for right < len(A)-1 && A[right] > A[right+1] {
					right++
				}
				// update the length of the longest mountain if the current mountain is longer
				if right-left+1 > longestMountain {
					longestMountain = right - left + 1
				}
			}
		}
	}
	return longestMountain
}
