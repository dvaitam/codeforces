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
	ref := filepath.Join(dir, "683H.go")
	return runProg(ref, input)
}

func randPermutation(n int) []int {
	p := rand.Perm(n)
	// ensure no fixed points
	for i := 0; i < n; i++ {
		if p[i] == i {
			j := (i + 1) % n
			p[i], p[j] = p[j], p[i]
		}
	}
	return p
}

func genCase() string {
	n := rand.Intn(9) + 2
	k := rand.Int63n(50) + 1
	perm := randPermutation(n)
	a := make([]int, n)
	for j, v := range perm {
		a[v] = j + 1
	}
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("%d %d\n", n, k))
	for i := 0; i < n; i++ {
		sb.WriteString(fmt.Sprintf("%d", a[i]))
		if i+1 < n {
			sb.WriteByte(' ')
		}
	}
	sb.WriteByte('\n')
	return sb.String()
}

func main() {
	if len(os.Args) < 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierH.go <binary>")
		os.Exit(1)
	}
	bin := os.Args[1]
	rand.Seed(time.Now().UnixNano())
	for i := 0; i < 100; i++ {
		in := genCase()
		expect, err := runRef(in)
		if err != nil {
			fmt.Fprintf(os.Stderr, "reference failed on test %d: %v\n", i+1, err)
			os.Exit(1)
		}
		got, err := runProg(bin, in)
		if err != nil {
			fmt.Fprintf(os.Stderr, "candidate failed on test %d: %v\n", i+1, err)
			os.Exit(1)
		}
		if strings.TrimSpace(got) != strings.TrimSpace(expect) {
			fmt.Printf("test %d failed\ninput:%sexpected:%s\nactual:%s\n", i+1, in, expect, got)
			os.Exit(1)
		}
	}
	fmt.Println("all tests passed")
}
