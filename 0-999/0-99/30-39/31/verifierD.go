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

type seg struct {
	x1, y1, x2, y2 int
}

type rect struct {
	x1, y1, x2, y2 int
}

type testCaseD struct {
	W, H  int
	segs  []seg
	areas []int
}

func generateCaseD(rng *rand.Rand) testCaseD {
	W := rng.Intn(9) + 2
	H := rng.Intn(9) + 2
	n := rng.Intn(4) + 1
	rects := []rect{{0, 0, W, H}}
	segs := make([]seg, 0, n)
	for i := 0; i < n; i++ {
		idx := rng.Intn(len(rects))
		r := rects[idx]
		var s seg
		if (r.x2-r.x1 > 1 && rng.Intn(2) == 0) || r.y2-r.y1 <= 1 {
			// vertical
			x := rng.Intn(r.x2-r.x1-1) + r.x1 + 1
			s = seg{x, r.y1, x, r.y2}
			rect1 := rect{r.x1, r.y1, x, r.y2}
			rect2 := rect{x, r.y1, r.x2, r.y2}
			rects[idx] = rect1
			rects = append(rects, rect2)
		} else {
			// horizontal
			y := rng.Intn(r.y2-r.y1-1) + r.y1 + 1
			s = seg{r.x1, y, r.x2, y}
			rect1 := rect{r.x1, r.y1, r.x2, y}
			rect2 := rect{r.x1, y, r.x2, r.y2}
			rects[idx] = rect1
			rects = append(rects, rect2)
		}
		segs = append(segs, s)
	}
	areas := make([]int, len(rects))
	for i, r := range rects {
		areas[i] = (r.x2 - r.x1) * (r.y2 - r.y1)
	}
	sort.Ints(areas)
	// shuffle segs
	perm := rng.Perm(len(segs))
	shuffled := make([]seg, len(segs))
	for i, p := range perm {
		shuffled[i] = segs[p]
	}
	segs = shuffled
	return testCaseD{W: W, H: H, segs: segs, areas: areas}
}

func runCaseD(bin string, tc testCaseD) error {
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("%d %d %d\n", tc.W, tc.H, len(tc.segs)))
	for _, s := range tc.segs {
		sb.WriteString(fmt.Sprintf("%d %d %d %d\n", s.x1, s.y1, s.x2, s.y2))
	}
	input := sb.String()

	cmd := exec.Command(bin)
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("runtime error: %v\n%s", err, out.String())
	}
	fields := strings.Fields(out.String())
	if len(fields) != len(tc.areas) {
		return fmt.Errorf("expected %d numbers got %d", len(tc.areas), len(fields))
	}
	got := make([]int, len(fields))
	for i, f := range fields {
		if _, err := fmt.Sscan(f, &got[i]); err != nil {
			return fmt.Errorf("bad output: %v", err)
		}
	}
	sort.Ints(got)
	for i := range got {
		if got[i] != tc.areas[i] {
			return fmt.Errorf("expected %v got %v", tc.areas, got)
		}
	}
	return nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierD.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < 100; i++ {
		tc := generateCaseD(rng)
		if err := runCaseD(bin, tc); err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\n", i+1, err)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
