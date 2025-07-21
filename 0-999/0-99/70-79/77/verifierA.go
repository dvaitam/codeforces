package main

import (
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"strings"
)

type likeEdge struct{ p, q int }

var names = []string{"Anka", "Chapay", "Cleo", "Troll", "Dracul", "Snowy", "Hexadecimal"}

func solveA(edges []likeEdge, xp [3]int64) (int64, int) {
	bestDiff := int64(1<<62 - 1)
	bestLikes := 0
	assign := [7]int{}
	for mask := 0; mask < 2187; mask++ { // 3^7
		m := mask
		sz := [3]int{}
		for i := 0; i < 7; i++ {
			assign[i] = m % 3
			sz[assign[i]]++
			m /= 3
		}
		if sz[0] == 0 || sz[1] == 0 || sz[2] == 0 {
			continue
		}
		xpGroup := [3]int64{}
		for i := 0; i < 3; i++ {
			xpGroup[i] = xp[i] / int64(sz[i])
		}
		minXP, maxXP := xpGroup[assign[0]], xpGroup[assign[0]]
		for i := 1; i < 7; i++ {
			x := xpGroup[assign[i]]
			if x < minXP {
				minXP = x
			}
			if x > maxXP {
				maxXP = x
			}
		}
		diff := maxXP - minXP
		likes := 0
		for _, e := range edges {
			if assign[e.p] == assign[e.q] {
				likes++
			}
		}
		if diff < bestDiff || (diff == bestDiff && likes > bestLikes) {
			bestDiff = diff
			bestLikes = likes
		}
	}
	return bestDiff, bestLikes
}

func run(bin string, input string) (string, error) {
	var cmd *exec.Cmd
	if strings.HasSuffix(bin, ".go") {
		cmd = exec.Command("go", "run", bin)
	} else {
		cmd = exec.Command(bin)
	}
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	var errBuf bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &errBuf
	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("runtime error: %v\n%s", err, errBuf.String())
	}
	return strings.TrimSpace(out.String()), nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierA.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rand.Seed(42)
	// generate 100 tests
	tests := 100
	allPairs := make([]likeEdge, 0, 42)
	for i := 0; i < 7; i++ {
		for j := 0; j < 7; j++ {
			if i == j {
				continue
			}
			allPairs = append(allPairs, likeEdge{i, j})
		}
	}
	for t := 1; t <= tests; t++ {
		n := rand.Intn(43) // 0..42
		perm := rand.Perm(len(allPairs))[:n]
		edges := make([]likeEdge, n)
		for i, idx := range perm {
			edges[i] = allPairs[idx]
		}
		xp := [3]int64{rand.Int63n(2_000_000_000) + 1, rand.Int63n(2_000_000_000) + 1, rand.Int63n(2_000_000_000) + 1}
		var input bytes.Buffer
		fmt.Fprintln(&input, n)
		for _, e := range edges {
			fmt.Fprintf(&input, "%s likes %s\n", names[e.p], names[e.q])
		}
		fmt.Fprintf(&input, "%d %d %d\n", xp[0], xp[1], xp[2])

		expDiff, expLikes := solveA(edges, xp)
		gotStr, err := run(bin, input.String())
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\n", t, err)
			os.Exit(1)
		}
		var gotDiff int64
		var gotLikes int
		fmt.Sscan(gotStr, &gotDiff, &gotLikes)
		if gotDiff != expDiff || gotLikes != expLikes {
			fmt.Fprintf(os.Stderr, "case %d failed: expected %d %d got %s\n", t, expDiff, expLikes, gotStr)
			os.Exit(1)
		}
	}
	fmt.Printf("All %d tests passed\n", tests)
}
