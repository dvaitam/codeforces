package main

import (
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"sort"
	"strings"
	"time"
)

func runBinary(path, input string) (string, error) {
	var cmd *exec.Cmd
	if strings.HasSuffix(path, ".go") {
		cmd = exec.Command("go", "run", path)
	} else {
		cmd = exec.Command(path)
	}
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	err := cmd.Run()
	return strings.TrimSpace(out.String()), err
}

func solveCase(n int, groups []int) bool {
	a := append([]int(nil), groups...)
	sort.Slice(a, func(i, j int) bool { return a[i] > a[j] })
	count4 := n
	count2 := 2 * n
	for _, ai := range a {
		bestCost := int(1e9)
		bestX, bestY := 0, 0
		maxX := ai/4 + 1
		if maxX > count4 {
			maxX = count4
		}
		for x := 0; x <= maxX; x++ {
			rem := ai - 4*x
			y := 0
			if rem > 0 {
				y = (rem + 1) / 2
			}
			if y > count2 {
				continue
			}
			cost := x + y
			if cost < bestCost || (cost == bestCost && x > bestX) {
				bestCost = cost
				bestX = x
				bestY = y
			}
		}
		if bestCost > count4+count2 {
			return false
		}
		count4 -= bestX
		count2 -= bestY
	}
	return true
}

func genCase(rng *rand.Rand) (int, []int) {
	n := rng.Intn(50) + 1
	k := rng.Intn(10) + 1
	groups := make([]int, k)
	remaining := 8 * n
	for i := 0; i < k; i++ {
		maxVal := remaining - (k - i - 1)
		if maxVal > 10000 {
			maxVal = 10000
		}
		if maxVal < 1 {
			maxVal = 1
		}
		groups[i] = rng.Intn(maxVal) + 1
		remaining -= groups[i]
		if remaining < 0 {
			remaining = 0
		}
	}
	return n, groups
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierB.go /path/to/binary")
		return
	}
	bin := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < 100; i++ {
		n, groups := genCase(rng)
		var sb strings.Builder
		fmt.Fprintf(&sb, "%d %d\n", n, len(groups))
		for j, v := range groups {
			if j > 0 {
				sb.WriteByte(' ')
			}
			fmt.Fprint(&sb, v)
		}
		sb.WriteByte('\n')
		input := sb.String()
		want := solveCase(n, groups)
		out, err := runBinary(bin, input)
		if err != nil {
			fmt.Printf("case %d: runtime error: %v\ninput:\n%s", i+1, err, input)
			return
		}
		got := strings.ToUpper(strings.TrimSpace(out))
		expect := "NO"
		if want {
			expect = "YES"
		}
		if got != expect {
			fmt.Printf("case %d failed\ninput:\n%sexpected: %s\ngot: %s\n", i+1, input, expect, out)
			return
		}
	}
	fmt.Println("All tests passed")
}
