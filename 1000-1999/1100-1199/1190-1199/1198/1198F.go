package main

import (
   "bufio"
   "fmt"
   "math/rand"
   "os"
   "time"
)

func gcd(a, b int) int {
   for b != 0 {
       a, b = b, a%b
   }
   return a
}

func main() {
   reader := bufio.NewReader(os.Stdin)
   writer := bufio.NewWriter(os.Stdout)
   defer writer.Flush()

   var n int
   if _, err := fmt.Fscan(reader, &n); err != nil {
       return
   }
   a := make([]int, n)
   for i := 0; i < n; i++ {
       fmt.Fscan(reader, &a[i])
   }

   rand.Seed(time.Now().UnixNano())
   primeSet := make(map[int]struct{})
   pList := make([]int, 0)
   // collect candidate primes
   for t := 0; t < 25; t++ {
       id := rand.Intn(n)
       x := a[id]
       j := 2
       for j*j <= x {
           if x%j == 0 {
               if _, ok := primeSet[j]; !ok {
                   occ := 0
                   for i := 0; i < n; i++ {
                       if a[i]%j == 0 {
                           occ++
                       }
                   }
                   if occ*2 >= n {
                       primeSet[j] = struct{}{}
                       pList = append(pList, j)
                   }
               }
               for x%j == 0 {
                   x /= j
               }
           }
           j++
       }
       if x > 1 {
           if _, ok := primeSet[x]; !ok {
               occ := 0
               for i := 0; i < n; i++ {
                   if a[i]%x == 0 {
                       occ++
                   }
               }
               if occ*2 >= n {
                   primeSet[x] = struct{}{}
                   pList = append(pList, x)
               }
           }
       }
   }

   m := len(pList)
   // build non-divisible index lists
   v := make([][]int, m)
   for i, p := range pList {
       for j := 0; j < n; j++ {
           if a[j]%p != 0 {
               v[i] = append(v[i], j)
           }
       }
       if len(v[i]) <= 1 {
           fmt.Fprintln(writer, "NO")
           return
       }
   }

   // prepare graph structures
   h := make([][]int, n)
   vis := make([]bool, n)
   color := make([]int, n)
   nodeSet := make(map[int]struct{})
   const E = 6000
   // trial
   for t := 0; t < E; t++ {
       // clear from previous
       for x := range nodeSet {
           h[x] = h[x][:0]
           vis[x] = false
       }
       nodeSet = make(map[int]struct{})

       scc := true
       // add random edges for each prime constraint
       for i := 0; i < m; i++ {
           sz := len(v[i])
           var x, y int
           for x == y {
               x = v[i][rand.Intn(sz)]
               y = v[i][rand.Intn(sz)]
           }
           h[x] = append(h[x], y)
           h[y] = append(h[y], x)
           nodeSet[x] = struct{}{}
           nodeSet[y] = struct{}{}
       }
       // dfs to check bipartiteness
       var dfs func(int, int)
       dfs = func(u, c int) {
           vis[u] = true
           color[u] = c
           for _, vtx := range h[u] {
               if vis[vtx] {
                   if color[vtx] == c {
                       scc = false
                   }
               } else {
                   dfs(vtx, c^1)
               }
               if !scc {
                   return
               }
           }
       }
       for x := range nodeSet {
           if !vis[x] {
               dfs(x, 0)
               if !scc {
                   break
               }
           }
       }
       if !scc {
           continue
       }
       // count colors in constrained nodes
       c0, c1 := 0, 0
       for x := range nodeSet {
           if color[x] == 0 {
               c0++
           } else {
               c1++
           }
       }
       half := n / 2
       if c0 <= half && c1 <= half {
           // assign remaining
           for i := 0; i < n; i++ {
               if _, ok := nodeSet[i]; !ok {
                   if c0 < half {
                       color[i] = 0
                       c0++
                   } else {
                       color[i] = 1
                       c1++
                   }
               }
           }
           // verify gcd
           g0, g1 := 0, 0
           for i := 0; i < n; i++ {
               if color[i] == 0 {
                   g0 = gcd(g0, a[i])
               } else {
                   g1 = gcd(g1, a[i])
               }
           }
           if g0 == 1 && g1 == 1 {
               fmt.Fprintln(writer, "YES")
               for i := 0; i < n; i++ {
                   fmt.Fprintf(writer, "%d ", color[i]+1)
               }
               fmt.Fprintln(writer)
               return
           }
       }
   }
   fmt.Fprintln(writer, "NO")
}
