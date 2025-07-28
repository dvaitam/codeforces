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

type matrix [][]int

func latinSquare(n int) matrix {
	m := make(matrix, n)
	for i := 0; i < n; i++ {
		m[i] = make([]int, n)
		for j := 0; j < n; j++ {
			m[i][j] = (i+j)%n + 1
		}
	}
	return m
}

func generateCase(rng *rand.Rand) string {
	n := rng.Intn(6) + 5 // 5..10
	base := latinSquare(n)
	used := make(map[string]bool)
	var extras matrix
	for len(extras) < n {
		idx := len(extras)
		arr := append([]int(nil), base[idx]...)
		p := rng.Intn(n)
		q := rng.Intn(n)
		for q == p {
			q = rng.Intn(n)
		}
		arr[p], arr[q] = arr[q], arr[p]
		key := fmt.Sprint(arr)
		if key == fmt.Sprint(base[idx]) || used[key] {
			continue
		}
		used[key] = true
		extras = append(extras, arr)
	}
	var sb strings.Builder
	sb.WriteString("1\n")
	sb.WriteString(fmt.Sprintf("%d\n", n))
	for i := 0; i < n; i++ {
		for j := 0; j < n; j++ {
			if j > 0 {
				sb.WriteByte(' ')
			}
			fmt.Fprintf(&sb, "%d", base[i][j])
		}
		sb.WriteByte('\n')
	}
	for i := 0; i < n; i++ {
		for j := 0; j < n; j++ {
			if j > 0 {
				sb.WriteByte(' ')
			}
			fmt.Fprintf(&sb, "%d", extras[i][j])
		}
		sb.WriteByte('\n')
	}
	return sb.String()
}

func runProg(path, input string) (string, error) {
	var cmd *exec.Cmd
	if strings.HasSuffix(path, ".go") {
		cmd = exec.Command("go", "run", path)
	} else {
		cmd = exec.Command(path)
	}
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	err := cmd.Run()
	return strings.TrimSpace(out.String()), err
}

func runRef(input string) (string, error) {
	_, self, _, _ := runtime.Caller(0)
	dir := filepath.Dir(self)
	ref := filepath.Join(dir, "1545C.go")
	return runProg(ref, input)
}

func main() {
	if len(os.Args) < 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierC.go <binary>")
		os.Exit(1)
	}
	bin := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < 100; i++ {
		input := generateCase(rng)
		expect, err := runRef(input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "reference failed on test %d: %v\n", i+1, err)
			os.Exit(1)
		}
		got, err := runProg(bin, input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "candidate failed on test %d: %v\n", i+1, err)
			os.Exit(1)
		}
		if got != expect {
			fmt.Printf("test %d failed\ninput:\n%sexpected:%s\nactual:%s\n", i+1, input, expect, got)
			os.Exit(1)
		}
	}
	fmt.Println("all tests passed")
}
