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

func parse(out string, t int) ([]string, error) {
	reader := strings.NewReader(out)
	res := make([]string, 0, t)
	for i := 0; i < t; i++ {
		var tok string
		if _, err := fmt.Fscan(reader, &tok); err != nil {
			return nil, fmt.Errorf("output ended early at case %d: %v", i+1, err)
		}
		up := strings.ToUpper(tok)
		if up != "YES" && up != "NO" {
			return nil, fmt.Errorf("invalid verdict '%s' on case %d", tok, i+1)
		}
		res = append(res, up)
	}
	var extra string
	if _, err := fmt.Fscan(reader, &extra); err == nil {
		return nil, fmt.Errorf("extra output after %d cases", t)
	}
	return res, nil
}

func genCases() ([]byte, []testCase) {
	t := rand.Intn(20) + 1
	totalN := 0
	cases := make([]testCase, 0, t)
	for len(cases) < t && totalN < 400 {
		n := rand.Intn(98) + 3
		if totalN+n > 500 {
			n = 3
		}
		a := make([]int, n)
		mode := rand.Intn(6)
		switch mode {
		case 0:
			// all -1
			for i := range a {
				a[i] = -1
			}
		case 1:
			// constant positive
			val := rand.Intn(100) + 1
			for i := range a {
				a[i] = val
			}
		case 2:
			// constant zero
			for i := range a {
				a[i] = 0
			}
		case 3:
			// mix different positives
			val1 := rand.Intn(100) + 1
			val2 := val1
			for val2 == val1 {
				val2 = rand.Intn(100) + 1
			}
			for i := range a {
				if rand.Intn(2) == 0 {
					a[i] = val1
				} else {
					a[i] = val2
				}
			}
		default:
			// random with -1 sprinkled
			for i := range a {
				switch rand.Intn(5) {
				case 0:
					a[i] = -1
				case 1:
					a[i] = 0
				default:
					a[i] = rand.Intn(101)
				}
			}
		}
		cases = append(cases, testCase{n: n, a: a})
		totalN += n
	}
	if len(cases) == 0 {
		cases = append(cases, testCase{n: 3, a: []int{-1, -1, -1}})
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
		fmt.Println("usage: go run verifierA.go <binary>")
		os.Exit(1)
	}
	cand := os.Args[1]
	ref := "./refA.bin"
	if err := exec.Command("go", "build", "-o", ref, "2127A.go").Run(); err != nil {
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

		var t int
		if _, err := fmt.Fscan(strings.NewReader(string(input)), &t); err != nil {
			fmt.Println("failed to parse generated input:", err)
			os.Exit(1)
		}

		refAns, err := parse(refOut, t)
		if err != nil {
			fmt.Println("failed to parse reference output on iteration", iter+1, ":", err)
			fmt.Println("output:\n", refOut)
			fmt.Println("input:\n", string(input))
			os.Exit(1)
		}
		candAns, err := parse(candOut, t)
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
				fmt.Println("input case n:", cases[i].n, "a:", cases[i].a)
				fmt.Println("reference:", refAns[i])
				fmt.Println("candidate:", candAns[i])
				os.Exit(1)
			}
		}
	}
	fmt.Println("All tests passed.")
}
