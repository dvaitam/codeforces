package main

import (
	"bytes"
	"fmt"
	"math"
	"math/big"
	"math/rand"
	"os"
	"os/exec"
	"sort"
	"strings"
	"time"
)

type testCase struct {
	n   int
	arr []int64
}

func solve(tc testCase) string {
	a := append([]int64(nil), tc.arr...)
	sort.Slice(a, func(i, j int) bool { return a[i] < a[j] })
	n := tc.n
	best := int64(0)
	for i := 0; i < n; i++ {
		diff := a[i] - 1
		if diff < 0 {
			diff = -diff
		}
		best += diff
	}
	if n > 1 {
		maxPower := best + a[n-1]
		maxC := int(math.Pow(float64(maxPower), 1.0/float64(n-1))) + 1
		for c := 2; c <= maxC; c++ {
			cost := int64(0)
			cur := big.NewInt(1)
			bigC := big.NewInt(int64(c))
			valid := true
			for i := 0; i < n; i++ {
				if cur.BitLen() > 62 {
					valid = false
					break
				}
				tar := cur.Int64()
				diff := a[i] - tar
				if diff < 0 {
					diff = -diff
				}
				cost += diff
				if cost > best {
					valid = false
					break
				}
				if i < n-1 {
					cur.Mul(cur, bigC)
				}
			}
			if valid && cost < best {
				best = cost
			}
		}
	}
	return fmt.Sprint(best)
}

func (tc testCase) input() string {
	var b strings.Builder
	b.WriteString(fmt.Sprintf("%d\n", tc.n))
	for i, v := range tc.arr {
		if i > 0 {
			b.WriteByte(' ')
		}
		b.WriteString(fmt.Sprintf("%d", v))
	}
	b.WriteByte('\n')
	return b.String()
}

func randomCase(rng *rand.Rand) testCase {
	n := rng.Intn(6) + 3
	arr := make([]int64, n)
	for i := 0; i < n; i++ {
		arr[i] = rng.Int63n(20) + 1
	}
	return testCase{n: n, arr: arr}
}

func runProgram(bin, input string) (string, error) {
	var cmd *exec.Cmd
	if strings.HasSuffix(bin, ".go") {
		cmd = exec.Command("go", "run", bin)
	} else {
		cmd = exec.Command(bin)
	}
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("runtime error: %v\n%s", err, out.String())
	}
	return strings.TrimSpace(out.String()), nil
}

func runCase(bin string, tc testCase) error {
	in := tc.input()
	expected := solve(tc)
	got, err := runProgram(bin, in)
	if err != nil {
		return err
	}
	if strings.TrimSpace(got) != expected {
		return fmt.Errorf("expected %s got %s", expected, got)
	}
	return nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierB.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	cases := []testCase{{n: 3, arr: []int64{1, 2, 3}}}
	for len(cases) < 100 {
		cases = append(cases, randomCase(rng))
	}
	for i, tc := range cases {
		if err := runCase(bin, tc); err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:\n%s", i+1, err, tc.input())
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
