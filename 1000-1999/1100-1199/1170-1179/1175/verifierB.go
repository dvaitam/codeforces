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

func solve(commands []string) string {
	const limit = uint64(1) << 32
	stack := []uint64{1}
	var x uint64
	for _, cmd := range commands {
		if cmd == "add" {
			x += stack[len(stack)-1]
			if x >= limit {
				return "OVERFLOW!!!"
			}
		} else if strings.HasPrefix(cmd, "for") {
			n, _ := strconv.ParseUint(strings.Fields(cmd)[1], 10, 64)
			top := stack[len(stack)-1]
			prod := top * n
			if top >= limit || prod >= limit {
				stack = append(stack, limit)
			} else {
				stack = append(stack, prod)
			}
		} else { // end
			stack = stack[:len(stack)-1]
		}
	}
	return fmt.Sprint(x)
}

func generateCase(rng *rand.Rand) (string, string) {
	l := rng.Intn(10) + 1
	cmds := make([]string, 0, l)
	open := 0
	for len(cmds) < l {
		choice := rng.Intn(3)
		if choice == 0 && open < 3 && len(cmds)+1 < l { // open loop
			n := rng.Intn(5) + 1
			cmds = append(cmds, fmt.Sprintf("for %d", n))
			open++
		} else if choice == 1 && open > 0 { // close loop
			cmds = append(cmds, "end")
			open--
		} else { // add
			cmds = append(cmds, "add")
		}
	}
	for open > 0 {
		cmds = append(cmds, "end")
		open--
	}
	input := fmt.Sprintf("%d\n%s\n", len(cmds), strings.Join(cmds, "\n"))
	expected := solve(cmds)
	return input, expected + "\n"
}

func runCase(bin string, in, exp string) error {
	cmd := exec.Command(bin)
	cmd.Stdin = strings.NewReader(in)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("runtime error: %v\n%s", err, out.String())
	}
	if strings.TrimSpace(out.String()) != strings.TrimSpace(exp) {
		return fmt.Errorf("expected:\n%s\ngot:\n%s", exp, out.String())
	}
	return nil
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
		if err := runCase(bin, in, exp); err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:\n%s", i+1, err, in)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
