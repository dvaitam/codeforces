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

   var s string
   if _, err := fmt.Fscan(reader, &s); err != nil {
       return
   }
   n := len(s)
   maxSum := -1

   for i := 1; i <= n-2; i++ {
       for j := i + 1; j <= n-1; j++ {
           a := s[:i]
           b := s[i:j]
           c := s[j:]
           va, ok := parseSegment(a)
           if !ok {
               continue
           }
           vb, ok := parseSegment(b)
           if !ok {
               continue
           }
           vc, ok := parseSegment(c)
           if !ok {
               continue
           }
           sum := va + vb + vc
           if sum > maxSum {
               maxSum = sum
           }
       }
   }

   fmt.Fprint(writer, maxSum)
}

// parseSegment checks if s represents a valid round score:
// no leading zeros (except "0"), value ≤ 1000000.
func parseSegment(s string) (int, bool) {
   if len(s) == 0 {
       return 0, false
   }
   if len(s) > 1 && s[0] == '0' {
       return 0, false
   }
   if len(s) > 7 {
       return 0, false
   }
   if len(s) == 7 && s > "1000000" {
       return 0, false
   }
   v, err := strconv.Atoi(s)
   if err != nil {
       return 0, false
   }
   return v, true
}
