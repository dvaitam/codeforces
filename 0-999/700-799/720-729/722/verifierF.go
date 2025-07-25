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

type testCaseF struct {
	n, m int
	seqs [][]int
}

func genTestsF() []testCaseF {
	rand.Seed(47)
	tests := make([]testCaseF, 100)
	for i := range tests {
		n := rand.Intn(3) + 1
		m := rand.Intn(10) + 1
		seqs := make([][]int, n)
		for j := 0; j < n; j++ {
			k := rand.Intn(3) + 1
			seq := make([]int, k)
			used := make(map[int]bool)
			for x := 0; x < k; x++ {
				v := rand.Intn(m) + 1
				for used[v] {
					v = rand.Intn(m) + 1
				}
				used[v] = true
				seq[x] = v
			}
			seqs[j] = seq
		}
		tests[i] = testCaseF{n, m, seqs}
	}
	return tests
}

func compileRef() (string, error) {
	exe := filepath.Join(os.TempDir(), "722F_ref")
	cmd := exec.Command("go", "build", "-o", exe, "722F.go")
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
		fmt.Println("Usage: go run verifierF.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	ref, err := compileRef()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	defer os.Remove(ref)
	tests := genTestsF()
	for i, tc := range tests {
		var sb strings.Builder
		fmt.Fprintf(&sb, "%d %d\n", tc.n, tc.m)
		for _, seq := range tc.seqs {
			fmt.Fprintf(&sb, "%d", len(seq))
			for _, v := range seq {
				fmt.Fprintf(&sb, " %d", v)
			}
			sb.WriteByte('\n')
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
