package main

import (
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"sort"
	"strconv"
	"strings"
	"time"
)

type testCase struct {
	bots [4][2]int
}

func genCase(rng *rand.Rand) testCase {
	side := rng.Intn(5) + 1
	x0 := rng.Intn(5)
	y0 := rng.Intn(5)
	// square corners
	corners := [4][2]int{{x0, y0}, {x0 + side, y0}, {x0, y0 + side}, {x0 + side, y0 + side}}
	// shuffle order and add small noise
	perm := rng.Perm(4)
	var bots [4][2]int
	for i := 0; i < 4; i++ {
		bots[i][0] = corners[perm[i]][0] + rng.Intn(3) - 1
		bots[i][1] = corners[perm[i]][1] + rng.Intn(3) - 1
	}
	return testCase{bots: bots}
}

var perms [][4]int
var cornerPositions = [4][2]int{{0, 0}, {1, 0}, {0, 1}, {1, 1}}

func init() {
	used := make([]bool, 4)
	var cur [4]int
	var dfs func(int)
	dfs = func(pos int) {
		if pos == 4 {
			tmp := cur
			perms = append(perms, tmp)
			return
		}
		for i := 0; i < 4; i++ {
			if !used[i] {
				used[i] = true
				cur[pos] = i
				dfs(pos + 1)
				used[i] = false
			}
		}
	}
	dfs(0)
}

func abs(x int) int {
	if x < 0 {
		return -x
	}
	return x
}

func solve(tc testCase) string {
	A := tc.bots
	INF := int(1e9)
	Dset := make(map[int]struct{})
	for i := 0; i < 4; i++ {
		for j := i + 1; j < 4; j++ {
			Dset[abs(A[i][0]-A[j][0])] = struct{}{}
			Dset[abs(A[i][1]-A[j][1])] = struct{}{}
		}
	}
	Ds := make([]int, 0, len(Dset))
	for k := range Dset {
		Ds = append(Ds, k)
	}
	Ds = uniqueInts(Ds)
	ans := INF
	var final [4][2]int
	for _, d := range Ds {
		X := make([]int, 0, 40)
		Y := make([]int, 0, 40)
		for j := 0; j < 4; j++ {
			X = append(X, A[j][0], A[j][0]-d, A[j][0]+d)
			Y = append(Y, A[j][1], A[j][1]-d, A[j][1]+d)
		}
		for _, p := range perms {
			lx, rx := INF, -INF
			ly, ry := INF, -INF
			for k := 0; k < 4; k++ {
				xx := A[k][0] - cornerPositions[p[k]][0]*d
				yy := A[k][1] - cornerPositions[p[k]][1]*d
				if xx < lx {
					lx = xx
				}
				if xx > rx {
					rx = xx
				}
				if yy < ly {
					ly = yy
				}
				if yy > ry {
					ry = yy
				}
			}
			X = append(X, (lx+rx)/2)
			Y = append(Y, (ly+ry)/2)
		}
		X = uniqueInts(X)
		Y = uniqueInts(Y)
		for _, x := range X {
			for _, y := range Y {
				for _, p := range perms {
					tmax := 0
					ok := true
					var pos [4][2]int
					for j := 0; j < 4; j++ {
						xx := x + cornerPositions[p[j]][0]*d
						yy := y + cornerPositions[p[j]][1]*d
						if xx != A[j][0] && yy != A[j][1] {
							ok = false
							break
						}
						move := abs(xx-A[j][0]) + abs(yy-A[j][1])
						if move > tmax {
							tmax = move
						}
						pos[j][0] = xx
						pos[j][1] = yy
					}
					if ok && tmax < ans {
						ans = tmax
						final = pos
					}
				}
			}
		}
	}
	if ans == INF {
		return "-1"
	}
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("%d\n", ans))
	for i := 0; i < 4; i++ {
		sb.WriteString(fmt.Sprintf("%d %d\n", final[i][0], final[i][1]))
	}
	return strings.TrimSpace(sb.String())
}

func uniqueInts(a []int) []int {
	sort.Ints(a)
	j := 0
	for i := 0; i < len(a); i++ {
		if i == 0 || a[i] != a[i-1] {
			a[j] = a[i]
			j++
		}
	}
	return a[:j]
}

func buildInput(tc testCase) string {
	var sb strings.Builder
	sb.WriteString("1\n")
	for i := 0; i < 4; i++ {
		sb.WriteString(fmt.Sprintf("%d %d\n", tc.bots[i][0], tc.bots[i][1]))
	}
	return sb.String()
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
	var stderr bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &stderr
	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("runtime error: %v\n%s", err, stderr.String())
	}
	return strings.TrimSpace(out.String()), nil
}

func parseOutput(s string) (int, [4][2]int, bool) {
	fields := strings.Fields(s)
	if len(fields) == 0 {
		return 0, [4][2]int{}, false
	}
	if fields[0] == "-1" {
		return -1, [4][2]int{}, true
	}
	cost, err := strconv.Atoi(fields[0])
	if err != nil || len(fields) != 9 {
		return 0, [4][2]int{}, false
	}
	var pos [4][2]int
	for i := 0; i < 4; i++ {
		x, err1 := strconv.Atoi(fields[1+2*i])
		y, err2 := strconv.Atoi(fields[2+2*i])
		if err1 != nil || err2 != nil {
			return 0, [4][2]int{}, false
		}
		pos[i][0] = x
		pos[i][1] = y
	}
	return cost, pos, true
}

