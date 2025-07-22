package main

import (
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"strconv"
	"strings"
)

type date struct {
	y, m, d int
}

func isLeap(y int) bool {
	if y%400 == 0 {
		return true
	}
	if y%100 == 0 {
		return false
	}
	return y%4 == 0
}

func daysInMonth(y, m int) int {
	switch m {
	case 1, 3, 5, 7, 8, 10, 12:
		return 31
	case 4, 6, 9, 11:
		return 30
	case 2:
		if isLeap(y) {
			return 29
		}
		return 28
	}
	return 30
}

func countDays(y, m, d int) int64 {
	yy := int64(y - 1)
	leaps := yy/4 - yy/100 + yy/400
	days := yy*365 + leaps
	for mm := 1; mm < m; mm++ {
		days += int64(daysInMonth(y, mm))
	}
	days += int64(d)
	return days
}

func diffDates(a, b date) int64 {
	da := countDays(a.y, a.m, a.d)
	db := countDays(b.y, b.m, b.d)
	if da > db {
		return da - db + 1
	}
	return db - da + 1
}

func formatDate(x date) string {
	return fmt.Sprintf("%04d:%02d:%02d", x.y, x.m, x.d)
}

func runCase(exe string, d1, d2 date) error {
	input := fmt.Sprintf("%s\n%s\n", formatDate(d1), formatDate(d2))
	cmd := exec.Command(exe)
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("runtime error: %v\n%s", err, out.String())
	}
	resStr := strings.TrimSpace(out.String())
	got, err := strconv.ParseInt(resStr, 10, 64)
	if err != nil {
		return fmt.Errorf("invalid output %q", resStr)
	}
	exp := diffDates(d1, d2)
	if got != exp {
		return fmt.Errorf("expected %d got %d", exp, got)
	}
	return nil
}

func randDate(rng *rand.Rand) date {
	y := rng.Intn(2038-1900+1) + 1900
	m := rng.Intn(12) + 1
	d := rng.Intn(daysInMonth(y, m)) + 1
	return date{y, m, d}
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierB.go /path/to/binary")
		os.Exit(1)
	}
	exe := os.Args[1]
	rng := rand.New(rand.NewSource(1))
	tests := make([][2]date, 100)
	for i := 0; i < len(tests); i++ {
		tests[i][0] = randDate(rng)
		tests[i][1] = randDate(rng)
	}
	// add some edge cases
	edge := [][2]date{
		{{1900, 1, 1}, {1900, 1, 1}},
		{{2000, 2, 28}, {2000, 3, 1}},
		{{1999, 12, 31}, {2000, 1, 1}},
		{{2038, 12, 31}, {2038, 12, 31}},
	}
	tests = append(tests, edge...)
	for idx, t := range tests {
		if err := runCase(exe, t[0], t[1]); err != nil {
			fmt.Printf("test %d failed: %v\n", idx+1, err)
			fmt.Printf("input:\n%s\n%s\n", formatDate(t[0]), formatDate(t[1]))
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
