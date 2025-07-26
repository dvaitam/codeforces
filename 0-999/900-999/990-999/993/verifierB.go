package main

import (
	"bytes"
	"fmt"
	"io"
	"math/rand"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"
)

func runBinary(path string, input string) (string, error) {
	cmd := exec.Command(path)
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = io.Discard
	if err := cmd.Run(); err != nil {
		return "", err
	}
	return strings.TrimSpace(out.String()), nil
}

func compileReference() (string, error) {
	refPath := filepath.Join(os.TempDir(), fmt.Sprintf("refB_%d", time.Now().UnixNano()))
	cmd := exec.Command("go", "build", "-o", refPath, "993B.go")
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		return "", err
	}
	return refPath, nil
}

func generatePair(rng *rand.Rand) [2]int {
	a := rng.Intn(9) + 1
	b := rng.Intn(9) + 1
	for b == a {
		b = rng.Intn(9) + 1
	}
	return [2]int{a, b}
}

func sharedExactlyOne(A, B [][2]int) bool {
	for _, a := range A {
		for _, b := range B {
			c := 0
			if a[0] == b[0] || a[0] == b[1] {
				c++
			}
			if a[1] == b[0] || a[1] == b[1] {
				c++
			}
			if c == 1 {
				return true
			}
		}
	}
	return false
}

func generateTests() []string {
	rng := rand.New(rand.NewSource(1))
	tests := make([]string, 100)
	for i := 0; i < 100; i++ {
		var A, B [][2]int
		n := rng.Intn(4) + 1
		m := rng.Intn(4) + 1
		for {
			A = make([][2]int, 0, n)
			B = make([][2]int, 0, m)
			seenA := map[[2]int]struct{}{}
			seenB := map[[2]int]struct{}{}
			for len(A) < n {
				p := generatePair(rng)
				if _, ok := seenA[p]; !ok {
					seenA[p] = struct{}{}
					A = append(A, p)
				}
			}
			for len(B) < m {
				p := generatePair(rng)
				if _, ok := seenB[p]; !ok {
					seenB[p] = struct{}{}
					B = append(B, p)
				}
			}
			if sharedExactlyOne(A, B) {
				break
			}
		}
		var sb strings.Builder
		sb.WriteString(fmt.Sprintf("%d %d\n", n, m))
		for i, p := range A {
			if i > 0 {
				sb.WriteByte(' ')
			}
			sb.WriteString(fmt.Sprintf("%d %d", p[0], p[1]))
		}
		sb.WriteByte('\n')
		for i, p := range B {
			if i > 0 {
				sb.WriteByte(' ')
			}
			sb.WriteString(fmt.Sprintf("%d %d", p[0], p[1]))
		}
		sb.WriteByte('\n')
		tests[i] = sb.String()
	}
	return tests
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierB.go /path/to/binary")
		return
	}
	candidate := os.Args[1]
	ref, err := compileReference()
	if err != nil {
		fmt.Println("failed to compile reference:", err)
		return
	}
	defer os.Remove(ref)
	tests := generateTests()
	for i, t := range tests {
		exp, err := runBinary(ref, t)
		if err != nil {
			fmt.Printf("reference failed on test %d: %v\n", i+1, err)
			return
		}
		out, err := runBinary(candidate, t)
		if err != nil {
			fmt.Printf("tested binary failed on test %d: %v\n", i+1, err)
			return
		}
		if strings.TrimSpace(out) != strings.TrimSpace(exp) {
			fmt.Printf("Wrong answer on test %d\nInput:\n%sExpected: %s\nGot: %s\n", i+1, t, exp, out)
			return
		}
	}
	fmt.Println("All tests passed!")
}
