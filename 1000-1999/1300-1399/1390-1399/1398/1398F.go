package main

import (
   "bufio"
   "bytes"
   "fmt"
   "os"
)

func main() {
   reader := bufio.NewReader(os.Stdin)
   var n int
   fmt.Fscan(reader, &n)
   sBytes := make([]byte, n)
   // read string
   var s string
   fmt.Fscan(reader, &s)
   for i := 0; i < n; i++ {
       sBytes[i] = s[i]
   }
   // compute maxlen: maximum run of same (0/1) or '?' starting at i
   maxlen := make([]int, n)
   var len0, len1 int
   for i := n - 1; i >= 0; i-- {
       if sBytes[i] != '1' {
           len0++
       } else {
           len0 = 0
       }
       if sBytes[i] != '0' {
           len1++
       } else {
           len1 = 0
       }
       if len0 > len1 {
           maxlen[i] = len0
       } else {
           maxlen[i] = len1
       }
   }
   // buckets for positions by maxlen
   head := make([]int, n+1)
   for i := range head {
       head[i] = -1
   }
   nextInBucket := make([]int, n)
   for i := 0; i < n; i++ {
       l := maxlen[i]
       nextInBucket[i] = head[l]
       head[l] = i
   }
   // DSU-based successor: maintain alive positions by deletions
   // parent for DSU find: next candidate
   parent := make([]int, n+1)
   for i := 0; i <= n; i++ {
       parent[i] = i
   }
   var find func(int) int
   find = func(u int) int {
       if parent[u] != u {
           parent[u] = find(parent[u])
       }
       return parent[u]
   }
   // remove positions with maxlen < x: union(i, i+1)
   // initial remove maxlen == 0
   for i := head[0]; i != -1; i = nextInBucket[i] {
       parent[find(i)] = find(i + 1)
   }
   ans := make([]int, n+1)
   // process x from 1 to n
   for x := 1; x <= n; x++ {
       // simulate greedy on alive positions (maxlen >= x)
       cnt := 0
       p := 0
       for {
           j := find(p)
           if j >= n {
               break
           }
           cnt++
           p = j + x
       }
       ans[x] = cnt
       // remove positions where maxlen == x (for next x+1)
       for i := head[x]; i != -1; i = nextInBucket[i] {
           parent[find(i)] = find(i + 1)
       }
   }
   // output
   var buf bytes.Buffer
   writer := bufio.NewWriter(os.Stdout)
   for i := 1; i <= n; i++ {
       buf.WriteString(fmt.Sprint(ans[i]))
       if i < n {
           buf.WriteByte(' ')
       }
   }
   buf.WriteByte('\n')
   writer.Write(buf.Bytes())
   writer.Flush()
}
