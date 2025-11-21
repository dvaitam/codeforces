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

const (
	refSource2171H = "2171H.go"
	refBinary2171H = "ref2171H.bin"
	totalTests     = 250
)

type testCase struct {
	g int
	c int
	l int
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("usage: go run verifierH.go /path/to/binary")
		return
	}
	candidate := os.Args[1]

	ref, err := buildReference()
	if err != nil {
		fmt.Println("failed to build reference:", err)
		return
	}
	defer os.Remove(ref)

	tests := generateTests()
	for i, tc := range tests {
		input := fmt.Sprintf("%d %d %d\n", tc.g, tc.c, tc.l)

		refOut, err := runProgram(ref, []byte(input))
		if err != nil {
			fmt.Printf("reference runtime error on test %d: %v\n", i+1, err)
			return
		}
		candOut, err := runProgram(candidate, []byte(input))
		if err != nil {
			fmt.Printf("candidate runtime error on test %d: %v\n", i+1, err)
			return
		}

		if normalize(refOut) != normalize(candOut) {
			fmt.Printf("Mismatch on test %d\ninput: %sexpected: %s\ngot: %s\n", i+1, input, refOut, candOut)
			return
		}
	}
	fmt.Printf("All %d tests passed\n", len(tests))
}

func buildReference() (string, error) {
	cmd := exec.Command("go", "build", "-o", refBinary2171H, refSource2171H)
	if out, err := cmd.CombinedOutput(); err != nil {
		return "", fmt.Errorf("%v\n%s", err, string(out))
	}
	return filepath.Join(".", refBinary2171H), nil
}

func runProgram(path string, input []byte) (string, error) {
	cmd := exec.Command(path)
	cmd.Stdin = bytes.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = os.Stderr
	err := cmd.Run()
	return out.String(), err
}

func normalize(s string) string {
	return strings.TrimSpace(strings.ToLower(s))
}

func generateTests() []testCase {
	rnd := rand.New(rand.NewSource(2171))
	tests := []testCase{
		{80, 80, 80},
		{100, 100, 100},
		{80, 90, 100},
		{95, 90, 85},
		{88, 94, 95},
		{100, 80, 81},
		{98, 99, 98},
		{95, 86, 85},
	}
	for len(tests) < totalTests {
		tests = append(tests, testCase{
			g: rnd.Intn(21) + 80,
			c: rnd.Intn(21) + 80,
			l: rnd.Intn(21) + 80,
		})
	}
	return tests
}
