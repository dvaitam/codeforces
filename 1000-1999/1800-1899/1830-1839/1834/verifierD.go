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
	n int
	m int
	l []int
	r []int
}

func expectedD(tc testCaseD) string {
	n, m := tc.n, tc.m
	best := 0
	subsets := 1 << m
	for mask := 0; mask < subsets; mask++ {
		heights := make([]int, n)
		for topic := 1; topic <= m; topic++ {
			if mask&(1<<(topic-1)) == 0 {
				continue
			}
			for i := 0; i < n; i++ {
				if tc.l[i] <= topic && topic <= tc.r[i] {
					heights[i]++
				} else {
					heights[i]--
				}
			}
		}
		hi := heights[0]
		lo := heights[0]
		for _, v := range heights {
			if v > hi {
				hi = v
			}
			if v < lo {
				lo = v
			}
		}
		diff := hi - lo
		if diff > best {
			best = diff
		}
	}
	return fmt.Sprint(best)
}

func genTestsD() []testCaseD {
	rand.Seed(4)
	tests := make([]testCaseD, 0, 100)
	for len(tests) < 100 {
		n := rand.Intn(3) + 2 // 2..4
		m := rand.Intn(5) + 1 // 1..5
		l := make([]int, n)
		r := make([]int, n)
		for i := 0; i < n; i++ {
			a := rand.Intn(m) + 1
			b := rand.Intn(m) + 1
			if a > b {
				a, b = b, a
			}
			l[i] = a
			r[i] = b
		}
		tests = append(tests, testCaseD{n: n, m: m, l: l, r: r})
	}
	return tests
}

func runCase(bin string, tc testCaseD) error {
	var input strings.Builder
	input.WriteString("1\n")
	input.WriteString(fmt.Sprintf("%d %d\n", tc.n, tc.m))
	for i := 0; i < tc.n; i++ {
		input.WriteString(fmt.Sprintf("%d %d\n", tc.l[i], tc.r[i]))
	}
	cmd := exec.Command(bin)
	if strings.HasSuffix(bin, ".go") {
		cmd = exec.Command("go", "run", bin)
	}
	cmd.Stdin = strings.NewReader(input.String())
	var out bytes.Buffer
	var errBuf bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &errBuf
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("runtime error: %v\n%s", err, errBuf.String())
	}
	got := strings.TrimSpace(out.String())
	expect := expectedD(tc)
	if got != expect {
		return fmt.Errorf("expected %s got %s", expect, got)
	}
	return nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("Usage: go run verifierD.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	tests := genTestsD()
	for i, tc := range tests {
		if err := runCase(bin, tc); err != nil {
			fmt.Printf("case %d failed: %v\n", i+1, err)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
