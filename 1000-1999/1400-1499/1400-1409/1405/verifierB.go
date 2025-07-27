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

type testCase struct {
	arr []int64
}

type supply struct {
	idx int
	rem int64
}

func runCandidate(bin, input string) (string, error) {
	var cmd *exec.Cmd
	if strings.HasSuffix(bin, ".go") {
		cmd = exec.Command("go", "run", bin)
	} else {
		cmd = exec.Command(bin)
	}
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	var errBuf bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &errBuf
	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("runtime error: %v\n%s", err, errBuf.String())
	}
	return strings.TrimSpace(out.String()), nil
}

func solve(a []int64) int64 {
	var total int64
	for _, v := range a {
		if v > 0 {
			total += v
		}
	}
	supplies := make([]supply, 0, len(a))
	head := 0
	var free int64
	for i, v := range a {
		if v > 0 {
			supplies = append(supplies, supply{idx: i, rem: v})
		} else if v < 0 {
			need := -v
			for need > 0 && head < len(supplies) && supplies[head].idx < i {
				cur := &supplies[head]
				d := cur.rem
				if d > need {
					d = need
				}
				free += d
				cur.rem -= d
				need -= d
				if cur.rem == 0 {
					head++
				}
			}
		}
	}
	return total - free
}

func genCase(rng *rand.Rand) testCase {
	n := rng.Intn(50) + 1
	arr := make([]int64, n)
	var sum int64
	for i := 0; i < n-1; i++ {
		v := int64(rng.Intn(21) - 10)
		arr[i] = v
		sum += v
	}
	arr[n-1] = -sum
	return testCase{arr: arr}
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierB.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	cases := []testCase{{arr: []int64{0}}, {arr: []int64{1, -1}}, {arr: []int64{-1, 1, 0}}}
	for len(cases) < 100 {
		cases = append(cases, genCase(rng))
	}
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("%d\n", len(cases)))
	for _, tc := range cases {
		sb.WriteString(fmt.Sprintf("%d\n", len(tc.arr)))
		for i, v := range tc.arr {
			if i > 0 {
				sb.WriteByte(' ')
			}
			sb.WriteString(fmt.Sprint(v))
		}
		sb.WriteByte('\n')
	}
	out, err := runCandidate(bin, sb.String())
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	fields := strings.Fields(out)
	if len(fields) != len(cases) {
		fmt.Fprintf(os.Stderr, "expected %d numbers, got %d\n", len(cases), len(fields))
		os.Exit(1)
	}
	for i, tc := range cases {
		v, err := strconv.ParseInt(fields[i], 10, 64)
		if err != nil {
			fmt.Fprintf(os.Stderr, "invalid integer on case %d: %v\n", i+1, err)
			os.Exit(1)
		}
		expect := solve(tc.arr)
		if v != expect {
			fmt.Fprintf(os.Stderr, "case %d failed\ninput: %v\nexpected: %d got: %d\n", i+1, tc.arr, expect, v)
			os.Exit(1)
		}
	}
	fmt.Printf("All %d tests passed\n", len(cases))
}
