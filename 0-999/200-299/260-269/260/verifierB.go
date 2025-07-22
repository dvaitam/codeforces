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

func expectedDate(s string) string {
	days := []int{31, 28, 31, 30, 31, 30, 31, 31, 30, 31, 30, 31}
	count := make(map[string]int)
	maxCount := 0
	result := ""
	n := len(s)
	for i := 0; i+10 <= n; i++ {
		if s[i+2] != '-' || s[i+5] != '-' {
			continue
		}
		sub := s[i : i+10]
		day, err1 := strconv.Atoi(sub[0:2])
		month, err2 := strconv.Atoi(sub[3:5])
		year, err3 := strconv.Atoi(sub[6:10])
		if err1 != nil || err2 != nil || err3 != nil {
			continue
		}
		if year < 2013 || year > 2015 {
			continue
		}
		if month < 1 || month > 12 {
			continue
		}
		if day < 1 || day > days[month-1] {
			continue
		}
		count[sub]++
		if count[sub] > maxCount {
			maxCount = count[sub]
			result = sub
		}
	}
	return result
}

func randomDate(rng *rand.Rand) string {
	year := rng.Intn(3) + 2013
	month := rng.Intn(12) + 1
	days := []int{31, 28, 31, 30, 31, 30, 31, 31, 30, 31, 30, 31}
	day := rng.Intn(days[month-1]) + 1
	return fmt.Sprintf("%02d-%02d-%04d", day, month, year)
}

func generateString(rng *rand.Rand) string {
	baseLen := rng.Intn(50) + 10
	var sb strings.Builder
	for i := 0; i < baseLen; i++ {
		if rng.Intn(5) == 0 {
			sb.WriteByte('-')
		} else {
			sb.WriteByte(byte('0' + rng.Intn(10)))
		}
	}
	date := randomDate(rng)
	pos := rng.Intn(sb.Len() + 1)
	str := sb.String()
	return str[:pos] + date + str[pos:]
}

func runCase(bin string, s string) error {
	input := s + "\n"
	expected := expectedDate(s)
	cmd := exec.Command(bin)
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	var stderr bytes.Buffer
	cmd.Stderr = &stderr
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("runtime error: %v\nstderr: %s", err, stderr.String())
	}
	got := strings.TrimSpace(out.String())
	if got != expected {
		return fmt.Errorf("expected %s got %s", expected, got)
	}
	return nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierB.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]

	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < 100; i++ {
		s := generateString(rng)
		if err := runCase(bin, s); err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:\n%s\n", i+1, err, s)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
