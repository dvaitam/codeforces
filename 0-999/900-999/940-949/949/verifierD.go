package main

import (
	"bytes"
	"context"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"time"
)

type test struct {
	n, d, b int
	a       []int
}

func genTests() []test {
	rand.Seed(4)
	tests := make([]test, 100)
	for i := range tests {
		n := rand.Intn(20) + 2
		d := rand.Intn(n-1) + 1
		b := rand.Intn(10) + 1
		a := make([]int, n)
		sum := n * b
		for j := 0; j < n-1; j++ {
			val := rand.Intn(2*b + 1)
			if val > sum {
				val = sum
			}
			a[j] = val
			sum -= val
		}
		a[n-1] = sum
		tests[i] = test{n, d, b, a}
	}
	return tests
}

func runBinary(path, input string) (string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	var cmd *exec.Cmd
	if strings.HasSuffix(path, ".go") {
		cmd = exec.CommandContext(ctx, "go", "run", path)
	} else {
		cmd = exec.CommandContext(ctx, path)
	}
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	err := cmd.Run()
	return out.String(), err
}

func compute(arr []int, rooms, d, b int) int {
	n := len(arr)
	p := 0
	var avail int64
	miss := 0
	for i := 1; i <= rooms; i++ {
		limit := i * (d + 1)
		if limit > n {
			limit = n
		}
		for p < limit {
			avail += int64(arr[p])
			p++
		}
		if avail >= int64(b) {
			avail -= int64(b)
		} else {
			miss++
		}
	}
	return miss
}

func solveRef(t test) int {
	a := make([]int, len(t.a))
	copy(a, t.a)
	leftRooms := (t.n + 1) / 2
	rightRooms := t.n / 2
	left := compute(a, leftRooms, t.d, t.b)
	rev := make([]int, len(a))
	for i := 0; i < len(a); i++ {
		rev[i] = a[len(a)-1-i]
	}
	right := compute(rev, rightRooms, t.d, t.b)
	if left > right {
		return left
	}
	return right
}

func buildInput(t test) string {
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("%d %d %d\n", t.n, t.d, t.b))
	for i, v := range t.a {
		if i > 0 {
			sb.WriteByte(' ')
		}
		sb.WriteString(strconv.Itoa(v))
	}
	sb.WriteByte('\n')
	return sb.String()
}

func parseOutput(out string) (int, bool) {
	s := strings.TrimSpace(out)
	v, err := strconv.Atoi(s)
	if err != nil {
		return 0, false
	}
	return v, true
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("usage: go run verifierD.go /path/to/binary")
		os.Exit(1)
	}
	candidate := os.Args[1]
	tests := genTests()
	for i, t := range tests {
		input := buildInput(t)
		out, err := runBinary(candidate, input)
		if err != nil {
			fmt.Printf("test %d: runtime error: %v\n", i+1, err)
			os.Exit(1)
		}
		got, ok := parseOutput(out)
		if !ok {
			fmt.Printf("test %d: bad output\n", i+1)
			os.Exit(1)
		}
		exp := solveRef(t)
		if got != exp {
			fmt.Printf("test %d failed\ninput:\n%s\nexpected:%d got:%d\n", i+1, input, exp, got)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
