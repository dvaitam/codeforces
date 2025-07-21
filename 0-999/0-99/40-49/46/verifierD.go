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

type car struct {
	id int
	x  int
	l  int
}

type op struct {
	typ int
	val int
}

func simulate(L, b, f int, ops []op) []int {
	parked := make([]car, 0)
	info := make(map[int]car)
	out := make([]int, 0)
	reqID := 1
	for _, op := range ops {
		if op.typ == 1 {
			length := op.val
			placed := -1
			if len(parked) == 0 {
				if length <= L {
					placed = 0
				}
			} else {
				first := parked[0]
				if first.x-f-length >= 0 {
					placed = 0
				}
			}
			if placed < 0 {
				for i := 0; i+1 < len(parked); i++ {
					cur := parked[i]
					next := parked[i+1]
					low := cur.x + cur.l + b
					high := next.x - f - length
					if high >= low {
						placed = low
						break
					}
				}
			}
			if placed < 0 && len(parked) > 0 {
				last := parked[len(parked)-1]
				low := last.x + last.l + b
				if low+length <= L {
					placed = low
				}
			}
			if placed >= 0 {
				c := car{id: reqID, x: placed, l: length}
				idx := 0
				for idx < len(parked) && parked[idx].x < placed {
					idx++
				}
				parked = append(parked, car{})
				copy(parked[idx+1:], parked[idx:])
				parked[idx] = c
				info[reqID] = c
			} else {
				info[reqID] = car{id: reqID, x: -1, l: length}
			}
			out = append(out, placed)
			reqID++
		} else {
			rid := op.val
			c, ok := info[rid]
			if ok && c.x >= 0 {
				for i, pc := range parked {
					if pc.id == rid {
						parked = append(parked[:i], parked[i+1:]...)
						break
					}
				}
				info[rid] = car{id: rid, x: -1, l: c.l}
			}
		}
	}
	return out
}

func generateCase(rng *rand.Rand) (string, []int) {
	L := rng.Intn(91) + 10
	b := rng.Intn(5) + 1
	f := rng.Intn(5) + 1
	n := rng.Intn(5) + 1
	ops := make([]op, n)
	nextID := 1
	for i := 0; i < n; i++ {
		if rng.Intn(2) == 0 || nextID == 1 {
			ops[i] = op{1, rng.Intn(9) + 1}
			nextID++
		} else {
			ops[i] = op{2, rng.Intn(nextID-1) + 1}
		}
	}
	sb := strings.Builder{}
	fmt.Fprintf(&sb, "%d %d %d\n%d\n", L, b, f, n)
	for _, op := range ops {
		fmt.Fprintf(&sb, "%d %d\n", op.typ, op.val)
	}
	expect := simulate(L, b, f, ops)
	return sb.String(), expect
}

func runCase(exe, input string, expect []int) error {
	cmd := exec.Command(exe)
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("runtime error: %v\n%s", err, out.String())
	}
	fields := strings.Fields(out.String())
	if len(fields) != len(expect) {
		return fmt.Errorf("expected %d outputs got %d", len(expect), len(fields))
	}
	for i, f := range fields {
		var v int
		if _, err := fmt.Sscan(f, &v); err != nil {
			return fmt.Errorf("bad output %q", f)
		}
		if v != expect[i] {
			return fmt.Errorf("output %d expected %d got %d", i+1, expect[i], v)
		}
	}
	return nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierD.go /path/to/binary")
		os.Exit(1)
	}
	exe := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < 100; i++ {
		in, exp := generateCase(rng)
		if err := runCase(exe, in, exp); err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:\n%s", i+1, err, in)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
