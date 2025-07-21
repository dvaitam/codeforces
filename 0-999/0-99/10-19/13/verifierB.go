package main

import (
	"bufio"
	"bytes"
	"fmt"
	"math"
	"math/rand"
	"os"
	"os/exec"
	"strings"
	"time"
)

type Point struct{ x, y int }

func run(bin, input string) (string, error) {
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
	return out.String(), nil
}

func onSegment(a, b, c Point) bool {
	cross := int64(b.x-a.x)*int64(c.y-a.y) - int64(b.y-a.y)*int64(c.x-a.x)
	if cross != 0 {
		return false
	}
	dot := int64(c.x-a.x)*int64(c.x-b.x) + int64(c.y-a.y)*int64(c.y-b.y)
	return dot <= 0
}

func ratioOK(p, e, q Point) bool {
	if !onSegment(p, e, q) {
		return false
	}
	ex, ey := int64(e.x-p.x), int64(e.y-p.y)
	dx, dy := int64(q.x-p.x), int64(q.y-p.y)
	dot := dx*ex + dy*ey
	lenSq := ex*ex + ey*ey
	if dot*5 < lenSq || dot*5 > 4*lenSq {
		return false
	}
	return true
}

func isLetterA(seg [3][2]Point) bool {
	for i := 0; i < 3; i++ {
		for j := i + 1; j < 3; j++ {
			a1, a2 := seg[i][0], seg[i][1]
			b1, b2 := seg[j][0], seg[j][1]
			var p, e1, e2 Point
			found := false
			switch {
			case a1 == b1:
				p, e1, e2 = a1, a2, b2
				found = true
			case a1 == b2:
				p, e1, e2 = a1, a2, b1
				found = true
			case a2 == b1:
				p, e1, e2 = a2, a1, b2
				found = true
			case a2 == b2:
				p, e1, e2 = a2, a1, b1
				found = true
			}
			if !found {
				continue
			}
			v1x, v1y := float64(e1.x-p.x), float64(e1.y-p.y)
			v2x, v2y := float64(e2.x-p.x), float64(e2.y-p.y)
			dot := v1x*v2x + v1y*v2y
			l1 := math.Hypot(v1x, v1y)
			l2 := math.Hypot(v2x, v2y)
			cross := v1x*v2y - v1y*v2x
			if cross == 0 || dot < 0 || dot > l1*l2 {
				continue
			}
			k := 3 - i - j
			c1, c2 := seg[k][0], seg[k][1]
			if (ratioOK(p, e1, c1) && ratioOK(p, e2, c2)) || (ratioOK(p, e1, c2) && ratioOK(p, e2, c1)) {
				return true
			}
		}
	}
	return false
}

func generateCaseB(rng *rand.Rand) (string, []string) {
	t := rng.Intn(3) + 1
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("%d\n", t))
	expect := make([]string, t)
	for i := 0; i < t; i++ {
		var seg [3][2]Point
		for j := 0; j < 3; j++ {
			seg[j][0] = Point{rng.Intn(11) - 5, rng.Intn(11) - 5}
			for {
				seg[j][1] = Point{rng.Intn(11) - 5, rng.Intn(11) - 5}
				if seg[j][1] != seg[j][0] {
					break
				}
			}
			sb.WriteString(fmt.Sprintf("%d %d %d %d\n", seg[j][0].x, seg[j][0].y, seg[j][1].x, seg[j][1].y))
		}
		if isLetterA(seg) {
			expect[i] = "YES"
		} else {
			expect[i] = "NO"
		}
	}
	return sb.String(), expect
}

func runCase(bin, input string, expected []string) error {
	out, err := run(bin, input)
	if err != nil {
		return err
	}
	scanner := bufio.NewScanner(strings.NewReader(out))
	for i := 0; i < len(expected); i++ {
		if !scanner.Scan() {
			return fmt.Errorf("not enough output for case, expected %d lines", len(expected))
		}
		line := strings.TrimSpace(scanner.Text())
		if line != expected[i] {
			return fmt.Errorf("line %d: expected '%s' got '%s'", i+1, expected[i], line)
		}
	}
	if scanner.Scan() {
		return fmt.Errorf("extra output: %s", scanner.Text())
	}
	return nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierB.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < 100; i++ {
		in, exp := generateCaseB(rng)
		if err := runCase(bin, in, exp); err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:\n%s", i+1, err, in)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
