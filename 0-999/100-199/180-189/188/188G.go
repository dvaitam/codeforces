package main

import (
   "bufio"
   "fmt"
   "os"
   "sort"
)

func main() {
   reader := bufio.NewReader(os.Stdin)
   writer := bufio.NewWriter(os.Stdout)
   defer writer.Flush()

   var n int
   if _, err := fmt.Fscan(reader, &n); err != nil {
       return
   }
   nums := make([]int, n)
   for i := 0; i < n; i++ {
       fmt.Fscan(reader, &nums[i])
   }
   sort.Ints(nums)
   for i, v := range nums {
       if i > 0 {
           writer.WriteByte(' ')
       }
       writer.WriteString(fmt.Sprint(v))
   }
}
