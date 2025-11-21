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

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

// Reference feasibility checker from the official solution logic.
func solveCase(n, m, k int) (bool, string) {
	if k > max(n, m) {
		return false, ""
	}
	lower := max(0, max(k-m, n-m))
	upper := min(k, min(n, n-m+k))
	if lower > upper {
		return false, ""
	}
	dMax := upper
	dMin := dMax - k

	cur := 0
	zeros := n
	ones := m
	desired := n - m
	res := make([]byte, 0, n+m)

	maxHit := dMax == 0
	minHit := dMin == 0

	for cur < dMax {
		if zeros == 0 {
			return false, ""
		}
		res = append(res, '0')
		cur++
		zeros--
	}
	if cur == dMax {
		maxHit = true
	}

	for zeros > 0 || ones > 0 {
		if !minHit && cur > dMin {
			if ones == 0 {
				return false, ""
			}
			res = append(res, '1')
			cur--
			ones--
			if cur == dMin {
				minHit = true
			}
			continue
		}
		if cur > desired {
			if ones == 0 {
				return false, ""
			}
			res = append(res, '1')
			cur--
			ones--
			if cur == dMin {
				minHit = true
			}
			continue
		}
		if cur < desired {
			if zeros == 0 || cur == dMax {
				return false, ""
			}
			res = append(res, '0')
			cur++
			zeros--
			if cur == dMax {
				maxHit = true
			}
			continue
		}
		if zeros > 0 && cur < dMax {
			res = append(res, '0')
			cur++
			zeros--
			if cur == dMax {
				maxHit = true
			}
			continue
		}
		if ones > 0 && cur > dMin {
			res = append(res, '1')
			cur--
			ones--
			if cur == dMin {
				minHit = true
			}
			continue
		}
		return false, ""
	}

	if cur != desired || !maxHit || !minHit {
		return false, ""
	}

	return true, string(res)
}

type testCase struct {
	n int
	m int
	k int
}

func generateCases(rng *rand.Rand, t int) []testCase {
	cases := make([]testCase, t)
	for i := 0; i < t; i++ {
		n := rng.Intn(50)
		m := rng.Intn(50)
		if n == 0 && m == 0 {
			n = 1
		}
		k := rng.Intn(n+m) + 1
		cases[i] = testCase{n: n, m: m, k: k}
	}
	return cases
}

func buildInput(cases []testCase) string {
	var sb strings.Builder
	fmt.Fprintln(&sb, len(cases))
	for _, tc := range cases {
		fmt.Fprintf(&sb, "%d %d %d\n", tc.n, tc.m, tc.k)
	}
	return sb.String()
}

func runCandidate(bin, input string) (string, error) {
	var cmd *exec.Cmd
	if strings.HasSuffix(bin, ".go") {
		cmd = exec.Command("go", "run", bin)
	} else {
		cmd = exec.Command(bin)
	}
	cmd.Stdin = strings.NewReader(input)
	var out, stderr bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &stderr
	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("runtime error: %v\n%s", err, stderr.String())
	}
	return out.String(), nil
}

func maxBalance(s string) int {
	pref := 0
	minPref, maxPref := 0, 0
	for _, ch := range s {
		if ch == '0' {
			pref++
		} else if ch == '1' {
			pref--
		} else {
			return -1
		}
		if pref < minPref {
			minPref = pref
		}
		if pref > maxPref {
			maxPref = pref
		}
	}
	return maxPref - minPref
}

func verifyCase(tc testCase, outLine string) error {
	outLine = strings.TrimSpace(outLine)
	possible, _ := solveCase(tc.n, tc.m, tc.k)
	if outLine == "-1" {
		if possible {
			return fmt.Errorf("case possible but candidate printed -1")
		}
		return nil
	}
	if !possible {
		return fmt.Errorf("case impossible but candidate printed solution")
	}
	if len(outLine) != tc.n+tc.m {
		return fmt.Errorf("wrong length: got %d expected %d", len(outLine), tc.n+tc.m)
	}
	zeros := 0
	ones := 0
	for i, ch := range outLine {
		if ch == '0' {
			zeros++
		} else if ch == '1' {
			ones++
		} else {
			return fmt.Errorf("invalid character at pos %d", i)
		}
	}
	if zeros != tc.n || ones != tc.m {
		return fmt.Errorf("count mismatch zeros %d/%d ones %d/%d", zeros, tc.n, ones, tc.m)
	}
	mb := maxBalance(outLine)
	if mb != tc.k {
		return fmt.Errorf("max balance %d expected %d", mb, tc.k)
	}
	return nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierE.go /path/to/2065E_binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	cases := generateCases(rng, 200)
	input := buildInput(cases)
	output, err := runCandidate(bin, input)
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to run candidate: %v\n", err)
		os.Exit(1)
	}
	lines := strings.Split(strings.TrimRight(output, "\n"), "\n")
	var cleaned []string
	for _, ln := range lines {
		if strings.TrimSpace(ln) != "" {
			cleaned = append(cleaned, ln)
		}
	}
	if len(cleaned) < len(cases) {
		fmt.Fprintf(os.Stderr, "not enough output lines: got %d expected %d\n", len(cleaned), len(cases))
		os.Exit(1)
	}
	for i, tc := range cases {
		if err := verifyCase(tc, cleaned[i]); err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput: %d %d %d\noutput: %s\n", i+1, err, tc.n, tc.m, tc.k, cleaned[i])
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
