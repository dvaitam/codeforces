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

func expectedAnswerB(s string, k int, weights []int64) int64 {
	var sum int64
	for i, ch := range s {
		sum += weights[ch-'a'] * int64(i+1)
	}
	maxW := weights[0]
	for _, w := range weights {
		if w > maxW {
			maxW = w
		}
	}
	n := int64(len(s))
	sum += maxW * (int64(k)*n + int64(k*(k+1))/2)
	return sum
}

func generateCaseB(rng *rand.Rand) (string, int, []int64) {
	n := rng.Intn(20) + 1 // length of s 1..20
	b := make([]byte, n)
	for i := range b {
		b[i] = byte('a' + rng.Intn(26))
	}
	k := rng.Intn(50) // 0..49
	weights := make([]int64, 26)
	for i := range weights {
		weights[i] = int64(rng.Intn(1000))
	}
	return string(b), k, weights
}

func runCaseB(bin string, s string, k int, weights []int64) error {
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("%s\n%d\n", s, k))
	for i, w := range weights {
		if i > 0 {
			sb.WriteByte(' ')
		}
		sb.WriteString(fmt.Sprint(w))
	}
	sb.WriteByte('\n')
	cmd := exec.Command(bin)
	cmd.Stdin = strings.NewReader(sb.String())
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("runtime error: %v\n%s", err, out.String())
	}
	got := strings.TrimSpace(out.String())
	expected := fmt.Sprint(expectedAnswerB(s, k, weights))
	if got != expected {
		return fmt.Errorf("expected %s got %s", expected, got)
	}
	return nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierB.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < 100; i++ {
		s, k, weights := generateCaseB(rng)
		if err := runCaseB(bin, s, k, weights); err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:\n%s\n%d\n%v\n", i+1, err, s, k, weights)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
