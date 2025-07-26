package main

import (
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"strings"
)

type test struct {
	input    string
	expected string
}

func runBinary(bin, input string) (string, error) {
	var cmd *exec.Cmd
	if strings.HasSuffix(bin, ".go") {
		cmd = exec.Command("go", "run", bin)
	} else {
		cmd = exec.Command(bin)
	}
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	err := cmd.Run()
	return strings.TrimSpace(out.String()), err
}

func bruteForce(a []int) int {
	n := len(a)
	mod := 1000000007
	count := 0
	for l1 := 0; l1 < n; l1++ {
		for r1 := l1; r1 < n; r1++ {
			for l2 := 0; l2 < n; l2++ {
				for r2 := l2; r2 < n; r2++ {
					if (l1 <= l2 && r2 <= r1) || (l2 <= l1 && r1 <= r2) {
						continue
					}
					L := l1
					if l2 > L {
						L = l2
					}
					R := r1
					if r2 < R {
						R = r2
					}
					if L > R {
						continue
					}
					ok := true
					for i := L; i <= R && ok; i++ {
						val := a[i]
						cnt1 := 0
						for j := l1; j <= r1; j++ {
							if a[j] == val {
								cnt1++
							}
						}
						cnt2 := 0
						for j := l2; j <= r2; j++ {
							if a[j] == val {
								cnt2++
							}
						}
						if cnt1 != 1 || cnt2 != 1 {
							ok = false
						}
					}
					if ok {
						count++
					}
				}
			}
		}
	}
	return count % mod
}

func solveD(input string) string {
	fields := strings.Fields(input)
	idx := 0
	n := 0
	fmt.Sscan(fields[idx], &n)
	idx++
	a := make([]int, n)
	for i := 0; i < n; i++ {
		fmt.Sscan(fields[idx+i], &a[i])
	}
	ans := bruteForce(a)
	return fmt.Sprintf("%d", ans)
}

func generateTests() []test {
	rng := rand.New(rand.NewSource(45))
	var tests []test
	for len(tests) < 100 {
		n := rng.Intn(5) + 1
		a := make([]int, n)
		for i := 0; i < n; i++ {
			a[i] = rng.Intn(5)
		}
		var sb strings.Builder
		fmt.Fprintf(&sb, "%d\n", n)
		for i := 0; i < n; i++ {
			if i > 0 {
				sb.WriteByte(' ')
			}
			fmt.Fprintf(&sb, "%d", a[i])
		}
		sb.WriteByte('\n')
		input := sb.String()
		expected := solveD(strings.TrimSpace(fmt.Sprintf("%d %s", n, strings.TrimSpace(strings.Join(strings.Fields(sb.String())[1:], " ")))))
		expected = solveD(input)
		tests = append(tests, test{input, expected})
	}
	return tests
}

func main() {
	if len(os.Args) < 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierD.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	tests := generateTests()
	for i, t := range tests {
		out, err := runBinary(bin, t.input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d runtime error: %v\n", i+1, err)
			os.Exit(1)
		}
		if strings.TrimSpace(out) != strings.TrimSpace(t.expected) {
			fmt.Fprintf(os.Stderr, "case %d failed\ninput:\n%sexpected:%s\ngot:%s\n", i+1, t.input, t.expected, out)
			os.Exit(1)
		}
	}
	fmt.Printf("All %d tests passed.\n", len(tests))
}
