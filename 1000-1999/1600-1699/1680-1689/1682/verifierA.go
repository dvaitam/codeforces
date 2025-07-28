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
	"time"
)

func prepareBinary(path string) (string, func(), error) {
	if strings.HasSuffix(path, ".go") {
		tmp := filepath.Join(os.TempDir(), "verifA_bin")
		cmd := exec.Command("go", "build", "-o", tmp, path)
		if out, err := cmd.CombinedOutput(); err != nil {
			return "", nil, fmt.Errorf("go build failed: %v: %s", err, out)
		}
		return tmp, func() { os.Remove(tmp) }, nil
	}
	return path, nil, nil
}

func expected(n int, s string) int {
	if n%2 == 1 {
		mid := n / 2
		ch := s[mid]
		ans := 1
		for i := mid - 1; i >= 0 && s[i] == ch; i-- {
			ans++
		}
		for j := mid + 1; j < n && s[j] == ch; j++ {
			ans++
		}
		return ans
	}
	left := n/2 - 1
	ch := s[left]
	ans := 0
	for i := left; i >= 0 && s[i] == ch; i-- {
		ans++
	}
	for j := left + 1; j < n && s[j] == ch; j++ {
		ans++
	}
	return ans
}

func generateCase(rng *rand.Rand) (string, int) {
	n := rng.Intn(10) + 2
	b := make([]byte, n)
	for i := 0; i < (n+1)/2; i++ {
		ch := byte('a' + rng.Intn(3))
		b[i] = ch
		b[n-1-i] = ch
	}
	s := string(b)
	var sb strings.Builder
	fmt.Fprintf(&sb, "1\n%d\n%s\n", n, s)
	return sb.String(), expected(n, s)
}

func runCase(bin, input string, exp int) error {
	cmd := exec.Command(bin)
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	var errBuf bytes.Buffer
	cmd.Stderr = &errBuf
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("runtime error: %v\n%s", err, errBuf.String())
	}
	gotStr := strings.TrimSpace(out.String())
	got, err := strconv.Atoi(gotStr)
	if err != nil {
		return fmt.Errorf("invalid output %q", gotStr)
	}
	if got != exp {
		return fmt.Errorf("expected %d got %d", exp, got)
	}
	return nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierA.go /path/to/binary")
		return
	}
	bin, clean, err := prepareBinary(os.Args[1])
	if err != nil {
		fmt.Println("compile error:", err)
		return
	}
	if clean != nil {
		defer clean()
	}
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < 100; i++ {
		input, exp := generateCase(rng)
		if err := runCase(bin, input, exp); err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:\n%s", i+1, err, input)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
