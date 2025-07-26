package main

import (
	"bytes"
	"fmt"
	"math"
	"math/rand"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
	"time"
)

type orderBook struct {
	buys  []int
	sells []int
	used  map[int]bool
}

func insertAsc(arr []int, v int) []int {
	i := 0
	for i < len(arr) && arr[i] < v {
		i++
	}
	arr = append(arr, 0)
	copy(arr[i+1:], arr[i:])
	arr[i] = v
	return arr
}

func insertDesc(arr []int, v int) []int {
	i := 0
	for i < len(arr) && arr[i] > v {
		i++
	}
	arr = append(arr, 0)
	copy(arr[i+1:], arr[i:])
	arr[i] = v
	return arr
}

func genCase(r *rand.Rand) string {
	steps := r.Intn(20) + 1
	ob := orderBook{used: make(map[int]bool)}
	var ops []string
	for i := 0; i < steps; i++ {
		if len(ob.buys)+len(ob.sells) > 0 && r.Intn(2) == 0 {
			// accept
			if len(ob.buys) > 0 && (len(ob.sells) == 0 || r.Intn(2) == 0) {
				p := ob.buys[0]
				ob.buys = ob.buys[1:]
				ops = append(ops, fmt.Sprintf("ACCEPT %d", p))
			} else {
				p := ob.sells[0]
				ob.sells = ob.sells[1:]
				ops = append(ops, fmt.Sprintf("ACCEPT %d", p))
			}
		} else {
			bestBuy := -1
			if len(ob.buys) > 0 {
				bestBuy = ob.buys[0]
			}
			bestSell := math.MaxInt32
			if len(ob.sells) > 0 {
				bestSell = ob.sells[0]
			}
			dir := r.Intn(2)
			p := 0
			for {
				p = r.Intn(1000) + 1
				if !ob.used[p] {
					break
				}
			}
			ob.used[p] = true
			if dir == 0 { // buy
				if p >= bestSell {
					p = bestSell - 1
				}
				if p <= 0 {
					p = 1
				}
				ob.buys = insertDesc(ob.buys, p)
			} else { // sell
				if p <= bestBuy {
					p = bestBuy + 1
				}
				ob.sells = insertAsc(ob.sells, p)
			}
			ops = append(ops, fmt.Sprintf("ADD %d", p))
		}
	}
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("%d\n", len(ops)))
	for _, op := range ops {
		sb.WriteString(op + "\n")
	}
	return sb.String()
}

func buildRef() (string, error) {
	_, cur, _, _ := runtime.Caller(0)
	dir := filepath.Dir(cur)
	src := filepath.Join(dir, "1028D.go")
	bin := filepath.Join(os.TempDir(), "ref1028D.bin")
	cmd := exec.Command("go", "build", "-o", bin, src)
	if out, err := cmd.CombinedOutput(); err != nil {
		return "", fmt.Errorf("build reference failed: %v\n%s", err, string(out))
	}
	return bin, nil
}

func run(bin string, input string) (string, error) {
	cmd := exec.Command(bin)
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	var errBuf bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &errBuf
	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("runtime error: %v\n%s", err, errBuf.String())
	}
	return strings.TrimSpace(out.String()), nil
}

func main() {
	rand.Seed(time.Now().UnixNano())
	if len(os.Args) < 2 {
		fmt.Println("usage: go run verifierD.go /path/to/binary")
		return
	}
	cand := os.Args[1]
	ref, err := buildRef()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	defer os.Remove(ref)
	for i := 1; i <= 100; i++ {
		in := genCase(rand.New(rand.NewSource(int64(i))))
		want, err := run(ref, in)
		if err != nil {
			fmt.Printf("reference failed on test %d: %v\n", i, err)
			os.Exit(1)
		}
		got, err := run(cand, in)
		if err != nil {
			fmt.Printf("runtime error on test %d: %v\n", i, err)
			os.Exit(1)
		}
		if strings.TrimSpace(want) != strings.TrimSpace(got) {
			fmt.Printf("wrong answer on test %d\ninput:\n%sexpected:%s\ngot:%s\n", i, in, want, got)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
