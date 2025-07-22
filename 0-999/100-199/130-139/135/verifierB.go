package main

import (
	"bytes"
	"fmt"
	"math"
	"math/rand"
	"os"
	"os/exec"
	"strings"
	"time"
)

type Point struct {
	x, y float64
}

func (p Point) Sub(q Point) Point   { return Point{p.x - q.x, p.y - q.y} }
func (p Point) Dot(q Point) float64 { return p.x*q.x + p.y*q.y }
func eq(v float64) bool             { return math.Abs(v) < 1e-7 }

func isSquare(a, b, c, d Point) bool {
	pts := []Point{b, c, d}
	idxs := [][3]int{{0, 1, 2}, {0, 2, 1}, {1, 0, 2}, {1, 2, 0}, {2, 0, 1}, {2, 1, 0}}
	for _, id := range idxs {
		ab := pts[id[0]].Sub(a)
		ac := pts[id[1]].Sub(a)
		ad := pts[id[2]].Sub(a)
		if eq(ab.Dot(ac)) && eq(ab.Dot(ad)) && eq(ac.Dot(pts[id[2]].Sub(pts[id[1]]))) {
			return true
		}
	}
	return false
}

func isRect(a, b, c, d Point) bool {
	pts := []Point{b, c, d}
	idxs := [][3]int{{0, 1, 2}, {0, 2, 1}, {1, 0, 2}, {1, 2, 0}, {2, 0, 1}, {2, 1, 0}}
	for _, id := range idxs {
		ab := pts[id[0]].Sub(a)
		ac := pts[id[1]].Sub(a)
		ad := pts[id[2]].Sub(a)
		if eq(ab.Dot(ac)) && eq(ab.Dot(ad)) && eq(ad.Sub(ac).Dot(pts[id[2]].Sub(pts[id[1]]))) {
			return true
		}
	}
	return false
}

func existsPartition(P [8]Point) bool {
	for i := 0; i < 5; i++ {
		for j := i + 1; j < 6; j++ {
			for k := j + 1; k < 7; k++ {
				for l := k + 1; l < 8; l++ {
					if isSquare(P[i], P[j], P[k], P[l]) {
						rem := make([]int, 0, 4)
						for m := 0; m < 8; m++ {
							if m != i && m != j && m != k && m != l {
								rem = append(rem, m)
							}
						}
						if isRect(P[rem[0]], P[rem[1]], P[rem[2]], P[rem[3]]) {
							return true
						}
					}
				}
			}
		}
	}
	return false
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
	var stderr bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &stderr
	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("runtime error: %v\n%s", err, stderr.String())
	}
	return strings.TrimSpace(out.String()), nil
}

func parseCandidate(out string) (bool, [4]int, [4]int, error) {
	fields := strings.Fields(out)
	if len(fields) == 1 && fields[0] == "NO" {
		return false, [4]int{}, [4]int{}, nil
	}
	if len(fields) != 9 || fields[0] != "YES" {
		return false, [4]int{}, [4]int{}, fmt.Errorf("bad output")
	}
	var sq [4]int
	var rc [4]int
	for i := 0; i < 4; i++ {
		fmt.Sscan(fields[1+i], &sq[i])
	}
	for i := 0; i < 4; i++ {
		fmt.Sscan(fields[5+i], &rc[i])
	}
	return true, sq, rc, nil
}

func checkSets(P [8]Point, sq, rc [4]int) bool {
	used := make([]bool, 8)
	for _, v := range sq {
		if v < 1 || v > 8 || used[v-1] {
			return false
		}
		used[v-1] = true
	}
	for _, v := range rc {
		if v < 1 || v > 8 || used[v-1] {
			return false
		}
		used[v-1] = true
	}
	var sqPts, rcPts [4]Point
	for i, idx := range sq {
		sqPts[i] = P[idx-1]
	}
	for i, idx := range rc {
		rcPts[i] = P[idx-1]
	}
	return isSquare(sqPts[0], sqPts[1], sqPts[2], sqPts[3]) && isRect(rcPts[0], rcPts[1], rcPts[2], rcPts[3])
}

func generateCase(rng *rand.Rand) (string, bool, [8]Point) {
	var P [8]Point
	for i := 0; i < 8; i++ {
		P[i] = Point{float64(rng.Intn(11) - 5), float64(rng.Intn(11) - 5)}
	}
	expect := existsPartition(P)
	var sb strings.Builder
	for i := 0; i < 8; i++ {
		sb.WriteString(fmt.Sprintf("%d %d\n", int(P[i].x), int(P[i].y)))
	}
	return sb.String(), expect, P
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierB.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < 100; i++ {
		in, exp, pts := generateCase(rng)
		out, err := runCandidate(bin, in)
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:\n%s", i+1, err, in)
			os.Exit(1)
		}
		has, sq, rc, err := parseCandidate(out)
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:\n%s", i+1, err, in)
			os.Exit(1)
		}
		if has != exp {
			fmt.Fprintf(os.Stderr, "case %d failed: expected %v got %v\ninput:\n%s", i+1, exp, has, in)
			os.Exit(1)
		}
		if has {
			if !checkSets(pts, sq, rc) {
				fmt.Fprintf(os.Stderr, "case %d failed: candidate produced invalid partition\ninput:\n%s", i+1, in)
				os.Exit(1)
			}
		}
	}
	fmt.Println("All tests passed")
}
