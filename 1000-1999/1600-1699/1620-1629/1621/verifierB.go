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

const INF int64 = 1 << 60

type segment struct {
	l int64
	r int64
	c int64
}

func solveCaseB(segs []segment) []string {
	lmin := int64(INF)
	rmax := int64(-INF)
	costL := int64(INF)
	costR := int64(INF)
	best := int64(INF)
	res := make([]string, len(segs))
	for i, s := range segs {
		if s.l < lmin {
			lmin = s.l
			costL = s.c
			best = INF
		} else if s.l == lmin && s.c < costL {
			costL = s.c
		}
		if s.r > rmax {
			rmax = s.r
			costR = s.c
			best = INF
		} else if s.r == rmax && s.c < costR {
			costR = s.c
		}
		if s.l == lmin && s.r == rmax && s.c < best {
			best = s.c
		}
		ans := costL + costR
		if best < ans {
			ans = best
		}
		res[i] = fmt.Sprintf("%d", ans)
	}
	return res
}

func runCandidate(bin, input string) (string, error) {
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

func generateCaseB(rng *rand.Rand) (string, []string) {
	t := rng.Intn(3) + 1
	input := fmt.Sprintf("%d\n", t)
	var outLines []string
	for ; t > 0; t-- {
		n := rng.Intn(6) + 1
		segs := make([]segment, n)
		input += fmt.Sprintf("%d\n", n)
		for i := 0; i < n; i++ {
			l := int64(rng.Intn(10) + 1)
			r := int64(rng.Intn(10-int(l)+1)) + l
			c := int64(rng.Intn(20) + 1)
			segs[i] = segment{l, r, c}
			input += fmt.Sprintf("%d %d %d\n", l, r, c)
		}
		out := solveCaseB(segs)
		outLines = append(outLines, out...)
	}
	return input, outLines
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierB.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < 100; i++ {
		in, exp := generateCaseB(rng)
		got, err := runCandidate(bin, in)
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:\n%s", i+1, err, in)
			os.Exit(1)
		}
        // Accept both line-separated and space-separated outputs by tokenizing
        tokens := strings.Fields(strings.TrimSpace(got))
        if len(tokens) != len(exp) {
            fmt.Fprintf(os.Stderr, "case %d failed: expected %d outputs got %d\ninput:\n%s", i+1, len(exp), len(tokens), in)
            os.Exit(1)
        }
        for j := range exp {
            if strings.TrimSpace(tokens[j]) != exp[j] {
                fmt.Fprintf(os.Stderr, "case %d failed: expected %s got %s\ninput:\n%s", i+1, exp[j], tokens[j], in)
                os.Exit(1)
            }
        }
	}
	fmt.Println("All tests passed")
}
