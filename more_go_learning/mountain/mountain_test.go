package main

import "testing"

func verifyMountain(expected, result int) bool {
	return expected == result
}

func TestLongestMountainArray1(t *testing.T) {
	// Test case 1 - Mountain array with a peak
	A := []int{0, 1, 0}
	expected := 3
	result := LongestMountain(A)
	if !verifyMountain(expected, result) {
		t.Errorf("Expected %v, but got %v", expected, result)
	}
}

func TestLongestMountainArray2(t *testing.T) {
	// Test case 2 - Mountain array with a peak
	A := []int{0, 2, 1, 0}
	expected := 4
	result := LongestMountain(A)
	if !verifyMountain(expected, result) {
		t.Errorf("Expected %v, but got %v", expected, result)
	}
}

func TestLongestMountainArray3(t *testing.T) {
	// Test case 3 - Mountain array with a peak
	A := []int{0, 1, 2, 3, 4, 5, 4, 3, 2, 1}
	expected := 10
	result := LongestMountain(A)
	if !verifyMountain(expected, result) {
		t.Errorf("Expected %v, but got %v", expected, result)
	}
}

func TestLongestMountainArray4(t *testing.T) {
	// Test case 4 - Mountain array with a peak
	A := []int{0, 1, 2, 3, 4, 5, 4, 3, 2, 1, 0}
	expected := 11
	result := LongestMountain(A)
	if !verifyMountain(expected, result) {
		t.Errorf("Expected %v, but got %v", expected, result)
	}
}

func TestLongestMountainArray5(t *testing.T) {
	// Test case 5 - Mountain array without a peak
	A := []int{0, 1, 2, 3, 4, 5}
	expected := 0
	result := LongestMountain(A)
	if !verifyMountain(expected, result) {
		t.Errorf("Expected %v, but got %v", expected, result)
	}
}

func TestLongestMountainArray6(t *testing.T) {
	// Test case 6 - Mountain array without a peak
	A := []int{0, 0}
	expected := 0
	result := LongestMountain(A)
	if !verifyMountain(expected, result) {
		t.Errorf("Expected %v, but got %v", expected, result)
	}
}

func TestLongestMountainArray7(t *testing.T) {
	// Test case 7 - Mountain array without a peak
	A := []int{6, 5, 4, 3, 2, 1, 0}
	expected := 0
	result := LongestMountain(A)
	if !verifyMountain(expected, result) {
		t.Errorf("Expected %v, but got %v", expected, result)
	}
}

func TestLongestMountainArray8(t *testing.T) {
	// Test case 8 - Mountain array without a peak
	A := []int{5, 5, 5, 5, 5, 5, 5, 5, 5, 5, 5, 5, 5, 5, 5, 5}
	expected := 0
	result := LongestMountain(A)
	if !verifyMountain(expected, result) {
		t.Errorf("Expected %v, but got %v", expected, result)
	}
}

func TestLongestMountainArray9(t *testing.T) {
	// Test case 9 - Mountain array without a peak
	var A []int
	expected := 0
	result := LongestMountain(A)
	if !verifyMountain(expected, result) {
		t.Errorf("Expected %v, but got %v", expected, result)
	}
}

func TestLongestMountainArray10(t *testing.T) {
	// Test case 10 - two peaks in a mountain array
	A := []int{0, 1, 2, 3, 4, 5, 4, 3, 2, 1, 0, 1, 2, 3, 4, 5, 0}
	expected := 11
	result := LongestMountain(A)
	if !verifyMountain(expected, result) {
		t.Errorf("Expected %v, but got %v", expected, result)
	}
}

func TestLongestMountainArray11(t *testing.T) {
	// Test case 11 - three peaks in a mountain array
	A := []int{0, 3, 0, 2, 3, 5, 3, 1, 4, 0}
	expected := 6
	result := LongestMountain(A)
	if !verifyMountain(expected, result) {
		t.Errorf("Expected %v, but got %v", expected, result)
	}
}
