package main

import (
	"fmt"
	"math/rand"
	"os"
	"runtime"
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
	
	matrixA, matrixB := makeMatrix(1024, 900), makeMatrix(900, 1500)

	rowsa, rowsb := len(matrixA), len(matrixB)
	colsa, colsb := len((matrixA)[0]), len((matrixB)[0])

	if colsa != rowsb {
		fmt.Println("Matrices cannot be multiplied!")
		os.Exit(3)
	}

	runtime.GOMAXPROCS(4)

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

	fmt.Print("\nFinished. Elapsed Time: ", elapsed.Microseconds(), " Microseconds")
	
	time.Sleep(4* time.Second)
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

func makeMatrix(rows, cols int) (m Matrix) {
	m = make([][]float64, rows)
	for i := range m {
		m[i] = make([]float64, cols)
		for j := range m[i] {
			m[i][j] = float64(rand.Intn(100))
		}
	}
	return
}

/* 
* first concurrent algorithm will be the following
* have Go routines each multiplying elements of the matrix
* before adding them together. i.e.
* 	Go routine: A[0][0] * B[0][0]
* 	Go routine: A[0][1] * B[1][0]
*	so on and once a row is done, add them together
*/