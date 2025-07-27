package main

import (
	"bufio"
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"strings"
	"time"
)

func expectedFirst(n, k int) int {
	if k%n == 0 {
		return 0
	}
	return 2
}

func genCase(rng *rand.Rand) (string, int, int) {
	n := rng.Intn(10) + 1 // limit to 1..10 for speed
	k := rng.Intn(n*n + 1)
	input := fmt.Sprintf("1\n%d %d\n", n, k)
	return input, n, k
}

func checkMatrix(n, k int, lines []string) error {
	if len(lines) != n {
		return fmt.Errorf("expected %d rows, got %d", n, len(lines))
	}
	ones := 0
	row := make([]int, n)
	col := make([]int, n)
	for i := 0; i < n; i++ {
		if len(lines[i]) != n {
			return fmt.Errorf("row %d length %d", i, len(lines[i]))
		}
		for j := 0; j < n; j++ {
			c := lines[i][j]
			if c == '1' {
				ones++
				row[i]++
				col[j]++
			} else if c != '0' {
				return fmt.Errorf("invalid char %c at %d,%d", c, i, j)
			}
		}
	}
	if ones != k {
		return fmt.Errorf("expected %d ones got %d", k, ones)
	}
	minRow, maxRow := row[0], row[0]
	minCol, maxCol := col[0], col[0]
	for i := 1; i < n; i++ {
		if row[i] < minRow {
			minRow = row[i]
		}
		if row[i] > maxRow {
			maxRow = row[i]
		}
		if col[i] < minCol {
			minCol = col[i]
		}
		if col[i] > maxCol {
			maxCol = col[i]
		}
	}
	if maxRow-minRow > 1 {
		return fmt.Errorf("row diff too big")
	}
	if maxCol-minCol > 1 {
		return fmt.Errorf("column diff too big")
	}
	return nil
}

func runCase(bin string, input string, n, k int) error {
	cmd := exec.Command(bin)
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("runtime error: %v\n%s", err, out.String())
	}
	scanner := bufio.NewScanner(strings.NewReader(out.String()))
	if !scanner.Scan() {
		return fmt.Errorf("missing first line")
	}
	var first int
	if _, err := fmt.Sscan(scanner.Text(), &first); err != nil {
		return fmt.Errorf("bad first line: %v", err)
	}
	expFirst := expectedFirst(n, k)
	if first != expFirst {
		return fmt.Errorf("expected first line %d got %d", expFirst, first)
	}
	matrix := make([]string, 0, n)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line != "" {
			matrix = append(matrix, line)
		}
	}
	if err := scanner.Err(); err != nil {
		return err
	}
	return checkMatrix(n, k, matrix)
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierD.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < 100; i++ {
		in, n, k := genCase(rng)
		if err := runCase(bin, in, n, k); err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:\n%s", i+1, err, in)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
