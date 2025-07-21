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

var alcohols = map[string]struct{}{
	"ABSINTH": {}, "BEER": {}, "BRANDY": {}, "CHAMPAGNE": {},
	"GIN": {}, "RUM": {}, "SAKE": {}, "TEQUILA": {},
	"VODKA": {}, "WHISKEY": {}, "WINE": {},
}

func generateCase(rng *rand.Rand) (string, int) {
	n := rng.Intn(100) + 1
	var b strings.Builder
	fmt.Fprintf(&b, "%d\n", n)
	count := 0
	drinks := []string{"ABSINTH", "BEER", "BRANDY", "CHAMPAGNE", "GIN", "RUM", "SAKE", "TEQUILA", "VODKA", "WHISKEY", "WINE"}
	for i := 0; i < n; i++ {
		if rng.Intn(2) == 0 { // age
			age := rng.Intn(1010) // 0..1009
			fmt.Fprintf(&b, "%d\n", age)
			if age < 18 {
				count++
			}
		} else {
			var drink string
			if rng.Intn(2) == 0 {
				drink = drinks[rng.Intn(len(drinks))]
				count++
			} else {
				// non-alcohol string
				drink = fmt.Sprintf("NONALC%d", rng.Intn(1000))
			}
			fmt.Fprintf(&b, "%s\n", drink)
		}
	}
	return b.String(), count
}

func runCase(bin string, input string, expected int) error {
	cmd := exec.Command(bin)
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("runtime error: %v\n%s", err, out.String())
	}
	gotStr := strings.TrimSpace(out.String())
	got, err := strconv.Atoi(gotStr)
	if err != nil {
		return fmt.Errorf("invalid output: %s", gotStr)
	}
	if got != expected {
		return fmt.Errorf("expected %d got %d", expected, got)
	}
	return nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierA.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < 100; i++ {
		input, expected := generateCase(rng)
		if err := runCase(bin, input, expected); err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:\n%s", i+1, err, input)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
