package main

import (
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"strings"
)

type testJ struct {
	n  int
	a  []int
	b  []int
	c  []int
	q  int
	qs [][2]int
}

func buildOracle() (string, error) {
	exe := "oracleJ.bin"
	cmd := exec.Command("go", "build", "-o", exe, "1571J.go")
	if out, err := cmd.CombinedOutput(); err != nil {
		return "", fmt.Errorf("build oracle failed: %v\n%s", err, out)
	}
	return "./" + exe, nil
}

func runProg(bin, input string) (string, error) {
	var cmd *exec.Cmd
	if strings.HasSuffix(bin, ".go") {
		cmd = exec.Command("go", "run", bin)
	} else {
		cmd = exec.Command(bin)
	}
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	var errBuf bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &errBuf
	err := cmd.Run()
	if err != nil {
		return out.String() + errBuf.String(), err
	}
	return strings.TrimSpace(out.String()), nil
}

func genTests() []testJ {
	rng := rand.New(rand.NewSource(1571))
	tests := make([]testJ, 0, 100)
	for len(tests) < 100 {
		n := rng.Intn(4) + 2
		a := make([]int, n-1)
		bArr := make([]int, n-1)
		c := make([]int, n)
		for i := 0; i < n-1; i++ {
			a[i] = rng.Intn(10) + 1
			bArr[i] = rng.Intn(10) + 1
		}
		for i := 0; i < n; i++ {
			c[i] = rng.Intn(10) + 1
		}
		q := rng.Intn(3) + 1
		qs := make([][2]int, q)
		for i := 0; i < q; i++ {
			l := rng.Intn(n-1) + 1
			r := l + rng.Intn(n-l) + 1
			qs[i] = [2]int{l, r}
		}
		tests = append(tests, testJ{n, a, bArr, c, q, qs})
	}
	return tests
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierJ.go /path/to/binary")
		return
	}
	bin := os.Args[1]
	oracle, err := buildOracle()
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	defer os.Remove(oracle)

	tests := genTests()
	for i, tc := range tests {
		var sb strings.Builder
		sb.WriteString(fmt.Sprintf("%d\n", tc.n))
		for idx, v := range tc.a {
			if idx > 0 {
				sb.WriteByte(' ')
			}
			sb.WriteString(fmt.Sprintf("%d", v))
		}
		sb.WriteByte('\n')
		for idx, v := range tc.b {
			if idx > 0 {
				sb.WriteByte(' ')
			}
			sb.WriteString(fmt.Sprintf("%d", v))
		}
		sb.WriteByte('\n')
		for idx, v := range tc.c {
			if idx > 0 {
				sb.WriteByte(' ')
			}
			sb.WriteString(fmt.Sprintf("%d", v))
		}
		sb.WriteByte('\n')
		sb.WriteString(fmt.Sprintf("%d\n", tc.q))
		for _, qr := range tc.qs {
			sb.WriteString(fmt.Sprintf("%d %d\n", qr[0], qr[1]))
		}
		input := sb.String()
		want, err := runProg(oracle, input)
		if err != nil {
			fmt.Printf("oracle runtime error on case %d: %v\n", i+1, err)
			os.Exit(1)
		}
		got, err := runProg(bin, input)
		if err != nil {
			fmt.Printf("candidate runtime error on case %d: %v\n%s", i+1, err, got)
			os.Exit(1)
		}
		if strings.TrimSpace(got) != strings.TrimSpace(want) {
			fmt.Printf("case %d failed\ninput:\n%sexpected: %s\ngot: %s\n", i+1, input, want, got)
			os.Exit(1)
		}
	}
	fmt.Printf("All %d tests passed\n", len(tests))
}
