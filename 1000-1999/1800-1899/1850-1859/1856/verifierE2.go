package main

import (
	"bytes"
	"fmt"
	"math/bits"
	"math/rand"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"time"
)

type testCase struct {
	name  string
	input string
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierE2.go /path/to/binary")
		os.Exit(1)
	}

	tests := generateTests()
	candidate := os.Args[1]

	for idx, tc := range tests {
		refVal := solveReference(tc.input)

		candOut, err := runCandidate(candidate, tc.input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "candidate runtime error on test %d (%s): %v\ninput:\n%soutput:\n%s\n", idx+1, tc.name, err, tc.input, candOut)
			os.Exit(1)
		}
		candVal, err := parseOutput(candOut)
		if err != nil {
			fmt.Fprintf(os.Stderr, "failed to parse candidate output on test %d (%s): %v\noutput:\n%s\n", idx+1, tc.name, err, candOut)
			os.Exit(1)
		}

		if refVal != candVal {
			fmt.Fprintf(os.Stderr, "wrong answer on test %d (%s): expected %d, got %d\ninput:\n%s\n", idx+1, tc.name, refVal, candVal, tc.input)
			os.Exit(1)
		}
	}

	fmt.Printf("All %d tests passed.\n", len(tests))
}

// Embedded correct solver for 1856E2
func solveReference(input string) int64 {
	data := []byte(input)
	pos := 0

	nextInt := func() int {
		for pos < len(data) && data[pos] <= ' ' {
			pos++
		}
		if pos >= len(data) {
			return 0
		}
		res := 0
		for pos < len(data) && data[pos] > ' ' {
			res = res*10 + int(data[pos]-'0')
			pos++
		}
		return res
	}

	n := nextInt()
	if n == 0 {
		return 0
	}

	parent := make([]int, n+1)
	for i := 2; i <= n; i++ {
		parent[i] = nextInt()
	}

	head := make([]int, n+1)
	next := make([]int, n+1)
	for i := n; i >= 2; i-- {
		p := parent[i]
		next[i] = head[p]
		head[p] = i
	}

	s := make([]int, n+1)
	for i := 1; i <= n; i++ {
		s[i] = 1
	}
	for i := n; i >= 2; i-- {
		s[parent[i]] += s[i]
	}

	count := make([]int, n+1)
	var items []int
	var distinct []int
	var bs []uint64

	var totalAns int64

	for u := 1; u <= n; u++ {
		S := s[u] - 1
		if S <= 0 {
			continue
		}

		maxChild := 0
		for v := head[u]; v != 0; v = next[v] {
			if s[v] > maxChild {
				maxChild = s[v]
			}
		}

		if maxChild >= S/2 {
			totalAns += int64(maxChild) * int64(S-maxChild)
			continue
		}

		distinct = distinct[:0]
		items = items[:0]

		for v := head[u]; v != 0; v = next[v] {
			c := s[v]
			if count[c] == 0 {
				distinct = append(distinct, c)
			}
			count[c]++
		}

		for _, w := range distinct {
			c := count[w]
			count[w] = 0
			k := 1
			for c > 0 {
				if k > c {
					k = c
				}
				items = append(items, w*k)
				c -= k
				k *= 2
			}
		}

		target := S / 2
		words := target/64 + 1
		if len(bs) < words {
			newSize := words * 2
			if newSize < 16 {
				newSize = 16
			}
			bs = make([]uint64, newSize)
		}
		for i := 0; i < words; i++ {
			bs[i] = 0
		}
		bs[0] = 1

		for _, v := range items {
			if v > target {
				continue
			}
			wordShift := v / 64
			bitShift := uint64(v % 64)

			if bitShift == 0 {
				for i := words - 1; i >= wordShift; i-- {
					bs[i] |= bs[i-wordShift]
				}
			} else {
				compShift := 64 - bitShift
				for i := words - 1; i > wordShift; i-- {
					bs[i] |= (bs[i-wordShift] << bitShift) | (bs[i-wordShift-1] >> compShift)
				}
				bs[wordShift] |= (bs[0] << bitShift)
			}
		}

		w := target / 64
		b := target % 64
		var mask uint64
		if b == 63 {
			mask = ^uint64(0)
		} else {
			mask = (uint64(1) << (b + 1)) - 1
		}

		val := bs[w] & mask
		if val != 0 {
			highestBit := bits.Len64(val) - 1
			ans := w*64 + highestBit
			totalAns += int64(ans) * int64(S-ans)
		} else {
			for i := w - 1; i >= 0; i-- {
				if bs[i] != 0 {
					highestBit := bits.Len64(bs[i]) - 1
					ans := i*64 + highestBit
					totalAns += int64(ans) * int64(S-ans)
					break
				}
			}
		}
	}

	return totalAns
}

func runCandidate(path, input string) (string, error) {
	cmd := commandFor(path)
	return runWithInput(cmd, input)
}

func commandFor(path string) *exec.Cmd {
	if strings.HasSuffix(path, ".go") {
		return exec.Command("go", "run", path)
	}
	return exec.Command(path)
}

func runWithInput(cmd *exec.Cmd, input string) (string, error) {
	cmd.Stdin = strings.NewReader(input)
	var stdout bytes.Buffer
	var stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	if err := cmd.Run(); err != nil {
		return stdout.String() + stderr.String(), err
	}
	return stdout.String(), nil
}

func parseOutput(out string) (int64, error) {
	lines := strings.Fields(out)
	if len(lines) == 0 {
		return 0, fmt.Errorf("empty output")
	}
	val, err := strconv.ParseInt(lines[0], 10, 64)
	if err != nil {
		return 0, fmt.Errorf("invalid integer %q: %v", lines[0], err)
	}
	return val, nil
}

func generateTests() []testCase {
	tests := []testCase{
		buildCase("n1", []int{1}),
		buildCase("chain", []int{5, 1, 2, 3, 4}),
	}

	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < 40; i++ {
		n := rng.Intn(10) + 1
		parents := make([]int, n-1)
		for j := 2; j <= n; j++ {
			parents[j-2] = rng.Intn(j-1) + 1
		}
		tests = append(tests, buildCase(fmt.Sprintf("random-%d", i+1), append([]int{n}, parents...)))
	}
	return tests
}

func buildCase(name string, dataSlice []int) testCase {
	var sb strings.Builder
	n := dataSlice[0]
	fmt.Fprintf(&sb, "%d\n", n)
	for i := 2; i <= n; i++ {
		fmt.Fprintf(&sb, "%d ", dataSlice[i-1])
	}
	if n > 1 {
		sb.WriteByte('\n')
	}
	return testCase{name: name, input: sb.String()}
}
