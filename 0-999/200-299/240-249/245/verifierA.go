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

type command struct {
	t int
	x int
	y int
}

func expected(commands []command) string {
	succA, succB := 0, 0
	cntA, cntB := 0, 0
	for _, c := range commands {
		if c.t == 1 {
			succA += c.x
			cntA += 10
		} else {
			succB += c.x
			cntB += 10
		}
	}
	var res strings.Builder
	if succA*2 >= cntA {
		res.WriteString("LIVE\n")
	} else {
		res.WriteString("DEAD\n")
	}
	if succB*2 >= cntB {
		res.WriteString("LIVE")
	} else {
		res.WriteString("DEAD")
	}
	return res.String()
}

func generateCase(rng *rand.Rand) (string, string) {
	n := rng.Intn(20) + 2
	cmds := make([]command, n)
	hasA, hasB := false, false
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("%d\n", n))
	for i := 0; i < n; i++ {
		t := 1 + rng.Intn(2)
		x := rng.Intn(11)
		y := 10 - x
		cmds[i] = command{t, x, y}
		if t == 1 {
			hasA = true
		} else {
			hasB = true
		}
		sb.WriteString(fmt.Sprintf("%d %d %d\n", t, x, y))
	}
	// ensure at least one of each type
	if !hasA {
		cmds[0].t = 1
	}
	if !hasB {
		cmds[0].t = 2
	}
	if !hasA || !hasB {
		sb.Reset()
		sb.WriteString(fmt.Sprintf("%d\n", n))
		for _, c := range cmds {
			sb.WriteString(fmt.Sprintf("%d %d %d\n", c.t, c.x, c.y))
		}
	}
	return sb.String(), expected(cmds)
}

func runCase(bin, input, expected string) error {
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
		return fmt.Errorf("runtime error: %v\n%s", err, out.String())
	}
	got := strings.TrimSpace(out.String())
	if got != expected {
		return fmt.Errorf("expected %s got %s", expected, got)
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
		in, exp := generateCase(rng)
		if err := runCase(bin, in, exp); err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:%s", i+1, err, in)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
