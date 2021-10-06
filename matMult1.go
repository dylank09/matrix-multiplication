package main

import (
	"fmt"
	"sync"
)

type Matrix [][]float64

func algorithm(row, col []float64, wg *sync.WaitGroup, elements chan<- float64) {
	sum := 0.0
	for i := range row {
		sum += row[i] * col[i]
	}
	elements <- sum
	wg.Done()
}

func main() {
	a := Matrix{{1, 2}, {2, 3}, {3, 4}}
    b := Matrix{{2, 3,4}, {5, 6,7}}

	fmt.Println("Matrix A")
	printMatrix(&a)
	
	fmt.Println("Matrix B")
	printMatrix(&b)

	//start

	rowsa, rowsb := len(a), len(b)
	_, colsb := len((a)[0]), len((b)[0])

	var wg sync.WaitGroup

	numElementsInResultMatrix := rowsa * colsb
	elements := make(chan float64, numElementsInResultMatrix)

	currentRow := 0
	currentCol := 0
	for i := 0; i < numElementsInResultMatrix; i++ {
		if i != 0 && i % colsb == 0 { //*
			currentRow += 1
		} 

		currentColData := make([]float64, rowsb)
		currentColData = nil
		if currentCol >= colsb {
			currentCol = 0
		}
		
		for i := 0; i < rowsb; i++ {
			currentColData = append(currentColData, b[i][currentCol])
		}
		wg.Add(1)
		go algorithm(a[currentRow], currentColData, &wg, elements)
		wg.Wait()

		currentCol += 1

	}

	fmt.Println("Multiply Matrices A and B to get Matrix C:")

	for i := 0; i < numElementsInResultMatrix; i++ {
		if i % colsb == 0 {
			fmt.Println()
		}
		fmt.Print(<-elements, " ")
	}
	
	//expected answer is:  Matrix{{12,15,18}, {19,24,29}, {26,33,40}}

	fmt.Print("\nFinished")
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


	
	// fmt.Println("Matrix A: ", rowsa, " x ", colsa)
	// fmt.Println("Matrix B: ", rowsb, " x ", colsb)
	// fmt.Println("Therefore Matrix C: ", rowsa, " x ", colsb)

	// c := Matrix{{}}
	// return &c 


/* 
* first concurrent algorithm will be the following
* have Go routines each multiplying elements of the matrix
* before adding them together. i.e.
* 	Go routine: A[0][0] * B[0][0]
* 	Go routine: A[0][1] * B[1][0]
*	so on and once a row is done, add them together
*/