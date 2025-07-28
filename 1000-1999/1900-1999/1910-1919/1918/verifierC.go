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

func abs(x int64) int64 {
	if x < 0 {
		return -x
	}
	return x
}

func solveCase(a, b, r int64) int64 {
	const maxBit = 60
	var w [maxBit + 1]int64
	var absW [maxBit + 1]int64
	for i := 0; i <= maxBit; i++ {
		aBit := (a >> i) & 1
		bBit := (b >> i) & 1
		if aBit != bBit {
			if aBit > bBit {
				w[i] = 1 << (i + 1)
			} else {
				w[i] = -1 << (i + 1)
			}
		} else {
			w[i] = 0
		}
		if w[i] >= 0 {
			absW[i] = w[i]
		} else {
			absW[i] = -w[i]
		}
	}

	var prefAll [maxBit + 2]int64
	var prefLim [maxBit + 2]int64
	for i := 1; i <= maxBit+1; i++ {
		prefAll[i] = prefAll[i-1] + absW[i-1]
		if (r>>(i-1))&1 == 1 {
			prefLim[i] = prefLim[i-1] + absW[i-1]
		} else {
			prefLim[i] = prefLim[i-1]
		}
	}

	delta := a - b
	prefixLess := false
	for i := maxBit; i >= 0; i-- {
		remAll := prefAll[i]
		remLim := prefLim[i]
		if prefixLess {
			delta0 := delta
			best0 := abs(delta0) - remAll
			if best0 < 0 {
				best0 = 0
			}
			delta1 := delta - w[i]
			best1 := abs(delta1) - remAll
			if best1 < 0 {
				best1 = 0
			}
			if best1 < best0 || (best1 == best0 && abs(delta1) < abs(delta0)) {
				delta = delta1
			}
		} else {
			if (r>>i)&1 == 0 {
				continue
			}
			deltaEqual := delta - w[i]
			bestEqual := abs(deltaEqual) - remLim
			if bestEqual < 0 {
				bestEqual = 0
			}
			deltaLower := delta
			bestLower := abs(deltaLower) - remAll
			if bestLower < 0 {
				bestLower = 0
			}
			if bestLower < bestEqual || (bestLower == bestEqual && abs(deltaLower) < abs(deltaEqual)) {
				prefixLess = true
			} else {
				delta = deltaEqual
			}
		}
	}
	if delta < 0 {
		delta = -delta
	}
	return delta
}

func runCase(bin string, a, b, r int64) error {
	input := fmt.Sprintf("1\n%d %d %d\n", a, b, r)
	cmd := exec.Command(bin)
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	var errBuf bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &errBuf
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("%v\n%s", err, errBuf.String())
	}
	got := strings.TrimSpace(out.String())
	expected := fmt.Sprintf("%d", solveCase(a, b, r))
	if got != expected {
		return fmt.Errorf("expected %s got %s", expected, got)
	}
	return nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierC.go /path/to/binary")
		os.Exit(1)
	}
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	bin := os.Args[1]
	// deterministic
	tests := [][3]int64{{4, 6, 0}}
	for i := 0; i < 99; i++ {
		a := rng.Int63n(1 << 60)
		b := rng.Int63n(1 << 60)
		r := rng.Int63n(1 << 60)
		tests = append(tests, [3]int64{a, b, r})
	}
	for i, t := range tests {
		if err := runCase(bin, t[0], t[1], t[2]); err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\n", i+1, err)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
