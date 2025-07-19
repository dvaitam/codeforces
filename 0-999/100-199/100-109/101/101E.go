package main

import (
   "bufio"
   "fmt"
   "os"
)

var (
   n, m, p int
   X, Y    []int
   F, G    []int
   ans     []byte
   sum     int64
)

func f(a, b int) int {
   return (X[a] + Y[b]) % p
}

func solve(x0, y0, x1, y1 int) {
   if x0 == x1 {
       sum += int64(f(x0, y0))
       for j := y0 + 1; j <= y1; j++ {
           sum += int64(f(x0, j))
           ans = append(ans, 'S')
       }
       return
   }
   for j := y0; j <= y1; j++ {
       F[j] = 0
       G[j] = 0
   }
   mid := (x0 + x1) >> 1
   // compute F for left half
   for i := x0; i <= mid; i++ {
       F[y0] += f(i, y0)
       for j := y0 + 1; j <= y1; j++ {
           if F[j-1] > F[j] {
               F[j] = F[j-1]
           }
           F[j] += f(i, j)
       }
   }
   // compute G for right half
   for i := x1; i >= mid+1; i-- {
       G[y1] += f(i, y1)
       for j := y1 - 1; j >= y0; j-- {
           if G[j+1] > G[j] {
               G[j] = G[j+1]
           }
           G[j] += f(i, j)
       }
   }
   bst := y1
   for j := y0; j < y1; j++ {
       if F[j]+G[j] > F[bst]+G[bst] {
           bst = j
       }
   }
   solve(x0, y0, mid, bst)
   ans = append(ans, 'C')
   solve(mid+1, bst, x1, y1)
}

func main() {
   reader := bufio.NewReader(os.Stdin)
   fmt.Fscan(reader, &n, &m, &p)
   X = make([]int, n)
   Y = make([]int, m)
   for i := 0; i < n; i++ {
       fmt.Fscan(reader, &X[i])
   }
   for j := 0; j < m; j++ {
       fmt.Fscan(reader, &Y[j])
   }
   F = make([]int, m)
   G = make([]int, m)
   ans = make([]byte, 0, n+m)
   solve(0, 0, n-1, m-1)
   writer := bufio.NewWriter(os.Stdout)
   defer writer.Flush()
   fmt.Fprintln(writer, sum)
   if len(ans) > 0 {
       writer.Write(ans)
   }
   writer.WriteByte('\n')
}
