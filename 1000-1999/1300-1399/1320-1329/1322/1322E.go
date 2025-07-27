package main

import (
   "bufio"
   "fmt"
   "os"
)

func median3(a, b, c int) int {
   // return median of three ints
   if a > b {
       if b > c {
           return b
       } else if a > c {
           return c
       }
       return a
   }
   // a <= b
   if a > c {
       return a
   } else if b > c {
       return c
   }
   return b
}

func main() {
   reader := bufio.NewReader(os.Stdin)
   writer := bufio.NewWriter(os.Stdout)
   defer writer.Flush()

   var n int
   fmt.Fscan(reader, &n)
   a0 := make([]int, n)
   for i := 0; i < n; i++ {
       fmt.Fscan(reader, &a0[i])
   }
   if n <= 2 {
       // no change ever for n<=2
       fmt.Fprintln(writer, 0)
       for i, v := range a0 {
           if i > 0 {
               writer.WriteByte(' ')
           }
           fmt.Fprint(writer, v)
       }
       writer.WriteByte('\n')
       return
   }
   // first iteration
   b1 := make([]int, n)
   b1[0] = a0[0]
   b1[n-1] = a0[n-1]
   changed := false
   for i := 1; i < n-1; i++ {
       v := median3(a0[i-1], a0[i], a0[i+1])
       b1[i] = v
       if v != a0[i] {
           changed = true
       }
   }
   if !changed {
       // already stable
       fmt.Fprintln(writer, 0)
       for i, v := range a0 {
           if i > 0 {
               writer.WriteByte(' ')
           }
           fmt.Fprint(writer, v)
       }
       writer.WriteByte('\n')
       return
   }
   // prepare BFS from b1 stable positions
   safe := make([]bool, n)
   for i := 0; i < n; i++ {
       if i == 0 || i == n-1 {
           safe[i] = true
       } else {
           bi := b1[i]
           if (bi > b1[i-1] && bi > b1[i+1]) || (bi < b1[i-1] && bi < b1[i+1]) {
               safe[i] = false
           } else {
               safe[i] = true
           }
       }
   }
   // time and final value
   time2 := make([]int, n)
   finalVal := make([]int, n)
   const INF = 1_000_000_000
   for i := 0; i < n; i++ {
       if safe[i] {
           time2[i] = 0
           finalVal[i] = b1[i]
       } else {
           time2[i] = -1
       }
   }
   // BFS queue
   q := make([]int, 0, n)
   for i := 0; i < n; i++ {
       if safe[i] {
           q = append(q, i)
       }
   }
   // BFS
   maxT := 0
   for head := 0; head < len(q); head++ {
       u := q[head]
       t := time2[u]
       for _, d := range []int{-1, 1} {
           v := u + d
           if v < 0 || v >= n {
               continue
           }
           // only propagate into unstable
           if time2[v] == -1 {
               time2[v] = t + 1
               finalVal[v] = finalVal[u]
               if time2[v] > maxT {
                   maxT = time2[v]
               }
               q = append(q, v)
           } else if time2[v] == t+1 {
               // collision: choose based on b1 initial extremum type
               // if v is peak in b1, pick max; if valley, pick min
               if v > 0 && v < n-1 {
                   bi := b1[v]
                   if (bi > b1[v-1] && bi > b1[v+1]) {
                       // peak
                       if finalVal[u] > finalVal[v] {
                           finalVal[v] = finalVal[u]
                       }
                   } else if (bi < b1[v-1] && bi < b1[v+1]) {
                       // valley
                       if finalVal[u] < finalVal[v] {
                           finalVal[v] = finalVal[u]
                       }
                   }
               }
           }
       }
   }
   // total c = 1 (first) + maxT
   c := 1 + maxT
   fmt.Fprintln(writer, c)
   // print finalVal
   for i := 0; i < n; i++ {
       if i > 0 {
           writer.WriteByte(' ')
       }
       fmt.Fprint(writer, finalVal[i])
   }
   writer.WriteByte('\n')
}
