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

func solveCase(s string) string {
	n := len(s)
	nextDiff := make([]int, n)
	for i := n - 1; i >= 0; i-- {
		if i+1 < n && s[i] == s[i+1] {
			nextDiff[i] = nextDiff[i+1]
		} else {
			nextDiff[i] = i + 1
		}
	}
	isEmpty := make([]bool, n+1)
	headChar := make([]byte, n+1)
	firstDiff := make([]byte, n+1)
	headRunLen := make([]int, n+1)
	length := make([]int, n+1)
	pref := make([]string, n+1)
	suff := make([]string, n+1)
	isEmpty[n] = true
	length[n] = 0
	for i := n - 1; i >= 0; i-- {
		c := s[i]
		k := nextDiff[i] - i
		r := i + k
		var chooseSmall bool
		if r >= n || isEmpty[r] {
			chooseSmall = true
		} else {
			hc := headChar[r]
			if hc != c {
				chooseSmall = hc < c
			} else {
				fd := firstDiff[r]
				if fd == 0 {
					chooseSmall = true
				} else {
					chooseSmall = fd < c
				}
			}
		}
		tMin := k & 1
		tMax := k
		t := tMax
		if chooseSmall {
			t = tMin
		}
		length[i] = t
		if r <= n {
			length[i] += length[r]
		}
		isEmpty[i] = (length[i] == 0)
		maxPref := 5
		if length[i] <= maxPref {
			var tmp []byte
			for j := 0; j < t; j++ {
				tmp = append(tmp, c)
			}
			if !isEmpty[r] {
				want := length[i] - t
				curPref := pref[r]
				if len(curPref) > want {
					curPref = curPref[:want]
				}
				tmp = append(tmp, curPref...)
			}
			pref[i] = string(tmp)
		} else {
			if t >= maxPref {
				buf := make([]byte, maxPref)
				for j := range buf {
					buf[j] = c
				}
				pref[i] = string(buf)
			} else {
				var buf []byte
				for j := 0; j < t; j++ {
					buf = append(buf, c)
				}
				need := maxPref - t
				add := pref[r]
				if len(add) > need {
					add = add[:need]
				}
				buf = append(buf, add...)
				pref[i] = string(buf)
			}
		}
		if length[i] <= 2 {
			var tmp []byte
			for j := 0; j < t; j++ {
				tmp = append(tmp, c)
			}
			if !isEmpty[r] {
				add := suff[r]
				tmp = append(tmp, add...)
			}
			suff[i] = string(tmp)
		} else {
			if r < n && length[r] >= 2 {
				suff[i] = suff[r]
			} else if r < n && length[r] == 1 {
				tmp := []byte{c, suff[r][0]}
				suff[i] = string(tmp)
			} else {
				suff[i] = string([]byte{c, c})
			}
		}
		if t > 0 {
			headChar[i] = c
			headRunLen[i] = t
			if r >= n || isEmpty[r] {
				firstDiff[i] = 0
			} else if headChar[r] != c {
				firstDiff[i] = headChar[r]
			} else {
				firstDiff[i] = firstDiff[r]
			}
		} else {
			headChar[i] = headChar[r]
			headRunLen[i] = headRunLen[r]
			firstDiff[i] = firstDiff[r]
		}
	}
	var out strings.Builder
	for i := 0; i < n; i++ {
		L := length[i]
		out.WriteString(fmt.Sprintf("%d ", L))
		if L <= 10 {
			out.WriteString(pref[i])
		} else {
			out.WriteString(pref[i])
			out.WriteString("...")
			out.WriteString(suff[i])
		}
		if i+1 < n {
			out.WriteByte('\n')
		}
	}
	return out.String()
}

func generateTests() []test {
	rng := rand.New(rand.NewSource(46))
	tests := []test{}
	fixed := []string{"a", "aaaa", "abab"}
	for _, f := range fixed {
		tests = append(tests, test{f, solveCase(f)})
	}
	for len(tests) < 100 {
		n := rng.Intn(8) + 1
		var sb strings.Builder
		for i := 0; i < n; i++ {
			sb.WriteByte(byte('a' + rng.Intn(3)))
		}
		s := sb.String()
		tests = append(tests, test{s, solveCase(s)})
	}
	return tests
}

func runBinary(bin, input string) (string, error) {
	var cmd *exec.Cmd
	if strings.HasSuffix(bin, ".go") {
		cmd = exec.Command("go", "run", bin)
	} else {
		cmd = exec.Command(bin)
	}
	cmd.Stdin = strings.NewReader(input + "\n")
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	err := cmd.Run()
	return strings.TrimSpace(out.String()), err
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierE.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	tests := generateTests()
	for i, t := range tests {
		got, err := runBinary(bin, t.input)
		if err != nil {
			fmt.Printf("Runtime error on test %d: %v\n", i+1, err)
			os.Exit(1)
		}
		if strings.TrimSpace(got) != strings.TrimSpace(t.expected) {
			fmt.Printf("Wrong answer on test %d\nInput:%s\nExpected:%s\nGot:%s\n", i+1, t.input, t.expected, got)
			os.Exit(1)
		}
	}
	fmt.Printf("All %d tests passed.\n", len(tests))
}
