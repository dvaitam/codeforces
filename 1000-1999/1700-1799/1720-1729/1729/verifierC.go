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

func solveCase(s string) (int, []int) {
	n := len(s)
	startVal := int(s[0] - 'a')
	endVal := int(s[n-1] - 'a')
	minv, maxv := startVal, endVal
	if minv > maxv {
		minv, maxv = maxv, minv
	}
	type pair struct{ val, idx int }
	arr := make([]pair, 0)
	for i := 1; i < n-1; i++ {
		v := int(s[i] - 'a')
		if v >= minv && v <= maxv {
			arr = append(arr, pair{v, i + 1})
		}
	}
	if startVal <= endVal {
		sort.Slice(arr, func(i, j int) bool {
			if arr[i].val == arr[j].val {
				return arr[i].idx < arr[j].idx
			}
			return arr[i].val < arr[j].val
		})
	} else {
		sort.Slice(arr, func(i, j int) bool {
			if arr[i].val == arr[j].val {
				return arr[i].idx < arr[j].idx
			}
			return arr[i].val > arr[j].val
		})
	}
	path := make([]int, 0, len(arr)+2)
	path = append(path, 1)
	for _, p := range arr {
		path = append(path, p.idx)
	}
	path = append(path, n)
	cost := abs(startVal - endVal)
	return cost, path
}

func abs(x int) int {
	if x < 0 {
		return -x
	}
	return x
}

func runCandidate(bin, input string) (string, error) {
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
		fmt.Fprintln(os.Stderr, "usage: go run verifierC.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	cases := []string{"ab", "ba", "az", "za", "abcdef", "fedcba"}
	for len(cases) < 120 {
		l := rng.Intn(10) + 2
		var sb strings.Builder
		for i := 0; i < l; i++ {
			sb.WriteByte(byte('a' + rng.Intn(26)))
		}
		cases = append(cases, sb.String())
	}

	for i, s := range cases {
		cost, path := solveCase(s)
		var sb strings.Builder
		sb.WriteString("1\n")
		sb.WriteString(fmt.Sprintf("%s\n", s))
		input := sb.String()
		expectedOut := fmt.Sprintf("%d %d\n", cost, len(path))
		for j, p := range path {
			if j > 0 {
				expectedOut += " "
			}
			expectedOut += fmt.Sprintf("%d", p)
		}
		expectedOut += "\n"

		out, err := runCandidate(bin, input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:\n%s", i+1, err, input)
			os.Exit(1)
		}
		if strings.TrimSpace(out) != strings.TrimSpace(expectedOut) {
			fmt.Fprintf(os.Stderr, "case %d failed: expected %q got %q\ninput:\n%s", i+1, expectedOut, out, input)
			os.Exit(1)
		}
	}
	fmt.Printf("All %d tests passed\n", len(cases))
}
