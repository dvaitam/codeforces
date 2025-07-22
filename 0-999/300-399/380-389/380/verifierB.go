package main

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"math/rand"
	"os"
	"os/exec"
	"strings"
	"time"
)

type Op struct {
	level, l, r int
}

func solveB(r io.Reader) string {
	reader := bufio.NewReader(r)
	var n, m int
	if _, err := fmt.Fscan(reader, &n, &m); err != nil {
		return ""
	}
	opsByX := make(map[int][]Op)
	var sb strings.Builder
	for i := 0; i < m; i++ {
		var tp int
		fmt.Fscan(reader, &tp)
		if tp == 1 {
			var t, l, r, x int
			fmt.Fscan(reader, &t, &l, &r, &x)
			opsByX[x] = append(opsByX[x], Op{level: t, l: l, r: r})
		} else {
			var t0, v0 int
			fmt.Fscan(reader, &t0, &v0)
			cnt := 0
			for _, ops := range opsByX {
				for _, op := range ops {
					if op.level < t0 {
						continue
					}
					maxPos := v0 + (op.level - t0)
					if op.l <= maxPos && op.r >= v0 {
						cnt++
						break
					}
				}
			}
			fmt.Fprintln(&sb, cnt)
		}
	}
	return sb.String()
}

func runCaseB(bin string, in, expect string) error {
	cmd := exec.Command(bin)
	cmd.Stdin = strings.NewReader(in)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("runtime error: %v\n%s", err, out.String())
	}
	if strings.TrimSpace(out.String()) != strings.TrimSpace(expect) {
		return fmt.Errorf("expected %q got %q", strings.TrimSpace(expect), strings.TrimSpace(out.String()))
	}
	return nil
}

func genCaseB(rng *rand.Rand) string {
	n := rng.Intn(4) + 1
	m := rng.Intn(8) + 1
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("%d %d\n", n, m))
	for i := 0; i < m; i++ {
		if rng.Intn(2) == 0 {
			t := rng.Intn(n) + 1
			l := rng.Intn(3) + 1
			r := l + rng.Intn(3)
			x := rng.Intn(10) + 1
			sb.WriteString(fmt.Sprintf("1 %d %d %d %d\n", t, l, r, x))
		} else {
			t := rng.Intn(n) + 1
			v := rng.Intn(3) + 1
			sb.WriteString(fmt.Sprintf("2 %d %d\n", t, v))
		}
	}
	return sb.String()
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierB.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < 100; i++ {
		in := genCaseB(rng)
		expect := solveB(strings.NewReader(in))
		if err := runCaseB(bin, in, expect); err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:\n%s", i+1, err, in)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
