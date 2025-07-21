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
   // Read n
   line, err := reader.ReadString('\n')
   if err != nil && err.Error() != "EOF" {
       fmt.Fprintln(os.Stderr, "failed to read n:", err)
       return
   }
   line = strings.TrimSpace(line)
   n, err := strconv.Atoi(line)
   if err != nil {
       fmt.Fprintln(os.Stderr, "invalid n:", err)
       return
   }
   // Read numbers line
   numsLine, err := reader.ReadString('\n')
   if err != nil && err.Error() != "EOF" {
       fmt.Fprintln(os.Stderr, "failed to read numbers:", err)
       return
   }
   numsLine = strings.TrimSpace(numsLine)
   parts := strings.Split(numsLine, ",")
   if len(parts) != n {
       // allow if extra spaces lead to empty parts
       // but strictly check
   }
   nums := make([]int, n)
   for i, p := range parts {
       p = strings.TrimSpace(p)
       v, err := strconv.Atoi(p)
       if err != nil {
           fmt.Fprintln(os.Stderr, "invalid number at position", i, err)
           return
       }
       nums[i] = v
   }
   // Check each pair
   for i := 0; i < n; i++ {
       for j := i + 1; j < n; j++ {
           a, b := nums[i], nums[j]
           // check if a divides b or b divides a
           ok := false
           if a != 0 && b%a == 0 {
               ok = true
           } else if b != 0 && a%b == 0 {
               ok = true
           }
           if !ok {
               fmt.Println("NOT FRIENDS")
               return
           }
       }
   }
   fmt.Println("FRIENDS")
}
