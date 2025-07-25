package main

import (
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"strings"
)

func solveB(n, m int) []int {
	res := make([]int, 0, m)
	for i := 1; i <= n; i++ {
		lNonWindow := 2*n + 2*(i-1) + 1
		lWindow := 2*(i-1) + 1
		rNonWindow := 2*n + 2*(i-1) + 2
		rWindow := 2*(i-1) + 2
		if lNonWindow <= m {
			res = append(res, lNonWindow)
		}
		if lWindow <= m {
			res = append(res, lWindow)
		}
		if rNonWindow <= m {
			res = append(res, rNonWindow)
		}
		if rWindow <= m {
			res = append(res, rWindow)
		}
	}
	return res
}

func runBinary(binPath string, input string) (string, error) {
	cmd := exec.Command(binPath)
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = os.Stderr
	err := cmd.Run()
	return out.String(), err
}

func main() {
	if len(os.Args) < 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierB.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rand.Seed(2)
	const tests = 100
	for t := 0; t < tests; t++ {
		n := rand.Intn(20) + 1
		m := rand.Intn(4*n) + 1
		var sb strings.Builder
		fmt.Fprintf(&sb, "%d %d\n", n, m)
		expected := solveB(n, m)
		var exp strings.Builder
		for i, v := range expected {
			if i > 0 {
				exp.WriteByte(' ')
			}
			fmt.Fprint(&exp, v)
		}
		exp.WriteByte('\n')
		output, err := runBinary(bin, sb.String())
		if err != nil {
			fmt.Fprintf(os.Stderr, "test %d: runtime error: %v\n", t+1, err)
			os.Exit(1)
		}
		got := strings.TrimSpace(output)
		want := strings.TrimSpace(exp.String())
		if got != want {
			fmt.Fprintf(os.Stderr, "test %d failed\ninput:\n%s\nexpected:\n%s\ngot:\n%s\n", t+1, sb.String(), want, got)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
