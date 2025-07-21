package main

import (
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"time"
)

func runCandidate(bin, input string) (string, error) {
	var cmd *exec.Cmd
	if strings.HasSuffix(bin, ".go") {
		cmd = exec.Command("go", "run", bin)
	} else {
		cmd = exec.Command(bin)
	}
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	var stderr bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &stderr
	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("runtime error: %v\n%s", err, stderr.String())
	}
	return strings.TrimSpace(out.String()), nil
}

func solveCase(n, m int, A, B [][]int) string {
	size := n * m
	rowA := make([]int, size+1)
	colA := make([]int, size+1)
	rowB := make([]int, size+1)
	colB := make([]int, size+1)
	for i := 0; i < n; i++ {
		for j := 0; j < m; j++ {
			x := A[i][j]
			rowA[x] = i
			colA[x] = j
		}
	}
	for i := 0; i < n; i++ {
		for j := 0; j < m; j++ {
			x := B[i][j]
			rowB[x] = i
			colB[x] = j
		}
	}
	rowMap := make([]int, n)
	colMap := make([]int, m)
	for i := range rowMap {
		rowMap[i] = -1
	}
	for i := range colMap {
		colMap[i] = -1
	}
	ok := true
	for val := 1; val <= size && ok; val++ {
		rA, cA := rowA[val], colA[val]
		rB, cB := rowB[val], colB[val]
		if rowMap[rA] == -1 {
			rowMap[rA] = rB
		} else if rowMap[rA] != rB {
			ok = false
			break
		}
		if colMap[cA] == -1 {
			colMap[cA] = cB
		} else if colMap[cA] != cB {
			ok = false
			break
		}
	}
	usedRow := make([]bool, n)
	usedCol := make([]bool, m)
	for _, v := range rowMap {
		if v == -1 || usedRow[v] {
			ok = false
			break
		}
		usedRow[v] = true
	}
	for _, v := range colMap {
		if v == -1 || usedCol[v] {
			ok = false
			break
		}
		usedCol[v] = true
	}
	if ok {
		return "YES"
	}
	return "NO"
}

func perm(rng *rand.Rand, n int) []int {
	arr := make([]int, n)
	for i := 0; i < n; i++ {
		arr[i] = i + 1
	}
	rng.Shuffle(n, func(i, j int) { arr[i], arr[j] = arr[j], arr[i] })
	return arr
}

func generateCase(rng *rand.Rand) (string, string) {
	n := rng.Intn(3) + 1
	m := rng.Intn(3) + 1
	size := n * m
	p1 := perm(rng, size)
	p2 := perm(rng, size)
	A := make([][]int, n)
	B := make([][]int, n)
	idx := 0
	for i := 0; i < n; i++ {
		A[i] = make([]int, m)
		for j := 0; j < m; j++ {
			A[i][j] = p1[idx]
			idx++
		}
	}
	idx = 0
	for i := 0; i < n; i++ {
		B[i] = make([]int, m)
		for j := 0; j < m; j++ {
			B[i][j] = p2[idx]
			idx++
		}
	}
	var sb strings.Builder
	sb.WriteString("1\n")
	sb.WriteString(fmt.Sprintf("%d %d\n", n, m))
	for i := 0; i < n; i++ {
		for j := 0; j < m; j++ {
			if j > 0 {
				sb.WriteByte(' ')
			}
			sb.WriteString(strconv.Itoa(A[i][j]))
		}
		sb.WriteByte('\n')
	}
	for i := 0; i < n; i++ {
		for j := 0; j < m; j++ {
			if j > 0 {
				sb.WriteByte(' ')
			}
			sb.WriteString(strconv.Itoa(B[i][j]))
		}
		sb.WriteByte('\n')
	}
	input := sb.String()
	expect := solveCase(n, m, A, B)
	return input, expect
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierE.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < 100; i++ {
		in, exp := generateCase(rng)
		out, err := runCandidate(bin, in)
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:\n%s", i+1, err, in)
			os.Exit(1)
		}
		if out != exp {
			fmt.Fprintf(os.Stderr, "case %d failed: expected %s got %s\ninput:\n%s", i+1, exp, out, in)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
