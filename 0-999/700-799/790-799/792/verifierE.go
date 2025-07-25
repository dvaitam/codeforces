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

func solveCaseE(a []int64) int64 {
	candMap := make(map[int64]struct{})
	for _, x := range a {
		for j := int64(1); j*j <= x; j++ {
			candMap[j] = struct{}{}
			if j-1 > 0 {
				candMap[j-1] = struct{}{}
			}
			q := x / j
			candMap[q] = struct{}{}
			if q-1 > 0 {
				candMap[q-1] = struct{}{}
			}
		}
	}
	cand := make([]int64, 0, len(candMap))
	for k := range candMap {
		if k > 0 {
			cand = append(cand, k)
		}
	}
	sort.Slice(cand, func(i, j int) bool { return cand[i] < cand[j] })
	best := int64(1<<63 - 1)
	for _, k := range cand {
		total := int64(0)
		feasible := true
		for _, x := range a {
			t := (x + k) / (k + 1)
			maxSet := x / k
			if t > maxSet {
				feasible = false
				break
			}
			total += t
			if total >= best {
				break
			}
		}
		if feasible && total < best {
			best = total
		}
	}
	return best
}

func genCaseE(rng *rand.Rand) (string, int64) {
	n := rng.Intn(5) + 1
	arr := make([]int64, n)
	for i := 0; i < n; i++ {
		arr[i] = int64(rng.Intn(20) + 1)
	}
	best := solveCaseE(arr)
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("%d\n", n))
	for i, v := range arr {
		if i > 0 {
			sb.WriteByte(' ')
		}
		sb.WriteString(strconv.FormatInt(v, 10))
	}
	sb.WriteByte('\n')
	return sb.String(), best
}

func runCandidate(bin string, input []byte) (string, error) {
	var cmd *exec.Cmd
	if strings.HasSuffix(bin, ".go") {
		cmd = exec.Command("go", "run", bin)
	} else {
		cmd = exec.Command(bin)
	}
	cmd.Stdin = bytes.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("runtime error: %v\n%s", err, out.String())
	}
	return strings.TrimSpace(out.String()), nil
}

func runCaseE(bin string, in string, exp int64) error {
	out, err := runCandidate(bin, []byte(in))
	if err != nil {
		return err
	}
	val, err := strconv.ParseInt(strings.TrimSpace(out), 10, 64)
	if err != nil {
		return fmt.Errorf("invalid output %q", out)
	}
	if val != exp {
		return fmt.Errorf("expected %d got %d", exp, val)
	}
	return nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierE.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < 100; i++ {
		in, exp := genCaseE(rng)
		if err := runCaseE(bin, in, exp); err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:\n%s", i+1, err, in)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
