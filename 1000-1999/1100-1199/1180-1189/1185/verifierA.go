package main

import (
	"bufio"
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"sort"
	"strconv"
	"strings"
)

func expected(a, b, c, d int) string {
	p := []int{a, b, c}
	sort.Ints(p)
	ans := 0
	if p[1]-p[0] < d {
		ans += d - (p[1] - p[0])
	}
	if p[2]-p[1] < d {
		ans += d - (p[2] - p[1])
	}
	return fmt.Sprintf("%d", ans)
}

func runCase(bin string, a, b, c, d int) error {
	input := fmt.Sprintf("%d %d %d %d\n", a, b, c, d)
	cmd := exec.Command(bin)
	if strings.HasSuffix(bin, ".go") {
		cmd = exec.Command("go", "run", bin)
	}
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	var errBuf bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &errBuf
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("runtime error: %v\n%s", err, errBuf.String())
	}
	got := strings.TrimSpace(out.String())
	exp := expected(a, b, c, d)
	if got != exp {
		return fmt.Errorf("expected %s got %s", exp, got)
	}
	return nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("Usage: go run verifierA.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	f, err := os.Open("testcasesA.txt")
	if err != nil {
		fmt.Println("could not open testcasesA.txt:", err)
		os.Exit(1)
	}
	defer f.Close()
	scanner := bufio.NewScanner(f)
	scanner.Split(bufio.ScanWords)
	if !scanner.Scan() {
		fmt.Println("invalid test file")
		os.Exit(1)
	}
	t, err := strconv.Atoi(scanner.Text())
	if err != nil {
		fmt.Println("invalid test count")
		os.Exit(1)
	}
	for i := 0; i < t; i++ {
		if !scanner.Scan() {
			fmt.Println("invalid test file")
			os.Exit(1)
		}
		a, _ := strconv.Atoi(scanner.Text())
		scanner.Scan()
		b, _ := strconv.Atoi(scanner.Text())
		scanner.Scan()
		c, _ := strconv.Atoi(scanner.Text())
		scanner.Scan()
		d, _ := strconv.Atoi(scanner.Text())
		if err := runCase(bin, a, b, c, d); err != nil {
			fmt.Printf("case %d failed: %v\n", i+1, err)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
