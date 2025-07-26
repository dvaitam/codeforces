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

const numTestsB = 100

func prepareBinary(path string) (string, func(), error) {
	if strings.HasSuffix(path, ".go") {
		tmp := filepath.Join(os.TempDir(), "verifB_bin")
		cmd := exec.Command("go", "build", "-o", tmp, path)
		if out, err := cmd.CombinedOutput(); err != nil {
			return "", nil, fmt.Errorf("go build failed: %v: %s", err, out)
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
	cmd.Stderr = &buf
	err := cmd.Run()
	return strings.TrimSpace(buf.String()), err
}

func solveB(s string) string {
	n := len(s)
	pref := make([]int64, n)
	for i := 1; i < n; i++ {
		pref[i] = pref[i-1]
		if s[i] == 'v' && s[i-1] == 'v' {
			pref[i]++
		}
	}
	total := int64(0)
	if n > 0 {
		total = pref[n-1]
	}
	var ans int64
	for i := 0; i < n; i++ {
		if s[i] == 'o' {
			left := pref[i]
			right := total - pref[i]
			ans += left * right
		}
	}
	return fmt.Sprint(ans)
}

func genCaseB(rng *rand.Rand) string {
	n := rng.Intn(50) + 1
	b := make([]byte, n)
	for i := range b {
		if rng.Intn(2) == 0 {
			b[i] = 'v'
		} else {
			b[i] = 'o'
		}
	}
	return string(b)
}

func runCaseB(bin, s string) error {
	input := s + "\n"
	expected := solveB(s)
	out, err := run(bin, input)
	if err != nil {
		return fmt.Errorf("runtime error: %v", err)
	}
	if out != expected {
		return fmt.Errorf("expected %s got %s", expected, out)
	}
	return nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierB.go /path/to/binary")
		return
	}
	bin, cleanup, err := prepareBinary(os.Args[1])
	if err != nil {
		fmt.Println("compile error:", err)
		return
	}
	if cleanup != nil {
		defer cleanup()
	}
	rng := rand.New(rand.NewSource(2))
	for t := 0; t < numTestsB; t++ {
		s := genCaseB(rng)
		if err := runCaseB(bin, s); err != nil {
			fmt.Printf("case %d failed: %v\n", t+1, err)
			return
		}
	}
	fmt.Println("All tests passed")
}
