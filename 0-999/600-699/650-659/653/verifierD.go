package main

import (
	"bytes"
	"fmt"
	"math"
	"math/rand"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"
	"time"
)

type Case struct{ input string }

func buildRef() (string, error) {
	refSrc := os.Getenv("REFERENCE_SOURCE_PATH")
	if refSrc == "" {
		return "", fmt.Errorf("REFERENCE_SOURCE_PATH not set")
	}
	ref := filepath.Join(os.TempDir(), "ref653D.bin")
	cmd := exec.Command("go", "build", "-o", ref, refSrc)
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
	var stdout bytes.Buffer
	var stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("runtime error: %v\n%s", err, stderr.String())
	}
	return strings.TrimSpace(stdout.String()), nil
}

func genCases() []Case {
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	cases := []Case{{"4 4 3\n1 2 2\n2 4 1\n1 3 1\n3 4 2\n"}}
	for len(cases) < 102 {
		n := rng.Intn(7) + 2
		x := rng.Intn(20) + 1
		// Generate a connected directed graph with no self-loops and no duplicate edges
		edges := make(map[[2]int]bool)
		var edgeList []struct{ u, v, c int }

		// Ensure connectivity: create a path 1->2->...->n
		for i := 1; i < n; i++ {
			key := [2]int{i, i + 1}
			if !edges[key] {
				edges[key] = true
				c := rng.Intn(1000) + 1
				edgeList = append(edgeList, struct{ u, v, c int }{i, i + 1, c})
			}
		}
		// Add some random edges (no self-loops, no duplicates)
		extraCount := rng.Intn(n*2) + 1
		for attempt := 0; attempt < extraCount*3; attempt++ {
			u := rng.Intn(n) + 1
			v := rng.Intn(n) + 1
			if u == v {
				continue
			}
			key := [2]int{u, v}
			if edges[key] {
				continue
			}
			edges[key] = true
			c := rng.Intn(1000) + 1
			edgeList = append(edgeList, struct{ u, v, c int }{u, v, c})
		}

		m := len(edgeList)
		var sb strings.Builder
		fmt.Fprintf(&sb, "%d %d %d\n", n, m, x)
		for _, e := range edgeList {
			fmt.Fprintf(&sb, "%d %d %d\n", e.u, e.v, e.c)
		}
		cases = append(cases, Case{sb.String()})
	}
	return cases
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierD.go /path/to/binary")
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
		expected, err := runBinary(ref, c.input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "reference failed on case %d: %v\n", i+1, err)
			os.Exit(1)
		}
		got, err := runBinary(bin, c.input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:\n%s", i+1, err, c.input)
			os.Exit(1)
		}

		// Parse as floats and compare with tolerance (relative or absolute error <= 1e-6)
		expVal, err1 := strconv.ParseFloat(strings.TrimSpace(expected), 64)
		gotVal, err2 := strconv.ParseFloat(strings.TrimSpace(got), 64)
		if err1 != nil {
			fmt.Fprintf(os.Stderr, "case %d: cannot parse reference output as float: %v\n", i+1, err1)
			os.Exit(1)
		}
		if err2 != nil {
			fmt.Fprintf(os.Stderr, "case %d: cannot parse candidate output as float: %v\n", i+1, err2)
			os.Exit(1)
		}
		diff := math.Abs(expVal - gotVal)
		denom := math.Max(1.0, math.Abs(expVal))
		if diff/denom > 1e-6 {
			fmt.Fprintf(os.Stderr, "case %d: expected %.10f got %.10f (rel error %.2e)\ninput:\n%s", i+1, expVal, gotVal, diff/denom, c.input)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
