package main

import (
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"sort"
	"strings"
)

type testCase struct {
	arr []int64
}

func solve(arr []int64) string {
	m := make(map[int64]bool, len(arr))
	for _, v := range arr {
		m[v] = true
	}
	sort.Slice(arr, func(i, j int) bool { return arr[i] < arr[j] })
	uniq := arr[:0]
	for _, v := range arr {
		if len(uniq) == 0 || uniq[len(uniq)-1] != v {
			uniq = append(uniq, v)
		}
	}
	arr = uniq
	last := arr[len(arr)-1]
	found2 := false
	var ansX, ansY int64
	for _, x := range arr {
		var prevDiff int64 = -1
		for diff := int64(1); x+diff <= last; diff <<= 1 {
			if m[x+diff] {
				if prevDiff != -1 && diff == prevDiff*2 {
					return fmt.Sprintf("3\n%d %d %d\n", x, x+prevDiff, x+diff)
				}
				prevDiff = diff
			}
		}
		if prevDiff != -1 && !found2 {
			found2 = true
			ansX = x
			ansY = x + prevDiff
		}
	}
	if found2 {
		return fmt.Sprintf("2\n%d %d\n", ansX, ansY)
	}
	return fmt.Sprintf("1\n%d\n", arr[0])
}

func buildInput(tc testCase) string {
	var sb strings.Builder
	fmt.Fprintf(&sb, "%d\n", len(tc.arr))
	for i, v := range tc.arr {
		if i > 0 {
			sb.WriteByte(' ')
		}
		fmt.Fprintf(&sb, "%d", v)
	}
	sb.WriteByte('\n')
	return sb.String()
}

func runCase(bin string, tc testCase) error {
	input := buildInput(tc)
	cmd := exec.Command(bin)
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("runtime error: %v\n%s", err, out.String())
	}
	expected := strings.TrimSpace(solve(append([]int64(nil), tc.arr...)))
	got := strings.TrimSpace(out.String())
	if expected != got {
		return fmt.Errorf("expected:\n%s\n-- got:\n%s", expected, got)
	}
	return nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierD.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]

	var cases []testCase
	cases = append(cases, testCase{arr: []int64{7, 3, 5}})
	cases = append(cases, testCase{arr: []int64{1}})
	cases = append(cases, testCase{arr: []int64{1, 3}})

	rng := rand.New(rand.NewSource(3))
	for i := 0; i < 100; i++ {
		n := rng.Intn(8) + 1
		arr := make([]int64, n)
		for j := 0; j < n; j++ {
			arr[j] = int64(rng.Intn(41) - 20)
		}
		sort.Slice(arr, func(i, j int) bool { return arr[i] < arr[j] })
		arr = arr[:n]
		cases = append(cases, testCase{arr: arr})
	}

	for i, tc := range cases {
		if err := runCase(bin, tc); err != nil {
			fmt.Fprintf(os.Stderr, "test %d failed: %v\ninput: %s", i+1, err, buildInput(tc))
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
