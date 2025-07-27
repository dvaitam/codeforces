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

type testCase struct {
	a []int
}

func expected(a []int) int {
	n := len(a)
	maxVal := 0
	for _, v := range a {
		if v > maxVal {
			maxVal = v
		}
	}
	counts := make([]int, maxVal+1)
	freqCount := make([]int, n+2)
	check := func(L int) bool {
		if L < 2 {
			return false
		}
		for i := range counts {
			counts[i] = 0
		}
		for i := range freqCount {
			freqCount[i] = 0
		}
		fmax := 0
		for i := 0; i < L; i++ {
			v := a[i]
			old := counts[v]
			if old > 0 {
				freqCount[old]--
			}
			counts[v] = old + 1
			freqCount[old+1]++
			if old+1 > fmax {
				fmax = old + 1
			}
		}
		if freqCount[fmax] >= 2 {
			return true
		}
		for i := L; i < n; i++ {
			v := a[i]
			old := counts[v]
			if old > 0 {
				freqCount[old]--
			}
			counts[v] = old + 1
			freqCount[old+1]++
			if old+1 > fmax {
				fmax = old + 1
			}
			u := a[i-L]
			old2 := counts[u]
			freqCount[old2]--
			counts[u] = old2 - 1
			if old2-1 > 0 {
				freqCount[old2-1]++
			}
			if old2 == fmax && freqCount[old2] == 0 {
				for fmax > 0 && freqCount[fmax] == 0 {
					fmax--
				}
			}
			if freqCount[fmax] >= 2 {
				return true
			}
		}
		return false
	}
	lo, hi := 0, n+1
	for lo+1 < hi {
		mid := (lo + hi) / 2
		if check(mid) {
			lo = mid
		} else {
			hi = mid
		}
	}
	return lo
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
		{a: []int{1, 1, 2, 2}},
		{a: []int{1, 2, 3, 4}},
		{a: []int{1, 1, 1, 2, 2, 2}},
		{a: []int{5, 5, 5, 5}},
	}
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < 100; i++ {
		n := r.Intn(20) + 1
		arr := make([]int, n)
		for j := 0; j < n; j++ {
			arr[j] = r.Intn(10) + 1
		}
		tests = append(tests, testCase{a: arr})
	}
	return tests
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("usage: go run verifierD1.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	tests := generateTests()
	for i, t := range tests {
		var sb strings.Builder
		sb.WriteString(fmt.Sprintf("%d\n", len(t.a)))
		for j, v := range t.a {
			if j > 0 {
				sb.WriteByte(' ')
			}
			sb.WriteString(fmt.Sprint(v))
		}
		sb.WriteByte('\n')
		want := expected(t.a)
		out, err := runBin(bin, sb.String())
		if err != nil {
			fmt.Fprintf(os.Stderr, "test %d: %v\n", i+1, err)
			os.Exit(1)
		}
		got, err := strconv.Atoi(strings.TrimSpace(out))
		if err != nil {
			fmt.Fprintf(os.Stderr, "test %d: cannot parse output\n", i+1)
			os.Exit(1)
		}
		if got != want {
			fmt.Fprintf(os.Stderr, "test %d failed: expected %d got %d\n", i+1, want, got)
			os.Exit(1)
		}
	}
	fmt.Printf("All %d tests passed\n", len(tests))
}
