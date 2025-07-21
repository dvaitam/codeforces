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

type Segment struct {
	speed int
	time  int
}

func countCrosses(a, b []Segment) int64 {
	i, j := 0, 0
	remA, remB := a[0].time, b[0].time
	vA, vB := a[0].speed, b[0].speed
	posA, posB := 0, 0
	var crosses int64
	for i < len(a) && j < len(b) {
		dur := remA
		if remB < dur {
			dur = remB
		}
		delta0 := posA - posB
		dv := vA - vB
		if dv > 0 {
			if delta0 < 0 && delta0+dv*dur > 0 {
				crosses++
			}
		} else if dv < 0 {
			if delta0 > 0 && delta0+dv*dur < 0 {
				crosses++
			}
		}
		posA += vA * dur
		posB += vB * dur
		remA -= dur
		remB -= dur
		if remA == 0 {
			i++
			if i < len(a) {
				remA = a[i].time
				vA = a[i].speed
			}
		}
		if remB == 0 {
			j++
			if j < len(b) {
				remB = b[j].time
				vB = b[j].speed
			}
		}
	}
	return crosses
}

func expected(cars [][]Segment) int64 {
	var total int64
	for i := 0; i < len(cars); i++ {
		for j := i + 1; j < len(cars); j++ {
			total += countCrosses(cars[i], cars[j])
		}
	}
	return total
}

func randomSegments(rng *rand.Rand, s int) []Segment {
	k := rng.Intn(3) + 1
	segs := make([]Segment, k)
	remaining := s
	for i := 0; i < k-1; i++ {
		v := rng.Intn(10) + 1
		maxLen := remaining - (k - i - 1)
		if maxLen < 1 {
			maxLen = 1
		}
		length := rng.Intn(maxLen) + 1
		t := length / v
		if t == 0 {
			t = 1
			length = v
		} else {
			length = v * t
		}
		if length > remaining-(k-i-1) {
			length = remaining - (k - i - 1)
			t = length / v
		}
		segs[i] = Segment{v, t}
		remaining -= length
	}
	segs[k-1] = Segment{1, remaining}
	return segs
}

func generateCase(rng *rand.Rand) (string, int64) {
	n := rng.Intn(4) + 2
	s := rng.Intn(40) + 10
	cars := make([][]Segment, n)
	for i := 0; i < n; i++ {
		cars[i] = randomSegments(rng, s)
	}
	var sb strings.Builder
	fmt.Fprintf(&sb, "%d %d\n", n, s)
	for i := 0; i < n; i++ {
		fmt.Fprintf(&sb, "%d", len(cars[i]))
		for _, seg := range cars[i] {
			fmt.Fprintf(&sb, " %d %d", seg.speed, seg.time)
		}
		sb.WriteByte('\n')
	}
	return sb.String(), expected(cars)
}

func runCase(exe, input string, exp int64) error {
	cmd := exec.Command(exe)
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("runtime error: %v\n%s", err, out.String())
	}
	var got int64
	if _, err := fmt.Fscan(strings.NewReader(out.String()), &got); err != nil {
		return fmt.Errorf("bad output: %v", err)
	}
	if got != exp {
		return fmt.Errorf("expected %d got %d", exp, got)
	}
	return nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierE.go /path/to/binary")
		os.Exit(1)
	}
	exe := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < 100; i++ {
		input, exp := generateCase(rng)
		if err := runCase(exe, input, exp); err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:\n%s", i+1, err, input)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
