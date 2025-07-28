package main

import (
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"strings"
	"time"
)

func solveF(n int, a, b [][]byte) string {
	row := make([]int, n)
	col := make([]int, n)
	for i := 0; i < n; i++ {
		if a[i][0] != b[i][0] {
			row[i] = 1
		}
	}
	for j := 0; j < n; j++ {
		if ((a[0][j] - '0') ^ byte(row[0])) != (b[0][j] - '0') {
			col[j] = 1
		}
	}
	for i := 0; i < n; i++ {
		for j := 0; j < n; j++ {
			v := int(a[i][j]-'0') ^ row[i] ^ col[j]
			if v != int(b[i][j]-'0') {
				return "NO"
			}
		}
	}
	return "YES"
}

func generateF(rng *rand.Rand) (string, string) {
	n := rng.Intn(5) + 1
	a := make([][]byte, n)
	b := make([][]byte, n)
	row := make([]int, n)
	col := make([]int, n)
	for i := 0; i < n; i++ {
		a[i] = make([]byte, n)
		b[i] = make([]byte, n)
		for j := 0; j < n; j++ {
			if rng.Intn(2) == 0 {
				a[i][j] = '0'
			} else {
				a[i][j] = '1'
			}
		}
	}
	for i := 0; i < n; i++ {
		row[i] = rng.Intn(2)
	}
	for j := 0; j < n; j++ {
		col[j] = rng.Intn(2)
	}
	for i := 0; i < n; i++ {
		for j := 0; j < n; j++ {
			v := (int(a[i][j]-'0') ^ row[i] ^ col[j])
			if v == 1 {
				b[i][j] = '1'
			} else {
				b[i][j] = '0'
			}
		}
	}
	// with probability 1/2 introduce error
	if rng.Intn(2) == 0 {
		i := rng.Intn(n)
		j := rng.Intn(n)
		if b[i][j] == '0' {
			b[i][j] = '1'
		} else {
			b[i][j] = '0'
		}
	}

	var sb strings.Builder
	sb.WriteString("1\n")
	sb.WriteString(fmt.Sprintf("%d\n", n))
	for i := 0; i < n; i++ {
		sb.WriteString(string(a[i]))
		sb.WriteByte('\n')
	}
	for i := 0; i < n; i++ {
		sb.WriteString(string(b[i]))
		sb.WriteByte('\n')
	}
	out := solveF(n, a, b)
	return sb.String(), out
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
	cmd.Stderr = &out
	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("runtime error: %v\n%s", err, out.String())
	}
	return strings.TrimSpace(out.String()), nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierF.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < 100; i++ {
		in, exp := generateF(rng)
		got, err := run(bin, in)
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:\n%s", i+1, err, in)
			os.Exit(1)
		}
		if got != exp {
			fmt.Fprintf(os.Stderr, "case %d failed: expected %s got %s\ninput:\n%s", i+1, exp, got, in)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
