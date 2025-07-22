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

func expectedD(mat [][]int64) int64 {
	n := len(mat)
	x := (n + 1) / 2
	var sumAbs int64
	cntNeg := 0
	minAbs := int64(-1)
	for i := 0; i < n; i++ {
		for j := 0; j < n; j++ {
			v := mat[i][j]
			if v < 0 {
				cntNeg++
				v = -v
			}
			sumAbs += v
			if minAbs < 0 || v < minAbs {
				minAbs = v
			}
		}
	}
	if x%2 == 0 && cntNeg%2 == 1 {
		sumAbs -= 2 * minAbs
	}
	return sumAbs
}

func generateCase(rng *rand.Rand) (string, int64) {
	n := rng.Intn(3)*2 + 1 // 1,3,5
	mat := make([][]int64, n)
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("%d\n", n))
	for i := 0; i < n; i++ {
		mat[i] = make([]int64, n)
		for j := 0; j < n; j++ {
			v := int64(rng.Intn(11) - 5)
			mat[i][j] = v
			sb.WriteString(fmt.Sprintf("%d", v))
			if j+1 < n {
				sb.WriteByte(' ')
			}
		}
		sb.WriteByte('\n')
	}
	expect := expectedD(mat)
	return sb.String(), expect
}

func runCase(bin string, input string, expect int64) error {
	cmd := exec.Command(bin)
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	var stderr bytes.Buffer
	cmd.Stderr = &stderr
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("runtime error: %v\nstderr: %s", err, stderr.String())
	}
	resStr := strings.TrimSpace(out.String())
	var got int64
	if _, err := fmt.Sscan(resStr, &got); err != nil {
		return fmt.Errorf("bad output %q", resStr)
	}
	if got != expect {
		return fmt.Errorf("expected %d got %d", expect, got)
	}
	return nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierD.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < 100; i++ {
		in, expect := generateCase(rng)
		if err := runCase(bin, in, expect); err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:\n%s", i+1, err, in)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
