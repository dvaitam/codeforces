package main

import (
	"bytes"
	"context"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"
)

type query struct {
	n int64
	k int
}
type test struct {
	q  int
	qs []query
}

func genTests() []test {
	rand.Seed(4)
	tests := make([]test, 0, 100)
	for len(tests) < 100 {
		q := rand.Intn(3) + 1
		qs := make([]query, q)
		for i := 0; i < q; i++ {
			n := int64(rand.Intn(100) + 1)
			k := rand.Intn(4)
			qs[i] = query{n, k}
		}
		tests = append(tests, test{q, qs})
	}
	return tests
}

func buildInput(t test) string {
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("%d\n", t.q))
	for _, qu := range t.qs {
		sb.WriteString(fmt.Sprintf("%d %d\n", qu.n, qu.k))
	}
	return sb.String()
}

func runBinary(path, input string) (string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	var cmd *exec.Cmd
	if strings.HasSuffix(path, ".go") {
		cmd = exec.CommandContext(ctx, "go", "run", path)
	} else {
		cmd = exec.CommandContext(ctx, path)
	}
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	err := cmd.Run()
	return out.String(), err
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("usage: go run verifierD.go /path/to/binary")
		os.Exit(1)
	}
	cand := os.Args[1]
	// build reference binary
	refBin := filepath.Join(os.TempDir(), "ref300D")
	build := exec.Command("go", "build", "-o", refBin, "300D.go")
	build.Dir = filepath.Dir(os.Args[0])
	if err := build.Run(); err != nil {
		fmt.Println("failed to build reference", err)
		os.Exit(1)
	}

	tests := genTests()
	for i, t := range tests {
		input := buildInput(t)
		candOut, err := runBinary(cand, input)
		if err != nil {
			fmt.Printf("test %d: run error %v\n", i+1, err)
			os.Exit(1)
		}
		refOut, err := runBinary(refBin, input)
		if err != nil {
			fmt.Printf("test %d: reference error %v\n", i+1, err)
			os.Exit(1)
		}
		if strings.TrimSpace(candOut) != strings.TrimSpace(refOut) {
			fmt.Printf("test %d failed. input:\n%sExpected:\n%sGot:\n%s\n", i+1, input, refOut, candOut)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
