package main

import (
   "bufio"
   "fmt"
   "os"
)

// check returns true if sum of digits in s[lastIndex+1..curIndex] has a prefix with sum % 3 == 0
func check(s string, curIndex, lastIndex int) bool {
   sum := 0
   for i := curIndex; i > lastIndex; i-- {
       sum += int(s[i] - '0')
       if sum%3 == 0 {
           return true
       }
   }
   return false
}

func main() {
   reader := bufio.NewReader(os.Stdin)
   var s string
   if _, err := fmt.Fscan(reader, &s); err != nil {
       return
   }

   lastIndex := -1
   ans := 0
   for i := 0; i < len(s); i++ {
       if check(s, i, lastIndex) {
           ans++
           lastIndex = i
       }
   }
   fmt.Println(ans)
}
