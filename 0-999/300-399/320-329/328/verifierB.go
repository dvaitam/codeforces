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

func expectedAnswerB(t string, pieces string) int {
	avail := make([]int, 10)
	for i := 0; i < len(pieces); i++ {
		avail[pieces[i]-'0']++
	}
	req := make([]int, 10)
	req69 := 0
	req25 := 0
	for i := 0; i < len(t); i++ {
		switch t[i] {
		case '6', '9':
			req69++
		case '2', '5':
			req25++
		default:
			req[t[i]-'0']++
		}
	}
	maxK := len(pieces)
	for d := 0; d <= 9; d++ {
		if req[d] > 0 {
			k := avail[d] / req[d]
			if k < maxK {
				maxK = k
			}
		}
	}
	if req69 > 0 {
		k := (avail[6] + avail[9]) / req69
		if k < maxK {
			maxK = k
		}
	}
	if req25 > 0 {
		k := (avail[2] + avail[5]) / req25
		if k < maxK {
			maxK = k
		}
	}
	return maxK
}

func generateCase(rng *rand.Rand) (string, string) {
	tLen := rng.Intn(15) + 1
	var tb strings.Builder
	for i := 0; i < tLen; i++ {
		tb.WriteByte(byte('0' + rng.Intn(10)))
	}
	piecesLen := rng.Intn(50) + 1
	var pb strings.Builder
	for i := 0; i < piecesLen; i++ {
		pb.WriteByte(byte('0' + rng.Intn(10)))
	}
	t := tb.String()
	pieces := pb.String()
	input := fmt.Sprintf("%s\n%s\n", t, pieces)
	expected := fmt.Sprintf("%d\n", expectedAnswerB(t, pieces))
	return input, expected
}

func runCase(bin, input, expected string) error {
	cmd := exec.Command(bin)
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("runtime error: %v\n%s", err, out.String())
	}
	got := strings.TrimSpace(out.String())
	exp := strings.TrimSpace(expected)
	if got != exp {
		return fmt.Errorf("expected %s got %s", exp, got)
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
		in, exp := generateCase(rng)
		if err := runCase(bin, in, exp); err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:\n%s", i+1, err, in)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
