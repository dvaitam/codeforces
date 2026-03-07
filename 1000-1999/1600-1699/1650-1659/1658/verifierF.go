package main

import (
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"strconv"
	"strings"
)

type Case struct{ input string }

func runExe(bin, input string) (string, error) {
	cmd := exec.Command(bin)
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("runtime error: %v\n%s", err, out.String())
	}
	return strings.TrimSpace(out.String()), nil
}

func genCases() []Case {
	rng := rand.New(rand.NewSource(7))
	cases := make([]Case, 0, 105)
	for i := 0; i < 100; i++ {
		t := rng.Intn(3) + 1
		var sb strings.Builder
		fmt.Fprintf(&sb, "%d\n", t)
		for j := 0; j < t; j++ {
			n := rng.Intn(50) + 1
			m := rng.Intn(n) + 1
			fmt.Fprintf(&sb, "%d %d\n", n, m)
			for k := 0; k < n; k++ {
				if rng.Intn(2) == 0 {
					sb.WriteByte('0')
				} else {
					sb.WriteByte('1')
				}
			}
			sb.WriteByte('\n')
		}
		cases = append(cases, Case{sb.String()})
	}
	cases = append(cases, Case{"1\n1 1\n0\n"})
	cases = append(cases, Case{"1\n4 2\n1100\n"})
	cases = append(cases, Case{"1\n8 6\n11000011\n"})
	cases = append(cases, Case{"1\n4 4\n0101\n"})
	cases = append(cases, Case{"2\n3 1\n101\n2 2\n01\n"})
	return cases
}

func count1s(s string, l, r int) int { // l,r 0-based inclusive
	c := 0
	for i := l; i <= r; i++ {
		if s[i] == '1' {
			c++
		}
	}
	return c
}

func canK1(s string, n, m, target int) bool {
	cur := count1s(s, 0, m-1)
	if cur == target {
		return true
	}
	for i := m; i < n; i++ {
		if s[i] == '1' {
			cur++
		}
		if s[i-m] == '1' {
			cur--
		}
		if cur == target {
			return true
		}
	}
	return false
}

// validate checks one test case's output. Returns error string or "".
func validate(s string, n, m int, output string) string {
	cnt1 := count1s(s, 0, n-1)
	if (cnt1*m)%n != 0 {
		// must be -1
		if strings.TrimSpace(output) == "-1" {
			return ""
		}
		return fmt.Sprintf("expected -1 (cnt1*m mod n != 0), got: %s", output)
	}
	target := cnt1 * m / n

	lines := strings.Fields(output)
	if len(lines) == 0 {
		return "empty output"
	}
	if lines[0] == "-1" {
		// candidate claims impossible; verify
		if canK1(s, n, m, target) {
			return "candidate said -1 but k=1 solution exists"
		}
		// Also check k=2 via complement window
		L := n - m
		if L > 0 {
			need := cnt1 - target
			cur := count1s(s, 0, L-1)
			if cur == need {
				return "candidate said -1 but k=2 solution exists"
			}
			for i := L; i < n; i++ {
				if s[i] == '1' {
					cur++
				}
				if s[i-L] == '1' {
					cur--
				}
				if cur == need {
					return "candidate said -1 but k=2 solution exists"
				}
			}
		}
		return ""
	}

	k, err := strconv.Atoi(lines[0])
	if err != nil || k < 1 {
		return fmt.Sprintf("invalid k: %s", lines[0])
	}
	if len(lines) != 1+2*k {
		return fmt.Sprintf("expected %d numbers after k=%d, got %d tokens total", 2*k, k, len(lines)-1)
	}

	segs := make([][2]int, k)
	for i := 0; i < k; i++ {
		l, e1 := strconv.Atoi(lines[1+2*i])
		r, e2 := strconv.Atoi(lines[2+2*i])
		if e1 != nil || e2 != nil {
			return fmt.Sprintf("invalid segment %d", i+1)
		}
		segs[i] = [2]int{l, r}
	}

	// Check ordering and bounds
	prev := 0
	totalLen := 0
	totalOnes := 0
	for i, seg := range segs {
		l, r := seg[0], seg[1]
		if l < 1 || r > n || l > r {
			return fmt.Sprintf("segment %d [%d,%d] out of bounds", i+1, l, r)
		}
		if l <= prev {
			return fmt.Sprintf("segment %d starts at %d but previous ended at %d", i+1, l, prev)
		}
		totalLen += r - l + 1
		totalOnes += count1s(s, l-1, r-1)
		prev = r
	}
	if totalLen != m {
		return fmt.Sprintf("total length %d != m=%d", totalLen, m)
	}
	if totalOnes != target {
		return fmt.Sprintf("total ones %d != target=%d", totalOnes, target)
	}

	// Check minimality: if k>1, verify k=1 is impossible
	if k > 1 && canK1(s, n, m, target) {
		return fmt.Sprintf("candidate used k=%d but k=1 is achievable", k)
	}

	return ""
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierF.go /path/to/binary")
		return
	}
	bin := os.Args[1]
	if strings.HasSuffix(bin, ".go") {
		tmp, err := os.CreateTemp("", "verifierF-bin-*")
		if err != nil {
			fmt.Fprintf(os.Stderr, "failed to create temp file: %v\n", err)
			os.Exit(1)
		}
		tmp.Close()
		defer os.Remove(tmp.Name())
		out, err := exec.Command("go", "build", "-o", tmp.Name(), bin).CombinedOutput()
		if err != nil {
			fmt.Fprintf(os.Stderr, "compile error: %v\n%s", err, out)
			os.Exit(1)
		}
		bin = tmp.Name()
	}

	cases := genCases()
	for i, c := range cases {
		// Parse all test cases in this batch
		lines := strings.Split(strings.TrimSpace(c.input), "\n")
		idx := 0
		t, _ := strconv.Atoi(strings.TrimSpace(lines[idx]))
		idx++

		// Run candidate on the whole batch
		got, err := runExe(bin, c.input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "candidate runtime error on case %d: %v\n", i+1, err)
			os.Exit(1)
		}
		outLines := strings.Split(strings.TrimSpace(got), "\n")
		outIdx := 0

		for tc := 0; tc < t; tc++ {
			parts := strings.Fields(lines[idx])
			idx++
			n, _ := strconv.Atoi(parts[0])
			m, _ := strconv.Atoi(parts[1])
			s := strings.TrimSpace(lines[idx])
			idx++

			if outIdx >= len(outLines) {
				fmt.Fprintf(os.Stderr, "case %d test %d: not enough output lines\n", i+1, tc+1)
				os.Exit(1)
			}

			// Collect output lines for this test case
			var tcOut []string
			tcOut = append(tcOut, outLines[outIdx])
			outIdx++
			if outLines[outIdx-1] != "-1" {
				k, _ := strconv.Atoi(outLines[outIdx-1])
				if k >= 1 {
					for j := 0; j < k && outIdx < len(outLines); j++ {
						tcOut = append(tcOut, outLines[outIdx])
						outIdx++
					}
				}
			}

			tcOutStr := strings.Join(tcOut, "\n")
			if msg := validate(s, n, m, tcOutStr); msg != "" {
				fmt.Fprintf(os.Stderr, "case %d test %d failed: %s\ninput: %d %d %s\noutput: %s\n",
					i+1, tc+1, msg, n, m, s, tcOutStr)
				os.Exit(1)
			}
		}
	}
	fmt.Println("All tests passed")
}
