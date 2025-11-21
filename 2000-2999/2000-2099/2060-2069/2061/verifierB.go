package main

import (
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"time"
)

func run(bin string, input []byte) (string, error) {
	cmd := exec.Command(bin)
	cmd.Stdin = bytes.NewReader(input)
	out, err := cmd.CombinedOutput()
	return string(out), err
}

func parseOutput(out string, t int) ([][]int, error) {
	reader := strings.NewReader(out)
	res := make([][]int, 0, t)
	for i := 0; i < t; i++ {
		var tok string
		if _, err := fmt.Fscan(reader, &tok); err != nil {
			return nil, fmt.Errorf("output ended early on test %d: %v", i+1, err)
		}
		if tok == "-1" {
			res = append(res, nil)
			continue
		}
		first, err := strconv.Atoi(tok)
		if err != nil {
			return nil, fmt.Errorf("invalid number '%s' on test %d", tok, i+1)
		}
		quad := []int{first}
		for k := 0; k < 3; k++ {
			var v int
			if _, err := fmt.Fscan(reader, &v); err != nil {
				return nil, fmt.Errorf("output ended early on test %d", i+1)
			}
			quad = append(quad, v)
		}
		res = append(res, quad)
	}
	// Ensure no extra tokens remain (ignore trailing whitespace).
	var extra string
	if _, err := fmt.Fscan(reader, &extra); err == nil {
		return nil, fmt.Errorf("extra output detected after %d testcases", t)
	}
	return res, nil
}

func abs(x int) int {
	if x < 0 {
		return -x
	}
	return x
}

func validSolution(arr []int, chosen []int) bool {
	if len(chosen) != 4 {
		return false
	}
	count := make(map[int]int)
	for _, v := range arr {
		count[v]++
	}
	for _, v := range chosen {
		if count[v] == 0 {
			return false
		}
		count[v]--
	}
	freq := make(map[int]int)
	for _, v := range chosen {
		freq[v]++
	}
	for l, c := range freq {
		if c < 2 {
			continue
		}
		rem := make([]int, 0, 2)
		for val, cnt := range freq {
			use := cnt
			if val == l {
				use -= 2
			}
			for k := 0; k < use; k++ {
				rem = append(rem, val)
			}
		}
		if len(rem) != 2 {
			continue
		}
		a, b := rem[0], rem[1]
		if a == b || 2*l > abs(a-b) {
			return true
		}
	}
	return false
}

func genSolutionCase() (int, []int) {
	l := rand.Intn(50) + 1
	var a, b int
	if rand.Intn(2) == 0 {
		a = rand.Intn(50) + 1
		b = a
	} else {
		diff := rand.Intn(2*l-1) + 1 // 1 .. 2*l-1 ensures positive area
		a = rand.Intn(50) + 1
		b = a + diff
	}
	arr := []int{l, l, a, b}
	dupLegs := rand.Intn(3)
	for i := 0; i < dupLegs; i++ {
		arr = append(arr, l)
	}
	n := rand.Intn(15) + len(arr)
	for len(arr) < n {
		arr = append(arr, rand.Intn(100)+1)
	}
	return n, arr
}

func genNoSolutionCase() (int, []int) {
	// Create a near-miss: legs duplicated but bases too far apart.
	l := rand.Intn(50) + 1
	a := rand.Intn(50) + 1
	b := a + 2*l // height would be zero
	arr := []int{l, l, a, b}
	n := rand.Intn(12) + 4
	used := make(map[int]bool)
	for _, v := range arr {
		used[v] = true
	}
	for len(arr) < n {
		val := rand.Intn(200) + 60 // keep away from small numbers to avoid accidental solutions
		if used[val] {
			continue
		}
		used[val] = true
		arr = append(arr, val)
	}
	// Occasionally make everything unique to guarantee impossibility.
	if rand.Intn(3) == 0 {
		arr = arr[:0]
		n = rand.Intn(20) + 4
		seen := make(map[int]bool)
		for len(arr) < n {
			val := rand.Intn(1_000_000_000) + 1
			if seen[val] {
				continue
			}
			seen[val] = true
			arr = append(arr, val)
		}
	}
	return len(arr), arr
}

