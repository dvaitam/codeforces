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

type testB struct {
	n, m, h int
	s       []int
}

func solveB(n, m, h int, s []int) string {
	total := 0
	for _, x := range s {
		total += x
	}
	if total < n {
		return "-1\n"
	}
	A := total - s[h-1]
	B := n - 1
	probNone := 1.0
	if A < B {
		probNone = 0.0
	} else {
		for i := 0; i < B; i++ {
			probNone *= float64(A-i) / float64(total-1-i)
		}
	}
	P := 1.0 - probNone
	return fmt.Sprintf("%.10f\n", P)
}

func run(bin, input string) (string, error) {
	cmd := exec.Command(bin)
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	var stderr bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &stderr
	err := cmd.Run()
	if err != nil {
		return out.String() + stderr.String(), err
	}
	return out.String(), nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierB.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < 100; i++ {
		m := rng.Intn(6) + 1
		n := rng.Intn(5) + 1
		h := rng.Intn(m) + 1
		s := make([]int, m)
		for j := 0; j < m; j++ {
			s[j] = rng.Intn(5) + 1
		}
		// occasionally make impossible
		if rng.Intn(4) == 0 {
			n = m + 100
		}
		var sb strings.Builder
		sb.WriteString(fmt.Sprintf("%d %d %d\n", n, m, h))
		for j, val := range s {
			if j > 0 {
				sb.WriteByte(' ')
			}
			sb.WriteString(fmt.Sprintf("%d", val))
		}
		sb.WriteByte('\n')
		input := sb.String()
		expect := solveB(n, m, h, s)
		out, err := run(bin, input)
		if err != nil {
			fmt.Printf("case %d: runtime error: %v\n%s", i+1, err, out)
			os.Exit(1)
		}
		out = strings.TrimSpace(out)
		exp := strings.TrimSpace(expect)
		if out != exp {
			fmt.Printf("case %d failed\ninput:\n%s\nexpected:%s\ngot:%s\n", i+1, input, exp, out)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
