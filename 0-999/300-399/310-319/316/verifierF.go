package main

import (
	"bytes"
	"fmt"
	"math"
	"math/rand"
	"os"
	"os/exec"
	"sort"
	"strings"
	"time"
)

type point struct{ x, y int }

func solveCase(h, w int, grid [][]int) string {
	visited := make([][]bool, h)
	for i := range visited {
		visited[i] = make([]bool, w)
	}
	var rayCounts []int
	dirs := [8][2]int{{1, 0}, {-1, 0}, {0, 1}, {0, -1}, {1, 1}, {1, -1}, {-1, 1}, {-1, -1}}
	for i := 0; i < h; i++ {
		for j := 0; j < w; j++ {
			if grid[i][j] == 1 && !visited[i][j] {
				var q []point
				q = append(q, point{i, j})
				visited[i][j] = true
				comp := make([]point, 0, 256)
				for qi := 0; qi < len(q); qi++ {
					p := q[qi]
					comp = append(comp, p)
					for _, d := range dirs {
						nx, ny := p.x+d[0], p.y+d[1]
						if nx >= 0 && nx < h && ny >= 0 && ny < w && grid[nx][ny] == 1 && !visited[nx][ny] {
							visited[nx][ny] = true
							q = append(q, point{nx, ny})
						}
					}
				}
				var sx, sy float64
				for _, p := range comp {
					sx += float64(p.x)
					sy += float64(p.y)
				}
				cnt := float64(len(comp))
				cx := sx / cnt
				cy := sy / cnt
				hist := make(map[int]int)
				for _, p := range comp {
					dx := float64(p.x) - cx
					dy := float64(p.y) - cy
					d := math.Hypot(dx, dy)
					r := int(d + 0.5)
					hist[r]++
				}
				r0, maxc := 0, 0
				for r, c := range hist {
					if c > maxc {
						maxc = c
						r0 = r
					}
				}
				cutoff := r0 + 2
				var rays []point
				for _, p := range comp {
					dx := float64(p.x) - cx
					dy := float64(p.y) - cy
					d := math.Hypot(dx, dy)
					if int(d+0.5) > cutoff {
						rays = append(rays, p)
					}
				}
				if len(rays) == 0 {
					rayCounts = append(rayCounts, 0)
					continue
				}
				minx, miny := rays[0].x, rays[0].y
				maxx, maxy := minx, miny
				for _, p := range rays {
					if p.x < minx {
						minx = p.x
					}
					if p.x > maxx {
						maxx = p.x
					}
					if p.y < miny {
						miny = p.y
					}
					if p.y > maxy {
						maxy = p.y
					}
				}
				hh := maxx - minx + 1
				ww := maxy - miny + 1
				mask := make([][]bool, hh)
				vis2 := make([][]bool, hh)
				for ii := 0; ii < hh; ii++ {
					mask[ii] = make([]bool, ww)
					vis2[ii] = make([]bool, ww)
				}
				for _, p := range rays {
					mask[p.x-minx][p.y-miny] = true
				}
				rc := 0
				for ii := 0; ii < hh; ii++ {
					for jj := 0; jj < ww; jj++ {
						if mask[ii][jj] && !vis2[ii][jj] {
							rc++
							var q2 []point
							q2 = append(q2, point{ii, jj})
							vis2[ii][jj] = true
							for qi := 0; qi < len(q2); qi++ {
								pp := q2[qi]
								for _, d := range dirs {
									nx, ny := pp.x+d[0], pp.y+d[1]
									if nx >= 0 && nx < hh && ny >= 0 && ny < ww && mask[nx][ny] && !vis2[nx][ny] {
										vis2[nx][ny] = true
										q2 = append(q2, point{nx, ny})
									}
								}
							}
						}
					}
				}
				rayCounts = append(rayCounts, rc)
			}
		}
	}
	sort.Ints(rayCounts)
	var out strings.Builder
	fmt.Fprintf(&out, "%d\n", len(rayCounts))
	for i, v := range rayCounts {
		if i > 0 {
			out.WriteByte(' ')
		}
		fmt.Fprintf(&out, "%d", v)
	}
	out.WriteByte('\n')
	return out.String()
}

func genCase(rng *rand.Rand) (string, string) {
	h := rng.Intn(3) + 2
	w := rng.Intn(3) + 2
	grid := make([][]int, h)
	var sb strings.Builder
	fmt.Fprintf(&sb, "%d %d\n", h, w)
	for i := 0; i < h; i++ {
		grid[i] = make([]int, w)
		for j := 0; j < w; j++ {
			v := rng.Intn(2)
			grid[i][j] = v
			if j > 0 {
				sb.WriteByte(' ')
			}
			fmt.Fprintf(&sb, "%d", v)
		}
		sb.WriteByte('\n')
	}
	expected := solveCase(h, w, grid)
	return sb.String(), expected
}

func runCase(bin, input, expected string) error {
	cmd := exec.Command(bin)
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("runtime error: %v\n%s", err, out.String())
	}
	if strings.TrimSpace(out.String()) != strings.TrimSpace(expected) {
		return fmt.Errorf("expected %s got %s", expected, out.String())
	}
	return nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierF.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < 100; i++ {
		in, exp := genCase(rng)
		if err := runCase(bin, in, exp); err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:\n%s", i+1, err, in)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
