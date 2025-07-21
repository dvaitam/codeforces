package main

import (
   "bufio"
   "fmt"
   "os"
)

func main() {
   reader := bufio.NewReader(os.Stdin)
   writer := bufio.NewWriter(os.Stdout)
   defer writer.Flush()

   var n, k int
   if _, err := fmt.Fscan(reader, &n, &k); err != nil {
       return
   }
   model := make([][]byte, n)
   for i := 0; i < n; i++ {
       var line string
       fmt.Fscan(reader, &line)
       model[i] = []byte(line)
   }

   // Initialize current fractal as the model (step 1)
   cur := make([][]byte, n)
   for i := 0; i < n; i++ {
       cur[i] = make([]byte, n)
       copy(cur[i], model[i])
   }
   curSize := n

   // Perform further steps 2..k
   for step := 2; step <= k; step++ {
       nextSize := curSize * n
       next := make([][]byte, nextSize)
       for i := 0; i < nextSize; i++ {
           next[i] = make([]byte, nextSize)
           for j := 0; j < nextSize; j++ {
               next[i][j] = '.'
           }
       }
       for i := 0; i < curSize; i++ {
           for j := 0; j < curSize; j++ {
               if cur[i][j] == '*' {
                   // fill entire block with black
                   for di := 0; di < n; di++ {
                       for dj := 0; dj < n; dj++ {
                           next[i*n+di][j*n+dj] = '*'
                       }
                   }
               } else {
                   // apply model to this block
                   for di := 0; di < n; di++ {
                       for dj := 0; dj < n; dj++ {
                           next[i*n+di][j*n+dj] = model[di][dj]
                       }
                   }
               }
           }
       }
       cur = next
       curSize = nextSize
   }

   // Output result
   for i := 0; i < curSize; i++ {
       writer.Write(cur[i])
       writer.WriteByte('\n')
   }
}
