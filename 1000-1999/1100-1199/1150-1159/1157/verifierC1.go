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

const numTestsC1 = 100

func main() {
	if len(os.Args) != 2 {
		fmt.Println("Usage: verifierC1.go <binary>")
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
	for t := 1; t <= numTestsC1; t++ {
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
		expected := solveC1(a)
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
		tmp := filepath.Join(os.TempDir(), "verify_binC1")
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

func solveC1(a []int) string {
	l, r := 0, len(a)-1
	now := 0
	res := make([]byte, 0, len(a))
	for l <= r {
		if a[l] <= now && a[r] <= now {
			break
		}
		if a[l] < a[r] {
			if now < a[l] {
				res = append(res, 'L')
				now = a[l]
				l++
			} else {
				res = append(res, 'R')
				now = a[r]
				r--
			}
		} else {
			if now < a[r] {
				res = append(res, 'R')
				now = a[r]
				r--
			} else {
				res = append(res, 'L')
				now = a[l]
				l++
			}
		}
	}
	return fmt.Sprintf("%d\n%s", len(res), string(res))
}
