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

type fenwick struct {
	f []uint8
	n int
}

func newFenwick(n int) *fenwick {
	return &fenwick{make([]uint8, n+2), n}
}

func (fw *fenwick) add(idx int) {
	idx++
	for idx <= fw.n+1 {
		fw.f[idx] ^= 1
		idx += idx & -idx
	}
}

func (fw *fenwick) sum(idx int) uint8 {
	idx++
	var res uint8
	for idx > 0 {
		res ^= fw.f[idx]
		idx -= idx & -idx
	}
	return res
}

func expectedD(n int, grid []string) string {
	g := make([][]byte, n)
	for i := range g {
		g[i] = []byte(grid[i])
	}
	offset := n - 1
	size := 2*n - 1
	fw := newFenwick(size)
	ans := 0
	for s := 2; s <= 2*n; s++ {
		iLow := 1
		if s-n > iLow {
			iLow = s - n
		}
		iHigh := n
		if s-1 < iHigh {
			iHigh = s - 1
		}
		for i := iLow; i <= iHigh; i++ {
			j := s - i
			idx := (i - j) + offset
			val := g[i-1][j-1] - '0'
			if fw.sum(idx)%2 == 1 {
				val ^= 1
			}
			if val == 1 {
				ans++
				fw.add(idx)
			}
		}
	}
	return fmt.Sprintf("%d\n", ans)
}

func generateCase(rng *rand.Rand) (string, string) {
	n := rng.Intn(3) + 2
	grid := make([]string, n)
	for i := 0; i < n; i++ {
		row := make([]byte, n)
		for j := 0; j < n; j++ {
			if rng.Intn(2) == 0 {
				row[j] = '0'
			} else {
				row[j] = '1'
			}
		}
		grid[i] = string(row)
	}
	var sb strings.Builder
	sb.WriteString("1\n")
	sb.WriteString(fmt.Sprintf("%d\n", n))
	for _, r := range grid {
		sb.WriteString(r)
		sb.WriteByte('\n')
	}
	expect := expectedD(n, grid)
	return sb.String(), expect
}

func runCase(bin, input, expected string) error {
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
		return fmt.Errorf("runtime error: %v\n%s", err, out.String())
	}
	outStr := strings.TrimSpace(out.String())
	expStr := strings.TrimSpace(expected)
	if outStr != expStr {
		return fmt.Errorf("expected %q got %q", expStr, outStr)
	}
	return nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierD.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < 100; i++ {
		in, exp := generateCase(rng)
		if err := runCase(bin, in, exp); err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:\n%s", i+1, err, in)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
