package main

import (
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

const numTestsE = 100

func main() {
	if len(os.Args) != 2 {
		fmt.Println("Usage: verifierE.go <binary>")
		os.Exit(1)
	}
	binPath, cleanup, err := prepareBinary(os.Args[1])
	if err != nil {
		fmt.Println("compile error:", err)
		os.Exit(1)
	}
	if cleanup != nil {
		defer cleanup()
	}
	r := rand.New(rand.NewSource(1))
	for t := 1; t <= numTestsE; t++ {
		n := r.Intn(10) + 1
		a := make([]int, n)
		b := make([]int, n)
		var sb strings.Builder
		sb.WriteString(fmt.Sprintf("%d\n", n))
		for i := 0; i < n; i++ {
			a[i] = r.Intn(n)
			if i > 0 {
				sb.WriteByte(' ')
			}
			sb.WriteString(fmt.Sprintf("%d", a[i]))
		}
		sb.WriteByte('\n')
		for i := 0; i < n; i++ {
			b[i] = r.Intn(n)
			if i > 0 {
				sb.WriteByte(' ')
			}
			sb.WriteString(fmt.Sprintf("%d", b[i]))
		}
		sb.WriteByte('\n')
		input := sb.String()
		expected := solveE(a, b)
		out, err := run(binPath, input)
		if err != nil {
			fmt.Printf("test %d: runtime error: %v\n", t, err)
			os.Exit(1)
		}
		out = strings.TrimSpace(out)
		if out != expected {
			fmt.Printf("test %d failed\ninput:%sexpected:%s got:%s\n", t, input, expected, out)
			os.Exit(1)
		}
	}
	fmt.Println("OK")
}

func prepareBinary(path string) (string, func(), error) {
	if strings.HasSuffix(path, ".go") {
		tmp := filepath.Join(os.TempDir(), "verify_binE")
		cmd := exec.Command("go", "build", "-o", tmp, path)
		if out, err := cmd.CombinedOutput(); err != nil {
			return "", nil, fmt.Errorf("go build failed: %v: %s", err, string(out))
		}
		return tmp, func() { os.Remove(tmp) }, nil
	}
	return path, nil, nil
}

func run(bin, input string) (string, error) {
	cmd := exec.Command(bin)
	cmd.Stdin = strings.NewReader(input)
	var buf bytes.Buffer
	cmd.Stdout = &buf
	cmd.Stderr = os.Stderr
	err := cmd.Run()
	return buf.String(), err
}

type Fenwick struct {
	n    int
	tree []int
}

func NewFenwick(n int) *Fenwick {
	return &Fenwick{n: n, tree: make([]int, n+1)}
}

func (f *Fenwick) Add(i, delta int) {
	for j := i + 1; j <= f.n; j += j & -j {
		f.tree[j] += delta
	}
}

func (f *Fenwick) Sum(i int) int {
	if i < 0 {
		return 0
	}
	res := 0
	for j := i + 1; j > 0; j -= j & -j {
		res += f.tree[j]
	}
	return res
}

func (f *Fenwick) FindKth(k int) int {
	idx := 0
	bit := 1
	for bit<<1 <= f.n {
		bit <<= 1
	}
	for ; bit > 0; bit >>= 1 {
		nxt := idx + bit
		if nxt <= f.n && f.tree[nxt] < k {
			idx = nxt
			k -= f.tree[nxt]
		}
	}
	return idx
}

func solveE(a, b []int) string {
	n := len(a)
	fenw := NewFenwick(n)
	for _, x := range b {
		fenw.Add(x, 1)
	}
	c := make([]int, n)
	for i := 0; i < n; i++ {
		want := (n - a[i]) % n
		s := fenw.Sum(want - 1)
		total := fenw.Sum(n - 1)
		var k int
		if total-s > 0 {
			k = s + 1
		} else {
			k = 1
		}
		idx := fenw.FindKth(k)
		c[i] = (a[i] + idx) % n
		fenw.Add(idx, -1)
	}
	var sb strings.Builder
	for i := 0; i < n; i++ {
		if i > 0 {
			sb.WriteByte(' ')
		}
		sb.WriteString(fmt.Sprintf("%d", c[i]))
	}
	return sb.String()
}
