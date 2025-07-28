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

const N = 120
const maxSteps = 300

type slime struct{ x, y int }

func precompute() []map[[2]int]bool {
	belts := make([][]int, N)
	for i := range belts {
		belts[i] = make([]int, N)
	}
	slimes := []slime{{0, 0}}
	history := make([]map[[2]int]bool, maxSteps+1)
	for t := 0; t <= maxSteps; t++ {
		grid := make(map[[2]int]bool)
		for _, s := range slimes {
			if s.x < N && s.y < N {
				grid[[2]int{s.x, s.y}] = true
			}
		}
		history[t] = grid
		if t == maxSteps {
			break
		}
		var next []slime
		for _, s := range slimes {
			nx, ny := s.x, s.y
			if belts[nx][ny] == 0 {
				ny++
			} else {
				nx++
			}
			if nx < N && ny < N {
				next = append(next, slime{nx, ny})
			}
		}
		for _, s := range slimes {
			belts[s.x][s.y] ^= 1
		}
		next = append(next, slime{0, 0})
		slimes = next
	}
	return history
}

var history = precompute()

func genCase(rng *rand.Rand) (string, []string) {
	q := rng.Intn(5) + 1
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("%d\n", q))
	answers := make([]string, q)
	for i := 0; i < q; i++ {
		t := rng.Intn(maxSteps + 1)
		x := rng.Intn(N)
		y := rng.Intn(N)
		sb.WriteString(fmt.Sprintf("%d %d %d\n", t, x, y))
		if history[t][[2]int{x, y}] {
			answers[i] = "YES"
		} else {
			answers[i] = "NO"
		}
	}
	return sb.String(), answers
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
		fmt.Println("usage: go run verifierE.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < 100; i++ {
		input, answers := genCase(rng)
		out, err := run(bin, input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:\n%s", i+1, err, input)
			os.Exit(1)
		}
		fields := strings.Fields(out)
		if len(fields) != len(answers) {
			fmt.Fprintf(os.Stderr, "case %d failed: expected %d lines got %d\ninput:\n%s", i+1, len(answers), len(fields), input)
			os.Exit(1)
		}
		for j, ans := range answers {
			if strings.ToUpper(fields[j]) != ans {
				fmt.Fprintf(os.Stderr, "case %d failed: expected %s got %s\ninput:\n%s", i+1, ans, fields[j], input)
				os.Exit(1)
			}
		}
	}
	fmt.Println("All tests passed")
}
