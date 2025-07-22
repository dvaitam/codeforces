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

func run(bin, input string) (string, error) {
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
		return "", fmt.Errorf("runtime error: %v\n%s", err, out.String())
	}
	return strings.TrimSpace(out.String()), nil
}

func solveD(n, m, k int, grid []string) int64 {
	psa := make([][]int, n+1)
	for i := range psa {
		psa[i] = make([]int, m+1)
	}
	for i := 1; i <= n; i++ {
		for j := 1; j <= m; j++ {
			psa[i][j] = psa[i-1][j] + psa[i][j-1] - psa[i-1][j-1]
			if grid[i-1][j-1] == 'a' {
				psa[i][j]++
			}
		}
	}
	var ans int64
	for x1 := 1; x1 < n; x1++ {
		for x2 := x1 + 1; x2 <= n; x2++ {
			var pos [26][]int
			for y := 1; y <= m; y++ {
				c1 := grid[x1-1][y-1]
				c2 := grid[x2-1][y-1]
				if c1 == c2 {
					pos[c1-'a'] = append(pos[c1-'a'], y)
				}
			}
			for c := 0; c < 26; c++ {
				P := pos[c]
				t := len(P)
				if t < 2 {
					continue
				}
				r := 1
				for l := 0; l < t-1; l++ {
					if r < l+1 {
						r = l + 1
					}
					for r < t {
						y1 := P[l]
						y2 := P[r]
						cnt := psa[x2][y2] - psa[x1-1][y2] - psa[x2][y1-1] + psa[x1-1][y1-1]
						if cnt <= k {
							r++
						} else {
							break
						}
					}
					ans += int64(r - l - 1)
				}
			}
		}
	}
	return ans
}

func randLetter(rng *rand.Rand) byte {
	return byte('a' + rng.Intn(3))
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierD.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for tc := 1; tc <= 100; tc++ {
		n := rng.Intn(4) + 2
		m := rng.Intn(4) + 2
		k := rng.Intn(n*m + 1)
		grid := make([]string, n)
		for i := 0; i < n; i++ {
			b := make([]byte, m)
			for j := range b {
				b[j] = randLetter(rng)
			}
			grid[i] = string(b)
		}
		var sb strings.Builder
		sb.WriteString(fmt.Sprintf("%d %d %d\n", n, m, k))
		for i := 0; i < n; i++ {
			sb.WriteString(grid[i])
			sb.WriteByte('\n')
		}
		input := sb.String()
		expect := fmt.Sprintf("%d", solveD(n, m, k, grid))
		out, err := run(bin, input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:\n%s", tc, err, input)
			os.Exit(1)
		}
		if strings.TrimSpace(out) != expect {
			fmt.Fprintf(os.Stderr, "case %d failed: expected %s got %s\ninput:\n%s", tc, expect, strings.TrimSpace(out), input)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
