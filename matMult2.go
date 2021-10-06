package matMult2

import (
	"fmt"
)

type Matrix [][]float64

/* 
* second concurrent algorithm will be the following
* have Go routines that go off and do a full row x column
* to get the corresponding value of the resulting matrix i.e.
* 	Go routine: A[0][0] * B[0][0] + A[0][1] * B[1][0] = C[0][0]
* 	Go routine: A[0][0] * B[0][1] + A[0][1] * B[1][1] = C[0][1]
*	so on and once an element is calculated, add them together
*/
func secondAlgorithm(a, b *Matrix) *Matrix {
	rowsa, rowsb := len(*a), len(*b)
	_, colsb := len((*a)[0]), len((*b)[0])
	n := numElementsInMatrix(rowsa, colsb)

	c := Matrix{{}}

	for i := 0; i < n; i++ {
		
		go func () {
			sum := 0
			for k := 0; k < rowsb; k++ {
				sum += *a[][k] * *b[k][]
			}
		}()

	}
	
	return &c 
}

func main() {
	a := Matrix{{1, 2}, {2, 3}, {3, 4}}
    b := Matrix{{2, 3,4}, {5, 6,7}}

	fmt.Println("Matrix A")
	printMatrix(&a)
	
	fmt.Println("Matrix B")
	printMatrix(&b)


	result := secondAlgorithm(&a, &b)
	//expected answer is:  Matrix{{12,15,18}, {19,24,29}, {26,33,40}}

	fmt.Println("Result is: ")
	printMatrix(result)
	fmt.Print("Finished")
}

//Helper functions
func printMatrix(m* Matrix) {
	for i := 0; i < len(*m); i++ {
        for j := 0; j < len((*m)[0]); j++ {
            fmt.Print((*m)[i][j])
			fmt.Print(" ")
        }
        fmt.Print("\n")
    } 
	fmt.Print("\n")
}

func numElementsInMatrix(rows, cols int) int {
	numElements := rows*cols
	// fmt.Printf("Matrix has %d elements\n", numElements)
	return numElements
}