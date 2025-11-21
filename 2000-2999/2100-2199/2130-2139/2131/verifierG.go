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

func run(bin string, input []byte) (string, error) {
	cmd := exec.Command(bin)
	cmd.Stdin = bytes.NewReader(input)
	out, err := cmd.CombinedOutput()
	return string(out), err
}

func parseOutputs(out string, t int) ([]int64, error) {
	reader := strings.NewReader(out)
	res := make([]int64, 0, t)
	for i := 0; i < t; i++ {
		var v int64
		if _, err := fmt.Fscan(reader, &v); err != nil {
			return nil, fmt.Errorf("output ended early at case %d: %v", i+1, err)
		}
		res = append(res, v)
	}
	var extra string
	if _, err := fmt.Fscan(reader, &extra); err == nil {
		return nil, fmt.Errorf("extra output after %d cases", t)
	}
	return res, nil
}

type testCase struct {
	n int
	k int64
	s []int64
}

func genCases() ([]byte, []testCase) {
	t := rand.Intn(6) + 1
	remainingN := 200000
	cases := make([]testCase, 0, t)
	for i := 0; i < t && remainingN > 0; i++ {
		maxN := remainingN
		if maxN > 20000 {
			maxN = 20000
		}
		n := rand.Intn(maxN) + 1
		remainingN -= n
		var k int64
		switch rand.Intn(4) {
		case 0:
			k = int64(rand.Intn(n) + 1)
		case 1:
			k = int64(rand.Intn(1_000_000_000) + 1)
		case 2:
			k = int64(n) + int64(rand.Intn(10))
		default:
			k = int64(rand.Intn(100000) + 1)
		}
		start := rand.Int63n(1_000_000_000-int64(n)) + 1
		s := make([]int64, n)
		for j := 0; j < n; j++ {
			s[j] = start + int64(j)
		}
		// Occasionally shuffle to add variety.
		if rand.Intn(3) == 0 {
			for j := n - 1; j > 0; j-- {
				p := rand.Intn(j + 1)
				s[j], s[p] = s[p], s[j]
			}
		}
		cases = append(cases, testCase{n: n, k: k, s: s})
	}
	if len(cases) == 0 {
		cases = append(cases, testCase{n: 1, k: 1, s: []int64{1}})
	}

	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("%d\n", len(cases)))
	for _, tc := range cases {
		sb.WriteString(fmt.Sprintf("%d %d\n", tc.n, tc.k))
		for i, v := range tc.s {
			if i > 0 {
				sb.WriteByte(' ')
			}
			sb.WriteString(fmt.Sprint(v))
		}
		sb.WriteByte('\n')
	}
	return []byte(sb.String()), cases
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierG.go <binary>")
		os.Exit(1)
	}
	cand := os.Args[1]
	ref := "./refG.bin"
	if err := exec.Command("go", "build", "-o", ref, "2131G.go").Run(); err != nil {
		fmt.Println("failed to build reference solution:", err)
		os.Exit(1)
	}
	defer os.Remove(ref)

	rand.Seed(time.Now().UnixNano())

	for iter := 0; iter < 200; iter++ {
		input, _ := genCases()
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

		var t int
		if _, err := fmt.Fscan(strings.NewReader(string(input)), &t); err != nil {
			fmt.Println("failed to parse generated input:", err)
			os.Exit(1)
		}

		refAns, err := parseOutputs(refOut, t)
		if err != nil {
			fmt.Println("failed to parse reference output on iteration", iter+1, ":", err)
			fmt.Println("output:\n", refOut)
			fmt.Println("input:\n", string(input))
			os.Exit(1)
		}
		candAns, err := parseOutputs(candOut, t)
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
				fmt.Printf("wrong answer on iteration %d case %d (expected %d got %d)\n", iter+1, i+1, refAns[i], candAns[i])
				fmt.Println("input:\n", string(input))
				os.Exit(1)
			}
		}
	}
	fmt.Println("All tests passed.")
}
