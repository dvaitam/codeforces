package main

import (
	"bufio"
	"bytes"
	"context"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"sort"
	"time"
)

type Test struct {
	n        int
	arr      []int
	possible bool
}

func hasSolution(arr []int) bool {
	a := append([]int(nil), arr...)
	sort.Ints(a)
	n := len(a) / 2
	N := len(a)
	orig := make(map[int]int, N)
	for _, v := range a {
		orig[v]++
	}
	for i := 0; i < N-1; i++ {
		x := a[N-1] + a[i]
		b := make(map[int]int, len(orig))
		for k, v := range orig {
			b[k] = v
		}
		tf := x
		idx := N - 1
		ok := true
		for k := 0; k < n; k++ {
			for idx >= 0 && b[a[idx]] == 0 {
				idx--
			}
			if idx < 0 {
				ok = false
				break
			}
			v := a[idx]
			b[v]--
			need := tf - v
			if b[need] == 0 {
				ok = false
				break
			}
			b[need]--
			tf = v
		}
		if ok {
			return true
		}
	}
	return false
}

func genTests() []Test {
	r := rand.New(rand.NewSource(42))
	tests := make([]Test, 0, 100)
	for len(tests) < 100 {
		n := r.Intn(4) + 2 // n from 2..5 -> array len 4..10
		arr := make([]int, 2*n)
		for i := range arr {
			arr[i] = r.Intn(20)
		}
		possible := hasSolution(arr)
		tests = append(tests, Test{n: n, arr: arr, possible: possible})
	}
	return tests
}

func verifyCase(t Test, scanner *bufio.Reader) error {
	var res string
	if _, err := fmt.Fscan(scanner, &res); err != nil {
		return fmt.Errorf("missing YES/NO")
	}
	if res == "NO" {
		if t.possible {
			return fmt.Errorf("should be YES")
		}
		return nil
	}
	if res != "YES" {
		return fmt.Errorf("invalid token %s", res)
	}
	if !t.possible {
		return fmt.Errorf("should be NO")
	}
	var x int
	if _, err := fmt.Fscan(scanner, &x); err != nil {
		return fmt.Errorf("missing x")
	}
	counts := make(map[int]int)
	for _, v := range t.arr {
		counts[v]++
	}
	tf := x
	for i := 0; i < t.n; i++ {
		var a, b int
		if _, err := fmt.Fscan(scanner, &a, &b); err != nil {
			return fmt.Errorf("missing pair %d", i+1)
		}
		if counts[a] == 0 || counts[b] == 0 {
			return fmt.Errorf("pair uses absent number")
		}
		counts[a]--
		counts[b]--
		if a+b != tf {
			return fmt.Errorf("sum mismatch")
		}
		if a > b {
			tf = a
		} else {
			tf = b
		}
	}
	for _, v := range counts {
		if v != 0 {
			return fmt.Errorf("numbers left over")
		}
	}
	return nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierC.go /path/to/binary")
		return
	}
	binary := os.Args[1]
	tests := genTests()

	var input bytes.Buffer
	fmt.Fprintln(&input, len(tests))
	for _, t := range tests {
		fmt.Fprintln(&input, t.n)
		for i, v := range t.arr {
			if i > 0 {
				fmt.Fprint(&input, " ")
			}
			fmt.Fprint(&input, v)
		}
		fmt.Fprintln(&input)
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
	// ensure no extra tokens
	if _, err := fmt.Fscan(scanner, new(string)); err == nil {
		fmt.Println("extra output detected")
		return
	}
	fmt.Println("OK")
}
