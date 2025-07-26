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

const numTestsA = 100

func main() {
	if len(os.Args) != 2 {
		fmt.Println("Usage: verifierA.go <binary>")
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
	for t := 1; t <= numTestsA; t++ {
		n := r.Intn(1_000_000_000) + 1
		input := fmt.Sprintf("%d\n", n)
		expected := solveA(n)
		out, err := run(binPath, input)
		if err != nil {
			fmt.Printf("test %d: runtime error: %v\n", t, err)
			os.Exit(1)
		}
		out = strings.TrimSpace(out)
		if out != expected {
			fmt.Printf("test %d failed: n=%d expected %s got %s\n", t, n, expected, out)
			os.Exit(1)
		}
	}
	fmt.Println("OK")
}

func prepareBinary(path string) (string, func(), error) {
	if strings.HasSuffix(path, ".go") {
		tmp := filepath.Join(os.TempDir(), "verify_binA")
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

func solveA(n int) string {
	seen := make(map[int]struct{})
	count := 0
	for {
		if _, ok := seen[n]; ok {
			break
		}
		seen[n] = struct{}{}
		count++
		n = f(n)
	}
	return strconv.Itoa(count)
}

func f(x int) int {
	x++
	for x%10 == 0 {
		x /= 10
	}
	return x
}
