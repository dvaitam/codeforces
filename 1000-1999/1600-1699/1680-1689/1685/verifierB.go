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

func runBinary(path, input string) (string, error) {
	cmd := exec.Command(path)
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	err := cmd.Run()
	return out.String(), err
}

func generateTests() []string {
	r := rand.New(rand.NewSource(2))
	tests := make([]string, 0, 108)

	// Fixed test cases from problem examples and edge cases
	fixed := []string{
		"1\n1 0 0 0\nB\n",
		"1\n0 0 1 0\nAB\n",
		"1\n1 1 0 1\nABAB\n",
		"1\n1 0 1 1\nABAAB\n",
		"1\n1 1 2 2\nBAABBABBAA\n",
		"1\n1 1 2 3\nABABABBAABAB\n",
		"1\n2 3 5 4\nAABAABBABAAABABBABBBABB\n",
		"1\n1 3 3 10\nBBABABABABBBABABABABABABAABABA\n",
	}
	tests = append(tests, fixed...)

	for len(tests) < 108 {
		n := r.Intn(10) + 1 // length of string
		c := r.Intn(n/2 + 1)
		d := r.Intn((n-2*c)/2 + 1)
		if 2*(c+d) > n {
			continue
		}
		leftover := n - 2*(c+d)
		a := r.Intn(leftover + 1)
		b := leftover - a
		numA := a + c + d
		numB := b + c + d
		arr := make([]byte, n)
		for i := 0; i < numA; i++ {
			arr[i] = 'A'
		}
		for i := 0; i < numB; i++ {
			arr[numA+i] = 'B'
		}
		r.Shuffle(n, func(i, j int) { arr[i], arr[j] = arr[j], arr[i] })
		s := string(arr)
		tests = append(tests, fmt.Sprintf("1\n%d %d %d %d\n%s\n", a, b, c, d, s))
	}
	return tests
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierB.go /path/to/binary")
		os.Exit(1)
	}
	candidate := os.Args[1]
	_, file, _, _ := runtime.Caller(0)
	dir := filepath.Dir(file)
	ref := filepath.Join(dir, "refB")
	cmd := exec.Command("go", "build", "-o", ref, filepath.Join(dir, "1685B.go"))
	if out, err := cmd.CombinedOutput(); err != nil {
		fmt.Fprintf(os.Stderr, "failed to build reference: %v\n%s", err, out)
		os.Exit(1)
	}
	defer os.Remove(ref)

	tests := generateTests()
	for i, input := range tests {
		candOut, cErr := runBinary(candidate, input)
		refOut, rErr := runBinary(ref, input)
		if cErr != nil {
			fmt.Fprintf(os.Stderr, "case %d: candidate error: %v\n", i+1, cErr)
			os.Exit(1)
		}
		if rErr != nil {
			fmt.Fprintf(os.Stderr, "case %d: reference error: %v\n", i+1, rErr)
			os.Exit(1)
		}
		if strings.TrimSpace(candOut) != strings.TrimSpace(refOut) {
			fmt.Fprintf(os.Stderr, "case %d failed\ninput:\n%sexpected:%sactual:%s", i+1, input, refOut, candOut)
			os.Exit(1)
		}
	}
	fmt.Printf("All %d tests passed\n", len(tests))
}
