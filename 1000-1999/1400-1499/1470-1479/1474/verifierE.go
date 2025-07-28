package main

import (
	"bufio"
	"bytes"
	"context"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"time"
)

type Test struct {
	n   int
	ans int64
}

func maxTime(n int) int64 {
	var ans int64
	for x := 3; x <= n+1; x++ {
		d := int64(n - x/2)
		ans += d * d
	}
	return ans
}

func genTests() []Test {
	r := rand.New(rand.NewSource(42))
	tests := make([]Test, 0, 100)
	for len(tests) < 100 {
		n := r.Intn(48) + 2 // n from 2..49
		tests = append(tests, Test{n: n, ans: maxTime(n)})
	}
	return tests
}

func verifyCase(t Test, scanner *bufio.Reader) error {
	var ans int64
	if _, err := fmt.Fscan(scanner, &ans); err != nil {
		return fmt.Errorf("missing answer")
	}
	if ans != t.ans {
		return fmt.Errorf("wrong time %d expected %d", ans, t.ans)
	}
	perm := make([]int, t.n)
	for i := 0; i < t.n; i++ {
		if _, err := fmt.Fscan(scanner, &perm[i]); err != nil {
			return fmt.Errorf("missing permutation")
		}
	}
	used := make([]bool, t.n+1)
	for _, v := range perm {
		if v < 1 || v > t.n || used[v] {
			return fmt.Errorf("invalid permutation")
		}
		used[v] = true
	}
	var m int
	if _, err := fmt.Fscan(scanner, &m); err != nil {
		return fmt.Errorf("missing m")
	}
	ops := make([][2]int, m)
	for i := 0; i < m; i++ {
		if _, err := fmt.Fscan(scanner, &ops[i][0], &ops[i][1]); err != nil {
			return fmt.Errorf("missing op %d", i+1)
		}
	}
	// simulate
	p := make([]int, t.n+1)
	for i, v := range perm {
		p[i+1] = v
	}
	var total int64
	for _, op := range ops {
		i, j := op[0], op[1]
		if i < 1 || i > t.n || j < 1 || j > t.n || i == j {
			return fmt.Errorf("bad swap")
		}
		if p[j] != i {
			return fmt.Errorf("invalid operation")
		}
		p[i], p[j] = p[j], p[i]
		d := i - j
		if d < 0 {
			d = -d
		}
		total += int64(d * d)
	}
	for i := 1; i <= t.n; i++ {
		if p[i] != i {
			return fmt.Errorf("not identity")
		}
	}
	if total != t.ans {
		return fmt.Errorf("time mismatch")
	}
	return nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierE.go /path/to/binary")
		return
	}
	binary := os.Args[1]
	tests := genTests()

	var input bytes.Buffer
	fmt.Fprintln(&input, len(tests))
	for _, t := range tests {
		fmt.Fprintln(&input, t.n)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	cmd := exec.CommandContext(ctx, binary)
	cmd.Stdin = bytes.NewReader(input.Bytes())
	out, err := cmd.Output()
	if ctx.Err() == context.DeadlineExceeded {
		fmt.Println("time limit exceeded")
		return
	}
	if err != nil {
		fmt.Println("execution error:", err)
		return
	}

	scanner := bufio.NewReader(bytes.NewReader(out))
	for i, t := range tests {
		if err := verifyCase(t, scanner); err != nil {
			fmt.Printf("case %d failed: %v\n", i+1, err)
			return
		}
	}
	if _, err := fmt.Fscan(scanner, new(int)); err == nil {
		fmt.Println("extra output detected")
		return
	}
	fmt.Println("OK")
}
