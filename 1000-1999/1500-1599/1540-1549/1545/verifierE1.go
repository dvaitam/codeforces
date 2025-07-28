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
	n := rng.Intn(5) + 1
	x := rng.Intn(100) + 1
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("%d\n%d\n", n, x))
	for i := 0; i < n; i++ {
		tl := rng.Intn(100) + 1
		tr := tl + rng.Intn(100)
		l := rng.Intn(100) + 1
		r := l + rng.Intn(100)
		fmt.Fprintf(&sb, "%d %d %d %d\n", tl, tr, l, r)
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
	ref := filepath.Join(dir, "1545E1.go")
	return runProg(ref, input)
}

func main() {
	if len(os.Args) < 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierE1.go <binary>")
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
