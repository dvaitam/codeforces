package main

import (
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"strings"
)

type Test struct {
	n, k int
	grid []string
}

func (t Test) Input() string {
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("%d %d\n", t.n, t.k))
	for _, row := range t.grid {
		sb.WriteString(row)
		sb.WriteByte('\n')
	}
	return sb.String()
}

func runExe(path, input string) (string, error) {
	cmd := exec.Command(path)
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	err := cmd.Run()
	return out.String(), err
}

func buildRef() (string, error) {
	ref := "./refB.bin"
	cmd := exec.Command("go", "build", "-o", ref, "965B.go")
	if out, err := cmd.CombinedOutput(); err != nil {
		return "", fmt.Errorf("build reference failed: %v: %s", err, string(out))
	}
	return ref, nil
}

func genTests() []Test {
	rand.Seed(2)
	tests := make([]Test, 0, 103)
	for i := 0; i < 100; i++ {
		n := rand.Intn(10) + 1
		k := rand.Intn(n) + 1
		grid := make([]string, n)
		for r := 0; r < n; r++ {
			row := make([]byte, n)
			for c := 0; c < n; c++ {
				ch := '.'
				if rand.Intn(3) == 0 {
					ch = '#'
				}
				row[c] = byte(ch)
			}
			grid[r] = string(row)
		}
		tests = append(tests, Test{n, k, grid})
	}
	tests = append(tests,
		Test{1, 1, []string{"."}},
		Test{3, 2, []string{"...", "...", "..."}},
		Test{2, 2, []string{"##", "##"}},
	)
	return tests
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("Usage: go run verifierB.go /path/to/binary")
		return
	}
	bin := os.Args[1]
	ref, err := buildRef()
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	defer os.Remove(ref)

	tests := genTests()
	for i, tc := range tests {
		input := tc.Input()
		exp, err := runExe(ref, input)
		if err != nil {
			fmt.Printf("reference runtime error on test %d: %v\n", i+1, err)
			os.Exit(1)
		}
		got, err := runExe(bin, input)
		if err != nil {
			fmt.Printf("candidate runtime error on test %d: %v\n", i+1, err)
			os.Exit(1)
		}
		if strings.TrimSpace(exp) != strings.TrimSpace(got) {
			fmt.Printf("Test %d failed\nInput:\n%sExpected:%sGot:%s\n", i+1, input, exp, got)
			os.Exit(1)
		}
	}
	fmt.Println("all tests passed")
}
