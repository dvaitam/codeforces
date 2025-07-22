package main

import (
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
	"time"
)

type event struct{ t, l, r int }

func solveB(n, m, s, f int, ev []event) string {
	dir := 1
	if s > f {
		dir = -1
	}
	cur := s
	idx := 0
	tstep := 1
	res := make([]byte, 0, abs(f-s)+m+5)
	for cur != f {
		if idx < m && ev[idx].t == tstep {
			L, R := ev[idx].l, ev[idx].r
			next := cur + dir
			if (cur >= L && cur <= R) || (next >= L && next <= R) {
				res = append(res, 'X')
			} else if dir == 1 {
				cur++
				res = append(res, 'R')
			} else {
				cur--
				res = append(res, 'L')
			}
			idx++
		} else {
			if dir == 1 {
				cur++
				res = append(res, 'R')
			} else {
				cur--
				res = append(res, 'L')
			}
		}
		tstep++
	}
	return string(res)
}

func abs(x int) int {
	if x < 0 {
		return -x
	}
	return x
}

func runBinary(path, stringInput string) (string, error) {
	cmd := exec.Command(path)
	cmd.Stdin = strings.NewReader(stringInput)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	err := cmd.Run()
	return out.String(), err
}

func generateCase(rng *rand.Rand) string {
	n := rng.Intn(10) + 2
	s := rng.Intn(n) + 1
	f := rng.Intn(n) + 1
	for f == s {
		f = rng.Intn(n) + 1
	}
	m := rng.Intn(10)
	ev := make([]event, m)
	curT := 1
	for i := 0; i < m; i++ {
		curT += rng.Intn(2) + 1
		l := rng.Intn(n) + 1
		r := l + rng.Intn(n-l+1)
		ev[i] = event{curT, l, r}
	}
	var b strings.Builder
	fmt.Fprintf(&b, "%d %d %d %d\n", n, m, s, f)
	for i := 0; i < m; i++ {
		fmt.Fprintf(&b, "%d %d %d\n", ev[i].t, ev[i].l, ev[i].r)
	}
	return b.String()
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierB.go /path/to/binary")
		os.Exit(1)
	}
	cand := os.Args[1]
	_, filename, _, _ := runtime.Caller(0)
	dir := filepath.Dir(filename)
	ref := filepath.Join(dir, "refB")
	cmd := exec.Command("go", "build", "-o", ref, filepath.Join(dir, "342B.go"))
	if out, err := cmd.CombinedOutput(); err != nil {
		fmt.Fprintf(os.Stderr, "failed to build reference: %v\n%s", err, out)
		os.Exit(1)
	}
	defer os.Remove(ref)
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for t := 0; t < 100; t++ {
		input := generateCase(rng)
		candOut, cErr := runBinary(cand, input)
		if cErr != nil {
			fmt.Printf("test %d: candidate error: %v\n", t+1, cErr)
			os.Exit(1)
		}
		evLines := strings.Split(strings.TrimSpace(input), "\n")
		var n, m, s, f int
		fmt.Sscanf(evLines[0], "%d %d %d %d", &n, &m, &s, &f)
		events := make([]event, m)
		for i := 0; i < m; i++ {
			fmt.Sscanf(evLines[i+1], "%d %d %d", &events[i].t, &events[i].l, &events[i].r)
		}
		expect := solveB(n, m, s, f, events)
		if strings.TrimSpace(candOut) != strings.TrimSpace(expect) {
			fmt.Printf("test %d failed\ninput:\n%sexpected:%s\nactual:%s\n", t+1, input, expect, candOut)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
