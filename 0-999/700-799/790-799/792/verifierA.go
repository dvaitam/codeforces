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

func solveCase(arr []int) (int, int) {
	sort.Ints(arr)
	minDiff := arr[1] - arr[0]
	cnt := 1
	for i := 1; i < len(arr)-1; i++ {
		diff := arr[i+1] - arr[i]
		if diff < minDiff {
			minDiff = diff
			cnt = 1
		} else if diff == minDiff {
			cnt++
		}
	}
	return minDiff, cnt
}

func genCase(rng *rand.Rand) (string, int, int) {
	n := rng.Intn(18) + 2 // n between 2 and 19
	m := make(map[int]struct{})
	arr := make([]int, n)
	for i := 0; i < n; i++ {
		x := rng.Intn(2001) - 1000
		for {
			if _, ok := m[x]; !ok {
				break
			}
			x = rng.Intn(2001) - 1000
		}
		arr[i] = x
		m[x] = struct{}{}
	}
	minDiff, cnt := solveCase(append([]int(nil), arr...))
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("%d\n", n))
	for i, v := range arr {
		if i > 0 {
			sb.WriteByte(' ')
		}
		sb.WriteString(strconv.Itoa(v))
	}
	sb.WriteByte('\n')
	return sb.String(), minDiff, cnt
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

func runCase(bin string, in string, expMin, expCnt int) error {
	out, err := runCandidate(bin, []byte(in))
	if err != nil {
		return err
	}
	tokens := strings.Fields(out)
	if len(tokens) < 2 {
		return fmt.Errorf("expected two integers, got %q", out)
	}
	gotMin, err1 := strconv.Atoi(tokens[0])
	gotCnt, err2 := strconv.Atoi(tokens[1])
	if err1 != nil || err2 != nil {
		return fmt.Errorf("invalid integers in output: %q", out)
	}
	if gotMin != expMin || gotCnt != expCnt {
		return fmt.Errorf("expected %d %d got %d %d", expMin, expCnt, gotMin, gotCnt)
	}
	return nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierA.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < 100; i++ {
		in, minDiff, cnt := genCase(rng)
		if err := runCase(bin, in, minDiff, cnt); err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:\n%s", i+1, err, in)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
