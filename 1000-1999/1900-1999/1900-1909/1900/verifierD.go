package main

import (
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"sort"
	"strings"
	"time"
)

type testCase struct {
	input    string
	expected string
}

const maxV = 100000

var phi []int64

func initPhi() {
	phi = make([]int64, maxV+1)
	for i := 0; i <= maxV; i++ {
		phi[i] = int64(i)
	}
	for i := 2; i <= maxV; i++ {
		if phi[i] == int64(i) {
			for j := i; j <= maxV; j += i {
				phi[j] -= phi[j] / int64(i)
			}
		}
	}
}

func divisors(x int) []int {
	res := []int{}
	for d := 1; d*d <= x; d++ {
		if x%d == 0 {
			res = append(res, d)
			if d*d != x {
				res = append(res, x/d)
			}
		}
	}
	return res
}

func solveOne(arr []int) int64 {
	n := len(arr)
	sort.Ints(arr)
	freq := make([]int64, maxV+1)
	ans := int64(0)
	for j, val := range arr {
		ds := divisors(val)
		total := int64(0)
		for _, d := range ds {
			total += phi[d] * freq[d]
		}
		ans += total * int64(n-j-1)
		for _, d := range ds {
			freq[d]++
		}
	}
	return ans
}

func generateCase(rng *rand.Rand) testCase {
	t := rng.Intn(3) + 1
	var in strings.Builder
	var out strings.Builder
	in.WriteString(fmt.Sprintf("%d\n", t))
	for i := 0; i < t; i++ {
		n := rng.Intn(8) + 3
		arr := make([]int, n)
		for j := 0; j < n; j++ {
			arr[j] = rng.Intn(1000) + 1
		}
		in.WriteString(fmt.Sprintf("%d\n", n))
		for j := 0; j < n; j++ {
			if j > 0 {
				in.WriteByte(' ')
			}
			in.WriteString(fmt.Sprintf("%d", arr[j]))
		}
		in.WriteByte('\n')
		out.WriteString(fmt.Sprintf("%d\n", solveOne(arr)))
	}
	return testCase{input: in.String(), expected: out.String()}
}

func runCase(bin string, tc testCase) error {
	cmd := exec.Command(bin)
	cmd.Stdin = strings.NewReader(tc.input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("runtime error: %v\n%s", err, out.String())
	}
	got := strings.TrimSpace(out.String())
	exp := strings.TrimSpace(tc.expected)
	if got != exp {
		return fmt.Errorf("expected:\n%s\ngot:\n%s", exp, got)
	}
	return nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierD.go /path/to/binary")
		os.Exit(1)
	}
	initPhi()
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	bin := os.Args[1]
	// simple deterministic test
	arr := []int{2, 3, 6}
	in := "1\n3\n2 3 6\n"
	out := fmt.Sprintf("%d\n", solveOne(arr))
	cases := []testCase{{input: in, expected: out}}
	for i := 0; i < 100; i++ {
		cases = append(cases, generateCase(rng))
	}
	for i, tc := range cases {
		if err := runCase(bin, tc); err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:\n%s", i+1, err, tc.input)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
