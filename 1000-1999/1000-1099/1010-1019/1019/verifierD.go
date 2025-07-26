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

func abs(x int64) int64 {
	if x < 0 {
		return -x
	}
	return x
}

func area2(a, b, c point) int64 { return abs((b.x-a.x)*(c.y-a.y) - (b.y-a.y)*(c.x-a.x)) }

func brute(points []point, p2 int64) (bool, [3]point) {
	n := len(points)
	var ans [3]point
	for i := 0; i < n; i++ {
		for j := i + 1; j < n; j++ {
			for k := j + 1; k < n; k++ {
				if area2(points[i], points[j], points[k]) == p2 {
					ans = [3]point{points[i], points[j], points[k]}
					return true, ans
				}
			}
		}
	}
	return false, ans
}

func checkTriple(points []point, p2 int64, coords []point) bool {
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
	return area2(coords[0], coords[1], coords[2]) == p2
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
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierD.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for tc := 0; tc < 100; tc++ {
		n := rng.Intn(5) + 3
		p2 := int64(rng.Intn(30) + 1)
		pts := make([]point, n)
		for i := 0; i < n; i++ {
			pts[i] = point{int64(rng.Intn(20)), int64(rng.Intn(20))}
		}
		var sb strings.Builder
		fmt.Fprintf(&sb, "%d %d\n", n, p2/2)
		for _, pt := range pts {
			fmt.Fprintf(&sb, "%d %d\n", pt.x, pt.y)
		}
		input := sb.String()
		exist, _ := brute(pts, p2)
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
				fmt.Fprintf(os.Stderr, "case %d failed: expected Yes\ninput:\n%s", tc+1, input)
				os.Exit(1)
			}
			coords := make([]point, 0, 3)
			for i := 0; i < 3; i++ {
				var x, y int64
				if !scan.Scan() {
					fmt.Fprintf(os.Stderr, "case %d failed: missing coord\n", tc+1)
					os.Exit(1)
				}
				xVal, err := strconv.ParseInt(scan.Text(), 10, 64)
				if err != nil {
					fmt.Fprintf(os.Stderr, "case %d bad output\n", tc+1)
					os.Exit(1)
				}
				if !scan.Scan() {
					fmt.Fprintf(os.Stderr, "case %d failed: missing coord\n", tc+1)
					os.Exit(1)
				}
				yVal, err2 := strconv.ParseInt(scan.Text(), 10, 64)
				if err2 != nil {
					fmt.Fprintf(os.Stderr, "case %d bad output\n", tc+1)
					os.Exit(1)
				}
				x = xVal
				y = yVal
				coords = append(coords, point{x, y})
			}
			if !checkTriple(pts, p2, coords) {
				fmt.Fprintf(os.Stderr, "case %d failed: invalid triangle\ninput:\n%s\noutput:\n%s", tc+1, input, out)
				os.Exit(1)
			}
		} else {
			if tok != "no" {
				fmt.Fprintf(os.Stderr, "case %d failed: expected No\ninput:\n%s", tc+1, input)
				os.Exit(1)
			}
		}
	}
	fmt.Println("All tests passed")
}
