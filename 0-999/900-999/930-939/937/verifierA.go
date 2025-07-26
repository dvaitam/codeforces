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

type testCaseA struct {
	scores []int
}

func generateA(rng *rand.Rand) testCaseA {
	n := rng.Intn(100) + 1 // 1..100
	arr := make([]int, n)
	hasPos := false
	for i := range arr {
		arr[i] = rng.Intn(601) // 0..600
		if arr[i] > 0 {
			hasPos = true
		}
	}
	if !hasPos {
		arr[0] = rng.Intn(600) + 1
	}
	return testCaseA{scores: arr}
}

func expectedA(tc testCaseA) int {
	uniq := make(map[int]struct{})
	for _, v := range tc.scores {
		if v > 0 {
			uniq[v] = struct{}{}
		}
	}
	return len(uniq)
}

func run(bin string, input string) (string, error) {
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
	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("runtime error: %v\n%s", err, errBuf.String())
	}
	return strings.TrimSpace(out.String()), nil
}

func runCase(bin string, tc testCaseA) error {
	var b strings.Builder
	fmt.Fprintf(&b, "%d\n", len(tc.scores))
	for i, v := range tc.scores {
		if i > 0 {
			b.WriteByte(' ')
		}
		fmt.Fprint(&b, v)
	}
	b.WriteByte('\n')
	expected := fmt.Sprintf("%d", expectedA(tc))
	got, err := run(bin, b.String())
	if err != nil {
		return err
	}
	if strings.TrimSpace(got) != expected {
		return fmt.Errorf("expected %s got %s", expected, got)
	}
	return nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierA.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < 100; i++ {
		tc := generateA(rng)
		if err := runCase(bin, tc); err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\n", i+1, err)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
