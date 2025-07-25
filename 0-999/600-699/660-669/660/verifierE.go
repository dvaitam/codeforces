package main

import (
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"strings"
)

const mod int64 = 1000000007

func solveE(n, m int) int64 {
	aPrev := int64(1)
	dPrev := int64(1)
	dPrevPrev := int64(0)
	for i := 1; i <= n; i++ {
		ai := (2*aPrev - dPrevPrev) % mod
		if ai < 0 {
			ai += mod
		}
		ai = ai * int64(m) % mod
		di := (ai + int64(m-1)*dPrev%mod) % mod
		dPrevPrev, dPrev, aPrev = dPrev, di, ai
	}
	return aPrev % mod
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
		fmt.Fprintln(os.Stderr, "usage: go run verifierE.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rand.Seed(5)
	const tests = 100
	for t := 0; t < tests; t++ {
		n := rand.Intn(20) + 1
		m := rand.Intn(20) + 1
		var sb strings.Builder
		fmt.Fprintf(&sb, "%d %d\n", n, m)
		expected := solveE(n, m)
		var exp strings.Builder
		fmt.Fprintln(&exp, expected)
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
