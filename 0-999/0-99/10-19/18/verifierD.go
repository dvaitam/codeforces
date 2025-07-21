package main

import (
	"bytes"
	"fmt"
	"math/big"
	"math/rand"
	"os"
	"os/exec"
	"sort"
	"strings"
	"time"
)

type event struct {
	typ string
	x   int
}

type interval struct {
	start, end int
	x          int
	weight     *big.Int
}

func solve(events []event) *big.Int {
	winsList := make(map[int][]int)
	intervals := make([]interval, 0)
	for i, e := range events {
		if e.typ == "win" {
			winsList[e.x] = append(winsList[e.x], i)
		} else {
			list := winsList[e.x]
			if len(list) > 0 {
				start := list[len(list)-1]
				intervals = append(intervals, interval{start: start, end: i, x: e.x})
			}
		}
	}
	if len(intervals) == 0 {
		return big.NewInt(0)
	}
	sort.Slice(intervals, func(i, j int) bool { return intervals[i].end < intervals[j].end })
	m := len(intervals)
	ends := make([]int, m)
	maxX := 0
	for i, iv := range intervals {
		ends[i] = iv.end
		if iv.x > maxX {
			maxX = iv.x
		}
	}
	pow2 := make([]*big.Int, maxX+1)
	pow2[0] = big.NewInt(1)
	for i := 1; i <= maxX; i++ {
		pow2[i] = new(big.Int).Lsh(pow2[i-1], 1)
	}
	for i := range intervals {
		intervals[i].weight = pow2[intervals[i].x]
	}
	dp := make([]*big.Int, m+1)
	dp[0] = big.NewInt(0)
	for j := 1; j <= m; j++ {
		best := new(big.Int).Set(dp[j-1])
		iv := intervals[j-1]
		k := sort.Search(m, func(i int) bool { return ends[i] >= iv.start })
		with := new(big.Int).Set(iv.weight)
		if k > 0 {
			with.Add(with, dp[k])
		}
		if with.Cmp(best) > 0 {
			dp[j] = with
		} else {
			dp[j] = best
		}
	}
	return dp[m]
}

func generateCase(rng *rand.Rand) (string, string) {
	n := rng.Intn(10) + 1
	events := make([]event, n)
	usedSell := make(map[int]bool)
	for i := 0; i < n; i++ {
		typ := "win"
		if rng.Intn(2) == 0 {
			typ = "sell"
		}
		x := rng.Intn(6)
		if typ == "sell" {
			if usedSell[x] {
				typ = "win"
			} else {
				usedSell[x] = true
			}
		}
		events[i] = event{typ: typ, x: x}
	}
	if solve(events).Cmp(big.NewInt(0)) == 0 {
		events = append(events, event{typ: "win", x: 1}, event{typ: "sell", x: 1})
	}
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("%d\n", len(events)))
	for _, e := range events {
		sb.WriteString(fmt.Sprintf("%s %d\n", e.typ, e.x))
	}
	ans := solve(events)
	return sb.String(), ans.String()
}

func runCase(bin, input, exp string) error {
	cmd := exec.Command(bin)
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("runtime error: %v\n%s", err, out.String())
	}
	got := strings.TrimSpace(out.String())
	if got != exp {
		return fmt.Errorf("expected %s got %s", exp, got)
	}
	return nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierD.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < 100; i++ {
		input, exp := generateCase(rng)
		if err := runCase(bin, input, exp); err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:\n%s", i+1, err, input)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
