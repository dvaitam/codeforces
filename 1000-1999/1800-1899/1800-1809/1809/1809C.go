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
   var t int
   if _, err := fmt.Fscan(reader, &t); err != nil {
       return
   }
   for ; t > 0; t-- {
       var n, k int
       fmt.Fscan(reader, &n, &k)
       // initialize array with large negatives
       a := make([]int, n)
       for i := 0; i < n; i++ {
           a[i] = -1000
       }
       // find maximum temp such that temp*(temp+1)/2 <= k
       temp := 0
       for i := 1; i <= n; i++ {
           if i*(i+1)/2 <= k {
               temp = i
           }
       }
       // remaining positive subarrays after full prefixes
       remK := k - temp*(temp+1)/2
       // set first temp values as positives 2,3,...
       for i := 0; i < temp; i++ {
           a[i] = i + 2
       }
       if remK > 0 {
           // sum from index remK-1 to temp-1
           sum := 0
           for i := remK - 1; i < temp; i++ {
               sum += a[i]
           }
           // place a negative number to adjust further subarrays
           a[temp] = -(sum - 1)
       }
       // output
       for i, v := range a {
           if i > 0 {
               writer.WriteByte(' ')
           }
           fmt.Fprint(writer, v)
       }
       writer.WriteByte('\n')
   }
}
