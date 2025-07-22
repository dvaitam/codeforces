package main

import (
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"sort"
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
	mp := make(map[pair]int)
	for _, p := range pts {
		mp[pair{p.x, p.y}]++
	}
	removed := 0
	for _, p := range pts {
		if mp[pair{p.x, p.y}] == 0 {
			continue
		}
		mp[pair{p.x, p.y}]--
		sym := pair{cx2 - p.x, cy2 - p.y}
		if mp[sym] > 0 {
			mp[sym]--
		} else {
			removed++
			if removed > k {
				return false
			}
		}
	}
	return true
}

func solveF(n, k int, pts []point) string {
	if k >= n {
		return "-1"
	}
	candMap := make(map[pair]struct{})
	for i := 0; i < n; i++ {
		candMap[pair{2 * pts[i].x, 2 * pts[i].y}] = struct{}{}
		for j := i + 1; j < n; j++ {
			candMap[pair{pts[i].x + pts[j].x, pts[i].y + pts[j].y}] = struct{}{}
		}
	}
	res := make([]pair, 0)
	for c := range candMap {
		if checkCenter(c.x, c.y, pts, k) {
			res = append(res, c)
		}
	}
	sort.Slice(res, func(i, j int) bool {
		if res[i].x != res[j].x {
			return res[i].x < res[j].x
		}
		return res[i].y < res[j].y
	})
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("%d\n", len(res)))
	for _, c := range res {
		sb.WriteString(fmt.Sprintf("%.1f %.1f\n", float64(c.x)/2.0, float64(c.y)/2.0))
	}
	return strings.TrimRight(sb.String(), "\n")
}

func generateCaseF(rng *rand.Rand) (string, string) {
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
	input := sb.String()
	expect := solveF(n, k, pts)
	return input, expect
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

func main() {
	if len(os.Args) != 2 && !(len(os.Args) == 3 && os.Args[1] == "--") {
		fmt.Fprintln(os.Stderr, "usage: go run verifierF.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[len(os.Args)-1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < 100; i++ {
		in, expect := generateCaseF(rng)
		got, err := run(bin, in)
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:\n%s", i+1, err, in)
			os.Exit(1)
		}
		if got != expect {
			fmt.Fprintf(os.Stderr, "case %d failed: expected %s got %s\ninput:\n%s", i+1, expect, got, in)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
