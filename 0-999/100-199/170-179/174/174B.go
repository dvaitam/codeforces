package main

import (
   "bufio"
   "fmt"
   "os"
)

func main() {
   reader := bufio.NewReader(os.Stdin)
   var s string
   if _, err := fmt.Fscan(reader, &s); err != nil {
       return
   }
   n := len(s)
   // collect positions of dots
   d := make([]int, 0)
   for i := 0; i < n; i++ {
       if s[i] == '.' {
           d = append(d, i)
       }
   }
   m := len(d)
   if m == 0 {
       fmt.Println("NO")
       return
   }
   // first name length must be 1..8
   if d[0] < 1 || d[0] > 8 {
       fmt.Println("NO")
       return
   }
   // ext lengths for each segment
   l := make([]int, m)
   // determine ext lengths for segments except last
   for k := 0; k < m-1; k++ {
       delta := d[k+1] - d[k]
       low := delta - 9
       if low < 1 {
           low = 1
       }
       high := delta - 2
       if high > 3 {
           high = 3
       }
       if low > high {
           fmt.Println("NO")
           return
       }
       // pick smallest valid ext length
       found := false
       for ll := low; ll <= high; ll++ {
           ok := true
           for j := 1; j <= ll; j++ {
               if s[d[k]+j] == '.' {
                   ok = false
                   break
               }
           }
           if ok {
               l[k] = ll
               found = true
               break
           }
       }
       if !found {
           fmt.Println("NO")
           return
       }
   }
   // last segment ext must reach end
   lastL := n - 1 - d[m-1]
   if lastL < 1 || lastL > 3 {
       fmt.Println("NO")
       return
   }
   for j := 1; j <= lastL; j++ {
       if s[d[m-1]+j] == '.' {
           fmt.Println("NO")
           return
       }
   }
   l[m-1] = lastL

   // output segments
   writer := bufio.NewWriter(os.Stdout)
   defer writer.Flush()
   fmt.Fprintln(writer, "YES")
   start := 0
   for k := 0; k < m; k++ {
       end := d[k] + l[k]
       fmt.Fprintln(writer, s[start:end+1])
       start = end + 1
   }
}
