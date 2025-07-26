package main

import (
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"strings"
)

type Test struct {
	input string
}

func runExe(path, input string) (string, error) {
	cmd := exec.Command(path)
	if strings.HasSuffix(path, ".go") {
		cmd = exec.Command("go", "run", path)
	}
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	err := cmd.Run()
	return out.String(), err
}

func buildRef() (string, error) {
	ref := "./refF.bin"
	cmd := exec.Command("go", "build", "-o", ref, "1110F.go")
	if out, err := cmd.CombinedOutput(); err != nil {
		return "", fmt.Errorf("build reference failed: %v: %s", err, string(out))
	}
	return ref, nil
}

func genTests() []Test {
	rand.Seed(6)
	var tests []Test
	for len(tests) < 100 {
		n := rand.Intn(13) + 3
		q := rand.Intn(5) + 1
		parent := make([]int, n+1)
		weight := make([]int, n+1)
		for i := 2; i <= n; i++ {
			parent[i] = rand.Intn(i-1) + 1
			weight[i] = rand.Intn(20) + 1
		}
		leaves := make([]bool, n+1)
		for i := 1; i <= n; i++ {
			leaves[i] = true
		}
		for i := 2; i <= n; i++ {
			leaves[parent[i]] = false
		}
		var leafIdx []int
		for i := 1; i <= n; i++ {
			if leaves[i] {
				leafIdx = append(leafIdx, i)
			}
		}
		var sb strings.Builder
		sb.WriteString(fmt.Sprintf("%d %d\n", n, q))
		for i := 2; i <= n; i++ {
			sb.WriteString(fmt.Sprintf("%d %d\n", parent[i], weight[i]))
		}
		for i := 0; i < q; i++ {
			v := rand.Intn(n) + 1
			lf := leafIdx[rand.Intn(len(leafIdx))]
			l := rand.Intn(lf) + 1
			r := lf + rand.Intn(n-lf+1)
			sb.WriteString(fmt.Sprintf("%d %d %d\n", v, l, r))
		}
		tests = append(tests, Test{sb.String()})
	}
	return tests
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierF.go /path/to/binary")
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
		exp, err := runExe(ref, tc.input)
		if err != nil {
			fmt.Printf("reference runtime error on test %d: %v\n", i+1, err)
			os.Exit(1)
		}
		got, err := runExe(bin, tc.input)
		if err != nil {
			fmt.Printf("candidate runtime error on test %d: %v\n", i+1, err)
			os.Exit(1)
		}
		if strings.TrimSpace(exp) != strings.TrimSpace(got) {
			fmt.Printf("Test %d failed\nInput:\n%sExpected:\n%sGot:\n%s\n", i+1, tc.input, exp, got)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
