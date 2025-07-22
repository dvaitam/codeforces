package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"strings"
)

type testCase struct {
	in  string
	out string
}

func gcd(a, b int) int {
	for b != 0 {
		a, b = b, a%b
	}
	return a
}

func lcm(a, b int) int {
	return a / gcd(a, b) * b
}

func solve(a []int) int64 {
	n := len(a)
	b := make([]int, n)
	var ans int64
	var dfs func(pos int)
	dfs = func(pos int) {
		if pos == n {
			m := 0
			l := b[0]
			for i := 0; i < n; i++ {
				if b[i] > m {
					m = b[i]
				}
				if i > 0 {
					l = lcm(l, b[i])
				}
			}
			if l == m {
				ans++
			}
			return
		}
		for v := 1; v <= a[pos]; v++ {
			b[pos] = v
			dfs(pos + 1)
		}
	}
	dfs(0)
	return ans % 1000000007
}

func generateTests() []testCase {
	var tests []testCase
	for n := 1; len(tests) < 100; n++ {
		if n > 3 {
			n = 1
		}
		arr := make([]int, n)
		for i := 0; i < n; i++ {
			arr[i] = (i+1)%4 + 1
		}
		out := solve(arr)
		sb := strings.Builder{}
		sb.WriteString(fmt.Sprintf("%d\n", n))
		for i, v := range arr {
			if i > 0 {
				sb.WriteByte(' ')
			}
			sb.WriteString(fmt.Sprint(v))
		}
		sb.WriteByte('\n')
		tests = append(tests, testCase{in: sb.String(), out: fmt.Sprint(out)})
	}
	return tests
}

func runTest(bin string, tc testCase) (string, error) {
	cmd := exec.Command(bin)
	stdin, err := cmd.StdinPipe()
	if err != nil {
		return "", err
	}
	go func() {
		defer stdin.Close()
		stdin.Write([]byte(tc.in))
	}()
	out, err := cmd.CombinedOutput()
	return strings.TrimSpace(string(out)), err
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("usage: go run verifierC.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	if strings.HasSuffix(bin, ".go") {
		tmp, err := ioutil.TempFile("", "solbin*")
		if err != nil {
			fmt.Println("cannot create temp file:", err)
			os.Exit(1)
		}
		tmp.Close()
		exec.Command("go", "build", "-o", tmp.Name(), bin).Run()
		bin = tmp.Name()
		defer os.Remove(bin)
	}
	tests := generateTests()
	for i, tc := range tests {
		got, err := runTest(bin, tc)
		if err != nil {
			fmt.Printf("runtime error on test %d: %v\n", i+1, err)
			os.Exit(1)
		}
		if got != tc.out {
			fmt.Printf("wrong answer on test %d\ninput: %sexpected: %s\ngot: %s\n", i+1, tc.in, tc.out, got)
			os.Exit(1)
		}
	}
	fmt.Println("all tests passed")
}
