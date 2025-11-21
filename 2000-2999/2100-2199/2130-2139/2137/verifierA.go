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
	k int
	x int64
}

func run(bin string, input []byte) (string, error) {
	cmd := exec.Command(bin)
	cmd.Stdin = bytes.NewReader(input)
	out, err := cmd.CombinedOutput()
	return string(out), err
}

func parseOutputs(out string, t int) ([]int64, error) {
	reader := strings.NewReader(out)
	ans := make([]int64, 0, t)
	for i := 0; i < t; i++ {
		var v int64
		if _, err := fmt.Fscan(reader, &v); err != nil {
			return nil, fmt.Errorf("output ended early at case %d: %v", i+1, err)
		}
		ans = append(ans, v)
	}
	var extra string
	if _, err := fmt.Fscan(reader, &extra); err == nil {
		return nil, fmt.Errorf("extra output detected after %d cases", t)
	}
	return ans, nil
}

func forwardCollatz(x int64, k int) int64 {
	for i := 0; i < k; i++ {
		if x%2 == 0 {
			x /= 2
		} else {
			x = 3*x + 1
		}
	}
	return x
}

func genTests() ([]byte, []testCase) {
	t := rand.Intn(20) + 1
	cases := make([]testCase, 0, t)
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("%d\n", t))
	for i := 0; i < t; i++ {
		k := rand.Intn(20) + 1
		x := rand.Intn(20) + 1
		if rand.Intn(5) == 0 {
			// Include sample-like cases
			k = 5
			x = 4
		}
		cases = append(cases, testCase{k: k, x: int64(x)})
		sb.WriteString(fmt.Sprintf("%d %d\n", k, x))
	}
	return []byte(sb.String()), cases
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierA.go <binary>")
		os.Exit(1)
	}
	cand := os.Args[1]
	ref := "./refA.bin"
	if err := exec.Command("go", "build", "-o", ref, "2137A.go").Run(); err != nil {
		fmt.Println("failed to build reference solution:", err)
		os.Exit(1)
	}
	defer os.Remove(ref)

	rand.Seed(time.Now().UnixNano())

	for iter := 0; iter < 200; iter++ {
		input, cases := genTests()
		if _, err := run(ref, input); err != nil {
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

		candAns, err := parseOutputs(candOut, t)
		if err != nil {
			fmt.Println("failed to parse candidate output on iteration", iter+1, ":", err)
			fmt.Println("output:\n", candOut)
			fmt.Println("input:\n", string(input))
			os.Exit(1)
		}

		for i, tc := range cases {
			val := candAns[i]
			if val < 1 {
				fmt.Printf("invalid value (non-positive) on iteration %d case %d\n", iter+1, i+1)
				fmt.Println("input case:", tc)
				os.Exit(1)
			}
			if forwardCollatz(val, tc.k) != tc.x {
				fmt.Printf("wrong answer on iteration %d case %d\n", iter+1, i+1)
				fmt.Println("input case:", tc)
				fmt.Println("candidate output:", val)
				os.Exit(1)
			}
		}
	}
	fmt.Println("All tests passed.")
}
