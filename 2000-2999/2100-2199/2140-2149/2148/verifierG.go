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
	name  string
	input string
}

func gcd(a, b int) int {
	for b != 0 {
		a, b = b, a%b
	}
	if a < 0 {
		return -a
	}
	return a
}

func buildDivisors(limit int) [][]int {
	divs := make([][]int, limit+1)
	for d := 1; d <= limit; d++ {
		for multiple := d; multiple <= limit; multiple += d {
			divs[multiple] = append(divs[multiple], d)
		}
	}
	return divs
}

func solveRef(input string) (string, error) {
	reader := bufio.NewReader(strings.NewReader(input))
	var t int
	if _, err := fmt.Fscan(reader, &t); err != nil {
		return "", err
	}
	var outputs []string
	for test := 0; test < t; test++ {
		var n int
		fmt.Fscan(reader, &n)
		arr := make([]int, n)
		for i := 0; i < n; i++ {
			fmt.Fscan(reader, &arr[i])
		}
		divs := buildDivisors(n)
		cnt := make([]int, n+1)
		gcdVal := make([]int, n+1)
		good := make([]bool, n+1)
		res := make([]int, n)
		globalG := 0
		best := 0
		for i := 0; i < n; i++ {
			x := arr[i]
			oldG := globalG
			if globalG == 0 {
				globalG = x
			} else {
				globalG = gcd(globalG, x)
			}
			if oldG != 0 && globalG != oldG {
				for _, d := range divs[oldG] {
					if d <= 1 {
						continue
					}
					if globalG%d == 0 {
						continue
					}
					if good[d] && cnt[d] > best {
						best = cnt[d]
					}
				}
			}
			for _, d := range divs[x] {
				if d <= 1 {
					continue
				}
				cnt[d]++
				if gcdVal[d] == 0 {
					gcdVal[d] = x
				} else {
					gcdVal[d] = gcd(gcdVal[d], x)
				}
				if !good[d] && gcdVal[d] == d {
					good[d] = true
				}
				if good[d] && globalG%d != 0 && cnt[d] > best {
					best = cnt[d]
				}
			}
			res[i] = best
		}
		var sb strings.Builder
		for i := 0; i < n; i++ {
			if i > 0 {
				sb.WriteByte(' ')
			}
			fmt.Fprintf(&sb, "%d", res[i])
		}
		outputs = append(outputs, sb.String())
	}
	return strings.Join(outputs, "\n"), nil
}

func runCandidate(bin, input string) (string, error) {
	cmd := exec.Command(bin)
	cmd.Stdin = strings.NewReader(input)
	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("runtime error: %v\nstderr: %s", err, stderr.String())
	}
	if stderr.Len() > 0 {
		return "", fmt.Errorf("unexpected stderr output: %s", stderr.String())
	}
	return strings.TrimSpace(stdout.String()), nil
}

func makeCase(name string, arrays [][]int) testCase {
	var sb strings.Builder
	fmt.Fprintf(&sb, "%d\n", len(arrays))
	for _, arr := range arrays {
		fmt.Fprintf(&sb, "%d\n", len(arr))
		for i, v := range arr {
			if i > 0 {
				sb.WriteByte(' ')
			}
			fmt.Fprintf(&sb, "%d", v)
		}
		sb.WriteByte('\n')
	}
	return testCase{name: name, input: sb.String()}
}

func randomTests() []testCase {
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	var tests []testCase
	for idx := 0; idx < 20; idx++ {
		tcCount := rng.Intn(3) + 1
		arrays := make([][]int, tcCount)
		for t := 0; t < tcCount; t++ {
			n := rng.Intn(6) + 1
			arr := make([]int, n)
			for i := 0; i < n; i++ {
				arr[i] = rng.Intn(n) + 1
			}
			arrays[t] = arr
		}
		tests = append(tests, makeCase(fmt.Sprintf("rand_%d", idx+1), arrays))
	}
	return tests
}

func handcraftedTests() []testCase {
	return []testCase{
		makeCase("simple", [][]int{{1}}),
		makeCase("increasing", [][]int{{1, 2, 3, 4}}),
		makeCase("duplicate", [][]int{{2, 2, 2, 2}}),
		makeCase("example", [][]int{{2, 4, 3, 6, 5, 7, 8}, {6, 6, 6, 6, 6, 6}, {8, 4, 2, 6, 3, 9, 5, 7, 8}}),
	}
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("usage: go run verifierG.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	tests := append(handcraftedTests(), randomTests()...)
	for idx, tc := range tests {
		expect, err := solveRef(tc.input)
		if err != nil {
			fmt.Printf("failed to compute reference for %s: %v\n", tc.name, err)
			os.Exit(1)
		}
		out, err := runCandidate(bin, tc.input)
		if err != nil {
			fmt.Printf("test %d (%s) runtime error: %v\ninput:\n%s\n", idx+1, tc.name, err, tc.input)
			os.Exit(1)
		}
		expectLines := strings.Split(strings.TrimSpace(expect), "\n")
		gotLines := strings.Split(strings.TrimSpace(out), "\n")
		if len(expectLines) != len(gotLines) {
			fmt.Printf("test %d (%s) line count mismatch\ninput:\n%s\nexpect:\n%s\nactual:\n%s\n", idx+1, tc.name, tc.input, expect, out)
			os.Exit(1)
		}
		for i := range expectLines {
			if strings.TrimSpace(expectLines[i]) != strings.TrimSpace(gotLines[i]) {
				fmt.Printf("test %d (%s) mismatch on case %d\ninput:\n%s\nexpect:\n%s\nactual:\n%s\n", idx+1, tc.name, i+1, tc.input, expect, out)
				os.Exit(1)
			}
		}
	}
	fmt.Println("all tests passed")
}
