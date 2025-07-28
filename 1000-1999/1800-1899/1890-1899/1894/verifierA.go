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

func winnerFor(X, Y int, s string) (byte, bool) {
	setsA, setsB := 0, 0
	winsA, winsB := 0, 0
	for i := 0; i < len(s); i++ {
		if setsA == Y || setsB == Y {
			return 0, false
		}
		if s[i] == 'A' {
			winsA++
		} else {
			winsB++
		}
		if winsA == X {
			setsA++
			winsA, winsB = 0, 0
			if setsA == Y {
				if i != len(s)-1 {
					return 0, false
				}
				return 'A', true
			}
		} else if winsB == X {
			setsB++
			winsA, winsB = 0, 0
			if setsB == Y {
				if i != len(s)-1 {
					return 0, false
				}
				return 'B', true
			}
		}
	}
	return 0, false
}

func expectedResult(s string) string {
	n := len(s)
	winners := make(map[byte]struct{})
	for x := 1; x <= n; x++ {
		for y := 1; y <= n; y++ {
			if w, ok := winnerFor(x, y, s); ok {
				winners[w] = struct{}{}
			}
		}
	}
	if len(winners) == 1 {
		for k := range winners {
			return string(k) + "\n"
		}
	}
	return "?\n"
}

func generateCase(rng *rand.Rand) (string, string) {
	for {
		X := rng.Intn(3) + 1
		Y := rng.Intn(3) + 1
		var sb strings.Builder
		setsA, setsB := 0, 0
		for setsA < Y && setsB < Y {
			winsA, winsB := 0, 0
			for winsA < X && winsB < X {
				if sb.Len() >= 20 {
					break
				}
				if rng.Intn(2) == 0 {
					sb.WriteByte('A')
					winsA++
				} else {
					sb.WriteByte('B')
					winsB++
				}
			}
			if winsA == X {
				setsA++
			} else if winsB == X {
				setsB++
			} else {
				break
			}
			if sb.Len() > 20 {
				break
			}
		}
		if (setsA == Y || setsB == Y) && sb.Len() <= 20 {
			s := sb.String()
			input := fmt.Sprintf("1\n%d\n%s\n", len(s), s)
			expected := expectedResult(s)
			return input, expected
		}
	}
}

func runCase(bin, input, expected string) error {
	var cmd *exec.Cmd
	if strings.HasSuffix(bin, ".go") {
		cmd = exec.Command("go", "run", bin)
	} else {
		cmd = exec.Command(bin)
	}
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("runtime error: %v\n%s", err, out.String())
	}
	outStr := strings.TrimSpace(out.String())
	exp := strings.TrimSpace(expected)
	if outStr != exp {
		return fmt.Errorf("expected %q got %q", exp, outStr)
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
