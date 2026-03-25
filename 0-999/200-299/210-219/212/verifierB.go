package main

import (
	"bufio"
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
	s    string
	sets []string
}

const maxStringLen = 100000

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierB.go /path/to/binary")
		os.Exit(1)
	}
	candidate := os.Args[1]

	tests := buildTests()

	// Compute expected answers using embedded solver
	expected := solveAll(tests)

	// Run candidate on each test case individually (the problem is single-case)
	got := make([]int64, 0, len(expected))
	for _, tc := range tests {
		input := tc.s + "\n" + strconv.Itoa(len(tc.sets)) + "\n"
		for _, c := range tc.sets {
			input += c + "\n"
		}
		out, err := runProgram(candidate, input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "candidate failed: %v\n", err)
			os.Exit(1)
		}
		fields := strings.Fields(out)
		if len(fields) != len(tc.sets) {
			fmt.Fprintf(os.Stderr, "expected %d answers, got %d (output: %q)\n", len(tc.sets), len(fields), out)
			os.Exit(1)
		}
		for _, f := range fields {
			val, err := strconv.ParseInt(f, 10, 64)
			if err != nil {
				fmt.Fprintf(os.Stderr, "invalid integer %q: %v\n", f, err)
				os.Exit(1)
			}
			got = append(got, val)
		}
	}

	idx := 0
	for ti, tc := range tests {
		for qi := range tc.sets {
			if expected[idx] != got[idx] {
				fmt.Fprintf(os.Stderr, "Mismatch on test %d query %d: expected %d got %d\ns=%s set=%s\n",
					ti+1, qi+1, expected[idx], got[idx], tc.s, tc.sets[qi])
				os.Exit(1)
			}
			idx++
		}
	}
	fmt.Println("All tests passed")
}

// solveAll uses the embedded CF-accepted solver logic to compute all answers.
func solveAll(tests []testCase) []int64 {
	var results []int64
	for _, tc := range tests {
		results = append(results, solveOne(tc)...)
	}
	return results
}

// solveOne is the embedded CF-accepted solver for 212B.
func solveOne(tc testCase) []int64 {
	s := []byte(tc.s)
	n := len(s)
	if n == 0 {
		res := make([]int64, len(tc.sets))
		return res
	}

	numWords := (n + 63) / 64
	bitsets := make([][]uint64, 26)
	for c := 0; c < 26; c++ {
		bitsets[c] = make([]uint64, numWords)
	}
	for i, b := range s {
		c := b - 'a'
		word := i / 64
		bit := i % 64
		bitsets[c][word] |= uint64(1) << bit
	}

	nextPos := make([][26]int32, n+1)
	for c := 0; c < 26; c++ {
		nextPos[n][c] = int32(n)
	}
	for i := n - 1; i >= 0; i-- {
		nextPos[i] = nextPos[i+1]
		nextPos[i][s[i]-'a'] = int32(i)
	}

	valid := make([]uint64, numWords)
	charsInMask := make([]int, 0, 26)
	memo := make(map[uint32]int)

	results := make([]int64, len(tc.sets))

	for qi, qstr := range tc.sets {
		q := []byte(qstr)
		maskC := uint32(0)
		for _, b := range q {
			if b >= 'a' && b <= 'z' {
				maskC |= 1 << (b - 'a')
			}
		}

		if val, ok := memo[maskC]; ok {
			results[qi] = int64(val)
			continue
		}

		charsInMask = charsInMask[:0]
		for c := 0; c < 26; c++ {
			if (maskC & (1 << c)) != 0 {
				charsInMask = append(charsInMask, c)
			}
		}
		numChars := len(charsInMask)

		for j := 0; j < numWords; j++ {
			valid[j] = 0
		}
		for _, c := range charsInMask {
			bc := bitsets[c]
			for j := 0; j < numWords; j++ {
				valid[j] |= bc[j]
			}
		}

		if n%64 != 0 {
			valid[numWords-1] &= (uint64(1) << (n % 64)) - 1
		}

		count := 0
		inBlock := false
		l := 0

		for j := 0; j < numWords; j++ {
			v := valid[j]
			b := 0
			for b < 64 {
				if inBlock {
					inv := (^v) & (^uint64(0) << b)
					if inv == 0 {
						break
					}
					tz := bits.TrailingZeros64(inv)
					r := j*64 + tz - 1
					if r-l+1 >= numChars {
						ok := true
						for _, c := range charsInMask {
							if nextPos[l][c] > int32(r) {
								ok = false
								break
							}
						}
						if ok {
							count++
						}
					}
					inBlock = false
					b = tz + 1
				} else {
					rem := v & (^uint64(0) << b)
					if rem == 0 {
						break
					}
					tz := bits.TrailingZeros64(rem)
					l = j*64 + tz
					inBlock = true
					b = tz + 1
				}
			}
		}

		if inBlock {
			r := n - 1
			if r-l+1 >= numChars {
				ok := true
				for _, c := range charsInMask {
					if nextPos[l][c] > int32(r) {
						ok = false
						break
					}
				}
				if ok {
					count++
				}
			}
		}

		memo[maskC] = count
		results[qi] = int64(count)
	}
	return results
}

func runProgram(target, input string) (string, error) {
	var cmd *exec.Cmd
	if strings.HasSuffix(target, ".go") {
		cmd = exec.Command("go", "run", target)
	} else {
		cmd = exec.Command(target)
	}
	cmd.Stdin = strings.NewReader(input)
	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("runtime error: %v\nstdout:\n%s\nstderr:\n%s", err, stdout.String(), stderr.String())
	}
	return stdout.String(), nil
}

func buildTests() []testCase {
	tests := deterministicTests()
	totalLen := totalStringLen(tests)
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))

	alphabet := "abcdefghijklmnopqrstuvwxyz"
	for totalLen < maxStringLen {
		length := rng.Intn(1000) + 1
		if totalLen+length > maxStringLen {
			length = maxStringLen - totalLen
		}
		var sb strings.Builder
		for i := 0; i < length; i++ {
			sb.WriteByte(alphabet[rng.Intn(len(alphabet))])
		}
		s := sb.String()
		q := rng.Intn(20) + 1
		setList := make([]string, q)
		for i := 0; i < q; i++ {
			setSize := rng.Intn(5) + 1
			mask := make(map[byte]struct{})
			var str strings.Builder
			for str.Len() < setSize {
				ch := alphabet[rng.Intn(len(alphabet))]
				if _, ok := mask[ch]; !ok {
					mask[ch] = struct{}{}
					str.WriteByte(ch)
				}
			}
			setList[i] = str.String()
		}
		tests = append(tests, testCase{s: s, sets: setList})
		totalLen += len(s)
	}
	return tests
}

func deterministicTests() []testCase {
	return []testCase{
		{s: "aaaaa", sets: []string{"a", "ab"}},
		{s: "abacaba", sets: []string{"ac", "ba", "abc"}},
		{s: "xyz", sets: []string{"x", "yz"}},
	}
}

func totalStringLen(tests []testCase) int {
	total := 0
	for _, tc := range tests {
		total += len(tc.s)
	}
	return total
}

// suppress unused import
var _ = bufio.NewReader
