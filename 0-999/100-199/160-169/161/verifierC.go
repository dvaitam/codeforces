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

var lenK [31]int64

func min(a, b int64) int64 {
	if a < b {
		return a
	}
	return b
}
func max(a, b int64) int64 {
	if a > b {
		return a
	}
	return b
}

var memo map[string]int64

func rec(k int, l1, r1, l2, r2 int64) int64 {
	if l1 > r1 || l2 > r2 || k <= 0 {
		return 0
	}
	key := fmt.Sprintf("%d_%d_%d_%d_%d", k, l1, r1, l2, r2)
	if v, ok := memo[key]; ok {
		return v
	}
	total := lenK[k]
	if l1 <= 1 && r1 >= total {
		return r2 - l2 + 1
	}
	if l2 <= 1 && r2 >= total {
		return r1 - l1 + 1
	}
	if l1 == l2 && r1 == r2 {
		return r1 - l1 + 1
	}
	mid := (total + 1) / 2
	var ans int64
	if l1 <= mid && mid <= r1 && l2 <= mid && mid <= r2 {
		left := min(mid-l1, mid-l2)
		right := min(r1-mid, r2-mid)
		ans = 1 + left + right
	}
	l1l, r1l := l1, min(r1, mid-1)
	l2l, r2l := l2, min(r2, mid-1)
	l1r, r1r := max(1, l1-mid), max(int64(0), r1-mid)
	l2r, r2r := max(1, l2-mid), max(int64(0), r2-mid)
	if l1l <= r1l && l2l <= r2l {
		if v := rec(k-1, l1l, r1l, l2l, r2l); v > ans {
			ans = v
		}
	}
	if l1r <= r1r && l2r <= r2r {
		if v := rec(k-1, l1r, r1r, l2r, r2r); v > ans {
			ans = v
		}
	}
	if l1l <= r1l && l2r <= r2r {
		if v := rec(k-1, l1l, r1l, l2r, r2r); v > ans {
			ans = v
		}
	}
	if l1r <= r1r && l2l <= r2l {
		if v := rec(k-1, l1r, r1r, l2l, r2l); v > ans {
			ans = v
		}
	}
	memo[key] = ans
	return ans
}

type Test struct{ l1, r1, l2, r2 int64 }

func generateTest() Test {
	l1 := int64(rand.Intn(100) + 1)
	r1 := l1 + int64(rand.Intn(100))
	l2 := int64(rand.Intn(100) + 1)
	r2 := l2 + int64(rand.Intn(100))
	return Test{l1, r1, l2, r2}
}

func (t Test) Input() string {
	return fmt.Sprintf("%d %d %d %d\n", t.l1, t.r1, t.l2, t.r2)
}

func solve(t Test) int64 {
	memo = make(map[string]int64)
	return rec(30, t.l1, t.r1, t.l2, t.r2)
}

func runBinary(path, input string) (string, error) {
	cmd := exec.Command(path)
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = os.Stderr
	err := cmd.Run()
	return strings.TrimSpace(out.String()), err
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierC.go <binary>")
		os.Exit(1)
	}
	bin := os.Args[1]
	lenK[0] = 0
	lenK[1] = 1
	for i := 2; i <= 30; i++ {
		lenK[i] = lenK[i-1]*2 + 1
	}
	rand.Seed(time.Now().UnixNano())
	for i := 0; i < 100; i++ {
		t := generateTest()
		inp := t.Input()
		exp := solve(t)
		out, err := runBinary(bin, inp)
		if err != nil {
			fmt.Fprintf(os.Stderr, "test %d runtime error: %v\n", i+1, err)
			os.Exit(1)
		}
		var got int64
		if _, e := fmt.Sscan(out, &got); e != nil {
			fmt.Fprintf(os.Stderr, "test %d failed to parse output: %v\n", i+1, e)
			os.Exit(1)
		}
		if got != exp {
			fmt.Fprintf(os.Stderr, "test %d failed: expected %d got %d\ninput:%s\n", i+1, exp, got, inp)
			os.Exit(1)
		}
	}
	fmt.Println("all tests passed")
}
