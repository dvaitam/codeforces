package main

import (
	"bytes"
	"context"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"strings"
	"time"
)

type edge struct{ u, v int }

func run(bin string, input string) (string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()
	cmd := exec.CommandContext(ctx, bin)
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	err := cmd.Run()
	if ctx.Err() == context.DeadlineExceeded {
		return "", fmt.Errorf("timeout")
	}
	if err != nil {
		return "", err
	}
	return strings.TrimSpace(out.String()), nil
}

func genGraph() (int, []edge) {
	c := rand.Intn(3) + 1
	edges := []edge{}
	node := 2
	for i := 0; i < c; i++ {
		L := rand.Intn(5) + 2
		prev := 1
		for j := 0; j < L-1; j++ {
			cur := node
			edges = append(edges, edge{prev, cur})
			prev = cur
			node++
		}
		edges = append(edges, edge{prev, 1})
	}
	n := node - 1
	return n, edges
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierB.go /path/to/binary")
		os.Exit(1)
	}
	userBin := os.Args[1]
	refBin := "./1425B_ref"
	if err := exec.Command("go", "build", "-o", refBin, "1425B.go").Run(); err != nil {
		fmt.Fprintln(os.Stderr, "failed to build reference solution:", err)
		os.Exit(1)
	}
	rand.Seed(42)
	for i := 0; i < 100; i++ {
		n, edges := genGraph()
		m := len(edges)
		var sb strings.Builder
		sb.WriteString(fmt.Sprintf("%d %d\n", n, m))
		for _, e := range edges {
			sb.WriteString(fmt.Sprintf("%d %d\n", e.u, e.v))
		}
		input := sb.String()
		expect, err := run(refBin, input)
		if err != nil {
			fmt.Fprintln(os.Stderr, "reference failed on test", i+1, ":", err)
			os.Exit(1)
		}
		got, err := run(userBin, input)
		if err != nil {
			fmt.Fprintln(os.Stderr, "program failed on test", i+1, ":", err)
			os.Exit(1)
		}
		if expect != strings.TrimSpace(got) {
			fmt.Fprintf(os.Stderr, "mismatch on test %d: expected %s got %s\n", i+1, expect, got)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed!")
}
