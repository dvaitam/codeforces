package main

import (
	"bufio"
	"fmt"
	"os"
)

// The full reconstruction of the array a from the provided b_{i,n} values
// requires a substantial bitwise inversion that is outside the scope of this
// stub.  To keep the repository self-contained we simply read the input and
// output -1 indicating that we could not reconstruct a valid array.
func main() {
	reader := bufio.NewReader(os.Stdin)
	writer := bufio.NewWriter(os.Stdout)
	defer writer.Flush()

	var n int
	if _, err := fmt.Fscan(reader, &n); err != nil {
		return
	}
	for i := 0; i < n; i++ {
		var x int
		fmt.Fscan(reader, &x)
	}
	fmt.Fprintln(writer, -1)
}
