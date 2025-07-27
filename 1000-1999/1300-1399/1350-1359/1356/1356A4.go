package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

// applyIX applies I \otimes X to the two-qubit state vector.
func applyIX(state [4]complex128) [4]complex128 {
	// swap basis states |00> <-> |01> and |10> <-> |11>
	return [4]complex128{state[1], state[0], state[3], state[2]}
}

// applyCNOT applies a CNOT gate with qubit 0 as control and qubit 1 as target.
func applyCNOT(state [4]complex128) [4]complex128 {
	// swap |10> <-> |11>
	return [4]complex128{state[0], state[1], state[3], state[2]}
}

func main() {
	reader := bufio.NewReader(os.Stdin)
	op, _ := reader.ReadString('\n')
	op = strings.TrimSpace(op)

	// initial state |00>
	var state [4]complex128
	state[0] = 1

	switch op {
	case "IX", "I X", "I_X", "I*X", "I\\otimesX":
		state = applyIX(state)
		fmt.Println(0)
	case "CNOT":
		state = applyCNOT(state)
		fmt.Println(1)
	}

	_ = state // avoid unused variable warning
}
