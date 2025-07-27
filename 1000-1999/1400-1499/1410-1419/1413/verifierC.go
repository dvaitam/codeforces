package main

import (
	"bufio"
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"sort"
	"strconv"
	"strings"
)

type testCase struct {
	input  string
	expect string
}

func runBinary(path, input string) (string, error) {
	cmd := exec.Command(path)
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	err := cmd.Run()
	return strings.TrimSpace(out.String()), err
}

func solveC(input string) string {
	in := bufio.NewScanner(strings.NewReader(input))
	in.Split(bufio.ScanWords)
	nextInt64 := func() int64 {
		if !in.Scan() {
			return 0
		}
		v, _ := strconv.ParseInt(in.Text(), 10, 64)
		return v
	}
	var a [6]int64
	for i := 0; i < 6; i++ {
		a[i] = nextInt64()
	}
	var n int
	if in.Scan() {
		n, _ = strconv.Atoi(in.Text())
	}
	b := make([]int64, n)
	for i := 0; i < n; i++ {
		if !in.Scan() {
			break
		}
		val, _ := strconv.ParseInt(in.Text(), 10, 64)
		b[i] = val
	}
	type pair struct {
		v   int64
		idx int
	}
	arr := make([]pair, 0, n*6)
	for i := 0; i < n; i++ {
		for j := 0; j < 6; j++ {
			arr = append(arr, pair{b[i] - a[j], i})
		}
	}
	sort.Slice(arr, func(i, j int) bool { return arr[i].v < arr[j].v })
	cnt := make([]int, n)
	missing := n
	ans := int64(1<<63 - 1)
	l := 0
	for r := 0; r < len(arr); r++ {
		p := arr[r]
		if cnt[p.idx] == 0 {
			missing--
		}
		cnt[p.idx]++
		for missing == 0 && l <= r {
			diff := arr[r].v - arr[l].v
			if diff < ans {
				ans = diff
			}
			li := arr[l].idx
			cnt[li]--
			if cnt[li] == 0 {
				missing++
			}
			l++
		}
	}
	return strconv.FormatInt(ans, 10)
}

func genTests() []testCase {
	rng := rand.New(rand.NewSource(3))
	tests := make([]testCase, 100)
	for i := 0; i < 100; i++ {
		var a [6]int64
		for j := 0; j < 6; j++ {
			a[j] = int64(rng.Intn(20) + 1)
		}
		n := rng.Intn(4) + 2
		b := make([]int64, n)
		for j := 0; j < n; j++ {
			b[j] = int64(rng.Intn(30) + 1)
		}
		var sb strings.Builder
		for j := 0; j < 6; j++ {
			if j > 0 {
				sb.WriteByte(' ')
			}
			sb.WriteString(strconv.FormatInt(a[j], 10))
		}
		sb.WriteByte('\n')
		sb.WriteString(strconv.Itoa(n))
		sb.WriteByte('\n')
		for j := 0; j < n; j++ {
			if j > 0 {
				sb.WriteByte(' ')
			}
			sb.WriteString(strconv.FormatInt(b[j], 10))
		}
		sb.WriteByte('\n')
		input := sb.String()
		expect := solveC(input)
		tests[i] = testCase{input: input, expect: expect}
	}
	return tests
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierC.go /path/to/binary")
		return
	}
	bin := os.Args[1]
	tests := genTests()
	for i, tc := range tests {
		out, err := runBinary(bin, tc.input)
		if err != nil {
			fmt.Printf("test %d: runtime error: %v\n", i+1, err)
			return
		}
		if out != tc.expect {
			fmt.Printf("test %d failed\ninput:\n%s\nexpected:\n%s\nactual:\n%s\n", i+1, tc.input, tc.expect, out)
			return
		}
	}
	fmt.Println("All tests passed")
}
