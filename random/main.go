package random

import (
	"fmt"
)

// Add is a simple function that adds two integers.
func Add(a, b int) int {
	return a + b
}
func Sub(a, b int) int {
	return a - b
}

func main() {
	result := Add(3, 5)
	fmt.Println("Result of Add function:", result)
}
