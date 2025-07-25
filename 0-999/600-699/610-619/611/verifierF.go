package main

import (
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"strconv"
	"strings"
)

func runCandidate(bin string, input string) (string, error) {
	var cmd *exec.Cmd
	if strings.HasSuffix(bin, ".go") {
		cmd = exec.Command("go", "run", bin)
	} else {
		cmd = exec.Command(bin)
	}
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	var stderr bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &stderr
	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("runtime error: %v\n%s", err, stderr.String())
	}
	return strings.TrimSpace(out.String()), nil
}

const MOD int64 = 1000000007

func expected(n, h, w int, s string) int64 {
	remW := int64(w)
	remH := int64(h)
	dx, dy := 0, 0
	minX, maxX := 0, 0
	minY, maxY := 0, 0
	ans := int64(0)
	step := 0
	idx := 0
	for remW > 0 && remH > 0 {
		ch := s[idx]
		idx++
		if idx == n {
			idx = 0
		}
		step++
		switch ch {
		case 'L':
			dx--
		case 'R':
			dx++
		case 'U':
			dy--
		case 'D':
			dy++
		}
		changed := false
		if dx < minX {
			minX = dx
			remW--
			ans = (ans + remH*int64(step)) % MOD
			changed = true
		}
		if dx > maxX && remW > 0 {
			maxX = dx
			remW--
			ans = (ans + remH*int64(step)) % MOD
			changed = true
		}
		if dy < minY && remH > 0 {
			minY = dy
			remH--
			ans = (ans + remW*int64(step)) % MOD
			changed = true
		}
		if dy > maxY && remH > 0 {
			maxY = dy
			remH--
			ans = (ans + remW*int64(step)) % MOD
			changed = true
		}
		if !changed && idx == 0 && dx == 0 && dy == 0 {
			return -1
		}
	}
	return ans % MOD
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierF.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rng := rand.New(rand.NewSource(6))
	letters := []byte("LRUD")
	for t := 0; t < 100; t++ {
		n := rng.Intn(20) + 1
		h := rng.Intn(10) + 1
		w := rng.Intn(10) + 1
		var sb strings.Builder
		for i := 0; i < n; i++ {
			sb.WriteByte(letters[rng.Intn(4)])
		}
		s := sb.String()
		input := fmt.Sprintf("%d %d %d\n%s\n", n, h, w, s)
		want := expected(n, h, w, s)
		out, err := runCandidate(bin, input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\n", t+1, err)
			os.Exit(1)
		}
		got, err := strconv.ParseInt(strings.TrimSpace(out), 10, 64)
		if err != nil || got != want {
			fmt.Fprintf(os.Stderr, "case %d failed: expected %d got %s\n", t+1, want, out)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
