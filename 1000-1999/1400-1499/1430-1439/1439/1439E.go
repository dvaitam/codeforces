package main

import (
   "bufio"
   "fmt"
   "os"
   "strconv"
   "strings"
   "math/bits"
)

func main() {
   reader := bufio.NewReader(os.Stdin)
   line, err := reader.ReadString('\n')
   if err != nil {
       panic(err)
   }
   line = strings.TrimSpace(line)
   m, err := strconv.Atoi(line)
   if err != nil {
       panic(err)
   }
   var T uint64 = 0
   for i := 0; i < m; i++ {
       s, err := reader.ReadString('\n')
       if err != nil {
           panic(err)
       }
       s = strings.TrimSpace(s)
       parts := strings.Fields(s)
       if len(parts) < 4 {
           i--
           continue
       }
       x1, _ := strconv.ParseUint(parts[0], 10, 64)
       y1, _ := strconv.ParseUint(parts[1], 10, 64)
       x2, _ := strconv.ParseUint(parts[2], 10, 64)
       y2, _ := strconv.ParseUint(parts[3], 10, 64)
       d1 := bits.OnesCount64(x1) + bits.OnesCount64(y1)
       d2 := bits.OnesCount64(x2) + bits.OnesCount64(y2)
       T ^= f(uint64(d1)) ^ f(uint64(d2))
   }
   // Determine minimum chops to make T zero
   if T == 0 {
       fmt.Println(0)
   } else if (T % 4 == 1 && T != 1) || (T % 4 == 2) {
       fmt.Println(2)
   } else {
       fmt.Println(1)
   }
}

// f returns XOR of all integers from 0 to d
func f(d uint64) uint64 {
   switch d & 3 {
   case 0:
       return d
   case 1:
       return 1
   case 2:
       return d + 1
   default:
       return 0
   }
}
