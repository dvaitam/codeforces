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

type rect struct{ r, c int }

func buildRef() (string, error) {
	_, cur, _, _ := runtime.Caller(0)
	dir := filepath.Dir(cur)
	src := filepath.Join(dir, "1028A.go")
	bin := filepath.Join(os.TempDir(), "ref1028A.bin")
	cmd := exec.Command("go", "build", "-o", bin, src)
	if out, err := cmd.CombinedOutput(); err != nil {
		return "", fmt.Errorf("build reference failed: %v\n%s", err, string(out))
	}
	return bin, nil
}

func run(bin string, input string) (string, error) {
	cmd := exec.Command(bin)
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	var errBuf bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &errBuf
	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("runtime error: %v\n%s", err, errBuf.String())
	}
	return strings.TrimSpace(out.String()), nil
}

func genCase(r *rand.Rand) string {
	n := r.Intn(10) + 1
	m := r.Intn(10) + 1
	maxSide := n
	if m < maxSide {
		maxSide = m
	}
	if maxSide%2 == 0 {
		maxSide--
	}
	if maxSide <= 0 {
		maxSide = 1
	}
	s := r.Intn((maxSide+1)/2)*2 + 1
	row := r.Intn(n - s + 1)
	col := r.Intn(m - s + 1)
	grid := make([][]byte, n)
	for i := 0; i < n; i++ {
		grid[i] = make([]byte, m)
		for j := 0; j < m; j++ {
			grid[i][j] = 'W'
		}
	}
	for i := 0; i < s; i++ {
		for j := 0; j < s; j++ {
			grid[row+i][col+j] = 'B'
		}
	}
	var sb strings.Builder
	fmt.Fprintf(&sb, "%d %d\n", n, m)
	for i := 0; i < n; i++ {
		sb.WriteString(string(grid[i]))
		sb.WriteByte('\n')
	}
	return sb.String()
}

func main() {
	rand.Seed(time.Now().UnixNano())
	if len(os.Args) < 2 {
		fmt.Println("usage: go run verifierA.go /path/to/binary")
		return
	}
	cand := os.Args[1]
	ref, err := buildRef()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	defer os.Remove(ref)
	for i := 1; i <= 100; i++ {
		in := genCase(rand.New(rand.NewSource(int64(i))))
		want, err := run(ref, in)
		if err != nil {
			fmt.Printf("reference failed on test %d: %v\n", i, err)
			os.Exit(1)
		}
		got, err := run(cand, in)
		if err != nil {
			fmt.Printf("runtime error on test %d: %v\n", i, err)
			os.Exit(1)
		}
		if strings.TrimSpace(want) != strings.TrimSpace(got) {
			fmt.Printf("wrong answer on test %d\ninput:\n%sexpected:%s\ngot:%s\n", i, in, want, got)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
