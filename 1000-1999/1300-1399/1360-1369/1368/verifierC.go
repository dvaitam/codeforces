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

type point struct{ x, y int }

func runCmd(path string, input []byte) (string, error) {
	var cmd *exec.Cmd
	if strings.HasSuffix(path, ".go") {
		cmd = exec.Command("go", "run", path)
	} else {
		cmd = exec.Command(path)
	}
	cmd.Stdin = bytes.NewReader(input)
	out, err := cmd.CombinedOutput()
	return string(out), err
}

func genTest() (int, []byte) {
	n := rand.Intn(50) + 1
	return n, []byte(fmt.Sprintf("%d\n", n))
}

func parseOutput(out string) (int, []point, error) {
	fields := strings.Fields(strings.TrimSpace(out))
	if len(fields) == 0 {
		return 0, nil, fmt.Errorf("empty output")
	}
	k, err := strconv.Atoi(fields[0])
	if err != nil {
		return 0, nil, err
	}
	if (len(fields)-1)%2 != 0 {
		return 0, nil, fmt.Errorf("incomplete coordinate pair")
	}
	cnt := (len(fields) - 1) / 2
	pts := make([]point, cnt)
	idx := 1
	for i := 0; i < cnt; i++ {
		x, err1 := strconv.Atoi(fields[idx])
		y, err2 := strconv.Atoi(fields[idx+1])
		if err1 != nil || err2 != nil {
			return 0, nil, fmt.Errorf("invalid integer")
		}
		pts[i] = point{x, y}
		idx += 2
	}
	if k != cnt {
		return 0, nil, fmt.Errorf("declared count %d but found %d", k, cnt)
	}
	return k, pts, nil
}

func check(n int, pts []point) error {
	if len(pts) > 500000 {
		return fmt.Errorf("too many points: %d", len(pts))
	}
	mp := make(map[point]struct{}, len(pts))
	for _, p := range pts {
		if p.x < -1e9 || p.x > 1e9 || p.y < -1e9 || p.y > 1e9 {
			return fmt.Errorf("coordinate out of bounds: %v", p)
		}
		if _, ok := mp[p]; ok {
			return fmt.Errorf("duplicate point: %v", p)
		}
		mp[p] = struct{}{}
	}
	if len(mp) == 0 {
		return fmt.Errorf("no points provided")
	}
	dirs := []point{{1, 0}, {-1, 0}, {0, 1}, {0, -1}}
	// connectivity
	visited := make(map[point]bool, len(mp))
	queue := make([]point, 0, len(mp))
	for p := range mp {
		queue = append(queue, p)
		visited[p] = true
		break
	}
	for len(queue) > 0 {
		cur := queue[0]
		queue = queue[1:]
		for _, d := range dirs {
			np := point{cur.x + d.x, cur.y + d.y}
			if _, ok := mp[np]; ok && !visited[np] {
				visited[np] = true
				queue = append(queue, np)
			}
		}
	}
	if len(visited) != len(mp) {
		return fmt.Errorf("picture is not connected")
	}
	countDeg4 := 0
	for p := range mp {
		deg := 0
		for _, d := range dirs {
			np := point{p.x + d.x, p.y + d.y}
			if _, ok := mp[np]; ok {
				deg++
			}
		}
		if deg%2 != 0 {
			return fmt.Errorf("cell %v has odd number of neighbours %d", p, deg)
		}
		if deg == 4 {
			countDeg4++
		}
	}
	if countDeg4 != n {
		return fmt.Errorf("expected %d cells with 4 neighbours, got %d", n, countDeg4)
	}
	return nil
}

func main() {
	var cand string
	if len(os.Args) == 2 {
		cand = os.Args[1]
	} else if len(os.Args) == 3 && os.Args[1] == "--" {
		cand = os.Args[2]
	} else {
		fmt.Println("usage: go run verifierC.go /path/to/binary")
		return
	}
	rand.Seed(time.Now().UnixNano())
	for i := 0; i < 100; i++ {
		n, input := genTest()
		gotRaw, err := runCmd(cand, input)
		if err != nil {
			fmt.Printf("candidate runtime error on test %d: %v\n", i+1, err)
			fmt.Println("input:\n", string(input))
			os.Exit(1)
		}
		_, pts, err := parseOutput(gotRaw)
		if err != nil {
			fmt.Printf("candidate output parse error on test %d: %v\n", i+1, err)
			fmt.Println("input:\n", string(input))
			fmt.Println("output:\n", gotRaw)
			os.Exit(1)
		}
		if err := check(n, pts); err != nil {
			fmt.Printf("wrong answer on test %d: %v\n", i+1, err)
			fmt.Println("input:\n", string(input))
			fmt.Println("output:\n", gotRaw)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed.")
}
