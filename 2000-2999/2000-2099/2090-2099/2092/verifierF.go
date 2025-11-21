package main

import (
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"strings"
	"time"
)

type testCase struct {
	n int
	s string
}

func run(bin string, input []byte) (string, error) {
	cmd := exec.Command(bin)
	cmd.Stdin = bytes.NewReader(input)
	out, err := cmd.CombinedOutput()
	return string(out), err
}

func parseOutput(out string, cases []testCase) ([][]int, error) {
	reader := strings.NewReader(out)
	results := make([][]int, 0, len(cases))
	for idx, tc := range cases {
		ans := make([]int, tc.n)
		for i := 0; i < tc.n; i++ {
			if _, err := fmt.Fscan(reader, &ans[i]); err != nil {
				return nil, fmt.Errorf("output ended early on test %d prefix %d: %v", idx+1, i+1, err)
			}
		}
		results = append(results, ans)
	}
	// ensure no extra tokens (ignore trailing whitespace)
	var extra string
	if _, err := fmt.Fscan(reader, &extra); err == nil {
		return nil, fmt.Errorf("extra output detected after all answers")
	}
	return results, nil
}

func genRandomString(n int) string {
	var sb strings.Builder
	sb.Grow(n)
	for i := 0; i < n; i++ {
		if rand.Intn(2) == 0 {
			sb.WriteByte('0')
		} else {
			sb.WriteByte('1')
		}
	}
	return sb.String()
}

func genTestCases() []testCase {
	t := rand.Intn(5) + 1
	// Occasionally stress with one huge case near limit.
	if rand.Intn(6) == 0 {
		t = 1
	}
	cases := make([]testCase, 0, t)
	remaining := 1_000_000
	for i := 0; i < t; i++ {
		if remaining <= 0 {
			break
		}
		mode := rand.Intn(8)
		n := 0
		s := ""
		switch mode {
		case 0:
			n = 1
			s = "0"
		case 1:
			n = rand.Intn(20) + 1
			s = strings.Repeat("0", n)
		case 2:
			n = rand.Intn(20) + 1
			s = strings.Repeat("1", n)
		case 3:
			n = rand.Intn(50) + 1
			var sb strings.Builder
			sb.Grow(n)
			for j := 0; j < n; j++ {
				if j%2 == 0 {
					sb.WriteByte('0')
				} else {
					sb.WriteByte('1')
				}
			}
			s = sb.String()
		case 4:
			n = rand.Intn(1000) + 1
			s = genRandomString(n)
		case 5:
			// Large random
			n = rand.Intn(50000) + 50000
			s = genRandomString(n)
		default:
			// Near maximum
			n = rand.Intn(200000) + 200000
			if rand.Intn(20) == 0 {
				n = 1_000_000
			}
			s = genRandomString(n)
		}
		if n > remaining {
			n = remaining
			s = s[:n]
		}
		remaining -= n
		cases = append(cases, testCase{n: n, s: s})
		if remaining == 0 {
			break
		}
	}
	return cases
}

func buildInput(cases []testCase) []byte {
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("%d\n", len(cases)))
	for _, tc := range cases {
		sb.WriteString(fmt.Sprintf("%d\n", tc.n))
		sb.WriteString(tc.s)
		sb.WriteByte('\n')
	}
	return []byte(sb.String())
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierF.go <binary>")
		os.Exit(1)
	}
	cand := os.Args[1]
	ref := "./refF.bin"
	if err := exec.Command("go", "build", "-o", ref, "2092F.go").Run(); err != nil {
		fmt.Println("failed to build reference solution:", err)
		os.Exit(1)
	}
	defer os.Remove(ref)

	rand.Seed(time.Now().UnixNano())

	for iter := 0; iter < 200; iter++ {
		cases := genTestCases()
		input := buildInput(cases)

		refOut, err := run(ref, input)
		if err != nil {
			fmt.Println("reference failed on iteration", iter+1, ":", err)
			fmt.Println("input:\n", string(input))
			os.Exit(1)
		}
		candOut, err := run(cand, input)
		if err != nil {
			fmt.Printf("candidate runtime error on iteration %d: %v\n", iter+1, err)
			fmt.Println("input:\n", string(input))
			os.Exit(1)
		}

		refParsed, err := parseOutput(refOut, cases)
		if err != nil {
			fmt.Println("failed to parse reference output on iteration", iter+1, ":", err)
			fmt.Println("output:\n", refOut)
			fmt.Println("input:\n", string(input))
			os.Exit(1)
		}
		candParsed, err := parseOutput(candOut, cases)
		if err != nil {
			fmt.Println("failed to parse candidate output on iteration", iter+1, ":", err)
			fmt.Println("output:\n", candOut)
			fmt.Println("input:\n", string(input))
			os.Exit(1)
		}

		if len(refParsed) != len(candParsed) {
			fmt.Printf("test count mismatch on iteration %d\n", iter+1)
			fmt.Println("input:\n", string(input))
			fmt.Println("reference:\n", refOut)
			fmt.Println("candidate:\n", candOut)
			os.Exit(1)
		}
		match := true
		for i := range refParsed {
			if len(refParsed[i]) != len(candParsed[i]) {
				match = false
				break
			}
			for j := range refParsed[i] {
				if refParsed[i][j] != candParsed[i][j] {
					match = false
					break
				}
			}
			if !match {
				break
			}
		}
		if !match {
			fmt.Printf("mismatch on iteration %d\n", iter+1)
			fmt.Println("input:\n", string(input))
			fmt.Println("reference:\n", refOut)
			fmt.Println("candidate:\n", candOut)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed.")
}
