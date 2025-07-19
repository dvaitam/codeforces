package main

import (
   "bufio"
   "fmt"
   "os"
   "strconv"
)

func main() {
   reader := bufio.NewReader(os.Stdin)
   writer := bufio.NewWriter(os.Stdout)
   defer writer.Flush()
   var n, k int
   if _, err := fmt.Fscan(reader, &n, &k); err != nil {
       return
   }
   a := make([]int, n+2)
   useb := make([]bool, n+2)
   for i := 1; i <= n; i++ {
       var s string
       fmt.Fscan(reader, &s)
       if s == "?" {
           useb[i] = true
       } else {
           v, _ := strconv.Atoi(s)
           a[i] = v
       }
   }
   const INF = 1001000000
   for i := 1; i <= k; i++ {
       a[0] = -INF
       a[n+1] = INF
       sg := 0
       var j int
       for j = i; j <= n; j += k {
           if !useb[j] {
               cnt := (j - sg) / k
               if a[j]-cnt < a[sg] {
                   fmt.Fprintln(writer, "Incorrect sequence")
                   return
               }
               if a[j]*2 < cnt {
                   for l := j - k; l > sg; l -= k {
                       a[l] = a[l+k] - 1
                   }
               } else if -a[sg]*2 < cnt {
                   for l := sg + k; l < j; l += k {
                       a[l] = a[l-k] + 1
                   }
               } else {
                   half := (j + sg) / k / 2
                   for l := j - k; l > sg; l -= k {
                       a[l] = l/k - half
                   }
               }
               sg = j
           }
       }
       cnt := (j - sg) / k
       if -a[sg]*2 < cnt {
           for l := sg + k; l < j; l += k {
               a[l] = a[l-k] + 1
           }
       } else {
           if sg == 0 {
               sg = i - k
           }
           half := (j + sg) / k / 2
           for l := sg + k; l < j; l += k {
               a[l] = l/k - half
           }
       }
   }
   for i := 1; i <= n; i++ {
       if i > 1 {
           writer.WriteByte(' ')
       }
       writer.WriteString(strconv.Itoa(a[i]))
   }
   writer.WriteByte('\n')
}
