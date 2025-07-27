package main

import (
	"bufio"
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"strings"
	"time"
)

const mod = 1000000007

func solveD(input string) string {
	reader := bufio.NewReader(strings.NewReader(input))
	var n, k int
	var L int
	fmt.Fscan(reader, &n, &k, &L)
	xs := make([]int, n)
	ys := make([]int, n)
	cs := make([]int, n)
	for i := 0; i < n; i++ {
		fmt.Fscan(reader, &xs[i], &ys[i], &cs[i])
		cs[i]--
	}
	var ans int64
	for x1 := 0; x1 < L; x1++ {
		for x2 := x1 + 1; x2 <= L; x2++ {
			for y1 := 0; y1 < L; y1++ {
				for y2 := y1 + 1; y2 <= L; y2++ {
					seen := make([]bool, k)
					cnt := 0
					for i := 0; i < n; i++ {
						if xs[i] >= x1 && xs[i] < x2 && ys[i] >= y1 && ys[i] < y2 {
							if !seen[cs[i]] {
								seen[cs[i]] = true
								cnt++
							}
						}
					}
					if cnt == k {
						ans++
					}
				}
			}
		}
	}
	ans %= mod
	return fmt.Sprint(ans)
}

func genTest(rng *rand.Rand) string {
	k := rng.Intn(3) + 1
	n := k + rng.Intn(3)
	if n < k {
		n = k
	}
	L := rng.Intn(4) + 2
	var buf strings.Builder
	fmt.Fprintf(&buf, "%d %d %d\n", n, k, L)
	for i := 0; i < n; i++ {
		x := rng.Intn(L)
		y := rng.Intn(L)
		c := rng.Intn(k) + 1
		fmt.Fprintf(&buf, "%d %d %d\n", x, y, c)
	}
	return buf.String()
}

func run(bin, input string) (string, error) {
	var cmd *exec.Cmd
	if strings.HasSuffix(bin, ".go") {
		cmd = exec.Command("go", "run", bin)
	} else {
		cmd = exec.Command(bin)
	}
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	var stderr bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &stderr
	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("runtime error: %v\n%s", err, stderr.String())
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
		input := genTest(rng)
		expect := solveD(input)
		got, err := run(bin, input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d: %v\n", i+1, err)
			os.Exit(1)
		}
		if strings.TrimSpace(got) != expect {
			fmt.Fprintf(os.Stderr, "case %d failed\ninput:\n%sexpected:\n%s\ngot:\n%s\n", i+1, input, expect, got)
			os.Exit(1)
		}
	}
	fmt.Println("all tests passed")
}
