package main

import (
	"bufio"
	"bytes"
	"context"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"strings"
	"time"
)

func runCase(bin string, arr []int) error {
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	var cmd *exec.Cmd
	if strings.HasSuffix(bin, ".go") {
		cmd = exec.CommandContext(ctx, "go", "run", bin)
	} else {
		cmd = exec.CommandContext(ctx, bin)
	}
	stdin, err := cmd.StdinPipe()
	if err != nil {
		return err
	}
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		return err
	}
	var stderr bytes.Buffer
	cmd.Stderr = &stderr
	if err := cmd.Start(); err != nil {
		return err
	}
	r := bufio.NewReader(stdout)
	w := bufio.NewWriter(stdin)
	fmt.Fprintf(w, "%d\n", len(arr))
	w.Flush()
	for {
		var op string
		if _, err := fmt.Fscan(r, &op); err != nil {
			cmd.Wait()
			return fmt.Errorf("failed reading command: %v\n%s", err, stderr.String())
		}
		if op == "!" {
			out := make([]int, len(arr))
			for i := range out {
				if _, err := fmt.Fscan(r, &out[i]); err != nil {
					cmd.Wait()
					return fmt.Errorf("failed to read answer: %v", err)
				}
			}
			if err := cmd.Wait(); err != nil {
				return fmt.Errorf("runtime error: %v\n%s", err, stderr.String())
			}
			for i := range arr {
				if arr[i] != out[i] {
					return fmt.Errorf("expected %v got %v", arr, out)
				}
			}
			break
		}
		var i, j int
		if _, err := fmt.Fscan(r, &i, &j); err != nil {
			cmd.Wait()
			return fmt.Errorf("failed reading indices: %v", err)
		}
		i--
		j--
		var ans int
		switch op {
		case "AND":
			ans = arr[i] & arr[j]
		case "OR":
			ans = arr[i] | arr[j]
		case "XOR":
			ans = arr[i] ^ arr[j]
		default:
			cmd.Wait()
			return fmt.Errorf("unknown operation %s", op)
		}
		fmt.Fprintf(w, "%d\n", ans)
		w.Flush()
	}
	return nil
}

func randomArray(rng *rand.Rand) []int {
	pow := rng.Intn(3) + 2 // 4..32
	n := 1 << pow
	arr := make([]int, n)
	for i := range arr {
		arr[i] = rng.Intn(n)
	}
	return arr
}

func main() {
	if len(os.Args) != 2 && !(len(os.Args) == 3 && os.Args[1] == "--") {
		fmt.Println("usage: go run verifierE1.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[len(os.Args)-1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < 100; i++ {
		arr := randomArray(rng)
		if err := runCase(bin, arr); err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\n", i+1, err)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
