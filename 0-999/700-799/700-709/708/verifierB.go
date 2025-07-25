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

func solve(a00, a01, a10, a11 int64) string {
	zero, one := int64(-1), int64(-1)
	maxv := a00
	if a11 > maxv {
		maxv = a11
	}
	for i := int64(1); ; i++ {
		v := i * (i - 1) / 2
		if v > maxv {
			break
		}
		if v == a00 {
			zero = i
		}
		if v == a11 {
			one = i
		}
	}
	if a00 == 0 && a01 == 0 && a10 == 0 && a11 == 0 {
		return "0"
	}
	if a00 == 0 && a01 == 0 && a10 == 0 {
		if one == -1 {
			return "Impossible"
		}
		return strings.Repeat("1", int(one))
	}
	if a11 == 0 && a01 == 0 && a10 == 0 {
		if zero == -1 {
			return "Impossible"
		}
		return strings.Repeat("0", int(zero))
	}
	if zero == -1 || one == -1 {
		return "Impossible"
	}
	if zero*one != a01+a10 {
		return "Impossible"
	}
	b01 := zero * one
	lol := int64(0)
	for b01-a01 >= one {
		lol++
		b01 -= one
	}
	extra := b01 - a01
	var sb strings.Builder
	preZeros := zero - lol
	if extra > 0 {
		preZeros--
	}
	for i := int64(0); i < preZeros; i++ {
		sb.WriteByte('0')
	}
	for i := int64(0); i < extra; i++ {
		sb.WriteByte('1')
	}
	if extra > 0 {
		sb.WriteByte('0')
	}
	for i := int64(0); i < one-extra; i++ {
		sb.WriteByte('1')
	}
	for i := int64(0); i < lol; i++ {
		sb.WriteByte('0')
	}
	return sb.String()
}

func computeCounts(s string) (int64, int64, int64, int64) {
	var a00, a01, a10, a11 int64
	for i := 0; i < len(s); i++ {
		for j := i + 1; j < len(s); j++ {
			if s[i] == '0' && s[j] == '0' {
				a00++
			} else if s[i] == '0' && s[j] == '1' {
				a01++
			} else if s[i] == '1' && s[j] == '0' {
				a10++
			} else if s[i] == '1' && s[j] == '1' {
				a11++
			}
		}
	}
	return a00, a01, a10, a11
}

func runCase(bin string, a00, a01, a10, a11 int64) error {
	input := fmt.Sprintf("%d %d %d %d\n", a00, a01, a10, a11)
	cmd := exec.Command(bin)
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("runtime error: %v\n%s", err, out.String())
	}
	got := strings.TrimSpace(out.String())
	expect := solve(a00, a01, a10, a11)
	if got != expect {
		return fmt.Errorf("expected %q got %q", expect, got)
	}
	return nil
}

func randomValid(rng *rand.Rand) (int64, int64, int64, int64) {
	n := rng.Intn(20) + 1
	var s strings.Builder
	for i := 0; i < n; i++ {
		if rng.Intn(2) == 0 {
			s.WriteByte('0')
		} else {
			s.WriteByte('1')
		}
	}
	a00, a01, a10, a11 := computeCounts(s.String())
	return a00, a01, a10, a11
}

func randomCase(rng *rand.Rand) (int64, int64, int64, int64) {
	if rng.Intn(2) == 0 {
		return randomValid(rng)
	}
	a00 := rng.Int63n(20)
	a01 := rng.Int63n(20)
	a10 := rng.Int63n(20)
	a11 := rng.Int63n(20)
	return a00, a01, a10, a11
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierB.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for t := 0; t < 100; t++ {
		a00, a01, a10, a11 := randomCase(rng)
		if err := runCase(bin, a00, a01, a10, a11); err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:\n%d %d %d %d\n", t+1, err, a00, a01, a10, a11)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
