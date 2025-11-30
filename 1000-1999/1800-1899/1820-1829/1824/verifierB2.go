package main

import (
    "bufio"
    "bytes"
    "encoding/base64"
    "fmt"
    "os"
    "os/exec"
    "strconv"
    "strings"
)

const mod int64 = 1_000_000_007

// Base64-encoded contents of testcasesB2.txt.
const testcasesB2 = "MTAwCjQgMQoxIDIKMiAzCjEgNAozIDEKMSAyCjEgMwo1IDEKMSAyCjEgMwozIDQKMyA1CjYgNgoxIDIKMiAzCjMgNAozIDUKMyA2CjEgMQozIDMKMSAyCjEgMwo1IDEKMSAyCjIgMwoyIDQKMSA1CjYgMwoxIDIKMiAzCjIgNAoxIDUKNCA2CjQgMgoxIDIKMSAzCjIgNAo2IDIKMSAyCjEgMwozIDQKNCA1CjEgNgo2IDYKMSAyCjIgMwoyIDQKNCA1CjUgNgo1IDIKMSAyCjIgMwoyIDQKMyA1CjUgMQoxIDIKMSAzCjEgNAoyIDUKMSAxCjMgMwoxIDIKMiAzCjEgMQozIDEKMSAyCjEgMwo0IDEKMSAyCjEgMwoxIDQKMiAxCjEgMgo2IDEKMSAyCjEgMwoxIDQKNCA1CjMgNgoyIDEKMSAyCjMgMgoxIDIKMSAzCjUgNQoxIDIKMiAzCjEgNAozIDUKNCAzCjEgMgoyIDMKMiA0CjUgNAoxIDIKMiAzCjIgNAoyIDUKMiAxCjEgMgo0IDQKMSAyCjIgMwoxIDQKMyAzCjEgMgoyIDMKNiA1CjEgMgoxIDMKMSA0CjIgNQo1IDYKNCAxCjEgMgoyIDMKMyA0CjIgMgoxIDIKMiAxCjEgMgo1IDEKMSAyCjIgMwoyIDQKMiA1CjEgMQo0IDQKMSAyCjIgMwoxIDQKMiAxCjEgMgo2IDEKMSAyCjIgMwoxIDQKMiA1CjUgNgoyIDEKMSAyCjUgMQoxIDIKMSAzCjEgNAo0IDUKNSA1CjEgMgoyIDMKMiA0CjIgNQozIDMKMSAyCjEgMwozIDMKMSAyCjIgMwoxIDEKMiAxCjEgMgoxIDEKMSAxCjUgMQoxIDIKMiAzCjMgNAozIDUKNCAzCjEgMgoyIDMKMSA0CjMgMgoxIDIKMSAzCjYgMgoxIDIKMSAzCjIgNAo0IDUKNCA2CjQgMwoxIDIKMiAzCjIgNAo0IDQKMSAyCjEgMwozIDQKMyAzCjEgMgoyIDMKMiAyCjEgMgozIDEKMSAyCjEgMwo0IDQKMSAyCjEgMwoyIDQKNSA0CjEgMgoxIDMKMSA0CjMgNQoxIDEKNiAyCjEgMgoyIDMKMiA0CjQgNQoxIDYKNSA1CjEgMgoyIDMKMSA0CjIgNQo1IDMKMSAyCjIgMwoxIDQKMSA1CjYgNAoxIDIKMSAzCjIgNAo0IDUKMSA2CjIgMgoxIDIKMSAxCjUgMQoxIDIKMSAzCjEgNAoxIDUKMSAxCjYgNAoxIDIKMiAzCjIgNAoyIDUKNSA2CjQgMQoxIDIKMSAzCjIgNAoxIDEKMSAxCjYgNQoxIDIKMSAzCjIgNAo0IDUKNSA2CjUgMgoxIDIKMiAzCjEgNAoyIDUKMiAxCjEgMgozIDIKMSAyCjIgMwoxIDEKNiA1CjEgMgoyIDMKMSA0CjQgNQoxIDYKNiA2CjEgMgoyIDMKMSA0CjEgNQo0IDYKMiAxCjEgMgo1IDMKMSAyCjIgMwoyIDQKMyA1CjEgMQozIDIKMSAyCjIgMwoyIDEKMSAyCjQgMQoxIDIKMiAzCjEgNAo2IDMKMSAyCjIgMwoxIDQKMSA1CjIgNgo0IDMKMSAyCjEgMwoyIDQKMiAyCjEgMgoxIDEKNCA0CjEgMgoxIDMKMyA0CjUgMgoxIDIKMiAzCjMgNAozIDUKMiAyCjEgMgo2IDUKMSAyCjIgMwozIDQKMSA1CjMgNgozIDEKMSAyCjEgMwo1IDMKMSAyCjIgMwozIDQKNCA1CjMgMwoxIDIKMSAzCjIgMgoxIDIKNiA2CjEgMgoxIDMKMiA0CjEgNQoxIDYKMyAyCjEgMgoxIDMKMiAxCjEgMgo2IDMKMSAyCjIgMwoxIDQKMSA1CjIgNgozIDIKMSAyCjIgMwo="

