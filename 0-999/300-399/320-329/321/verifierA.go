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

func reachable(rx, ry, dx, dy int64) bool {
	if dx == 0 && dy == 0 {
		return rx == 0 && ry == 0
	}
	if dx == 0 {
		if rx != 0 || dy == 0 || ry%dy != 0 {
			return false
		}
		k := ry / dy
		return k >= 0
	}
	if dy == 0 {
		if ry != 0 || dx == 0 || rx%dx != 0 {
			return false
		}
		k := rx / dx
		return k >= 0
	}
	if rx%dx != 0 || ry%dy != 0 {
		return false
	}
	kx := rx / dx
	ky := ry / dy
	return kx == ky && kx >= 0
}

func expectedA(a, b int64, s string) bool {
	var dx, dy int64
	for _, ch := range s {
		switch ch {
		case 'U':
			dy++
		case 'D':
			dy--
		case 'L':
			dx--
		case 'R':
			dx++
		}
	}
	var px, py int64
	for i := 0; i <= len(s); i++ {
		rx := a - px
		ry := b - py
		if reachable(rx, ry, dx, dy) {
			return true
		}
		if i < len(s) {
			switch s[i] {
			case 'U':
				py++
			case 'D':
				py--
			case 'L':
				px--
			case 'R':
				px++
			}
		}
	}
	return false
}

func runCase(bin string, input string, expect bool) error {
	cmd := exec.Command(bin)
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	var stderr bytes.Buffer
	cmd.Stderr = &stderr
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("runtime error: %v\nstderr: %s", err, stderr.String())
	}
	res := strings.TrimSpace(out.String())
	want := "No"
	if expect {
		want = "Yes"
	}
	if res != want {
		return fmt.Errorf("expected %q got %q", want, res)
	}
	return nil
}

func generateCase(rng *rand.Rand) (string, bool) {
	a := int64(rng.Intn(21) - 10)
	b := int64(rng.Intn(21) - 10)
	l := rng.Intn(15) + 1
	var sb strings.Builder
	moves := make([]byte, l)
	for i := 0; i < l; i++ {
		switch rng.Intn(4) {
		case 0:
			moves[i] = 'U'
		case 1:
			moves[i] = 'D'
		case 2:
			moves[i] = 'L'
		default:
			moves[i] = 'R'
		}
	}
	s := string(moves)
	sb.WriteString(fmt.Sprintf("%d %d\n", a, b))
	sb.WriteString(s)
	sb.WriteByte('\n')
	expect := expectedA(a, b, s)
	return sb.String(), expect
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierA.go /path/to/binary")
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
