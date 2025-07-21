package main

import (
   "bufio"
   "fmt"
   "os"
   "strconv"
   "strings"
)

func main() {
   reader := bufio.NewReader(os.Stdin)
   var n int
   if _, err := fmt.Fscan(reader, &n); err != nil {
       return
   }
   sN := strconv.Itoa(n)
   if strings.Contains(sN, "13") {
       fmt.Println(-1)
       return
   }
   count := 0
   for b := 2; b <= n; b++ {
       x := n
       // collect digits in base b
       var digits []int
       for x > 0 {
           digits = append(digits, x%b)
           x /= b
       }
       // build representation string
       var sb strings.Builder
       // estimate: average digit length ~2
       sb.Grow(len(digits) * 2)
       for i := len(digits) - 1; i >= 0; i-- {
           sb.WriteString(strconv.Itoa(digits[i]))
       }
       if strings.Contains(sb.String(), "13") {
           count++
       }
   }
   fmt.Println(count)
}
