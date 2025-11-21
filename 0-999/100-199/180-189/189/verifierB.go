package main

import (
	"bufio"
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"strings"
)

func countRhombi(w, h int64) int64 {
	halfW := w / 2
	ceilHalfW := (w + 1) / 2
	halfH := h / 2
	ceilHalfH := (h + 1) / 2
	return halfW * ceilHalfW * halfH * ceilHalfH
}

func solveRef(input string) string {
	reader := bufio.NewReader(strings.NewReader(input))
	var w, h int64
	if _, err := fmt.Fscan(reader, &w, &h); err != nil {
		return ""
	}
	return fmt.Sprintf("%d", countRhombi(w, h))
}

type testCase struct {
	name   string
	input  string
	expect string
}

func handcraftedTests() []testCase {
	cases := []struct {
		w, h int64
		name string
	}{
		{1, 1, "min_zero"},
		{2, 2, "simple"},
		{3, 3, "odd"},
		{4, 7, "rectangular"},
		{5, 2, "wide"},
		{3999, 4000, "near_limit"},
	}
	var tests []testCase
	for _, c := range cases {
		input := fmt.Sprintf("%d %d\n", c.w, c.h)
		tests = append(tests, testCase{
			name:   c.name,
			input:  input,
			expect: solveRef(input),
		})
	}
	return tests
}

func randomTests() []testCase {
	r := rand.New(rand.NewSource(189))
	var tests []testCase
	for i := 0; i < 300; i++ {
		w := int64(r.Intn(4000) + 1)
		h := int64(r.Intn(4000) + 1)
		input := fmt.Sprintf("%d %d\n", w, h)
		tests = append(tests, testCase{
			name:   fmt.Sprintf("rand_%d", i+1),
			input:  input,
			expect: solveRef(input),
		})
	}
	return tests
}

func runCandidate(bin, input string) (string, error) {
	cmd := exec.Command(bin)
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	err := cmd.Run()
	return strings.TrimSpace(out.String()), err
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("usage: go run verifierB.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	tests := append(handcraftedTests(), randomTests()...)
	for idx, tc := range tests {
		got, err := runCandidate(bin, tc.input)
		if err != nil {
			fmt.Printf("test %d (%s) runtime error: %v\ninput:\n%s", idx+1, tc.name, err, tc.input)
			os.Exit(1)
		}
		if got != tc.expect {
			fmt.Printf("test %d (%s) failed\ninput:\n%s\nexpect:%s\nactual:%s\n", idx+1, tc.name, tc.input, tc.expect, got)
			os.Exit(1)
		}
	}
	fmt.Println("all tests passed")
}
