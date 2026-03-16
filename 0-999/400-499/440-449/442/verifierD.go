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

// Embedded correct solver for 442D
func solveD(input string) string {
	r := strings.NewReader(input)
	var n int
	fmt.Fscan(r, &n)

	parent := make([]int, n+2)
	h := make([]uint8, n+2)
	maxChild := make([]uint8, n+2)
	countMax := make([]uint32, n+2)

	out := make([]byte, 0, n*4)

	for i := 1; i <= n; i++ {
		var p int
		fmt.Fscan(r, &p)
		v := i + 1
		parent[v] = p

		oldH := h[p]
		changed := false

		if countMax[p] == 0 {
			countMax[p] = 1
			changed = true
		} else if maxChild[p] == 0 {
			countMax[p]++
			changed = true
		}

		if changed {
			newH := maxChild[p]
			if countMax[p] >= 2 {
				newH++
			}
			if newH != oldH {
				h[p] = newH
				child := p
				oldVal := oldH
				newVal := newH

				for {
					par := parent[child]
					if par == 0 {
						break
					}

					oldParH := h[par]
					mp := maxChild[par]

					if oldVal == mp {
						maxChild[par] = newVal
						countMax[par] = 1
					} else if newVal == mp {
						countMax[par]++
					} else {
						break
					}

					newParH := maxChild[par]
					if countMax[par] >= 2 {
						newParH++
					}
					if newParH == oldParH {
						break
					}

					h[par] = newParH
					child = par
					oldVal = oldParH
					newVal = newParH
				}
			}
		}

		ans := 0
		if countMax[1] != 0 {
			ans = int(maxChild[1]) + 1
		}

		if i > 1 {
			out = append(out, ' ')
		}
		out = appendInt442(out, ans)
	}

	return string(out)
}

func appendInt442(out []byte, x int) []byte {
	if x >= 100 {
		out = append(out, byte('0'+x/100))
		x %= 100
		out = append(out, byte('0'+x/10), byte('0'+x%10))
	} else if x >= 10 {
		out = append(out, byte('0'+x/10), byte('0'+x%10))
	} else {
		out = append(out, byte('0'+x))
	}
	return out
}

func runBinary(path, input string) (string, error) {
	cmd := exec.Command(path)
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("runtime error: %v\n%s", err, out.String())
	}
	return out.String(), nil
}

func generateCase(rng *rand.Rand) string {
	n := rng.Intn(20) + 1
	var b strings.Builder
	fmt.Fprintf(&b, "%d\n", n)
	for i := 1; i <= n; i++ {
		if i > 1 {
			b.WriteByte(' ')
		}
		fmt.Fprintf(&b, "%d", rng.Intn(i)+1)
	}
	b.WriteByte('\n')
	return b.String()
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierD.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]

	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < 100; i++ {
		tc := generateCase(rng)
		expect := solveD(tc)
		out, err := runBinary(bin, tc)
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:\n%s", i+1, err, tc)
			os.Exit(1)
		}
		if strings.TrimSpace(out) != strings.TrimSpace(expect) {
			fmt.Fprintf(os.Stderr, "case %d failed: expected %s got %sinput:\n%s", i+1, expect, out, tc)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
