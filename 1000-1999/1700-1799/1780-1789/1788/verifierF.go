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

func buildSolution() (string, error) {
	exe := "./_verifier_solF"
	cmd := exec.Command("go", "build", "-o", exe, "1788F.go")
	if err := cmd.Run(); err != nil {
		return "", err
	}
	return exe, nil
}

func runBinary(path, input string) (string, error) {
	var cmd *exec.Cmd
	if strings.HasSuffix(path, ".go") {
		cmd = exec.Command("go", "run", path)
	} else {
		cmd = exec.Command(path)
	}
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	var errBuf bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &errBuf
	timer := time.AfterFunc(4*time.Second, func() { cmd.Process.Kill() })
	err := cmd.Run()
	timer.Stop()
	if err != nil {
		return out.String(), fmt.Errorf("runtime error: %v\n%s", err, errBuf.String())
	}
	return strings.TrimSpace(out.String()), nil
}

type testCase struct{ input string }

func genTreeEdges(rng *rand.Rand, n int) [][2]int {
	edges := make([][2]int, n-1)
	for i := 2; i <= n; i++ {
		p := rng.Intn(i-1) + 1
		edges[i-2] = [2]int{p, i}
	}
	return edges
}

func genTest(rng *rand.Rand) testCase {
	n := rng.Intn(6) + 2
	q := rng.Intn(5)
	edges := genTreeEdges(rng, n)
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("%d %d\n", n, q))
	for _, e := range edges {
		sb.WriteString(fmt.Sprintf("%d %d\n", e[0], e[1]))
	}
	for i := 0; i < q; i++ {
		u := rng.Intn(n) + 1
		v := rng.Intn(n) + 1
		for v == u {
			v = rng.Intn(n) + 1
		}
		x := rng.Intn(16)
		sb.WriteString(fmt.Sprintf("%d %d %d\n", u, v, x))
	}
	return testCase{input: sb.String()}
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("usage: go run verifierF.go /path/to/binary")
		return
	}
	bin := os.Args[1]
	rand.Seed(time.Now().UnixNano())
	solExe, err := buildSolution()
	if err != nil {
		fmt.Println("failed to build reference solution:", err)
		return
	}
	defer os.Remove(solExe)

	fixed := []testCase{{input: "2 0\n1 2\n"}, {input: "3 1\n1 2\n2 3\n1 3 7\n"}}
	tests := make([]testCase, 0, 100)
	tests = append(tests, fixed...)
	for len(tests) < 100 {
		tests = append(tests, genTest(rand.New(rand.NewSource(rand.Int63()))))
	}
	for i, tc := range tests {
		exp, err := runBinary(solExe, tc.input)
		if err != nil {
			fmt.Printf("reference solution failed on test %d: %v\n", i+1, err)
			return
		}
		out, err := runBinary(bin, tc.input)
		if err != nil {
			fmt.Printf("candidate runtime error on test %d: %v\n", i+1, err)
			return
		}
		if strings.TrimSpace(exp) != strings.TrimSpace(out) {
			fmt.Printf("wrong answer on test %d\nexpected:\n%s\nactual:\n%s\n", i+1, exp, out)
			return
		}
	}
	fmt.Printf("All %d tests passed.\n", len(tests))
}
