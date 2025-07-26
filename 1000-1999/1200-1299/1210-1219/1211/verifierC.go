package main

import (
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"sort"
	"strings"
	"time"
)

type day struct {
	a int64
	b int64
	c int64
}

func runProgram(bin string, input []byte) (string, error) {
	cmd := exec.Command(bin)
	cmd.Stdin = bytes.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	err := cmd.Run()
	return strings.TrimSpace(out.String()), err
}

func expectedC(days []day, k int64) string {
	var minNeed, maxPossible int64
	for _, d := range days {
		minNeed += d.a
		maxPossible += d.b
	}
	if k < minNeed || k > maxPossible {
		return "-1"
	}
	totalCost := int64(0)
	for i := range days {
		totalCost += days[i].a * days[i].c
		days[i].b -= days[i].a
	}
	remaining := k - minNeed
	sort.Slice(days, func(i, j int) bool { return days[i].c < days[j].c })
	for _, d := range days {
		if remaining == 0 {
			break
		}
		take := d.b
		if take > remaining {
			take = remaining
		}
		totalCost += take * d.c
		remaining -= take
	}
	return fmt.Sprintf("%d", totalCost)
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("usage: go run verifierC.go /path/to/binary")
		return
	}
	rand.Seed(time.Now().UnixNano())
	bin := os.Args[1]

	for t := 0; t < 100; t++ {
		n := rand.Intn(5) + 1
		days := make([]day, n)
		var input strings.Builder
		var minNeed, maxPossible int64
		for i := 0; i < n; i++ {
			a := int64(rand.Intn(5))
			b := a + int64(rand.Intn(5))
			c := int64(rand.Intn(5) + 1)
			days[i] = day{a, b, c}
			minNeed += a
			maxPossible += b
		}
		k := minNeed
		if maxPossible > minNeed {
			k += int64(rand.Intn(int(maxPossible - minNeed + 1)))
		}
		input.WriteString(fmt.Sprintf("%d %d\n", n, k))
		for i := 0; i < n; i++ {
			if i > 0 {
				input.WriteByte('\n')
			}
			d := days[i]
			input.WriteString(fmt.Sprintf("%d %d %d", d.a, d.b, d.c))
		}
		input.WriteByte('\n')
		inBytes := []byte(input.String())
		expected := expectedC(append([]day(nil), days...), k)
		out, err := runProgram(bin, inBytes)
		if err != nil || strings.TrimSpace(out) != expected {
			fmt.Println("Test", t+1, "failed")
			fmt.Println("Input:\n", input.String())
			fmt.Println("Expected:", expected)
			fmt.Println("Output:", out)
			if err != nil {
				fmt.Println("Error:", err)
			}
			return
		}
	}
	fmt.Println("All tests passed")
}
