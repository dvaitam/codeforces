package main

import (
    "fmt"
)

var vis [256]int
var edges [256]rune
var ans []rune

func dfs(u rune) {
    vis[u] = 3
    v := edges[u]
    if v >= 'a' && vis[v] != 3 {
        dfs(v)
    }
    ans = append(ans, u)
}

func main() {
    var n int
    if _, err := fmt.Scan(&n); err != nil {
        return
    }
    var s string
    for i := 0; i < n; i++ {
        fmt.Scan(&s)
        for j := 0; j+1 < len(s); j++ {
            u := rune(s[j])
            v := rune(s[j+1])
            edges[u] = v
            vis[v] = 2
        }
        u0 := rune(s[0])
        if vis[u0] != 2 {
            vis[u0] = 1
        }
    }
    for u := 'a'; u <= 'z'; u++ {
        if vis[u] == 1 {
            dfs(u)
        }
    }
    // reverse ans
    for i, j := 0, len(ans)-1; i < j; i, j = i+1, j-1 {
        ans[i], ans[j] = ans[j], ans[i]
    }
    fmt.Println(string(ans))
}
