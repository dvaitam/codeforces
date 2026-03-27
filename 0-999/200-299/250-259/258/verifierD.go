package main

import (
	"bytes"
	"fmt"
	"io"
	"math"
	"math/rand"
	"os"
	"os/exec"
	"strconv"
	"strings"
)

func solve(input string) string {
	data := []byte(input)
	pos := 0
	nextInt := func() int {
		for pos < len(data) && (data[pos] < '0' || data[pos] > '9') {
			pos++
		}
		val := 0
		for pos < len(data) && data[pos] >= '0' && data[pos] <= '9' {
			val = val*10 + int(data[pos]-'0')
			pos++
		}
		return val
	}

	n := nextInt()
	m := nextInt()

	p := make([]int, n)
	for i := 0; i < n; i++ {
		p[i] = nextInt()
	}

	f := make([][]float64, n)
	for i := 0; i < n; i++ {
		f[i] = make([]float64, n)
		for j := 0; j < n; j++ {
			if p[i] > p[j] {
				f[i][j] = 1.0
			}
		}
	}

	for t := 0; t < m; t++ {
		a := nextInt() - 1
		b := nextInt() - 1

		for k := 0; k < n; k++ {
			if k == a || k == b {
				continue
			}
			v1 := (f[a][k] + f[b][k]) * 0.5
			v2 := (f[k][a] + f[k][b]) * 0.5
			f[a][k] = v1
			f[b][k] = v1
			f[k][a] = v2
			f[k][b] = v2
		}
		f[a][b] = 0.5
		f[b][a] = 0.5
	}

	ans := 0.0
	for i := 0; i < n; i++ {
		for j := i + 1; j < n; j++ {
			ans += f[i][j]
		}
	}

	return fmt.Sprintf("%.10f", ans)
}

func generateTests() []string {
	rng := rand.New(rand.NewSource(42))
	var inputs []string
	for len(inputs) < 100 {
		n := rng.Intn(6) + 2
		m := rng.Intn(5) + 1
		var sb strings.Builder
		sb.WriteString(fmt.Sprintf("%d %d\n", n, m))
		perm := rng.Perm(n)
		for i, v := range perm {
			if i > 0 {
				sb.WriteByte(' ')
			}
			sb.WriteString(strconv.Itoa(v + 1))
		}
		sb.WriteByte('\n')
		for i := 0; i < m; i++ {
			a := rng.Intn(n) + 1
			b := a
			for b == a {
				b = rng.Intn(n) + 1
			}
			if a > b {
				a, b = b, a
			}
			sb.WriteString(fmt.Sprintf("%d %d\n", a, b))
		}
		inputs = append(inputs, sb.String())
	}
	return inputs
}

func runCandidate(bin, input string) (string, error) {
	cmd := exec.Command(bin)
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	var stderr bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &stderr
	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("runtime error: %v\nstderr: %s", err, stderr.String())
	}
	return strings.TrimSpace(out.String()), nil
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("usage: go run verifierD.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	if strings.HasSuffix(bin, ".go") {
		tmp, err := os.CreateTemp("", "solbin*")
		if err != nil {
			fmt.Println("cannot create temp file:", err)
			os.Exit(1)
		}
		tmp.Close()
		out, err := exec.Command("go", "build", "-o", tmp.Name(), bin).CombinedOutput()
		if err != nil {
			fmt.Printf("build error: %v\n%s\n", err, string(out))
			os.Remove(tmp.Name())
			os.Exit(1)
		}
		bin = tmp.Name()
		defer os.Remove(bin)
	}

	_ = io.Discard
	tests := generateTests()
	for i, input := range tests {
		expected := solve(input)
		got, err := runCandidate(bin, input)
		if err != nil {
			fmt.Printf("runtime error on test %d: %v\n", i+1, err)
			os.Exit(1)
		}
		expVal, err1 := strconv.ParseFloat(expected, 64)
		gotVal, err2 := strconv.ParseFloat(got, 64)
		if err1 != nil || err2 != nil {
			fmt.Printf("parse error on test %d: expected=%q got=%q\n", i+1, expected, got)
			os.Exit(1)
		}
		if math.Abs(expVal-gotVal) > 1e-6 {
			fmt.Printf("wrong answer on test %d\ninput: %sexpected: %s\ngot: %s\n", i+1, input, expected, got)
			os.Exit(1)
		}
	}
	fmt.Println("all tests passed")
}
