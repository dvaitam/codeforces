package main

import (
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"time"
)

func generateChain(n int) [][]int {
	perm := make([]int, n)
	for i := 0; i < n; i++ {
		perm[i] = i + 1
	}
	res := make([][]int, 0, n)
	copyPerm := func() []int {
		p := make([]int, n)
		copy(p, perm)
		return p
	}
	res = append(res, copyPerm())
	perm[0], perm[n-1] = perm[n-1], perm[0]
	res = append(res, copyPerm())
	for i := 3; i <= n; i++ {
		j1 := i - 3
		j2 := i - 2
		perm[j1], perm[j2] = perm[j2], perm[j1]
		res = append(res, copyPerm())
	}
	return res
}

func runCandidate(bin string, input string) (string, error) {
	var cmd *exec.Cmd
	if strings.HasSuffix(bin, ".go") {
		cmd = exec.Command("go", "run", bin)
	} else {
		cmd = exec.Command(bin)
	}
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	var errBuf bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &errBuf
	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("runtime error: %v\n%s", err, errBuf.String())
	}
	return strings.TrimSpace(out.String()), nil
}

func verifyCase(bin string, n int) error {
	input := fmt.Sprintf("1\n%d\n", n)
	out, err := runCandidate(bin, input)
	if err != nil {
		return err
	}
	fields := strings.Fields(out)
	if len(fields) < 1 {
		return fmt.Errorf("no output")
	}
	k, err := strconv.Atoi(fields[0])
	if err != nil {
		return fmt.Errorf("invalid k: %v", err)
	}
	if k != n {
		return fmt.Errorf("expected k=%d got %d", n, k)
	}
	if len(fields) != 1+k*n {
		return fmt.Errorf("expected %d numbers, got %d", 1+k*n, len(fields))
	}
	perms := generateChain(n)
	idx := 1
	for i := 0; i < k; i++ {
		for j := 0; j < n; j++ {
			v, err := strconv.Atoi(fields[idx])
			if err != nil {
				return fmt.Errorf("invalid int: %v", err)
			}
			if v != perms[i][j] {
				return fmt.Errorf("mismatch at perm %d index %d: expected %d got %d", i+1, j+1, perms[i][j], v)
			}
			idx++
		}
	}
	return nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierB.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < 100; i++ {
		n := rng.Intn(99) + 2 // 2..100
		if err := verifyCase(bin, n); err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\n", i+1, err)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
