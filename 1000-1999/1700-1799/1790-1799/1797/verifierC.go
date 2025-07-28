package main

import (
	"bufio"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

func abs(x int) int {
	if x < 0 {
		return -x
	}
	return x
}

func runInteractive(bin string, n, m, hx, hy int) error {
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

	// send n m
	fmt.Fprintf(stdin, "%d %d\n", n, m)
	reader := bufio.NewReader(stdout)

	var r1, c1 int
	if _, err := fmt.Fscan(reader, &r1, &c1); err != nil {
		return fmt.Errorf("failed to read query1: %v", err)
	}
	fmt.Fprintf(stdin, "%d\n", abs(r1-hx)+abs(c1-hy))

	var r2, c2 int
	if _, err := fmt.Fscan(reader, &r2, &c2); err != nil {
		return fmt.Errorf("failed to read query2: %v", err)
	}
	fmt.Fprintf(stdin, "%d\n", abs(r2-hx)+abs(c2-hy))

	var fr, fc int
	if _, err := fmt.Fscan(reader, &fr, &fc); err != nil {
		return fmt.Errorf("failed to read final: %v", err)
	}

	stdin.Close()
	err = cmd.Wait()
	if err != nil {
		return fmt.Errorf("binary error: %v", err)
	}

	if fr != hx || fc != hy {
		return fmt.Errorf("expected %d %d got %d %d", hx, hy, fr, fc)
	}
	return nil
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("usage: verifierC.go path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	if b, err := filepath.Abs(bin); err == nil {
		bin = b
	}

	rand.Seed(3)
	const T = 100
	for i := 0; i < T; i++ {
		n := rand.Intn(4) + 2 // 2..5
		m := rand.Intn(4) + 2
		hx := rand.Intn(n) + 1
		hy := rand.Intn(m) + 1
		if err := runInteractive(bin, n, m, hx, hy); err != nil {
			fmt.Printf("test %d failed: %v\n", i+1, err)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
