package main

import (
    "bufio"
    "fmt"
    "os"
)

const mod int64 = 998244353

// TODO: implement a proper algorithm for large k.
// Currently handles only very small k via brute force enumeration of seeds.

func bruteForce(k int, fixed []int) int64 {
    n := 1 << k
    teams := make([]int, n)
    for i := 0; i < n; i++ {
        teams[i] = i + 1
    }
    used := make([]bool, n)
    var ans int64
    var dfs func(int)
    dfs = func(pos int) {
        if pos == n {
            if checkPermutation(teams, k) {
                ans = (ans + 1) % mod
            }
            return
        }
        if fixed[pos] != -1 {
            teams[pos] = fixed[pos]
            used[fixed[pos]-1] = true
            dfs(pos + 1)
            used[fixed[pos]-1] = false
            return
        }
        for i := 0; i < n; i++ {
            if used[i] {
                continue
            }
            teams[pos] = i + 1
            used[i] = true
            dfs(pos + 1)
            used[i] = false
        }
    }
    dfs(0)
    return ans
}

func checkPermutation(p []int, k int) bool {
    n := 1 << k
    rounds := [][]int{append([]int{}, p...)}
    for size := n; size > 1; size >>= 1 {
        prev := rounds[len(rounds)-1]
        nxt := make([]int, size>>1)
        for i := 0; i < size; i += 2 {
            if prev[i] < prev[i+1] {
                nxt[i>>1] = prev[i]
            } else {
                nxt[i>>1] = prev[i+1]
            }
        }
        rounds = append(rounds, nxt)
    }
    elimination := make([]int, n+1)
    for r := 0; r < k; r++ {
        prev := rounds[r]
        for i := 0; i < len(prev); i += 2 {
            a, b := prev[i], prev[i+1]
            if a < b {
                elimination[b] = k - r
            } else {
                elimination[a] = k - r
            }
        }
    }
    champion := rounds[len(rounds)-1][0]
    elimination[champion] = 0
    for i := 1; i <= n; i++ {
        place := 1
        if elimination[i] > 0 {
            r := elimination[i] - 1
            place = (1 << r) + 1
        }
        expected := 1
        if i > 1 {
            m := 1
            for (1 << m) < i {
                m++
            }
            expected = (1 << (m-1)) + 1
        }
        if place != expected {
            return false
        }
    }
    return true
}

func main() {
    reader := bufio.NewReader(os.Stdin)
    writer := bufio.NewWriter(os.Stdout)
    defer writer.Flush()

    var k int
    if _, err := fmt.Fscan(reader, &k); err != nil {
        return
    }
    n := 1 << k
    fixed := make([]int, n)
    for i := 0; i < n; i++ {
        fmt.Fscan(reader, &fixed[i])
    }

    if k > 6 {
        fmt.Println(0)
        return
    }
    ans := bruteForce(k, fixed)
    fmt.Fprintln(writer, ans%mod)
}

