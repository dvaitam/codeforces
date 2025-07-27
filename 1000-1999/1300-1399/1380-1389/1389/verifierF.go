package main

import (
	"bytes"
	"fmt"
	"math/bits"
	"math/rand"
	"os"
	"os/exec"
	"strings"
	"time"
)

type Segment struct {
	l, r int
	t    int
}

func intersect(a, b Segment) bool {
	if a.t == b.t {
		return false
	}
	return !(a.r < b.l || b.r < a.l)
}

func solveCase(segs []Segment) int {
	n := len(segs)
	best := 0
	for mask := 1; mask < (1 << n); mask++ {
		ok := true
		for i := 0; i < n && ok; i++ {
			if mask&(1<<i) == 0 {
				continue
			}
			for j := i + 1; j < n; j++ {
				if mask&(1<<j) == 0 {
					continue
				}
				if intersect(segs[i], segs[j]) {
					ok = false
					break
				}
			}
		}
		if ok {
			cnt := bits.OnesCount(uint(mask))
			if cnt > best {
				best = cnt
			}
		}
	}
	return best
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

func genTest(rng *rand.Rand) (string, string) {
	n := rng.Intn(6) + 1
	segs := make([]Segment, n)
	for i := 0; i < n; i++ {
		l := rng.Intn(10)
		r := l + rng.Intn(10)
		t := rng.Intn(2) + 1
		segs[i] = Segment{l, r, t}
	}
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("1\n%d\n", n))
	for i := 0; i < n; i++ {
		sb.WriteString(fmt.Sprintf("%d %d %d\n", segs[i].l, segs[i].r, segs[i].t))
	}
	exp := fmt.Sprintf("%d", solveCase(segs))
	return sb.String(), exp
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierF.go /path/to/binary")
		return
	}
	bin := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 1; i <= 100; i++ {
		input, exp := genTest(rng)
		out, err := run(bin, input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "test %d: %v\n", i, err)
			os.Exit(1)
		}
		if strings.TrimSpace(out) != strings.TrimSpace(exp) {
			fmt.Fprintf(os.Stderr, "test %d failed\ninput:\n%sexpected:%s\ngot:%s\n", i, input, exp, out)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
