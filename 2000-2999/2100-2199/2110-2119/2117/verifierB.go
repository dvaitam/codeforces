package main

import (
	"bufio"
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
}

func run(bin string, input []byte) (string, error) {
	cmd := exec.Command(bin)
	cmd.Stdin = bytes.NewReader(input)
	out, err := cmd.CombinedOutput()
	return string(out), err
}

func parsePerms(out string, cases []testCase) ([][]int, error) {
	reader := bufio.NewReader(strings.NewReader(out))
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
		return nil, fmt.Errorf("extra output detected after all cases")
	}
	return res, nil
}

func isPerm(arr []int, n int) bool {
	if len(arr) != n {
		return false
	}
	seen := make([]bool, n+1)
	for _, v := range arr {
		if v < 1 || v > n || seen[v] {
			return false
		}
		seen[v] = true
	}
	return true
}

func shrinkScore(a []int) int {
	n := len(a)
	if n < 3 {
		return 0
	}
	prev := make([]int, n)
	next := make([]int, n)
	for i := 0; i < n; i++ {
		prev[i] = i - 1
		next[i] = i + 1
	}
	next[n-1] = -1
	inQ := make([]bool, n)
	queue := make([]int, 0)
	checkPeak := func(i int) bool {
		l := prev[i]
		r := next[i]
		if l == -1 || r == -1 {
			return false
		}
		return a[i] > a[l] && a[i] > a[r]
	}
	pushIfPeak := func(i int) {
		if i == -1 || inQ[i] {
			return
		}
		if checkPeak(i) {
			inQ[i] = true
			queue = append(queue, i)
		}
	}
	for i := 1; i <= n-2; i++ {
		pushIfPeak(i)
	}
	removed := make([]bool, n)
	score := 0
	for len(queue) > 0 {
		i := queue[0]
		queue = queue[1:]
		inQ[i] = false
		if removed[i] || !checkPeak(i) {
			continue
		}
		score++
		removed[i] = true
		l := prev[i]
		r := next[i]
		if l != -1 {
			next[l] = r
		}
		if r != -1 {
			prev[r] = l
		}
		if l != -1 && r != -1 {
			pushIfPeak(l)
			pushIfPeak(r)
		} else {
			if l != -1 {
				pushIfPeak(l)
			}
			if r != -1 {
				pushIfPeak(r)
			}
		}
	}
	return score
}

func genCases() ([]byte, []testCase) {
	// Keep total n moderate.
	t := rand.Intn(8) + 1
	remaining := 20000
	var sb strings.Builder
	cases := make([]testCase, 0, t)
	for i := 0; i < t && remaining >= 3; i++ {
		var n int
		switch rand.Intn(6) {
		case 0:
			n = 3
		case 1:
			n = rand.Intn(10) + 3
		case 2:
			n = rand.Intn(100) + 3
		case 3:
			n = rand.Intn(1000) + 3
		default:
			n = rand.Intn(5000) + 3
		}
		if n > remaining {
			n = remaining
		}
		remaining -= n
		cases = append(cases, testCase{n: n})
	}
	if len(cases) == 0 {
		cases = append(cases, testCase{n: 3})
	}
	sb.WriteString(fmt.Sprintf("%d\n", len(cases)))
	for _, c := range cases {
		sb.WriteString(fmt.Sprintf("%d\n", c.n))
	}
	return []byte(sb.String()), cases
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierB.go <binary>")
		os.Exit(1)
	}
	cand := os.Args[1]
	ref := "./refB.bin"
	if err := exec.Command("go", "build", "-o", ref, "2117B.go").Run(); err != nil {
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

		refPerms, err := parsePerms(refOut, cases)
		if err != nil {
			fmt.Println("failed to parse reference output on iteration", iter+1, ":", err)
			fmt.Println("output:\n", refOut)
			fmt.Println("input:\n", string(input))
			os.Exit(1)
		}
		candPerms, err := parsePerms(candOut, cases)
		if err != nil {
			fmt.Println("failed to parse candidate output on iteration", iter+1, ":", err)
			fmt.Println("output:\n", candOut)
			fmt.Println("input:\n", string(input))
			os.Exit(1)
		}

		for i, tc := range cases {
			if !isPerm(refPerms[i], tc.n) {
				fmt.Printf("reference produced invalid permutation on case %d\n", i+1)
				os.Exit(1)
			}
			if !isPerm(candPerms[i], tc.n) {
				fmt.Printf("candidate produced invalid permutation on iteration %d case %d\n", iter+1, i+1)
				fmt.Println("input n:", tc.n)
				os.Exit(1)
			}
			refScore := shrinkScore(refPerms[i])
			candScore := shrinkScore(candPerms[i])
			if refScore != candScore {
				fmt.Printf("wrong score on iteration %d case %d (expected %d, got %d)\n", iter+1, i+1, refScore, candScore)
				fmt.Println("input n:", tc.n)
				os.Exit(1)
			}
		}
	}
	fmt.Println("All tests passed.")
}
