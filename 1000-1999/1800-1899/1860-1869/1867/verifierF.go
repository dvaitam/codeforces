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

func run(bin, input string) (string, error) {
	var cmd *exec.Cmd
	if strings.HasSuffix(bin, ".go") {
		cmd = exec.Command("go", "run", bin)
	} else {
		cmd = exec.Command(bin)
	}
	cmd.Stdin = strings.NewReader(input)
	var out, stderr bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &stderr
	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("runtime error: %v\n%s", err, stderr.String())
	}
	return strings.TrimSpace(out.String()), nil
}

func expected(n int) string {
	var sb strings.Builder
	for i := 1; i < n; i++ {
		fmt.Fprintf(&sb, "%d %d\n", i, i+1)
	}
	return strings.TrimSpace(sb.String())
}

func genCase(r *rand.Rand) (string, string) {
	n := r.Intn(10) + 2
	edges := make([][2]int, n-1)
	for i := 2; i <= n; i++ {
		p := r.Intn(i-1) + 1
		edges[i-2] = [2]int{p, i}
	}
	var sb strings.Builder
	fmt.Fprintf(&sb, "%d\n", n)
	for _, e := range edges {
		fmt.Fprintf(&sb, "%d %d\n", e[0], e[1])
	}
	return sb.String(), expected(n)
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierF.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 1; i <= 100; i++ {
		in, exp := genCase(r)
		out, err := run(bin, in)
		if err != nil {
			fmt.Printf("test %d: %v\n", i, err)
			os.Exit(1)
		}
		if strings.TrimSpace(out) != exp {
			fmt.Printf("test %d failed\ninput:\n%sexpected:\n%s\ngot:\n%s\n", i, in, exp, out)
			os.Exit(1)
		}
	}
	fmt.Println("All 100 tests passed")
}
