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

func runCandidate(bin, input string) (string, error) {
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
		return "", fmt.Errorf("runtime error: %v\n%s", err, out.String())
	}
	return strings.TrimSpace(out.String()), nil
}

func computeExpected(x int, lines []string) (int, int) {
	used := make([]bool, x)
	for _, line := range lines {
		parts := strings.Fields(line)
		if len(parts) == 0 {
			continue
		}
		if parts[0] == "1" {
			num2, _ := strconv.Atoi(parts[1])
			num1, _ := strconv.Atoi(parts[2])
			if num2 < x {
				used[num2] = true
			}
			if num1 < x {
				used[num1] = true
			}
		} else {
			num, _ := strconv.Atoi(parts[1])
			if num < x {
				used[num] = true
			}
		}
	}
	sumMin, sumMax := 0, 0
	for i := 1; i < x; {
		if used[i] {
			i++
			continue
		}
		j := i
		for j < x && !used[j] {
			j++
		}
		L := j - i
		sumMax += L
		sumMin += (L + 1) / 2
		i = j
	}
	return sumMin, sumMax
}

func generateCase(rng *rand.Rand) (string, string) {
	x := rng.Intn(30) + 2 // 2..31
	maxK := x - 1
	if maxK > 15 {
		maxK = 15
	}
	k := rng.Intn(maxK + 1)
	used := make(map[int]bool)
	lines := make([]string, 0, k)
	for len(lines) < k {
		typ := rng.Intn(2) + 1
		if typ == 1 {
			if x <= 2 {
				continue
			}
			num2 := rng.Intn(x-2) + 1 // ensure num2+1 < x
			if used[num2] || used[num2+1] {
				continue
			}
			used[num2] = true
			used[num2+1] = true
			lines = append(lines, fmt.Sprintf("1 %d %d", num2, num2+1))
		} else {
			num := rng.Intn(x-1) + 1
			if used[num] {
				continue
			}
			used[num] = true
			lines = append(lines, fmt.Sprintf("2 %d", num))
		}
	}

	expMin, expMax := computeExpected(x, lines)

	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("%d %d\n", x, len(lines)))
	for _, l := range lines {
		sb.WriteString(l)
		sb.WriteByte('\n')
	}
	exp := fmt.Sprintf("%d %d", expMin, expMax)
	return sb.String(), exp
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierB.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < 100; i++ {
		in, exp := generateCase(rng)
		out, err := runCandidate(bin, in)
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:\n%s", i+1, err, in)
			os.Exit(1)
		}
		if strings.TrimSpace(out) != exp {
			fmt.Fprintf(os.Stderr, "case %d failed: expected %s got %s\ninput:\n%s", i+1, exp, out, in)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
