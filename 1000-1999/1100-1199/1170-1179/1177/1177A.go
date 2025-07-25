package main

import (
   "bufio"
   "fmt"
   "os"
)

func main() {
   reader := bufio.NewReader(os.Stdin)
   var k int
   if _, err := fmt.Fscan(reader, &k); err != nil {
       return
   }
   length := 1
   count := 9
   start := 1
   // find the digit length where k falls
   for k > length*count {
       k -= length * count
       length++;
       count *= 10
       start *= 10
   }
   // determine the exact number
   offset := (k - 1) / length
   num := start + offset
   s := fmt.Sprint(num)
   // position in the number
   pos := (k - 1) % length
   fmt.Printf("%c", s[pos])
}
