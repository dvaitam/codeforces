package main

import (
	"bufio"
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"sort"
	"strings"
	"time"
)

func runCandidate(bin, input string) (string, error) {
	var cmd *exec.Cmd
	if strings.HasSuffix(bin, ".go") {
		cmd = exec.Command("go", "run", bin)
	} else if strings.HasSuffix(bin, ".py") {
		cmd = exec.Command("python3", bin)
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

func solve(input string) string {
	in := bufio.NewReader(strings.NewReader(input))
	var t int
	fmt.Fscan(in, &t)
	var sb strings.Builder
	for ; t > 0; t-- {
		var w, h int64
		fmt.Fscan(in, &w, &h)
		ans := int64(0)
		for side := 0; side < 4; side++ {
			var k int
			fmt.Fscan(in, &k)
			var first, last int64
			for i := 0; i < k; i++ {
				var v int64
				fmt.Fscan(in, &v)
				if i == 0 {
					first = v
				}
				if i == k-1 {
					last = v
				}
			}
			val := (last - first)
			if side < 2 {
				val *= h
			} else {
				val *= w
			}
			if val > ans {
				ans = val
			}
		}
		sb.WriteString(fmt.Sprint(ans))
		if t > 1 {
			sb.WriteByte('\n')
		}
	}
	return sb.String()
}

func genLine(rng *rand.Rand, limit int64) (string, int64, int64) {
	k := rng.Intn(4) + 2
	arr := make([]int64, k)
	for i := range arr {
		arr[i] = rng.Int63n(limit + 1)
	}
	sort.Slice(arr, func(i, j int) bool { return arr[i] < arr[j] })
	var sb strings.Builder
	fmt.Fprintf(&sb, "%d\n", k)
	for i, v := range arr {
		if i > 0 {
			sb.WriteByte(' ')
		}
		fmt.Fprintf(&sb, "%d", v)
	}
	sb.WriteByte('\n')
	return sb.String(), arr[0], arr[len(arr)-1]
}

func generateCase(rng *rand.Rand) string {
	w := rng.Int63n(1000) + 1
	h := rng.Int63n(1000) + 1
	var sb strings.Builder
	fmt.Fprintf(&sb, "1\n%d %d\n", w, h)
	line1, _, _ := genLine(rng, w)
	sb.WriteString(line1)
	line2, _, _ := genLine(rng, w)
	sb.WriteString(line2)
	line3, _, _ := genLine(rng, h)
	sb.WriteString(line3)
	line4, _, _ := genLine(rng, h)
	sb.WriteString(line4)
	return sb.String()
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierB.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 1; i <= 100; i++ {
		input := generateCase(rng)
		expect := solve(input)
		got, err := runCandidate(bin, input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:\n%s", i, err, input)
			os.Exit(1)
		}
		if strings.TrimSpace(got) != strings.TrimSpace(expect) {
			fmt.Fprintf(os.Stderr, "case %d failed: expected %s got %s\ninput:\n%s", i, expect, got, input)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
