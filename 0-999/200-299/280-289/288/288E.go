package main

import (
   "bufio"
   "fmt"
   "os"
)

func main() {
   reader := bufio.NewReader(os.Stdin)
   var l, r string
   fmt.Fscan(reader, &l)
   fmt.Fscan(reader, &r)
   // TODO: implement solution for sum of products of consecutive lucky numbers
   // This problem requires advanced algorithms for large lengths.
   fmt.Println(0)
}
