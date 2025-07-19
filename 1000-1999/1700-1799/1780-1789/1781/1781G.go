package main

import (
   "bufio"
   "fmt"
   "os"
)

type nodeInfo struct { f, s int }

var (
   lc, rc []int
   ans []int
   p []nodeInfo
   reader *bufio.Reader
   writer *bufio.Writer
)

func dfs(x int) {
   if lc[x] == 0 {
      p[x] = nodeInfo{1, 0}
      return
   }
   if rc[x] == 0 {
      dfs(lc[x])
      c := lc[x]
      if p[c].s == 0 {
         ans[c] ^= 1
         p[x] = nodeInfo{1 - p[c].f, 1}
      } else {
         if p[c].f == 2 {
            ans[lc[c]] ^= 1
            ans[c] ^= 1
            p[x] = nodeInfo{1, 1}
         } else {
            ans[c] ^= 1
            p[x] = nodeInfo{1 - p[c].f, 1}
         }
      }
      return
   }
   dfs(lc[x])
   dfs(rc[x])
   // handle f == 2 cases
   if p[lc[x]].f == 2 {
      ans[lc[lc[x]]] ^= 1
      p[lc[x]] = nodeInfo{0, 0}
   }
   if p[rc[x]].f == 2 {
      ans[lc[rc[x]]] ^= 1
      p[rc[x]] = nodeInfo{0, 0}
   }
   // both s == 0
   if p[lc[x]].s == 0 && p[rc[x]].s == 0 {
      ans[lc[x]] ^= 1
      ans[rc[x]] ^= 1
      val := 1 - p[lc[x]].f - p[rc[x]].f
      p[x] = nodeInfo{val, 1}
      return
   }
   for f1 := 0; f1 < 2; f1++ {
      for f2 := 0; f2 < 2; f2++ {
         if p[lc[x]].s == 0 && f1 == 0 {
            continue
         }
         if p[rc[x]].s == 0 && f2 == 0 {
            continue
         }
         val := 1
         if f1 == 1 {
            val -= p[lc[x]].f
         } else {
            val += p[lc[x]].f
         }
         if f2 == 1 {
            val -= p[rc[x]].f
         } else {
            val += p[rc[x]].f
         }
         if (val == 0 || val == 1) && (f1 == 1 || f2 == 1) {
            if f1 == 1 {
               ans[lc[x]] ^= 1
            }
            if f2 == 1 {
               ans[rc[x]] ^= 1
            }
            p[x] = nodeInfo{val, 1}
            return
         }
      }
   }
}

func solve() {
   var n int
   fmt.Fscan(reader, &n)
   lc = make([]int, n+3)
   rc = make([]int, n+3)
   ans = make([]int, n+3)
   p = make([]nodeInfo, n+3)
   for i := 2; i <= n; i++ {
      var x int
      fmt.Fscan(reader, &x)
      if lc[x] == 0 {
         lc[x] = i
      } else {
         rc[x] = i
      }
      // interactive response
      if i == 4 && lc[2] != 0 && rc[2] != 0 {
         writer.WriteString("2\n")
      } else {
         writer.WriteString(fmt.Sprintf("%d\n", i&1))
      }
      writer.Flush()
   }
   dfs(1)
   for i := 1; i <= n; i++ {
      if lc[i] != 0 {
         ans[lc[i]] ^= ans[i]
      }
      if rc[i] != 0 {
         ans[rc[i]] ^= ans[i]
      }
      if ans[i] != 0 {
         writer.WriteByte('w')
      } else {
         writer.WriteByte('b')
      }
   }
   writer.WriteByte('\n')
   writer.Flush()
}

func main() {
   reader = bufio.NewReader(os.Stdin)
   writer = bufio.NewWriter(os.Stdout)
   var T int
   fmt.Fscan(reader, &T)
   for T > 0 {
      solve()
      T--
   }
}
