package main

import (
	"bufio"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"strings"
	"time"
)

func runGame(bin string, secret int64) error {
	var cmd *exec.Cmd
	if strings.HasSuffix(bin, ".go") {
		cmd = exec.Command("go", "run", bin)
	} else {
		cmd = exec.Command(bin)
	}
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		return err
	}
	stdin, err := cmd.StdinPipe()
	if err != nil {
		return err
	}
	cmd.Stderr = os.Stderr
	if err := cmd.Start(); err != nil {
		return err
	}
	r := bufio.NewReader(stdout)
	w := bufio.NewWriter(stdin)
	queries := 0
	for {
		var token string
		if _, err := fmt.Fscan(r, &token); err != nil {
			cmd.Process.Kill()
			return fmt.Errorf("read error: %v", err)
		}
		if token == "?" {
			var x, y int64
			if _, err := fmt.Fscan(r, &x, &y); err != nil {
				cmd.Process.Kill()
				return fmt.Errorf("bad query: %v", err)
			}
			queries++
			if queries > 60 {
				cmd.Process.Kill()
				return fmt.Errorf("too many queries")
			}
			resp := "y"
			if x%secret >= y%secret {
				resp = "x"
			}
			fmt.Fprintln(w, resp)
			w.Flush()
		} else if token == "!" {
			var ans int64
			if _, err := fmt.Fscan(r, &ans); err != nil {
				cmd.Process.Kill()
				return fmt.Errorf("bad answer: %v", err)
			}
			if ans != secret {
				cmd.Process.Kill()
				return fmt.Errorf("wrong answer: got %d expected %d", ans, secret)
			}
			stdin.Close()
			err := cmd.Wait()
			return err
		} else {
			cmd.Process.Kill()
			return fmt.Errorf("unexpected token %s", token)
		}
	}
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierB.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	tests := []int64{1, 2, 3}
	for len(tests) < 100 {
		tests = append(tests, rng.Int63n(1_000_000_000)+1)
	}
	for i, secret := range tests {
		if err := runGame(bin, secret); err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\n", i+1, err)
			os.Exit(1)
		}
	}
	fmt.Printf("All %d tests passed\n", len(tests))
}
