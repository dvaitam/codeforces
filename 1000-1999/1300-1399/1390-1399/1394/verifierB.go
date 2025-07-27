package main

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"math/rand"
	"os"
	"os/exec"
	"strings"
	"time"
)

func generateTests() []string {
	r := rand.New(rand.NewSource(42))
	tests := make([]string, 100)
	for i := 0; i < 100; i++ {
		n := r.Intn(5) + 2
		k := r.Intn(3) + 1
		type Edge struct{ u, v, w int }
		edges := make([]Edge, 0)
		weight := 1
		for u := 1; u <= n; u++ {
			out := r.Intn(k) + 1
			if out > n-1 {
				out = n - 1
			}
			perm := r.Perm(n)
			count := 0
			idx := 0
			for count < out && idx < n {
				v := perm[idx] + 1
				idx++
				if v == u {
					continue
				}
				edges = append(edges, Edge{u, v, weight})
				weight++
				count++
			}
		}
		m := len(edges)
		var b strings.Builder
		fmt.Fprintf(&b, "%d %d %d\n", n, m, k)
		for _, e := range edges {
			fmt.Fprintf(&b, "%d %d %d\n", e.u, e.v, e.w)
		}
		tests[i] = b.String()
	}
	return tests
}

func runCmd(cmdPath string, input string) (string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()
	c := exec.CommandContext(ctx, cmdPath)
	c.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	c.Stdout = &out
	c.Stderr = io.Discard
	err := c.Run()
	return strings.TrimSpace(out.String()), err
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: go run verifierB.go /path/to/binary")
		os.Exit(1)
	}
	solPath := "./solB.bin"
	if err := exec.Command("go", "build", "-o", solPath, "1394B.go").Run(); err != nil {
		fmt.Println("failed to build reference solution:", err)
		os.Exit(1)
	}
	defer os.Remove(solPath)
	bin := os.Args[1]

	tests := generateTests()
	for i, t := range tests {
		expect, err := runCmd(solPath, t)
		if err != nil {
			fmt.Printf("reference solution failed on test %d: %v\n", i+1, err)
			os.Exit(1)
		}
		got, err := runCmd(bin, t)
		if err != nil {
			fmt.Printf("binary failed on test %d: %v\n", i+1, err)
			os.Exit(1)
		}
		if expect != strings.TrimSpace(got) {
			fmt.Printf("mismatch on test %d\nexpected: %s\n got: %s\n", i+1, expect, got)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
