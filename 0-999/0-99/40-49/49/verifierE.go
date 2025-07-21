package main

import (
	"bufio"
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"strings"
)

type testCase struct {
	input string
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("usage: go run verifierE.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	tests := generateTests()
	for i, t := range tests {
		expect := solveE(strings.NewReader(t.input))
		out, err := runBinary(bin, t.input)
		if err != nil {
			fmt.Printf("test %d: execution failed: %v\n", i+1, err)
			os.Exit(1)
		}
		if strings.TrimSpace(out) != strings.TrimSpace(expect) {
			fmt.Printf("test %d failed: expected %q got %q\n", i+1, strings.TrimSpace(expect), strings.TrimSpace(out))
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}

func runBinary(path, input string) (string, error) {
	cmd := exec.Command(path)
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = os.Stderr
	err := cmd.Run()
	return out.String(), err
}

type prod struct{ a, b, c int }

func solveE(r *strings.Reader) string {
	in := bufio.NewReader(r)
	var s1, s2 string
	var n int
	fmt.Fscan(in, &s1, &s2, &n)
	prods := make([]prod, 0, n)
	for i := 0; i < n; i++ {
		var rule string
		fmt.Fscan(in, &rule)
		if len(rule) == 5 && rule[1] == '-' && rule[2] == '>' {
			A := int(rule[0] - 'a')
			B := int(rule[3] - 'a')
			C := int(rule[4] - 'a')
			prods = append(prods, prod{A, B, C})
		}
	}
	n1 := len(s1)
	n2 := len(s2)
	dp1 := make([][]uint32, n1)
	for i := range dp1 {
		dp1[i] = make([]uint32, n1)
	}
	dp2 := make([][]uint32, n2)
	for i := range dp2 {
		dp2[i] = make([]uint32, n2)
	}
	for i := 0; i < n1; i++ {
		dp1[i][i] = 1 << (s1[i] - 'a')
	}
	for i := 0; i < n2; i++ {
		dp2[i][i] = 1 << (s2[i] - 'a')
	}
	for length := 2; length <= n1; length++ {
		for i := 0; i+length <= n1; i++ {
			j := i + length - 1
			var mask uint32
			for k := i; k < j; k++ {
				left := dp1[i][k]
				right := dp1[k+1][j]
				if left == 0 || right == 0 {
					continue
				}
				for _, p := range prods {
					if (left>>p.b)&1 != 0 && (right>>p.c)&1 != 0 {
						mask |= 1 << p.a
					}
				}
			}
			dp1[i][j] = mask
		}
	}
	for length := 2; length <= n2; length++ {
		for i := 0; i+length <= n2; i++ {
			j := i + length - 1
			var mask uint32
			for k := i; k < j; k++ {
				left := dp2[i][k]
				right := dp2[k+1][j]
				if left == 0 || right == 0 {
					continue
				}
				for _, p := range prods {
					if (left>>p.b)&1 != 0 && (right>>p.c)&1 != 0 {
						mask |= 1 << p.a
					}
				}
			}
			dp2[i][j] = mask
		}
	}
	const INF = 1e9
	dist := make([][]int, n1+1)
	for i := range dist {
		dist[i] = make([]int, n2+1)
		for j := range dist[i] {
			dist[i][j] = INF
		}
	}
	type pair struct{ i, j int }
	q := make([]pair, 0, (n1+1)*(n2+1))
	dist[0][0] = 0
	q = append(q, pair{0, 0})
	for head := 0; head < len(q); head++ {
		u := q[head]
		d := dist[u.i][u.j]
		for i2 := u.i + 1; i2 <= n1; i2++ {
			for j2 := u.j + 1; j2 <= n2; j2++ {
				m1 := dp1[u.i][i2-1]
				if m1 == 0 {
					continue
				}
				m2 := dp2[u.j][j2-1]
				if m2 == 0 {
					continue
				}
				if (m1 & m2) == 0 {
					continue
				}
				if dist[i2][j2] > d+1 {
					dist[i2][j2] = d + 1
					q = append(q, pair{i2, j2})
				}
			}
		}
	}
	ans := dist[n1][n2]
	if ans >= INF {
		return "-1\n"
	}
	return fmt.Sprintf("%d\n", ans)
}

func generateTests() []testCase {
	rand.Seed(5)
	tests := make([]testCase, 0, 100)
	letters := []byte{'a', 'b', 'c'}
	for i := 0; i < 20; i++ {
		tests = append(tests, testCase{input: "a a 0\n"})
	}
	for len(tests) < 100 {
		n1 := rand.Intn(3) + 1
		n2 := rand.Intn(3) + 1
		var s1, s2 strings.Builder
		for i := 0; i < n1; i++ {
			s1.WriteByte(letters[rand.Intn(len(letters))])
		}
		for i := 0; i < n2; i++ {
			s2.WriteByte(letters[rand.Intn(len(letters))])
		}
		rules := rand.Intn(3)
		var sb strings.Builder
		fmt.Fprintf(&sb, "%s %s %d\n", s1.String(), s2.String(), rules)
		for i := 0; i < rules; i++ {
			a := letters[rand.Intn(len(letters))]
			b := letters[rand.Intn(len(letters))]
			c := letters[rand.Intn(len(letters))]
			fmt.Fprintf(&sb, "%c->%c%c\n", a, b, c)
		}
		tests = append(tests, testCase{input: sb.String()})
	}
	return tests
}
