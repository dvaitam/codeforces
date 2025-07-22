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

type testC struct {
	n   int
	k   int
	p   int
	arr []int
}

func feasible(tc testC) bool {
	odd := 0
	for _, v := range tc.arr {
		if v%2 != 0 {
			odd++
		}
	}
	even := tc.n - odd
	needOdd := tc.k - tc.p
	if odd < needOdd {
		return false
	}
	if (odd-needOdd)%2 != 0 {
		return false
	}
	if (odd-needOdd)/2+even < tc.p {
		return false
	}
	return true
}

func genC(rng *rand.Rand) testC {
	n := rng.Intn(20) + 1
	k := rng.Intn(n) + 1
	p := rng.Intn(k + 1)
	arr := make([]int, n)
	used := map[int]bool{}
	for i := 0; i < n; i++ {
		for {
			v := rng.Intn(1000) + 1
			if !used[v] {
				arr[i] = v
				used[v] = true
				break
			}
		}
	}
	return testC{n: n, k: k, p: p, arr: arr}
}

func runCase(bin string, tc testC) (string, error) {
	var sb strings.Builder
	fmt.Fprintf(&sb, "%d %d %d\n", tc.n, tc.k, tc.p)
	for i, v := range tc.arr {
		if i > 0 {
			sb.WriteByte(' ')
		}
		fmt.Fprint(&sb, v)
	}
	sb.WriteByte('\n')
	cmd := exec.Command(bin)
	cmd.Stdin = strings.NewReader(sb.String())
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	err := cmd.Run()
	return strings.TrimSpace(out.String()), err
}

func checkOutput(out string, tc testC) error {
	lines := strings.Split(strings.TrimSpace(out), "\n")
	if len(lines) == 0 {
		return fmt.Errorf("no output")
	}
	ans := strings.TrimSpace(lines[0])
	possible := feasible(tc)
	if ans == "NO" {
		if possible {
			return fmt.Errorf("should be YES")
		}
		return nil
	}
	if ans != "YES" {
		return fmt.Errorf("first line must be YES or NO")
	}
	if !possible {
		return fmt.Errorf("output YES but impossible")
	}
	if len(lines)-1 != tc.k {
		return fmt.Errorf("expected %d parts, got %d", tc.k, len(lines)-1)
	}
	used := make(map[int]bool)
	evenParts := 0
	for i := 1; i <= tc.k; i++ {
		f := strings.Fields(lines[i])
		if len(f) < 1 {
			return fmt.Errorf("line %d empty", i+1)
		}
		cnt, err := strconv.Atoi(f[0])
		if err != nil {
			return fmt.Errorf("line %d bad count", i+1)
		}
		if cnt != len(f)-1 {
			return fmt.Errorf("line %d count mismatch", i+1)
		}
		if cnt == 0 {
			return fmt.Errorf("line %d empty part", i+1)
		}
		sum := 0
		for j := 1; j < len(f); j++ {
			v, err := strconv.Atoi(f[j])
			if err != nil {
				return fmt.Errorf("line %d bad number", i+1)
			}
			present := false
			for _, a := range tc.arr {
				if a == v {
					present = true
					break
				}
			}
			if !present {
				return fmt.Errorf("line %d uses unknown value", i+1)
			}
			if used[v] {
				return fmt.Errorf("value %d repeated", v)
			}
			used[v] = true
			sum += v
		}
		if sum%2 == 0 {
			evenParts++
		}
	}
	if len(used) != tc.n {
		return fmt.Errorf("not all values used")
	}
	if evenParts != tc.p {
		return fmt.Errorf("expected %d even parts got %d", tc.p, evenParts)
	}
	return nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierC.go /path/to/binary")
		os.Exit(1)
	}
	exe := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	cases := []testC{}
	cases = append(cases, testC{n: 1, k: 1, p: 1, arr: []int{2}})
	cases = append(cases, testC{n: 3, k: 2, p: 1, arr: []int{1, 2, 3}})
	for i := 0; i < 100; i++ {
		cases = append(cases, genC(rng))
	}
	for i, tc := range cases {
		out, err := runCase(exe, tc)
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d: runtime error: %v\n", i+1, err)
			os.Exit(1)
		}
		if err := checkOutput(out, tc); err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\n", i+1, err)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
