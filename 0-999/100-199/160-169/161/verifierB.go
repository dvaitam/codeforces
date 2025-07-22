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

type Test struct {
	n     int
	votes []int
}

func generateTest() Test {
	n := rand.Intn(20) + 1
	votes := make([]int, n)
	for i := 0; i < n; i++ {
		votes[i] = rand.Intn(10)
	}
	return Test{n, votes}
}

func (t Test) Input() string {
	var sb strings.Builder
	fmt.Fprintf(&sb, "%d\n", t.n)
	for i, v := range t.votes {
		if i > 0 {
			sb.WriteByte(' ')
		}
		fmt.Fprintf(&sb, "%d", v)
	}
	sb.WriteByte('\n')
	return sb.String()
}

func solve(t Test) string {
	sum := 0
	for _, v := range t.votes {
		sum += v
	}
	if t.votes[0] > sum-t.votes[0] {
		return "Yes"
	}
	return "No"
}

func runBinary(path, input string) (string, error) {
	cmd := exec.Command(path)
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = os.Stderr
	err := cmd.Run()
	return strings.TrimSpace(out.String()), err
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierB.go <binary>")
		os.Exit(1)
	}
	bin := os.Args[1]
	rand.Seed(time.Now().UnixNano())
	for i := 0; i < 100; i++ {
		t := generateTest()
		inp := t.Input()
		exp := solve(t)
		out, err := runBinary(bin, inp)
		if err != nil {
			fmt.Fprintf(os.Stderr, "test %d runtime error: %v\n", i+1, err)
			os.Exit(1)
		}
		if strings.TrimSpace(out) != exp {
			fmt.Fprintf(os.Stderr, "test %d failed: expected %s got %s\ninput:\n%s\n", i+1, exp, out, inp)
			os.Exit(1)
		}
	}
	fmt.Println("all tests passed")
}
