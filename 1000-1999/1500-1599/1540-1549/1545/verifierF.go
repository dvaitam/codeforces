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

func generateCase(rng *rand.Rand) string {
	n := rng.Intn(10) + 1
	m := rng.Intn(20) + 1
	a := make([]int, n)
	b := make([]int, n)
	c := make([]int, n)
	for i := 0; i < n; i++ {
		a[i] = rng.Intn(n) + 1
		b[i] = rng.Intn(n) + 1
		c[i] = rng.Intn(n) + 1
	}
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("%d %d\n", n, m))
	for i := 0; i < n; i++ {
		if i > 0 {
			sb.WriteByte(' ')
		}
		fmt.Fprintf(&sb, "%d", a[i])
	}
	sb.WriteByte('\n')
	for i := 0; i < n; i++ {
		if i > 0 {
			sb.WriteByte(' ')
		}
		fmt.Fprintf(&sb, "%d", b[i])
	}
	sb.WriteByte('\n')
	for i := 0; i < n; i++ {
		if i > 0 {
			sb.WriteByte(' ')
		}
		fmt.Fprintf(&sb, "%d", c[i])
	}
	sb.WriteByte('\n')
	secondExists := false
	for op := 0; op < m; op++ {
		typ := rng.Intn(2) + 1
		if op == m-1 && !secondExists {
			typ = 2
		}
		if typ == 1 {
			k := rng.Intn(n) + 1
			x := rng.Intn(n) + 1
			fmt.Fprintf(&sb, "1 %d %d\n", k, x)
		} else {
			r := rng.Intn(n) + 1
			fmt.Fprintf(&sb, "2 %d\n", r)
			secondExists = true
		}
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
	ref := filepath.Join(dir, "1545F.go")
	return runProg(ref, input)
}

func main() {
	if len(os.Args) < 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierF.go <binary>")
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
