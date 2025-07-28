package main

import (
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"strings"
)

func run(bin, input string) (string, error) {
	var cmd *exec.Cmd
	if strings.HasSuffix(bin, ".go") {
		cmd = exec.Command("go", "run", bin)
	} else {
		cmd = exec.Command(bin)
	}
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	var errb bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &errb
	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("runtime error: %v\n%s", err, errb.String())
	}
	return strings.TrimSpace(out.String()), nil
}

func expected(a []int) string {
	n := len(a)
	b := make([]int, n)
	copy(b, a)
	start := n % 2
	for i := start; i+1 < n; i += 2 {
		if b[i] > b[i+1] {
			b[i], b[i+1] = b[i+1], b[i]
		}
	}
	for i := 1; i < n; i++ {
		if b[i-1] > b[i] {
			return "NO"
		}
	}
	return "YES"
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierD.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rng := rand.New(rand.NewSource(45))
	for t := 0; t < 100; t++ {
		n := rng.Intn(10) + 1
		arr := make([]int, n)
		var sb strings.Builder
		sb.WriteString("1\n")
		sb.WriteString(fmt.Sprintf("%d\n", n))
		for i := 0; i < n; i++ {
			arr[i] = rng.Intn(100)
			sb.WriteString(fmt.Sprintf("%d ", arr[i]))
		}
		sb.WriteByte('\n')
		input := sb.String()
		exp := expected(arr)
		got, err := run(bin, input)
		if err != nil {
			fmt.Printf("runtime error on test %d: %v\ninput:\n%s\n", t+1, err, input)
			os.Exit(1)
		}
		if strings.TrimSpace(strings.ToUpper(got)) != exp {
			fmt.Printf("wrong answer on test %d\ninput:\n%s\nexpected: %s\ngot: %s\n", t+1, input, exp, got)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed.")
}
