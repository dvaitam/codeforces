package main

import (
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"strings"
)

func buildRef() (string, error) {
	ref := "refE.bin"
	cmd := exec.Command("go", "build", "-o", ref, "1391E.go")
	if out, err := cmd.CombinedOutput(); err != nil {
		return "", fmt.Errorf("failed to build reference: %v\n%s", err, out)
	}
	return ref, nil
}

func runBinary(bin, input string) (string, string, error) {
	var cmd *exec.Cmd
	if strings.HasSuffix(bin, ".go") {
		cmd = exec.Command("go", "run", bin)
	} else {
		cmd = exec.Command(bin)
	}
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	var errBuf bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &errBuf
	err := cmd.Run()
	return strings.TrimSpace(out.String()), errBuf.String(), err
}

type testCase struct{ input string }

func genCase(rng *rand.Rand) testCase {
	n := rng.Intn(5) + 2
	maxEdges := n * (n - 1) / 2
	m := rng.Intn(maxEdges-(n-1)+1) + (n - 1)
	edges := make(map[[2]int]bool)
	// build a tree to ensure connectivity
	for i := 2; i <= n; i++ {
		u := i - 1
		v := i
		if u > v {
			u, v = v, u
		}
		edges[[2]int{u, v}] = true
	}
	for len(edges) < m {
		u := rng.Intn(n) + 1
		v := rng.Intn(n) + 1
		if u == v {
			continue
		}
		if u > v {
			u, v = v, u
		}
		pair := [2]int{u, v}
		if edges[pair] {
			continue
		}
		edges[pair] = true
	}
	var sb strings.Builder
	sb.WriteString("1\n")
	sb.WriteString(fmt.Sprintf("%d %d\n", n, len(edges)))
	for e := range edges {
		sb.WriteString(fmt.Sprintf("%d %d\n", e[0], e[1]))
	}
	return testCase{sb.String()}
}

func runCase(bin, ref string, tc testCase, idx int) error {
	expOut, expErr, err := runBinary(ref, tc.input)
	if err != nil {
		return fmt.Errorf("reference failed: %v\nstderr: %s", err, expErr)
	}
	gotOut, gotErr, err := runBinary(bin, tc.input)
	if err != nil {
		return fmt.Errorf("test %d: runtime error: %v\nstderr: %s", idx, err, gotErr)
	}
	if strings.TrimSpace(gotOut) != strings.TrimSpace(expOut) {
		return fmt.Errorf("test %d failed\nexpected:\n%s\n got:\n%s", idx, expOut, gotOut)
	}
	return nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierE.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	ref, err := buildRef()
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	defer os.Remove(ref)

	rng := rand.New(rand.NewSource(99))
	for i := 1; i <= 100; i++ {
		tc := genCase(rng)
		if err := runCase(bin, ref, tc, i); err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
