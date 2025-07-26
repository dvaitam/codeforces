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

const numTestsC2 = 100

func main() {
	if len(os.Args) != 2 {
		fmt.Println("Usage: verifierC2.go <binary>")
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
	for t := 1; t <= numTestsC2; t++ {
		n := r.Intn(20) + 1
		a := make([]int, n)
		var sb strings.Builder
		sb.WriteString(fmt.Sprintf("%d\n", n))
		for i := 0; i < n; i++ {
			a[i] = r.Intn(100) + 1
			if i > 0 {
				sb.WriteByte(' ')
			}
			sb.WriteString(fmt.Sprintf("%d", a[i]))
		}
		sb.WriteByte('\n')
		input := sb.String()
		expected := solveC2(a)
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
		tmp := filepath.Join(os.TempDir(), "verify_binC2")
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

func solveC2(a []int) string {
	n := len(a)
	g := make([]int, n)
	g[0] = 1
	for i := 1; i < n; i++ {
		if a[i-1] > a[i] {
			g[i] = g[i-1] + 1
		} else {
			g[i] = 1
		}
	}
	f := make([]int, n)
	f[n-1] = 1
	for i := n - 2; i >= 0; i-- {
		if a[i+1] > a[i] {
			f[i] = f[i+1] + 1
		} else {
			f[i] = 1
		}
	}
	l, r := 0, n-1
	now, ans := 0, 0
	var sb strings.Builder
	for l <= r && (a[l] > now || a[r] > now) {
		j := -1
		if a[l] > now && a[r] > now {
			if a[l] == a[r] {
				if f[l] > g[r] {
					j = l
				} else {
					j = r
				}
			} else if a[l] < a[r] {
				j = l
			} else {
				j = r
			}
		} else if a[l] > now {
			j = l
		} else {
			j = r
		}
		now = a[j]
		ans++
		if j == l {
			sb.WriteByte('L')
			l++
		} else {
			sb.WriteByte('R')
			r--
		}
	}
	return fmt.Sprintf("%d\n%s", ans, sb.String())
}
