package skyline

func getSkyline(buildings [][]int) [][]int {
	/*
		Get the skyline of a city represented by a list of buildings

		Inputs: buildings [][]int - a list of buildings represented by [left, right, height]
		outputs: [][]int - a list of points representing the skyline

		Divide, Conquer, and Merge approach:
		1. If the list of buildings is empty, return an empty list
		2. If the list of buildings has only one building, return the skyline of that building
			2.1. The skyline of a building is represented by the points [left, height] and [right, 0]
		3. Otherwise, divide the list of buildings into two halves
		4. Recursively get the skyline of the left half and the right half
		5. Merge the two skylines
		6. Cleanup any redundant points in the merged skyline and return the result

		To merge two skylines:
		1. Initialize two pointers, i and j, to 0
		2. Initialize two variables, h1 and h2, to 0
		3. Initialize an empty list, skyline, to store the merged skyline
		4. While i < len(skyline1) and j < len(skyline2):
			4.1. If the x-coordinate of the point in skyline1 is less than the x-coordinate of the point in skyline2:
				4.1.1. Update h1 to the y-coordinate of the point in skyline1
				4.1.2. Calculate the maximum height, h, as the maximum of h1 and h2
				4.1.3. Append the point [x-coordinate of the point in skyline1, h] to the merged skyline
				4.1.4. Increment i by 1
			4.2. If the x-coordinate of the point in skyline1 is greater than the x-coordinate of the point in skyline2:
				4.2.1. Update h2 to the y-coordinate of the point in skyline2
				4.2.2. Calculate the maximum height, h, as the maximum of h1 and h2
				4.2.3. Append the point [x-coordinate of the point in skyline2, h] to the merged skyline
				4.2.4. Increment j by 1
			4.3. If the x-coordinate of the point in skyline1 is equal to the x-coordinate of the point in skyline2:
				4.3.1. Update h1 and h2 to the y-coordinates of the points in skyline1 and skyline2
				4.3.2. Calculate the maximum height, h, as the maximum of h1 and h2
				4.3.3. Append the point [x-coordinate of the point in skyline1, h] to the merged skyline
				4.3.4. Increment i and j by 1
		5. Append the remaining points in skyline1 and skyline2 to the merged skyline
		6. Return the merged skyline
	*/

	results := getSkylineBeforeCleaning(buildings)
	// remove redundant points
	var skyline [][]int
	for i := 0; i < len(results); i++ {
		if i == 0 || results[i][1] != results[i-1][1] {
			skyline = append(skyline, results[i])
		}
	}
	return skyline
}

func getSkylineBeforeCleaning(buildings [][]int) [][]int {
	if len(buildings) == 0 {
		return [][]int{}
	}
	if len(buildings) == 1 {
		return [][]int{{buildings[0][0], buildings[0][2]}, {buildings[0][1], 0}} // skyline of a building
	}
	mid := len(buildings) / 2
	leftSkyline := getSkyline(buildings[:mid])
	rightSkyline := getSkyline(buildings[mid:])

	return mergeSkyline(leftSkyline, rightSkyline)
}

func mergeSkyline(skyline1, skyline2 [][]int) [][]int {
	i, j := 0, 0
	h1, h2 := 0, 0
	var skyline [][]int
	for i < len(skyline1) && j < len(skyline2) { // for each point in skyline1 and skyline2
		if skyline1[i][0] < skyline2[j][0] { // if the x-coordinate of the point in skyline1 is less than the x-coordinate of the point in skyline2
			h1 = skyline1[i][1]                                 // update h1 to the y-coordinate of the point in skyline1
			h := max(h1, h2)                                    // calculate the maximum height, h, as the maximum of h1 and h2
			skyline = append(skyline, []int{skyline1[i][0], h}) // append the point [x-coordinate of the point in skyline1, h] to the merged skyline
			i++                                                 // increment i by 1
		} else if skyline1[i][0] > skyline2[j][0] { // if the x-coordinate of the point in skyline1 is greater than the x-coordinate of the point in skyline2
			h2 = skyline2[j][1]                                 // update h2 to the y-coordinate of the point in skyline2
			h := max(h1, h2)                                    // calculate the maximum height, h, as the maximum of h1 and h2
			skyline = append(skyline, []int{skyline2[j][0], h}) // append the point [x-coordinate of the point in skyline2, h] to the merged skyline
			j++                                                 // increment j by 1
		} else {
			h1 = skyline1[i][1]                                 // update h1 to the y-coordinate of the point in skyline1
			h2 = skyline2[j][1]                                 // update h2 to the y-coordinate of the point in skyline2
			h := max(h1, h2)                                    // calculate the maximum height, h, as the maximum of h1 and h2
			skyline = append(skyline, []int{skyline1[i][0], h}) // append the point [x-coordinate of the point in skyline1, h] to the merged skyline
			i++                                                 // increment i by 1
			j++                                                 // increment j by 1
		}
	}
	skyline = append(skyline, skyline1[i:]...) // append the remaining points in skyline1 to the merged skyline
	skyline = append(skyline, skyline2[j:]...) // append the remaining points in skyline2 to the merged skyline
	return skyline
}
