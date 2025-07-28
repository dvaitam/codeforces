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
	ref := "./refC.bin"
	cmd := exec.Command("go", "build", "-o", ref, "1922C.go")
	if out, err := cmd.CombinedOutput(); err != nil {
		return "", fmt.Errorf("failed to build reference: %v: %s", err, out)
	}
	return ref, nil
}

func genTests() []string {
	rand.Seed(3)
	tests := make([]string, 100)
	for i := range tests {
		n := rand.Intn(10) + 2
		a := make([]int, n)
		cur := rand.Intn(5)
		a[0] = cur
		prevDiff := rand.Intn(10) + 1
		for j := 1; j < n; j++ {
			diff := rand.Intn(10) + 1
			for diff == prevDiff {
				diff = rand.Intn(10) + 1
			}
			cur += diff
			a[j] = cur
			prevDiff = diff
		}
		m := rand.Intn(5) + 1
		var sb strings.Builder
		sb.WriteString("1\n")
		sb.WriteString(fmt.Sprintf("%d\n", n))
		for j := 0; j < n; j++ {
			if j > 0 {
				sb.WriteByte(' ')
			}
			sb.WriteString(fmt.Sprintf("%d", a[j]))
		}
		sb.WriteByte('\n')
		sb.WriteString(fmt.Sprintf("%d\n", m))
		for j := 0; j < m; j++ {
			x := rand.Intn(n) + 1
			y := rand.Intn(n) + 1
			for y == x {
				y = rand.Intn(n) + 1
			}
			sb.WriteString(fmt.Sprintf("%d %d\n", x, y))
		}
		tests[i] = sb.String()
	}
	return tests
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierC.go /path/to/binary")
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
