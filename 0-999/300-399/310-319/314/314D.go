package main

import (
   "bufio"
   "fmt"
   "os"
   "sort"
)

type point struct {
   X, Y int64
}

func max(a, b int64) int64 {
   if a > b {
       return a
   }
   return b
}

func min(a, b int64) int64 {
   if a < b {
       return a
   }
   return b
}

var rd = bufio.NewReader(os.Stdin)

func readInt() int64 {
   var x int64
   var c byte
   var sign int64 = 1
   for {
       b, err := rd.ReadByte()
       if err != nil {
           return 0
       }
       c = b
       if c == '-' || (c >= '0' && c <= '9') {
           break
       }
   }
   if c == '-' {
       sign = -1
   } else {
       x = int64(c - '0')
   }
   for {
       b, err := rd.ReadByte()
       if err != nil || b < '0' || b > '9' {
           break
       }
       x = x*10 + int64(b-'0')
   }
   return x * sign
}

func main() {
   N := int(readInt())
   A := make([]point, N)
   for i := 0; i < N; i++ {
       p := readInt()
       q := readInt()
       A[i].X = p + q
       A[i].Y = p - q
   }
   // sort points by X coordinate
   sort.Slice(A, func(i, j int) bool { return A[i].X < A[j].X })
   // prefix on Y
   INF := int64(1e18)
   LMax := make([]int64, N+1)
   LMin := make([]int64, N+1)
   RMax := make([]int64, N+1)
   RMin := make([]int64, N+1)
   LMax[0] = -INF
   LMin[0] = INF
   for i := 0; i < N; i++ {
       LMax[i+1] = max(LMax[i], A[i].Y)
       LMin[i+1] = min(LMin[i], A[i].Y)
   }
   RMax[N] = -INF
   RMin[N] = INF
   for i := N - 1; i >= 0; i-- {
       RMax[i] = max(RMax[i+1], A[i].Y)
       RMin[i] = min(RMin[i+1], A[i].Y)
   }
   Ans := INF
   j := 1
   for i := 0; i < N; i++ {
       if j < i+1 {
           j = i + 1
       }
       // function P: span in X
       P := func(l, r int) int64 {
           return A[r-1].X - A[l].X
       }
       // function Q: span in Y
       Q := func(l, r int) int64 {
           return max(LMax[l], RMax[r]) - min(LMin[l], RMin[r])
       }
       // initial at current j
       if j <= N {
           p := P(i, j)
           q := Q(i, j)
           if m := max(p, q); m < Ans {
               Ans = m
           }
       }
       // increase j while P < Q
       for j < N {
           p := P(i, j)
           q := Q(i, j)
           if p >= q {
               break
           }
           j++
           p = P(i, j)
           q = Q(i, j)
           if m := max(p, q); m < Ans {
               Ans = m
           }
       }
   }
   // answer is Ans/2.0
   real := float64(Ans) / 2.0
   writer := bufio.NewWriter(os.Stdout)
   defer writer.Flush()
   fmt.Fprintf(writer, "%.10f\n", real)
}
