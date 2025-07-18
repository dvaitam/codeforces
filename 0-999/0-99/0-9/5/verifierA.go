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
	line         string
	participants int
	add          bool
	remove       bool
}

func generateCase(rng *rand.Rand) (string, int) {
	names := []string{"alice", "bob", "carol", "dave", "eve"}
	used := make(map[string]bool)
	var cmds []string
	participants := 0
	total := 0
	n := rng.Intn(20) + 1
	for i := 0; i < n; i++ {
		choice := rng.Intn(3)
		if participants == 0 {
			choice = 0 // must add
		}
		switch choice {
		case 0: // add
			var name string
			for {
				name = names[rng.Intn(len(names))]
				if !used[name] {
					break
				}
			}
			used[name] = true
			participants++
			cmds = append(cmds, "+"+name)
		case 1: // message
			// pick any participant
			idx := rng.Intn(participants)
			var name string
			cnt := 0
			for n := range used {
				if cnt == idx {
					name = n
					break
				}
				cnt++
			}
			msgLen := rng.Intn(10)
			msg := make([]byte, msgLen)
			for j := 0; j < msgLen; j++ {
				msg[j] = byte('a' + rng.Intn(26))
			}
			cmds = append(cmds, fmt.Sprintf("%s:%s", name, string(msg)))
			total += participants * msgLen
		case 2: // remove
			idx := rng.Intn(participants)
			var name string
			cnt := 0
			for n := range used {
				if cnt == idx {
					name = n
					break
				}
				cnt++
			}
			delete(used, name)
			participants--
			cmds = append(cmds, "-"+name)
		}
	}
	input := strings.Join(cmds, "\n") + "\n"
	return input, total
}

func runCase(bin string, input string, expected int) error {
	cmd := exec.Command(bin)
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("runtime error: %v\n%s", err, out.String())
	}
	var got int
	if _, err := fmt.Sscan(strings.TrimSpace(out.String()), &got); err != nil {
		return fmt.Errorf("failed to parse output: %v", err)
	}
	if got != expected {
		return fmt.Errorf("expected %d got %d", expected, got)
	}
	return nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierA.go /path/to/binary")
		os.Exit(1)
	}
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	bin := os.Args[1]
	for i := 0; i < 100; i++ {
		in, exp := generateCase(rng)
		if err := runCase(bin, in, exp); err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:\n%s", i+1, err, in)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
