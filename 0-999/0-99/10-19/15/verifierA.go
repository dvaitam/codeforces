package main

import (
    "bytes"
    "fmt"
    "math/rand"
    "os"
    "os/exec"
    "sort"
    "strings"
    "time"
)

func run(bin, input string) (string, error) {
    var cmd *exec.Cmd
    if strings.HasSuffix(bin, ".go") {
        cmd = exec.Command("go", "run", bin)
    } else {
        cmd = exec.Command(bin)
    }
    cmd.Stdin = strings.NewReader(input)
    var out bytes.Buffer
    cmd.Stdout = &out
    cmd.Stderr = &out
    if err := cmd.Run(); err != nil {
        return "", fmt.Errorf("runtime error: %v\n%s", err, out.String())
    }
    return strings.TrimSpace(out.String()), nil
}

func expected(n, t int, arr [][2]int) string {
    type interval struct{ l, r int }
    ints := make([]interval, n)
    for i := 0; i < n; i++ {
        x := arr[i][0]
        a := arr[i][1]
        ints[i].l = 2*x - a
        ints[i].r = 2*x + a
    }
    sort.Slice(ints, func(i, j int) bool {
        if ints[i].l != ints[j].l {
            return ints[i].l < ints[j].l
        }
        return ints[i].r < ints[j].r
    })
    ans := 2
    t2 := 2 * t
    for i := 0; i+1 < n; i++ {
        gap := ints[i+1].l - ints[i].r
        if gap == t2 {
            ans++
        } else if gap > t2 {
            ans += 2
        }
    }
    return fmt.Sprintf("%d", ans)
}

func genCase(rng *rand.Rand) (string, string) {
    n := rng.Intn(10) + 1
    t := rng.Intn(20) + 1
    arr := make([][2]int, n)
    // generate non-overlapping intervals by sorting centers
    xs := make([]int, n)
    for i := 0; i < n; i++ {
        xs[i] = rng.Intn(200) - 100
    }
    sort.Ints(xs)
    for i := 0; i < n; i++ {
        a := rng.Intn(10) + 1
        arr[i][0] = xs[i]
        arr[i][1] = a
        if i > 0 {
            prevR := 2*arr[i-1][0] + arr[i-1][1]
            curL := 2*arr[i][0] - arr[i][1]
            if curL < prevR {
                diff := prevR - curL
                arr[i][0] += (diff+1)/2
            }
        }
    }
    var sb strings.Builder
    sb.WriteString(fmt.Sprintf("%d %d\n", n, t))
    for i := 0; i < n; i++ {
        sb.WriteString(fmt.Sprintf("%d %d\n", arr[i][0], arr[i][1]))
    }
    exp := expected(n, t, arr)
    return sb.String(), exp
}

func main() {
    if len(os.Args) != 2 {
        fmt.Fprintln(os.Stderr, "usage: go run verifierA.go /path/to/binary")
        os.Exit(1)
    }
    bin := os.Args[1]
    rng := rand.New(rand.NewSource(time.Now().UnixNano()))
    for i := 0; i < 100; i++ {
        in, exp := genCase(rng)
        out, err := run(bin, in)
        if err != nil {
            fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:\n%s", i+1, err, in)
            os.Exit(1)
        }
        if out != exp {
            fmt.Fprintf(os.Stderr, "case %d failed: expected %s got %s\ninput:\n%s", i+1, exp, out, in)
            os.Exit(1)
        }
    }
    fmt.Println("All tests passed")
}

