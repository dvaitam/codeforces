package main

import (
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"
)

const numTestsB = 100

func main() {
	if len(os.Args) != 2 {
		fmt.Println("Usage: verifierB.go <binary>")
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
	for t := 1; t <= numTestsB; t++ {
		n := r.Intn(20) + 1
		digits := make([]byte, n)
		for i := range digits {
			digits[i] = byte('1' + r.Intn(9))
		}
		f := make([]int, 10)
		for i := 1; i <= 9; i++ {
			f[i] = r.Intn(9) + 1
		}
		var sb strings.Builder
		sb.WriteString(fmt.Sprintf("%d\n", n))
		sb.WriteString(fmt.Sprintf("%s\n", string(digits)))
		for i := 1; i <= 9; i++ {
			if i > 1 {
				sb.WriteByte(' ')
			}
			sb.WriteString(strconv.Itoa(f[i]))
		}
		sb.WriteByte('\n')
		input := sb.String()
		expected := solveB(n, string(digits), f)
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
		tmp := filepath.Join(os.TempDir(), "verify_binB")
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

func solveB(n int, s string, f []int) string {
	b := []byte(s)
	started := false
	for i := 0; i < n; i++ {
		d := int(b[i] - '0')
		if !started {
			if f[d] > d {
				started = true
				b[i] = byte('0' + f[d])
			}
		} else {
			if f[d] >= d {
				b[i] = byte('0' + f[d])
			} else {
				break
			}
		}
	}
	return string(b)
}
