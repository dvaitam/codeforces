package main

import (
	"bytes"
	"context"
	"fmt"
	"math/big"
	"math/rand"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"
)

const oracleSource = `package main

import (
	"fmt"
)

const Y = 201
const offset = 100

func main() {
	var s string
	if _, err := fmt.Scan(&s); err != nil {
		return
	}

	n := make([]int, 65)
	for i := 0; i < len(s); i++ {
		n[i] = int(s[len(s)-1-i] - '0')
	}

	var dp [Y][Y]int64
	var next_dp [Y][Y]int64

	for i := 0; i < Y; i++ {
		for j := 0; j < Y; j++ {
			dp[i][j] = 1e18
		}
	}

	for y1 := 0; y1 < Y; y1++ {
		dp[offset][y1] = 0
	}

	for i := 0; i < 63; i++ {
		for j := 0; j < Y; j++ {
			for k := 0; k < Y; k++ {
				next_dp[j][k] = 1e18
			}
		}

		for y_i := 0; y_i < Y; y_i++ {
			for y_i1 := 0; y_i1 < Y; y_i1++ {
				if dp[y_i][y_i1] >= 1e15 {
					continue
				}
				xi := int64(n[i]) + int64(10*(y_i1-offset)) - int64(y_i-offset)
				for y_i2 := 0; y_i2 < Y; y_i2++ {
					xi1 := int64(n[i+1]) + int64(10*(y_i2-offset)) - int64(y_i1-offset)
					diff := xi - xi1
					if diff < 0 {
						diff = -diff
					}
					cost := diff * int64(i+1)

					if dp[y_i][y_i1]+cost < next_dp[y_i1][y_i2] {
						next_dp[y_i1][y_i2] = dp[y_i][y_i1] + cost
					}
				}
			}
		}

		for j := 0; j < Y; j++ {
			for k := 0; k < Y; k++ {
				dp[j][k] = next_dp[j][k]
			}
		}
	}

	fmt.Println(dp[offset][offset])
}
`

func buildOracle() (string, error) {
	dir := os.TempDir()
	src := filepath.Join(dir, fmt.Sprintf("oracle1487F_%d.go", time.Now().UnixNano()))
	if err := os.WriteFile(src, []byte(oracleSource), 0644); err != nil {
		return "", fmt.Errorf("write oracle source: %v", err)
	}
	defer os.Remove(src)
	oracle := src[:len(src)-3]
	cmd := exec.Command("go", "build", "-o", oracle, src)
	if out, err := cmd.CombinedOutput(); err != nil {
		return "", fmt.Errorf("build oracle failed: %v\n%s", err, out)
	}
	return oracle, nil
}

func run(bin, input string) (string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	cmd := exec.CommandContext(ctx, bin)
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	err := cmd.Run()
	return strings.TrimSpace(out.String()), err
}

func randBig(rng *rand.Rand) string {
	v := rng.Intn(50) + 1
	return fmt.Sprintf("%d", v)
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("Usage: go run verifierF.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	oracle, err := buildOracle()
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	defer os.Remove(oracle)

	rng := rand.New(rand.NewSource(1))
	for t := 0; t < 10; t++ {
		s := randBig(rng)
		input := s + "\n"
		exp, err := run(oracle, input)
		if err != nil {
			fmt.Printf("oracle runtime error: %v\n", err)
			os.Exit(1)
		}
		got, err := run(bin, input)
		if err != nil {
			fmt.Printf("case %d: runtime error: %v\n%s", t+1, err, got)
			os.Exit(1)
		}
		if strings.TrimSpace(got) != strings.TrimSpace(exp) {
			fmt.Printf("case %d failed:\ninput:%sexpected %s got %s\n", t+1, input, exp, got)
			os.Exit(1)
		}
		// ensure output is numeric
		if _, ok := new(big.Int).SetString(strings.TrimSpace(got), 10); !ok {
			fmt.Printf("case %d: output not integer: %s\n", t+1, got)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
