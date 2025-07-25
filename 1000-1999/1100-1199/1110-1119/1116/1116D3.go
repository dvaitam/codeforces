package main

import (
	"fmt"
	"math"
)

// main prints the 8x8 unitary matrix for N=3 qubits
// with non-zero elements on both main diagonal and anti-diagonal.
// Non-zero entries are set to 1/sqrt(2) to ensure unitarity.
func main() {
	const N = 3
	dim := 1 << N
	// Value for non-zero entries: 1/sqrt(2)
	val := 1.0 / math.Sqrt(2.0)
	for i := 0; i < dim; i++ {
		for j := 0; j < dim; j++ {
			var x float64
			if i == j || i+j == dim-1 {
				x = val
			} else {
				x = 0.0
			}
			if j > 0 {
				fmt.Print(" ")
			}
			fmt.Printf("%f", x)
		}
		fmt.Println()
	}
}
