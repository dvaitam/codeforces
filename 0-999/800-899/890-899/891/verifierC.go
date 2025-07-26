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
)

func buildReference() (string, error) {
	_, file, _, _ := runtime.Caller(0)
	dir := filepath.Dir(file)
	src := filepath.Join(dir, "891C.go")
	ref := filepath.Join(os.TempDir(), "ref891C")
	cmd := exec.Command("go", "build", "-o", ref, src)
	if out, err := cmd.CombinedOutput(); err != nil {
		return "", fmt.Errorf("failed to build reference: %v\n%s", err, out)
	}
	return ref, nil
}

func runBinary(bin, input string) (string, error) {
	cmd := exec.Command(bin)
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = new(bytes.Buffer)
	err := cmd.Run()
	return strings.TrimSpace(out.String()), err
}

func generateTests() []string {
	rand.Seed(3)
	tests := make([]string, 100)
	for t := 0; t < 100; t++ {
		n := rand.Intn(4) + 2
		m := n - 1 + rand.Intn(3)
		edges := make([][3]int, m)
		for i := 1; i < n; i++ {
			p := rand.Intn(i)
			w := rand.Intn(5) + 1
			edges[i-1] = [3]int{p, i, w}
		}
		for i := n - 1; i < m; i++ {
			u := rand.Intn(n)
			v := rand.Intn(n)
			if u == v {
				v = (v + 1) % n
			}
			w := rand.Intn(5) + 1
			edges[i] = [3]int{u, v, w}
		}
		q := rand.Intn(3) + 1
		var sb strings.Builder
		fmt.Fprintf(&sb, "%d %d\n", n, m)
		for i := 0; i < m; i++ {
			fmt.Fprintf(&sb, "%d %d %d\n", edges[i][0]+1, edges[i][1]+1, edges[i][2])
		}
		fmt.Fprintf(&sb, "%d\n", q)
		for qi := 0; qi < q; qi++ {
			k := rand.Intn(m) + 1
			idx := rand.Perm(m)[:k]
			fmt.Fprintf(&sb, "%d", k)
			for _, id := range idx {
				fmt.Fprintf(&sb, " %d", id+1)
			}
			sb.WriteByte('\n')
		}
		tests[t] = sb.String()
	}
	return tests
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierC.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	ref, err := buildReference()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	defer os.Remove(ref)

	tests := generateTests()
	for i, t := range tests {
		exp, err := runBinary(ref, t)
		if err != nil {
			fmt.Printf("reference failed on test %d: %v\n", i+1, err)
			os.Exit(1)
		}
		got, err := runBinary(bin, t)
		if err != nil {
			fmt.Printf("Test %d: runtime error: %v\n", i+1, err)
			os.Exit(1)
		}
		if got != exp {
			fmt.Printf("Test %d failed:\ninput:\n%sexpected:%s\ngot:%s\n", i+1, t, exp, got)
			os.Exit(1)
		}
	}
	fmt.Printf("All %d tests passed.\n", len(tests))
}