type testCase struct {
    n, k int
    edges [][2]int
}

// Embedded solver logic from 1824B2.go.
func modPow(a, e int64) int64 {
    r := int64(1)
    for e > 0 {
        if e&1 == 1 {
            r = r * a % mod
        }
        a = a * a % mod
        e >>= 1
    }
    return r
}

var fact, invFact []int64

func initComb(n int) {
    fact = make([]int64, n+1)
    invFact = make([]int64, n+1)
    fact[0] = 1
    for i := 1; i <= n; i++ {
        fact[i] = fact[i-1] * int64(i) % mod
    }
    invFact[n] = modPow(fact[n], mod-2)
    for i := n; i > 0; i-- {
        invFact[i-1] = invFact[i] * int64(i) % mod
    }
}

func comb(n, k int64) int64 {
    if k < 0 || k > n {
        return 0
    }
    return fact[n] * invFact[k] % mod * invFact[n-k] % mod
}

func solve(tc testCase) string {
    n, k := tc.n, tc.k
    initComb(n)

    g := make([][]int, n)
    for _, e := range tc.edges {
        u := e[0] - 1
        v := e[1] - 1
        g[u] = append(g[u], v)
        g[v] = append(g[v], u)
    }

    parent := make([]int, n)
    order := make([]int, 0, n)
    stack := []int{0}
    parent[0] = -1
    for len(stack) > 0 {
        v := stack[len(stack)-1]
        stack = stack[:len(stack)-1]
        order = append(order, v)
        for _, to := range g[v] {
            if to == parent[v] {
                continue
            }
            parent[to] = v
            stack = append(stack, to)
        }
    }

    size := make([]int, n)
    for i := n - 1; i >= 0; i-- {
        v := order[i]
        size[v] = 1
        for _, to := range g[v] {
            if to == parent[v] {
                continue
            }
            size[v] += size[to]
        }
    }

    half := k / 2
    combNK := comb(int64(n), int64(k))
    invCombNK := modPow(combNK, mod-2)
    cache := make(map[int]int64)

    calc := func(s int) int64 {
        if s <= half {
            return 0
        }
        if v, ok := cache[s]; ok {
            return v
        }
        var sum int64
        upper := k
        if s < upper {
            upper = s
        }
        for j := half + 1; j <= upper; j++ {
            sum = (sum + comb(int64(s), int64(j))*comb(int64(n-s), int64(k-j))) % mod
        }
        cache[s] = sum
        return sum
    }

    var total int64
    for v := 1; v < n; v++ {
        s := size[v]
        total = (total + calc(s)) % mod
        total = (total + calc(n-s)) % mod
    }

    ans := (int64(n)%mod*combNK%mod - total + mod) % mod
    ans = ans * invCombNK % mod
    return strconv.FormatInt(ans, 10)
}

