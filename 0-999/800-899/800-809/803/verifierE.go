package main

import (
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
)

func buildReference() (string, error) {
	_, file, _, _ := runtime.Caller(0)
	dir := filepath.Dir(file)
	src := filepath.Join(dir, "803E.go")
	bin := filepath.Join(os.TempDir(), "ref803E.bin")
	cmd := exec.Command("go", "build", "-o", bin, src)
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		return "", err
	}
	return bin, nil
}

func runBinary(bin, input string) (string, error) {
	cmd := exec.Command(bin)
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &bytes.Buffer{}
	err := cmd.Run()
	return strings.TrimSpace(out.String()), err
}

func randChar() byte {
	switch rand.Intn(4) {
	case 0:
		return 'W'
	case 1:
		return 'L'
	case 2:
		return 'D'
	default:
		return '?'
	}
}

func genCase() string {
	n := rand.Intn(20) + 1
	k := rand.Intn(5) + 1
	b := make([]byte, n)
	for i := range b {
		b[i] = randChar()
	}
	return fmt.Sprintf("%d %d\n%s\n", n, k, string(b))
}

func abs(x int) int {
	if x < 0 {
		return -x
	}
	return x
}

func valid(n, k int, orig, out string) bool {
	if len(out) != n {
		return false
	}
	diff := 0
	for i := 0; i < n; i++ {
		c := out[i]
		if c != 'W' && c != 'L' && c != 'D' {
			return false
		}
		if orig[i] != '?' && orig[i] != c {
			return false
		}
		switch c {
		case 'W':
			diff++
		case 'L':
			diff--
		}
		if i != n-1 && abs(diff) >= k {
			return false
		}
	}
	return abs(diff) == k
}

func parseInput(input string) (int, int, string) {
	var n, k int
	var s string
	fmt.Fscan(strings.NewReader(input), &n, &k, &s)
	return n, k, s
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("usage: verifierE.go <path-to-binary>")
		return
	}
	userBin := os.Args[1]
	rand.Seed(1)
	refBin, err := buildReference()
	if err != nil {
		fmt.Println("failed to build reference:", err)
		return
	}
	defer os.Remove(refBin)

	for i := 0; i < 100; i++ {
		input := genCase()
		want, err1 := runBinary(refBin, input)
		if err1 != nil {
			fmt.Println("reference solution failed:", err1)
			return
		}
		got, err2 := runBinary(userBin, input)
		if err2 != nil {
			fmt.Printf("test %d: runtime error: %v\n", i+1, err2)
			fmt.Println("input:\n" + input)
			return
		}

		n, k, s := parseInput(input)
		if want == "NO" {
			if got != "NO" {
				fmt.Printf("test %d failed\ninput:\n%sexpected NO but got:\n%s\n", i+1, input, got)
				return
			}
		} else {
			if got == "NO" || !valid(n, k, s, got) {
				fmt.Printf("test %d failed\ninput:\n%sreference output:\n%s\nyour output:\n%s\n", i+1, input, want, got)
				return
			}
		}
	}
	fmt.Println("All tests passed")
}
