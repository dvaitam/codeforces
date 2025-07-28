package main

import (
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"strings"
)

func runExe(path string, input []byte) (string, error) {
	var cmd *exec.Cmd
	if strings.HasSuffix(path, ".go") {
		cmd = exec.Command("go", "run", path)
	} else {
		cmd = exec.Command(path)
	}
	cmd.Stdin = bytes.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	err := cmd.Run()
	return out.String(), err
}

func buildRef() (string, error) {
	ref := "./refE.bin"
	cmd := exec.Command("go", "build", "-o", ref, "1635E.go")
	if out, err := cmd.CombinedOutput(); err != nil {
		return "", fmt.Errorf("build reference failed: %v: %s", err, string(out))
	}
	return ref, nil
}

func buildCase(n, m int, rels [][3]int) []byte {
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("%d %d\n", n, m))
	for _, r := range rels {
		sb.WriteString(fmt.Sprintf("%d %d %d\n", r[0], r[1], r[2]))
	}
	return []byte(sb.String())
}

func genRandomCase(rng *rand.Rand) []byte {
	n := rng.Intn(4) + 1
	m := rng.Intn(n*n) + 1
	rels := make([][3]int, m)
	for i := 0; i < m; i++ {
		typ := rng.Intn(2) + 1
		u := rng.Intn(n) + 1
		v := rng.Intn(n) + 1
		for v == u {
			v = rng.Intn(n) + 1
		}
		rels[i] = [3]int{typ, u, v}
	}
	return buildCase(n, m, rels)
}

func genTests() [][]byte {
	rng := rand.New(rand.NewSource(46))
	tests := make([][]byte, 0, 100)
	tests = append(tests, buildCase(1, 0, nil))
	tests = append(tests, buildCase(2, 1, [][3]int{{1, 1, 2}}))
	for len(tests) < 100 {
		tests = append(tests, genRandomCase(rng))
	}
	return tests
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("Usage: go run verifierE.go /path/to/binary")
		return
	}
	bin := os.Args[1]
	ref, err := buildRef()
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	defer os.Remove(ref)

	tests := genTests()
	for i, tc := range tests {
		exp, err := runExe(ref, tc)
		if err != nil {
			fmt.Printf("reference runtime error on test %d: %v\n%s", i+1, err, exp)
			os.Exit(1)
		}
		got, err := runExe(bin, tc)
		if err != nil {
			fmt.Printf("candidate runtime error on test %d: %v\n%s", i+1, err, got)
			os.Exit(1)
		}
		if strings.TrimSpace(exp) != strings.TrimSpace(got) {
			fmt.Printf("test %d failed\ninput:\n%sexpected:\n%s\ngot:\n%s\n", i+1, string(tc), exp, got)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
