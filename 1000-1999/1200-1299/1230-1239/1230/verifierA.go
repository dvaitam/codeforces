package main

import (
	"bytes"
	"context"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"sort"
	"strings"
	"time"
)

type testCase struct {
	a [4]int
}

func generateTests() []testCase {
	r := rand.New(rand.NewSource(1))
	tests := make([]testCase, 100)
	for i := 0; i < 100; i++ {
		tests[i] = testCase{a: [4]int{r.Intn(100) + 1, r.Intn(100) + 1, r.Intn(100) + 1, r.Intn(100) + 1}}
	}
	return tests
}

func expected(t testCase) string {
	a := []int{t.a[0], t.a[1], t.a[2], t.a[3]}
	sort.Ints(a)
	if a[0]+a[3] == a[1]+a[2] || a[0]+a[1]+a[2] == a[3] {
		return "YES"
	}
	return "NO"
}

func runBinary(bin, input string) (string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	cmd := exec.CommandContext(ctx, bin)
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	var errb bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &errb
	if err := cmd.Run(); err != nil {
		if ctx.Err() == context.DeadlineExceeded {
			return "", fmt.Errorf("timeout")
		}
		return "", fmt.Errorf("run failed: %v\n%s", err, errb.String())
	}
	if ctx.Err() == context.DeadlineExceeded {
		return "", fmt.Errorf("timeout")
	}
	return strings.TrimSpace(out.String()), nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierA.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	tests := generateTests()
	for i, t := range tests {
		input := fmt.Sprintf("%d %d %d %d\n", t.a[0], t.a[1], t.a[2], t.a[3])
		want := expected(t)
		got, err := runBinary(bin, input)
		if err != nil {
			fmt.Printf("test %d error: %v\n", i+1, err)
			os.Exit(1)
		}
		if strings.ToUpper(strings.TrimSpace(got)) != want {
			fmt.Printf("test %d failed: expected %s got %s\n", i+1, want, strings.TrimSpace(got))
			os.Exit(1)
		}
	}
	fmt.Printf("All %d tests passed\n", len(tests))
}
