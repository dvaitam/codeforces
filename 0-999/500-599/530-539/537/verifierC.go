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

type note struct{ d, h int }

func abs(x int) int {
	if x < 0 {
		return -x
	}
	return x
}
func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}

func expectedAnswerC(n int, notes []note) (string, int) {
	ans := notes[0].h + (notes[0].d - 1)
	for i := 0; i < len(notes)-1; i++ {
		d1, h1 := notes[i].d, notes[i].h
		d2, h2 := notes[i+1].d, notes[i+1].h
		dist := d2 - d1
		diff := abs(h2 - h1)
		if diff > dist {
			return "IMPOSSIBLE", 0
		}
		cand := max(h1, h2) + (dist-diff)/2
		if cand > ans {
			ans = cand
		}
	}
	last := notes[len(notes)-1]
	if last.h+(n-last.d) > ans {
		ans = last.h + (n - last.d)
	}
	return "OK", ans
}

func generateCaseC(rng *rand.Rand) (int, []note) {
	n := rng.Intn(30) + 1
	m := rng.Intn(n) + 1
	ds := rand.Perm(n)[:m]
	sort.Ints(ds)
	notes := make([]note, m)
	for i, d := range ds {
		notes[i] = note{d + 1, rng.Intn(20)}
	}
	return n, notes
}

func runCaseC(bin string, n int, notes []note) error {
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("%d %d\n", n, len(notes)))
	for _, nt := range notes {
		sb.WriteString(fmt.Sprintf("%d %d\n", nt.d, nt.h))
	}
	cmd := exec.Command(bin)
	cmd.Stdin = strings.NewReader(sb.String())
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("runtime error: %v\n%s", err, out.String())
	}
	got := strings.TrimSpace(out.String())
	status, ans := expectedAnswerC(n, notes)
	var expected string
	if status == "IMPOSSIBLE" {
		expected = "IMPOSSIBLE"
	} else {
		expected = fmt.Sprint(ans)
	}
	if got != expected {
		return fmt.Errorf("expected %s got %s", expected, got)
	}
	return nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierC.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < 100; i++ {
		n, notes := generateCaseC(rng)
		if err := runCaseC(bin, n, notes); err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\n", i+1, err)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
