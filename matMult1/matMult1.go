package main

import (
	"fmt"
	"os"
	"time"
)

type Matrix [][]float64

func rowByCol(currentRow, currentCol *[]float64, resultMatrix *Matrix, resultRowIndex, resultColIndex int) {
	sum := 0.0
	for i := range *currentRow {
		sum += (*currentRow)[i] * (*currentCol)[i]
	}

	(*resultMatrix)[resultRowIndex][resultColIndex] = sum //fills in the element in the resulting matrix by multiplying row of first matrix by column in second matrix
}

func main() {
	matrixA := Matrix{{7, 8, 2}, {1, 9, 21}, {34, 14, 8}, {1, 4, 11}, {21, 4, 2}}
    matrixB := Matrix{{2, 11, 17, 21}, {3, 6, 8, 91}, {3, 4, 5, 2}}

	fmt.Println("Matrix A")
	printMatrix(&matrixA)
	
	fmt.Println("Matrix B")
	printMatrix(&matrixB)

	rowsa, rowsb := len(matrixA), len(matrixB)
	colsa, colsb := len((matrixA)[0]), len((matrixB)[0])

	if colsa != rowsb {
		fmt.Println("Matrices cannot be multiplied!")
		os.Exit(3)
	}

	//start
	start := time.Now()

	matrixC := make([][]float64, rowsa)
	for i := range matrixC {
		matrixC[i] = make([]float64, colsb)
	}

	resultMatrix := Matrix(matrixC)

	numElementsInResultMatrix := rowsa * colsb
	currentRowIndex := 0
	currentCol := 0
	for i := 0; i < numElementsInResultMatrix; i++ {
		if i != 0 && i % colsb == 0 { //*
			currentRowIndex += 1
		} 

		currentColData := make([]float64, rowsb)
		if currentCol >= colsb {
			currentCol = 0
		}
		
		for i := 0; i < rowsb; i++ {
			currentColData[i] = matrixB[i][currentCol]
		}

		go rowByCol(&matrixA[currentRowIndex], &currentColData, &resultMatrix, currentRowIndex, currentCol)

		currentCol += 1
	}

	//end
	elapsed := time.Since(start)

	fmt.Println("Multiply Matrices A and B to get Matrix C:")
	printMatrix(&resultMatrix)
	
	//expected answer is: Matrix{{44, 	133, 	193, 	879}, 
	//							 {92, 	149, 	194, 	882}, 
	//							 {134,  490, 	730, 	2004}, 
	//							 {47, 	79, 	104, 	407}, 
	//							 {60, 	263, 	399, 	809}}

	fmt.Print("\nFinished. Elapsed Time: ", elapsed)
	
	time.Sleep(1* time.Second)
}

//Helper functions
func printMatrix(m* Matrix) {
	for i := 0; i < len(*m); i++ {
        for j := 0; j < len((*m)[0]); j++ {
            fmt.Print((*m)[i][j])
			fmt.Print("\t")
        }
        fmt.Print("\n")
    } 
	fmt.Print("\n")
}

/* 
* first concurrent algorithm will be the following
* have Go routines each multiplying elements of the matrix
* before adding them together. i.e.
* 	Go routine: A[0][0] * B[0][0]
* 	Go routine: A[0][1] * B[1][0]
*	so on and once a row is done, add them together
*/