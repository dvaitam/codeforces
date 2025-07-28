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
	"time"
)

func buildIfGo(path string) (string, func(), error) {
	if strings.HasSuffix(path, ".go") {
		tmp, err := os.CreateTemp("", "cand*")
		if err != nil {
			return "", nil, err
		}
		tmp.Close()
		if out, err := exec.Command("go", "build", "-o", tmp.Name(), path).CombinedOutput(); err != nil {
			os.Remove(tmp.Name())
			return "", nil, fmt.Errorf("build failed: %v\n%s", err, out)
		}
		return tmp.Name(), func() { os.Remove(tmp.Name()) }, nil
	}
	return path, func() {}, nil
}

func buildOracle() (string, func(), error) {
	_, file, _, _ := runtime.Caller(0)
	dir := filepath.Dir(file)
	src := filepath.Join(dir, "1924B.go")
	tmp, err := os.CreateTemp("", "oracle*")
	if err != nil {
		return "", nil, err
	}
	tmp.Close()
	if out, err := exec.Command("go", "build", "-o", tmp.Name(), src).CombinedOutput(); err != nil {
		os.Remove(tmp.Name())
		return "", nil, fmt.Errorf("build oracle failed: %v\n%s", err, out)
	}
	return tmp.Name(), func() { os.Remove(tmp.Name()) }, nil
}

func runProg(exe, input string) (string, error) {
	cmd := exec.Command(exe)
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("runtime error: %v\n%s", err, out.String())
	}
	return strings.TrimSpace(out.String()), nil
}

func generateCase(rng *rand.Rand) string {
	n := rng.Intn(20) + 2
	m := rng.Intn(n-1) + 2
	q := rng.Intn(5) + 1

	used := map[int]bool{1: true, n: true}
	pos := []int{1, n}
	for len(pos) < m {
		x := rng.Intn(n-2) + 2
		if !used[x] {
			used[x] = true
			pos = append(pos, x)
		}
	}
	sort.Ints(pos)
	vals := make([]int, m)
	for i := range vals {
		vals[i] = rng.Intn(100) + 1
	}

	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("%d %d %d\n", n, m, q))
	for i, x := range pos {
		if i > 0 {
			sb.WriteByte(' ')
		}
		sb.WriteString(fmt.Sprintf("%d", x))
	}
	sb.WriteByte('\n')
	for i, v := range vals {
		if i > 0 {
			sb.WriteByte(' ')
		}
		sb.WriteString(fmt.Sprintf("%d", v))
	}
	sb.WriteByte('\n')

	for j := 0; j < q; j++ {
		if j == 0 || rng.Intn(2) == 0 {
			l := rng.Intn(n) + 1
			r := rng.Intn(n-l+1) + l
			sb.WriteString(fmt.Sprintf("2 %d %d\n", l, r))
		} else {
			x := rng.Intn(n-2) + 2
			for used[x] {
				x = rng.Intn(n-2) + 2
			}
			used[x] = true
			v := rng.Intn(100) + 1
			sb.WriteString(fmt.Sprintf("1 %d %d\n", x, v))
		}
	}
	return sb.String()
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierB.go /path/to/binary")
		os.Exit(1)
	}
	cand, cleanup, err := buildIfGo(os.Args[1])
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	defer cleanup()

	oracle, ocleanup, err := buildOracle()
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	defer ocleanup()

	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < 100; i++ {
		input := generateCase(rng)
		exp, err := runProg(oracle, input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "oracle failure on case %d: %v\ninput:\n%s", i+1, err, input)
			os.Exit(1)
		}
		got, err := runProg(cand, input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d: %v\ninput:\n%s", i+1, err, input)
			os.Exit(1)
		}
		if exp != got {
			fmt.Fprintf(os.Stderr, "case %d mismatch\nexpected:\n%s\ngot:\n%s\ninput:\n%s", i+1, exp, got, input)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
