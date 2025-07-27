package main

import (
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"sort"
	"strings"
	"time"
)

type testCaseA struct {
	n int
	x int
	a []int
	b []int
}

func generateA(rng *rand.Rand) testCaseA {
	n := rng.Intn(50) + 1   // 1..50
	x := rng.Intn(1000) + 1 // 1..1000
	a := make([]int, n)
	b := make([]int, n)
	for i := 0; i < n; i++ {
		a[i] = rng.Intn(x) + 1
	}
	for i := 0; i < n; i++ {
		b[i] = rng.Intn(x) + 1
	}
	sort.Ints(a)
	sort.Ints(b)
	return testCaseA{n: n, x: x, a: a, b: b}
}

func expectedA(tc testCaseA) string {
	for i := 0; i < tc.n; i++ {
		if tc.a[i]+tc.b[tc.n-1-i] > tc.x {
			return "No"
		}
	}
	return "Yes"
}

func run(bin, input string) (string, error) {
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
	fmt.Fprintf(&b, "1\n%d %d\n", tc.n, tc.x)
	for i, v := range tc.a {
		if i > 0 {
			b.WriteByte(' ')
		}
		fmt.Fprint(&b, v)
	}
	b.WriteByte('\n')
	for i, v := range tc.b {
		if i > 0 {
			b.WriteByte(' ')
		}
		fmt.Fprint(&b, v)
	}
	b.WriteByte('\n')
	expected := strings.ToLower(expectedA(tc))
	got, err := run(bin, b.String())
	if err != nil {
		return err
	}
	if strings.ToLower(strings.TrimSpace(got)) != expected {
		return fmt.Errorf("expected %s got %s", expected, strings.TrimSpace(got))
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
