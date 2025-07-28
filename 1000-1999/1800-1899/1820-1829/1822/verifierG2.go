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

var primes []int

func initPrimes(limit int) {
	sieve := make([]bool, limit+1)
	for i := 2; i <= limit; i++ {
		if !sieve[i] {
			primes = append(primes, i)
			for j := i * i; j <= limit; j += i {
				sieve[j] = true
			}
		}
	}
}

func divisors(n int) []int {
	res := []int{1}
	x := n
	for _, p := range primes {
		if p*p > x {
			break
		}
		if x%p == 0 {
			cnt := 0
			for x%p == 0 {
				x /= p
				cnt++
			}
			base := res
			res = make([]int, 0, len(base)*(cnt+1))
			pow := 1
			for i := 0; i <= cnt; i++ {
				for _, d := range base {
					res = append(res, d*pow)
				}
				pow *= p
			}
		}
	}
	if x > 1 {
		base := res
		res = make([]int, 0, len(base)*2)
		for _, d := range base {
			res = append(res, d)
			res = append(res, d*x)
		}
	}
	return res
}

func solveCase(arr []int) string {
	freq := make(map[int]int)
	for _, v := range arr {
		freq[v]++
	}
	var ans int64
	for _, c := range freq {
		if c >= 3 {
			ans += int64(c) * int64(c-1) * int64(c-2)
		}
	}
	for y, cy := range freq {
		ds := divisors(y)
		for _, i := range ds {
			if i == y {
				continue
			}
			k64 := int64(y) * int64(y) / int64(i)
			if k64 > 1000000000 {
				continue
			}
			k := int(k64)
			ci, ok1 := freq[i]
			ck, ok2 := freq[k]
			if ok1 && ok2 {
				ans += int64(ci) * int64(cy) * int64(ck)
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
	initPrimes(31623)
	if len(os.Args) != 2 && !(len(os.Args) == 3 && os.Args[1] == "--") {
		fmt.Println("usage: go run verifierG2.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[len(os.Args)-1]
	data, err := os.ReadFile("testcasesG2.txt")
	if err != nil {
		fmt.Println("could not read testcasesG2.txt:", err)
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
