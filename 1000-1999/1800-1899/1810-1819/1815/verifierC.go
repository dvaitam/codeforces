package main

import (
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"strings"
)

type Case struct {
	n     int
	m     int
	edges [][2]int
}

func genCases() []Case {
	rng := rand.New(rand.NewSource(1815))
	cases := make([]Case, 100)
	for i := range cases {
		n := rng.Intn(20) + 1
		maxEdges := n * (n - 1)
		m := 0
		if maxEdges > 0 {
			m = rng.Intn(min(20, maxEdges) + 1)
		}
		edges := make([][2]int, 0, m)
		seen := make(map[[2]int]bool)
		for len(edges) < m {
			a := rng.Intn(n) + 1
			b := rng.Intn(n) + 1
			if a == b {
				continue
			}
			key := [2]int{a, b}
			if seen[key] {
				continue
			}
			seen[key] = true
			edges = append(edges, key)
		}
		cases[i] = Case{n, m, edges}
	}
	return cases
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func buildRef() (string, error) {
	ref := "./refC.bin"
	cmd := exec.Command("go", "build", "-o", ref, "1815C.go")
	if out, err := cmd.CombinedOutput(); err != nil {
		return "", fmt.Errorf("failed to build reference: %v\n%s", err, out)
	}
	return ref, nil
}

func runBinary(bin, input string) (string, error) {
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

func runCase(bin, ref string, c Case) error {
	input := fmt.Sprintf("1\n%d %d\n", c.n, c.m)
	for _, e := range c.edges {
		input += fmt.Sprintf("%d %d\n", e[0], e[1])
	}
	expected, err := runBinary(ref, input)
	if err != nil {
		return fmt.Errorf("reference failed: %v", err)
	}
	got, err := runBinary(bin, input)
	if err != nil {
		return err
	}
	if strings.TrimSpace(got) != expected {
		return fmt.Errorf("expected %q got %q", expected, got)
	}
	return nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierC.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	ref, err := buildRef()
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	defer os.Remove(ref)

	cases := genCases()
	for i, c := range cases {
		if err := runCase(bin, ref, c); err != nil {
			input := fmt.Sprintf("1\n%d %d\n", c.n, c.m)
			for _, e := range c.edges {
				input += fmt.Sprintf("%d %d\n", e[0], e[1])
			}
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:\n%s", i+1, err, input)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
