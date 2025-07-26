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

func solve(xs []int64, cs []byte) int64 {
	n := len(xs)
	var pIdx []int
	for i, c := range cs {
		if c == 'P' {
			pIdx = append(pIdx, i)
		}
	}
	if len(pIdx) == 0 {
		var ans int64
		var lastB int64
		var haveB bool
		var lastR int64
		var haveR bool
		for i := 0; i < n; i++ {
			switch cs[i] {
			case 'B':
				if haveB {
					ans += xs[i] - lastB
				}
				lastB = xs[i]
				haveB = true
			case 'R':
				if haveR {
					ans += xs[i] - lastR
				}
				lastR = xs[i]
				haveR = true
			}
		}
		return ans
	}
	ans := int64(0)
	firstP := pIdx[0]
	lastP := pIdx[len(pIdx)-1]
	firstB, firstR := int64(0), int64(0)
	haveB := false
	haveR := false
	for i := 0; i < firstP; i++ {
		if cs[i] == 'B' && !haveB {
			firstB = xs[i]
			haveB = true
		}
		if cs[i] == 'R' && !haveR {
			firstR = xs[i]
			haveR = true
		}
	}
	if haveB {
		ans += xs[firstP] - firstB
	}
	if haveR {
		ans += xs[firstP] - firstR
	}
	lastBpos, lastRpos := int64(0), int64(0)
	haveB = false
	haveR = false
	for i := n - 1; i > lastP; i-- {
		if cs[i] == 'B' && !haveB {
			lastBpos = xs[i]
			haveB = true
		}
		if cs[i] == 'R' && !haveR {
			lastRpos = xs[i]
			haveR = true
		}
	}
	if haveB {
		ans += lastBpos - xs[lastP]
	}
	if haveR {
		ans += lastRpos - xs[lastP]
	}
	for idx := 0; idx < len(pIdx)-1; idx++ {
		i := pIdx[idx]
		j := pIdx[idx+1]
		L := xs[i]
		R := xs[j]
		seg := R - L
		maxB := int64(0)
		prev := L
		var hasB bool
		for t := i + 1; t < j; t++ {
			if cs[t] == 'B' {
				hasB = true
				if xs[t]-prev > maxB {
					maxB = xs[t] - prev
				}
				prev = xs[t]
			}
		}
		if R-prev > maxB {
			maxB = R - prev
		}
		maxR := int64(0)
		prev = L
		var hasR bool
		for t := i + 1; t < j; t++ {
			if cs[t] == 'R' {
				hasR = true
				if xs[t]-prev > maxR {
					maxR = xs[t] - prev
				}
				prev = xs[t]
			}
		}
		if R-prev > maxR {
			maxR = R - prev
		}
		cost2 := 3*seg - maxB - maxR
		if hasB && hasR {
			cost1 := 2 * seg
			if cost1 < cost2 {
				ans += cost1
			} else {
				ans += cost2
			}
		} else {
			ans += cost2
		}
	}
	return ans
}

func genCase(rng *rand.Rand) (string, string) {
	n := rng.Intn(6) + 2
	xs := make([]int64, n)
	cs := make([]byte, n)
	cur := int64(0)
	for i := 0; i < n; i++ {
		cur += int64(rng.Intn(5) + 1)
		xs[i] = cur
		switch rng.Intn(3) {
		case 0:
			cs[i] = 'B'
		case 1:
			cs[i] = 'R'
		default:
			cs[i] = 'P'
		}
	}
	input := fmt.Sprintf("%d\n", n)
	for i := 0; i < n; i++ {
		input += fmt.Sprintf("%d %c\n", xs[i], cs[i])
	}
	out := solve(xs, cs)
	expected := fmt.Sprintf("%d\n", out)
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
	if strings.TrimSpace(out.String()) != strings.TrimSpace(expected) {
		return fmt.Errorf("expected %s got %s", expected, out.String())
	}
	return nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierE.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < 100; i++ {
		in, exp := genCase(rng)
		if err := runCase(bin, in, exp); err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:\n%s", i+1, err, in)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
