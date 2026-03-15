package main

import (
	"bytes"
	"fmt"
	"math"
	"math/rand"
	"sort"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"time"
)

type point struct {
	x, y int
}

type pair struct {
	x, y int
}

func checkCenter(cx2, cy2 int, pts []point, k int) bool {
	set := make(map[pair]struct{}, len(pts))
	for _, p := range pts {
		set[pair{p.x, p.y}] = struct{}{}
	}
	miss := 0
	for _, p := range pts {
		sym := pair{cx2 - p.x, cy2 - p.y}
		if _, ok := set[sym]; !ok {
			miss++
			if miss > k {
				return false
			}
		}
	}
	return true
}

func solveF(n, k int, pts []point) (int, []pair) {
	if k >= n {
		return -1, nil
	}

	// Sort points by a unique key to determine which points must be paired
	sorted := make([]point, n)
	copy(sorted, pts)
	sort.Slice(sorted, func(i, j int) bool {
		if sorted[i].x != sorted[j].x {
			return sorted[i].x < sorted[j].x
		}
		return sorted[i].y < sorted[j].y
	})

	t := k + 1
	if t > n {
		t = n
	}

	// Any valid center must pair at least one of the first t points with one of the last t points
	candMap := make(map[pair]struct{})
	for i := 0; i < t; i++ {
		for j := n - t; j < n; j++ {
			candMap[pair{sorted[i].x + sorted[j].x, sorted[i].y + sorted[j].y}] = struct{}{}
		}
	}

	var res []pair
	for c := range candMap {
		if checkCenter(c.x, c.y, pts, k) {
			res = append(res, c)
		}
	}
	return len(res), res
}

func parseCandidateOutput(output string) (int, []pair, error) {
	fields := strings.Fields(output)
	if len(fields) == 0 {
		return 0, nil, fmt.Errorf("empty output")
	}
	cnt, err := strconv.Atoi(fields[0])
	if err != nil {
		return 0, nil, fmt.Errorf("invalid count %q", fields[0])
	}
	if cnt == -1 {
		if len(fields) != 1 {
			return -1, nil, fmt.Errorf("expected only -1, got extra tokens")
		}
		return -1, nil, nil
	}
	if len(fields) != 1+2*cnt {
		return 0, nil, fmt.Errorf("expected %d coordinate pairs after count %d, got %d tokens total", cnt, cnt, len(fields))
	}
	var centers []pair
	for i := 0; i < cnt; i++ {
		xf, err1 := strconv.ParseFloat(fields[1+2*i], 64)
		yf, err2 := strconv.ParseFloat(fields[1+2*i+1], 64)
		if err1 != nil || err2 != nil {
			return 0, nil, fmt.Errorf("invalid coordinate at center %d", i+1)
		}
		// Convert to doubled integer coordinates
		cx2 := int(math.Round(xf * 2))
		cy2 := int(math.Round(yf * 2))
		centers = append(centers, pair{cx2, cy2})
	}
	return cnt, centers, nil
}

func main() {
	if len(os.Args) != 2 && !(len(os.Args) == 3 && os.Args[1] == "--") {
		fmt.Fprintln(os.Stderr, "usage: go run verifierF.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[len(os.Args)-1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < 100; i++ {
		n, k, pts, input := generateCaseF(rng)
		refCount, refCenters := solveF(n, k, pts)

		got, err := run(bin, input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:\n%s", i+1, err, input)
			os.Exit(1)
		}

		gotCount, gotCenters, err := parseCandidateOutput(got)
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: invalid output: %v\ninput:\n%soutput:\n%s\n", i+1, err, input, got)
			os.Exit(1)
		}

		// Check count matches
		if gotCount != refCount {
			fmt.Fprintf(os.Stderr, "case %d failed: expected count %d, got %d\ninput:\n%sexpected centers: %v\ngot output:\n%s\n",
				i+1, refCount, gotCount, input, refCenters, got)
			os.Exit(1)
		}

		if refCount == -1 {
			continue
		}

		// Verify each candidate center is valid
		for ci, c := range gotCenters {
			if !checkCenter(c.x, c.y, pts, k) {
				fmt.Fprintf(os.Stderr, "case %d failed: candidate center %d (%.1f, %.1f) is not valid\ninput:\n%s",
					i+1, ci+1, float64(c.x)/2, float64(c.y)/2, input)
				os.Exit(1)
			}
		}

		// Verify no duplicates
		seen := make(map[pair]bool)
		for _, c := range gotCenters {
			if seen[c] {
				fmt.Fprintf(os.Stderr, "case %d failed: duplicate center (%.1f, %.1f)\ninput:\n%s",
					i+1, float64(c.x)/2, float64(c.y)/2, input)
				os.Exit(1)
			}
			seen[c] = true
		}
	}
	fmt.Println("All tests passed")
}

func generateCaseF(rng *rand.Rand) (int, int, []point, string) {
	n := rng.Intn(7) + 2
	k := rng.Intn(3)
	if k >= n {
		k = n - 1
	}
	pts := make([]point, 0, n)
	used := make(map[pair]bool)
	for len(pts) < n {
		x := rng.Intn(21) - 10
		y := rng.Intn(21) - 10
		p := pair{x, y}
		if !used[p] {
			used[p] = true
			pts = append(pts, point{x, y})
		}
	}
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("%d %d\n", n, k))
	for i := 0; i < n; i++ {
		sb.WriteString(fmt.Sprintf("%d %d\n", pts[i].x, pts[i].y))
	}
	return n, k, pts, sb.String()
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
