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

func expectedMoves(n int, d int64, arr []int64) int64 {
	var moves int64
	var last int64 = -1
	for i := 0; i < n; i++ {
		cur := arr[i]
		if last >= cur {
			diff := last - cur
			k := diff/d + 1
			moves += k
			cur += k * d
		}
		last = cur
	}
	return moves
}

func generateCase(rng *rand.Rand) (string, int64) {
	n := rng.Intn(19) + 2 // 2..20
	d := int64(rng.Intn(1000) + 1)
	arr := make([]int64, n)
	for i := range arr {
		arr[i] = int64(rng.Intn(2000))
	}
	arrCopy := make([]int64, n)
	copy(arrCopy, arr)
	exp := expectedMoves(n, d, arrCopy)
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("%d %d\n", n, d))
	for i, v := range arr {
		if i > 0 {
			sb.WriteByte(' ')
		}
		sb.WriteString(fmt.Sprintf("%d", v))
	}
	sb.WriteByte('\n')
	return sb.String(), exp
}

func runCase(bin, input string, expected int64) error {
	cmd := exec.Command(bin)
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("runtime error: %v\n%s", err, out.String())
	}
	var got int64
	if _, err := fmt.Fscan(strings.NewReader(out.String()), &got); err != nil {
		return fmt.Errorf("bad output: %v", err)
	}
	if got != expected {
		return fmt.Errorf("expected %d got %d", expected, got)
	}
	return nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierA.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < 100; i++ {
		in, exp := generateCase(rng)
		if err := runCase(bin, in, exp); err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:\n%s", i+1, err, in)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
