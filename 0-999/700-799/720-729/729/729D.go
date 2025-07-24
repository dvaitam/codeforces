package main

import (
   "bufio"
   "fmt"
   "os"
)

func main() {
   reader := bufio.NewReader(os.Stdin)
   var n, a, b, k int
   if _, err := fmt.Fscan(reader, &n, &a, &b, &k); err != nil {
       return
   }
   var s string
   if _, err := fmt.Fscan(reader, &s); err != nil {
       return
   }
   var candidates []int
   // scan zero segments
   for i := 0; i < n; {
       if s[i] == '1' {
           i++
           continue
       }
       j := i
       for j < n && s[j] == '0' {
           j++
       }
       L := j - i
       cnt := L / b
       for t := 1; t <= cnt; t++ {
           // position to shoot: i + t*b (1-based index)
           pos := i + t*b
           candidates = append(candidates, pos)
       }
       i = j
   }
   S := len(candidates)
   // need to reduce S to a-1, so need R shots
   R := S - a + 1
   if R < 0 {
       R = 0
   }
   // output
   writer := bufio.NewWriter(os.Stdout)
   defer writer.Flush()
   fmt.Fprintln(writer, R)
   for i := 0; i < R; i++ {
       if i > 0 {
           writer.WriteByte(' ')
       }
       fmt.Fprint(writer, candidates[i])
   }
   // second line end
   writer.WriteByte('\n')
}
