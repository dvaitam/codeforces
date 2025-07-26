package main

import (
	"bytes"
	"fmt"
	"math"
	"math/rand"
	"os"
	"os/exec"
	"sort"
	"strconv"
	"strings"
	"time"
)

type node struct{ x, w int }

func bestCopies(s []int, k int) int {
	freq := map[int]int{}
	maxF := 0
	for _, v := range s {
		freq[v]++
		if freq[v] > maxF {
			maxF = freq[v]
		}
	}
	nodes := make([]node, 0, len(freq))
	for val, w := range freq {
		nodes = append(nodes, node{val, w})
	}
	sort.Slice(nodes, func(i, j int) bool { return nodes[i].w > nodes[j].w })
	l, r := 1, maxF
	for l <= r {
		mid := (l + r) / 2
		sum := 0
		for _, nd := range nodes {
			sum += nd.w / mid
			if sum >= k {
				break
			}
		}
		if sum >= k {
			l = mid + 1
		} else {
			r = mid - 1
		}
	}
	return r
}

func copiesPossible(s []int, t []int) int {
	freqS := map[int]int{}
	for _, v := range s {
		freqS[v]++
	}
	freqT := map[int]int{}
	for _, v := range t {
		freqT[v]++
	}
	res := math.MaxInt32
	for val, c := range freqT {
		q := freqS[val] / c
		if q < res {
			res = q
		}
	}
	return res
}

func runCase(exe string, input string, s []int, k int, expR int) error {
	cmd := exec.Command(exe)
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("runtime error: %v\n%s", err, out.String())
	}
	fields := strings.Fields(out.String())
	if len(fields) != k {
		return fmt.Errorf("expected %d numbers, got %d", k, len(fields))
	}
	tArr := make([]int, k)
	for i, f := range fields {
		v, err := strconv.Atoi(f)
		if err != nil {
			return fmt.Errorf("invalid integer %q", f)
		}
		tArr[i] = v
	}
	rActual := copiesPossible(s, tArr)
	if rActual != expR {
		return fmt.Errorf("expected copies %d, got %d", expR, rActual)
	}
	return nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierD.go /path/to/binary")
		os.Exit(1)
	}
	exe := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for tcase := 0; tcase < 100; tcase++ {
		n := rng.Intn(100) + 1
		k := rng.Intn(n) + 1
		s := make([]int, n)
		var sb strings.Builder
		fmt.Fprintf(&sb, "%d %d\n", n, k)
		for i := 0; i < n; i++ {
			s[i] = rng.Intn(20) + 1
			if i > 0 {
				sb.WriteByte(' ')
			}
			sb.WriteString(strconv.Itoa(s[i]))
		}
		sb.WriteByte('\n')
		input := sb.String()
		expR := bestCopies(s, k)
		if err := runCase(exe, input, s, k, expR); err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:\n%s", tcase+1, err, input)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
