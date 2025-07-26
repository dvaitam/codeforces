package main

import (
   "bufio"
   "fmt"
   "os"
)

func main() {
   reader := bufio.NewReader(os.Stdin)
   // fast integer input
   readInt := func() (int, error) {
       var x int
       var sign int = 1
       // read next non-space
       for {
           b, err := reader.ReadByte()
           if err != nil {
               return 0, err
           }
           if b == '-' {
               sign = -1
               break
           }
           if b >= '0' && b <= '9' {
               x = int(b - '0')
               break
           }
       }
       // read rest digits
       for {
           b, err := reader.ReadByte()
           if err != nil {
               break
           }
           if b < '0' || b > '9' {
               break
           }
           x = x*10 + int(b-'0')
       }
       return sign * x, nil
   }

   n, err := readInt()
   if err != nil {
       return
   }
   m, _ := readInt()

   // row minima stats
   rowMinVal := make([]int, n)
   rowMinCnt := make([]int, n)
   rowMinPos := make([]int, n)
   // column maxima stats
   colMaxVal := make([]int, m)
   colMaxCnt := make([]int, m)
   colMaxPos := make([]int, m)

   // initialize colMaxVal to minimal
   // values are >=1, so zero is fine

   for i := 0; i < n; i++ {
       // track row min
       const INF = int(1e9 + 5)
       minVal := INF
       cnt := 0
       pos := 0
       for j := 0; j < m; j++ {
           a, _ := readInt()
           // row min
           if a < minVal {
               minVal = a
               cnt = 1
               pos = j
           } else if a == minVal {
               cnt++
           }
           // col max
           if a > colMaxVal[j] {
               colMaxVal[j] = a
               colMaxCnt[j] = 1
               colMaxPos[j] = i
           } else if a == colMaxVal[j] {
               colMaxCnt[j]++
           }
       }
       rowMinVal[i] = minVal
       rowMinCnt[i] = cnt
       rowMinPos[i] = pos
   }

   // find lexicographically smallest equilibrium
   for i := 0; i < n; i++ {
       if rowMinCnt[i] != 1 {
           continue
       }
       j := rowMinPos[i]
       if colMaxCnt[j] == 1 && colMaxPos[j] == i {
           // output 1-based
           fmt.Printf("%d %d", i+1, j+1)
           return
       }
   }
   fmt.Print("0 0")
}
