package main

import (
   "bufio"
   "fmt"
   "os"
   "strconv"
)

func main() {
   reader := bufio.NewReader(os.Stdin)
   var n int
   if _, err := fmt.Fscan(reader, &n); err != nil {
       return
   }
   alcohols := map[string]struct{}{
       "ABSINTH": {}, "BEER": {}, "BRANDY": {}, "CHAMPAGNE": {},
       "GIN": {}, "RUM": {}, "SAKE": {}, "TEQUILA": {},
       "VODKA": {}, "WHISKEY": {}, "WINE": {},
   }
   count := 0
   for i := 0; i < n; i++ {
       var s string
       fmt.Fscan(reader, &s)
       // try parse as integer age
       if age, err := strconv.Atoi(s); err == nil {
           if age < 18 {
               count++
           }
       } else {
           // treat as drink
           if _, ok := alcohols[s]; ok {
               count++
           }
       }
   }
   fmt.Println(count)
}
