package main

import (
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"strconv"
	"strings"
)

type testCaseC struct {
	n   int
	arr []int
}

func genTestsC() []testCaseC {
	rand.Seed(3)
	tests := make([]testCaseC, 100)
	for i := range tests {
		n := rand.Intn(6) + 1
		arr := make([]int, n)
		for j := range arr {
			arr[j] = rand.Intn(10) + 1
		}
		tests[i] = testCaseC{n, arr}
	}
	return tests
}

// solver copied from 671C.go
func solveC(tc testCaseC) int64 {
	n := tc.n
	arr := tc.arr
	maxV := 0
	for _, v := range arr {
		if v > maxV {
			maxV = v
		}
	}
	first1 := make([]int, maxV+1)
	first2 := make([]int, maxV+1)
	last1 := make([]int, maxV+1)
	last2 := make([]int, maxV+1)
	cnt := make([]int, maxV+1)
	for i := 0; i <= maxV; i++ {
		first1[i] = n + 1
		first2[i] = n + 1
	}
	for idx, v := range arr {
		pos := idx + 1
		for d := 1; d*d <= v; d++ {
			if v%d == 0 {
				update := func(div int) {
					cnt[div]++
					if pos < first1[div] {
						first2[div] = first1[div]
						first1[div] = pos
					} else if pos < first2[div] {
						first2[div] = pos
					}
					if pos > last1[div] {
						last2[div] = last1[div]
						last1[div] = pos
					} else if pos > last2[div] {
						last2[div] = pos
					}
				}
				update(d)
				if d*d != v {
					update(v / d)
				}
			}
		}
	}
	total := int64(n) * int64(n+1) / 2
	S := make([]int64, maxV+1)
	for g := 1; g <= maxV; g++ {
		if cnt[g] >= 2 {
			c := int64(first2[g])*(int64(n)-int64(last1[g])+1) + int64(first1[g])*(int64(last1[g])-int64(last2[g]))
			S[g] = total - c
		}
	}
	F := make([]int64, maxV+1)
	var ans int64
	for g := maxV; g >= 1; g-- {
		val := S[g]
		for m := g * 2; m <= maxV; m += g {
			val -= F[m]
		}
		if val < 0 {
			val = 0
		}
		F[g] = val
		ans += val * int64(g)
	}
	return ans
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
	tests := genTestsC()
	for i, tc := range tests {
		var input bytes.Buffer
		fmt.Fprintf(&input, "%d\n", tc.n)
		for j, v := range tc.arr {
			if j > 0 {
				input.WriteByte(' ')
			}
			fmt.Fprint(&input, v)
		}
		input.WriteByte('\n')
		expect := solveC(tc)
		out, err := run(bin, input.String())
		if err != nil {
			fmt.Fprintf(os.Stderr, "test %d: %v\n", i+1, err)
			os.Exit(1)
		}
		val, err := strconv.ParseInt(strings.TrimSpace(out), 10, 64)
		if err != nil {
			fmt.Fprintf(os.Stderr, "test %d: invalid output %q\n", i+1, out)
			os.Exit(1)
		}
		if val != expect {
			fmt.Fprintf(os.Stderr, "test %d: expected %d got %d\n", i+1, expect, val)
			os.Exit(1)
		}
	}
	fmt.Printf("All %d tests passed\n", len(tests))
}
