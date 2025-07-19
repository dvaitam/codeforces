package main

import (
   "bufio"
   "fmt"
   "os"
)

const N = 200500

var (
   a   = make([]int, N)
   p   = make([]int, N)
   cnt = make([]int, N)
   par = make([]int, N)
   last = make([]int, N)
)

type pair struct{ x, y int }

func find(u int) int {
   if par[u] < 0 {
       return u
   }
   par[u] = find(par[u])
   return par[u]
}

func unite(u, v int) {
   u = find(u)
   v = find(v)
   if u == v {
       return
   }
   if par[u] > par[v] {
       u, v = v, u
   }
   par[u] += par[v]
   par[v] = u
}

func add(u, v int) {
   p[u] = v
   unite(u, v)
}

func doIt(u, v, x, y int) {
   if find(u) == find(x) || find(v) == find(y) {
       // swap x and y
       x, y = y, x
   }
   add(u, x)
   add(v, y)
}

func main() {
   in := bufio.NewReader(os.Stdin)
   out := bufio.NewWriter(os.Stdout)
   defer out.Flush()
   var t, n int
   fmt.Fscan(in, &t)
   // initialize par
   for i := 0; i < N; i++ {
       par[i] = -1
   }
   for t > 0 {
       t--
       fmt.Fscan(in, &n)
       // reset counters
       for i := 1; i <= n; i++ {
           cnt[i] = 0
           par[i] = -1
           last[i] = 0
       }
       // read array
       for i := 1; i <= n; i++ {
           fmt.Fscan(in, &a[i])
           cnt[a[i]]++
       }
       odd := 0
       for i := 1; i <= n; i++ {
           if cnt[i]&1 != 0 {
               odd++
           }
       }
       mid := n/2 + 1
       if odd > 1 || (n&1 == 1 && cnt[a[mid]] == 1) {
           fmt.Fprintln(out, "NO")
           continue
       }
       fmt.Fprintln(out, "YES")
       b := 0
       if n&1 == 1 {
           for i := 1; i <= n; i++ {
               if i != mid && cnt[a[i]]&1 == 1 {
                   b = i
                   add(mid, i)
                   break
               }
           }
       }
       // pair up same values
       vec := make([]pair, 0, n)
       for i := 1; i <= n; i++ {
           if i == b {
               continue
           }
           v := a[i]
           if last[v] != 0 {
               vec = append(vec, pair{last[v], i})
               last[v] = 0
           } else {
               last[v] = i
           }
       }
       // match positions
       half := n/2
       for i := 1; i <= half; i++ {
           u := i
           v := n - i + 1
           L := len(vec)
           p1 := vec[L-1]
           vec = vec[:L-1]
           if len(vec) > 0 {
               // check bad pairing
               if (find(p1.x) == find(u) && find(p1.y) == find(v)) ||
                  (find(p1.y) == find(u) && find(p1.x) == find(v)) {
                   L = len(vec)
                   p2 := vec[L-1]
                   vec = vec[:L-1]
                   vec = append(vec, p1)
                   doIt(u, v, p2.x, p2.y)
               } else {
                   doIt(u, v, p1.x, p1.y)
               }
           } else {
               doIt(u, v, p1.x, p1.y)
           }
       }
       // output permutation
       for i := 1; i <= n; i++ {
           if i > 1 {
               out.WriteByte(' ')
           }
           fmt.Fprint(out, p[i])
       }
       out.WriteByte('\n')
   }
}
