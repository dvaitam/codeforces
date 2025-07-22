package main

import (
	"bufio"
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"time"
)

func solve(data string) string {
	scanner := bufio.NewScanner(strings.NewReader(data))
	scanner.Split(bufio.ScanWords)
	next := func() string {
		if scanner.Scan() {
			return scanner.Text()
		}
		return ""
	}
	n, _ := strconv.Atoi(next())
	names := make([]string, n)
	taxi := make([]int, n)
	pizza := make([]int, n)
	girls := make([]int, n)
	for i := 0; i < n; i++ {
		si, _ := strconv.Atoi(next())
		name := next()
		names[i] = name
		for j := 0; j < si; j++ {
			num := next()
			digits := make([]int, 0, 6)
			for _, c := range num {
				if c >= '0' && c <= '9' {
					digits = append(digits, int(c-'0'))
				}
			}
			allEqual := true
			for k := 1; k < 6; k++ {
				if digits[k] != digits[0] {
					allEqual = false
					break
				}
			}
			if allEqual {
				taxi[i]++
				continue
			}
			dec := true
			for k := 1; k < 6; k++ {
				if digits[k] >= digits[k-1] {
					dec = false
					break
				}
			}
			if dec {
				pizza[i]++
			} else {
				girls[i]++
			}
		}
	}
	maxTaxi, maxPizza, maxGirls := 0, 0, 0
	for i := 0; i < n; i++ {
		if taxi[i] > maxTaxi {
			maxTaxi = taxi[i]
		}
		if pizza[i] > maxPizza {
			maxPizza = pizza[i]
		}
		if girls[i] > maxGirls {
			maxGirls = girls[i]
		}
	}
	var taxiNames, pizzaNames, girlNames []string
	for i := 0; i < n; i++ {
		if taxi[i] == maxTaxi {
			taxiNames = append(taxiNames, names[i])
		}
		if pizza[i] == maxPizza {
			pizzaNames = append(pizzaNames, names[i])
		}
		if girls[i] == maxGirls {
			girlNames = append(girlNames, names[i])
		}
	}
	var sb strings.Builder
	fmt.Fprintf(&sb, "If you want to call a taxi, you should call: %s.\n", strings.Join(taxiNames, ", "))
	fmt.Fprintf(&sb, "If you want to order a pizza, you should call: %s.\n", strings.Join(pizzaNames, ", "))
	fmt.Fprintf(&sb, "If you want to go to a cafe with a wonderful girl, you should call: %s.\n", strings.Join(girlNames, ", "))
	return sb.String()
}

func randomTaxiNumber(rng *rand.Rand) string {
	d := rng.Intn(10)
	return fmt.Sprintf("%d%d-%d%d-%d%d", d, d, d, d, d, d)
}

func randomPizzaNumber(rng *rand.Rand) string {
	start := rng.Intn(5) + 5 // 5..9
	digits := []int{start, start - 1, start - 2, start - 3, start - 4, start - 5}
	return fmt.Sprintf("%d%d-%d%d-%d%d", digits[0], digits[1], digits[2], digits[3], digits[4], digits[5])
}

func randomGirlNumber(rng *rand.Rand) string {
	for {
		d := [6]int{}
		for i := 0; i < 6; i++ {
			d[i] = rng.Intn(10)
		}
		allEqual := true
		for i := 1; i < 6; i++ {
			if d[i] != d[0] {
				allEqual = false
			}
		}
		dec := true
		for i := 1; i < 6; i++ {
			if d[i] >= d[i-1] {
				dec = false
			}
		}
		if !allEqual && !dec {
			return fmt.Sprintf("%d%d-%d%d-%d%d", d[0], d[1], d[2], d[3], d[4], d[5])
		}
	}
}

func generateCase(rng *rand.Rand) (string, string) {
	n := rng.Intn(4) + 1
	var sb strings.Builder
	fmt.Fprintf(&sb, "%d\n", n)
	for i := 0; i < n; i++ {
		si := rng.Intn(4) + 1
		name := fmt.Sprintf("f%d", i+1)
		fmt.Fprintf(&sb, "%d %s\n", si, name)
		for j := 0; j < si; j++ {
			t := rng.Intn(3)
			var num string
			if t == 0 {
				num = randomTaxiNumber(rng)
			} else if t == 1 {
				num = randomPizzaNumber(rng)
			} else {
				num = randomGirlNumber(rng)
			}
			sb.WriteString(num + "\n")
		}
	}
	input := sb.String()
	expected := solve(input)
	return input, expected
}

func runCase(exe, input, expect string) error {
	cmd := exec.Command(exe)
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("runtime error: %v\n%s", err, out.String())
	}
	if strings.TrimSpace(out.String()) != strings.TrimSpace(expect) {
		return fmt.Errorf("expected %q got %q", strings.TrimSpace(expect), strings.TrimSpace(out.String()))
	}
	return nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierB.go /path/to/binary")
		os.Exit(1)
	}
	exe := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < 100; i++ {
		in, exp := generateCase(rng)
		if err := runCase(exe, in, exp); err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:\n%s", i+1, err, in)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
