package main

import (
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"strings"
)

func run(bin, in string) (string, error) {
	cmd := exec.Command(bin)
	cmd.Stdin = strings.NewReader(in)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = os.Stderr
	err := cmd.Run()
	return strings.TrimSpace(out.String()), err
}

func runRef(in string) (string, error) {
	cmd := exec.Command("go", "run", "1659F.go")
	cmd.Stdin = strings.NewReader(in)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = os.Stderr
	err := cmd.Run()
	return strings.TrimSpace(out.String()), err
}

func genCase() string {
	n := rand.Intn(8) + 3
	x := rand.Intn(n) + 1
	type edge struct{ u, v int }
	edges := make([]edge, 0, n-1)
	for i := 2; i <= n; i++ { // build tree
		j := rand.Intn(i-1) + 1
		edges = append(edges, edge{j, i})
	}
	perm := rand.Perm(n)
	var sb strings.Builder
	sb.WriteString("1\n")
	fmt.Fprintf(&sb, "%d %d\n", n, x)
	for _, e := range edges {
		fmt.Fprintf(&sb, "%d %d\n", e.u, e.v)
	}
	for i, v := range perm {
		if i > 0 {
			sb.WriteByte(' ')
		}
		fmt.Fprintf(&sb, "%d", v+1)
	}
	sb.WriteByte('\n')
	return sb.String()
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierF.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rand.Seed(1)
	for t := 1; t <= 100; t++ {
		input := genCase()
		exp, err := runRef(input)
		if err != nil {
			fmt.Println("reference failed:", err)
			os.Exit(1)
		}
		got, err := run(bin, input)
		if err != nil {
			fmt.Printf("test %d exec failed: %v\n", t, err)
			os.Exit(1)
		}
		if got != exp {
			fmt.Printf("test %d failed\nexpected: %s\ngot: %s\n", t, exp, got)
			os.Exit(1)
		}
	}
	fmt.Println("ok")
}
