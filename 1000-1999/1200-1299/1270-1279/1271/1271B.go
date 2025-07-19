package main

import (
   "bufio"
   "fmt"
   "os"
)

func solve() {
   var n int
   in := bufio.NewReader(os.Stdin)
   out := bufio.NewWriter(os.Stdout)
   defer out.Flush()

   if _, err := fmt.Fscan(in, &n); err != nil {
       return
   }
   var s string
   fmt.Fscan(in, &s)
   // convert to byte slice
   orig := []byte(s)

   try := func(target byte) ([]int, bool) {
       a := make([]byte, n)
       copy(a, orig)
       ops := make([]int, 0, n)
       for i := 0; i < n-1; i++ {
           if a[i] != target {
               // flip i and i+1
               ops = append(ops, i+1)
               a[i] = target
               // toggle next
               if a[i+1] == 'W' {
                   a[i+1] = 'B'
               } else {
                   a[i+1] = 'W'
               }
           }
       }
       // check all equal target
       for i := 0; i < n; i++ {
           if a[i] != target {
               return nil, false
           }
       }
       return ops, true
   }

   // try all W
   if ops, ok := try('W'); ok {
       fmt.Fprintln(out, len(ops))
       for i, v := range ops {
           if i > 0 {
               out.WriteByte(' ')
           }
           fmt.Fprint(out, v)
       }
       fmt.Fprintln(out)
       return
   }
   // try all B
   if ops, ok := try('B'); ok {
       fmt.Fprintln(out, len(ops))
       for i, v := range ops {
           if i > 0 {
               out.WriteByte(' ')
           }
           fmt.Fprint(out, v)
       }
       fmt.Fprintln(out)
       return
   }
   // impossible
   fmt.Fprintln(out, -1)
}

func main() {
   solve()
}
