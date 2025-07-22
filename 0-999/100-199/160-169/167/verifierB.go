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

func solve(n, l, k int, p []int, c []int) float64 {
	var b, b_ [201]float64
	b1 := &b
	b2 := &b_
	b1[0] = 1.0
	kb := 0
	for i := 0; i < n; i++ {
		if c[i] == -1 {
			p1 := float64(p[i]) / 100.0
			p0 := 1 - p1
			for j := 0; j <= kb+1; j++ {
				b2[j] = 0
			}
			for j := 0; j <= kb; j++ {
				b2[j] += p0 * b1[j]
				b2[j+1] += p1 * b1[j]
			}
			b1, b2 = b2, b1
			kb++
		}
	}
	ma := kb - k
	if ma < 0 {
		ma = 0
	}
	var a, a_ [201][201]float64
	a1 := &a
	a2 := &a_
	a1[0][0] = 1.0
	ka := 0
	for i := 0; i < n; i++ {
		if c[i] != -1 {
			p1 := float64(p[i]) / 100.0
			p0 := 1 - p1
			for x := 0; x <= ka+1; x++ {
				for y := 0; y <= ma; y++ {
					a2[x][y] = 0
				}
			}
			for x := 0; x <= ka; x++ {
				for y := 0; y <= ma; y++ {
					a2[x][y] += p0 * a1[x][y]
					idx := y + c[i]
					if idx > ma {
						idx = ma
					}
					a2[x+1][idx] += p1 * a1[x][y]
				}
			}
			a1, a2 = a2, a1
			ka++
		}
	}
	ans := 0.0
	for i := 0; i <= ka; i++ {
		for j := 0; j <= ma; j++ {
			for t := 0; t <= kb; t++ {
				if t+i >= l && j+k >= t {
					ans += a1[i][j] * b1[t]
				}
			}
		}
	}
	return ans
}

func generateCase(rng *rand.Rand) (string, float64) {
	n := rng.Intn(5) + 1
	l := rng.Intn(n + 1)
	k := rng.Intn(n + 1)
	p := make([]int, n)
	c := make([]int, n)
	for i := 0; i < n; i++ {
		p[i] = rng.Intn(101)
		if rng.Intn(2) == 0 {
			c[i] = -1
		} else {
			c[i] = rng.Intn(3)
		}
	}
	ans := solve(n, l, k, p, c)
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("%d %d %d\n", n, l, k))
	for i := 0; i < n; i++ {
		if i > 0 {
			sb.WriteByte(' ')
		}
		sb.WriteString(fmt.Sprintf("%d", p[i]))
	}
	sb.WriteByte('\n')
	for i := 0; i < n; i++ {
		if i > 0 {
			sb.WriteByte(' ')
		}
		sb.WriteString(fmt.Sprintf("%d", c[i]))
	}
	sb.WriteByte('\n')
	return sb.String(), ans
}

func run(bin, input string) (float64, error) {
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
		return 0, fmt.Errorf("runtime error: %v\n%s", err, stderr.String())
	}
	val, err := strconv.ParseFloat(strings.TrimSpace(out.String()), 64)
	if err != nil {
		return 0, fmt.Errorf("invalid float output: %v", err)
	}
	return val, nil
}

func main() {
	if len(os.Args) != 2 && !(len(os.Args) == 3 && os.Args[1] == "--") {
		fmt.Println("usage: go run verifierB.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[len(os.Args)-1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < 100; i++ {
		in, exp := generateCase(rng)
		out, err := run(bin, in)
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:\n%s", i+1, err, in)
			os.Exit(1)
		}
		if mathAbs(out-exp) > 1e-6 {
			fmt.Fprintf(os.Stderr, "case %d mismatch: expected %.8f got %.8f\ninput:\n%s", i+1, exp, out, in)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}

func mathAbs(f float64) float64 {
	if f < 0 {
		return -f
	}
	return f
}
