package main

import (
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"sort"
	"strings"
)

func buildOracle() (string, error) {
	_, file, _, _ := runtime.Caller(0)
	dir := filepath.Dir(file)
	src := filepath.Join(dir, "1569D.go")
	bin := filepath.Join(os.TempDir(), "oracle1569D.bin")
	cmd := exec.Command("go", "build", "-o", bin, src)
	if out, err := cmd.CombinedOutput(); err != nil {
		return "", fmt.Errorf("build oracle failed: %v\n%s", err, out)
	}
	return bin, nil
}

func run(bin, input string) (string, error) {
	cmd := exec.Command(bin)
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	var stderr bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &stderr
	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("runtime error: %v\n%s", err, stderr.String())
	}
	return strings.TrimSpace(out.String()), nil
}

func genCase(r *rand.Rand) string {
	n := r.Intn(3) + 2
	m := r.Intn(3) + 2
	k := r.Intn(5) + 2
	xs := make([]int, n)
	ys := make([]int, m)
	xs[0] = 0
	xs[n-1] = 1000000
	for i := 1; i < n-1; i++ {
		xs[i] = r.Intn(999999) + 1
	}
	ys[0] = 0
	ys[m-1] = 1000000
	for i := 1; i < m-1; i++ {
		ys[i] = r.Intn(999999) + 1
	}
	sort.Ints(xs)
	sort.Ints(ys)
	type pt struct{ x, y int }
	pts := make(map[pt]struct{}, k)
	var arr []pt
	for len(arr) < k {
		var x, y int
		if r.Intn(2) == 0 {
			x = xs[r.Intn(n)]
			y = r.Intn(1000001)
		} else {
			x = r.Intn(1000001)
			y = ys[r.Intn(m)]
		}
		p := pt{x, y}
		if _, ok := pts[p]; !ok {
			pts[p] = struct{}{}
			arr = append(arr, p)
		}
	}
	var sb strings.Builder
	fmt.Fprintf(&sb, "1\n%d %d %d\n", n, m, k)
	for i, v := range xs {
		if i > 0 {
			sb.WriteByte(' ')
		}
		fmt.Fprintf(&sb, "%d", v)
	}
	sb.WriteByte('\n')
	for i, v := range ys {
		if i > 0 {
			sb.WriteByte(' ')
		}
		fmt.Fprintf(&sb, "%d", v)
	}
	sb.WriteByte('\n')
	for _, p := range arr {
		fmt.Fprintf(&sb, "%d %d\n", p.x, p.y)
	}
	return sb.String()
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierD.go /path/to/binary")
		os.Exit(1)
	}
	userBin := os.Args[1]
	oracle, err := buildOracle()
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	defer os.Remove(oracle)
	r := rand.New(rand.NewSource(1))
	for i := 0; i < 100; i++ {
		input := genCase(r)
		want, err := run(oracle, input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "oracle failed on test %d: %v\ninput:\n%s", i+1, err, input)
			os.Exit(1)
		}
		got, err := run(userBin, input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "test %d: %v\ninput:\n%s", i+1, err, input)
			os.Exit(1)
		}
		if want != got {
			fmt.Printf("test %d failed\ninput:\n%sexpected: %s\ngot: %s\n", i+1, input, want, got)
			os.Exit(1)
		}
	}
	fmt.Printf("All %d tests passed\n", 100)
}
