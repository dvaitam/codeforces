package main

import (
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"strings"
	"time"
)

func solveD(x1, x2, a, b int64) (string, string) {
	d := x2 - x1
	if d >= a && d <= b {
		return "FIRST", fmt.Sprintf("%d", x2)
	}
	if a <= 0 && b >= 0 {
		return "DRAW", ""
	}
	A := a
	B := b
	if A < 0 {
		A = -A
	}
	if B < 0 {
		B = -B
	}
	if A > B {
		A, B = B, A
	}
	D := d
	if D < 0 {
		D = -D
	}
	c := A + B
	if c != 0 {
		r := D % c
		if r >= A && r <= B {
			move := r
			if d < 0 {
				move = -r
			}
			return "FIRST", fmt.Sprintf("%d", x1+move)
		}
	}
	var intervals [][2]int64
	l1 := d - B
	r1 := d - A
	if l1 <= r1 {
		if l1 < a {
			l1 = a
		}
		if r1 > b {
			r1 = b
		}
		if l1 <= r1 {
			intervals = append(intervals, [2]int64{l1, r1})
		}
	}
	l2 := d + A
	r2 := d + B
	if l2 <= r2 {
		if l2 < a {
			l2 = a
		}
		if r2 > b {
			r2 = b
		}
		if l2 <= r2 {
			intervals = append(intervals, [2]int64{l2, r2})
		}
	}
	covered := int64(0)
	if len(intervals) > 0 {
		for i := 0; i < len(intervals)-1; i++ {
			for j := i + 1; j < len(intervals); j++ {
				if intervals[j][0] < intervals[i][0] {
					intervals[i], intervals[j] = intervals[j], intervals[i]
				}
			}
		}
		cur := intervals[0]
		for i := 1; i < len(intervals); i++ {
			nxt := intervals[i]
			if nxt[0] <= cur[1]+1 {
				if nxt[1] > cur[1] {
					cur[1] = nxt[1]
				}
			} else {
				covered += cur[1] - cur[0] + 1
				cur = nxt
			}
		}
		covered += cur[1] - cur[0] + 1
	}
	total := b - a + 1
	if covered == total {
		return "SECOND", ""
	}
	return "DRAW", ""
}

func genCaseD(rng *rand.Rand) (int64, int64, int64, int64) {
	x1 := int64(rng.Intn(41) - 20)
	x2 := int64(rng.Intn(41) - 20)
	for x2 == x1 {
		x2 = int64(rng.Intn(41) - 20)
	}
	a := int64(rng.Intn(21) - 10)
	b := a + int64(rng.Intn(21))
	if a > b {
		a, b = b, a
	}
	return x1, x2, a, b
}

func runCaseD(bin string, x1, x2, a, b int64) error {
	input := fmt.Sprintf("%d %d %d %d\n", x1, x2, a, b)
	cmd := exec.Command(bin)
	cmd.Stdin = strings.NewReader(input)
	out, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("runtime error: %v", err)
	}
	exp1, exp2 := solveD(x1, x2, a, b)
	parts := strings.Fields(strings.TrimSpace(string(out)))
	if exp1 == "FIRST" {
		if len(parts) != 2 || parts[0] != "FIRST" || parts[1] != exp2 {
			return fmt.Errorf("expected FIRST %s got %q", exp2, strings.Join(parts, " "))
		}
	} else {
		if len(parts) != 1 || parts[0] != exp1 {
			return fmt.Errorf("expected %s got %q", exp1, strings.Join(parts, " "))
		}
	}
	return nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierD.go /path/to/binary")
		return
	}
	bin := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for t := 0; t < 100; t++ {
		x1, x2, a, b := genCaseD(rng)
		if err := runCaseD(bin, x1, x2, a, b); err != nil {
			fmt.Printf("case %d failed: %v\n", t+1, err)
			return
		}
	}
	fmt.Println("All tests passed")
}
