package main

// Solution for problemC1: distinguish between |0> and |+> qubit states with >=80% success
// Using the optimal Helstrom measurement: apply RY(pi/4) rotation, then measure in Z basis.
// The operations are printed as a sequence of gates for the quantum interactor.
// It is assumed the interactor supports RY and MEASURE commands, and interprets them accordingly.
import (
	"fmt"
)

func main() {
	// Two operations: rotation and measurement
	fmt.Println(2)
	// Rotate qubit 1 by pi/4 around Y axis (parameter is fraction of pi)
	fmt.Println("RY 1 0.25")
	// Measure qubit 1 in computational basis
	fmt.Println("MEASURE 1")
}
