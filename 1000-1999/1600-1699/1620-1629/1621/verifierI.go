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

func lexSmaller(a []int, i1 int, b []int, i2 int, length int) bool {
	for k := 0; k < length; k++ {
		if a[i1+k] < b[i2+k] {
			return true
		} else if a[i1+k] > b[i2+k] {
			return false
		}
	}
	return false
}

func op(arr []int) []int {
	n := len(arr)
	D := make([]int, n)
	copy(D, arr)
	for i := 1; i <= n; i++ {
		bestIdx := 0
		for s := 1; s <= n-i; s++ {
			if lexSmaller(D, s, D, bestIdx, i) {
				bestIdx = s
			}
		}
		copy(D[n-i:], D[bestIdx:bestIdx+i])
	}
	return D
}

func solveCaseI(A []int, queries [][2]int) []int {
	n := len(A)
	B := make([][]int, n+1)
	B[0] = make([]int, n)
	copy(B[0], A)
	for i := 1; i <= n; i++ {
		B[i] = op(B[i-1])
	}
	res := make([]int, len(queries))
	for idx, q := range queries {
		i, j := q[0], q[1]
		if i > n {
			i = n
		}
		if j > n {
			j = n
		}
		res[idx] = B[i][j-1]
	}
	return res
}

func runCandidate(bin, input string) (string, error) {
	var cmd *exec.Cmd
	if strings.HasSuffix(bin, ".go") {
		cmd = exec.Command("go", "run", bin)
	} else {
		cmd = exec.Command(bin)
	}
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("runtime error: %v\n%s", err, out.String())
	}
	return strings.TrimSpace(out.String()), nil
}

func generateCaseI(rng *rand.Rand) (string, []string) {
	n := rng.Intn(5) + 1
	A := make([]int, n)
	for i := range A {
		A[i] = rng.Intn(10) + 1
	}
	qn := rng.Intn(4) + 1
	queries := make([][2]int, qn)
	input := fmt.Sprintf("1\n%d\n", n)
	for i, v := range A {
		if i > 0 {
			input += " "
		}
		input += fmt.Sprintf("%d", v)
	}
	input += "\n"
	input += fmt.Sprintf("%d\n", qn)
	for i := 0; i < qn; i++ {
		x := rng.Intn(n+2) + 1
		y := rng.Intn(n+2) + 1
		queries[i] = [2]int{x, y}
		input += fmt.Sprintf("%d %d\n", x, y)
	}
	resVals := solveCaseI(A, queries)
	outLines := make([]string, qn)
	for i, v := range resVals {
		outLines[i] = fmt.Sprintf("%d", v)
	}
	return input, outLines
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierI.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < 100; i++ {
		in, expLines := generateCaseI(rng)
		got, err := runCandidate(bin, in)
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:\n%s", i+1, err, in)
			os.Exit(1)
		}
		gotLines := strings.Split(strings.TrimSpace(got), "\n")
		if len(gotLines) != len(expLines) {
			fmt.Fprintf(os.Stderr, "case %d failed: expected %d lines got %d\ninput:\n%s", i+1, len(expLines), len(gotLines), in)
			os.Exit(1)
		}
		for j := range expLines {
			if strings.TrimSpace(gotLines[j]) != expLines[j] {
				fmt.Fprintf(os.Stderr, "case %d failed: expected %s got %s\ninput:\n%s", i+1, expLines[j], gotLines[j], in)
				os.Exit(1)
			}
		}
	}
	fmt.Println("All tests passed")
}
