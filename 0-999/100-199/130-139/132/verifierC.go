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

func expectedDistance(cmds string, flips int) int {
	m := len(cmds)
	offset := m
	maxP := 2*m + 1
	dpPrev := make([][][]bool, flips+1)
	dpCurr := make([][][]bool, flips+1)
	for j := 0; j <= flips; j++ {
		dpPrev[j] = make([][]bool, 2)
		dpCurr[j] = make([][]bool, 2)
		for d := 0; d < 2; d++ {
			dpPrev[j][d] = make([]bool, maxP)
			dpCurr[j][d] = make([]bool, maxP)
		}
	}
	dpPrev[0][0][offset] = true
	for i := 0; i < m; i++ {
		for j := 0; j <= flips; j++ {
			for d := 0; d < 2; d++ {
				for p := 0; p < maxP; p++ {
					dpCurr[j][d][p] = false
				}
			}
		}
		orig := cmds[i]
		for j := 0; j <= flips; j++ {
			for d := 0; d < 2; d++ {
				for p := 0; p < maxP; p++ {
					if !dpPrev[j][d][p] {
						continue
					}
					for flip := 0; flip < 2; flip++ {
						nj := j + flip
						if nj > flips {
							continue
						}
						cmd := orig
						if flip == 1 {
							if orig == 'T' {
								cmd = 'F'
							} else {
								cmd = 'T'
							}
						}
						nd, np := d, p
						if cmd == 'T' {
							nd = 1 - d
						} else {
							if d == 0 {
								np = p + 1
							} else {
								np = p - 1
							}
						}
						if np >= 0 && np < maxP {
							dpCurr[nj][nd][np] = true
						}
					}
				}
			}
		}
		dpPrev, dpCurr = dpCurr, dpPrev
	}
	ans := 0
	for j := 0; j <= flips; j++ {
		if (flips-j)%2 != 0 {
			continue
		}
		for d := 0; d < 2; d++ {
			for p := 0; p < maxP; p++ {
				if dpPrev[j][d][p] {
					pos := p - offset
					if pos < 0 {
						pos = -pos
					}
					if pos > ans {
						ans = pos
					}
				}
			}
		}
	}
	return ans
}

func run(bin, input string) (string, error) {
	cmd := exec.Command(bin)
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("runtime error: %v\n%s", err, out.String())
	}
	return strings.TrimSpace(out.String()), nil
}

func genCase(rng *rand.Rand) (string, int, string) {
	m := rng.Intn(15) + 1
	b := make([]byte, m)
	for i := range b {
		if rng.Intn(2) == 0 {
			b[i] = 'T'
		} else {
			b[i] = 'F'
		}
	}
	flips := rng.Intn(m + 1)
	input := fmt.Sprintf("%s\n%d\n", string(b), flips)
	return string(b), flips, input
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierC.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < 100; i++ {
		cmds, flips, input := genCase(rng)
		exp := expectedDistance(cmds, flips)
		out, err := run(bin, input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:\n%s", i+1, err, input)
			os.Exit(1)
		}
		got := strings.TrimSpace(out)
		if got != fmt.Sprint(exp) {
			fmt.Fprintf(os.Stderr, "case %d mismatch: expected %d got %s\ninput:\n%s", i+1, exp, got, input)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
