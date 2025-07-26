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

func buildRef() (string, error) {
	tmp, err := os.CreateTemp("", "refF-*")
	if err != nil {
		return "", err
	}
	tmp.Close()
	path := tmp.Name()
	cmd := exec.Command("go", "build", "-o", path, "1042F.go")
	if out, err := cmd.CombinedOutput(); err != nil {
		return "", fmt.Errorf("build ref failed: %v\n%s", err, out)
	}
	return path, nil
}

func runExe(path, input string) (string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	cmd := exec.CommandContext(ctx, path)
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("runtime error: %v\n%s", err, out.String())
	}
	return strings.TrimSpace(out.String()), nil
}

type testCaseF struct {
	n, k  int
	edges [][2]int
}

func generateTests() []testCaseF {
	rng := rand.New(rand.NewSource(47))
	tests := make([]testCaseF, 100)
	for i := range tests {
		n := rng.Intn(15) + 3
		k := rng.Intn(5) + 1
		edges := make([][2]int, n-1)
		for j := 0; j < n-1; j++ {
			u := rng.Intn(n) + 1
			v := rng.Intn(n) + 1
			for v == u {
				v = rng.Intn(n) + 1
			}
			edges[j] = [2]int{u, v}
		}
		tests[i] = testCaseF{n, k, edges}
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
		fmt.Println(err)
		return
	}
	defer os.Remove(ref)
	tests := generateTests()
	for idx, tc := range tests {
		var sb strings.Builder
		fmt.Fprintf(&sb, "%d %d\n", tc.n, tc.k)
		for _, e := range tc.edges {
			fmt.Fprintf(&sb, "%d %d\n", e[0], e[1])
		}
		input := sb.String()
		exp, err := runExe(ref, input)
		if err != nil {
			fmt.Printf("ref error on test %d: %v\n", idx+1, err)
			return
		}
		got, err := runExe(bin, input)
		if err != nil {
			fmt.Printf("test %d runtime error: %v\n", idx+1, err)
			return
		}
		if strings.TrimSpace(got) != strings.TrimSpace(exp) {
			fmt.Printf("test %d failed\ninput:\n%sexpected %s got %s\n", idx+1, input, exp, got)
			return
		}
	}
	fmt.Printf("all %d tests passed\n", len(tests))
}
