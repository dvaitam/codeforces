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

func solve(s string) (string, bool) {
	L := len(s)
	digits := make([]int, L)
	for i := 0; i < L; i++ {
		c := s[L-1-i]
		switch {
		case c >= '0' && c <= '9':
			digits[i] = int(c - '0')
		case c >= 'a' && c <= 'f':
			digits[i] = int(c-'a') + 10
		case c >= 'A' && c <= 'F':
			digits[i] = int(c-'A') + 10
		}
	}
	best := uint64(^uint64(0))
	found := false
	for _, sign := range []int{1, -1} {
		type key struct {
			pos, carry int
			diff       [16]int
		}
		memo := make(map[key]uint64)
		var dfs func(int, int, [16]int) (uint64, bool)
		dfs = func(pos, carry int, d [16]int) (uint64, bool) {
			k := key{pos, carry, d}
			if v, ok := memo[k]; ok {
				if v == ^uint64(0) {
					return 0, false
				}
				return v, true
			}
			if pos == L {
				if carry == 0 {
					for _, x := range d {
						if x != 0 {
							memo[k] = ^uint64(0)
							return 0, false
						}
					}
					memo[k] = 0
					return 0, true
				}
				memo[k] = ^uint64(0)
				return 0, false
			}
			rem := L - pos
			sumAbs := 0
			for _, x := range d {
				if x > rem || -x > rem {
					memo[k] = ^uint64(0)
					return 0, false
				}
				if x >= 0 {
					sumAbs += x
				} else {
					sumAbs -= x
				}
			}
			if sumAbs > 2*rem {
				memo[k] = ^uint64(0)
				return 0, false
			}
			bestHere := uint64(^uint64(0))
			good := false
			for a := 0; a < 16; a++ {
				t := 0
				nextCarry := 0
				var b int
				if sign == 1 {
					t = a + digits[pos] + carry
					b = t % 16
					nextCarry = t / 16
				} else {
					t = a - digits[pos] - carry
					if t < 0 {
						t += 16
						nextCarry = 1
					}
					b = t
				}
				d[a]++
				d[b]--
				if val, ok := dfs(pos+1, nextCarry, d); ok {
					cand := val + uint64(a)<<uint(4*pos)
					if cand < bestHere {
						bestHere = cand
					}
					good = true
				}
				d[a]--
				d[b]++
			}
			if good {
				memo[k] = bestHere
				return bestHere, true
			}
			memo[k] = ^uint64(0)
			return 0, false
		}
		val, ok := dfs(0, 0, [16]int{})
		if ok {
			if !found || val < best {
				best = val
				found = true
			}
		}
	}
	if !found {
		return "", false
	}
	out := make([]byte, L)
	for i := 0; i < L; i++ {
		d := (best >> uint(4*i)) & 15
		if d < 10 {
			out[L-1-i] = byte('0' + d)
		} else {
			out[L-1-i] = byte('a' + d - 10)
		}
	}
	return string(out), true
}

func run(bin, input string) (string, error) {
	var cmd *exec.Cmd
	if strings.HasSuffix(bin, ".go") {
		cmd = exec.Command("go", "run", bin)
	} else {
		cmd = exec.Command(bin)
	}
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	var stderr bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &stderr
	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("runtime error: %v\n%s", err, stderr.String())
	}
	return strings.TrimSpace(out.String()), nil
}

func main() {
	if len(os.Args) == 3 && os.Args[1] == "--" {
		os.Args = append([]string{os.Args[0]}, os.Args[2])
	}
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierE.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	digits := []byte("0123456789abcdef")
	for t := 0; t < 100; t++ {
		L := rng.Intn(5) + 1
		b := make([]byte, L)
		for i := 0; i < L; i++ {
			b[i] = digits[rng.Intn(16)]
		}
		if strings.TrimLeft(string(b), "0") == "" {
			b[0] = digits[rng.Intn(15)+1]
		}
		input := string(b)
		expectedVal, ok := solve(input)
		expected := ""
		if ok {
			expected = expectedVal
		} else {
			expected = "NO"
		}
		got, err := run(bin, input+"\n")
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\n", t+1, err)
			os.Exit(1)
		}
		if strings.TrimSpace(got) != strings.TrimSpace(expected) {
			fmt.Fprintf(os.Stderr, "case %d failed\ninput:\n%s\nexpected:\n%s\n---\ngot:\n%s\n", t+1, input, expected, got)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
