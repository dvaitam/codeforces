package main

import (
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"time"
)

type plane [][]byte

func solveB(n, k int, rows []string) (int, []string) {
	s := make([][]byte, n)
	for i := 0; i < n; i++ {
		s[i] = []byte(rows[i])
	}
	kLeft := k
	for i := 0; i < n; i++ {
		if kLeft == 0 {
			break
		}
		if s[i][0] == '.' && s[i][1] != 'S' {
			s[i][0] = 'x'
			kLeft--
		}
		if kLeft == 0 {
			break
		}
		if s[i][11] == '.' && s[i][10] != 'S' {
			s[i][11] = 'x'
			kLeft--
		}
		if kLeft == 0 {
			break
		}
		for j := 1; j < 11 && kLeft > 0; j++ {
			if s[i][j] == '.' && s[i][j-1] != 'S' && s[i][j+1] != 'S' {
				s[i][j] = 'x'
				kLeft--
			}
		}
	}

	for i := 0; i < n; i++ {
		if kLeft == 0 {
			break
		}
		if s[i][0] == '.' {
			s[i][0] = 'x'
			kLeft--
		}
		if kLeft == 0 {
			break
		}
		if s[i][11] == '.' {
			s[i][11] = 'x'
			kLeft--
		}
		if kLeft == 0 {
			break
		}
		for j := 1; j < 11 && kLeft > 0; j++ {
			if s[i][j] == '.' && !(s[i][j-1] == 'S' && s[i][j+1] == 'S') {
				s[i][j] = 'x'
				kLeft--
			}
		}
	}

	for i := 0; i < n && kLeft > 0; i++ {
		for j := 0; j < 12 && kLeft > 0; j++ {
			if s[i][j] == '.' {
				s[i][j] = 'x'
				kLeft--
			}
		}
	}

	ans := 0
	for i := 0; i < n; i++ {
		if s[i][0] == 'S' {
			if s[i][1] != '.' && s[i][1] != '-' {
				ans++
			}
		}
		if s[i][11] == 'S' {
			if s[i][10] != '.' && s[i][10] != '-' {
				ans++
			}
		}
		for j := 1; j < 11; j++ {
			if s[i][j] == 'S' {
				if s[i][j-1] != '.' && s[i][j-1] != '-' {
					ans++
				}
				if s[i][j+1] != '.' && s[i][j+1] != '-' {
					ans++
				}
			}
		}
	}

	out := make([]string, n)
	for i := 0; i < n; i++ {
		out[i] = string(s[i])
	}
	return ans, out
}

func generateCase(rng *rand.Rand) (string, string) {
	n := rng.Intn(20) + 1
	k := rng.Intn(10*n + 1)
	rows := make([]string, n)
	for i := 0; i < n; i++ {
		b := make([]byte, 12)
		for j := 0; j < 12; j++ {
			if j == 3 || j == 8 {
				b[j] = '-'
			} else {
				if rng.Intn(5) == 0 {
					b[j] = 'S'
				} else {
					b[j] = '.'
				}
			}
		}
		rows[i] = string(b)
	}
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("%d %d\n", n, k))
	for i := 0; i < n; i++ {
		sb.WriteString(rows[i])
		sb.WriteByte('\n')
	}
	ans, outRows := solveB(n, k, rows)
	var expect strings.Builder
	expect.WriteString(strconv.Itoa(ans))
	expect.WriteByte('\n')
	for i := 0; i < n; i++ {
		expect.WriteString(outRows[i])
		expect.WriteByte('\n')
	}
	return sb.String(), expect.String()
}

func runCase(exe, input, expected string) error {
	cmd := exec.Command(exe)
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("runtime error: %v\n%s", err, out.String())
	}
	if strings.TrimSpace(out.String()) != strings.TrimSpace(expected) {
		return fmt.Errorf("expected:\n%s\ngot:\n%s", expected, out.String())
	}
	return nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("Usage: go run verifierB.go /path/to/binary")
		os.Exit(1)
	}
	exe := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < 100; i++ {
		in, exp := generateCase(rng)
		if err := runCase(exe, in, exp); err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:\n%s", i+1, err, in)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
