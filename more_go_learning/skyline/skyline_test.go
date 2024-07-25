package skyline

import (
	"testing"
)

func verifySkyline(expected, result [][]int) bool {
	if len(result) != len(expected) {
		return false
	}
	for i := 0; i < len(result); i++ {
		if result[i][0] != expected[i][0] || result[i][1] != expected[i][1] {
			return false
		}
	}
	return true
}

func TestGetSkyline1(t *testing.T) {
	// Test case 1 - Multiple buildings with different heights
	buildings := [][]int{{2, 9, 10}, {3, 7, 15}, {5, 12, 12}, {15, 20, 10}, {19, 24, 8}}
	expected := [][]int{{2, 10}, {3, 15}, {7, 12}, {12, 0}, {15, 10}, {20, 8}, {24, 0}}
	result := getSkyline(buildings)
	if !verifySkyline(expected, result) {
		t.Errorf("Expected %v, but got %v", expected, result)
	}
}

func TestGetSkyline2(t *testing.T) {
	// Test case 2 - Multiple buildings with the same height
	buildings := [][]int{{0, 2, 3}, {2, 5, 3}}
	expected := [][]int{{0, 3}, {5, 0}}
	result := getSkyline(buildings)
	if !verifySkyline(expected, result) {
		t.Errorf("Expected %v, but got %v", expected, result)
	}
}

func TestGetSkyline3(t *testing.T) {
	// Test case 3 - Multiple buildings with same start and end points
	buildings := [][]int{{1, 2, 1}, {1, 2, 2}, {1, 2, 3}}
	expected := [][]int{{1, 3}, {2, 0}}
	result := getSkyline(buildings)
	if !verifySkyline(expected, result) {
		t.Errorf("Expected %v, but got %v", expected, result)
	}
}

func TestGetSkyline4(t *testing.T) {
	// Test case 4 - Empty building
	buildings := [][]int{{0, 0, 0}, {0, 0, 0}, {0, 0, 0}}
	expected := [][]int{{0, 0}}
	result := getSkyline(buildings)
	if !verifySkyline(expected, result) {
		t.Errorf("Expected %v, but got %v", expected, result)
	}
}

func TestGetSkyline5(t *testing.T) {
	//test case 5 - Zero height building
	buildings := [][]int{{0, 1, 0}, {1, 2, 0}, {2, 3, 0}}
	expected := [][]int{{0, 0}}
	result := getSkyline(buildings)
	if !verifySkyline(expected, result) {
		t.Errorf("Expected %v, but got %v", expected, result)
	}
}

func TestGetSkyline6(t *testing.T) {
	//test case 6 - Multiple buildings with different heights but overlapping x-coordinates
	buildings := [][]int{{0, 1, 4}, {0, 1, 2}, {0, 1, 3}}
	expected := [][]int{{0, 4}, {1, 0}}
	result := getSkyline(buildings)
	if !verifySkyline(expected, result) {
		t.Errorf("Expected %v, but got %v", expected, result)
	}
}
