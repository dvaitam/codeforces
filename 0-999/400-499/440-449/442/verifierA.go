package main

import (
	"bufio"
	"bytes"
	"fmt"
	"math/bits"
	"math/rand"
	"os"
	"os/exec"
	"strings"
	"time"
)

func runBinary(path, input string) (string, error) {
	var cmd *exec.Cmd
	if strings.HasSuffix(path, ".go") {
		cmd = exec.Command("go", "run", path)
	} else {
		cmd = exec.Command(path)
	}
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("runtime error: %v\n%s", err, out.String())
	}
	return out.String(), nil
}

func colorIndex(c byte) int {
	switch c {
	case 'R':
		return 0
	case 'G':
		return 1
	case 'B':
		return 2
	case 'Y':
		return 3
	case 'W':
		return 4
	}
	return -1
}

func solveA(r *bufio.Reader) string {
	var n int
	fmt.Fscan(r, &n)
	cards := make([]string, n)
	for i := 0; i < n; i++ {
		fmt.Fscan(r, &cards[i])
	}
	present := make([][]bool, 5)
	for i := range present {
		present[i] = make([]bool, 5)
	}
	for _, s := range cards {
		ci := colorIndex(s[0])
		vi := int(s[1] - '1')
		if ci >= 0 && vi >= 0 && ci < 5 && vi < 5 {
			present[ci][vi] = true
		}
	}
	types := make([][2]int, 0, 25)
	for ci := 0; ci < 5; ci++ {
		for vi := 0; vi < 5; vi++ {
			if present[ci][vi] {
				types = append(types, [2]int{ci, vi})
			}
		}
	}
	best := 10
	for mask := 0; mask < (1 << 10); mask++ {
		cnt := bits.OnesCount(uint(mask))
		if cnt >= best {
			continue
		}
		ok := true
		for i := 0; i < len(types) && ok; i++ {
			for j := i + 1; j < len(types); j++ {
				c1, v1 := types[i][0], types[i][1]
				c2, v2 := types[j][0], types[j][1]
				distinguished := false
				for bit := 0; bit < 10; bit++ {
					if mask&(1<<bit) == 0 {
						continue
					}
					if bit < 5 {
						if (c1 == bit) != (c2 == bit) {
							distinguished = true
							break
						}
					} else {
						vb := bit - 5
						if (v1 == vb) != (v2 == vb) {
							distinguished = true
							break
						}
					}
				}
				if !distinguished {
					ok = false
					break
				}
			}
		}
		if ok {
			best = cnt
		}
	}
	return fmt.Sprintf("%d\n", best)
}

func generateCase(rng *rand.Rand) string {
	n := rng.Intn(20) + 1
	colors := []byte{'R', 'G', 'B', 'Y', 'W'}
	var b strings.Builder
	fmt.Fprintf(&b, "%d\n", n)
	for i := 0; i < n; i++ {
		if i > 0 {
			b.WriteByte(' ')
		}
		c := colors[rng.Intn(5)]
		v := rng.Intn(5) + 1
		fmt.Fprintf(&b, "%c%d", c, v)
	}
	b.WriteByte('\n')
	return b.String()
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierA.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	cases := make([]string, 0, 100)
	for i := 0; i < 100; i++ {
		cases = append(cases, generateCase(rng))
	}
	for i, tc := range cases {
		expect := solveA(bufio.NewReader(strings.NewReader(tc)))
		out, err := runBinary(bin, tc)
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:\n%s", i+1, err, tc)
			os.Exit(1)
		}
		if strings.TrimSpace(out) != strings.TrimSpace(expect) {
			fmt.Fprintf(os.Stderr, "case %d failed: expected %s got %sinput:\n%s", i+1, expect, out, tc)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
