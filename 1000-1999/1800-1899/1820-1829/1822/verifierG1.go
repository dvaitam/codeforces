package main

import (
	"bufio"
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"strings"
)

func runExe(path, input string) (string, error) {
	var cmd *exec.Cmd
	if strings.HasSuffix(path, ".go") {
		cmd = exec.Command("go", "run", path)
	} else {
		cmd = exec.Command(path)
	}
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("%v\n%s", err, out.String())
	}
	return strings.TrimSpace(out.String()), nil
}

func divisors(x int) []int {
	res := []int{}
	for d := 1; d*d <= x; d++ {
		if x%d == 0 {
			res = append(res, d)
			if d != x/d {
				res = append(res, x/d)
			}
		}
	}
	return res
}

func solveCase(arr []int) string {
	freq := make(map[int]int)
	used := []int{}
	for _, v := range arr {
		if freq[v] == 0 {
			used = append(used, v)
		}
		freq[v]++
	}
	var ans int64
	const maxA = 1000000
	for _, x := range used {
		cx := freq[x]
		divs := divisors(x)
		for _, b := range divs {
			if x*b > maxA {
				continue
			}
			y := x / b
			z := x * b
			cy := freq[y]
			cz := freq[z]
			if cy == 0 || cz == 0 {
				continue
			}
			if y == x && z == x {
				if cx >= 3 {
					ans += int64(cx) * int64(cx-1) * int64(cx-2)
				}
			} else if y == x {
				if cx >= 2 {
					ans += int64(cx) * int64(cx-1) * int64(cz)
				}
			} else if z == x {
				if cx >= 2 {
					ans += int64(cy) * int64(cx) * int64(cx-1)
				}
			} else {
				ans += int64(cx) * int64(cy) * int64(cz)
			}
		}
	}
	return fmt.Sprintf("%d", ans)
}

func join(arr []int) string {
	var sb strings.Builder
	for i, v := range arr {
		if i > 0 {
			sb.WriteByte(' ')
		}
		fmt.Fprintf(&sb, "%d", v)
	}
	return sb.String()
}

func main() {
	if len(os.Args) != 2 && !(len(os.Args) == 3 && os.Args[1] == "--") {
		fmt.Println("usage: go run verifierG1.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[len(os.Args)-1]
	data, err := os.ReadFile("testcasesG1.txt")
	if err != nil {
		fmt.Println("could not read testcasesG1.txt:", err)
		os.Exit(1)
	}
	scan := bufio.NewScanner(bytes.NewReader(data))
	scan.Split(bufio.ScanWords)
	if !scan.Scan() {
		fmt.Println("invalid test file")
		os.Exit(1)
	}
	var tcase int
	fmt.Sscan(scan.Text(), &tcase)
	for idx := 0; idx < tcase; idx++ {
		if !scan.Scan() {
			fmt.Println("bad test file")
			os.Exit(1)
		}
		var n int
		fmt.Sscan(scan.Text(), &n)
		arr := make([]int, n)
		for i := 0; i < n; i++ {
			if !scan.Scan() {
				fmt.Println("bad test file")
				os.Exit(1)
			}
			fmt.Sscan(scan.Text(), &arr[i])
		}
		input := fmt.Sprintf("1\n%d\n%s\n", n, join(arr))
		exp := solveCase(arr)
		got, err := runExe(bin, input)
		if err != nil {
			fmt.Printf("case %d: %v\n", idx+1, err)
			os.Exit(1)
		}
		if strings.TrimSpace(got) != exp {
			fmt.Printf("case %d failed: expected %s got %s\n", idx+1, exp, got)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
