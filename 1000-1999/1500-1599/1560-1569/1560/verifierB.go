package main

import (
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"strings"
)

type testCase struct {
	a int
	b int
	c int
}

func solve(a, b, c int) int {
	diff := abs(a - b)
	n := diff * 2
	if diff == 0 || max3(a, b, c) > n {
		return -1
	}
	if c > diff {
		return c - diff
	}
	return c + diff
}

func abs(x int) int {
	if x < 0 {
		return -x
	}
	return x
}

func max3(x, y, z int) int {
	if x > y {
		if x > z {
			return x
		}
		return z
	}
	if y > z {
		return y
	}
	return z
}

func genTests() []testCase {
	r := rand.New(rand.NewSource(2))
	tests := make([]testCase, 100)
	for i := 0; i < 50; i++ {
		diff := r.Intn(100) + 1
		n := diff * 2
		a := r.Intn(diff) + 1
		b := a + diff
		c := r.Intn(n) + 1
		tests[i] = testCase{a, b, c}
	}
	for i := 50; i < 100; i++ {
		a := r.Intn(200) + 1
		b := r.Intn(200) + 1
		c := r.Intn(300) + 1
		tests[i] = testCase{a, b, c}
	}
	return tests
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
		return out.String() + errBuf.String(), fmt.Errorf("runtime error: %v", err)
	}
	return strings.TrimSpace(out.String()), nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierB.go /path/to/binary")
		return
	}
	bin := os.Args[1]
	tests := genTests()
	for i, tc := range tests {
		input := fmt.Sprintf("1\n%d %d %d\n", tc.a, tc.b, tc.c)
		want := fmt.Sprintf("%d", solve(tc.a, tc.b, tc.c))
		got, err := run(bin, input)
		if err != nil {
			fmt.Printf("case %d: %v\n", i+1, err)
			return
		}
		if strings.TrimSpace(got) != want {
			fmt.Printf("case %d failed: expected %s got %s\n", i+1, want, got)
			fmt.Printf("input:\n%s", input)
			return
		}
	}
	fmt.Println("All tests passed")
}
