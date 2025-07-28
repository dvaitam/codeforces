package main

import (
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"strings"
)

const mod int64 = 998244353
const maxn = 200000

var inv [maxn + 1]int64
var pref [maxn + 1]int64

func init() {
	inv[1] = 1
	for i := 2; i <= maxn; i++ {
		inv[i] = mod - mod/int64(i)*inv[mod%int64(i)]%mod
	}
	for i := 1; i <= maxn; i++ {
		pref[i] = (pref[i-1] + inv[i]*inv[i]) % mod
	}
}

type test struct {
	input    string
	expected string
}

func solveCase(a []int) string {
	n := len(a)
	ones := 0
	for _, v := range a {
		if v == 1 {
			ones++
		}
	}
	zeros := n - ones
	k := 0
	for i := 0; i < zeros; i++ {
		if a[i] == 1 {
			k++
		}
	}
	C := int64(n) * int64(n-1) / 2 % mod
	ans := C * pref[k] % mod
	return fmt.Sprintf("%d", ans)
}

func generateTests() []test {
	rng := rand.New(rand.NewSource(44))
	var tests []test
	tests = append(tests, test{input: "1\n1\n1\n", expected: solveCase([]int{1})})
	tests = append(tests, test{input: "2\n1 0\n", expected: solveCase([]int{1, 0})})
	for len(tests) < 100 {
		n := rng.Intn(6) + 1
		arr := make([]int, n)
		for i := 0; i < n; i++ {
			arr[i] = rng.Intn(2)
		}
		var sb strings.Builder
		sb.WriteString(fmt.Sprintf("%d\n", n))
		for i := 0; i < n; i++ {
			if i > 0 {
				sb.WriteByte(' ')
			}
			sb.WriteString(fmt.Sprintf("%d", arr[i]))
		}
		sb.WriteByte('\n')
		tests = append(tests, test{input: sb.String(), expected: solveCase(arr)})
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
	cmd.Stdin = strings.NewReader("1\n" + input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	err := cmd.Run()
	return strings.TrimSpace(out.String()), err
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierC.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	tests := generateTests()
	for i, t := range tests {
		out, err := runBinary(bin, t.input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Runtime error on test %d: %v\n", i+1, err)
			os.Exit(1)
		}
		if strings.TrimSpace(out) != strings.TrimSpace(t.expected) {
			fmt.Printf("Wrong answer on test %d\nInput:\n%sExpected:%s\nGot:%s\n", i+1, t.input, t.expected, out)
			os.Exit(1)
		}
	}
	fmt.Printf("All %d tests passed.\n", len(tests))
}
