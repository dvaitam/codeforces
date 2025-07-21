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

func runCandidate(bin, input string) (string, error) {
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
		return "", fmt.Errorf("runtime error: %v\n%s", err, stderr.String())
	}
	return strings.TrimSpace(out.String()), nil
}

type point struct{ x, y int }

type op struct {
	kind string
	x, y int
}

func solveCase(ops []op) []string {
	pts := make(map[point]struct{})
	var res []string
	for _, o := range ops {
		switch o.kind {
		case "add":
			pts[point{o.x, o.y}] = struct{}{}
		case "remove":
			delete(pts, point{o.x, o.y})
		case "find":
			found := false
			best := point{0, 0}
			for p := range pts {
				if p.x > o.x && p.y > o.y {
					if !found || p.x < best.x || (p.x == best.x && p.y < best.y) {
						found = true
						best = p
					}
				}
			}
			if found {
				res = append(res, fmt.Sprintf("%d %d", best.x, best.y))
			} else {
				res = append(res, "-1")
			}
		}
	}
	return res
}

func generateCase(rng *rand.Rand) (string, string) {
	n := rng.Intn(15) + 1
	var ops []op
	pts := make([]point, 0)
	inUse := make(map[point]bool)
	for i := 0; i < n; i++ {
		t := rng.Intn(3)
		switch t {
		case 0:
			// add
			var p point
			for {
				p = point{rng.Intn(10), rng.Intn(10)}
				if !inUse[p] {
					break
				}
			}
			inUse[p] = true
			pts = append(pts, p)
			ops = append(ops, op{"add", p.x, p.y})
		case 1:
			if len(pts) == 0 {
				i--
				continue
			}
			idx := rng.Intn(len(pts))
			p := pts[idx]
			pts = append(pts[:idx], pts[idx+1:]...)
			delete(inUse, p)
			ops = append(ops, op{"remove", p.x, p.y})
		case 2:
			ops = append(ops, op{"find", rng.Intn(10), rng.Intn(10)})
		}
	}
	var sb strings.Builder
	fmt.Fprintf(&sb, "%d\n", len(ops))
	for _, o := range ops {
		fmt.Fprintf(&sb, "%s %d %d\n", o.kind, o.x, o.y)
	}
	outVals := solveCase(ops)
	var out strings.Builder
	for i, v := range outVals {
		if i > 0 {
			out.WriteByte('\n')
		}
		out.WriteString(v)
	}
	if len(outVals) > 0 {
		out.WriteByte('\n')
	}
	return sb.String(), out.String()
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierD.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < 100; i++ {
		in, exp := generateCase(rng)
		out, err := runCandidate(bin, in)
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:\n%s", i+1, err, in)
			os.Exit(1)
		}
		if strings.TrimSpace(out) != strings.TrimSpace(exp) {
			fmt.Fprintf(os.Stderr, "case %d failed: expected\n%s\ngot\n%s\ninput:\n%s", i+1, exp, out, in)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
