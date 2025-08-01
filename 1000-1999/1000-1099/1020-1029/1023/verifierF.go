package main

import (
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"sort"
	"strings"
	"time"
)

func buildRef() (string, error) {
	ref := "refF.bin"
	cmd := exec.Command("go", "build", "-o", ref, "1023F.go")
	if out, err := cmd.CombinedOutput(); err != nil {
		return "", fmt.Errorf("failed to build reference: %v\n%s", err, out)
	}
	return ref, nil
}

func run(bin, input string) (string, error) {
	var cmd *exec.Cmd
	if strings.HasSuffix(bin, ".go") {
		cmd = exec.Command("go", "run", bin)
	} else {
		if !strings.Contains(bin, "/") {
			bin = "./" + bin
		}
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

func deterministicCases() []string {
	return []string{
		"2 1 0\n1 2\n\n",
	}
}

func randomCase(rng *rand.Rand) string {
	n := rng.Intn(5) + 2
	k := n - 1
	m := rng.Intn(3)
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("%d %d %d\n", n, k, m))
	for i := 1; i <= k; i++ {
		sb.WriteString(fmt.Sprintf("%d %d\n", i, i+1))
	}
	type edge struct{ a, b, w int }
	edges := make([]edge, m)
	for i := 0; i < m; i++ {
		a := rng.Intn(n) + 1
		b := rng.Intn(n) + 1
		for b == a {
			b = rng.Intn(n) + 1
		}
		w := rng.Intn(10) + 1
		edges[i] = edge{a, b, w}
	}
	sort.Slice(edges, func(i, j int) bool { return edges[i].w < edges[j].w })
	for _, e := range edges {
		sb.WriteString(fmt.Sprintf("%d %d %d\n", e.a, e.b, e.w))
	}
	return sb.String()
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("usage: go run verifierF.go /path/to/binary")
		os.Exit(1)
	}
	candidate := os.Args[len(os.Args)-1]
	ref, err := buildRef()
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	defer os.Remove(ref)
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	cases := deterministicCases()
	for len(cases) < 100 {
		cases = append(cases, randomCase(rng))
	}
	for i, in := range cases {
		exp, err := run(ref, in)
		if err != nil {
			fmt.Fprintf(os.Stderr, "reference failed on case %d: %v\n", i+1, err)
			os.Exit(1)
		}
		out, err := run(candidate, in)
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:\n%s", i+1, err, in)
			os.Exit(1)
		}
		if out != exp {
			fmt.Fprintf(os.Stderr, "case %d wrong answer\nexpected:\n%s\ngot:\n%s\ninput:\n%s", i+1, exp, out, in)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
