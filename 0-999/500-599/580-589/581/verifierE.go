package main

import (
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"strings"
)

func run(bin string, input string) (string, error) {
	cmd := exec.Command(bin)
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	err := cmd.Run()
	return strings.TrimSpace(out.String()), err
}

func gen() string {
	e := rand.Intn(1000) + 1
	s := rand.Intn(100) + 1
	n := rand.Intn(5) + 1
	m := rand.Intn(5) + 1
	var b strings.Builder
	fmt.Fprintf(&b, "%d %d %d %d\n", e, s, n, m)
	used := make(map[int]bool)
	for i := 0; i < n; i++ {
		var x int
		for {
			x = rand.Intn(2*e) - e
			if !used[x] {
				used[x] = true
				break
			}
		}
		t := rand.Intn(3) + 1
		fmt.Fprintf(&b, "%d %d\n", t, x)
	}
	for i := 0; i < m; i++ {
		if i > 0 {
			b.WriteByte(' ')
		}
		f := rand.Intn(e*2) - e
		if f >= e {
			f = e - 1
		}
		fmt.Fprintf(&b, "%d", f)
	}
	b.WriteByte('\n')
	return b.String()
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierE.go /path/to/binary")
		os.Exit(1)
	}
	candidate := os.Args[1]
	ref := "./_refE.bin"
	if err := exec.Command("go", "build", "-o", ref, "581E.go").Run(); err != nil {
		fmt.Fprintln(os.Stderr, "failed to build reference:", err)
		os.Exit(1)
	}
	defer os.Remove(ref)

	rand.Seed(1)
	for i := 1; i <= 100; i++ {
		input := gen()
		exp, err := run(ref, input)
		if err != nil {
			fmt.Fprintln(os.Stderr, "reference runtime error:", err)
			os.Exit(1)
		}
		got, err := run(candidate, input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "candidate runtime error on test %d: %v\n", i, err)
			os.Exit(1)
		}
		if exp != got {
			fmt.Printf("Test %d failed\nInput:%sExpected:%s\nGot:%s\n", i, input, exp, got)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
