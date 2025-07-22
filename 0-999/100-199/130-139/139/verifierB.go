package main

import (
	"context"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"strings"
	"time"
)

func runBinary(bin, input string) (string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()
	var cmd *exec.Cmd
	if strings.HasSuffix(bin, ".go") {
		cmd = exec.CommandContext(ctx, "go", "run", bin)
	} else {
		cmd = exec.CommandContext(ctx, bin)
	}
	cmd.Stdin = strings.NewReader(input)
	out, err := cmd.CombinedOutput()
	if ctx.Err() == context.DeadlineExceeded {
		return "", fmt.Errorf("time limit")
	}
	if err != nil {
		return "", fmt.Errorf("%v: %s", err, out)
	}
	return strings.TrimSpace(string(out)), nil
}

type room struct{ l, w, h int }
type wallpaper struct{ L, W, P int }

func solveB(rooms []room, types []wallpaper) int64 {
	const inf int64 = 1<<63 - 1
	var total int64
	for _, r := range rooms {
		perim := 2 * (r.l + r.w)
		best := inf
		for _, t := range types {
			strips := t.L / r.h
			if strips <= 0 {
				continue
			}
			coverage := strips * t.W
			rolls := (perim + coverage - 1) / coverage
			cost := int64(rolls) * int64(t.P)
			if cost < best {
				best = cost
			}
		}
		total += best
	}
	return total
}

func generateCase(rng *rand.Rand) (string, string) {
	n := rng.Intn(5) + 1
	rooms := make([]room, n)
	for i := 0; i < n; i++ {
		rooms[i] = room{
			l: rng.Intn(10) + 1,
			w: rng.Intn(10) + 1,
			h: rng.Intn(10) + 1,
		}
	}
	m := rng.Intn(5) + 1
	types := make([]wallpaper, m)
	for i := 0; i < m; i++ {
		types[i] = wallpaper{
			L: rng.Intn(15) + 1,
			W: rng.Intn(5) + 1,
			P: rng.Intn(20) + 1,
		}
	}
	// ensure each room can be papered
	for _, r := range rooms {
		ok := false
		for _, t := range types {
			if t.L >= r.h {
				ok = true
				break
			}
		}
		if !ok {
			idx := rng.Intn(m)
			types[idx].L = r.h + rng.Intn(5)
		}
	}

	var sb strings.Builder
	fmt.Fprintf(&sb, "%d\n", n)
	for _, rm := range rooms {
		fmt.Fprintf(&sb, "%d %d %d\n", rm.l, rm.w, rm.h)
	}
	fmt.Fprintf(&sb, "%d\n", m)
	for _, tp := range types {
		fmt.Fprintf(&sb, "%d %d %d\n", tp.L, tp.W, tp.P)
	}

	expected := fmt.Sprintf("%d", solveB(rooms, types))
	return sb.String(), expected
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierB.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < 100; i++ {
		in, exp := generateCase(rng)
		out, err := runBinary(bin, in)
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:\n%s", i+1, err, in)
			os.Exit(1)
		}
		if strings.TrimSpace(out) != strings.TrimSpace(exp) {
			fmt.Fprintf(os.Stderr, "case %d failed: expected %s got %s\ninput:\n%s", i+1, exp, out, in)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
