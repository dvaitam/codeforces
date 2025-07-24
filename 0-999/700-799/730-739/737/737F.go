package main

import (
   "bufio"
   "fmt"
   "os"
)

func main() {
   in := bufio.NewReader(os.Stdin)
   var n, a, b int
   if _, err := fmt.Fscan(in, &n, &a, &b); err != nil {
       return
   }
   s := make([]int, n)
   for i := 0; i < n; i++ {
       fmt.Fscan(in, &s[i])
   }
   // simulate stack-sort via one stack
   raw := make([]int, 0, 2*n)
   stack := make([]int, 0, n)
   need := 1
   // process dirty stack: top is s[n-1]
   for i := n - 1; i >= 0; i-- {
       // push one
       stack = append(stack, s[i])
       raw = append(raw, +1)
       // pop as long as top matches need
       for len(stack) > 0 && stack[len(stack)-1] == need {
           stack = stack[:len(stack)-1]
           raw = append(raw, -1)
           need++
       }
   }
   // after processing, check if sorted
   if need != n+1 {
       fmt.Println("NO")
       return
   }
   // group raw ops into operations within capacity
   type op struct{ t, c int }
   ops := make([]op, 0, len(raw))
   // iterate raw
   for i := 0; i < len(raw); {
       typ := raw[i]
       j := i + 1
       for j < len(raw) && raw[j] == typ {
           j++
       }
       cnt := j - i
       if typ > 0 {
           // push operations, type 1, capacity a
           for cnt > 0 {
               k := cnt
               if k > a {
                   k = a
               }
               ops = append(ops, op{1, k})
               cnt -= k
           }
       } else {
           // pop operations, type 2, capacity b
           for cnt > 0 {
               k := cnt
               if k > b {
                   k = b
               }
               ops = append(ops, op{2, k})
               cnt -= k
           }
       }
       i = j
   }
   // output result
   w := bufio.NewWriter(os.Stdout)
   defer w.Flush()
   fmt.Fprintln(w, "YES")
   fmt.Fprintln(w, len(ops))
   for _, o := range ops {
       fmt.Fprintf(w, "%d %d\n", o.t, o.c)
   }
}
