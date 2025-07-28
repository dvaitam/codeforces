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

func solveB(n, x, y int) ([]int, bool) {
	if (x > 0 && y > 0) || (x == 0 && y == 0) {
		return nil, false
	}
	if x == 0 {
		x, y = y, x
	}
	if x == 0 || (n-1)%x != 0 {
		return nil, false
	}
	res := make([]int, n-1)
	a, b, cnt := 1, 2, 0
	for i := 0; i < n-1; i++ {
		if cnt < x {
			res[i] = a
			b++
			cnt++
		} else {
			res[i] = b
			a = b + 1
			a, b = b, a
			cnt = 1
		}
	}
	return res, true
}

func checkSequence(n, x, y int, seq []int) error {
	if len(seq) != n-1 {
		return fmt.Errorf("expected %d numbers, got %d", n-1, len(seq))
	}
	wins := make([]int, n+1)
	prevWinner := seq[0]
	if prevWinner != 1 && prevWinner != 2 {
		return fmt.Errorf("game 1 winner must be 1 or 2")
	}
	wins[prevWinner]++
	for i := 1; i < n-1; i++ {
		nextPlayer := i + 2
		w := seq[i]
		if w != prevWinner && w != nextPlayer {
			return fmt.Errorf("game %d invalid winner", i+1)
		}
		wins[w]++
		prevWinner = w
	}
	for i := 1; i <= n; i++ {
		if wins[i] != x && wins[i] != y {
			return fmt.Errorf("player %d wins %d not %d/%d", i, wins[i], x, y)
		}
	}
	return nil
}

func genCase(rng *rand.Rand) (string, int, int, int) {
	n := rng.Intn(10) + 2
	x := rng.Intn(n)
	y := rng.Intn(n)
	input := fmt.Sprintf("1\n%d %d %d\n", n, x, y)
	return input, n, x, y
}

func run(bin, input string) (string, error) {
	var cmd *exec.Cmd
	if strings.HasSuffix(bin, ".go") {
		cmd = exec.Command("go", "run", bin)
	} else if strings.HasSuffix(bin, ".py") {
		cmd = exec.Command("python3", bin)
	} else {
		cmd = exec.Command(bin)
	}
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	var errBuf bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &errBuf
	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("runtime error: %v\n%s", err, errBuf.String())
	}
	return strings.TrimSpace(out.String()), nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierB.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < 100; i++ {
		input, n, x, y := genCase(rng)
		out, err := run(bin, input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:\n%s", i+1, err, input)
			os.Exit(1)
		}
		fields := strings.Fields(out)
		if len(fields) == 1 && fields[0] == "-1" {
			if _, ok := solveB(n, x, y); ok {
				fmt.Fprintf(os.Stderr, "case %d failed: output -1 but solution exists\ninput:\n%s", i+1, input)
				os.Exit(1)
			}
			continue
		}
		seq := make([]int, len(fields))
		for j, f := range fields {
			v, err := strconv.Atoi(f)
			if err != nil {
				fmt.Fprintf(os.Stderr, "case %d failed: invalid integer %q\n", i+1, f)
				os.Exit(1)
			}
			seq[j] = v
		}
		if _, ok := solveB(n, x, y); !ok {
			fmt.Fprintf(os.Stderr, "case %d failed: solution impossible but output provided\ninput:\n%s", i+1, input)
			os.Exit(1)
		}
		if err := checkSequence(n, x, y, seq); err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:\n%soutput:%s", i+1, err, input, out)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
