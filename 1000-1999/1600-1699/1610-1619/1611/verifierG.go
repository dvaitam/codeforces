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

// Embedded solver (same logic as the accepted solution cf_t24_1611_G.go)

type Point struct {
	u, v int
}

func minPaths(pts []Point) int {
	if len(pts) == 0 {
		return 0
	}
	sort.Slice(pts, func(i, j int) bool {
		if pts[i].u != pts[j].u {
			return pts[i].u < pts[j].u
		}
		return pts[i].v < pts[j].v
	})
	tails := make([]int, 0)
	for _, p := range pts {
		x := -p.v
		l, r := 0, len(tails)
		for l < r {
			m := l + (r-l)/2
			if tails[m] < x {
				l = m + 1
			} else {
				r = m
			}
		}
		if l == len(tails) {
			tails = append(tails, x)
		} else {
			tails[l] = x
		}
	}
	return len(tails)
}

func oracleSolve(input string) (string, error) {
	sc := strings.Fields(input)
	pos := 0
	nextInt := func() int {
		v, _ := strconv.Atoi(sc[pos])
		pos++
		return v
	}
	t := nextInt()
	var results []string
	for tc := 0; tc < t; tc++ {
		n := nextInt()
		m := nextInt()
		var evenPts, oddPts []Point
		for i := 1; i <= n; i++ {
			row := sc[pos]
			pos++
			for j := 1; j <= m; j++ {
				if row[j-1] == '1' {
					u := i + j
					v := i - j
					if u%2 == 0 {
						evenPts = append(evenPts, Point{u, v})
					} else {
						oddPts = append(oddPts, Point{u, v})
					}
				}
			}
		}
		ans := minPaths(evenPts) + minPaths(oddPts)
		results = append(results, strconv.Itoa(ans))
	}
	return strings.Join(results, "\n"), nil
}

func run(bin, input string) (string, error) {
	cmd := exec.Command(bin)
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	var errBuf bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &errBuf
	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("runtime error: %v\n%s", err, errBuf.String())
	}
	return strings.TrimSpace(out.String()), nil
}

func genCase(r *rand.Rand) string {
	t := r.Intn(3) + 1
	var sb strings.Builder
	fmt.Fprintf(&sb, "%d\n", t)
	for i := 0; i < t; i++ {
		n := r.Intn(5) + 1
		m := r.Intn(5) + 2
		if n*m > 50 {
			m = 50/n + 1
		}
		fmt.Fprintf(&sb, "%d %d\n", n, m)
		for x := 0; x < n; x++ {
			for y := 0; y < m; y++ {
				if r.Intn(2) == 0 {
					sb.WriteByte('0')
				} else {
					sb.WriteByte('1')
				}
			}
			sb.WriteByte('\n')
		}
	}
	return sb.String()
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierG.go /path/to/binary")
		os.Exit(1)
	}
	userBin := os.Args[1]

	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < 100; i++ {
		input := genCase(r)
		want, err := oracleSolve(input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "oracle failed on test %d: %v\ninput:\n%s", i+1, err, input)
			os.Exit(1)
		}
		got, err := run(userBin, input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "test %d: %v\ninput:\n%s", i+1, err, input)
			os.Exit(1)
		}
		if want != got {
			fmt.Printf("test %d failed\ninput:\n%sexpected: %s\ngot: %s\n", i+1, input, want, got)
			os.Exit(1)
		}
	}
	fmt.Printf("All %d tests passed\n", 100)
}
