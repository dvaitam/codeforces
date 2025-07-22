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

func run(bin, input string) (string, error) {
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

func runCase(bin string, date time.Time, shift int) error {
	input := fmt.Sprintf("%s %d\n", date.Format("02.01.2006"), shift)
	out, err := run(bin, input)
	if err != nil {
		return err
	}
	expected := date.AddDate(0, 0, shift).Format("02.01.2006")
	if out != expected {
		return fmt.Errorf("expected %s got %s", expected, out)
	}
	return nil
}

func randomDate(rng *rand.Rand) time.Time {
	year := rng.Intn(41) + 1980 // 1980-2020
	// choose a random day in the year
	start := time.Date(year, 1, 1, 0, 0, 0, 0, time.UTC)
	day := rng.Intn(365)
	if start.AddDate(0, 0, day).Year() != year { // handle leap year
		day = rng.Intn(366)
	}
	return start.AddDate(0, 0, day)
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierD.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	edge := []struct {
		date  time.Time
		shift int
	}{
		{time.Date(1980, 1, 1, 0, 0, 0, 0, time.UTC), 0},
		{time.Date(2020, 12, 31, 0, 0, 0, 0, time.UTC), -1},
		{time.Date(2016, 2, 28, 0, 0, 0, 0, time.UTC), 1}, // leap year
		{time.Date(2019, 12, 31, 0, 0, 0, 0, time.UTC), 1},
	}
	idx := 0
	for ; idx < len(edge); idx++ {
		e := edge[idx]
		if err := runCase(bin, e.date, e.shift); err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v (date=%s shift=%d)\n", idx+1, err, e.date.Format("02.01.2006"), e.shift)
			os.Exit(1)
		}
	}
	for ; idx < 100; idx++ {
		d := randomDate(rng)
		shift := rng.Intn(2001) - 1000
		if err := runCase(bin, d, shift); err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v (date=%s shift=%d)\n", idx+1, err, d.Format("02.01.2006"), shift)
			os.Exit(1)
		}
	}
	fmt.Printf("All %d tests passed\n", idx)
}
