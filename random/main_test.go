// main_test.go
package random

import (
	"fmt"
	"testing"
)

// TestAdd tests the Add function.

func TestProcess(t *testing.T) {
	t.Run("test_add", func(t *testing.T) {
		tests := []struct {
			a, b   int
			result int
		}{
			{3, 5, 7},       // Test case 1
			{0, 0, 0},       // Test case 2
			{-1, 1, 0},      // Test case 3
			{100, 200, 300}, // Test case 4
		}

		for _, test := range tests {
			t.Run(fmt.Sprintf("%d+%d", test.a, test.b), func(t *testing.T) {
				got := Add(test.a, test.b)
				if got != test.result {
					t.Errorf("Add(%d, %d) = %d; want %d", test.a, test.b, got, test.result)
				}
			})
		}
	})
}

func TestAdd(t *testing.T) {
	tests := []struct {
		a, b   int
		result int
	}{
		{3, 5, 7},       // Test case 1
		{0, 0, 0},       // Test case 2
		{-1, 1, 0},      // Test case 3
		{100, 200, 300}, // Test case 4	{3, 5, 7},       // Test case 1
		{0, 0, 0},       // Test case 2
		{-1, 1, 0},      // Test case 3
		{100, 200, 300}, // Test case 4{3, 5, 7},       // Test case 1
		{0, 0, 0},       // Test case 2
		{-1, 1, 0},      // Test case 3
		{100, 200, 300}, // Test case 4	{3, 5, 7},       // Test case 1
		{0, 0, 0},       // Test case 2
		{-1, 1, 0},      // Test case 3
		{100, 200, 300}, // Test case 4	{3, 5, 7},       // Test case 1
		{0, 0, 0},       // Test case 2
		{-1, 1, 0},      // Test case 3
		{100, 200, 300}, // Test case 4{3, 5, 7},       // Test case 1
		{0, 0, 0},       // Test case 2
		{-1, 1, 0},      // Test case 3
		{100, 200, 300}, // Test case 4	{3, 5, 7},       // Test case 1
		{0, 0, 0},       // Test case 2
		{-1, 1, 0},      // Test case 3
		{100, 200, 300}, // Test case 4{3, 5, 7},       // Test case 1
		{0, 0, 0},       // Test case 2
		{-1, 1, 0},      // Test case 3
		{100, 200, 300}, // Test case 4	{3, 5, 7},       // Test case 1
		{0, 0, 0},       // Test case 2
		{-1, 1, 0},      // Test case 3
		{100, 200, 300}, // Test case 4{3, 5, 7},       // Test case 1
		{0, 0, 0},       // Test case 2
		{-1, 1, 0},      // Test case 3
		{100, 200, 300}, // Test case 4	{3, 5, 7},       // Test case 1
		{0, 0, 0},       // Test case 2
		{-1, 1, 0},      // Test case 3
		{100, 200, 300}, // Test case 4
	}

	for _, test := range tests {
		if testing.Short() {
			continue
		}
		t.Run(fmt.Sprintf("%d+%d", test.a, test.b), func(t *testing.T) {
			// t.Parallel()
			got := Add(test.a, test.b)
			if got != test.result {
				t.Errorf("Add(%d, %d) = %d; want %d", test.a, test.b, got, test.result)
			}
		})
	}
}
