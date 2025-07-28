package main

import (
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"strings"
	"time"
)

func solveD(l, r []int64, k int64) int64 {
	var pref int64
	for i := 0; i < len(l); i++ {
		seg := r[i] - l[i] + 1
		if pref+seg >= k {
			q := l[i] + (k - pref) - 1
			return q + 2
		}
		pref += seg
	}
	return -1
}

func genCaseD(rng *rand.Rand) (int, []int64, []int64, int64) {
	n := rng.Intn(5) + 1
	l := make([]int64, n)
	r := make([]int64, n)
	cur := int64(rng.Intn(5) + 1)
	for i := 0; i < n; i++ {
		lenSeg := int64(rng.Intn(5) + 1)
		l[i] = cur
		r[i] = cur + lenSeg - 1
		cur = r[i] + int64(rng.Intn(3)+2)
	}
	total := int64(0)
	for i := 0; i < n; i++ {
		total += r[i] - l[i] + 1
	}
	k := int64(rng.Intn(int(total)+5) + 1)
	return n, l, r, k
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
	cmd.Stdout = &out
	cmd.Stderr = &out
	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("runtime error: %v\n%s", err, out.String())
	}
	return strings.TrimSpace(out.String()), nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierD.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < 100; i++ {
		n, larr, rarr, k := genCaseD(rng)
		input := fmt.Sprintf("1\n%d %d\n", n, k)
		for j, v := range larr {
			if j > 0 {
				input += " "
			}
			input += fmt.Sprint(v)
		}
		input += "\n"
		for j, v := range rarr {
			if j > 0 {
				input += " "
			}
			input += fmt.Sprint(v)
		}
		input += "\n"
		expect := fmt.Sprintf("%d", solveD(larr, rarr, k))
		out, err := runCandidate(bin, input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\n", i+1, err)
			os.Exit(1)
		}
		if out != expect {
			fmt.Fprintf(os.Stderr, "case %d failed: expected %s got %s\n", i+1, expect, out)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
