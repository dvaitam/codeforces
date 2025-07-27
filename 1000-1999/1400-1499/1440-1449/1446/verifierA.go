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
	"time"
)

type testCase struct {
	n int
	W int64
	w []int64
}

func expected(n int, W int64, w []int64) []int {
	half := (W + 1) / 2
	for i := 0; i < n; i++ {
		if w[i] >= half && w[i] <= W {
			return []int{i + 1}
		}
	}
	type pair struct {
		w   int64
		idx int
	}
	arr := make([]pair, n)
	for i := range arr {
		arr[i] = pair{w[i], i + 1}
	}
	sort.Slice(arr, func(i, j int) bool { return arr[i].w < arr[j].w })
	var sum int64
	var res []int
	for _, p := range arr {
		if sum+p.w <= W {
			sum += p.w
			res = append(res, p.idx)
			if sum >= half {
				return res
			}
		}
	}
	return nil
}

func runBin(bin, input string) (string, error) {
	cmd := exec.Command(bin)
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	var errb bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &errb
	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("runtime error: %v\n%s", err, errb.String())
	}
	return strings.TrimSpace(out.String()), nil
}

func generateTests() []testCase {
	tests := []testCase{
		{n: 1, W: 10, w: []int64{5}},
		{n: 3, W: 10, w: []int64{1, 2, 8}},
		{n: 4, W: 7, w: []int64{1, 2, 3, 4}},
		{n: 5, W: 100, w: []int64{60, 70, 80, 90, 95}},
		{n: 2, W: 1, w: []int64{2, 3}},
	}
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < 100; i++ {
		n := r.Intn(10) + 1
		W := r.Int63n(1000) + 1
		w := make([]int64, n)
		for j := 0; j < n; j++ {
			w[j] = r.Int63n(1000) + 1
		}
		tests = append(tests, testCase{n: n, W: W, w: w})
	}
	return tests
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("usage: go run verifierA.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	tests := generateTests()
	for i, t := range tests {
		input := fmt.Sprintf("1\n%d %d\n", t.n, t.W)
		for j, v := range t.w {
			if j > 0 {
				input += " "
			}
			input += fmt.Sprint(v)
		}
		input += "\n"
		exp := expected(t.n, t.W, t.w)
		out, err := runBin(bin, input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "test %d: %v\n", i+1, err)
			os.Exit(1)
		}
		fields := strings.Fields(out)
		if len(fields) == 0 {
			fmt.Fprintf(os.Stderr, "test %d: no output\n", i+1)
			os.Exit(1)
		}
		k, err := strconv.Atoi(fields[0])
		if err != nil {
			fmt.Fprintf(os.Stderr, "test %d: cannot parse k\n", i+1)
			os.Exit(1)
		}
		half := (t.W + 1) / 2
		if k == -1 {
			if exp != nil {
				fmt.Fprintf(os.Stderr, "test %d: expected solution but got -1\n", i+1)
				os.Exit(1)
			}
			continue
		}
		if len(fields) != k+1 {
			fmt.Fprintf(os.Stderr, "test %d: expected %d indices, got %d\n", i+1, k, len(fields)-1)
			os.Exit(1)
		}
		used := make(map[int]bool)
		var sum int64
		for j := 0; j < k; j++ {
			idx, err := strconv.Atoi(fields[j+1])
			if err != nil || idx < 1 || idx > t.n || used[idx] {
				fmt.Fprintf(os.Stderr, "test %d: invalid index %q\n", i+1, fields[j+1])
				os.Exit(1)
			}
			used[idx] = true
			sum += t.w[idx-1]
		}
		if sum < half || sum > t.W {
			fmt.Fprintf(os.Stderr, "test %d: invalid subset sum %d\n", i+1, sum)
			os.Exit(1)
		}
	}
	fmt.Printf("All %d tests passed\n", len(tests))
}
