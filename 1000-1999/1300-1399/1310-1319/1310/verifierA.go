package main

import (
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

func buildRef() (string, error) {
	tmp := filepath.Join(os.TempDir(), "refA")
	cmd := exec.Command("go", "build", "-o", tmp, "1310A.go")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		return "", err
	}
	return tmp, nil
}

func runProg(path string, input string) (string, error) {
	cmd := exec.Command(path)
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = os.Stderr
	err := cmd.Run()
	return strings.TrimSpace(out.String()), err
}

func genTest(rng *rand.Rand) string {
	n := rng.Intn(5) + 1
	a := make([]int, n)
	t := make([]int, n)
	for i := 0; i < n; i++ {
		a[i] = rng.Intn(5) + 1
		t[i] = rng.Intn(5) + 1
	}
	var buf bytes.Buffer
	fmt.Fprintln(&buf, n)
	for i, v := range a {
		if i > 0 {
			fmt.Fprint(&buf, " ")
		}
		fmt.Fprint(&buf, v)
	}
	fmt.Fprintln(&buf)
	for i, v := range t {
		if i > 0 {
			fmt.Fprint(&buf, " ")
		}
		fmt.Fprint(&buf, v)
	}
	fmt.Fprintln(&buf)
	return buf.String()
}

func main() {
	if len(os.Args) < 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierA.go /path/to/bin")
		os.Exit(1)
	}
	target := os.Args[1]
	ref, err := buildRef()
	if err != nil {
		fmt.Fprintln(os.Stderr, "failed to build reference:", err)
		os.Exit(1)
	}
	defer os.Remove(ref)

	rng := rand.New(rand.NewSource(42))
	for i := 0; i < 100; i++ {
		test := genTest(rng)
		expected, err1 := runProg(ref, test)
		if err1 != nil {
			fmt.Fprintln(os.Stderr, "reference run error:", err1)
			os.Exit(1)
		}
		got, err2 := runProg(target, test)
		if err2 != nil {
			fmt.Fprintf(os.Stderr, "test %d: execution error: %v\n", i+1, err2)
			os.Exit(1)
		}
		if expected != got {
			fmt.Printf("test %d failed\ninput:\n%s\nexpected:%s\ngot:%s\n", i+1, test, expected, got)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
