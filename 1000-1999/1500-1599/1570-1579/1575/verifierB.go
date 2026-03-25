package main

import (
	"context"
	"fmt"
	"math"
	"math/rand"
	"os"
	"os/exec"
	"sort"
	"strconv"
	"strings"
	"time"
)

func runBinary(bin, input string) (string, error) {
	if !strings.Contains(bin, "/") {
		bin = "./" + bin
	}
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()
	cmd := exec.CommandContext(ctx, bin)
	cmd.Stdin = strings.NewReader(input)
	out, err := cmd.CombinedOutput()
	if ctx.Err() == context.DeadlineExceeded {
		return "", fmt.Errorf("timeout")
	}
	if err != nil {
		return "", fmt.Errorf("%v: %s", err, out)
	}
	return strings.TrimSpace(string(out)), nil
}

type Event struct {
	pos float64
	val int
}

type EventSlice []Event

func (s EventSlice) Len() int { return len(s) }
func (s EventSlice) Less(i, j int) bool {
	if s[i].pos != s[j].pos {
		return s[i].pos < s[j].pos
	}
	return s[i].val > s[j].val
}
func (s EventSlice) Swap(i, j int) { s[i], s[j] = s[j], s[i] }

type Point struct {
	d     float64
	alpha float64
}

func solveB(n, k int, pts []Point) float64 {
	if k <= 0 {
		return 0.0
	}

	events := make(EventSlice, 0, 4*len(pts))

	check := func(r float64) bool {
		events = events[:0]
		for _, p := range pts {
			if p.d > 2*r {
				continue
			}
			ratio := p.d / (2 * r)
			if ratio > 1.0 {
				ratio = 1.0
			}
			delta := math.Acos(ratio) + 1e-11
			l := p.alpha - delta
			ri := p.alpha + delta

			events = append(events, Event{l, 1}, Event{ri, -1}, Event{l + 2*math.Pi, 1}, Event{ri + 2*math.Pi, -1})
		}

		sort.Sort(events)

		cnt := 0
		maxCnt := 0
		for _, e := range events {
			cnt += e.val
			if cnt > maxCnt {
				maxCnt = cnt
			}
		}

		return maxCnt >= k
	}

	low, high := 0.0, 3e5
	for iter := 0; iter < 55; iter++ {
		mid := (low + high) / 2
		if check(mid) {
			high = mid
		} else {
			low = mid
		}
	}

	return high
}

func genTest() (string, int, int, []Point) {
	n := rand.Intn(6) + 1
	k := rand.Intn(n) + 1
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("%d %d\n", n, k))

	originCount := 0
	var pts []Point
	for i := 0; i < n; i++ {
		x := rand.Intn(11) - 5
		y := rand.Intn(11) - 5
		sb.WriteString(fmt.Sprintf("%d %d\n", x, y))
		if x == 0 && y == 0 {
			originCount++
		} else {
			pts = append(pts, Point{
				d:     math.Hypot(float64(x), float64(y)),
				alpha: math.Atan2(float64(y), float64(x)),
			})
		}
	}
	return sb.String(), k - originCount, originCount, pts
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintf(os.Stderr, "usage: go run verifierB.go /path/to/binary\n")
		os.Exit(1)
	}
	cand := os.Args[1]
	rand.Seed(1)
	for i := 0; i < 100; i++ {
		input, effK, _, pts := genTest()

		expected := solveB(len(pts), effK, pts)

		got, err := runBinary(cand, input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "test %d failed to run: %v\n", i+1, err)
			os.Exit(1)
		}

		gotVal, err := strconv.ParseFloat(strings.TrimSpace(got), 64)
		if err != nil {
			fmt.Fprintf(os.Stderr, "test %d: could not parse output %q: %v\n", i+1, got, err)
			os.Exit(1)
		}

		// Compare with relative/absolute tolerance
		diff := math.Abs(gotVal - expected)
		tol := 1e-4
		if diff > tol && diff > tol*math.Abs(expected) {
			fmt.Fprintf(os.Stderr, "test %d failed:\ninput:\n%sexpected: %.10f\ngot: %.10f\n", i+1, input, expected, gotVal)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
