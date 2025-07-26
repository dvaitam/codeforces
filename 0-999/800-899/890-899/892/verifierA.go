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

type testCase struct {
	n int
	a []int64
	b []int64
}

func generateTest(i int) testCase {
	r := rand.New(rand.NewSource(int64(i + 1)))
	n := r.Intn(50) + 2
	a := make([]int64, n)
	b := make([]int64, n)
	for j := 0; j < n; j++ {
		a[j] = r.Int63n(1000)
		b[j] = a[j] + r.Int63n(1000)
	}
	return testCase{n: n, a: a, b: b}
}

func buildTests() []testCase {
	tests := make([]testCase, 100)
	for i := range tests {
		tests[i] = generateTest(i)
	}
	return tests
}

func expected(t testCase) string {
	var total int64
	for _, x := range t.a {
		total += x
	}
	// find two largest in b
	var max1, max2 int64
	for _, x := range t.b {
		if x > max1 {
			max2 = max1
			max1 = x
		} else if x > max2 {
			max2 = x
		}
	}
	if total <= max1+max2 {
		return "YES"
	}
	return "NO"
}

func runBinary(bin string, input string) (string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()
	cmd := exec.CommandContext(ctx, bin)
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	err := cmd.Run()
	if ctx.Err() == context.DeadlineExceeded {
		return out.String(), fmt.Errorf("timeout")
	}
	if err != nil {
		return out.String(), err
	}
	return out.String(), nil
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("usage: verifierA <binary>")
		os.Exit(1)
	}
	bin := os.Args[1]
	tests := buildTests()
	for i, t := range tests {
		input := fmt.Sprintf("%d\n", t.n)
		for j, v := range t.a {
			if j > 0 {
				input += " "
			}
			input += fmt.Sprint(v)
		}
		input += "\n"
		for j, v := range t.b {
			if j > 0 {
				input += " "
			}
			input += fmt.Sprint(v)
		}
		input += "\n"

		out, err := runBinary(bin, input)
		exp := expected(t)
		got := strings.TrimSpace(strings.ToUpper(out))
		if err != nil {
			fmt.Printf("Test %d: runtime error: %v\nOutput:%s\n", i+1, err, out)
			os.Exit(1)
		}
		if got != exp {
			fmt.Printf("Test %d failed: expected %s got %s\n", i+1, exp, got)
			fmt.Printf("Input:\n%s\n", input)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
