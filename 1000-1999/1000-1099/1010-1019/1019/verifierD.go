package main

import (
	"bufio"
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"time"
)

type point struct{ x, y int64 }

func abs64(x int64) int64 {
	if x < 0 {
		return -x
	}
	return x
}

func area2(a, b, c point) int64 {
	return abs64((b.x-a.x)*(c.y-a.y) - (b.y-a.y)*(c.x-a.x))
}

func brute(points []point, twoS int64) (bool, [3]point) {
	n := len(points)
	var ans [3]point
	for i := 0; i < n; i++ {
		for j := i + 1; j < n; j++ {
			for k := j + 1; k < n; k++ {
				if area2(points[i], points[j], points[k]) == twoS {
					ans = [3]point{points[i], points[j], points[k]}
					return true, ans
				}
			}
		}
	}
	return false, ans
}

func checkTriple(points []point, twoS int64, coords []point) bool {
	used := make([]bool, len(points))
	for _, pt := range coords {
		found := false
		for i, p := range points {
			if !used[i] && p == pt {
				used[i] = true
				found = true
				break
			}
		}
		if !found {
			return false
		}
	}
	return area2(coords[0], coords[1], coords[2]) == twoS
}

func run(bin, input string) (string, error) {
	cmd := exec.Command(bin)
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
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierD.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for tc := 0; tc < 100; tc++ {
		n := rng.Intn(5) + 3
		S := int64(rng.Intn(15) + 1)
		twoS := 2 * S
		pts := make([]point, n)
		for i := 0; i < n; i++ {
			pts[i] = point{int64(rng.Intn(20)), int64(rng.Intn(20))}
		}
		var sb strings.Builder
		fmt.Fprintf(&sb, "%d %d\n", n, S)
		for _, pt := range pts {
			fmt.Fprintf(&sb, "%d %d\n", pt.x, pt.y)
		}
		input := sb.String()
		exist, _ := brute(pts, twoS)
		out, err := run(bin, input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:\n%s", tc+1, err, input)
			os.Exit(1)
		}
		scan := bufio.NewScanner(strings.NewReader(out))
		scan.Split(bufio.ScanWords)
		if !scan.Scan() {
			fmt.Fprintf(os.Stderr, "case %d failed: no output\n", tc+1)
			os.Exit(1)
		}
		tok := strings.ToLower(scan.Text())
		if exist {
			if tok != "yes" {
				// The candidate might say No if it can't find a solution with its algorithm,
				// but the brute force found one. However, since the candidate is CF-accepted,
				// if brute says Yes, candidate should too. Unless points are duplicated /
				// collinear (which the problem says won't happen, but our random gen might).
				// Just check: if candidate says Yes, validate the triangle.
				// If candidate says No but brute says Yes, that's a failure.
				fmt.Fprintf(os.Stderr, "case %d failed: expected Yes, got %s\ninput:\n%s", tc+1, tok, input)
				os.Exit(1)
			}
			coords := make([]point, 0, 3)
			for i := 0; i < 3; i++ {
				if !scan.Scan() {
					fmt.Fprintf(os.Stderr, "case %d failed: missing coord\n", tc+1)
					os.Exit(1)
				}
				xVal, err := strconv.ParseInt(scan.Text(), 10, 64)
				if err != nil {
					fmt.Fprintf(os.Stderr, "case %d bad output x\n", tc+1)
					os.Exit(1)
				}
				if !scan.Scan() {
					fmt.Fprintf(os.Stderr, "case %d failed: missing coord\n", tc+1)
					os.Exit(1)
				}
				yVal, err := strconv.ParseInt(scan.Text(), 10, 64)
				if err != nil {
					fmt.Fprintf(os.Stderr, "case %d bad output y\n", tc+1)
					os.Exit(1)
				}
				coords = append(coords, point{xVal, yVal})
			}
			if !checkTriple(pts, twoS, coords) {
				fmt.Fprintf(os.Stderr, "case %d failed: invalid triangle\ninput:\n%s\noutput:\n%s", tc+1, input, out)
				os.Exit(1)
			}
		} else {
			if tok == "yes" {
				// Candidate says Yes but brute says No - validate the triangle anyway
				coords := make([]point, 0, 3)
				valid := true
				for i := 0; i < 3; i++ {
					if !scan.Scan() {
						valid = false
						break
					}
					xVal, err := strconv.ParseInt(scan.Text(), 10, 64)
					if err != nil {
						valid = false
						break
					}
					if !scan.Scan() {
						valid = false
						break
					}
					yVal, err := strconv.ParseInt(scan.Text(), 10, 64)
					if err != nil {
						valid = false
						break
					}
					coords = append(coords, point{xVal, yVal})
				}
				if !valid || !checkTriple(pts, twoS, coords) {
					fmt.Fprintf(os.Stderr, "case %d failed: candidate says Yes but triangle invalid\ninput:\n%s\noutput:\n%s", tc+1, input, out)
					os.Exit(1)
				}
			}
			// If both say No, that's fine
		}
	}
	fmt.Println("All tests passed")
}
