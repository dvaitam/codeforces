package main

import (
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

const numTestsC1 = 100

func main() {
	if len(os.Args) != 2 {
		fmt.Println("Usage: verifierC1.go <binary>")
		os.Exit(1)
	}
	binPath, cleanup, err := prepareBinary(os.Args[1])
	if err != nil {
		fmt.Println("compile error:", err)
		os.Exit(1)
	}
	if cleanup != nil {
		defer cleanup()
	}
	r := rand.New(rand.NewSource(1))
	for t := 1; t <= numTestsC1; t++ {
		n := r.Intn(20) + 1
		a := make([]int, n)
		var sb strings.Builder
		sb.WriteString(fmt.Sprintf("%d\n", n))
		for i := 0; i < n; i++ {
			a[i] = r.Intn(100) + 1
			if i > 0 {
				sb.WriteByte(' ')
			}
			sb.WriteString(fmt.Sprintf("%d", a[i]))
		}
		sb.WriteByte('\n')
		input := sb.String()
		expectedLen := solveC1Len(a)
		out, err := run(binPath, input)
		if err != nil {
			fmt.Printf("test %d: runtime error: %v\n", t, err)
			os.Exit(1)
		}
		out = strings.TrimSpace(out)
		if err := verifyC1(a, out, expectedLen); err != nil {
			fmt.Printf("test %d failed: %v\ninput:%s\noutput:%s\n", t, err, input, out)
			os.Exit(1)
		}
	}
	fmt.Println("OK")
}

func prepareBinary(path string) (string, func(), error) {
	if strings.HasSuffix(path, ".go") {
		tmp := filepath.Join(os.TempDir(), "verify_binC1")
		cmd := exec.Command("go", "build", "-o", tmp, path)
		if out, err := cmd.CombinedOutput(); err != nil {
			return "", nil, fmt.Errorf("go build failed: %v: %s", err, string(out))
		}
		return tmp, func() { os.Remove(tmp) }, nil
	}
	return path, nil, nil
}

func run(bin, input string) (string, error) {
	cmd := exec.Command(bin)
	cmd.Stdin = strings.NewReader(input)
	var buf bytes.Buffer
	cmd.Stdout = &buf
	cmd.Stderr = os.Stderr
	err := cmd.Run()
	return buf.String(), err
}

// solveC1Len computes the optimal (maximum) length using the correct greedy
// from the accepted solution.
func solveC1Len(a []int) int {
	n := len(a)
	l, r := 0, n-1
	last := 0
	count := 0

	for l <= r {
		leftOk := a[l] > last
		rightOk := a[r] > last

		if !leftOk && !rightOk {
			break
		}

		if leftOk && !rightOk {
			last = a[l]
			l++
			count++
		} else if !leftOk && rightOk {
			last = a[r]
			r--
			count++
		} else {
			if a[l] < a[r] {
				last = a[l]
				l++
				count++
			} else if a[r] < a[l] {
				last = a[r]
				r--
				count++
			} else {
				// Both equal and > last: simulate both directions, pick longer
				llast := last
				i := l
				cl := 0
				for i <= r && a[i] > llast {
					cl++
					llast = a[i]
					i++
				}

				rlast := last
				i = r
				cr := 0
				for i >= l && a[i] > rlast {
					cr++
					rlast = a[i]
					i--
				}

				if cl >= cr {
					count += cl
				} else {
					count += cr
				}
				return count
			}
		}
	}
	return count
}

// verifyC1 checks that the candidate output is valid and optimal.
func verifyC1(a []int, out string, expectedLen int) error {
	lines := strings.Split(out, "\n")
	if len(lines) < 1 {
		return fmt.Errorf("no output")
	}
	var k int
	if _, err := fmt.Sscanf(lines[0], "%d", &k); err != nil {
		return fmt.Errorf("cannot parse length: %v", err)
	}
	if k != expectedLen {
		return fmt.Errorf("length mismatch: expected %d got %d", expectedLen, k)
	}
	moves := ""
	if len(lines) > 1 {
		moves = strings.TrimSpace(lines[1])
	}
	if len(moves) != k {
		return fmt.Errorf("move string length %d != declared length %d", len(moves), k)
	}

	n := len(a)
	l, r := 0, n-1
	last := 0
	for i := 0; i < k; i++ {
		ch := moves[i]
		var val int
		if ch == 'L' {
			if l > r {
				return fmt.Errorf("move %d: L but deque empty", i+1)
			}
			val = a[l]
			l++
		} else if ch == 'R' {
			if l > r {
				return fmt.Errorf("move %d: R but deque empty", i+1)
			}
			val = a[r]
			r--
		} else {
			return fmt.Errorf("move %d: invalid char '%c'", i+1, ch)
		}
		if val <= last {
			return fmt.Errorf("move %d: value %d not strictly greater than previous %d", i+1, val, last)
		}
		last = val
	}
	return nil
}
