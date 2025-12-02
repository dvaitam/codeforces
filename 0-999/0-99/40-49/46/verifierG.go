package main

import (
	"bytes"
	"fmt"
	"math"
	"os"
	"os/exec"
	"sort"
	"strings"
)

const M = 200010

var vis [M]bool

type pair struct{ first, second int }

var p [M]pair

func init() {
	for i := 0; i < 310; i++ {
		for j := 0; j < 310; j++ {
			s := i*i + j*j
			if s < M {
				vis[s] = true
				p[s] = pair{i, j}
			}
		}
	}
}

func getSign(arr []int) {
	type pi struct{ val, idx int }
	n := len(arr)
	arrp := make([]pi, n)
	for i := 0; i < n; i++ {
		arrp[i] = pi{arr[i], i}
	}
	sort.Slice(arrp, func(i, j int) bool { return arrp[i].val > arrp[j].val })
	d := 0
	for i := 0; i+1 < n; i += 2 {
		if arrp[i].val == arrp[i+1].val {
			idx := arrp[i].idx
			arr[idx] = -arr[idx]
		} else {
			diff := arrp[i].val - arrp[i+1].val
			if d <= 0 {
				idx := arrp[i+1].idx
				arr[idx] = -arr[idx]
				d += diff
			} else {
				idx := arrp[i].idx
				arr[idx] = -arr[idx]
				d -= diff
			}
		}
	}
	if n%2 == 1 {
		idx := arrp[n-1].idx
		arr[idx] = -arr[idx]
	}
}

// Generates a reference polygon to determine the expected minimal max-side-length.
func polygon(n int) [][2]int {
	x := make([]int, n)
	y := make([]int, n)
	cnt := 0
	sum := 0
	i := 1
	for cnt < n-1 {
		if i < M && vis[i] {
			x[cnt] = p[i].first
			y[cnt] = p[i].second
			sum += i
			cnt++
		}
		i++
	}
	var a, b, tcount int
	for j := i; tcount < 2; j++ {
		if j < M && vis[j] {
			if tcount == 0 {
				a = j
			} else {
				b = j
			}
			tcount++
		}
	}
	ok := false
	if (sum+a)%2 == 0 {
		x[cnt] = p[a].first
		y[cnt] = p[a].second
	} else if (sum+b)%2 == 0 {
		x[cnt] = p[b].first
		y[cnt] = p[b].second
	} else {
		x[cnt] = p[a].first
		y[cnt] = p[a].second
		if sum%2 != 0 {
			x[0] = p[b].first
			y[0] = p[b].second
		} else {
			if n >= 4 {
				x[3] = p[b].first
				y[3] = p[b].second
			} else {
				x[0] = p[b].first
				y[0] = p[b].second
			}
		}
		ok = true
	}
	cnt++
	sum = 0
	for k := 0; k < n; k++ {
		sum += x[k]
	}
	if sum%2 != 0 {
		if !ok {
			x[0], y[0] = y[0], x[0]
		} else {
			if n >= 4 {
				x[3], y[3] = y[3], x[3]
			} else {
				x[0], y[0] = y[0], x[0]
			}
		}
	}
	getSign(x)
	getSign(y)
	ang := make([]float64, n)
	for k := 0; k < n; k++ {
		ang[k] = math.Atan2(float64(y[k]), float64(x[k]))
	}
	ids := make([]int, n)
	for k := 0; k < n; k++ {
		ids[k] = k
	}
	sort.Slice(ids, func(i, j int) bool { return ang[ids[i]] < ang[ids[j]] })
	coords := make([][2]int, n)
	tx, ty := 0, 0
	for idx, id := range ids {
		tx += x[id]
		ty += y[id]
		coords[idx] = [2]int{tx, ty}
	}
	return coords
}

func distSq(p1, p2 [2]int) int {
	dx := p1[0] - p2[0]
	dy := p1[1] - p2[1]
	return dx*dx + dy*dy
}

func getMaxSqDist(coords [][2]int) int {
	maxD := 0
	n := len(coords)
	for i := 0; i < n; i++ {
		d := distSq(coords[i], coords[(i+1)%n])
		if d > maxD {
			maxD = d
		}
	}
	return maxD
}

func runCase(exe string, n int) error {
	input := fmt.Sprintf("%d\n", n)
	cmd := exec.Command(exe)
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("runtime error: %v\n%s", err, out.String())
	}
	
	fields := strings.Fields(out.String())
	if len(fields) == 0 {
		return fmt.Errorf("empty output")
	}
	if fields[0] != "YES" {
		return fmt.Errorf("expected YES")
	}
	if len(fields) < 1 + 2*n {
		return fmt.Errorf("insufficient output tokens")
	}
	
	userCoords := make([][2]int, n)
	idx := 1
	for i := 0; i < n; i++ {
		var x, y int
		fmt.Sscan(fields[idx], &x)
		idx++
		fmt.Sscan(fields[idx], &y)
		idx++
		userCoords[i] = [2]int{x, y}
	}
	
	// Calculate squared lengths
	userLens := make(map[int]bool)
	maxUserLen := 0
	for i := 0; i < n; i++ {
		d := distSq(userCoords[i], userCoords[(i+1)%n])
		if d == 0 {
			return fmt.Errorf("zero length side at index %d", i)
		}
		if userLens[d] {
			return fmt.Errorf("duplicate side length squared %d", d)
		}
		userLens[d] = true
		if d > maxUserLen {
			maxUserLen = d
		}
	}
	
	// Reference
	refPoly := polygon(n)
	maxRefLen := getMaxSqDist(refPoly)
	
	if maxUserLen > maxRefLen {
		return fmt.Errorf("user max length %d > expected %d", maxUserLen, maxRefLen)
	}
	
	// Convexity check
	type Vec struct { x, y int64 }
	edges := make([]Vec, n)
	for i := 0; i < n; i++ {
		edges[i] = Vec{
			x: int64(userCoords[(i+1)%n][0] - userCoords[i][0]),
			y: int64(userCoords[(i+1)%n][1] - userCoords[i][1]),
		}
	}
	
	pos := 0
	neg := 0
	for i := 0; i < n; i++ {
		v1 := edges[i]
		v2 := edges[(i+1)%n]
		cp := v1.x*v2.y - v1.y*v2.x
		if cp > 0 {
			pos++
		} else if cp < 0 {
			neg++
		}
	}
	
	if pos > 0 && neg > 0 {
		return fmt.Errorf("polygon is not convex (turns both left and right)")
	}
	
	return nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierG.go /path/to/binary")
		os.Exit(1)
	}
	exe := os.Args[1]
	for n := 3; n <= 102; n++ {
		if err := runCase(exe, n); err != nil {
			fmt.Fprintf(os.Stderr, "case n=%d failed: %v\n", n, err)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}