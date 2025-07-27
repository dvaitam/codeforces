package main

import (
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"strings"
)

type testCase struct {
	n, k int
	arr  []int
	caps []int
}

func (tc testCase) Input() string {
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("%d %d\n", tc.n, tc.k))
	for i, v := range tc.arr {
		if i > 0 {
			sb.WriteByte(' ')
		}
		sb.WriteString(fmt.Sprint(v))
	}
	sb.WriteByte('\n')
	for i, v := range tc.caps {
		if i > 0 {
			sb.WriteByte(' ')
		}
		sb.WriteString(fmt.Sprint(v))
	}
	sb.WriteByte('\n')
	return sb.String()
}

func runBinary(path, input string) (string, error) {
	var cmd *exec.Cmd
	if strings.HasSuffix(path, ".go") {
		cmd = exec.Command("go", "run", path)
	} else {
		cmd = exec.Command(path)
	}
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	err := cmd.Run()
	return strings.TrimSpace(out.String()), err
}

func solve(tc testCase) string {
	k := tc.k
	a := make([]int, k+1)
	for _, v := range tc.arr {
		if v >= 1 && v <= k {
			a[v]++
		}
	}
	b := append([]int{0}, tc.caps...)
	suf := 0
	num := 0
	for i := k; i >= 1; i-- {
		suf += a[i]
		t := suf / b[i]
		if suf%b[i] != 0 {
			t++
		}
		if t > num {
			num = t
		}
	}
	if num == 0 {
		return "0"
	}
	ans := make([][]int, num)
	id := 0
	for i := 1; i <= k; i++ {
		for j := 0; j < a[i]; j++ {
			ans[id%num] = append(ans[id%num], i)
			id++
		}
	}
	var sb strings.Builder
	sb.WriteString(fmt.Sprintln(num))
	for _, g := range ans {
		sb.WriteString(fmt.Sprint(len(g)))
		for _, v := range g {
			sb.WriteString(" ")
			sb.WriteString(fmt.Sprint(v))
		}
		sb.WriteByte('\n')
	}
	return strings.TrimSpace(sb.String())
}

func genCaps(rng *rand.Rand, k, n int) []int {
	caps := make([]int, k)
	max := rng.Intn(n) + 1
	caps[0] = max
	for i := 1; i < k; i++ {
		max = rng.Intn(max) + 1
		caps[i] = max
	}
	return caps
}

func genTests() []testCase {
	rng := rand.New(rand.NewSource(4))
	tests := make([]testCase, 100)
	for i := range tests {
		k := rng.Intn(5) + 1
		n := rng.Intn(20) + 1
		arr := make([]int, n)
		for j := range arr {
			arr[j] = rng.Intn(k) + 1
		}
		caps := genCaps(rng, k, n)
		tests[i] = testCase{n: n, k: k, arr: arr, caps: caps}
	}
	return tests
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierD.go /path/to/binary")
		return
	}
	bin := os.Args[1]
	tests := genTests()
	for i, tc := range tests {
		input := tc.Input()
		want := solve(tc)
		out, err := runBinary(bin, input)
		if err != nil {
			fmt.Printf("case %d: runtime error: %v\ninput:\n%s", i+1, err, input)
			return
		}
		if strings.TrimSpace(out) != want {
			fmt.Printf("case %d failed\ninput:\n%sexpected:\n%s\ngot:\n%s\n", i+1, input, want, out)
			return
		}
	}
	fmt.Println("All tests passed")
}
