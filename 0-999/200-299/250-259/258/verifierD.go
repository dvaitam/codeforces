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

func inversions(p []int) int {
	c := 0
	n := len(p)
	for i := 0; i < n; i++ {
		for j := i + 1; j < n; j++ {
			if p[i] > p[j] {
				c++
			}
		}
	}
	return c
}

func solve(n, m int, w []int, ops [][2]int) float64 {
	total := 0.0
	states := 1 << m
	for mask := 0; mask < states; mask++ {
		p := make([]int, n)
		copy(p, w)
		for i := 0; i < m; i++ {
			if mask&(1<<i) != 0 {
				a := ops[i][0] - 1
				b := ops[i][1] - 1
				p[a], p[b] = p[b], p[a]
			}
		}
		total += float64(inversions(p))
	}
	return total / float64(states)
}

func generateTests() []testCase {
	var tests []testCase
	n := 2
	m := 1
	for len(tests) < 100 {
		w := make([]int, n)
		for i := 0; i < n; i++ {
			w[i] = i + 1
		}
		ops := make([][2]int, m)
		for i := 0; i < m; i++ {
			ops[i] = [2]int{1, n}
		}
		ans := solve(n, m, w, ops)
		sb := strings.Builder{}
		sb.WriteString(fmt.Sprintf("%d %d\n", n, m))
		for i, v := range w {
			if i > 0 {
				sb.WriteByte(' ')
			}
			sb.WriteString(fmt.Sprint(v))
		}
		sb.WriteByte('\n')
		for _, op := range ops {
			sb.WriteString(fmt.Sprintf("%d %d\n", op[0], op[1]))
		}
		tests = append(tests, testCase{in: sb.String(), out: fmt.Sprintf("%.6f", ans)})
		n++
		if n > 3 {
			n = 2
			m++
			if m > 2 {
				m = 1
			}
		}
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
		fmt.Println("usage: go run verifierD.go /path/to/binary")
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
