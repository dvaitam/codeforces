package main

import (
	"bytes"
	"context"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"time"
)

type test struct {
	k int
	a []int
}

func genTests() []test {
	rand.Seed(5)
	tests := make([]test, 0, 100)
	for len(tests) < 100 {
		k := rand.Intn(5) + 1
		a := make([]int, k)
		for i := 0; i < k; i++ {
			a[i] = rand.Intn(20) + 1
		}
		tests = append(tests, test{k, a})
	}
	return tests
}

func buildInput(t test) string {
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("%d\n", t.k))
	for i, v := range t.a {
		if i > 0 {
			sb.WriteByte(' ')
		}
		sb.WriteString(strconv.Itoa(v))
	}
	sb.WriteByte('\n')
	return sb.String()
}

func runBinary(path, input string) (string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	var cmd *exec.Cmd
	if strings.HasSuffix(path, ".go") {
		cmd = exec.CommandContext(ctx, "go", "run", path)
	} else {
		cmd = exec.CommandContext(ctx, path)
	}
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	err := cmd.Run()
	return out.String(), err
}

func sieve(max int) []int {
	isComp := make([]bool, max+1)
	primes := []int{}
	for i := 2; i <= max; i++ {
		if !isComp[i] {
			primes = append(primes, i)
			if i <= max/i {
				for j := i * i; j <= max; j += i {
					isComp[j] = true
				}
			}
		}
	}
	return primes
}

func requiredN(a []int) int64 {
	maxa := 0
	var sumA int64
	for _, v := range a {
		if v > maxa {
			maxa = v
		}
		sumA += int64(v)
	}
	freq := make([]int, maxa+2)
	for _, v := range a {
		freq[v]++
	}
	suf := make([]int, maxa+2)
	for v := maxa; v >= 0; v-- {
		suf[v] = freq[v] + suf[v+1]
	}
	primes := sieve(maxa)
	B := make([]int64, len(primes))
	for idx, p := range primes {
		var bp int64
		pPow := p
		for pPow <= maxa {
			var fx int64
			for j := pPow; j <= maxa; j += pPow {
				fx += int64(suf[j])
			}
			bp += fx
			if pPow > maxa/p {
				break
			}
			pPow *= p
		}
		B[idx] = bp
	}
	lo, hi := int64(1), sumA
	ans := sumA
	for lo <= hi {
		mid := (lo + hi) / 2
		ok := true
		for idx, p := range primes {
			need := B[idx]
			if need == 0 {
				continue
			}
			var have int64
			n := mid
			for n > 0 && have < need {
				n /= int64(p)
				have += n
			}
			if have < need {
				ok = false
				break
			}
		}
		if ok {
			ans = mid
			hi = mid - 1
		} else {
			lo = mid + 1
		}
	}
	return ans
}

func verifyOutput(out string, t test) bool {
	val, err := strconv.ParseInt(strings.TrimSpace(out), 10, 64)
	if err != nil {
		return false
	}
	expected := requiredN(t.a)
	return val == expected
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("usage: go run verifierE.go /path/to/binary")
		os.Exit(1)
	}
	cand := os.Args[1]
	tests := genTests()
	for i, t := range tests {
		input := buildInput(t)
		out, err := runBinary(cand, input)
		if err != nil {
			fmt.Printf("test %d: run error %v\n", i+1, err)
			os.Exit(1)
		}
		if !verifyOutput(out, t) {
			fmt.Printf("test %d failed. input:\n%soutput:\n%s\n", i+1, input, out)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
