package main

import (
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"time"
)

func runBinary(bin, input string) (string, error) {
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

// solveCase201E ports the accepted C++ binary search oracle used by candidates.
func solveCase201E(n, m int64) int64 {
	l, r := int64(0), n
	for l <= r {
		s := (l + r) >> 1
		var num, sum int64 = 0, 0
		var c int64 = 1 // C(s,0)
		found := false
		for i := int64(0); i <= s; i++ {
			if i > 0 {
				c = c * (s - i + 1) / i // C(s,i)
			}
			num += c
			sum += c * i
			if num >= n {
				if sum-(num-n)*i <= s*m {
					r = s - 1
				} else {
					l = s + 1
				}
				found = true
				break
			}
		}
		if !found {
			l = s + 1
		}
	}
	return l
}

func generateCase(rng *rand.Rand) string {
	t := rng.Intn(5) + 1
	var b strings.Builder
	fmt.Fprintf(&b, "%d\n", t)
	for i := 0; i < t; i++ {
		n := rng.Intn(20) + 1
		m := rng.Intn(20) + 1
		fmt.Fprintf(&b, "%d %d\n", int64(n), int64(m))
	}
	return b.String()
}

func parseAllInts(s string) ([]int64, error) {
	fs := strings.Fields(s)
	res := make([]int64, 0, len(fs))
	for _, f := range fs {
		v, err := strconv.ParseInt(f, 10, 64)
		if err != nil {
			return nil, err
		}
		res = append(res, v)
	}
	return res, nil
}

func runAndCheck(bin string, input string) error {
	gotRaw, err := runBinary(bin, input)
	if err != nil {
		return err
	}
	vals, err := parseAllInts(gotRaw)
	if err != nil {
		return fmt.Errorf("cannot parse output: %v (output=%q)", err, gotRaw)
	}
	inVals, err := parseAllInts(input)
	if err != nil || len(inVals) < 1 {
		return fmt.Errorf("internal parse error on input")
	}
	T := int(inVals[0])
	if len(vals) != T {
		return fmt.Errorf("wrong number of lines: expected %d got %d\nexpected:\n<computed>\ngot:\n%s", T, len(vals), strings.TrimSpace(gotRaw))
	}
	idx := 1
	var expLines []string
	for i := 0; i < T; i++ {
		n := inVals[idx]
		m := inVals[idx+1]
		idx += 2
		exp := solveCase201E(n, m)
		if vals[i] != exp {
			// build expected output string
			expVals := make([]string, 0, T)
			idx2 := 1
			for j := 0; j < T; j++ {
				nj := inVals[idx2]
				mj := inVals[idx2+1]
				idx2 += 2
				expVals = append(expVals, fmt.Sprint(solveCase201E(nj, mj)))
			}
			return fmt.Errorf("wrong answer on case %d\nexpected:\n%s\n\ngot:\n%s\ninput:\n%s", i+1, strings.Join(expVals, "\n"), strings.TrimSpace(gotRaw), input)
		}
		expLines = append(expLines, fmt.Sprint(exp))
	}
	return nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierE.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < 100; i++ {
		input := generateCase(rng)
		if err := runAndCheck(bin, input); err != nil {
			fmt.Fprintf(os.Stderr, "case %d wrong answer\n%v\n", i+1, err)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
