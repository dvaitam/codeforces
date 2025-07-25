package main

import (
	"bufio"
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"strconv"
	"strings"
)

func solve(n, m, k int, cells []string) int {
	grid := make([][][]bool, n)
	idx := 0
	for i := 0; i < n; i++ {
		grid[i] = make([][]bool, m)
		for j := 0; j < m; j++ {
			s := cells[idx]
			idx++
			grid[i][j] = make([]bool, k)
			for t := 0; t < k && t < len(s); t++ {
				if s[t] == '1' {
					grid[i][j][t] = true
				}
			}
		}
	}

	pred := [3][3]int{{-1, 0, 0}, {0, -1, 0}, {0, 0, -1}}
	succ := [3][3]int{{1, 0, 0}, {0, 1, 0}, {0, 0, 1}}
	count := 0
	for x := 0; x < n; x++ {
		for y := 0; y < m; y++ {
			for z := 0; z < k; z++ {
				if !grid[x][y][z] {
					continue
				}
				preds := [][3]int{}
				succs := [][3]int{}
				for i := 0; i < 3; i++ {
					px, py, pz := x+pred[i][0], y+pred[i][1], z+pred[i][2]
					if px >= 0 && py >= 0 && pz >= 0 && grid[px][py][pz] {
						preds = append(preds, [3]int{px, py, pz})
					}
				}
				for j := 0; j < 3; j++ {
					qx, qy, qz := x+succ[j][0], y+succ[j][1], z+succ[j][2]
					if qx < n && qy < m && qz < k && grid[qx][qy][qz] {
						succs = append(succs, [3]int{qx, qy, qz})
					}
				}
				if len(preds) == 0 || len(succs) == 0 {
					continue
				}
				critical := false
				for pi, p := range preds {
					for qi := range succs {
						if critical {
							break
						}
						if pi == qi {
							critical = true
							break
						}
						rx := p[0] + succ[qi][0]
						ry := p[1] + succ[qi][1]
						rz := p[2] + succ[qi][2]
						if rx < 0 || ry < 0 || rz < 0 || rx >= n || ry >= m || rz >= k || !grid[rx][ry][rz] {
							critical = true
						}
					}
				}
				if critical {
					count++
				}
			}
		}
	}
	return count
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("usage: go run verifierD.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	f, err := os.Open("testcasesD.txt")
	if err != nil {
		panic(err)
	}
	defer f.Close()
	scanner := bufio.NewScanner(f)
	idx := 0
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" {
			continue
		}
		idx++
		parts := strings.Fields(line)
		n, _ := strconv.Atoi(parts[0])
		m, _ := strconv.Atoi(parts[1])
		k, _ := strconv.Atoi(parts[2])
		cells := parts[3:]
		expect := solve(n, m, k, cells)

		var input strings.Builder
		fmt.Fprintf(&input, "%d %d %d\n", n, m, k)
		idx2 := 0
		for i := 0; i < n; i++ {
			for j := 0; j < m; j++ {
				fmt.Fprintln(&input, cells[idx2])
				idx2++
			}
		}
		cmd := exec.Command(bin)
		cmd.Stdin = bytes.NewBufferString(input.String())
		out, err := cmd.CombinedOutput()
		if err != nil {
			fmt.Printf("Test %d: runtime error: %v\n", idx, err)
			os.Exit(1)
		}
		gotStr := strings.TrimSpace(string(out))
		got, err2 := strconv.Atoi(gotStr)
		if err2 != nil || got != expect {
			fmt.Printf("Test %d failed: expected %d got %s\n", idx, expect, gotStr)
			os.Exit(1)
		}
	}
	fmt.Printf("All %d tests passed\n", idx)
}
