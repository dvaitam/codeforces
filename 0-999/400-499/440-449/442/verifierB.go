package main

import (
	"bufio"
	"bytes"
	"fmt"
	"math"
	"math/rand"
	"os"
	"os/exec"
	"sort"
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
	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("runtime error: %v\n%s", err, out.String())
	}
	return out.String(), nil
}

func solveB(r *bufio.Reader) string {
	var n int
	fmt.Fscan(r, &n)
	p := make([]float64, n)
	for i := 0; i < n; i++ {
		fmt.Fscan(r, &p[i])
	}
	sort.Float64s(p)
	const eps = 1e-9
	if n > 0 && math.Abs(p[n-1]-1.0) < eps {
		return fmt.Sprintf("%.15f\n", 1.0)
	}
	prod := make([]float64, n+1)
	rev := make([]float64, n+1)
	prod[0] = 1.0
	for i := 0; i < n; i++ {
		prod[i+1] = prod[i] * (1.0 - p[i])
	}
	rev[n] = 1.0
	for i := n - 1; i >= 0; i-- {
		rev[i] = rev[i+1] * (1.0 - p[i])
	}
	best := 0.0
	for i := 0; i <= n; i++ {
		for j := n; j >= i; j-- {
			pr := prod[i] * rev[j]
			cnd := 0.0
			for k := 0; k < i; k++ {
				if 1.0-p[k] > eps {
					cnd += p[k] * pr / (1.0 - p[k])
				}
			}
			for k := n - 1; k >= j; k-- {
				if 1.0-p[k] > eps {
					cnd += p[k] * pr / (1.0 - p[k])
				}
			}
			if cnd > best {
				best = cnd
			}
		}
	}
	return fmt.Sprintf("%.15f\n", best)
}

func generateCase(rng *rand.Rand) string {
	n := rng.Intn(6) + 1
	var b strings.Builder
	fmt.Fprintf(&b, "%d\n", n)
	for i := 0; i < n; i++ {
		if i > 0 {
			b.WriteByte(' ')
		}
		val := rng.Float64()
		fmt.Fprintf(&b, "%.6f", val)
	}
	b.WriteByte('\n')
	return b.String()
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierB.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < 100; i++ {
		tc := generateCase(rng)
		expect := solveB(bufio.NewReader(strings.NewReader(tc)))
		out, err := runBinary(bin, tc)
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:\n%s", i+1, err, tc)
			os.Exit(1)
		}
		got, err := strconv.ParseFloat(strings.TrimSpace(out), 64)
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: cannot parse output\n", i+1)
			os.Exit(1)
		}
		exp, _ := strconv.ParseFloat(strings.TrimSpace(expect), 64)
		if math.Abs(got-exp) > 1e-6 {
			fmt.Fprintf(os.Stderr, "case %d failed: expected %.6f got %.6f\ninput:\n%s", i+1, exp, got, tc)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
