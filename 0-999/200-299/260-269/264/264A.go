package main

import (
   "bufio"
   "fmt"
   "io"
   "os"
   "strings"
)

func main() {
   reader := bufio.NewReader(os.Stdin)
   // Read entire input as string
   input, err := io.ReadAll(reader)
   if err != nil {
       panic(err)
   }
   s := strings.TrimSpace(string(input))
   n := len(s)
   // Use a deque with enough space
   // ans[left:right] will contain result
   size := 2*n + 1
   ans := make([]int, size)
   left, right := n, n
   // Build deque by reverse iteration
   for i := n - 1; i >= 0; i-- {
       num := i + 1
       if s[i] == 'r' {
           left--
           ans[left] = num
       } else {
           ans[right] = num
           right++
       }
   }
   // Output result
   writer := bufio.NewWriter(os.Stdout)
   defer writer.Flush()
   for i := left; i < right; i++ {
       fmt.Fprintln(writer, ans[i])
   }
}
