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

func isImpossible(n int, s string) bool {
	if n%4 != 0 {
		return true
	}
	count := map[rune]int{'A': 0, 'C': 0, 'G': 0, 'T': 0}
	q := 0
	for _, ch := range s {
		if ch == '?' {
			q++
		} else {
			count[ch]++
		}
	}
	target := n / 4
	deficit := 0
	for _, c := range []rune{'A', 'C', 'G', 'T'} {
		if count[c] > target {
			return true
		}
		deficit += target - count[c]
	}
	return deficit != q
}

func checkCase(bin string, n int, s string) error {
	input := fmt.Sprintf("%d\n%s\n", n, s)
	cmd := exec.Command(bin)
	cmd.Stdin = strings.NewReader(input)
	out, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("runtime error: %v", err)
	}
	res := strings.TrimSpace(string(out))
	impossible := isImpossible(n, s)
	if res == "===" {
		if impossible {
			return nil
		}
		return fmt.Errorf("expected valid string, got ===")
	}
	if impossible {
		return fmt.Errorf("expected ===, got %s", res)
	}
	if len(res) != n {
		return fmt.Errorf("length mismatch: expected %d got %d", n, len(res))
	}
	count := map[rune]int{'A': 0, 'C': 0, 'G': 0, 'T': 0}
	for i, ch := range res {
		if ch != 'A' && ch != 'C' && ch != 'G' && ch != 'T' {
			return fmt.Errorf("invalid character %c", ch)
		}
		if s[i] != '?' && rune(s[i]) != ch {
			return fmt.Errorf("changed fixed position")
		}
		count[ch]++
	}
	target := n / 4
	for _, c := range []rune{'A', 'C', 'G', 'T'} {
		if count[c] != target {
			return fmt.Errorf("count mismatch for %c", c)
		}
	}
	return nil
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: go run verifierB.go /path/to/binary")
		return
	}
	bin := os.Args[1]
	rand.Seed(1)
	for tcase := 1; tcase <= 120; tcase++ {
		n := rand.Intn(50) + 4
		letters := []byte{'A', 'C', 'G', 'T', '?'}
		sb := make([]byte, n)
		for i := range sb {
			sb[i] = letters[rand.Intn(len(letters))]
		}
		if err := checkCase(bin, n, string(sb)); err != nil {
			fmt.Printf("Test %d failed: %v\n", tcase, err)
			return
		}
		time.Sleep(0) // yield for determinism
	}
	fmt.Println("OK")
}