func genTest() []byte {
	t := rand.Intn(6) + 1
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("%d\n", t))
	for i := 0; i < t; i++ {
		var n int
		var arr []int
		switch rand.Intn(20) {
		case 0:
			// Large uniform array: always solvable.
			n = 200000
			val := rand.Intn(100000000) + 1
			arr = make([]int, n)
			for j := 0; j < n; j++ {
				arr[j] = val
			}
		case 1, 2, 3, 4, 5:
			n, arr = genSolutionCase()
		case 6, 7:
			n, arr = genNoSolutionCase()
		default:
			// Mixed random; let the reference answer decide.
			n = rand.Intn(50) + 4
			arr = make([]int, n)
			for j := 0; j < n; j++ {
				arr[j] = rand.Intn(100000000) + 1
			}
		}
		sb.WriteString(fmt.Sprintf("%d\n", n))
		for j, v := range arr {
			if j > 0 {
				sb.WriteByte(' ')
			}
			sb.WriteString(fmt.Sprint(v))
		}
		sb.WriteByte('\n')
	}
	return []byte(sb.String())
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierB.go <binary>")
		os.Exit(1)
	}
	cand := os.Args[1]
	ref := "./refB.bin"
	if err := exec.Command("go", "build", "-o", ref, "2061B.go").Run(); err != nil {
		fmt.Println("failed to build reference solution:", err)
		os.Exit(1)
	}
	defer os.Remove(ref)

	rand.Seed(time.Now().UnixNano())

	for i := 0; i < 200; i++ {
		input := genTest()
		wantOut, err := run(ref, input)
		if err != nil {
			fmt.Println("reference solution failed:", err)
			fmt.Println("input:\n", string(input))
			os.Exit(1)
		}
		gotOut, err := run(cand, input)
		if err != nil {
			fmt.Printf("candidate runtime error on iteration %d: %v\n", i+1, err)
			fmt.Println("input:\n", string(input))
			os.Exit(1)
		}

		var t int
		if _, err := fmt.Fscan(strings.NewReader(string(input)), &t); err != nil {
			fmt.Println("failed to re-read test count:", err)
			os.Exit(1)
		}

		wantParsed, err := parseOutput(wantOut, t)
		if err != nil {
			fmt.Println("failed to parse reference output:", err)
			fmt.Println("input:\n", string(input))
			fmt.Println("output:\n", wantOut)
			os.Exit(1)
		}
		gotParsed, err := parseOutput(gotOut, t)
		if err != nil {
			fmt.Printf("failed to parse candidate output on iteration %d: %v\n", i+1, err)
			fmt.Println("input:\n", string(input))
			fmt.Println("output:\n", gotOut)
			os.Exit(1)
		}

		// Validate per test case.
		s := strings.Split(strings.TrimSpace(string(input)), "\n")
		ptr := 1
		for tc := 0; tc < t; tc++ {
			// Read n and array for this test case.
			var n int
			fmt.Sscan(s[ptr], &n)
			ptr++
			arrStr := strings.Fields(s[ptr])
			ptr++
			arr := make([]int, n)
			for idx := 0; idx < n; idx++ {
				arr[idx], _ = strconv.Atoi(arrStr[idx])
			}

			wantCase := wantParsed[tc]
			gotCase := gotParsed[tc]

			if wantCase == nil {
				if gotCase != nil {
					fmt.Printf("candidate found a solution when none expected on test %d of iteration %d\n", tc+1, i+1)
					fmt.Println("input case:", s[ptr-2], s[ptr-1])
					fmt.Println("candidate output:", gotCase)
					fmt.Println("reference output: -1")
					os.Exit(1)
				}
				continue
			}

			if gotCase == nil || !validSolution(arr, gotCase) {
				fmt.Printf("invalid solution on test %d of iteration %d\n", tc+1, i+1)
				fmt.Println("array:", arr)
				fmt.Println("candidate output:", gotCase)
				fmt.Println("reference output:", wantCase)
				os.Exit(1)
			}
		}
	}
	fmt.Println("All tests passed.")
}
