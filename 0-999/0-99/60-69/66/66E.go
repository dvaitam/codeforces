package main

import (
   "bufio"
   "fmt"
   "os"
)

func main() {
   reader := bufio.NewReader(os.Stdin)
   var n int
   fmt.Fscan(reader, &n)
   a := make([]int64, n)
   b := make([]int64, n)
   for i := 0; i < n; i++ {
       fmt.Fscan(reader, &a[i])
   }
   for i := 0; i < n; i++ {
       fmt.Fscan(reader, &b[i])
   }
   // forward net gains
   c := make([]int64, n)
   for i := 0; i < n; i++ {
       c[i] = a[i] - b[i]
   }
   // reverse net gains: consume b_prev
   cr := make([]int64, n)
   for i := 0; i < n; i++ {
       prev := i - 1
       if prev < 0 {
           prev = n - 1
       }
       cr[i] = a[i] - b[prev]
   }
   // compute valid starts for forward and reverse
   okf := check(c)
  okr := check(cr)
   // collect union
   out := make([]int, 0, n)
   for i := 0; i < n; i++ {
       if okf[i] || okr[i] {
           out = append(out, i+1)
       }
   }
   // print
   w := bufio.NewWriter(os.Stdout)
   defer w.Flush()
   fmt.Fprintln(w, len(out))
   if len(out) > 0 {
       for i, v := range out {
           if i > 0 {
               w.WriteString(" ")
           }
           fmt.Fprint(w, v)
       }
       w.WriteByte('\n')
   }
}

// check returns a bool slice of length n indicating valid start positions
func check(c []int64) []bool {
   n := len(c)
   // prefix sums size 2n+1
   S := make([]int64, 2*n+1)
   for i := 1; i <= 2*n; i++ {
       S[i] = S[i-1] + c[(i-1)%n]
   }
   ok := make([]bool, n)
   // deque for indices of S, maintain increasing S values
   D := make([]int, 2*n+1)
   head, tail := 0, 0
   for j := 1; j <= 2*n; j++ {
       // push j: maintain S[D[*]] increasing
       for head < tail && S[D[tail-1]] >= S[j] {
           tail--
       }
       D[tail] = j
       tail++
       // remove indices out of window [j-(n-1), j]
       start := j - (n - 1)
       if head < tail && D[head] < start {
           head++
       }
       // when window full, check start
       if start >= 1 {
           // minimal S in window is S[D[head]]
           if S[D[head]]-S[start-1] >= 0 {
               idx := start - 1 // 0-based start index
               if idx < n {
                   ok[idx] = true
               }
           }
       }
   }
   return ok
}
