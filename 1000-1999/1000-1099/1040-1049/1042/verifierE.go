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
	tmp, err := os.CreateTemp("", "refE-*")
	if err != nil {
		return "", err
	}
	tmp.Close()
	path := tmp.Name()
	cmd := exec.Command("go", "build", "-o", path, "1042E.go")
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

type testCaseE struct {
	n, m   int
	matrix [][]int
	r, c   int
}

func generateTests() []testCaseE {
	rng := rand.New(rand.NewSource(46))
	tests := make([]testCaseE, 100)
	for i := range tests {
		n := rng.Intn(4) + 1
		m := rng.Intn(4) + 1
		matrix := make([][]int, n)
		for x := 0; x < n; x++ {
			row := make([]int, m)
			for y := 0; y < m; y++ {
				row[y] = rng.Intn(10)
			}
			matrix[x] = row
		}
		r := rng.Intn(n) + 1
		c := rng.Intn(m) + 1
		tests[i] = testCaseE{n, m, matrix, r, c}
	}
	return tests
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierE.go /path/to/binary")
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
		fmt.Fprintf(&sb, "%d %d\n", tc.n, tc.m)
		for _, row := range tc.matrix {
			for j, v := range row {
				if j > 0 {
					sb.WriteByte(' ')
				}
				sb.WriteString(fmt.Sprint(v))
			}
			sb.WriteByte('\n')
		}
		fmt.Fprintf(&sb, "%d %d\n", tc.r, tc.c)
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
