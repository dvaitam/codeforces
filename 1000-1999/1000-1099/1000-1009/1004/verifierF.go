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

func runCandidate(bin, input string) (string, error) {
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
	err := cmd.Run()
	return strings.TrimSpace(out.String()), err
}

func countSub(arr []int, l, r, x int) int64 {
	var res int64
	for i := l; i <= r; i++ {
		or := 0
		for j := i; j <= r; j++ {
			or |= arr[j-1]
			if or >= x {
				res++
			}
		}
	}
	return res
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierF.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))

	type query struct{ typ, x, y int }
	type test struct {
		n, m, x int
		arr     []int
		qs      []query
	}
	var cases []test
	cases = append(cases, test{1, 1, 1, []int{0}, []query{{2, 1, 1}}})
	for i := 0; i < 99; i++ {
		n := rng.Intn(10) + 1
		m := rng.Intn(10) + 1
		x := rng.Intn(16)
		arr := make([]int, n)
		for j := 0; j < n; j++ {
			arr[j] = rng.Intn(16)
		}
		qs := make([]query, m)
		for j := 0; j < m; j++ {
			typ := rng.Intn(2) + 1
			if typ == 1 {
				pos := rng.Intn(n) + 1
				val := rng.Intn(16)
				qs[j] = query{1, pos, val}
			} else {
				l := rng.Intn(n) + 1
				r := rng.Intn(n-l+1) + l
				qs[j] = query{2, l, r}
			}
		}
		cases = append(cases, test{n, m, x, arr, qs})
	}

	for idx, tc := range cases {
		arr := make([]int, len(tc.arr))
		copy(arr, tc.arr)
		input := fmt.Sprintf("%d %d %d\n", tc.n, tc.m, tc.x)
		for i, v := range tc.arr {
			if i > 0 {
				input += " "
			}
			input += fmt.Sprintf("%d", v)
		}
		input += "\n"
		for _, q := range tc.qs {
			input += fmt.Sprintf("%d %d %d\n", q.typ, q.x, q.y)
		}
		got, err := runCandidate(bin, input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d runtime error: %v\n", idx+1, err)
			os.Exit(1)
		}
		outputs := strings.Fields(got)
		outIdx := 0
		for _, q := range tc.qs {
			if q.typ == 1 {
				arr[q.x-1] = q.y
			} else {
				if outIdx >= len(outputs) {
					fmt.Fprintf(os.Stderr, "case %d missing output\n", idx+1)
					os.Exit(1)
				}
				want := fmt.Sprintf("%d", countSub(arr, q.x, q.y, tc.x))
				if outputs[outIdx] != want {
					fmt.Fprintf(os.Stderr, "case %d failed: expected %s got %s\n", idx+1, want, outputs[outIdx])
					os.Exit(1)
				}
				outIdx++
			}
		}
		if outIdx != len(outputs) {
			fmt.Fprintf(os.Stderr, "case %d extra output\n", idx+1)
			os.Exit(1)
		}
	}
	fmt.Printf("All %d tests passed\n", len(cases))
}
