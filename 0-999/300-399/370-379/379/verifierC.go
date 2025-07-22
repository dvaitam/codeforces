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

type pair struct {
	val int
	idx int
}

func minimal(a []int) ([]int64, int64) {
	n := len(a)
	P := make([]pair, n)
	for i := 0; i < n; i++ {
		P[i] = pair{a[i], i}
	}
	sort.Slice(P, func(i, j int) bool { return P[i].val < P[j].val })
	res := make([]int64, n)
	var curr int64
	var sum int64
	for _, p := range P {
		ai := int64(p.val)
		if ai > curr {
			curr = ai
		} else {
			curr++
		}
		res[p.idx] = curr
		sum += curr
	}
	return res, sum
}

func runCase(exe string, a []int) error {
	n := len(a)
	input := fmt.Sprintf("%d\n", n)
	for i, v := range a {
		if i > 0 {
			input += " "
		}
		input += fmt.Sprintf("%d", v)
	}
	input += "\n"
	cmd := exec.Command(exe)
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("runtime error: %v\n%s", err, out.String())
	}
	fields := strings.Fields(out.String())
	if len(fields) != n {
		return fmt.Errorf("expected %d numbers got %d", n, len(fields))
	}
	res := make([]int64, n)
	for i, f := range fields {
		var v int64
		if _, err := fmt.Sscan(f, &v); err != nil {
			return fmt.Errorf("parse error: %v", err)
		}
		res[i] = v
	}
	_, expSum := minimal(a)
	seen := make(map[int64]bool)
	var sum int64
	for i := 0; i < n; i++ {
		if res[i] < int64(a[i]) {
			return fmt.Errorf("value %d less than a[i]", res[i])
		}
		if seen[res[i]] {
			return fmt.Errorf("values not distinct")
		}
		seen[res[i]] = true
		sum += res[i]
	}
	if sum != expSum {
		return fmt.Errorf("sum not minimal")
	}
	return nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierC.go /path/to/binary")
		os.Exit(1)
	}
	exe := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	cases := [][]int{{1, 1}, {5, 3, 3}, {1, 2, 3}}
	for i := 0; i < 100; i++ {
		n := rng.Intn(6) + 1
		arr := make([]int, n)
		for j := 0; j < n; j++ {
			arr[j] = rng.Intn(10) + 1
		}
		cases = append(cases, arr)
	}
	for idx, arr := range cases {
		if err := runCase(exe, arr); err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\n", idx+1, err)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
