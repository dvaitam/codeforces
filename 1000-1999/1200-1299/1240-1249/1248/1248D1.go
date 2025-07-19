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

   var n int
   if _, err := fmt.Fscan(reader, &n); err != nil {
       return
   }
   var s string
   fmt.Fscan(reader, &s)

   id, hmin, h, count := 0, 0, 0, 0
   for i := 0; i < n; i++ {
       if h < hmin {
           id = i
           hmin = h
           count = 0
       }
       if h == hmin {
           count++
       }
       if s[i] == '(' {
           h++
       } else {
           h--
       }
   }
   if h != 0 {
       fmt.Fprintln(writer, 0)
       fmt.Fprintln(writer, "1 1")
       return
   }
   // rotate string by id
   s2 := s[id:] + s[:id]
   best := count
   curr1, curr2 := 0, 0
   a, b := 0, 0
   a1, a2 := 0, 0
   h = 0
   for i := 0; i < n; i++ {
       if s2[i] == '(' {
           h++
       } else {
           h--
       }
       switch {
       case h == 0:
           if curr1 > best {
               best = curr1
               a = a1
               b = i
           }
           curr1 = 0
           a1 = i + 1
       case h == 1:
           curr1++
           if curr2+count > best {
               best = curr2 + count
               a = a2
               b = i
           }
           curr2 = 0
           a2 = i + 1
       case h == 2:
           curr2++
       }
   }
   // output result, convert to 1-based with original index
   fmt.Fprintln(writer, best)
   // note: a and b are 0-based on s2; shift back
   start := (a + id) % n
   end := (b + id) % n
   fmt.Fprintf(writer, "%d %d\n", start+1, end+1)
}
