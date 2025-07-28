package main

import (
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"strings"
)

const mod int64 = 1_000_000_007

type test struct {
	input    string
	expected string
}

func modPow(base, exp, m int64) int64 {
	res := int64(1)
	for exp > 0 {
		if exp&1 == 1 {
			res = res * base % m
		}
		base = base * base % m
		exp >>= 1
	}
	return res
}

func maxSub(arr []int64) int64 {
	best := arr[0]
	cur := arr[0]
	for i := 1; i < len(arr); i++ {
		if cur+arr[i] > arr[i] {
			cur += arr[i]
		} else {
			cur = arr[i]
		}
		if cur > best {
			best = cur
		}
	}
	return best
}

func solve(input string) string {
	reader := strings.NewReader(input)
	var t int
	fmt.Fscan(reader, &t)
	var out strings.Builder
	for ; t > 0; t-- {
		var n, k int
		fmt.Fscan(reader, &n, &k)
		arr := make([]int64, n)
		var sum int64
		for i := 0; i < n; i++ {
			fmt.Fscan(reader, &arr[i])
			sum += arr[i]
		}
		mx := maxSub(arr)
		if mx <= 0 {
			out.WriteString(fmt.Sprintf("%d\n", (sum%mod+mod)%mod))
			continue
		}
		pow2 := modPow(2, int64(k), mod)
		if sum >= mx {
			ans := ((sum%mod + mod) % mod) * pow2 % mod
			out.WriteString(fmt.Sprintf("%d\n", ans))
		} else {
			inc := (mx%mod + mod) % mod
			inc = inc * ((pow2 - 1 + mod) % mod) % mod
			ans := ((sum%mod+mod)%mod + inc) % mod
			out.WriteString(fmt.Sprintf("%d\n", ans))
		}
	}
	return strings.TrimSpace(out.String())
}

func generateTests() []test {
	rng := rand.New(rand.NewSource(47))
	var tests []test
	fixed := []string{
		"1\n1 0\n5\n",
		"1\n3 1\n1 2 3\n",
	}
	for _, f := range fixed {
		tests = append(tests, test{f, solve(f)})
	}
	for len(tests) < 100 {
		n := rng.Intn(5) + 1
		k := rng.Intn(4)
		var sb strings.Builder
		sb.WriteString("1\n")
		sb.WriteString(fmt.Sprintf("%d %d\n", n, k))
		for i := 0; i < n; i++ {
			if i > 0 {
				sb.WriteByte(' ')
			}
			fmt.Fprintf(&sb, "%d", rng.Intn(11)-5)
		}
		sb.WriteByte('\n')
		inp := sb.String()
		tests = append(tests, test{inp, solve(inp)})
	}
	return tests
}

func runBinary(bin, input string) (string, error) {
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

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierB.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	tests := generateTests()
	for i, t := range tests {
		got, err := runBinary(bin, t.input)
		if err != nil {
			fmt.Printf("Runtime error on test %d: %v\n", i+1, err)
			os.Exit(1)
		}
		if strings.TrimSpace(got) != strings.TrimSpace(t.expected) {
			fmt.Printf("Wrong answer on test %d\nInput:\n%sExpected:%s\nGot:%s\n", i+1, t.input, t.expected, got)
			os.Exit(1)
		}
	}
	fmt.Printf("All %d tests passed.\n", len(tests))
}
