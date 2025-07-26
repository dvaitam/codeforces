package main

import (
	"bufio"
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"strings"
	"time"
)

type caseE struct {
	n     int
	L     int64
	funcs [][]int64
}

func genFuncs(n int, L int64, rng *rand.Rand) [][]int64 {
	funcs := make([][]int64, n)
	for i := 0; i < n; i++ {
		vals := make([]int64, L+1)
		remain := L
		for x := int64(1); x <= L; x++ {
			stepsLeft := L - (x - 1)
			inc := int64(rng.Intn(2))
			if remain == stepsLeft {
				inc = 1
			} else if remain == 0 {
				inc = 0
			}
			if inc > remain {
				inc = remain
			}
			vals[x] = vals[x-1] + inc
			remain -= inc
		}
		vals[L] = L
		funcs[i] = vals
	}
	return funcs
}

func generateCase(rng *rand.Rand) caseE {
	n := rng.Intn(3) + 2
	per := rng.Intn(5) + 1
	L := int64(n * per)
	funcs := genFuncs(n, int64(L), rng)
	return caseE{n: n, L: L, funcs: funcs}
}

func value(f []int64, x int64, L int64) int64 {
	if x < 0 {
		return 0
	}
	if x >= L {
		return L
	}
	return f[x]
}

func runCase(bin string, c caseE) error {
	var cmd *exec.Cmd
	if strings.HasSuffix(bin, ".go") {
		cmd = exec.Command("go", "run", bin)
	} else {
		cmd = exec.Command(bin)
	}
	stdin, err := cmd.StdinPipe()
	if err != nil {
		return err
	}
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		return err
	}
	var stderr bytes.Buffer
	cmd.Stderr = &stderr
	if err := cmd.Start(); err != nil {
		return err
	}
	inw := bufio.NewWriter(stdin)
	outr := bufio.NewReader(stdout)
	fmt.Fprintf(inw, "%d %d\n", c.n, c.L)
	inw.Flush()
	queries := 0
	segs := make([][2]int64, c.n)
	for {
		line, err := outr.ReadString('\n')
		if err != nil {
			cmd.Process.Kill()
			return fmt.Errorf("read error: %v %s", err, stderr.String())
		}
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}
		if strings.HasPrefix(line, "?") {
			var idx int
			var x int64
			fmt.Sscanf(line, "? %d %d", &idx, &x)
			if idx < 1 || idx > c.n {
				cmd.Process.Kill()
				return fmt.Errorf("bad query index")
			}
			queries++
			if queries > 200000 {
				cmd.Process.Kill()
				return fmt.Errorf("too many queries")
			}
			val := value(c.funcs[idx-1], x, c.L)
			fmt.Fprintln(inw, val)
			inw.Flush()
		} else if line == "!" {
			for i := 0; i < c.n; i++ {
				l, _ := outr.ReadString('\n')
				l = strings.TrimSpace(l)
				fmt.Sscanf(l, "%d %d", &segs[i][0], &segs[i][1])
			}
			inw.Flush()
			stdin.Close()
			err := cmd.Wait()
			if err != nil {
				return fmt.Errorf("runtime error: %v %s", err, stderr.String())
			}
			break
		} else {
			cmd.Process.Kill()
			return fmt.Errorf("unexpected output: %q", line)
		}
	}
	per := c.L / int64(c.n)
	for i := 0; i < c.n; i++ {
		l := segs[i][0]
		r := segs[i][1]
		if l < 0 || r < l {
			return fmt.Errorf("bad segment")
		}
		if value(c.funcs[i], r, c.L)-value(c.funcs[i], l, c.L) < per {
			return fmt.Errorf("segment %d too small", i+1)
		}
	}
	for i := 0; i < c.n; i++ {
		for j := i + 1; j < c.n; j++ {
			l1, r1 := segs[i][0], segs[i][1]
			l2, r2 := segs[j][0], segs[j][1]
			if max64(l1, l2) < min64(r1, r2) {
				return fmt.Errorf("segments overlap")
			}
		}
	}
	return nil
}

func max64(a, b int64) int64 {
	if a > b {
		return a
	}
	return b
}
func min64(a, b int64) int64 {
	if a < b {
		return a
	}
	return b
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierE.go /path/to/binary")
		return
	}
	bin := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < 100; i++ {
		c := generateCase(rng)
		if err := runCase(bin, c); err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\n", i+1, err)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
