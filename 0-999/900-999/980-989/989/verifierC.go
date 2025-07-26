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
	in  string
	out string
}

func expected(a, b, c, d int) string {
	m := make([][]byte, 50)
	for i := 0; i < 50; i++ {
		m[i] = make([]byte, 50)
		for j := 0; j < 50; j++ {
			m[i][j] = '#'
		}
	}
	f := func(x, y int, ch byte) {
		w := 0
		for i := 1; i <= x; i++ {
			row := y - 1 + w%3
			col := i
			m[row][col] = ch
			w++
		}
	}
	ff := func(k, x, y int, ch byte) {
		w := 0
		for i := 1; i <= x; i++ {
			row := y - 1 + w%3
			col := i + k - 1
			m[row][col] = ch
			w++
		}
	}
	for i := 1; i <= 5; i++ {
		m[i-1][0] = 'A'
	}
	for i := 1; i <= 50; i++ {
		m[0][i-1] = 'A'
	}
	for i := 1; i <= 40; i++ {
		m[i-1][49] = 'A'
	}
	for i := 2; i <= 50; i++ {
		m[30][i-1] = 'A'
	}
	for i := 31; i <= 45; i++ {
		m[i-1][0] = 'D'
	}
	if a > 90 {
		ff(3, a-90, 27, 'A')
		a = 90
	}
	if b > 90 {
		ff(14, b-90, 27, 'B')
		b = 90
	}
	if c > 90 {
		ff(25, c-90, 27, 'C')
		c = 90
	}
	if a > 45 {
		f(45, 7, 'A')
		a -= 45
	}
	f(a, 3, 'A')
	if b > 45 {
		f(45, 11, 'B')
		b -= 45
	}
	f(b, 15, 'B')
	if c > 45 {
		f(45, 19, 'C')
		c -= 45
	}
	f(c, 23, 'C')
	if d > 90 {
		ff(3, d-90, 32, 'D')
		d = 90
	}
	if d > 45 {
		ff(3, 45, 36, 'D')
		d -= 45
	}
	f(d, 45, 'D')
	for i := 1; i <= 31; i++ {
		for j := 1; j <= 50; j++ {
			if m[i-1][j-1] == '#' {
				m[i-1][j-1] = 'D'
			}
		}
	}
	for i := 31; i <= 50; i++ {
		for j := 1; j <= 50; j++ {
			if m[i-1][j-1] == '#' {
				m[i-1][j-1] = 'A'
			}
		}
	}
	var sb strings.Builder
	fmt.Fprintf(&sb, "50 50\n")
	for i := 0; i < 50; i++ {
		sb.Write(m[i])
		sb.WriteByte('\n')
	}
	return strings.TrimSpace(sb.String())
}

func generate() []testCase {
	const T = 100
	rand.Seed(3)
	cases := make([]testCase, T)
	for i := 0; i < T; i++ {
		a := rand.Intn(100) + 1
		b := rand.Intn(100) + 1
		c := rand.Intn(100) + 1
		d := rand.Intn(100) + 1
		cases[i] = testCase{
			in:  fmt.Sprintf("%d %d %d %d\n", a, b, c, d),
			out: expected(a, b, c, d),
		}
	}
	return cases
}

func run(bin, input string) (string, error) {
	var cmd *exec.Cmd
	if strings.HasSuffix(bin, ".go") {
		cmd = exec.Command("go", "run", bin)
	} else {
		cmd = exec.Command(bin)
	}
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = os.Stderr
	err := cmd.Run()
	return strings.TrimSpace(out.String()), err
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierC.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	cases := generate()
	for idx, tc := range cases {
		got, err := run(bin, tc.in)
		if err != nil {
			fmt.Printf("case %d: runtime error: %v\n", idx+1, err)
			os.Exit(1)
		}
		if strings.TrimSpace(got) != tc.out {
			fmt.Printf("case %d failed\ninput:\n%s\nexpected:\n%s\n\ngot:\n%s\n", idx+1, tc.in, tc.out, got)
			os.Exit(1)
		}
	}
	fmt.Printf("All %d tests passed\n", len(cases))
}
