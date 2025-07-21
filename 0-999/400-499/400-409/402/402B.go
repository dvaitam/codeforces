package main

import (
   "bufio"
   "fmt"
   "os"
)

func main() {
   reader := bufio.NewReader(os.Stdin)
   writer := bufio.NewWriter(os.Stdout)
   defer writer.Flush()

   var n, k int
   if _, err := fmt.Fscan(reader, &n, &k); err != nil {
       return
   }
   a := make([]int, n)
   for i := 0; i < n; i++ {
       fmt.Fscan(reader, &a[i])
   }
   // For arithmetic progression a[i] = h + i*k, compute b[i] = a[i] - i*k
   freq := make(map[int]int)
   for i := 0; i < n; i++ {
       bi := a[i] - i*k
       if bi >= 1 {
           freq[bi]++
       }
   }
   // choose h maximizing frequency
   bestH := 1
   bestCnt := 0
   for h, cnt := range freq {
       if cnt > bestCnt {
           bestCnt = cnt
           bestH = h
       }
   }
   // record operations for mismatches
   ops := make([]string, 0, n)
   for i := 0; i < n; i++ {
       desired := bestH + i*k
       if a[i] != desired {
           if a[i] < desired {
               ops = append(ops, fmt.Sprintf("+ %d %d", i+1, desired-a[i]))
           } else {
               ops = append(ops, fmt.Sprintf("- %d %d", i+1, a[i]-desired))
           }
       }
   }
   // output
   fmt.Fprintln(writer, len(ops))
   for _, op := range ops {
       fmt.Fprintln(writer, op)
   }
}
