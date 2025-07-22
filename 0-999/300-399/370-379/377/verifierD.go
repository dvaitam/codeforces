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

type seg struct{ l, x, r int }

type testCaseD struct {
	n    int
	segs []seg
}

func generateCase(rng *rand.Rand) (string, testCaseD) {
	n := rng.Intn(8) + 1
	maxC := 20
	segs := make([]seg, n)
	for i := 0; i < n; i++ {
		l := rng.Intn(maxC-1) + 1
		r := l + rng.Intn(maxC-l)
		x := l + rng.Intn(r-l+1)
		segs[i] = seg{l, x, r}
	}
	var b strings.Builder
	fmt.Fprintf(&b, "%d\n", n)
	for _, s := range segs {
		fmt.Fprintf(&b, "%d %d %d\n", s.l, s.x, s.r)
	}
	return b.String(), testCaseD{n, segs}
}

func brute(tc testCaseD) int {
	maxC := 0
	for _, s := range tc.segs {
		if s.r > maxC {
			maxC = s.r
		}
	}
	best := 0
	for L := 1; L <= maxC; L++ {
		for R := L; R <= maxC; R++ {
			cnt := 0
			for _, s := range tc.segs {
				if s.l <= L && s.x >= L && s.x <= R && s.r >= R {
					cnt++
				}
			}
			if cnt > best {
				best = cnt
			}
		}
	}
	return best
}

func verify(tc testCaseD, out string) error {
	reader := bufio.NewReader(strings.NewReader(out))
	var ans int
	if _, err := fmt.Fscan(reader, &ans); err != nil {
		return fmt.Errorf("parse ans")
	}
	var ids []int
	for {
		var v int
		if _, err := fmt.Fscan(reader, &v); err != nil {
			break
		}
		ids = append(ids, v)
	}
	best := brute(tc)
	if ans != best {
		return fmt.Errorf("expected %d got %d", best, ans)
	}
	if len(ids) != ans {
		return fmt.Errorf("expected %d ids", ans)
	}
	seen := make(map[int]bool)
	Llow := 0
	Lhigh := 1 << 30
	Rlow := 0
	Rhigh := 1 << 30
	for _, id := range ids {
		if id < 1 || id > tc.n {
			return fmt.Errorf("bad id %d", id)
		}
		if seen[id] {
			return fmt.Errorf("duplicate id %d", id)
		}
		seen[id] = true
		s := tc.segs[id-1]
		if s.l > Llow {
			Llow = s.l
		}
		if s.x < Lhigh {
			Lhigh = s.x
		}
		if s.x > Rlow {
			Rlow = s.x
		}
		if s.r < Rhigh {
			Rhigh = s.r
		}
	}
	if Llow > Lhigh || Rlow > Rhigh || Llow > Rhigh {
		return fmt.Errorf("infeasible interval")
	}
	return nil
}

func runCase(bin string, input string, tc testCaseD) error {
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
	return verify(tc, out.String())
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierD.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < 100; i++ {
		input, tc := generateCase(rng)
		if err := runCase(bin, input, tc); err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:\n%s", i+1, err, input)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
