package main

import (
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
	"time"
)

type BIT struct {
	n    int
	tree []int
}

func NewBIT(n int) *BIT {
	b := &BIT{n: n + 2, tree: make([]int, n+3)}
	return b
}

func (b *BIT) Add(idx, val int) {
	idx++
	for idx <= b.n {
		b.tree[idx] += val
		idx += idx & -idx
	}
}

func (b *BIT) Sum(idx int) int {
	if idx < 0 {
		return 0
	}
	if idx >= b.n {
		idx = b.n - 1
	}
	idx++
	res := 0
	for idx > 0 {
		res += b.tree[idx]
		idx -= idx & -idx
	}
	return res
}

func buildBinary(src, tag string) (string, error) {
	if strings.HasSuffix(src, ".go") {
		out := filepath.Join(os.TempDir(), tag)
		cmd := exec.Command("go", "build", "-o", out, src)
		if outb, err := cmd.CombinedOutput(); err != nil {
			return "", fmt.Errorf("build %s: %v\n%s", src, err, string(outb))
		}
		return out, nil
	}
	return src, nil
}

func runCase(exe, input, expected string) error {
	cmd := exec.Command(exe)
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("runtime error: %v\n%s", err, out.String())
	}
	got := strings.TrimSpace(out.String())
	exp := strings.TrimSpace(expected)
	if got != exp {
		return fmt.Errorf("expected %q got %q", exp, got)
	}
	return nil
}

func solveCase(n int, s string) int64 {
	size := 2*n + 5
	bits := []*BIT{NewBIT(size), NewBIT(size), NewBIT(size)}
	offset := n + 2
	sum := 0
	bits[sum%3].Add(sum+offset, 1)
	var ans int64
	for i := 0; i < n; i++ {
		if s[i] == '+' {
			sum--
		} else {
			sum++
		}
		mod := ((sum % 3) + 3) % 3
		idx := sum + offset
		ans += int64(bits[mod].Sum(idx))
		bits[mod].Add(idx, 1)
	}
	return ans
}

func randString(r *rand.Rand, n int) string {
	b := make([]byte, n)
	for i := range b {
		if r.Intn(2) == 0 {
			b[i] = '+'
		} else {
			b[i] = '-'
		}
	}
	return string(b)
}

func generateCase(r *rand.Rand) (string, string) {
	n := r.Intn(50) + 1
	s := randString(r, n)
	input := fmt.Sprintf("1\n%d\n%s\n", n, s)
	expect := fmt.Sprintf("%d\n", solveCase(n, s))
	return input, expect
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierF2.go /path/to/binary")
		return
	}
	candSrc := os.Args[1]
	_, file, _, _ := runtime.Caller(0)
	dir := filepath.Dir(file)
	refSrc := filepath.Join(dir, "1660F2.go")

	cand, err := buildBinary(candSrc, "candF2.bin")
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	ref, err := buildBinary(refSrc, "refF2.bin")
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < 100; i++ {
		in, exp := generateCase(rng)
		if err := runCase(ref, in, exp); err != nil {
			fmt.Fprintf(os.Stderr, "reference failed: %v\n", err)
			os.Exit(1)
		}
		if err := runCase(cand, in, exp); err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:\n%s", i+1, err, in)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
