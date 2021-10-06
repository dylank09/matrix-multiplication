package matMult3

import (
	"fmt"
)

type Matrix [][]float64

/*
*
*
*
*
*
*
*
*/
func thirdAlgorithm(a, b *Matrix) *Matrix {
	rowsa, _ := len(*a), len(*b)
	_, colsb := len((*a)[0]), len((*b)[0])
	n := numElementsInMatrix(rowsa, colsb)
	
	c := Matrix{{}}
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