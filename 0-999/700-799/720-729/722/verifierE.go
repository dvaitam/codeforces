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

type testCaseE struct {
	n, m  int
	k     int
	s     int
	cells [][2]int
}

func genTestsE() []testCaseE {
	rand.Seed(46)
	tests := make([]testCaseE, 100)
	for i := range tests {
		n := rand.Intn(4) + 1
		m := rand.Intn(4) + 1
		maxCells := n * m
		k := rand.Intn(maxCells + 1)
		s := rand.Intn(100) + 1
		used := make(map[[2]int]bool)
		cells := make([][2]int, 0, k)
		for j := 0; j < k; j++ {
			var r, c int
			for {
				r = rand.Intn(n) + 1
				c = rand.Intn(m) + 1
				if !used[[2]int{r, c}] {
					used[[2]int{r, c}] = true
					break
				}
			}
			cells = append(cells, [2]int{r, c})
		}
		tests[i] = testCaseE{n, m, k, s, cells}
	}
	return tests
}

func compileRef() (string, error) {
	exe := filepath.Join(os.TempDir(), "722E_ref")
	cmd := exec.Command("go", "build", "-o", exe, "722E.go")
	cmd.Dir = "."
	out, err := cmd.CombinedOutput()
	if err != nil {
		return "", fmt.Errorf("compile reference failed: %v\n%s", err, out)
	}
	return exe, nil
}

func runBinary(bin, input string) (string, error) {
	cmd := exec.Command(bin)
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	if err := cmd.Run(); err != nil {
		return out.String(), fmt.Errorf("runtime error: %v", err)
	}
	return out.String(), nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("Usage: go run verifierE.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	ref, err := compileRef()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	defer os.Remove(ref)
	tests := genTestsE()
	for i, tc := range tests {
		var sb strings.Builder
		fmt.Fprintf(&sb, "%d %d %d %d\n", tc.n, tc.m, tc.k, tc.s)
		for _, p := range tc.cells {
			fmt.Fprintf(&sb, "%d %d\n", p[0], p[1])
		}
		input := sb.String()
		exp, err := runBinary(ref, input)
		if err != nil {
			fmt.Printf("reference failed on test %d: %v\n", i+1, err)
			os.Exit(1)
		}
		got, err := runBinary(bin, input)
		if err != nil {
			fmt.Printf("test %d: %v\noutput:\n%s", i+1, err, got)
			os.Exit(1)
		}
		if strings.TrimSpace(got) != strings.TrimSpace(exp) {
			fmt.Printf("test %d failed: expected %q got %q\ninput:\n%s", i+1, strings.TrimSpace(exp), strings.TrimSpace(got), input)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
