package main

import (
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

const numTestsG = 100

func main() {
	if len(os.Args) != 2 {
		fmt.Println("Usage: verifierG.go <binary>")
		os.Exit(1)
	}
	binPath, cleanup, err := prepareBinary(os.Args[1])
	if err != nil {
		fmt.Println("compile error:", err)
		os.Exit(1)
	}
	if cleanup != nil {
		defer cleanup()
	}
	r := rand.New(rand.NewSource(1))
	for t := 1; t <= numTestsG; t++ {
		n := r.Intn(4) + 1
		m := r.Intn(4) + 1
		a := make([][]int, n)
		var sb strings.Builder
		sb.WriteString(fmt.Sprintf("%d %d\n", n, m))
		for i := 0; i < n; i++ {
			a[i] = make([]int, m)
			for j := 0; j < m; j++ {
				a[i][j] = r.Intn(2)
				if j > 0 {
					sb.WriteByte(' ')
				}
				sb.WriteString(fmt.Sprintf("%d", a[i][j]))
			}
			sb.WriteByte('\n')
		}
		input := sb.String()
		expected := solveG(a)
		out, err := run(binPath, input)
		if err != nil {
			fmt.Printf("test %d: runtime error: %v\n", t, err)
			os.Exit(1)
		}
		out = strings.TrimSpace(out)
		if out != expected {
			fmt.Printf("test %d failed\ninput:%sexpected:%s got:%s\n", t, input, expected, out)
			os.Exit(1)
		}
	}
	fmt.Println("OK")
}

func prepareBinary(path string) (string, func(), error) {
	if strings.HasSuffix(path, ".go") {
		tmp := filepath.Join(os.TempDir(), "verify_binG")
		cmd := exec.Command("go", "build", "-o", tmp, path)
		if out, err := cmd.CombinedOutput(); err != nil {
			return "", nil, fmt.Errorf("go build failed: %v: %s", err, string(out))
		}
		return tmp, func() { os.Remove(tmp) }, nil
	}
	return path, nil, nil
}

func run(bin, input string) (string, error) {
	cmd := exec.Command(bin)
	cmd.Stdin = strings.NewReader(input)
	var buf bytes.Buffer
	cmd.Stdout = &buf
	cmd.Stderr = os.Stderr
	err := cmd.Run()
	return buf.String(), err
}

func solveG(a [][]int) string {
	n := len(a)
	m := len(a[0])
	F := make([][2][2]bool, m-1)
	for i := 0; i < n; i++ {
		for j := 0; j+1 < m; j++ {
			x := a[i][j]
			y := a[i][j+1]
			u := 1 - x
			v := y
			F[j][u][v] = true
		}
	}
	dp := make([][2]bool, m)
	parent := make([][2]int, m)
	dp[0][0] = true
	dp[0][1] = true
	for j := 0; j+1 < m; j++ {
		for b := 0; b < 2; b++ {
			if !dp[j][b] {
				continue
			}
			for nb := 0; nb < 2; nb++ {
				if !F[j][b][nb] {
					if !dp[j+1][nb] {
						dp[j+1][nb] = true
						parent[j+1][nb] = b
					}
				}
			}
		}
	}
	c := make([]int, m)
	if m > 0 {
		if dp[m-1][0] {
			c[m-1] = 0
		} else if dp[m-1][1] {
			c[m-1] = 1
		} else {
			return "NO"
		}
		for j := m - 1; j > 0; j-- {
			c[j-1] = parent[j][c[j]]
		}
	}
	rarr := make([]int, n)
	if n > 0 {
		rarr[0] = 0
		for i := 0; i+1 < n; i++ {
			d := a[i][m-1] ^ c[m-1] ^ a[i+1][0] ^ c[0]
			if d == 0 {
				rarr[i+1] = rarr[i]
			} else {
				rarr[i+1] = 0
			}
		}
	}
	var sb strings.Builder
	sb.WriteString("YES\n")
	for i := 0; i < n; i++ {
		sb.WriteByte(byte('0' + rarr[i]))
	}
	sb.WriteByte('\n')
	for j := 0; j < m; j++ {
		sb.WriteByte(byte('0' + c[j]))
	}
	return strings.TrimSpace(sb.String())
}
