package main

import (
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"strings"
)

func solveB(n int, k int) string {
	s := make([]byte, n)
	for i := 0; i < n; i++ {
		s[i] = 'a'
	}
	for i := n - 2; i >= 0; i-- {
		cnt := n - i - 1
		if k > cnt {
			k -= cnt
		} else {
			s[i] = 'b'
			s[n-k] = 'b'
			break
		}
	}
	return string(s)
}

func run(bin, input string) (string, error) {
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

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierB.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rand.Seed(2)
	for t := 1; t <= 100; t++ {
		n := rand.Intn(50) + 3
		maxK := n * (n - 1) / 2
		k := rand.Intn(maxK) + 1
		input := fmt.Sprintf("1\n%d %d\n", n, k)
		expect := solveB(n, k)
		out, err := run(bin, input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "test %d failed: %v\ninput:\n%s", t, err, input)
			os.Exit(1)
		}
		if strings.TrimSpace(out) != expect {
			fmt.Fprintf(os.Stderr, "test %d failed: expected %s got %s\ninput:\n%s", t, expect, out, input)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
