package main

import (
    "bytes"
    "fmt"
    "math/rand"
    "os"
    "os/exec"
    "sort"
    "strconv"
    "strings"
)

type testCase struct {
    arr []int64
}

func isPow2(x int64) bool {
    return x > 0 && (x&(x-1)) == 0
}

func hasTriple(arr []int64, m map[int64]int) bool {
    // check existence of x, x+d, x+2d for some power-of-two d
    if len(arr) == 0 {
        return false
    }
    sort.Slice(arr, func(i, j int) bool { return arr[i] < arr[j] })
    uniq := arr[:0]
    for _, v := range arr {
        if len(uniq) == 0 || uniq[len(uniq)-1] != v {
            uniq = append(uniq, v)
        }
    }
    arr = uniq
    last := arr[len(arr)-1]
    for _, x := range arr {
        for d := int64(1); x+2*d <= last; d <<= 1 {
            if m[x+d] > 0 && m[x+2*d] > 0 {
                return true
            }
        }
    }
    return false
}

func hasPair(arr []int64, m map[int64]int) bool {
    if len(arr) == 0 {
        return false
    }
    sort.Slice(arr, func(i, j int) bool { return arr[i] < arr[j] })
    uniq := arr[:0]
    for _, v := range arr {
        if len(uniq) == 0 || uniq[len(uniq)-1] != v {
            uniq = append(uniq, v)
        }
    }
    arr = uniq
    last := arr[len(arr)-1]
    for _, x := range arr {
        for d := int64(1); x+d <= last; d <<= 1 {
            if m[x+d] > 0 {
                return true
            }
        }
    }
    return false
}

func checkOutput(arr []int64, out string) error {
    // Build multiset counts for presence checks
    cnt := make(map[int64]int, len(arr))
    for _, v := range arr {
        cnt[v]++
    }

    fields := strings.Fields(out)
    if len(fields) == 0 {
        return fmt.Errorf("empty output")
    }
    c, err := strconv.Atoi(fields[0])
    if err != nil {
        return fmt.Errorf("first token not an integer: %v", err)
    }
    if c < 1 || c > 3 {
        return fmt.Errorf("invalid count %d", c)
    }
    if len(fields) < 1+c {
        return fmt.Errorf("expected %d numbers got %d", c, len(fields)-1)
    }
    vals := make([]int64, c)
    for i := 0; i < c; i++ {
        v, err := strconv.ParseInt(fields[1+i], 10, 64)
        if err != nil {
            return fmt.Errorf("invalid number: %v", err)
        }
        vals[i] = v
    }
    // Validate presence and constraints
    // presence: ensure each selected number exists at least once in arr
    // Use a temp copy of counts to handle duplicates robustly
    tmp := make(map[int64]int, len(cnt))
    for k, v := range cnt {
        tmp[k] = v
    }
    for _, v := range vals {
        if tmp[v] == 0 {
            return fmt.Errorf("value %d not present in input", v)
        }
        tmp[v]--
    }

    // Determine maximum possible size for this arr
    max := 1
    if hasPair(arr, cnt) {
        max = 2
    }
    if hasTriple(arr, cnt) {
        max = 3
    }
    if c != max {
        return fmt.Errorf("reported size %d but maximum is %d", c, max)
    }
    if c == 1 {
        return nil
    }
    // For pair/triple, validate power-of-two differences
    sort.Slice(vals, func(i, j int) bool { return vals[i] < vals[j] })
    if c == 2 {
        d := vals[1] - vals[0]
        if !isPow2(d) {
            return fmt.Errorf("pair difference %d is not power of two", d)
        }
        return nil
    }
    // c == 3
    if !(isPow2(vals[1]-vals[0]) && vals[1]-vals[0] == vals[2]-vals[1]) {
        return fmt.Errorf("triple must be equally spaced by power of two")
    }
    return nil
}

func buildInput(tc testCase) string {
	var sb strings.Builder
	fmt.Fprintf(&sb, "%d\n", len(tc.arr))
	for i, v := range tc.arr {
		if i > 0 {
			sb.WriteByte(' ')
		}
		fmt.Fprintf(&sb, "%d", v)
	}
	sb.WriteByte('\n')
	return sb.String()
}

func runCase(bin string, tc testCase) error {
    input := buildInput(tc)
    cmd := exec.Command(bin)
    cmd.Stdin = strings.NewReader(input)
    var out bytes.Buffer
    cmd.Stdout = &out
    cmd.Stderr = &out
    if err := cmd.Run(); err != nil {
        return fmt.Errorf("runtime error: %v\n%s", err, out.String())
    }
    got := strings.TrimSpace(out.String())
    if err := checkOutput(append([]int64(nil), tc.arr...), got); err != nil {
        return err
    }
    return nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierD.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]

	var cases []testCase
	cases = append(cases, testCase{arr: []int64{7, 3, 5}})
	cases = append(cases, testCase{arr: []int64{1}})
	cases = append(cases, testCase{arr: []int64{1, 3}})

	rng := rand.New(rand.NewSource(3))
	for i := 0; i < 100; i++ {
		n := rng.Intn(8) + 1
		arr := make([]int64, n)
		for j := 0; j < n; j++ {
			arr[j] = int64(rng.Intn(41) - 20)
		}
		sort.Slice(arr, func(i, j int) bool { return arr[i] < arr[j] })
		arr = arr[:n]
		cases = append(cases, testCase{arr: arr})
	}

	for i, tc := range cases {
		if err := runCase(bin, tc); err != nil {
			fmt.Fprintf(os.Stderr, "test %d failed: %v\ninput: %s", i+1, err, buildInput(tc))
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
