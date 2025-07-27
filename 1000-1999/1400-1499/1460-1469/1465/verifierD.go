package main

import (
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"strings"
)

type testCaseD struct {
	s    string
	x, y int
}

func genTestsD() []testCaseD {
	rand.Seed(45)
	tests := make([]testCaseD, 100)
	for i := range tests {
		n := rand.Intn(6) + 1
		var sb strings.Builder
		for j := 0; j < n; j++ {
			r := rand.Intn(3)
			if r == 0 {
				sb.WriteByte('0')
			} else if r == 1 {
				sb.WriteByte('1')
			} else {
				sb.WriteByte('?')
			}
		}
		x := rand.Intn(5)
		y := rand.Intn(5)
		tests[i] = testCaseD{sb.String(), x, y}
	}
	tests = append(tests, testCaseD{"0?1", 2, 3})
	tests = append(tests, testCaseD{"?????", 13, 37})
	tests = append(tests, testCaseD{"?10?", 239, 7})
	tests = append(tests, testCaseD{"01101001", 5, 7})
	return tests
}

func cost(str string, x, y int) int64 {
	zeros, ones := 0, 0
	var c int64
	for i := 0; i < len(str); i++ {
		if str[i] == '0' {
			c += int64(y * ones)
			zeros++
		} else {
			c += int64(x * zeros)
			ones++
		}
	}
	return c
}

func solveD(tc testCaseD) int64 {
	pos := []int{}
	for i, ch := range tc.s {
		if ch == '?' {
			pos = append(pos, i)
		}
	}
	best := int64(1<<63 - 1)
	total := 1 << len(pos)
	buf := []byte(tc.s)
	for mask := 0; mask < total; mask++ {
		for j, p := range pos {
			if mask>>j&1 == 1 {
				buf[p] = '1'
			} else {
				buf[p] = '0'
			}
		}
		c := cost(string(buf), tc.x, tc.y)
		if c < best {
			best = c
		}
	}
	return best
}

func runBinary(bin, input string) (string, error) {
	cmd := exec.Command(bin)
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	cmd.Env = append(os.Environ(), "GOMAXPROCS=1")
	if err := cmd.Run(); err != nil {
		return out.String(), fmt.Errorf("runtime error: %v", err)
	}
	return out.String(), nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("Usage: go run verifierD.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	tests := genTestsD()
	for i, tc := range tests {
		input := fmt.Sprintf("%s\n%d %d\n", tc.s, tc.x, tc.y)
		exp := solveD(tc)
		got, err := runBinary(bin, input)
		if err != nil {
			fmt.Printf("test %d: %v\noutput:\n%s", i+1, err, got)
			os.Exit(1)
		}
		got = strings.TrimSpace(got)
		if got != fmt.Sprintf("%d", exp) {
			fmt.Printf("test %d failed: expected %d got %s\ninput:\n%s", i+1, exp, got, input)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
