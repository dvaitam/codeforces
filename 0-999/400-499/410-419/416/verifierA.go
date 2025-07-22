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

func solve(input string) string {
	scan := bufio.NewScanner(strings.NewReader(input))
	scan.Split(bufio.ScanWords)
	if !scan.Scan() {
		return ""
	}
	n, _ := strconv.Atoi(scan.Text())
	l := int64(-2000000000)
	r := int64(2000000000)
	for i := 0; i < n; i++ {
		if !scan.Scan() {
			return ""
		}
		op := scan.Text()
		scan.Scan()
		num, _ := strconv.ParseInt(scan.Text(), 10, 64)
		scan.Scan()
		ans := scan.Text()
		switch op {
		case ">":
			if ans == "Y" {
				if num+1 > l {
					l = num + 1
				}
			} else {
				if num < r {
					r = num
				}
			}
		case "<":
			if ans == "Y" {
				if num-1 < r {
					r = num - 1
				}
			} else {
				if num > l {
					l = num
				}
			}
		case ">=":
			if ans == "Y" {
				if num > l {
					l = num
				}
			} else {
				if num-1 < r {
					r = num - 1
				}
			}
		case "<=":
			if ans == "Y" {
				if num < r {
					r = num
				}
			} else {
				if num+1 > l {
					l = num + 1
				}
			}
		}
		if l > r {
			return "Impossible"
		}
	}
	return fmt.Sprintf("%d", l)
}

func generateCase(rng *rand.Rand) (string, string) {
	n := rng.Intn(10) + 1
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("%d\n", n))
	for i := 0; i < n; i++ {
		op := []string{">", "<", ">=", "<="}[rng.Intn(4)]
		num := rng.Int63n(21) - 10
		ans := "N"
		if rng.Intn(2) == 0 {
			ans = "Y"
		}
		sb.WriteString(fmt.Sprintf("%s %d %s\n", op, num, ans))
	}
	input := sb.String()
	exp := solve(input)
	return input, strings.TrimSpace(exp)
}

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

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierA.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < 100; i++ {
		in, exp := generateCase(rng)
		out, err := run(bin, in)
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:\n%s", i+1, err, in)
			os.Exit(1)
		}
		if out != exp {
			fmt.Fprintf(os.Stderr, "case %d failed: expected %s got %s\ninput:\n%s", i+1, exp, out, in)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
