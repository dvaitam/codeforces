package main

import (
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"strings"
)

func runBinary(path, input string) (string, error) {
	cmd := exec.Command(path)
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	err := cmd.Run()
	return strings.TrimSpace(out.String()), err
}

func zigzagSeq(z, length int) []int64 {
	t := 2 * (z - 1)
	res := make([]int64, length)
	for i := 0; i < length; i++ {
		m := (i + 1) % t
		if m == 0 {
			res[i] = 2
		} else if m <= z {
			res[i] = int64(m)
		} else {
			res[i] = int64(2*z - m)
		}
	}
	return res
}

func solveD(input string) string {
	rd := strings.Fields(input)
	idx := 0
	n := atoi(rd[idx])
	idx++
	a := make([]int64, n)
	for i := 0; i < n; i++ {
		a[i] = int64(atoi(rd[idx]))
		idx++
	}
	m := atoi(rd[idx])
	idx++
	var out strings.Builder
	for qi := 0; qi < m; qi++ {
		ty := atoi(rd[idx])
		idx++
		if ty == 1 {
			p := atoi(rd[idx])
			idx++
			v := int64(atoi(rd[idx]))
			idx++
			a[p-1] = v
		} else {
			l := atoi(rd[idx])
			idx++
			r := atoi(rd[idx])
			idx++
			z := atoi(rd[idx])
			idx++
			seq := zigzagSeq(z, r-l+1)
			var sum int64
			for i := l - 1; i < r; i++ {
				sum += a[i] * seq[i-(l-1)]
			}
			out.WriteString(fmt.Sprintln(sum))
		}
	}
	return strings.TrimSpace(out.String())
}

func atoi(s string) int {
	var n int
	for i := 0; i < len(s); i++ {
		n = n*10 + int(s[i]-'0')
	}
	return n
}

func genTests() []string {
	r := rand.New(rand.NewSource(42))
	tests := []string{}
	for len(tests) < 100 {
		n := r.Intn(3) + 3 // 3..5
		var sb strings.Builder
		fmt.Fprintf(&sb, "%d\n", n)
		for i := 0; i < n; i++ {
			fmt.Fprintf(&sb, "%d ", r.Intn(5)+1)
		}
		sb.WriteByte('\n')
		m := r.Intn(3) + 3 // 3..5 ops
		fmt.Fprintf(&sb, "%d\n", m)
		for i := 0; i < m; i++ {
			if r.Intn(2) == 0 {
				p := r.Intn(n) + 1
				v := r.Intn(5) + 1
				fmt.Fprintf(&sb, "1 %d %d\n", p, v)
			} else {
				l := r.Intn(n) + 1
				rgt := r.Intn(n-l+1) + l
				z := r.Intn(5) + 2
				fmt.Fprintf(&sb, "2 %d %d %d\n", l, rgt, z)
			}
		}
		tests = append(tests, sb.String())
	}
	return tests
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierD.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	tests := genTests()
	for i, in := range tests {
		exp := solveD(strings.TrimSpace(in))
		out, err := runBinary(bin, in)
		if err != nil {
			fmt.Printf("test %d runtime error: %v\n", i+1, err)
			os.Exit(1)
		}
		if strings.TrimSpace(out) != exp {
			fmt.Printf("test %d failed:\ninput:\n%sexpected=%s got=%s\n", i+1, in, exp, strings.TrimSpace(out))
			os.Exit(1)
		}
	}
	fmt.Printf("ok %d tests\n", len(tests))
}
