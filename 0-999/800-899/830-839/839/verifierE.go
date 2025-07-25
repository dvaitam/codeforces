package main

import (
	"bytes"
	"fmt"
	"math/bits"
	"math/rand"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"time"
)

func runBinary(path, input string) (string, error) {
	var cmd *exec.Cmd
	if strings.HasSuffix(path, ".go") {
		cmd = exec.Command("go", "run", path)
	} else {
		cmd = exec.Command(path)
	}
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	err := cmd.Run()
	return strings.TrimSpace(out.String()), err
}

func maxClique(n int, g [][]bool) int {
	best := 1
	for mask := 1; mask < (1 << n); mask++ {
		ok := true
		cnt := bits.OnesCount(uint(mask))
		if cnt <= best {
			continue
		}
		for i := 0; i < n && ok; i++ {
			if mask&(1<<i) == 0 {
				continue
			}
			for j := i + 1; j < n; j++ {
				if mask&(1<<j) != 0 && !g[i][j] {
					ok = false
					break
				}
			}
		}
		if ok && cnt > best {
			best = cnt
		}
	}
	return best
}

func solveCase(n, k int, g [][]bool) float64 {
	m := maxClique(n, g)
	if m == 0 {
		return 0
	}
	return float64(k*k) * float64(m-1) / float64(2*m)
}

func genCase(rng *rand.Rand) (int, int, [][]bool) {
	n := rng.Intn(6) + 1
	k := rng.Intn(100) + 1
	g := make([][]bool, n)
	for i := 0; i < n; i++ {
		g[i] = make([]bool, n)
	}
	for i := 0; i < n; i++ {
		for j := i + 1; j < n; j++ {
			if rng.Intn(2) == 1 {
				g[i][j] = true
				g[j][i] = true
			}
		}
	}
	return n, k, g
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierE.go /path/to/binary")
		return
	}
	bin := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < 100; i++ {
		n, k, g := genCase(rng)
		var sb strings.Builder
		fmt.Fprintf(&sb, "%d %d\n", n, k)
		for x := 0; x < n; x++ {
			for y := 0; y < n; y++ {
				if y > 0 {
					sb.WriteByte(' ')
				}
				if g[x][y] {
					sb.WriteByte('1')
				} else {
					sb.WriteByte('0')
				}
			}
			sb.WriteByte('\n')
		}
		input := sb.String()
		want := solveCase(n, k, g)
		out, err := runBinary(bin, input)
		if err != nil {
			fmt.Printf("case %d: runtime error: %v\ninput:\n%s", i+1, err, input)
			return
		}
		got, err := strconv.ParseFloat(strings.TrimSpace(out), 64)
		if err != nil {
			fmt.Printf("case %d: bad float output %q\n", i+1, out)
			return
		}
		diff := got - want
		if diff < -1e-6 || diff > 1e-6 {
			fmt.Printf("case %d failed\ninput:\n%sexpected: %.6f\ngot: %s\n", i+1, input, want, out)
			return
		}
	}
	fmt.Println("All tests passed")
}
