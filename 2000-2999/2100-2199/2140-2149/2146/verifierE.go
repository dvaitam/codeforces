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
	a []int
}

func run(bin string, input []byte) (string, error) {
	cmd := exec.Command(bin)
	cmd.Stdin = bytes.NewReader(input)
	out, err := cmd.CombinedOutput()
	return string(out), err
}

func parseOutputs(out string, cases []testCase) ([][]int, error) {
	reader := strings.NewReader(out)
	res := make([][]int, 0, len(cases))
	for idx, tc := range cases {
		arr := make([]int, tc.n)
		for i := 0; i < tc.n; i++ {
			if _, err := fmt.Fscan(reader, &arr[i]); err != nil {
				return nil, fmt.Errorf("output ended early on case %d element %d: %v", idx+1, i+1, err)
			}
		}
		res = append(res, arr)
	}
	var extra string
	if _, err := fmt.Fscan(reader, &extra); err == nil {
		return nil, fmt.Errorf("extra output detected after all answers")
	}
	return res, nil
}

func genCases() ([]byte, []testCase) {
	t := rand.Intn(5) + 1
	totalN := 0
	cases := make([]testCase, 0, t)
	for len(cases) < t && totalN < 200 {
		n := rand.Intn(40) + 1
		if totalN+n > 300 {
			n = 1
		}
		a := make([]int, n)
		for i := 0; i < n; i++ {
			a[i] = rand.Intn(50)
		}
		cases = append(cases, testCase{n: n, a: a})
		totalN += n
	}
	if len(cases) == 0 {
		cases = append(cases, testCase{n: 1, a: []int{0}})
	}

	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("%d\n", len(cases)))
	for _, tc := range cases {
		sb.WriteString(fmt.Sprintf("%d\n", tc.n))
		for i, v := range tc.a {
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
		fmt.Println("usage: go run verifierE.go <binary>")
		os.Exit(1)
	}
	cand := os.Args[1]
	ref := "./refE.bin"
	if err := exec.Command("go", "build", "-o", ref, "2146E.go").Run(); err != nil {
		fmt.Println("failed to build reference solution:", err)
		os.Exit(1)
	}
	defer os.Remove(ref)

	rand.Seed(time.Now().UnixNano())

	for iter := 0; iter < 200; iter++ {
		input, cases := genCases()
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

		refAns, err := parseOutputs(refOut, cases)
		if err != nil {
			fmt.Println("failed to parse reference output on iteration", iter+1, ":", err)
			fmt.Println("output:\n", refOut)
			fmt.Println("input:\n", string(input))
			os.Exit(1)
		}
		candAns, err := parseOutputs(candOut, cases)
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
			if len(refAns[i]) != len(candAns[i]) {
				fmt.Printf("length mismatch on iteration %d case %d\n", iter+1, i+1)
				fmt.Println("input:\n", string(input))
				os.Exit(1)
			}
			for j := range refAns[i] {
				if refAns[i][j] != candAns[i][j] {
					fmt.Printf("wrong answer on iteration %d case %d element %d (expected %d got %d)\n", iter+1, i+1, j+1, refAns[i][j], candAns[i][j])
					fmt.Println("input:\n", string(input))
					os.Exit(1)
				}
			}
		}
	}
	fmt.Println("All tests passed.")
}
