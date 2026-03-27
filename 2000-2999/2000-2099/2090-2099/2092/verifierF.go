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

// solveReference is the correct solver for 2092F, embedded directly.
func solveReference(input []byte) string {
	data := input
	idx := 0

	nextInt := func() int {
		n := len(data)
		for idx < n && data[idx] <= ' ' {
			idx++
		}
		val := 0
		for idx < n && data[idx] > ' ' {
			val = val*10 + int(data[idx]-'0')
			idx++
		}
		return val
	}

	nextBytes := func() []byte {
		n := len(data)
		for idx < n && data[idx] <= ' ' {
			idx++
		}
		start := idx
		for idx < n && data[idx] > ' ' {
			idx++
		}
		return data[start:idx]
	}

	appendInt := func(dst []byte, x int) []byte {
		if x == 0 {
			return append(dst, '0')
		}
		var buf [20]byte
		i := len(buf)
		for x > 0 {
			i--
			buf[i] = byte('0' + x%10)
			x /= 10
		}
		return append(dst, buf[i:]...)
	}

	t := nextInt()
	out := make([]byte, 0, len(data)*4)

	for ; t > 0; t-- {
		n := nextInt()
		s := nextBytes()

		runs := make([]int, 1, n+1)
		cnt := 1
		for i := 1; i < n; i++ {
			if s[i] == s[i-1] {
				cnt++
			} else {
				runs = append(runs, cnt)
				cnt = 1
			}
		}
		runs = append(runs, cnt)

		R := len(runs) - 1
		diff := make([]int, R+3)

		for b := 1; b < R; b++ {
			low := b + 1
			diff[low]++
			diff[low+1]--

			L := 1
			high := b + 1
			for {
				nxt := L + b
				if nxt <= R && runs[nxt] >= 2 {
					L = nxt
				} else {
					L = nxt + 1
				}
				low = L + b
				if low > R {
					break
				}
				high += b + 1
				if high > R {
					high = R
				}
				diff[low]++
				diff[high+1]--
			}
		}

		curPos := 0
		pref := 0
		first := true
		for r := 1; r <= R; r++ {
			curPos += diff[r]
			base := pref - r + 1 + curPos
			for j := 1; j <= runs[r]; j++ {
				if !first {
					out = append(out, ' ')
				} else {
					first = false
				}
				out = appendInt(out, base+j)
			}
			pref += runs[r]
		}
		out = append(out, '\n')
	}

	return string(out)
}

func run(bin string, input []byte) (string, error) {
	cmd := exec.Command(bin)
	cmd.Stdin = bytes.NewReader(input)
	out, err := cmd.CombinedOutput()
	return string(out), err
}

func runCandidate(path string, input []byte) (string, error) {
	var cmd *exec.Cmd
	if strings.HasSuffix(path, ".go") {
		cmd = exec.Command("go", "run", path)
	} else {
		cmd = exec.Command(path)
	}
	cmd.Stdin = bytes.NewReader(input)
	var stdout bytes.Buffer
	var stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	err := cmd.Run()
	if err != nil {
		msg := stderr.String()
		if msg == "" {
			msg = stdout.String()
		}
		return "", fmt.Errorf("%v\n%s", err, msg)
	}
	return stdout.String(), nil
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
			n = rand.Intn(50000) + 50000
			s = genRandomString(n)
		default:
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

	rand.Seed(time.Now().UnixNano())

	for iter := 0; iter < 200; iter++ {
		cases := genTestCases()
		input := buildInput(cases)

		refOut := solveReference(input)

		candOut, err := runCandidate(cand, input)
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
