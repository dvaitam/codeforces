package main

import (
	"bufio"
	"fmt"
	"os"
)

// The program distinguishes between two 2-qubit gates:
// CNOT12 (control qubit 1, target qubit 2) and
// CNOT21 (control qubit 2, target qubit 1).
// It reads a 4x4 matrix of floats representing the unitary in
// the computational basis and prints 0 if it is CNOT12 or
// 1 if it is CNOT21.
func main() {
	reader := bufio.NewReader(os.Stdin)
	var m [4][4]float64
	for i := 0; i < 4; i++ {
		for j := 0; j < 4; j++ {
			if _, err := fmt.Fscan(reader, &m[i][j]); err != nil {
				return
			}
		}
	}

	row := -1
	for i := 0; i < 4; i++ {
		if m[i][2] != 0 {
			row = i
			break
		}
	}
	if row == 3 {
		fmt.Println(0)
	} else {
		fmt.Println(1)
	}
}
