package main

import (
	"fmt"
	"math/rand"
	"os"
	"runtime"
	"sync"
	"time"
)

type Matrix [][]float64

func rowByColAlgo1(currentRow, currentCol *[]float64, resultMatrix *Matrix, resultRowIndex, resultColIndex int, wg *sync.WaitGroup) {
	sum := 0.0
	for i := range *currentRow {
		sum += (*currentRow)[i] * (*currentCol)[i]
	}

	(*resultMatrix)[resultRowIndex][resultColIndex] = sum //fills in the element in the resulting matrix by multiplying row of first matrix by column in second matrix
	wg.Done()
}

func rowByFullMatrixAlgo2(row *[]float64, cols, resultMatrix *Matrix, rowNum int, wg *sync.WaitGroup) {
	
	result := make([] float64, len((*cols)[0])) //new row will be same size as row in mat A
	sum := 0.0
	for i := 0; i < len((*cols)[0]); i++ {

		currentColData := make([]float64, len(*cols))
		for j := 0; j < len(*cols); j++ {
			currentColData[j] = (*cols)[j][i]
		}

		//multiply the row by the current column to get an element
		sum = 0.0
		for k, r := range *row {
			sum += r * currentColData[k]
		}
		result[i] = sum
		
	}
	(*resultMatrix)[rowNum] = result //returns full returning row on the new matrix
	wg.Done()

}

func main() {
	a, b := makeMatrix(1024, 900), makeMatrix(900, 1500)

	rowsa, rowsb := len(a), len(b)
	colsa, colsb := len((a)[0]), len((b)[0])

	if colsa != rowsb {
		fmt.Println("Matrices cannot be multiplied!")
		os.Exit(3)
	}

	runtime.GOMAXPROCS(4)

	//start sequential
	fmt.Println("Sequential Algorithm")
	
	start0 := time.Now()

	result0 := make([][]float64, rowsa)
	for i := range result0 {
		result0[i] = make([]float64, colsb)
	}

	resultMatrix0 := Matrix(result0)
	sum := 0.0
	for i := 0; i < rowsa; i++ {
		for j := 0; j < colsb; j++ {
			sum = 0.0
			for k := 0; k < rowsb; k++ {
				sum += a[i][k]*b[k][j]
			}
			resultMatrix0[i][j] = sum
		}
	}
	
	elapsed0 := time.Since(start0)
	
	fmt.Println("Finished sequential. Elapsed Time: ", elapsed0)

	//end sequential

	//start algorithm 1
	fmt.Println("\nFirst Algorithm")
	
	start1 := time.Now()

	result1 := make([][]float64, rowsa)
	for i := range result1 {
		result1[i] = make([]float64, colsb)
	}

	resultMatrix1 := Matrix(result1)

	numElementsInResultMatrix := rowsa * colsb
	currentRowIndex := 0
	currentCol := 0

	var wg1 sync.WaitGroup

	wg1.Add(numElementsInResultMatrix)

	for i := 0; i < numElementsInResultMatrix; i++ {
		if i != 0 && i % colsb == 0 { //*
			currentRowIndex += 1
		} 

		currentColData := make([]float64, rowsb)
		if currentCol >= colsb {
			currentCol = 0
		}
		
		for i := 0; i < rowsb; i++ {
			currentColData[i] = b[i][currentCol]
		}

		go rowByColAlgo1(&a[currentRowIndex], &currentColData, &resultMatrix1, currentRowIndex, currentCol, &wg1)

		currentCol += 1

	}

	wg1.Wait()

	elapsed1 := time.Since(start1)
	
	fmt.Println("Finished algorithm 1. Elapsed Time: ", elapsed1)

	//end algorithm 1

	//start algorithm 2
	
	fmt.Println("\n\nSecond Algorithm: ")

	start2 := time.Now()

	result2 := make([][]float64, rowsa)
	for i := range result2 {
		result2[i] = make([]float64, colsb)
	}

	resultMatrix2 := Matrix(result2)

	numRowsInResultMatrix := rowsa
	
	var wg2 sync.WaitGroup

	wg2.Add(numRowsInResultMatrix)
	for i := 0; i < numRowsInResultMatrix; i++ {
		
		go rowByFullMatrixAlgo2(&a[i], &b, &resultMatrix2, i, &wg2)

	}
	wg2.Wait()
	elapsed2 := time.Since(start2)

	fmt.Println("Finished algorithm 2. Elapsed Time: ", elapsed2)

	//end algorithm 2

	//start algorithm 3
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
