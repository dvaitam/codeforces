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

func solveCaseC(n, k int, str string) (int, string) {
	s := []byte(str)
	if k == 2 {
		pattern1 := make([]byte, n)
		pattern2 := make([]byte, n)
		changes1, changes2 := 0, 0
		for i := 0; i < n; i++ {
			pattern1[i] = byte('A' + byte(i%2))
			pattern2[i] = byte('A' + byte((i+1)%2))
			if s[i] != pattern1[i] {
				changes1++
			}
			if s[i] != pattern2[i] {
				changes2++
			}
		}
		if changes1 <= changes2 {
			return changes1, string(pattern1)
		}
		return changes2, string(pattern2)
	}
	changes := 0
	for i := 0; i+1 < n; i++ {
		if s[i] == s[i+1] {
			changes++
			for c := byte('A'); c < byte('A'+k); c++ {
				if c != s[i] && (i+2 >= n || c != s[i+2]) {
					s[i+1] = c
					break
				}
			}
		}
	}
	return changes, string(s)
}

func generateCaseC(rng *rand.Rand) (int, int, string) {
	n := rng.Intn(20) + 1
	k := rng.Intn(5) + 2 // at least 2
	b := make([]byte, n)
	for i := range b {
		b[i] = byte('A' + rng.Intn(k))
	}
	return n, k, string(b)
}

func runCaseC(bin string, n, k int, s string) error {
	input := fmt.Sprintf("%d %d\n%s\n", n, k, s)
	cmd := exec.Command(bin)
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("runtime error: %v\n%s", err, out.String())
	}
	fields := strings.Fields(strings.TrimSpace(out.String()))
	if len(fields) < 2 {
		return fmt.Errorf("output too short: %q", out.String())
	}
	var gotChanges int
	fmt.Sscan(fields[0], &gotChanges)
	repainted := fields[1]
	expChanges, expStr := solveCaseC(n, k, s)
	if gotChanges != expChanges || repainted != expStr {
		return fmt.Errorf("expected %d %s got %d %s", expChanges, expStr, gotChanges, repainted)
	}
	return nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierC.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))

	// deterministic case
	if err := runCaseC(bin, 1, 2, "A"); err != nil {
		fmt.Fprintln(os.Stderr, "deterministic case failed:", err)
		os.Exit(1)
	}

	for i := 0; i < 100; i++ {
		n, k, s := generateCaseC(rng)
		if err := runCaseC(bin, n, k, s); err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:\n%d %d\n%s\n", i+1, err, n, k, s)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
