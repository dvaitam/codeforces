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

type Test struct {
	q  int
	ns []int64
}

func buildRef() (string, error) {
	_, file, _, _ := runtime.Caller(0)
	dir := filepath.Dir(file)
	src := filepath.Join(dir, "818F.go")
	bin := filepath.Join(os.TempDir(), "ref818F.bin")
	cmd := exec.Command("go", "build", "-o", bin, src)
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		return "", err
	}
	return bin, nil
}

func runExe(path, input string) (string, error) {
	cmd := exec.Command(path)
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &bytes.Buffer{}
	err := cmd.Run()
	return strings.TrimSpace(out.String()), err
}

func genTests() []Test {
	rand.Seed(6)
	tests := make([]Test, 100)
	for i := 0; i < 100; i++ {
		q := rand.Intn(3) + 1
		ns := make([]int64, q)
		for j := range ns {
			ns[j] = rand.Int63n(1_000_000_000) + 1
		}
		tests[i] = Test{q, ns}
	}
	return tests
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierF.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	ref, err := buildRef()
	if err != nil {
		fmt.Fprintln(os.Stderr, "failed to build reference:", err)
		os.Exit(1)
	}
	defer os.Remove(ref)

	tests := genTests()
	for i, t := range tests {
		var sb strings.Builder
		sb.WriteString(fmt.Sprintf("%d\n", t.q))
		for _, n := range t.ns {
			sb.WriteString(fmt.Sprintf("%d\n", n))
		}
		input := sb.String()
		want, err := runExe(ref, input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "reference runtime error on test %d: %v\n", i+1, err)
			os.Exit(1)
		}
		got, err := runExe(bin, input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "candidate runtime error on test %d: %v\n", i+1, err)
			os.Exit(1)
		}
		if want != got {
			fmt.Fprintf(os.Stderr, "test %d failed\ninput:%sexpected:%s\nactual:%s\n", i+1, input, want, got)
			os.Exit(1)
		}
	}
	fmt.Println("all tests passed")
}
