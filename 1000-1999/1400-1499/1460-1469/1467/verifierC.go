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

type testCaseC struct {
	n1, n2, n3 int
	a1, a2, a3 []int64
}

func solveC(tc testCaseC) int64 {
	sums := []int64{0, 0, 0}
	mins := []int64{1<<63 - 1, 1<<63 - 1, 1<<63 - 1}
	arrs := [][]int64{tc.a1, tc.a2, tc.a3}
	for i := 0; i < 3; i++ {
		for _, v := range arrs[i] {
			sums[i] += v
			if v < mins[i] {
				mins[i] = v
			}
		}
	}
	total := sums[0] + sums[1] + sums[2]
	cand := sums[0]
	if sums[1] < cand {
		cand = sums[1]
	}
	if sums[2] < cand {
		cand = sums[2]
	}
	pair := mins[0] + mins[1]
	if pair < cand {
		cand = pair
	}
	pair = mins[0] + mins[2]
	if pair < cand {
		cand = pair
	}
	pair = mins[1] + mins[2]
	if pair < cand {
		cand = pair
	}
	return total - 2*cand
}

func buildInputC(tc testCaseC) string {
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("%d %d %d\n", tc.n1, tc.n2, tc.n3))
	for i, arr := range [][]int64{tc.a1, tc.a2, tc.a3} {
		for j, v := range arr {
			if j > 0 {
				sb.WriteByte(' ')
			}
			sb.WriteString(fmt.Sprint(v))
		}
		if i < 2 {
			sb.WriteByte('\n')
		} else {
			sb.WriteByte('\n')
		}
	}
	return sb.String()
}

func runCaseC(bin string, tc testCaseC) error {
	cmd := exec.Command(bin)
	cmd.Stdin = strings.NewReader(buildInputC(tc))
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("runtime error: %v\n%s", err, out.String())
	}
	var got int64
	if _, err := fmt.Fscan(strings.NewReader(out.String()), &got); err != nil {
		return fmt.Errorf("bad output: %v", err)
	}
	expect := solveC(tc)
	if got != expect {
		return fmt.Errorf("expected %d got %d", expect, got)
	}
	return nil
}

func generateCasesC() []testCaseC {
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	cases := make([]testCaseC, 0, 100)
	for len(cases) < 100 {
		n1 := rng.Intn(5) + 1
		n2 := rng.Intn(5) + 1
		n3 := rng.Intn(5) + 1
		a1 := make([]int64, n1)
		a2 := make([]int64, n2)
		a3 := make([]int64, n3)
		for i := 0; i < n1; i++ {
			a1[i] = int64(rng.Intn(20) + 1)
		}
		for i := 0; i < n2; i++ {
			a2[i] = int64(rng.Intn(20) + 1)
		}
		for i := 0; i < n3; i++ {
			a3[i] = int64(rng.Intn(20) + 1)
		}
		cases = append(cases, testCaseC{n1, n2, n3, a1, a2, a3})
	}
	return cases
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierC.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	cases := generateCasesC()
	for i, tc := range cases {
		if err := runCaseC(bin, tc); err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\n", i+1, err)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
