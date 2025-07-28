package main

import (
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
)

type edge struct{ u, v int }

func genEdges(n int) [][2]int {
	N := 1 << n
	edges := make([][2]int, 0, n*N/2)
	for i := 0; i < N; i++ {
		for b := 0; b < n; b++ {
			j := i ^ (1 << b)
			if i < j {
				edges = append(edges, [2]int{i, j})
			}
		}
	}
	return edges
}

func permEdges(edges [][2]int, perm []int) [][2]int {
	res := make([][2]int, len(edges))
	for i, e := range edges {
		res[i] = [2]int{perm[e[0]], perm[e[1]]}
	}
	return res
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierE.go /path/to/binary")
		os.Exit(1)
	}
	binary := os.Args[1]
	rand.Seed(42)
	const T = 100
	type tc struct {
		n     int
		edges [][2]int
	}
	tests := make([]tc, T)
	for i := 0; i < T; i++ {
		n := rand.Intn(3) + 1
		base := genEdges(n)
		perm := rand.Perm(1 << n)
		pe := permEdges(base, perm)
		tests[i] = tc{n: n, edges: pe}
	}
	var input bytes.Buffer
	fmt.Fprintln(&input, T)
	for _, tc := range tests {
		fmt.Fprintln(&input, tc.n)
		for _, e := range tc.edges {
			fmt.Fprintf(&input, "%d %d\n", e[0], e[1])
		}
	}
	inpBytes := input.Bytes()
	// run official solver
	solCmd := exec.Command("go", "run", "1543E.go")
	solCmd.Stdin = bytes.NewReader(inpBytes)
	var solOut bytes.Buffer
	solCmd.Stdout = &solOut
	solCmd.Stderr = os.Stderr
	if err := solCmd.Run(); err != nil {
		fmt.Fprintln(os.Stderr, "failed to run official solver:", err)
		os.Exit(1)
	}
	// run candidate
	cmd := exec.Command(binary)
	cmd.Stdin = bytes.NewReader(inpBytes)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		fmt.Fprintln(os.Stderr, "error running binary:", err)
		os.Exit(1)
	}
	if !bytes.Equal(bytes.TrimSpace(out.Bytes()), bytes.TrimSpace(solOut.Bytes())) {
		fmt.Fprintln(os.Stderr, "output mismatch with official solution")
		os.Exit(1)
	}
	fmt.Println("all tests passed")
}
