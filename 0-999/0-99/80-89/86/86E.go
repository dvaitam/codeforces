package main

import (
   "bufio"
   "fmt"
   "os"
)

var answer = [][]int{
   {}, // 0
   {}, // 1
   {2, 1}, // 2
   {3, 2}, // 3
   {4, 3},
   {5, 3},
   {6, 5},
   {7, 6},
   {8, 6, 5, 4},
   {9, 5},
   {10, 7},
   {11, 9},
   {12, 11, 10, 4},
   {13, 12, 11, 8},
   {14, 13, 12, 2},
   {15, 14},
   {16, 14, 13, 11},
   {17, 14},
   {18, 11},
   {19, 6, 2, 1},
   {20, 17},
   {21, 19},
   {22, 21},
   {23, 18},
   {24, 23, 22, 17},
   {25, 22},
   {26, 6, 2, 1},
   {27, 5, 2, 1},
   {28, 25},
   {29, 27},
   {30, 6, 4, 1},
   {31, 28},
   {32, 22, 2, 1},
   {33, 20},
   {34, 27, 2, 1},
   {35, 33},
   {36, 25},
   {37, 5, 4, 3, 2, 1},
   {38, 6, 5, 1},
   {39, 35},
   {40, 38, 21, 19},
   {41, 38},
   {42, 41, 20, 19},
   {43, 42, 38, 37},
   {44, 43, 18, 17},
   {45, 44, 42, 41},
   {46, 45, 26, 25},
   {47, 42},
   {48, 47, 21, 20},
   {49, 40},
   {50, 49, 24, 23},
   {}, // 51
}

func main() {
   reader := bufio.NewReader(os.Stdin)
   writer := bufio.NewWriter(os.Stdout)
   defer writer.Flush()

   for {
       var n int
       if _, err := fmt.Fscan(reader, &n); err != nil {
           break
       }
       arr := make([]int, n)
       for _, x := range answer[n] {
           if x > 0 && x <= n {
               arr[x-1] = 1
           }
       }
       for i, v := range arr {
           if i > 0 {
               writer.WriteByte(' ')
           }
           writer.WriteString(fmt.Sprint(v))
       }
       writer.WriteByte('\n')
       for i, v := range arr {
           if i > 0 {
               writer.WriteByte(' ')
           }
           writer.WriteString(fmt.Sprint(v))
       }
       writer.WriteByte('\n')
   }
}
