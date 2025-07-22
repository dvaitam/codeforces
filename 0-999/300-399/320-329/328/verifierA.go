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

func expectedAnswerA(a [4]int64) int64 {
	d1 := a[1] - a[0]
	d2 := a[2] - a[1]
	d3 := a[3] - a[2]
	if d1 == d2 && d2 == d3 {
		return a[3] + d1
	}
	if a[0] != 0 && a[1]*a[1] == a[0]*a[2] && a[2]*a[2] == a[1]*a[3] {
		num := a[1]
		den := a[0]
		if (a[3]*num)%den == 0 {
			return a[3] * num / den
		}
	}
	return 42
}

func generateCase(rng *rand.Rand) (string, string) {
	typ := rng.Intn(4)
	var seq [4]int64
	for {
		switch typ {
		case 0: // arithmetic progression
			start := int64(rng.Intn(1000) + 1)
			d := int64(rng.Intn(21) - 10)
			seq[0] = start
			ok := true
			for i := 1; i < 4; i++ {
				start += d
				if start < 1 || start > 1000 {
					ok = false
					break
				}
				seq[i] = start
			}
			if !ok {
				continue
			}
		case 1: // geometric progression with integer ratio
			start := int64(rng.Intn(10) + 1)
			r := int64(rng.Intn(5) + 1)
			seq[0] = start
			ok := true
			for i := 1; i < 4; i++ {
				start = start * r
				if start < 1 || start > 1000 {
					ok = false
					break
				}
				seq[i] = start
			}
			if !ok {
				continue
			}
		case 2: // geometric progression ratio <1 so next term not integer
			options := [][4]int64{{8, 4, 2, 1}, {16, 8, 4, 2}, {81, 27, 9, 3}, {125, 25, 5, 1}, {64, 16, 4, 1}}
			seq = options[rng.Intn(len(options))]
		default: // random sequence
			for i := 0; i < 4; i++ {
				seq[i] = int64(rng.Intn(1000) + 1)
			}
		}
		break
	}
	input := fmt.Sprintf("%d %d %d %d\n", seq[0], seq[1], seq[2], seq[3])
	expected := fmt.Sprintf("%d\n", expectedAnswerA(seq))
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
