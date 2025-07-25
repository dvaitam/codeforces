package main

import (
	"bytes"
	"fmt"
	"math"
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

func abs(a int) int {
	if a < 0 {
		return -a
	}
	return a
}

func bestPoint(points [][3]int) (int, int, int) {
	minX, maxX := points[0][0], points[0][0]
	minY, maxY := points[0][1], points[0][1]
	minZ, maxZ := points[0][2], points[0][2]
	for _, p := range points {
		if p[0] < minX {
			minX = p[0]
		}
		if p[0] > maxX {
			maxX = p[0]
		}
		if p[1] < minY {
			minY = p[1]
		}
		if p[1] > maxY {
			maxY = p[1]
		}
		if p[2] < minZ {
			minZ = p[2]
		}
		if p[2] > maxZ {
			maxZ = p[2]
		}
	}
	bestD := math.MaxInt32
	bx, by, bz := 0, 0, 0
	for x := minX - 2; x <= maxX+2; x++ {
		for y := minY - 2; y <= maxY+2; y++ {
			for z := minZ - 2; z <= maxZ+2; z++ {
				md := 0
				for _, p := range points {
					d := abs(x-p[0]) + abs(y-p[1]) + abs(z-p[2])
					if d > md {
						md = d
					}
				}
				if md < bestD || (md == bestD && (x < bx || (x == bx && (y < by || (y == by && z < bz))))) {
					bestD = md
					bx, by, bz = x, y, z
				}
			}
		}
	}
	return bx, by, bz
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("usage: go run verifierC.go <binary>")
		return
	}
	bin := os.Args[1]
	rand.Seed(42)
	for tcase := 0; tcase < 100; tcase++ {
		n := rand.Intn(4) + 1
		points := make([][3]int, n)
		for i := 0; i < n; i++ {
			points[i][0] = rand.Intn(7) - 3
			points[i][1] = rand.Intn(7) - 3
			points[i][2] = rand.Intn(7) - 3
		}
		var sb strings.Builder
		sb.WriteString("1\n")
		sb.WriteString(strconv.Itoa(n) + "\n")
		for i := 0; i < n; i++ {
			sb.WriteString(fmt.Sprintf("%d %d %d\n", points[i][0], points[i][1], points[i][2]))
		}
		bx, by, bz := bestPoint(points)
		out, err := runBinary(bin, sb.String())
		if err != nil {
			fmt.Printf("runtime error on test %d: %v\noutput:%s", tcase+1, err, out)
			return
		}
		outFields := strings.Fields(out)
		if len(outFields) != 3 {
			fmt.Printf("invalid output on test %d\ninput:%soutput:%s", tcase+1, sb.String(), out)
			return
		}
		gx, _ := strconv.Atoi(outFields[0])
		gy, _ := strconv.Atoi(outFields[1])
		gz, _ := strconv.Atoi(outFields[2])
		if gx != bx || gy != by || gz != bz {
			fmt.Printf("wrong answer on test %d\ninput:%sexpected:%d %d %d\noutput:%s", tcase+1, sb.String(), bx, by, bz, out)
			return
		}
	}
	fmt.Println("All tests passed")
}
