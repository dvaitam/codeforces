package main

import (
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"
)

func compileRef() (string, error) {
	exe := filepath.Join(os.TempDir(), "ref935F")
	cmd := exec.Command("go", "build", "-o", exe, "935F.go")
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		return "", err
	}
	return exe, nil
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

func generateCase(r *rand.Rand) string {
	n := r.Intn(6) + 3
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("%d\n", n))
	for i := 0; i < n; i++ {
		sb.WriteString(fmt.Sprintf("%d", r.Intn(21)))
		if i+1 < n {
			sb.WriteByte(' ')
		}
	}
	sb.WriteByte('\n')
	q := r.Intn(6) + 1
	sb.WriteString(fmt.Sprintf("%d\n", q))
	hasType1 := false
	for i := 0; i < q; i++ {
		t := r.Intn(2) + 1
		if i == q-1 && !hasType1 {
			t = 1
		}
		if t == 1 {
			hasType1 = true
		}
		l := r.Intn(n-2) + 2
		rgt := l + r.Intn(n-l)
		x := r.Intn(21)
		sb.WriteString(fmt.Sprintf("%d %d %d %d\n", t, l, rgt, x))
	}
	return sb.String()
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierF.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	ref, err := compileRef()
	if err != nil {
		fmt.Fprintln(os.Stderr, "failed to build reference:", err)
		os.Exit(1)
	}
	defer os.Remove(ref)

	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < 100; i++ {
		in := generateCase(rng)
		expect, err := runProg(ref, in)
		if err != nil {
			fmt.Fprintf(os.Stderr, "reference runtime error on test %d: %v\n", i+1, err)
			os.Exit(1)
		}
		out, err := runProg(bin, in)
		if err != nil {
			fmt.Fprintf(os.Stderr, "test %d runtime error: %v\n", i+1, err)
			os.Exit(1)
		}
		if strings.TrimSpace(out) != strings.TrimSpace(expect) {
			fmt.Fprintf(os.Stderr, "test %d failed.\nInput:\n%sExpected:\n%s\nGot:\n%s\n", i+1, in, expect, out)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
