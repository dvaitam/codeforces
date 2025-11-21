package main

import (
	"bufio"
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"strings"
)

func solveRef(input string) string {
	reader := bufio.NewReader(strings.NewReader(input))
	var b, d int
	var a, c string
	if _, err := fmt.Fscan(reader, &b, &d); err != nil {
		return ""
	}
	fmt.Fscan(reader, &a)
	fmt.Fscan(reader, &c)
	aLen := len(a)
	cLen := len(c)
	count := make([]int, cLen)
	nextIdx := make([]int, cLen)
	for start := 0; start < cLen; start++ {
		j := start
		cnt := 0
		for i := 0; i < aLen; i++ {
			if a[i] == c[j] {
				j++
				if j == cLen {
					j = 0
					cnt++
				}
			}
		}
		count[start] = cnt
		nextIdx[start] = j
	}
	firstIdx := make([]int, cLen)
	firstCount := make([]int64, cLen)
	cur := 0
	var total int64
	for rep := 1; rep <= b; rep++ {
		if firstIdx[cur] != 0 {
			prevRep := firstIdx[cur]
			prevCount := firstCount[cur]
			cycleLen := rep - prevRep
			cycleGain := total - prevCount
			remaining := b - rep + 1
			cycles := remaining / cycleLen
			total += int64(cycles) * cycleGain
			left := remaining % cycleLen
			for i := 0; i < left; i++ {
				total += int64(count[cur])
				cur = nextIdx[cur]
			}
			return fmt.Sprintf("%d", total/int64(d))
		}
		firstIdx[cur] = rep
		firstCount[cur] = total
		total += int64(count[cur])
		cur = nextIdx[cur]
	}
	return fmt.Sprintf("%d", total/int64(d))
}

type testCase struct {
	name   string
	input  string
	expect string
}

func makeCase(name string, b, d int, a, c string) testCase {
	input := fmt.Sprintf("%d %d\n%s\n%s\n", b, d, a, c)
	return testCase{name: name, input: input, expect: solveRef(input)}
}

func handcraftedTests() []testCase {
	return []testCase{
		makeCase("single_char_match", 1, 1, "a", "a"),
		makeCase("single_char_mismatch", 1, 1, "a", "b"),
		makeCase("simple_repeat", 4, 2, "ab", "a"),
		makeCase("repeat_c_longer", 5, 1, "abc", "ac"),
		makeCase("no_possible", 10, 1, "bbbb", "a"),
		makeCase("multiple_cycles", 7, 2, "abcabc", "abc"),
		makeCase("partial_cycle", 3, 1, "aba", "aa"),
		makeCase("long_a_short_c", 8, 3, "abcabcabc", "abc"),
		makeCase("c_longer_than_a", 10, 2, "ab", "aba"),
		makeCase("mixed_chars", 12, 4, "abacaba", "aba"),
	}
}

func randomTests() []testCase {
	rng := rand.New(rand.NewSource(314))
	var tests []testCase
	gen := func(prefix string, count, maxLen int, maxBD int) {
		for i := 0; i < count; i++ {
			aLen := rng.Intn(maxLen) + 1
			cLen := rng.Intn(maxLen) + 1
			a := make([]byte, aLen)
			c := make([]byte, cLen)
			for idx := range a {
				a[idx] = byte('a' + rng.Intn(4))
			}
			for idx := range c {
				c[idx] = byte('a' + rng.Intn(4))
			}
			b := rng.Intn(maxBD) + 1
			d := rng.Intn(maxBD) + 1
			tests = append(tests, makeCase(fmt.Sprintf("%s_%d", prefix, i+1), b, d, string(a), string(c)))
		}
	}
	gen("tiny", 150, 5, 10)
	gen("small", 120, 10, 100)
	return tests
}

func largeTests() []testCase {
	var tests []testCase
	tests = append(tests,
		makeCase("large_bd1", 1_000_000, 10, strings.Repeat("abc", 10), "abc"),
		makeCase("large_bd2", 10_000_000, 5, "abcd", "abcd"),
		makeCase("large_bd3", 5_000_000, 7, "abcabcabcabc", "cab"),
	)
	return tests
}

func runCandidate(bin, input string) (string, error) {
	cmd := exec.Command(bin)
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	err := cmd.Run()
	return strings.TrimSpace(out.String()), err
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("usage: go run verifierB.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	tests := append(handcraftedTests(), append(randomTests(), largeTests()...)...)
	for idx, tc := range tests {
		got, err := runCandidate(bin, tc.input)
		if err != nil {
			fmt.Printf("test %d (%s) runtime error: %v\ninput:\n%s", idx+1, tc.name, err, tc.input)
			os.Exit(1)
		}
		if got != tc.expect {
			fmt.Printf("test %d (%s) failed\ninput:\n%s\nexpect:%s\nactual:%s\n", idx+1, tc.name, tc.input, tc.expect, got)
			os.Exit(1)
		}
	}
	fmt.Println("all tests passed")
}
