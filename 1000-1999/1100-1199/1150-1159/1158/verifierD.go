package main

import (
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"time"
)

type Vec struct {
	x, y int64
	id   int
}

func solveD(points []Vec, s string) []int {
	n := len(points)
	used := make([]bool, n)
	ans := make([]int, n)
	start := 0
	for i := 1; i < n; i++ {
		if points[i].x < points[start].x || (points[i].x == points[start].x && points[i].y < points[start].y) {
			start = i
		}
	}
	used[start] = true
	ans[0] = points[start].id
	curr := points[start]
	prev := Vec{x: 1, y: 0}
	for k := 0; k < len(s); k++ {
		best := -1
		var bestVec Vec
		for i := 0; i < n; i++ {
			if used[i] {
				continue
			}
			v := Vec{x: points[i].x - curr.x, y: points[i].y - curr.y, id: points[i].id}
			if best == -1 {
				best = i
				bestVec = v
				continue
			}
			u := bestVec
			cdu := prev.x*u.y - prev.y*u.x
			cdv := prev.x*v.y - prev.y*v.x
			if s[k] == 'L' {
				useV := false
				if (cdu >= 0) != (cdv >= 0) {
					if cdv >= 0 {
						useV = true
					}
				} else if v.x*u.y-v.y*u.x > 0 {
					useV = true
				}
				if useV {
					best = i
					bestVec = v
				}
			} else {
				useV := false
				if (cdu >= 0) != (cdv >= 0) {
					if cdv < 0 {
						useV = true
					}
				} else if u.x*v.y-u.y*v.x > 0 {
					useV = true
				}
				if useV {
					best = i
					bestVec = v
				}
			}
		}
		used[best] = true
		ans[k+1] = points[best].id
		curr = points[best]
		prev = bestVec
	}
	for i := 0; i < n; i++ {
		if !used[i] {
			ans[n-1] = points[i].id
			break
		}
	}
	return ans
}

func run(bin, input string) (string, error) {
	cmd := exec.Command(bin)
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("runtime error: %v\n%s", err, out.String())
	}
	return strings.TrimSpace(out.String()), nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierD.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for t := 0; t < 100; t++ {
		n := rng.Intn(5) + 3
		pts := make([]Vec, 0, n)
		for len(pts) < n {
			x := int64(rng.Intn(21) - 10)
			y := int64(rng.Intn(21) - 10)
			dup := false
			for _, p := range pts {
				if p.x == x && p.y == y {
					dup = true
					break
				}
			}
			if dup {
				continue
			}
			pts = append(pts, Vec{x: x, y: y, id: len(pts) + 1})
		}
		// ensure no three collinear by regeneration
		ok := false
		for !ok {
			ok = true
			for i := 0; i < n && ok; i++ {
				for j := i + 1; j < n && ok; j++ {
					for k := j + 1; k < n; k++ {
						if (pts[j].x-pts[i].x)*(pts[k].y-pts[i].y)-(pts[j].y-pts[i].y)*(pts[k].x-pts[i].x) == 0 {
							// regenerate one point
							pts[k].x = int64(rng.Intn(21) - 10)
							pts[k].y = int64(rng.Intn(21) - 10)
							k = -1
							ok = false
							break
						}
					}
				}
			}
		}
		sLen := n - 2
		var sb strings.Builder
		for i := 0; i < sLen; i++ {
			if rng.Intn(2) == 0 {
				sb.WriteByte('L')
			} else {
				sb.WriteByte('R')
			}
		}
		str := sb.String()
		input := fmt.Sprintf("%d\n", n)
		for _, p := range pts {
			input += fmt.Sprintf("%d %d\n", p.x, p.y)
		}
		input += str + "\n"
		exp := solveD(pts, str)
		expStr := strings.TrimSpace(strings.Join(func() []string {
			tmp := make([]string, len(exp))
			for i, v := range exp {
				tmp[i] = strconv.Itoa(v)
			}
			return tmp
		}(), " "))
		got, err := run(bin, input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:%s", t+1, err, input)
			os.Exit(1)
		}
		if strings.TrimSpace(got) != expStr {
			fmt.Fprintf(os.Stderr, "case %d failed: expected %s got %s\ninput:%s", t+1, expStr, got, input)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
