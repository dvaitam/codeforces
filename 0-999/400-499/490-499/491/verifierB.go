package main

import (
	"bytes"
	"fmt"
	"math"
	"math/rand"
	"os"
	"os/exec"
	"strings"
	"time"
)

type point struct{ x, y int64 }

func runBinary(bin string, input string) (string, error) {
	cmd := exec.Command(bin)
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	err := cmd.Run()
	return strings.TrimSpace(out.String()), err
}

func abs(a int64) int64 {
	if a < 0 {
		return -a
	}
	return a
}

func expected(n, m int64, hotels, rests []point) (int64, int) {
	bestD := int64(math.MaxInt64)
	bestIdx := 1
	for i, r := range rests {
		maxd := int64(0)
		for _, h := range hotels {
			d := abs(r.x-h.x) + abs(r.y-h.y)
			if d > maxd {
				maxd = d
			}
		}
		if maxd < bestD {
			bestD = maxd
			bestIdx = i + 1
		}
	}
	return bestD, bestIdx
}

func parseOutput(out string) (int64, int, error) {
	fields := strings.Fields(out)
	if len(fields) < 2 {
		return 0, 0, fmt.Errorf("output should contain two values")
	}
	var d int64
	var idx int
	_, err := fmt.Sscan(fields[0], &d)
	if err != nil {
		return 0, 0, err
	}
	_, err = fmt.Sscan(fields[1], &idx)
	return d, idx, err
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("Usage: go run verifierB.go /path/to/binary")
		return
	}
	bin := os.Args[1]
	rand.Seed(time.Now().UnixNano())
	const tests = 100
	for t := 0; t < tests; t++ {
		n := int64(rand.Intn(50) + 1)
		m := int64(rand.Intn(50) + 1)
		C := rand.Intn(5) + 1
		H := rand.Intn(5) + 1
		hotels := make([]point, C)
		rests := make([]point, H)
		for i := 0; i < C; i++ {
			hotels[i] = point{int64(rand.Intn(int(n)) + 1), int64(rand.Intn(int(m)) + 1)}
		}
		for i := 0; i < H; i++ {
			rests[i] = point{int64(rand.Intn(int(n)) + 1), int64(rand.Intn(int(m)) + 1)}
		}
		input := fmt.Sprintf("%d %d\n%d\n", n, m, C)
		for _, h := range hotels {
			input += fmt.Sprintf("%d %d\n", h.x, h.y)
		}
		input += fmt.Sprintf("%d\n", H)
		for _, r := range rests {
			input += fmt.Sprintf("%d %d\n", r.x, r.y)
		}
		wantD, wantIdx := expected(n, m, hotels, rests)
		out, err := runBinary(bin, input)
		if err != nil {
			fmt.Printf("Test %d runtime error: %v\nOutput:\n%s\n", t+1, err, out)
			return
		}
		gotD, gotIdx, err := parseOutput(out)
		if err != nil {
			fmt.Printf("Test %d invalid output: %v\nGot: %s\n", t+1, err, out)
			return
		}
		if gotD != wantD || gotIdx < 1 || gotIdx > H || gotIdx != wantIdx && gotD == wantD {
			// allow multiple optimal indexes
			if gotD == wantD {
				// index mismatch but distance correct - check if index is also optimal
				maxd := int64(0)
				for _, h := range hotels {
					d := abs(rests[gotIdx-1].x-h.x) + abs(rests[gotIdx-1].y-h.y)
					if d > maxd {
						maxd = d
					}
				}
				if maxd == wantD {
					continue
				}
			}
			fmt.Printf("Test %d failed.\nInput:\n%sExpected: %d %d\nGot: %d %d\n", t+1, input, wantD, wantIdx, gotD, gotIdx)
			return
		}
	}
	fmt.Println("All tests passed.")
}
