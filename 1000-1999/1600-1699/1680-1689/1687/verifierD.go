package main

import (
	"bytes"
	"fmt"
	"math"
	"math/rand"
	"os"
	"os/exec"
	"strings"
)

const MAXK = 2000

func expected(a []int) int {
	n := len(a)
	diff := make([]int, MAXK+2)
	for _, x := range a {
		sMin := int(math.Sqrt(float64(max(0, x-MAXK)))) - 2
		if sMin < 0 {
			sMin = 0
		}
		sMax := int(math.Sqrt(float64(x+MAXK))) + 2
		for s := sMin; s <= sMax; s++ {
			sq := s * s
			l := sq - x
			r := sq + s - x
			if r < 0 || l > MAXK {
				continue
			}
			if l < 0 {
				l = 0
			}
			if r > MAXK {
				r = MAXK
			}
			diff[l]++
			diff[r+1]--
		}
	}
	cnt := 0
	for k := 0; k <= MAXK; k++ {
		cnt += diff[k]
		if cnt == n {
			return k
		}
	}
	return -1
}

func genTest(rng *rand.Rand) []int {
	n := rng.Intn(6) + 1
	a := make([]int, n)
	for i := 0; i < n; i++ {
		a[i] = rng.Intn(50)
	}
	return a
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
		fmt.Println("usage: go run verifierD.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rng := rand.New(rand.NewSource(3))
	for i := 0; i < 100; i++ {
		a := genTest(rng)
		n := len(a)
		var sb strings.Builder
		sb.WriteString(fmt.Sprintf("%d\n", n))
		for j, v := range a {
			if j > 0 {
				sb.WriteByte(' ')
			}
			sb.WriteString(fmt.Sprint(v))
		}
		sb.WriteByte('\n')
		expect := fmt.Sprint(expected(a))
		got, err := runCandidate(bin, sb.String())
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\n", i+1, err)
			os.Exit(1)
		}
		if strings.TrimSpace(got) != expect {
			fmt.Fprintf(os.Stderr, "case %d failed: expected %s got %s\n", i+1, expect, got)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}
