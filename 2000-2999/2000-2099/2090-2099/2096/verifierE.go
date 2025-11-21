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

func parseOutputs(out string, t int) ([]string, error) {
	reader := strings.NewReader(out)
	res := make([]string, 0, t)
	for i := 0; i < t; i++ {
		var tok string
		if _, err := fmt.Fscan(reader, &tok); err != nil {
			return nil, fmt.Errorf("output ended early at case %d: %v", i+1, err)
		}
		res = append(res, tok)
	}
	// ensure no extra tokens beyond whitespace
	var extra string
	if _, err := fmt.Fscan(reader, &extra); err == nil {
		return nil, fmt.Errorf("extra output after %d cases", t)
	}
	return res, nil
}

func genCase(maxRemain int) testCase {
	mode := rand.Intn(6)
	n := 0
	s := ""
	switch mode {
	case 0:
		n = 3
		s = "BBB"
	case 1:
		n = rand.Intn(8) + 3
		s = strings.Repeat("P", n)
	case 2:
		n = rand.Intn(8) + 3
		s = strings.Repeat("B", n)
	case 3:
		n = rand.Intn(50) + 3
		var sb strings.Builder
		sb.Grow(n)
		for i := 0; i < n; i++ {
			if i%2 == 0 {
				sb.WriteByte('B')
			} else {
				sb.WriteByte('P')
			}
		}
		s = sb.String()
	default:
		n = rand.Intn(50000) + 3
		if n > maxRemain {
			n = maxRemain
		}
		var sb strings.Builder
		sb.Grow(n)
		for i := 0; i < n; i++ {
			if rand.Intn(2) == 0 {
				sb.WriteByte('B')
			} else {
				sb.WriteByte('P')
			}
		}
		s = sb.String()
	}
	if n > maxRemain {
		n = maxRemain
		s = s[:n]
	}
	return testCase{n: n, s: s}
}

func buildInput() ([]byte, []testCase) {
	remaining := 200000
	t := rand.Intn(6) + 1
	cases := make([]testCase, 0, t)
	for i := 0; i < t; i++ {
		if remaining <= 0 {
			break
		}
		tc := genCase(remaining)
		if tc.n == 0 {
			continue
		}
		remaining -= tc.n
		cases = append(cases, tc)
	}
	if len(cases) == 0 {
		cases = append(cases, testCase{n: 3, s: "BPB"})
	}

	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("%d\n", len(cases)))
	for _, tc := range cases {
		sb.WriteString(fmt.Sprintf("%d\n", tc.n))
		sb.WriteString(tc.s)
		sb.WriteByte('\n')
	}
	return []byte(sb.String()), cases
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierE.go <binary>")
		os.Exit(1)
	}
	cand := os.Args[1]
	ref := "./refE.bin"
	if err := exec.Command("go", "build", "-o", ref, "2096E.go").Run(); err != nil {
		fmt.Println("failed to build reference solution:", err)
		os.Exit(1)
	}
	defer os.Remove(ref)

	rand.Seed(time.Now().UnixNano())

	for iter := 0; iter < 200; iter++ {
		input, cases := buildInput()
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

		refAns, err := parseOutputs(refOut, len(cases))
		if err != nil {
			fmt.Println("failed to parse reference output on iteration", iter+1, ":", err)
			fmt.Println("output:\n", refOut)
			fmt.Println("input:\n", string(input))
			os.Exit(1)
		}
		candAns, err := parseOutputs(candOut, len(cases))
		if err != nil {
			fmt.Println("failed to parse candidate output on iteration", iter+1, ":", err)
			fmt.Println("output:\n", candOut)
			fmt.Println("input:\n", string(input))
			os.Exit(1)
		}

		if len(refAns) != len(candAns) {
			fmt.Printf("answer count mismatch on iteration %d\n", iter+1)
			fmt.Println("input:\n", string(input))
			fmt.Println("reference:\n", refOut)
			fmt.Println("candidate:\n", candOut)
			os.Exit(1)
		}
		for i := range refAns {
			if refAns[i] != candAns[i] {
				fmt.Printf("wrong answer on iteration %d case %d\n", iter+1, i+1)
				fmt.Println("input case:", cases[i].n, cases[i].s)
				fmt.Println("reference:", refAns[i])
				fmt.Println("candidate:", candAns[i])
				os.Exit(1)
			}
		}
	}
	fmt.Println("All tests passed.")
}
