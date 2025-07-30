package main

import (
	"bufio"
	"context"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"sort"
	"strings"
	"time"
)

func runBinary(bin, input string) (string, error) {
	if !strings.Contains(bin, "/") {
		bin = "./" + bin
	}
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()
	cmd := exec.CommandContext(ctx, bin)
	cmd.Stdin = strings.NewReader(input)
	out, err := cmd.CombinedOutput()
	if ctx.Err() == context.DeadlineExceeded {
		return "", fmt.Errorf("timeout")
	}
	if err != nil {
		return "", fmt.Errorf("%v: %s", err, out)
	}
	return strings.TrimSpace(string(out)), nil
}

type outData struct {
	m     int
	sum   int
	pairs [][2]int
}

func parseOut(s string) (outData, error) {
	var d outData
	sc := bufio.NewScanner(strings.NewReader(strings.TrimSpace(s)))
	if !sc.Scan() {
		return d, fmt.Errorf("empty output")
	}
	if _, err := fmt.Sscan(sc.Text(), &d.m, &d.sum); err != nil {
		return d, fmt.Errorf("bad header: %v", err)
	}
	for i := 0; i < d.m; i++ {
		if !sc.Scan() {
			return d, fmt.Errorf("expected %d lines, got %d", d.m, i)
		}
		var a, b int
		if _, err := fmt.Sscan(sc.Text(), &a, &b); err != nil {
			return d, fmt.Errorf("bad pair: %v", err)
		}
		d.pairs = append(d.pairs, [2]int{a, b})
	}
	sort.Slice(d.pairs, func(i, j int) bool {
		if d.pairs[i][0] == d.pairs[j][0] {
			return d.pairs[i][1] < d.pairs[j][1]
		}
		return d.pairs[i][0] < d.pairs[j][0]
	})
	return d, nil
}

func genTest() string {
	n := rand.Intn(4) + 1
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("%d\n", n))
	for i := 0; i < n; i++ {
		sb.WriteString(fmt.Sprintf("%d %d\n", rand.Intn(5)+1, rand.Intn(20)+1))
	}
	k := rand.Intn(4) + 1
	sb.WriteString(fmt.Sprintf("%d\n", k))
	for i := 0; i < k; i++ {
		sb.WriteString(fmt.Sprintf("%d\n", rand.Intn(5)+1))
	}
	return sb.String()
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintf(os.Stderr, "Usage: %s /path/to/binary\n", os.Args[0])
		os.Exit(1)
	}
	cand := os.Args[1]
	ref := "refE"
	if err := exec.Command("go", "build", "-o", ref, "1212E.go").Run(); err != nil {
		fmt.Fprintln(os.Stderr, "failed to build reference:", err)
		os.Exit(1)
	}
	defer os.Remove(ref)

	rand.Seed(1)
	for i := 0; i < 100; i++ {
		input := genTest()
		expStr, err := runBinary(ref, input)
		if err != nil {
			fmt.Fprintln(os.Stderr, "reference failed:", err)
			os.Exit(1)
		}
		gotStr, err := runBinary(cand, input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "test %d failed to run: %v\n", i+1, err)
			os.Exit(1)
		}
		exp, err := parseOut(expStr)
		if err != nil {
			fmt.Fprintln(os.Stderr, "bad reference output:", err)
			os.Exit(1)
		}
		got, err := parseOut(gotStr)
		if err != nil {
			fmt.Fprintf(os.Stderr, "test %d bad output: %v\n", i+1, err)
			os.Exit(1)
		}
		if exp.m != got.m || exp.sum != got.sum || len(exp.pairs) != len(got.pairs) {
			fmt.Fprintf(os.Stderr, "test %d failed:\ninput:\n%sexpected: %s\ngot: %s\n", i+1, input, expStr, gotStr)
			os.Exit(1)
		}
		for j := range exp.pairs {
			if exp.pairs[j] != got.pairs[j] {
				fmt.Fprintf(os.Stderr, "test %d failed:\ninput:\n%sexpected: %s\ngot: %s\n", i+1, input, expStr, gotStr)
				os.Exit(1)
			}
		}
	}
	fmt.Println("All tests passed")
}
