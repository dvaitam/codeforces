package main

import (
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"strings"
	"time"
)

type testCaseI struct {
	n     int
	edges [][2]int
}

func generateTestsI(num int) []testCaseI {
	rand.Seed(time.Now().UnixNano())
	tests := make([]testCaseI, num)
	for i := 0; i < num; i++ {
		n := rand.Intn(10) + 2
		edges := make([][2]int, n-1)
		for j := 0; j < n-1; j++ {
			u := j + 2
			v := rand.Intn(u-1) + 1
			edges[j] = [2]int{u, v}
		}
		tests[i] = testCaseI{n: n, edges: edges}
	}
	return tests
}

func solveI(tc testCaseI) string {
	n := tc.n
	deg := make([]int, n+1)
	for _, e := range tc.edges {
		u := e[0]
		v := e[1]
		deg[u]++
		deg[v]++
	}
	leaves := 0
	for i := 1; i <= n; i++ {
		if deg[i] == 1 {
			leaves++
		}
	}
	ans := n - 1
	if leaves*(leaves-1)/2 > ans {
		ans = leaves * (leaves - 1) / 2
	}
	return fmt.Sprint(ans)
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("usage: go run verifierI.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	tests := generateTestsI(100)
	var input bytes.Buffer
	for _, tc := range tests {
		fmt.Fprintln(&input, tc.n)
		for _, e := range tc.edges {
			fmt.Fprintf(&input, "%d %d\n", e[0], e[1])
		}
	}
	cmd := exec.Command(bin)
	cmd.Stdin = bytes.NewReader(input.Bytes())
	out, err := cmd.CombinedOutput()
	if err != nil {
		fmt.Println("binary execution failed:", err)
		os.Exit(1)
	}
	outputs := strings.Split(strings.TrimSpace(string(out)), "\n")
	if len(outputs) != len(tests) {
		fmt.Printf("expected %d lines of output, got %d\n", len(tests), len(outputs))
		os.Exit(1)
	}
	for i, tc := range tests {
		expected := solveI(tc)
		if strings.TrimSpace(outputs[i]) != expected {
			fmt.Printf("mismatch on test %d: expected %s got %s\n", i+1, expected, outputs[i])
			os.Exit(1)
		}
	}
	fmt.Println("all tests passed")
}
