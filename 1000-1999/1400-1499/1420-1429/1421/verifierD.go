package main

import (
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"time"
)

type testCase struct {
	x, y int64
	c    [7]int64
}

func buildCase(x, y int64, c [7]int64) testCase {
	return testCase{x: x, y: y, c: c}
}

func generateCase(rng *rand.Rand) testCase {
	x := rng.Int63n(2_000_000_001) - 1_000_000_000
	y := rng.Int63n(2_000_000_001) - 1_000_000_000
	var c [7]int64
	for i := 1; i <= 6; i++ {
		c[i] = rng.Int63n(1_000_000_000) + 1
	}
	return buildCase(x, y, c)
}

func solve(tc testCase) int64 {
	x := tc.x
	y := tc.y
	c := tc.c
	for it := 0; it < 6; it++ {
		for i := 1; i <= 6; i++ {
			left := i - 1
			if left == 0 {
				left = 6
			}
			right := i + 1
			if right == 7 {
				right = 1
			}
			if c[i] > c[left]+c[right] {
				c[i] = c[left] + c[right]
			}
		}
	}
	var ans int64
	switch {
	case x >= 0 && y >= 0:
		k := x
		if y < k {
			k = y
		}
		ans += k * c[1]
		x -= k
		y -= k
		if x > 0 {
			ans += x * c[6]
		} else if y > 0 {
			ans += y * c[2]
		}
	case x <= 0 && y <= 0:
		xx := -x
		yy := -y
		k := xx
		if yy < k {
			k = yy
		}
		ans += k * c[4]
		xx -= k
		yy -= k
		if xx > 0 {
			ans += xx * c[3]
		} else if yy > 0 {
			ans += yy * c[5]
		}
	case x >= 0 && y <= 0:
		ans += x * c[6]
		ans += -y * c[5]
	case x <= 0 && y >= 0:
		ans += y * c[2]
		ans += -x * c[3]
	}
	return ans
}

func runCase(bin string, tc testCase) error {
	var sb strings.Builder
	sb.WriteString("1\n")
	fmt.Fprintf(&sb, "%d %d\n", tc.x, tc.y)
	for i := 1; i <= 6; i++ {
		fmt.Fprintf(&sb, "%d ", tc.c[i])
	}
	sb.WriteString("\n")
	input := sb.String()
	cmd := exec.Command(bin)
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("runtime error: %v\n%s", err, out.String())
	}
	gotStr := strings.TrimSpace(out.String())
	got, err := strconv.ParseInt(gotStr, 10, 64)
	if err != nil {
		return fmt.Errorf("bad output: %v", err)
	}
	expect := solve(tc)
	if got != expect {
		return fmt.Errorf("expected %d got %d", expect, got)
	}
	return nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierD.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))

	var cases []testCase
	c0 := [7]int64{}
	for i := 1; i <= 6; i++ {
		c0[i] = int64(i)
	}
	cases = append(cases, buildCase(1, 1, c0))
	for i := 0; i < 100; i++ {
		cases = append(cases, generateCase(rng))
	}

	for idx, tc := range cases {
		if err := runCase(bin, tc); err != nil {
			var sb strings.Builder
			sb.WriteString("1\n")
			fmt.Fprintf(&sb, "%d %d\n", tc.x, tc.y)
			for i := 1; i <= 6; i++ {
				fmt.Fprintf(&sb, "%d ", tc.c[i])
			}
			sb.WriteString("\n")
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:\n%s", idx+1, err, sb.String())
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
