package main

import (
   "bufio"
   "fmt"
   "os"
   "sort"
)

func main() {
   reader := bufio.NewReader(os.Stdin)
   writer := bufio.NewWriter(os.Stdout)
   defer writer.Flush()

   var n, k int
   if _, err := fmt.Fscan(reader, &n, &k); err != nil {
       return
   }
   var s string
   fmt.Fscan(reader, &s)
   if k > n || k > 13 {
       fmt.Fprintln(writer, -1)
       return
   }
   // convert string to byte slice and sort
   b := []byte(s)
   sort.Slice(b, func(i, j int) bool { return b[i] < b[j] })
   const inf = 1<<60
   mn := inf
   found := false
   for i := 0; i < n; i++ {
       cnt := 1
       prev := b[i]
       sum := int(b[i]-'a'+1)
       for j := i + 1; j < n && cnt < k; j++ {
           if b[j]-prev >= 2 {
               sum += int(b[j] - 'a' + 1)
               prev = b[j]
               cnt++
           }
       }
       if cnt == k {
           found = true
           if sum < mn {
               mn = sum
           }
       }
   }
   if !found {
       fmt.Fprintln(writer, -1)
   } else {
       fmt.Fprintln(writer, mn)
   }
}
