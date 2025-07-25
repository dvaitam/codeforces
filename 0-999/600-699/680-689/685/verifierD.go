package main

import (
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"strconv"
	"strings"
)

func runBinary(bin, input string) (string, error) {
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
	err := cmd.Run()
	return out.String(), err
}

func countSquares(n, k int, pts [][2]int) []int64 {
	res := make([]int64, n+1)
	if n == 0 {
		return res
	}
	minX, maxX := pts[0][0]-k+1, pts[0][0]
	minY, maxY := pts[0][1]-k+1, pts[0][1]
	for _, p := range pts {
		if p[0]-k+1 < minX {
			minX = p[0] - k + 1
		}
		if p[0] > maxX {
			maxX = p[0]
		}
		if p[1]-k+1 < minY {
			minY = p[1] - k + 1
		}
		if p[1] > maxY {
			maxY = p[1]
		}
	}
	for x := minX; x <= maxX; x++ {
		for y := minY; y <= maxY; y++ {
			cnt := 0
			for _, p := range pts {
				if p[0] >= x && p[0] <= x+k-1 && p[1] >= y && p[1] <= y+k-1 {
					cnt++
				}
			}
			if cnt >= 1 {
				res[cnt]++
			}
		}
	}
	return res
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("usage: go run verifierD.go <binary>")
		return
	}
	bin := os.Args[1]
	rand.Seed(42)
	for tcase := 0; tcase < 100; tcase++ {
		n := rand.Intn(5) + 1
		k := rand.Intn(4) + 1
		pts := make([][2]int, n)
		for i := 0; i < n; i++ {
			pts[i][0] = rand.Intn(11) - 5
			pts[i][1] = rand.Intn(11) - 5
		}
		var sb strings.Builder
		sb.WriteString(fmt.Sprintf("%d %d\n", n, k))
		for i := 0; i < n; i++ {
			sb.WriteString(fmt.Sprintf("%d %d\n", pts[i][0], pts[i][1]))
		}
		want := countSquares(n, k, pts)
		out, err := runBinary(bin, sb.String())
		if err != nil {
			fmt.Printf("runtime error on test %d: %v\noutput:%s", tcase+1, err, out)
			return
		}
		fields := strings.Fields(out)
		if len(fields) != n {
			fmt.Printf("invalid output on test %d\ninput:%soutput:%s", tcase+1, sb.String(), out)
			return
		}
		for i := 0; i < n; i++ {
			got, err := strconv.ParseInt(fields[i], 10, 64)
			if err != nil || got != want[i+1] {
				fmt.Printf("wrong answer on test %d\ninput:%sexpected:%v\noutput:%s", tcase+1, sb.String(), want[1:], out)
				return
			}
		}
	}
	fmt.Println("All tests passed")
}
