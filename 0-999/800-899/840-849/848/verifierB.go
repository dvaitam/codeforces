package main

import (
	"bufio"
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"sort"
	"strings"
	"time"
)

type dancer struct {
	idx int
	g   int
	p   int
	t   int
}

func runSolution(bin, input string) (string, error) {
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

func solveB(n, w, h int, arr []dancer) [][2]int {
	groups := make(map[int][]dancer)
	for _, d := range arr {
		key := d.p - d.t
		groups[key] = append(groups[key], d)
	}
	ans := make([][2]int, n)
	for _, group := range groups {
		var horiz, vert []dancer
		for _, d := range group {
			if d.g == 2 {
				horiz = append(horiz, d)
			} else {
				vert = append(vert, d)
			}
		}
		sort.Slice(horiz, func(i, j int) bool { return horiz[i].p < horiz[j].p })
		sort.Slice(vert, func(i, j int) bool { return vert[i].p < vert[j].p })
		start := append(append([]dancer{}, horiz...), vert...)
		var fin [][2]int
		tempV := append([]dancer{}, vert...)
		sort.Slice(tempV, func(i, j int) bool { return tempV[i].p < tempV[j].p })
		for _, d := range tempV {
			fin = append(fin, [2]int{d.p, h})
		}
		tempH := append([]dancer{}, horiz...)
		sort.Slice(tempH, func(i, j int) bool { return tempH[i].p < tempH[j].p })
		for _, d := range tempH {
			fin = append(fin, [2]int{w, d.p})
		}
		for i, d := range start {
			ans[d.idx] = fin[i]
		}
	}
	return ans
}

func parseInputB(in string) (int, int, int, []dancer, error) {
	r := bufio.NewReader(strings.NewReader(in))
	var n, w, h int
	if _, err := fmt.Fscan(r, &n, &w, &h); err != nil {
		return 0, 0, 0, nil, err
	}
	arr := make([]dancer, n)
	for i := 0; i < n; i++ {
		fmt.Fscan(r, &arr[i].g, &arr[i].p, &arr[i].t)
		arr[i].idx = i
	}
	return n, w, h, arr, nil
}

func verifyB(input, output string) error {
	n, w, h, arr, err := parseInputB(input)
	if err != nil {
		return fmt.Errorf("input parse: %v", err)
	}
	expected := solveB(n, w, h, arr)
	r := bufio.NewReader(strings.NewReader(output))
	for i := 0; i < n; i++ {
		var x, y int
		if _, err := fmt.Fscan(r, &x, &y); err != nil {
			return fmt.Errorf("parse output line %d: %v", i+1, err)
		}
		if x != expected[i][0] || y != expected[i][1] {
			return fmt.Errorf("dancer %d expected %d %d got %d %d", i+1, expected[i][0], expected[i][1], x, y)
		}
	}
	return nil
}

func generateCaseB(rng *rand.Rand) string {
	n := rng.Intn(5) + 1
	w := rng.Intn(9) + 2
	h := rng.Intn(9) + 2
	var sb strings.Builder
	fmt.Fprintf(&sb, "%d %d %d\n", n, w, h)
	used := make(map[string]bool)
	for i := 0; i < n; i++ {
		for {
			g := rng.Intn(2) + 1
			var p int
			if g == 1 {
				p = rng.Intn(w-1) + 1
			} else {
				p = rng.Intn(h-1) + 1
			}
			t := rng.Intn(10)
			key := fmt.Sprintf("%d-%d-%d", g, p, t)
			if !used[key] {
				used[key] = true
				fmt.Fprintf(&sb, "%d %d %d\n", g, p, t)
				break
			}
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
	cases := []string{}
	for len(cases) < 100 {
		cases = append(cases, generateCaseB(rng))
	}
	for i, tc := range cases {
		out, err := runSolution(bin, tc)
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:\n%s", i+1, err, tc)
			os.Exit(1)
		}
		if err := verifyB(tc, out); err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:\n%soutput:\n%s", i+1, err, tc, out)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
