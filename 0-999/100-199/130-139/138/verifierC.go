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

type Tree struct {
	a, h int
	l, r int
}

type Mush struct {
	b, z int
}

func expected(trees []Tree, mush []Mush) float64 {
	n := len(trees)
	m := len(mush)
	var dfs func(idx int, prob float64, mask int)
	total := 0.0
	dfs = func(idx int, prob float64, mask int) {
		if idx == n {
			sum := 0
			for j := 0; j < m; j++ {
				if mask&(1<<j) != 0 {
					sum += mush[j].z
				}
			}
			total += prob * float64(sum)
			return
		}
		t := trees[idx]
		// stand
		pStand := float64(100-t.l-t.r) / 100.0
		if pStand > 0 {
			dfs(idx+1, prob*pStand, mask)
		}
		// left
		if t.l > 0 {
			newMask := mask
			for j := 0; j < m; j++ {
				if newMask&(1<<j) == 0 {
					continue
				}
				if mush[j].b >= t.a-t.h && mush[j].b < t.a {
					newMask &^= 1 << j
				}
			}
			dfs(idx+1, prob*float64(t.l)/100.0, newMask)
		}
		// right
		if t.r > 0 {
			newMask := mask
			for j := 0; j < m; j++ {
				if newMask&(1<<j) == 0 {
					continue
				}
				if mush[j].b > t.a && mush[j].b <= t.a+t.h {
					newMask &^= 1 << j
				}
			}
			dfs(idx+1, prob*float64(t.r)/100.0, newMask)
		}
	}
	dfs(0, 1.0, (1<<m)-1)
	return total
}

func generateCase(rng *rand.Rand) (string, float64) {
	n := rng.Intn(3) + 1
	m := rng.Intn(3) + 1
	trees := make([]Tree, n)
	mush := make([]Mush, m)
	var sb strings.Builder
	fmt.Fprintf(&sb, "%d %d\n", n, m)
	for i := 0; i < n; i++ {
		trees[i] = Tree{a: rng.Intn(11) - 5, h: rng.Intn(5) + 1, l: rng.Intn(101)}
		trees[i].r = rng.Intn(101 - trees[i].l)
		fmt.Fprintf(&sb, "%d %d %d %d\n", trees[i].a, trees[i].h, trees[i].l, trees[i].r)
	}
	for j := 0; j < m; j++ {
		mush[j] = Mush{b: rng.Intn(11) - 5, z: rng.Intn(10) + 1}
		fmt.Fprintf(&sb, "%d %d\n", mush[j].b, mush[j].z)
	}
	exp := expected(trees, mush)
	return sb.String(), exp
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
	var errBuf bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &errBuf
	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("runtime error: %v\n%s", err, errBuf.String())
	}
	return strings.TrimSpace(out.String()), nil
}

func parseFloat(s string) (float64, error) {
	var v float64
	_, err := fmt.Sscan(s, &v)
	return v, err
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierC.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for t := 1; t <= 100; t++ {
		in, exp := generateCase(rng)
		out, err := run(bin, in)
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:\n%s", t, err, in)
			os.Exit(1)
		}
		got, err := parseFloat(out)
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: invalid float output %q\n", t, out)
			os.Exit(1)
		}
		if diff := got - exp; diff < -1e-4 || diff > 1e-4 {
			fmt.Fprintf(os.Stderr, "case %d failed: expected %.6f got %.6f\ninput:\n%s", t, exp, got, in)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
