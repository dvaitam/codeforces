package main

import (
"bufio"
"fmt"
"os"
)

const mod = 998244353

// This solution implements a very naive brute force enumerator for the
// problem. It generates all sequences of length n consisting of numbers
// from 1..n and checks whether they satisfy the condition from the
// statement. For every valid sequence it accumulates how many times each
// value appears. The algorithm is exponentially slow and is intended only
// as a placeholder illustrating the required computations.

func main() {
in := bufio.NewReader(os.Stdin)
out := bufio.NewWriter(os.Stdout)
defer out.Flush()

var n int
if _, err := fmt.Fscan(in, &n); err != nil {
return
}

seq := make([]int, n)
ans := make([]int64, n+1)

var dfs func(int)
dfs = func(pos int) {
if pos == n {
if goodSequence(seq) {
for _, v := range seq {
if v <= n {
ans[v]++
}
}
}
return
}
for val := 1; val <= n; val++ {
seq[pos] = val
dfs(pos + 1)
}
}

dfs(0)

for i := 1; i <= n; i++ {
fmt.Fprintf(out, "%d", ans[i]%mod)
if i < n {
fmt.Fprint(out, " ")
}
}
fmt.Fprintln(out)
}

func goodSequence(p []int) bool {
n := len(p)
seen := make(map[int]bool)
for _, v := range p {
seen[v] = true
}
for k := range seen {
if k <= 1 {
continue
}
ok := false
for j := 0; j < n; j++ {
if p[j] == k {
for i := 0; i < j; i++ {
if p[i] == k-1 {
ok = true
break
}
}
if ok {
break
}
}
}
if !ok {
return false
}
}
return true
}

