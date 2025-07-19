package main

import (
   "bufio"
   "fmt"
   "os"
   "sort"
)

type Point struct { x, y int }
type Edge struct { to, nxt int }

const MAXN = 1005

var (
   n, k   int
   p       [MAXN]Point
   dis     [MAXN][MAXN]int
   D       []int
   tot     int
   MaxDist int

   b   []Edge
   fst [MAXN]int
   T   int

   cd    [MAXN]int
   use   [MAXN]int
   q     []int
   fro   int
   ans   bool
   hFlag bool

   reader = bufio.NewReader(os.Stdin)
   writer = bufio.NewWriter(os.Stdout)
)

func build(f, t int) {
   b = append(b, Edge{to: t, nxt: fst[f]})
   fst[f] = len(b) - 1
}

func sq(x int) int { return x * x }

func initData() {
   fmt.Fscan(reader, &n, &k)
   for i := 1; i <= n; i++ {
       fmt.Fscan(reader, &p[i].x, &p[i].y)
   }
   D = D[:0]
   MaxDist = 0
   for i := 1; i < n; i++ {
       for j := i + 1; j <= n; j++ {
           d := sq(p[i].x-p[j].x) + sq(p[i].y-p[j].y)
           dis[i][j] = d
           if d > MaxDist {
               MaxDist = d
           }
           D = append(D, d)
       }
   }
   D = append(D, MaxDist+1)
   sort.Ints(D)
   // unique
   tot = 0
   for i := 0; i < len(D); i++ {
       if tot == 0 || D[i] != D[tot-1] {
           D[tot] = D[i]
           tot++
       }
   }
   D = D[:tot]
}

func printAns() {
   rem := k
   for i := n; i >= 1 && rem > 0; i-- {
       if use[i] > 0 {
           fmt.Fprintf(writer, "%d ", i)
           rem--
       }
   }
   for i := n; i >= 1 && rem > 0; i-- {
       if use[i] == 0 {
           fmt.Fprintf(writer, "%d ", i)
           rem--
       }
   }
}

func dfs(x, u int) {
   if u > k || ans {
       return
   }
   if x >= fro {
       ans = true
       if hFlag {
           printAns()
           hFlag = false
       }
       return
   }
   f := q[x]
   if use[f] > 0 {
       dfs(x+1, u)
   } else {
       a := u
       // include neighbors
       for i := fst[f]; i != 0; i = b[i].nxt {
           v := b[i].to
           if use[v] == 0 {
               a++
           }
           use[v]++
       }
       dfs(x+1, a)
       // rollback
       for i := fst[f]; i != 0; i = b[i].nxt {
           use[b[i].to]--
       }
       // try select f
       if cd[f] != 1 && u+1 < a {
           use[f] = 1
           dfs(x+1, u+1)
           use[f] = 0
       }
   }
}

func check(x int) bool {
   // reset
   for i := 1; i <= n; i++ {
       fst[i] = 0
       cd[i] = 0
       use[i] = 0
   }
   b = b[:0]
   fro = 0
   // build graph
   for i := 1; i < n; i++ {
       for j := i + 1; j <= n; j++ {
           if dis[i][j] >= x {
               build(i, j)
               build(j, i)
               if cd[i] == 0 {
                   q = append(q, i)
                   fro++
               }
               if cd[j] == 0 {
                   q = append(q, j)
                   fro++
               }
               cd[i]++
               cd[j]++
           }
       }
   }
   ans = false
   dfs(0, 0)
   return ans
}

func work() {
   initData()
   l, r := 0, tot-1
   for r-l > 1 {
       mid := (l + r) >> 1
       q = q[:0]
       if check(D[mid]) {
           r = mid
       } else {
           l = mid
       }
   }
   hFlag = true
   q = q[:0]
   check(D[r])
}

func main() {
   defer writer.Flush()
   work()
}
