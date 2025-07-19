package main

import (
   "bufio"
   "fmt"
   "os"
)

const maxT = 1000000

var (
   ft   []int64
   fc   []int64
   V    []int64
   xArr []int64
   tArr []int64
   nodes [][]pair
)

type pair struct {
   to int
   w  int64
}

// add adds v to BIT f at index idx
func add(f []int64, idx int, v int64) {
   for i := idx; i <= maxT; i += i & -i {
       f[i] += v
   }
}

// sum returns prefix sum of BIT f up to idx
func sum(f []int64, idx int) int64 {
   var s int64
   for i := idx; i > 0; i -= i & -i {
       s += f[i]
   }
   return s
}

// bs returns the maximum count under total time T
func bs(T int64) int64 {
   if sum(ft, maxT) <= T {
       return sum(fc, maxT)
   }
   l, r := 1, maxT
   for l < r {
       mid := (l + r) >> 1
       if sum(ft, mid) >= T {
           r = mid
       } else {
           l = mid + 1
       }
   }
   s := sum(ft, l)
   rem := s - T
   base := sum(fc, l)
   if rem%int64(l) == 0 {
       return base - rem/int64(l)
   }
   return base - rem/int64(l) - 1
}

// F performs DFS and computes answer
func F(u int, T int64) int64 {
   if T <= 0 {
       return 0
   }
   idx := int(tArr[u])
   add(ft, idx, xArr[u]*tArr[u])
   add(fc, idx, xArr[u])
   V[u] = bs(T)
   var M1, M2 int64
   for _, p := range nodes[u] {
       z := F(p.to, T - p.w*2)
       if z >= M1 {
           M2 = M1
           M1 = z
       } else if z > M2 {
           M2 = z
       }
   }
   // rollback
   add(ft, idx, -xArr[u]*tArr[u])
   add(fc, idx, -xArr[u])
   if u == 1 {
       if M1 > V[u] {
           return M1
       }
       return V[u]
   }
   if M2 > V[u] {
       return M2
   }
   return V[u]
}

func main() {
   reader := bufio.NewReader(os.Stdin)
   writer := bufio.NewWriter(os.Stdout)
   defer writer.Flush()
   var n int
   var T int64
   fmt.Fscan(reader, &n, &T)
   ft = make([]int64, maxT+1)
   fc = make([]int64, maxT+1)
   V = make([]int64, n+1)
   xArr = make([]int64, n+1)
   tArr = make([]int64, n+1)
   nodes = make([][]pair, n+1)
   for i := 1; i <= n; i++ {
       fmt.Fscan(reader, &xArr[i])
   }
   for i := 1; i <= n; i++ {
       fmt.Fscan(reader, &tArr[i])
   }
   for i := 2; i <= n; i++ {
       var p int
       var l int64
       fmt.Fscan(reader, &p, &l)
       nodes[p] = append(nodes[p], pair{to: i, w: l})
   }
   res := F(1, T)
   fmt.Fprint(writer, res)
}
