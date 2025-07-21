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

const NEG_INF = -1000000000

func euphony(s string, k int, c [26][26]int) int {
	m := len(s)
	dpPrev := make([][]int, k+1)
	dpCurr := make([][]int, k+1)
	for j := 0; j <= k; j++ {
		dpPrev[j] = make([]int, 26)
		dpCurr[j] = make([]int, 26)
		for l := 0; l < 26; l++ {
			dpPrev[j][l] = NEG_INF
			dpCurr[j][l] = NEG_INF
		}
	}
	orig0 := int(s[0] - 'a')
	for l := 0; l < 26; l++ {
		chg := 0
		if l != orig0 {
			chg = 1
		}
		if chg <= k {
			dpPrev[chg][l] = 0
		}
	}
	for pos := 1; pos < m; pos++ {
		for j := 0; j <= k; j++ {
			for l := 0; l < 26; l++ {
				dpCurr[j][l] = NEG_INF
			}
		}
		curOrig := int(s[pos] - 'a')
		for used := 0; used <= k; used++ {
			for prevL := 0; prevL < 26; prevL++ {
				prevVal := dpPrev[used][prevL]
				if prevVal <= NEG_INF {
					continue
				}
				for l := 0; l < 26; l++ {
					nj := used
					if l != curOrig {
						nj++
					}
					if nj > k {
						continue
					}
					val := prevVal + c[prevL][l]
					if val > dpCurr[nj][l] {
						dpCurr[nj][l] = val
					}
				}
			}
		}
		dpPrev, dpCurr = dpCurr, dpPrev
	}
	ans := NEG_INF
	for used := 0; used <= k; used++ {
		for l := 0; l < 26; l++ {
			if dpPrev[used][l] > ans {
				ans = dpPrev[used][l]
			}
		}
	}
	if m <= 1 {
		ans = 0
	}
	return ans
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

func randString(rng *rand.Rand, n int) string {
	b := make([]byte, n)
	for i := 0; i < n; i++ {
		b[i] = byte('a' + rng.Intn(4))
	}
	return string(b)
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierC.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < 100; i++ {
		s := randString(rng, rng.Intn(5)+1)
		k := rng.Intn(len(s) + 1)
		n := rng.Intn(5)
		var c [26][26]int
		for j := 0; j < n; j++ {
			x := rng.Intn(4)
			y := rng.Intn(4)
			val := rng.Intn(11) - 5
			c[x][y] = val
		}
		expected := fmt.Sprintf("%d", euphony(s, k, c))
		var sb strings.Builder
		sb.WriteString(fmt.Sprintf("%s %d\n", s, k))
		sb.WriteString(fmt.Sprintf("%d\n", n))
		for x := 0; x < 26; x++ {
			for y := 0; y < 26; y++ {
				if c[x][y] != 0 {
					sb.WriteString(fmt.Sprintf("%c %c %d\n", 'a'+x, 'a'+y, c[x][y]))
				}
			}
		}
		input := sb.String()
		got, err := run(bin, input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:\n%s", i+1, err, input)
			os.Exit(1)
		}
		if strings.TrimSpace(got) != expected {
			fmt.Fprintf(os.Stderr, "case %d failed: expected %s got %s\ninput:\n%s", i+1, expected, got, input)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
