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
	input    string
	expected string
}

func digits(x int) int {
	c := 0
	for x > 0 {
		c++
		x /= 10
	}
	return c
}

func solveCase(a, b []int) int {
	countA := map[int]int{}
	countB := map[int]int{}
	for _, v := range a {
		countA[v]++
	}
	for _, v := range b {
		countB[v]++
	}
	for val, ca := range countA {
		if cb, ok := countB[val]; ok {
			if ca < cb {
				countB[val] -= ca
				delete(countA, val)
			} else if cb < ca {
				countA[val] -= cb
				delete(countB, val)
			} else {
				delete(countA, val)
				delete(countB, val)
			}
		}
	}
	ops := 0
	newA := map[int]int{}
	for val, cnt := range countA {
		if val >= 10 {
			d := digits(val)
			newA[d] += cnt
			ops += cnt
		} else {
			newA[val] += cnt
		}
	}
	countA = newA
	newB := map[int]int{}
	for val, cnt := range countB {
		if val >= 10 {
			d := digits(val)
			newB[d] += cnt
			ops += cnt
		} else {
			newB[val] += cnt
		}
	}
	countB = newB
	for val, ca := range countA {
		if cb, ok := countB[val]; ok {
			if ca < cb {
				countB[val] -= ca
				delete(countA, val)
			} else if cb < ca {
				countA[val] -= cb
				delete(countB, val)
			} else {
				delete(countA, val)
				delete(countB, val)
			}
		}
	}
	for val, cnt := range countA {
		if val > 1 {
			ops += cnt
			countA[val] -= cnt
			countA[1] += cnt
			delete(countA, val)
		}
	}
	for val, cnt := range countB {
		if val > 1 {
			ops += cnt
			countB[val] -= cnt
			countB[1] += cnt
			delete(countB, val)
		}
	}
	return ops
}

func generateTests() []testCase {
	rng := rand.New(rand.NewSource(3))
	tests := []testCase{}
	tests = append(tests, testCase{input: "1\n1\n1\n1\n", expected: "0"})
	tests = append(tests, testCase{input: "1\n3\n1 2 3\n3 2 1\n", expected: "2"})
	for len(tests) < 100 {
		n := rng.Intn(10) + 1
		a := make([]int, n)
		b := make([]int, n)
		for i := 0; i < n; i++ {
			a[i] = rng.Intn(1000) + 1
			b[i] = rng.Intn(1000) + 1
		}
		ops := solveCase(append([]int{}, a...), append([]int{}, b...))
		var sb strings.Builder
		sb.WriteString(fmt.Sprintf("1\n%d\n", n))
		for i, v := range a {
			if i > 0 {
				sb.WriteByte(' ')
			}
			sb.WriteString(fmt.Sprintf("%d", v))
		}
		sb.WriteByte('\n')
		for i, v := range b {
			if i > 0 {
				sb.WriteByte(' ')
			}
			sb.WriteString(fmt.Sprintf("%d", v))
		}
		sb.WriteByte('\n')
		tests = append(tests, testCase{input: sb.String(), expected: fmt.Sprint(ops)})
	}
	return tests
}

func run(bin, input string) (string, error) {
	var cmd *exec.Cmd
	if strings.HasSuffix(bin, ".go") {
		cmd = exec.Command("go", "run", bin)
	} else {
		cmd = exec.Command(bin)
	}
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	var stderr bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &stderr
	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("runtime error: %v\n%s", err, stderr.String())
	}
	return strings.TrimSpace(out.String()), nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierC.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	_ = rand.New(rand.NewSource(time.Now().UnixNano()))
	tests := generateTests()
	for i, tc := range tests {
		got, err := run(bin, tc.input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:\n%s", i+1, err, tc.input)
			os.Exit(1)
		}
		if got != tc.expected {
			fmt.Fprintf(os.Stderr, "case %d failed: expected %s got %s\ninput:\n%s", i+1, tc.expected, got, tc.input)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
