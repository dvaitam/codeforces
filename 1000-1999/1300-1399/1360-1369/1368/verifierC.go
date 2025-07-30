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

func genTest() []byte {
	n := rand.Intn(50) + 1
	return []byte(fmt.Sprintf("%d\n", n))
}

type point struct{ x, y int }

func parseOutput(out string) (int, map[point]int, error) {
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
	pts := make(map[point]int, cnt)
	idx := 1
	for i := 0; i < cnt; i++ {
		x, err1 := strconv.Atoi(fields[idx])
		y, err2 := strconv.Atoi(fields[idx+1])
		if err1 != nil || err2 != nil {
			return 0, nil, fmt.Errorf("invalid integer")
		}
		pts[point{x, y}]++
		idx += 2
	}
	if k != cnt {
		return 0, nil, fmt.Errorf("declared count %d but found %d", k, cnt)
	}
	return k, pts, nil
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
	ref := "./refC.bin"
	if err := exec.Command("go", "build", "-o", ref, "1368C.go").Run(); err != nil {
		fmt.Println("failed to build reference solution:", err)
		return
	}
	defer os.Remove(ref)
	rand.Seed(time.Now().UnixNano())
	for i := 0; i < 100; i++ {
		input := genTest()
		wantRaw, err := runCmd(ref, input)
		if err != nil {
			fmt.Println("reference failed:", err)
			os.Exit(1)
		}
		gotRaw, err := runCmd(cand, input)
		if err != nil {
			fmt.Printf("candidate runtime error on test %d: %v\n", i+1, err)
			fmt.Println("input:\n", string(input))
			os.Exit(1)
		}

		wantCount, wantPts, err := parseOutput(wantRaw)
		if err != nil {
			fmt.Printf("reference output parse error on test %d: %v\n", i+1, err)
			os.Exit(1)
		}
		gotCount, gotPts, err := parseOutput(gotRaw)
		if err != nil {
			fmt.Printf("candidate output parse error on test %d: %v\n", i+1, err)
			fmt.Println("output:\n", gotRaw)
			os.Exit(1)
		}
		if wantCount != gotCount || len(wantPts) != len(gotPts) {
			fmt.Printf("wrong answer on test %d\n", i+1)
			fmt.Println("input:\n", string(input))
			fmt.Println("expected:\n", wantRaw)
			fmt.Println("got:\n", gotRaw)
			os.Exit(1)
		}
		for p := range wantPts {
			if gotPts[p] == 0 {
				fmt.Printf("wrong answer on test %d\n", i+1)
				fmt.Println("input:\n", string(input))
				fmt.Println("expected:\n", wantRaw)
				fmt.Println("got:\n", gotRaw)
				os.Exit(1)
			}
		}
		for p := range gotPts {
			if wantPts[p] == 0 {
				fmt.Printf("wrong answer on test %d\n", i+1)
				fmt.Println("input:\n", string(input))
				fmt.Println("expected:\n", wantRaw)
				fmt.Println("got:\n", gotRaw)
				os.Exit(1)
			}
		}
	}
	fmt.Println("All tests passed.")
}
