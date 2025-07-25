package main

import (
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"sort"
	"strconv"
	"strings"
	"time"
)

type testCase struct {
	ts, tf, t int64
	arrivals  []int64
}

func solve(ts, tf, t int64, arr []int64) string {
	sort.Slice(arr, func(i, j int) bool { return arr[i] < arr[j] })
	arr = append(arr, int64(1)<<62)
	cur := make([]int64, len(arr))
	cur[0] = ts
	for i := 0; i < len(arr)-1; i++ {
		start := cur[i]
		if arr[i] > start {
			start = arr[i]
		}
		cur[i+1] = start + t
	}
	for i := 0; i < len(arr)-1; i++ {
		if cur[i]+t <= tf && cur[i] < arr[i] {
			return fmt.Sprintf("%d", cur[i])
		}
	}
	bestTime := ts
	bestWait := int64(1<<63 - 1)
	j := sort.Search(len(arr), func(i int) bool { return arr[i] > ts })
	start := cur[j]
	if start < ts {
		start = ts
	}
	if start+t <= tf {
		bestWait = start - ts
		bestTime = ts
	}
	for i := 0; i < len(arr)-1; i++ {
		x := arr[i] - 1
		if x < 0 {
			continue
		}
		j = sort.Search(len(arr), func(k int) bool { return arr[k] > x })
		start = cur[j]
		if start < x {
			start = x
		}
		if start+t <= tf {
			wait := start - x
			if wait < bestWait {
				bestWait = wait
				bestTime = x
			}
		}
	}
	if cur[len(arr)-1]+t <= tf && 0 < bestWait {
		bestTime = cur[len(arr)-1]
	}
	return fmt.Sprintf("%d", bestTime)
}

func (tc testCase) input() string {
	var sb strings.Builder
	fmt.Fprintf(&sb, "%d %d %d\n", tc.ts, tc.tf, tc.t)
	sb.WriteString(fmt.Sprintf("%d\n", len(tc.arrivals)))
	for i, v := range tc.arrivals {
		if i > 0 {
			sb.WriteByte(' ')
		}
		sb.WriteString(strconv.FormatInt(v, 10))
	}
	sb.WriteByte('\n')
	return sb.String()
}

func randomCase(rng *rand.Rand) testCase {
	ts := int64(rng.Intn(10))
	tf := ts + int64(rng.Intn(20)+10)
	t := int64(rng.Intn(5) + 1)
	n := rng.Intn(8) + 1
	arr := make([]int64, n)
	for i := range arr {
		arr[i] = int64(rng.Intn(int(tf)))
	}
	return testCase{ts: ts, tf: tf, t: t, arrivals: arr}
}

func runProgram(bin, input string) (string, error) {
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
	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("runtime error: %v\n%s", err, out.String())
	}
	return strings.TrimSpace(out.String()), nil
}

func runCase(bin string, tc testCase) error {
	in := tc.input()
	expected := solve(tc.ts, tc.tf, tc.t, append([]int64(nil), tc.arrivals...))
	got, err := runProgram(bin, in)
	if err != nil {
		return err
	}
	if strings.TrimSpace(got) != strings.TrimSpace(expected) {
		return fmt.Errorf("expected %s got %s", expected, got)
	}
	return nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierB.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	cases := []testCase{randomCase(rng)}
	for len(cases) < 100 {
		cases = append(cases, randomCase(rng))
	}
	for i, tc := range cases {
		if err := runCase(bin, tc); err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:\n%s", i+1, err, tc.input())
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
