package main

import (
	"bufio"
	"fmt"
	"math"
	"math/rand"
	"os"
	"os/exec"
	"sort"
	"strings"
	"time"
)

type Point struct{ x, y int64 }

func cross(o, a, b Point) int64 {
	return (a.x-o.x)*(b.y-o.y) - (a.y-o.y)*(b.x-o.x)
}

func triArea2(a, b, c Point) int64 {
	v := cross(a, b, c)
	if v < 0 {
		v = -v
	}
	return v // 2 * triangle area
}

func polygonArea2(pts []Point) int64 {
	n := len(pts)
	var s int64
	for i := 0; i < n; i++ {
		j := (i + 1) % n
		s += pts[i].x*pts[j].y - pts[i].y*pts[j].x
	}
	if s < 0 {
		s = -s
	}
	return s // 2 * polygon area
}

// genCase generates a strictly convex CCW polygon with integer coords, n >= 4.
func genCase(rng *rand.Rand, n int) ([]Point, string) {
	for {
		angles := make([]float64, n)
		for i := range angles {
			angles[i] = rng.Float64() * 2 * math.Pi
		}
		sort.Float64s(angles)
		// Ensure angles are well-separated to avoid collinear integer points
		ok := true
		for i := 0; i < n; i++ {
			next := (i + 1) % n
			diff := angles[next] - angles[i]
			if next == 0 {
				diff = 2*math.Pi - angles[n-1] + angles[0]
			}
			if diff < 0.3 {
				ok = false
				break
			}
		}
		if !ok {
			continue
		}

		R := 100000.0
		pts := make([]Point, n)
		for i, ang := range angles {
			pts[i] = Point{int64(math.Round(R * math.Cos(ang))), int64(math.Round(R * math.Sin(ang)))}
		}

		// Verify strictly convex and CCW
		convex := true
		for i := 0; i < n; i++ {
			a, b, c := pts[i], pts[(i+1)%n], pts[(i+2)%n]
			if cross(a, b, c) <= 0 {
				convex = false
				break
			}
		}
		if !convex {
			continue
		}

		var sb strings.Builder
		sb.WriteString(fmt.Sprintf("%d\n", n))
		for _, p := range pts {
			sb.WriteString(fmt.Sprintf("%d %d\n", p.x, p.y))
		}
		return pts, sb.String()
	}
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierI.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))

	for test := 1; test <= 200; test++ {
		n := rng.Intn(7) + 4 // 4..10
		pts, inputStr := genCase(rng, n)
		total2 := polygonArea2(pts) // 2 * total area

		cmd := exec.Command(bin)
		stdinPipe, err := cmd.StdinPipe()
		if err != nil {
			fmt.Fprintln(os.Stderr, "stdin pipe:", err)
			os.Exit(1)
		}
		stdoutPipe, err := cmd.StdoutPipe()
		if err != nil {
			fmt.Fprintln(os.Stderr, "stdout pipe:", err)
			os.Exit(1)
		}
		if err := cmd.Start(); err != nil {
			fmt.Fprintln(os.Stderr, "start:", err)
			os.Exit(1)
		}

		reader := bufio.NewReader(stdoutPipe)

		// Send the polygon input
		fmt.Fprint(stdinPipe, inputStr)

		fail := func(msg string) {
			fmt.Printf("Test %d failed\nInput:\n%s%s\n", test, inputStr, msg)
			cmd.Process.Kill()
			os.Exit(1)
		}

		// Read which player the solution is helping
		line, err := reader.ReadString('\n')
		if err != nil {
			fail(fmt.Sprintf("error reading first output: %v", err))
		}
		helping := strings.TrimSpace(line)
		if helping != "Alberto" && helping != "Beatrice" && helping != "Either" {
			fail(fmt.Sprintf("invalid first output: %q", helping))
		}
		helpedPlayer := helping
		if helpedPlayer == "Either" {
			helpedPlayer = "Alberto"
		}

		// Doubly-linked list for remaining vertices (1-indexed)
		nxt := make([]int, n+1)
		prv := make([]int, n+1)
		removed := make([]bool, n+1)
		for i := 1; i <= n; i++ {
			nxt[i] = i%n + 1
			prv[i%n+1] = i
		}

		doRemove := func(v int) int64 {
			p, nx := prv[v], nxt[v]
			ar2 := triArea2(pts[p-1], pts[v-1], pts[nx-1])
			nxt[p] = nx
			prv[nx] = p
			removed[v] = true
			return ar2
		}

		// Judge picks the min-area ear for the opponent
		judgeMove := func() int {
			bestV := -1
			var bestAr int64 = math.MaxInt64
			for i := 1; i <= n; i++ {
				if !removed[i] {
					p, nx := prv[i], nxt[i]
					ar2 := triArea2(pts[p-1], pts[i-1], pts[nx-1])
					if ar2 < bestAr {
						bestAr = ar2
						bestV = i
					}
				}
			}
			return bestV
		}

		// areas2[0]=Alberto, areas2[1]=Beatrice (stored as 2*area)
		areas2 := [2]int64{}
		playerIdx := map[string]int{"Alberto": 0, "Beatrice": 1}

		totalTurns := n - 2
		for turn := 0; turn < totalTurns; turn++ {
			curPlayer := "Alberto"
			if turn%2 == 1 {
				curPlayer = "Beatrice"
			}

			if curPlayer == helpedPlayer {
				// Solution makes this move
				line, err = reader.ReadString('\n')
				if err != nil {
					fail(fmt.Sprintf("error reading solution move on turn %d: %v", turn, err))
				}
				var v int
				if _, scanErr := fmt.Sscan(strings.TrimSpace(line), &v); scanErr != nil || v < 1 || v > n || removed[v] {
					fail(fmt.Sprintf("invalid move %q on turn %d", strings.TrimSpace(line), turn))
				}
				ar2 := doRemove(v)
				areas2[playerIdx[curPlayer]] += ar2
			} else {
				// Judge makes this move
				v := judgeMove()
				ar2 := doRemove(v)
				areas2[playerIdx[curPlayer]] += ar2
				fmt.Fprintf(stdinPipe, "%d\n", v)
			}
		}

		stdinPipe.Close()
		cmd.Wait()

		// Verify: helped player must have eaten <= half (i.e. 2*area <= total area)
		helpedIdx := playerIdx[helpedPlayer]
		if areas2[helpedIdx] > total2 {
			fail(fmt.Sprintf("helped player %s ate %d/2 > total %d/2",
				helpedPlayer, areas2[helpedIdx], total2))
		}
	}
	fmt.Println("All tests passed")
}
