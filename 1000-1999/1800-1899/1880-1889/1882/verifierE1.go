package main

import (
	"bufio"
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"strings"
	"time"
)

func runCandidate(bin, input string) (string, error) {
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
	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("runtime error: %v\n%s", err, errBuf.String())
	}
	return strings.TrimSpace(out.String()), nil
}

func readInt(reader *bufio.Reader) (int, error) {
	x := 0
	sign := 1
	b, err := reader.ReadByte()
	for err == nil && (b == ' ' || b == '\n' || b == '\r' || b == '\t') {
		b, err = reader.ReadByte()
	}
	if err != nil {
		return 0, err
	}
	if b == '-' {
		sign = -1
		b, err = reader.ReadByte()
		if err != nil {
			return 0, err
		}
	}
	for b >= '0' && b <= '9' {
		x = x*10 + int(b-'0')
		b, err = reader.ReadByte()
		if err != nil {
			break
		}
	}
	return x * sign, nil
}

func solvePerm(n int, a []int) []int {
	var ops []int
	for i := 1; i <= n; i++ {
		for j := i; j <= n; j++ {
			if a[j] == i {
				if j != i {
					ops = append(ops, i, j-i, n-j+1)
				}
				a[i], a[j] = a[j], a[i]
				break
			}
		}
	}
	return ops
}

func expectedOutput(n, m int, a, b []int) string {
	opa := solvePerm(n, append([]int(nil), a...))
	opb := solvePerm(m, append([]int(nil), b...))
	la, lb := len(opa), len(opb)
	if (la+lb)%2 == 1 {
		if n%2 == 1 {
			for i := 0; i < n; i++ {
				opa = append(opa, n)
			}
		} else if m%2 == 1 {
			for i := 0; i < m; i++ {
				opb = append(opb, m)
			}
		} else {
			return "-1"
		}
		la, lb = len(opa), len(opb)
	}
	for la < lb {
		opa = append(opa, 1, n)
		la += 2
	}
	for lb < la {
		opb = append(opb, 1, m)
		lb += 2
	}
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("%d\n", la))
	for i := 0; i < la; i++ {
		sb.WriteString(fmt.Sprintf("%d %d\n", opa[i], opb[i]))
	}
	return strings.TrimSpace(sb.String())
}

func genCase(rng *rand.Rand) (string, string) {
	n := rng.Intn(5) + 1
	m := rng.Intn(5) + 1
	p := rng.Perm(n)
	q := rng.Perm(m)
	a := make([]int, n+1)
	b := make([]int, m+1)
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("%d %d\n", n, m))
	for i, v := range p {
		a[i+1] = v + 1
		if i > 0 {
			sb.WriteByte(' ')
		}
		sb.WriteString(fmt.Sprintf("%d", v+1))
	}
	sb.WriteByte('\n')
	for i, v := range q {
		b[i+1] = v + 1
		if i > 0 {
			sb.WriteByte(' ')
		}
		sb.WriteString(fmt.Sprintf("%d", v+1))
	}
	sb.WriteByte('\n')
	expect := expectedOutput(n, m, a, b)
	input := sb.String()
	input = fmt.Sprintf("%s", input)
	return input, expect
}

func fixedCases() [][2]string {
	return [][2]string{
		{"1 1\n1\n1\n", "0"},
	}
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierE1.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	for idx, tc := range fixedCases() {
		out, err := runCandidate(bin, tc[0])
		if err != nil {
			fmt.Fprintf(os.Stderr, "fixed case %d failed: %v\ninput:\n%s", idx+1, err, tc[0])
			os.Exit(1)
		}
		if strings.TrimSpace(out) != tc[1] {
			fmt.Fprintf(os.Stderr, "fixed case %d failed: expected %s got %s\ninput:\n%s", idx+1, tc[1], out, tc[0])
			os.Exit(1)
		}
	}
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < 100; i++ {
		in, exp := genCase(rng)
		out, err := runCandidate(bin, in)
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:\n%s", i+1, err, in)
			os.Exit(1)
		}
		if strings.TrimSpace(out) != exp {
			fmt.Fprintf(os.Stderr, "case %d failed: expected %s got %s\ninput:\n%s", i+1, exp, out, in)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
