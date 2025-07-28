package main

import (
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"strings"
)

func run(path, input string) (string, error) {
	cmd := exec.Command(path)
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = os.Stderr
	err := cmd.Run()
	return strings.TrimSpace(out.String()), err
}

func buildRef() (string, error) {
	ref := "./refH.bin"
	cmd := exec.Command("go", "build", "-o", ref, "1844H.go")
	if out, err := cmd.CombinedOutput(); err != nil {
		return "", fmt.Errorf("failed to build reference: %v: %s", err, out)
	}
	return ref, nil
}

func genTests() []string {
	rand.Seed(9)
	tests := make([]string, 100)
	for i := range tests {
		n := (rand.Intn(5) + 1) * 3
		xs := rand.Perm(n)
		ys := rand.Perm(n)
		var sb strings.Builder
		sb.WriteString(fmt.Sprintf("%d\n", n))
		for j := 0; j < n; j++ {
			fmt.Fprintf(&sb, "%d %d\n", xs[j]+1, ys[j]+1)
		}
		tests[i] = sb.String()
	}
	return tests
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierH.go /path/to/binary")
		return
	}
	candidate := os.Args[1]
	ref, err := buildRef()
	if err != nil {
		fmt.Println(err)
		return
	}
	defer os.Remove(ref)

	tests := genTests()
	for i, t := range tests {
		exp, err := run(ref, t)
		if err != nil {
			fmt.Printf("reference runtime error on test %d: %v\n", i+1, err)
			return
		}
		out, err := run(candidate, t)
		if err != nil {
			fmt.Printf("candidate runtime error on test %d: %v\n", i+1, err)
			return
		}
		if exp != out {
			fmt.Printf("wrong answer on test %d\nexpected: %s\ngot: %s\n", i+1, exp, out)
			return
		}
	}
	fmt.Println("All tests passed!")
}
