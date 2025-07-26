package main

import (
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"sort"
	"strconv"
	"strings"
)

type Test struct {
	n, m int
	x    []int
	t    []int
}

func (tc Test) Input() string {
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("%d %d\n", tc.n, tc.m))
	for i, v := range tc.x {
		if i > 0 {
			sb.WriteByte(' ')
		}
		sb.WriteString(strconv.Itoa(v))
	}
	sb.WriteByte('\n')
	for i, v := range tc.t {
		if i > 0 {
			sb.WriteByte(' ')
		}
		sb.WriteString(strconv.Itoa(v))
	}
	sb.WriteByte('\n')
	return sb.String()
}

func runExe(path, input string) (string, error) {
	cmd := exec.Command(path)
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	err := cmd.Run()
	return out.String(), err
}

func expected(tc Test) []int {
	drivers := make([]int, 0, tc.m)
	for i, v := range tc.t {
		if v == 1 {
			drivers = append(drivers, tc.x[i])
		}
	}
	res := make([]int, tc.m)
	for i, v := range tc.t {
		if v == 0 {
			xi := tc.x[i]
			idx := sort.SearchInts(drivers, xi)
			if idx == 0 {
				res[0]++
			} else if idx == len(drivers) {
				res[len(drivers)-1]++
			} else {
				l := xi - drivers[idx-1]
				r := drivers[idx] - xi
				if l <= r {
					res[idx-1]++
				} else {
					res[idx]++
				}
			}
		}
	}
	return res
}

func parseOutput(out string, m int) ([]int, error) {
	fields := strings.Fields(out)
	if len(fields) != m {
		return nil, fmt.Errorf("expected %d numbers, got %d", m, len(fields))
	}
	arr := make([]int, m)
	for i := 0; i < m; i++ {
		v, err := strconv.Atoi(fields[i])
		if err != nil {
			return nil, fmt.Errorf("invalid integer")
		}
		arr[i] = v
	}
	return arr, nil
}

func genTests() []Test {
	rand.Seed(1)
	tests := make([]Test, 0, 105)
	for i := 0; i < 100; i++ {
		n := rand.Intn(6) + 1
		m := rand.Intn(6) + 1
		tot := n + m
		coordsSet := map[int]struct{}{}
		coords := make([]int, 0, tot)
		for len(coords) < tot {
			v := rand.Intn(1000) + 1
			if _, ok := coordsSet[v]; ok {
				continue
			}
			coordsSet[v] = struct{}{}
			coords = append(coords, v)
		}
		sort.Ints(coords)
		t := make([]int, tot)
		idxs := rand.Perm(tot)[:m]
		for _, id := range idxs {
			t[id] = 1
		}
		tests = append(tests, Test{n: n, m: m, x: coords, t: t})
	}
	tests = append(tests,
		Test{n: 1, m: 1, x: []int{1, 2}, t: []int{0, 1}},
		Test{n: 2, m: 1, x: []int{1, 2, 3}, t: []int{0, 0, 1}},
		Test{n: 3, m: 2, x: []int{1, 5, 6, 10, 11}, t: []int{0, 0, 1, 0, 1}},
		Test{n: 1, m: 4, x: []int{1, 2, 3, 4, 5}, t: []int{0, 1, 1, 1, 1}},
		Test{n: 4, m: 1, x: []int{10, 20, 30, 40, 50}, t: []int{0, 0, 1, 0, 0}},
	)
	return tests
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("Usage: go run verifierB.go /path/to/binary")
		return
	}
	bin := os.Args[1]
	tests := genTests()
	for i, tc := range tests {
		input := tc.Input()
		out, err := runExe(bin, input)
		if err != nil {
			fmt.Printf("runtime error on test %d: %v\n", i+1, err)
			os.Exit(1)
		}
		got, err := parseOutput(out, tc.m)
		if err != nil {
			fmt.Printf("output parse failed on test %d: %v\n", i+1, err)
			os.Exit(1)
		}
		exp := expected(tc)
		for j := 0; j < tc.m; j++ {
			if got[j] != exp[j] {
				fmt.Printf("Test %d failed\nInput:%sExpected:%v\nGot:%v\n", i+1, input, exp, got)
				os.Exit(1)
			}
		}
	}
	fmt.Println("all tests passed")
}
