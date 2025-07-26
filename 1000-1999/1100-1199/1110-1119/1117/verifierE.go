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

// applyOps applies swap operations to string s
func applyOps(s string, ops [][2]int) string {
	b := []byte(s)
	for _, op := range ops {
		i := op[0] - 1
		j := op[1] - 1
		b[i], b[j] = b[j], b[i]
	}
	return string(b)
}

type caseE struct {
	s   string
	ops [][2]int
}

func generateCase(rng *rand.Rand) caseE {
	n := rng.Intn(6) + 1
	letters := []byte{'a', 'b', 'c', 'd', 'e'}
	b := make([]byte, n)
	for i := 0; i < n; i++ {
		b[i] = letters[rng.Intn(len(letters))]
	}
	k := rng.Intn(n + 1)
	ops := make([][2]int, k)
	for i := 0; i < k; i++ {
		a := rng.Intn(n) + 1
		b2 := rng.Intn(n) + 1
		ops[i] = [2]int{a, b2}
	}
	return caseE{string(b), ops}
}

func runCase(bin string, c caseE) error {
	t := applyOps(c.s, c.ops)
	var cmd *exec.Cmd
	if strings.HasSuffix(bin, ".go") {
		cmd = exec.Command("go", "run", bin)
	} else {
		cmd = exec.Command(bin)
	}
	stdin, err := cmd.StdinPipe()
	if err != nil {
		return err
	}
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		return err
	}
	cmd.Stderr = os.Stderr
	if err := cmd.Start(); err != nil {
		return err
	}
	inWriter := bufio.NewWriter(stdin)
	outReader := bufio.NewReader(stdout)
	fmt.Fprintln(inWriter, t)
	inWriter.Flush()
	for i := 0; i < 3; i++ {
		line, err := outReader.ReadString('\n')
		if err != nil {
			cmd.Process.Kill()
			return fmt.Errorf("expected query %d: %v", i+1, err)
		}
		line = strings.TrimSpace(line)
		if !strings.HasPrefix(line, "? ") {
			cmd.Process.Kill()
			return fmt.Errorf("expected '? query', got %q", line)
		}
		q := line[2:]
		if len(q) != len(c.s) {
			cmd.Process.Kill()
			return fmt.Errorf("query length mismatch")
		}
		resp := applyOps(q, c.ops)
		fmt.Fprintln(inWriter, resp)
		inWriter.Flush()
	}
	ansLine, err := outReader.ReadString('\n')
	if err != nil {
		cmd.Process.Kill()
		return fmt.Errorf("failed to read answer: %v", err)
	}
	ansLine = strings.TrimSpace(ansLine)
	if !strings.HasPrefix(ansLine, "! ") {
		cmd.Process.Kill()
		return fmt.Errorf("expected final answer, got %q", ansLine)
	}
	ans := ansLine[2:]
	if ans != c.s {
		cmd.Process.Kill()
		return fmt.Errorf("expected %s got %s", c.s, ans)
	}
	cmd.Wait()
	return nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierE.go /path/to/binary")
		os.Exit(1)
	}
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	bin := os.Args[1]
	for i := 0; i < 100; i++ {
		c := generateCase(rng)
		if err := runCase(bin, c); err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\n", i+1, err)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
