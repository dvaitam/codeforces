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

// bruteForce computes the optimal efficiency and returns it.
func bruteForce(vids [][2]int, chans []struct{ a, b, c int }) int64 {
	var best int64
	for _, v := range vids {
		for _, ch := range chans {
			lo := v[0]
			if ch.a > lo {
				lo = ch.a
			}
			hi := v[1]
			if ch.b < hi {
				hi = ch.b
			}
			if hi > lo {
				eff := int64(hi-lo) * int64(ch.c)
				if eff > best {
					best = eff
				}
			}
		}
	}
	return best
}

func verify(expEff int64, got string, vids [][2]int, chans []struct{ a, b, c int }) error {
	gotLines := strings.Fields(got)
	if len(gotLines) == 0 {
		return fmt.Errorf("empty output")
	}
	var gotEff int64
	if _, err := fmt.Sscan(gotLines[0], &gotEff); err != nil {
		return fmt.Errorf("cannot parse efficiency: %v", err)
	}
	if gotEff != expEff {
		return fmt.Errorf("expected %d got %d", expEff, gotEff)
	}
	if gotEff == 0 {
		return nil
	}
	// Verify the claimed (vi, cj) actually achieves gotEff
	if len(gotLines) < 3 {
		return fmt.Errorf("missing video/channel line")
	}
	var vi, cj int
	if _, err := fmt.Sscan(gotLines[1], &vi); err != nil {
		return fmt.Errorf("cannot parse video index: %v", err)
	}
	if _, err := fmt.Sscan(gotLines[2], &cj); err != nil {
		return fmt.Errorf("cannot parse channel index: %v", err)
	}
	if vi < 1 || vi > len(vids) {
		return fmt.Errorf("video index %d out of range", vi)
	}
	if cj < 1 || cj > len(chans) {
		return fmt.Errorf("channel index %d out of range", cj)
	}
	v := vids[vi-1]
	ch := chans[cj-1]
	lo := v[0]
	if ch.a > lo {
		lo = ch.a
	}
	hi := v[1]
	if ch.b < hi {
		hi = ch.b
	}
	if hi <= lo {
		return fmt.Errorf("video %d and channel %d do not intersect", vi, cj)
	}
	actual := int64(hi-lo) * int64(ch.c)
	if actual != gotEff {
		return fmt.Errorf("claimed efficiency %d but pair (%d,%d) achieves %d", gotEff, vi, cj, actual)
	}
	return nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierA.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < 100; i++ {
		n := rng.Intn(5) + 1
		m := rng.Intn(5) + 1
		vids := make([][2]int, n)
		for j := 0; j < n; j++ {
			l := rng.Intn(10)
			r := l + rng.Intn(5)
			vids[j] = [2]int{l, r}
		}
		chans := make([]struct{ a, b, c int }, m)
		for j := 0; j < m; j++ {
			a := rng.Intn(10)
			b := a + rng.Intn(5)
			c := rng.Intn(5) + 1
			chans[j] = struct{ a, b, c int }{a, b, c}
		}
		var sb strings.Builder
		sb.WriteString(fmt.Sprintf("%d %d\n", n, m))
		for j := 0; j < n; j++ {
			sb.WriteString(fmt.Sprintf("%d %d\n", vids[j][0], vids[j][1]))
		}
		for j := 0; j < m; j++ {
			sb.WriteString(fmt.Sprintf("%d %d %d\n", chans[j].a, chans[j].b, chans[j].c))
		}
		input := sb.String()
		expEff := bruteForce(vids, chans)
		got, err := run(bin, input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:\n%s", i+1, err, input)
			os.Exit(1)
		}
		if err := verify(expEff, got, vids, chans); err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ngot %s\ninput:\n%s", i+1, err, got, input)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