func parseTestcases() ([]testCase, error) {
    raw, err := base64.StdEncoding.DecodeString(testcasesB2)
    if err != nil {
        return nil, err
    }
    sc := bufio.NewScanner(bytes.NewReader(raw))
    sc.Split(bufio.ScanWords)
    if !sc.Scan() {
        return nil, fmt.Errorf("invalid test data")
    }
    t, err := strconv.Atoi(sc.Text())
    if err != nil {
        return nil, fmt.Errorf("parse t: %v", err)
    }
    cases := make([]testCase, 0, t)
    for i := 0; i < t; i++ {
        if !sc.Scan() {
            return nil, fmt.Errorf("case %d missing n", i+1)
        }
        n, err := strconv.Atoi(sc.Text())
        if err != nil {
            return nil, fmt.Errorf("case %d n: %v", i+1, err)
        }
        if !sc.Scan() {
            return nil, fmt.Errorf("case %d missing k", i+1)
        }
        k, err := strconv.Atoi(sc.Text())
        if err != nil {
            return nil, fmt.Errorf("case %d k: %v", i+1, err)
        }
        edges := make([][2]int, n-1)
        for j := 0; j < n-1; j++ {
            if !sc.Scan() {
                return nil, fmt.Errorf("case %d truncated edge u", i+1)
            }
            u, err := strconv.Atoi(sc.Text())
            if err != nil {
                return nil, fmt.Errorf("case %d edge u: %v", i+1, err)
            }
            if !sc.Scan() {
                return nil, fmt.Errorf("case %d truncated edge v", i+1)
            }
            v, err := strconv.Atoi(sc.Text())
            if err != nil {
                return nil, fmt.Errorf("case %d edge v: %v", i+1, err)
            }
            edges[j] = [2]int{u, v}
        }
        cases = append(cases, testCase{n: n, k: k, edges: edges})
    }
    return cases, nil
}

func buildIfGo(path string) (string, func(), error) {
    if strings.HasSuffix(path, ".go") {
        tmp, err := os.CreateTemp("", "bin*")
        if err != nil {
            return "", nil, err
        }
        tmp.Close()
        if out, err := exec.Command("go", "build", "-o", tmp.Name(), path).CombinedOutput(); err != nil {
            os.Remove(tmp.Name())
            return "", nil, fmt.Errorf("build failed: %v\n%s", err, out)
        }
        return tmp.Name(), func() { os.Remove(tmp.Name()) }, nil
    }
    return path, func() {}, nil
}

func runCandidate(bin, input string) (string, error) {
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
    if len(os.Args) != 2 {
        fmt.Println("Usage: go run verifierB2.go /path/to/binary")
        os.Exit(1)
    }

    cases, err := parseTestcases()
    if err != nil {
        fmt.Fprintln(os.Stderr, err)
        os.Exit(1)
    }

    bin, cleanup, err := buildIfGo(os.Args[1])
    if err != nil {
        fmt.Fprintln(os.Stderr, err)
        os.Exit(1)
    }
    defer cleanup()

    for idx, tc := range cases {
        var input strings.Builder
        fmt.Fprintf(&input, "%d %d\n", tc.n, tc.k)
        for _, e := range tc.edges {
            fmt.Fprintf(&input, "%d %d\n", e[0], e[1])
        }
        want := solve(tc)
        got, err := runCandidate(bin, input.String())
        if err != nil {
            fmt.Fprintf(os.Stderr, "case %d failed: %v\n", idx+1, err)
            os.Exit(1)
        }
        if got != want {
            fmt.Printf("case %d failed\ninput:\n%sexpected: %s\ngot: %s\n", idx+1, input.String(), want, got)
            os.Exit(1)
        }
    }
    fmt.Printf("All %d tests passed\n", len(cases))
}
