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

var daysBefore = [13]int{0, 0, 31, 60, 91, 121, 152, 182, 213, 244, 274, 305, 335}

func parseTS(ts string) int {
	month, _ := strconv.Atoi(ts[5:7])
	day, _ := strconv.Atoi(ts[8:10])
	hour, _ := strconv.Atoi(ts[11:13])
	minute, _ := strconv.Atoi(ts[14:16])
	second, _ := strconv.Atoi(ts[17:19])
	days := daysBefore[month] + day - 1
	return days*86400 + hour*3600 + minute*60 + second
}

func solveF(window, threshold int, lines []string) string {
	var q []int
	for _, rec := range lines {
		ts := rec[:19]
		cur := parseTS(ts)
		q = append(q, cur)
		limit := cur - window + 1
		i := 0
		for i < len(q) && q[i] < limit {
			i++
		}
		if i > 0 {
			q = q[i:]
		}
		if len(q) >= threshold {
			return ts
		}
	}
	return "-1"
}

func randomTS(base int) string {
	sec := base
	day := sec / 86400
	sec %= 86400
	hour := sec / 3600
	sec %= 3600
	minute := sec / 60
	second := sec % 60
	// convert day to month/day
	md := []int{31, 29, 31, 30, 31, 30, 31, 31, 30, 31, 30, 31}
	month := 1
	for day >= md[month-1] {
		day -= md[month-1]
		month++
	}
	return fmt.Sprintf("2012-%02d-%02d %02d:%02d:%02d", month, day+1, hour, minute, second)
}

func generateCase(rng *rand.Rand) (string, string) {
	window := rng.Intn(50) + 1
	threshold := rng.Intn(10) + 1
	linesCount := rng.Intn(30) + threshold
	base := 0
	lines := make([]string, linesCount)
	for i := 0; i < linesCount; i++ {
		base += rng.Intn(5) + 1
		ts := randomTS(base)
		lines[i] = ts + ":A"
	}
	input := fmt.Sprintf("%d %d\n%s\n", window, threshold, strings.Join(lines, "\n"))
	ans := solveF(window, threshold, lines)
	return input, ans
}

func runCase(bin, input, expected string) error {
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
		return fmt.Errorf("runtime error: %v\n%s", err, out.String())
	}
	got := strings.TrimSpace(out.String())
	if got != expected {
		return fmt.Errorf("expected %s got %s", expected, got)
	}
	return nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierF.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < 100; i++ {
		in, exp := generateCase(rng)
		if err := runCase(bin, in, exp); err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:%s", i+1, err, in)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
