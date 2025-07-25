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

func runBinary(path, input string) (string, error) {
	cmd := exec.Command(path)
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	err := cmd.Run()
	return strings.TrimSpace(out.String()), err
}

func generateCase() string {
	n := rand.Intn(10) + 1
	m := rand.Intn(n) + n // ensure connected by extra edges
	k := rand.Intn(n) + 1
	edges := make([][2]int, 0, m)
	for i := 2; i <= n; i++ {
		p := rand.Intn(i-1) + 1
		edges = append(edges, [2]int{p, i})
	}
	for len(edges) < m {
		a := rand.Intn(n) + 1
		b := rand.Intn(n) + 1
		edges = append(edges, [2]int{a, b})
	}
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("%d %d %d\n", n, len(edges), k))
	for _, e := range edges {
		sb.WriteString(fmt.Sprintf("%d %d\n", e[0], e[1]))
	}
	return sb.String()
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierE.go /path/to/binary")
		os.Exit(1)
	}
	target := os.Args[1]
	rand.Seed(time.Now().UnixNano())

	ref := "refE_bin"
	if err := exec.Command("go", "build", "-o", ref, "780E.go").Run(); err != nil {
		fmt.Fprintln(os.Stderr, "failed to build reference solution:", err)
		os.Exit(1)
	}
	defer os.Remove(ref)

	for i := 0; i < 100; i++ {
		input := generateCase()
		want, err := runBinary("./"+ref, input)
		if err != nil {
			fmt.Fprintln(os.Stderr, "failed to run reference solution:", err)
			os.Exit(1)
		}
		got, err := runBinary(target, input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "test %d: candidate error: %v\n", i+1, err)
			os.Exit(1)
		}
		if want != got {
			fmt.Printf("test %d failed\ninput:\n%sexpected:\n%s\ngot:\n%s\n", i+1, input, want, got)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
