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

func run(bin, input string) (string, error) {
	var cmd *exec.Cmd
	if strings.HasSuffix(bin, ".go") {
		cmd = exec.Command("go", "run", bin)
	} else {
		cmd = exec.Command(bin)
	}
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	var errb bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &errb
	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("runtime error: %v\n%s", err, errb.String())
	}
	return strings.TrimSpace(out.String()), nil
}

// Implementation of the reference solution from 1242C.go

type cycle struct {
	mask  int
	nodes []int64
}

type move struct {
	num  int64
	dest int
}

func solveC(k int, boxes [][]int64) (string, []move) {
	boxsum := make([]int64, k)
	whichbox := make(map[int64]int)
	var total int64
	var nodes []int64
	for i := 0; i < k; i++ {
		for _, v := range boxes[i] {
			whichbox[v] = i
			nodes = append(nodes, v)
			total += v
			boxsum[i] += v
		}
	}
	if total%int64(k) != 0 {
		return "NO", nil
	}
	target := total / int64(k)
	nxt := make(map[int64]int64, len(nodes))
	for i := 0; i < k; i++ {
		for _, v := range boxes[i] {
			needed := target - boxsum[i] + v
			if _, ok := whichbox[needed]; ok {
				nxt[v] = needed
			} else {
				nxt[v] = -1
			}
		}
	}
	processed := make(map[int64]bool)
	var validcycles []cycle
	for _, start := range nodes {
		if processed[start] {
			continue
		}
		position := make(map[int64]int)
		var path []int64
		cur := start
		found := false
		for cur != -1 {
			if _, seen := position[cur]; seen {
				found = true
				break
			}
			position[cur] = len(path)
			path = append(path, cur)
			cur = nxt[cur]
		}
		if found {
			pos := position[cur]
			var mask int
			var cyc []int64
			ok := true
			for _, v := range path[pos:] {
				b := whichbox[v]
				if mask&(1<<b) != 0 {
					ok = false
					break
				}
				mask |= 1 << b
				cyc = append(cyc, v)
			}
			if ok {
				validcycles = append(validcycles, cycle{mask: mask, nodes: cyc})
			}
		}
		for _, v := range path {
			processed[v] = true
		}
	}
	fullMask := (1 << k) - 1
	cyclesByMask := make([][]int, fullMask+1)
	for i, c := range validcycles {
		cyclesByMask[c.mask] = append(cyclesByMask[c.mask], i)
	}
	dp := make([]bool, fullMask+1)
	parent := make([]int, fullMask+1)
	used := make([]int, fullMask+1)
	f := make([]bool, fullMask+1)
	dp[0] = true
	for _, c := range validcycles {
		f[c.mask] = true
	}
	for mask := 0; mask <= fullMask; mask++ {
		if !dp[mask] {
			continue
		}
		remain := fullMask ^ mask
		for s := remain; s > 0; s = (s - 1) & remain {
			if !f[s] {
				continue
			}
			newMask := mask | s
			if !dp[newMask] {
				dp[newMask] = true
				parent[newMask] = mask
				used[newMask] = s
			}
		}
	}
	if !dp[fullMask] {
		return "NO", nil
	}
	var chain []int
	for m := fullMask; m != 0; m = parent[m] {
		chain = append(chain, used[m])
	}
	for i, j := 0, len(chain)-1; i < j; i, j = i+1, j-1 {
		chain[i], chain[j] = chain[j], chain[i]
	}
	var ans []move
	for _, sm := range chain {
		ci := cyclesByMask[sm][0]
		cyc := validcycles[ci].nodes
		for i, j := 0, len(cyc)-1; i < j; i, j = i+1, j-1 {
			cyc[i], cyc[j] = cyc[j], cyc[i]
		}
		t := len(cyc)
		for i := 0; i < t; i++ {
			num := cyc[i]
			destBox := whichbox[cyc[(i+1)%t]] + 1
			ans = append(ans, move{num: num, dest: destBox})
		}
	}
	sort.Slice(ans, func(i, j int) bool {
		bi := whichbox[ans[i].num]
		bj := whichbox[ans[j].num]
		return bi < bj
	})
	return "YES", ans
}

func formatInput(k int, boxes [][]int64) string {
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("%d\n", k))
	for _, box := range boxes {
		sb.WriteString(fmt.Sprintf("%d ", len(box)))
		for i, v := range box {
			if i+1 == len(box) {
				sb.WriteString(fmt.Sprintf("%d\n", v))
			} else {
				sb.WriteString(fmt.Sprintf("%d ", v))
			}
		}
	}
	return sb.String()
}

func formatOutput(status string, moves []move) string {
	if status == "NO" {
		return "NO"
	}
	var sb strings.Builder
	sb.WriteString("YES\n")
	for _, mv := range moves {
		sb.WriteString(fmt.Sprintf("%d %d\n", mv.num, mv.dest))
	}
	return strings.TrimSpace(sb.String())
}

func genCase(rng *rand.Rand) (int, [][]int64) {
	k := rng.Intn(3) + 1
	boxes := make([][]int64, k)
	used := make(map[int64]bool)
	for i := 0; i < k; i++ {
		n := rng.Intn(3) + 1
		boxes[i] = make([]int64, n)
		for j := 0; j < n; j++ {
			var v int64
			for {
				v = int64(rng.Intn(20) - 10)
				if !used[v] {
					break
				}
			}
			used[v] = true
			boxes[i][j] = v
		}
	}
	return k, boxes
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierC.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))

	var cases []struct {
		k     int
		boxes [][]int64
	}
	k1 := 1
	boxes1 := [][]int64{{5}}
	cases = append(cases, struct {
		k     int
		boxes [][]int64
	}{k1, boxes1})
	for len(cases) < 100 {
		k, b := genCase(rng)
		cases = append(cases, struct {
			k     int
			boxes [][]int64
		}{k, b})
	}

	for i, tc := range cases {
		in := formatInput(tc.k, tc.boxes)
		status, moves := solveC(tc.k, tc.boxes)
		exp := formatOutput(status, moves)
		out, err := run(bin, in)
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:\n%s", i+1, err, in)
			os.Exit(1)
		}
		if strings.TrimSpace(out) != exp {
			fmt.Fprintf(os.Stderr, "case %d failed: expected:\n%s\n got:\n%s\ninput:\n%s", i+1, exp, out, in)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
