package main

import (
	"bufio"
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"strings"
)

func runBinary(bin, input string) (string, error) {
	var cmd *exec.Cmd
	if strings.HasSuffix(bin, ".go") {
		cmd = exec.Command("go", "run", bin)
	} else {
		cmd = exec.Command(bin)
	}
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = os.Stderr
	err := cmd.Run()
	return out.String(), err
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}

func solveCaseF(a []int) (bool, int, int, int) {
	n := len(a)
	L := make([]int, n+2)
	R := make([]int, n+2)
	maxL := make([]int, n+2)
	maxR := make([]int, n+2)
	stack := make([]int, 0, n+2)
	a2 := make([]int, n+2)
	copy(a2[1:], a)
	a2[0] = 0
	stack = append(stack, 0)
	for i := 1; i <= n; i++ {
		for len(stack) > 0 && a2[stack[len(stack)-1]] >= a2[i] {
			stack = stack[:len(stack)-1]
		}
		L[i] = stack[len(stack)-1] + 1
		stack = append(stack, i)
	}
	stack = stack[:0]
	a2[n+1] = 0
	stack = append(stack, n+1)
	for i := n; i >= 1; i-- {
		for len(stack) > 0 && a2[stack[len(stack)-1]] >= a2[i] {
			stack = stack[:len(stack)-1]
		}
		R[i] = stack[len(stack)-1] - 1
		stack = append(stack, i)
	}
	maxL[0] = 0
	for i := 1; i <= n; i++ {
		maxL[i] = max(maxL[i-1], a2[i])
	}
	maxR[n+1] = 0
	for i := n; i >= 1; i-- {
		maxR[i] = max(maxR[i+1], a2[i])
	}
	for i := 1; i <= n; i++ {
		u := maxL[L[i]-1]
		v := maxR[R[i]+1]
		Li, Ri := L[i], R[i]
		if u < a2[i] && Li != i && a2[Li] == a2[i] {
			u = a2[i]
			Li++
		}
		if v < a2[i] && Ri != i && a2[Ri] == a2[i] {
			v = a2[i]
			Ri--
		}
		if u == a2[i] && v == a2[i] {
			x := Li - 1
			y := Ri - Li + 1
			z := n - x - y
			return true, x, y, z
		}
	}
	return false, 0, 0, 0
}

func generateTests() ([][]int, string) {
	const t = 100
	r := rand.New(rand.NewSource(6))
	arrays := make([][]int, t)
	var sb strings.Builder
	fmt.Fprintf(&sb, "%d\n", t)
	for i := 0; i < t; i++ {
		n := r.Intn(20) + 3
		fmt.Fprintf(&sb, "%d\n", n)
		arr := make([]int, n)
		for j := 0; j < n; j++ {
			arr[j] = r.Intn(n) + 1
			fmt.Fprintf(&sb, "%d ", arr[j])
		}
		fmt.Fprintln(&sb)
		arrays[i] = arr
	}
	return arrays, sb.String()
}

func verify(arrays [][]int, output string) error {
	scanner := bufio.NewScanner(strings.NewReader(output))
	scanner.Split(bufio.ScanWords)
	for idx, arr := range arrays {
		if !scanner.Scan() {
			return fmt.Errorf("case %d: missing result", idx+1)
		}
		token := scanner.Text()
		success, x, y, z := solveCaseF(arr)
		if token == "NO" {
			if success {
				return fmt.Errorf("case %d: expected YES but got NO", idx+1)
			}
		} else if token == "YES" {
			if !success {
				return fmt.Errorf("case %d: expected NO but got YES", idx+1)
			}
			var X, Y, Z int
			if !scanner.Scan() {
				return fmt.Errorf("case %d: missing x", idx+1)
			}
			fmt.Sscan(scanner.Text(), &X)
			if !scanner.Scan() {
				return fmt.Errorf("case %d: missing y", idx+1)
			}
			fmt.Sscan(scanner.Text(), &Y)
			if !scanner.Scan() {
				return fmt.Errorf("case %d: missing z", idx+1)
			}
			fmt.Sscan(scanner.Text(), &Z)
			if X != x || Y != y || Z != z {
				return fmt.Errorf("case %d: expected %d %d %d got %d %d %d", idx+1, x, y, z, X, Y, Z)
			}
		} else {
			return fmt.Errorf("case %d: unexpected token %s", idx+1, token)
		}
	}
	if scanner.Scan() {
		return fmt.Errorf("extra output: %s", scanner.Text())
	}
	return nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierF.go <binary>")
		os.Exit(1)
	}
	arrays, input := generateTests()
	out, err := runBinary(os.Args[1], input)
	if err != nil {
		fmt.Fprintln(os.Stderr, "error running binary:", err)
		os.Exit(1)
	}
	if err := verify(arrays, out); err != nil {
		fmt.Fprintln(os.Stderr, "verification failed:", err)
		os.Exit(1)
	}
	fmt.Println("All tests passed for problem F")
}
