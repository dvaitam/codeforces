package main

import (
    "bufio"
    "bytes"
    "fmt"
    "math/rand"
    "os"
    "os/exec"
    "strconv"
    "strings"
)

type movie struct {
    title string
    actors []int
}

type testD struct {
    m int
    fav []int
    movies []movie
}

func genTestsD() []testD {
    rng := rand.New(rand.NewSource(45))
    tests := make([]testD, 100)
    for i := range tests {
        m := rng.Intn(10) + 1
        k := rng.Intn(m) + 1
        fav := rng.Perm(m)[:k]
        for j := range fav { fav[j] += 1 }
        n := rng.Intn(5) + 1
        movies := make([]movie, n)
        for j := 0; j < n; j++ {
            d := rng.Intn(m) + 1
            used := make(map[int]bool)
            actors := make([]int, d)
            for x := 0; x < d; x++ {
                if rng.Intn(3) == 0 {
                    actors[x] = 0
                } else {
                    for {
                        v := rng.Intn(m) + 1
                        if !used[v] {
                            used[v] = true
                            actors[x] = v
                            break
                        }
                    }
                }
            }
            movies[j] = movie{title: fmt.Sprintf("mv%d_%d", i, j), actors: actors}
        }
        tests[i] = testD{m: m, fav: fav, movies: movies}
    }
    return tests
}

func solveD(tc testD) []int {
    m := tc.m
    k := len(tc.fav)
    fav := make([]bool, m+1)
    for _, x := range tc.fav {
        fav[x] = true
    }
    n := len(tc.movies)
    minFav := make([]int, n)
    maxFav := make([]int, n)
    for i, mv := range tc.movies {
        knownFav := 0
        knownNonFav := 0
        zeros := 0
        for _, b := range mv.actors {
            if b == 0 {
                zeros++
            } else if fav[b] {
                knownFav++
            } else {
                knownNonFav++
            }
        }
        availableNonFav := (m - k) - knownNonFav
        extraMin := 0
        if zeros > availableNonFav {
            extraMin = zeros - availableNonFav
        }
        minFav[i] = knownFav + extraMin
        availableFav := k - knownFav
        extraMax := zeros
        if extraMax > availableFav {
            extraMax = availableFav
        }
        maxFav[i] = knownFav + extraMax
    }
    res := make([]int, n)
    for i := 0; i < n; i++ {
        isSureFav := true
        isSureNotFav := false
        for j := 0; j < n; j++ {
            if i == j { continue }
            if maxFav[j] > minFav[i] {
                isSureFav = false
            }
            if minFav[j] > maxFav[i] {
                isSureNotFav = true
            }
        }
        if isSureFav {
            res[i] = 0
        } else if isSureNotFav {
            res[i] = 1
        } else {
            res[i] = 2
        }
    }
    return res
}

func run(bin string, input string) (string, error) {
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
    return out.String(), nil
}

func runCase(bin string, tc testD) error {
    var sb strings.Builder
    sb.WriteString(fmt.Sprintf("%d %d\n", tc.m, len(tc.fav)))
    for i, v := range tc.fav {
        if i > 0 { sb.WriteByte(' ') }
        sb.WriteString(strconv.Itoa(v))
    }
    sb.WriteByte('\n')
    sb.WriteString(fmt.Sprintf("%d\n", len(tc.movies)))
    for _, mv := range tc.movies {
        sb.WriteString(fmt.Sprintf("%s\n", mv.title))
        sb.WriteString(fmt.Sprintf("%d\n", len(mv.actors)))
        for i, v := range mv.actors {
            if i > 0 { sb.WriteByte(' ') }
            sb.WriteString(strconv.Itoa(v))
        }
        sb.WriteByte('\n')
    }
    out, err := run(bin, sb.String())
    if err != nil { return err }
    expected := solveD(tc)
    scanner := bufio.NewScanner(strings.NewReader(strings.TrimSpace(out)))
    for _, exp := range expected {
        if !scanner.Scan() { return fmt.Errorf("missing output line") }
        val, err := strconv.Atoi(strings.TrimSpace(scanner.Text()))
        if err != nil { return fmt.Errorf("bad integer") }
        if val != exp { return fmt.Errorf("expected %d got %d", exp, val) }
    }
    if scanner.Scan() { return fmt.Errorf("extra output") }
    return nil
}

func main() {
    if len(os.Args) != 2 {
        fmt.Fprintln(os.Stderr, "usage: go run verifierD.go /path/to/binary")
        os.Exit(1)
    }
    bin := os.Args[1]
    tests := genTestsD()
    for i, tc := range tests {
        if err := runCase(bin, tc); err != nil {
            fmt.Fprintf(os.Stderr, "case %d failed: %v\n", i+1, err)
            os.Exit(1)
        }
    }
    fmt.Printf("All %d tests passed\n", len(tests))
}