func validateOutput(tc testCase, cost int, pos [4][2]int) error {
	if cost < 0 {
		return fmt.Errorf("negative cost")
	}
	// Check that the 4 positions form an axis-aligned square with side > 0
	// Collect all x and y coords
	xs := map[int]bool{pos[0][0]: true, pos[1][0]: true, pos[2][0]: true, pos[3][0]: true}
	ys := map[int]bool{pos[0][1]: true, pos[1][1]: true, pos[2][1]: true, pos[3][1]: true}
	if len(xs) != 2 || len(ys) != 2 {
		return fmt.Errorf("positions do not form a square: need exactly 2 distinct x and 2 distinct y")
	}
	xvals := make([]int, 0, 2)
	for v := range xs {
		xvals = append(xvals, v)
	}
	yvals := make([]int, 0, 2)
	for v := range ys {
		yvals = append(yvals, v)
	}
	if xvals[0] > xvals[1] {
		xvals[0], xvals[1] = xvals[1], xvals[0]
	}
	if yvals[0] > yvals[1] {
		yvals[0], yvals[1] = yvals[1], yvals[0]
	}
	side := xvals[1] - xvals[0]
	if side <= 0 {
		return fmt.Errorf("side length must be positive")
	}
	if yvals[1]-yvals[0] != side {
		return fmt.Errorf("not a square: x-side=%d y-side=%d", side, yvals[1]-yvals[0])
	}
	// Each of the 4 corners must be present
	corners := map[[2]int]bool{
		{xvals[0], yvals[0]}: false,
		{xvals[0], yvals[1]}: false,
		{xvals[1], yvals[0]}: false,
		{xvals[1], yvals[1]}: false,
	}
	for i := 0; i < 4; i++ {
		key := [2]int{pos[i][0], pos[i][1]}
		if _, ok := corners[key]; !ok {
			return fmt.Errorf("position %d (%d,%d) is not a corner of the square", i, pos[i][0], pos[i][1])
		}
		corners[key] = true
	}
	for k, v := range corners {
		if !v {
			return fmt.Errorf("corner (%d,%d) not assigned to any bot", k[0], k[1])
		}
	}
	// Check movement costs: each bot must move via L1 along one axis (share x or y with original)
	maxMove := 0
	for i := 0; i < 4; i++ {
		ox, oy := tc.bots[i][0], tc.bots[i][1]
		nx, ny := pos[i][0], pos[i][1]
		if ox != nx && oy != ny {
			return fmt.Errorf("bot %d moves diagonally from (%d,%d) to (%d,%d): must share x or y", i, ox, oy, nx, ny)
		}
		move := abs(ox-nx) + abs(oy-ny)
		if move > maxMove {
			maxMove = move
		}
	}
	if maxMove != cost {
		return fmt.Errorf("claimed cost %d but actual max move is %d", cost, maxMove)
	}
	return nil
}

func main() {
	if len(os.Args) == 3 && os.Args[1] == "--" {
		os.Args = append([]string{os.Args[0]}, os.Args[2])
	}
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierD.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < 100; i++ {
		tc := genCase(rng)
		input := buildInput(tc)
		expStr := solve(tc)
		got, err := run(bin, input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\n", i+1, err)
			os.Exit(1)
		}
		// Parse expected to get optimal cost
		expCost, _, expOk := parseOutput(expStr)
		if !expOk {
			fmt.Fprintf(os.Stderr, "case %d: cannot parse reference output: %s\n", i+1, expStr)
			os.Exit(1)
		}
		// Parse candidate output
		gotCost, gotPos, gotOk := parseOutput(strings.TrimSpace(got))
		if !gotOk {
			fmt.Fprintf(os.Stderr, "case %d failed: cannot parse output\ninput:\n%s\ngot:\n%s\n", i+1, input, got)
			os.Exit(1)
		}
		if expCost == -1 {
			if gotCost != -1 {
				fmt.Fprintf(os.Stderr, "case %d failed: expected -1 got cost %d\ninput:\n%s\n", i+1, gotCost, input)
				os.Exit(1)
			}
			continue
		}
		if gotCost == -1 {
			fmt.Fprintf(os.Stderr, "case %d failed: candidate says -1 but solution exists with cost %d\ninput:\n%s\n", i+1, expCost, input)
			os.Exit(1)
		}
		if gotCost != expCost {
			fmt.Fprintf(os.Stderr, "case %d failed: optimal cost is %d but candidate claims %d\ninput:\n%s\ngot:\n%s\n", i+1, expCost, gotCost, input, got)
			os.Exit(1)
		}
		if err := validateOutput(tc, gotCost, gotPos); err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: invalid output: %v\ninput:\n%s\ngot:\n%s\n", i+1, err, input, got)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
