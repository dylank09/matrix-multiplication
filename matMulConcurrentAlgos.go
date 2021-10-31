package main

import (
	"fmt"
	"math/rand"
	"os"
	"sync"
	"time"
)

type Matrix [][]float64

func rowByColAlgo1(currentRow, currentCol *[]float64, resultMatrix *Matrix, resultRowIndex, resultColIndex int, wg *sync.WaitGroup) {
	defer wg.Done()

	sum := 0.0
	for i := range *currentRow {
		sum += (*currentRow)[i] * (*currentCol)[i]
	}

	(*resultMatrix)[resultRowIndex][resultColIndex] = sum //fills in the element in the resulting matrix by multiplying row of first matrix by column in second matrix
}

func rowByFullMatrixAlgo2(row *[]float64, cols, resultMatrix *Matrix, rowNum int, wg *sync.WaitGroup) {
	
	defer wg.Done()

	result := make([] float64, len(*cols)) //new row will be same size as row in mat A
	sum := 0.0
	for i := 0; i < len(*cols); i++ {

		currentColData := (*cols)[i]

		//multiply the row by the current column to get an element
		sum = 0.0
		for k, r := range *row {
			sum += r * currentColData[k]
		}
		result[i] = sum
		
	}
	(*resultMatrix)[rowNum] = result //returns full returning row on the new matrix

}

func rowByColGetMatrixAlgo3(resultMatrix *Matrix, colA, rowB *[]float64, wg *sync.WaitGroup) {
	defer wg.Done()

	for i := range *rowB {
		for j := range *colA {
			
			(*resultMatrix)[j][i] += (*rowB)[i] * (*colA)[j]
			
		}
	}
}

func main() {
	// a, b := makeMatrix(1024, 900), makeMatrix(900, 1400)
	// a, b := makeMatrix(10, 5), makeMatrix(5, 11)
	a := Matrix{{2,7}, {3,8}, {4,9}}
	b := Matrix{{7,9,2,3}, {8,6,1,4}}

	rowsa, rowsb := len(a), len(b)
	colsa, colsb := len((a)[0]), len((b)[0])

	if colsa != rowsb {
		fmt.Println("Matrices cannot be multiplied!")
		os.Exit(3)
	}

	// runtime.GOMAXPROCS(4)

	//start sequential
	fmt.Println("Sequential Algorithm")
	
	start0 := time.Now()

	resultMatrix0 := makeEmptyMatrix(rowsa, colsb)

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

	transposedB := transposeMat(b) //transpose matrix b to make columns into rows

	resultMatrix1 := makeEmptyMatrix(rowsa, colsb)

	numElementsInResultMatrix := rowsa * colsb
	currentRowIndex := 0
	currentCol := 0

	var wg1 sync.WaitGroup

	wg1.Add(numElementsInResultMatrix)

	for i := 0; i < numElementsInResultMatrix; i++ {
		if i != 0 && i % colsb == 0 { //*
			currentRowIndex += 1
		} 

		if currentCol >= colsb {
			currentCol = 0
		}

		currentColData := transposedB[currentCol]

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

	transposedB = transposeMat(b) //transpose matrix b to make columns into rows

	resultMatrix2 := makeEmptyMatrix(rowsa, colsb)

	numRowsInResultMatrix := rowsa
	
	var wg2 sync.WaitGroup

	wg2.Add(numRowsInResultMatrix)
	for i := 0; i < numRowsInResultMatrix; i++ {
		
		go rowByFullMatrixAlgo2(&a[i], &transposedB, &resultMatrix2, i, &wg2)

	}
	wg2.Wait()
	elapsed2 := time.Since(start2)

	fmt.Println("Finished algorithm 2. Elapsed Time: ", elapsed2)

	//end algorithm 2

	//start algorithm 3

	fmt.Println("\nThird Algorithm: ")

	start3 := time.Now()

	transposedA := transposeMat(a)

	var wg3 sync.WaitGroup

	resultMatrix3 := makeEmptyMatrix(rowsa, colsb)

	wg3.Add(colsa)

	for i := 0; i < colsa; i++ {
		
		colA := transposedA[i]
		rowB := b[i]
		go rowByColGetMatrixAlgo3(&resultMatrix3, &colA, &rowB, &wg3)

	}
	wg3.Wait()

	// time.Sleep(time.Second * 2)

	//wait?

	elapsed3 := time.Since(start3)
	fmt.Println("Finished algorithm 3. Elapsed Time: ", elapsed3)

	//end algorithm 3

	//test equality

	fmt.Println("Check result 1: ", compareMatrices(&resultMatrix0, &resultMatrix1))
	fmt.Println("Check result 2: ", compareMatrices(&resultMatrix0, &resultMatrix2))
	fmt.Println("Check result 3: ", compareMatrices(&resultMatrix0, &resultMatrix3))

	if !compareMatrices(&resultMatrix0, &resultMatrix3) {
		printMatrix(&resultMatrix3)
	}
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

func compareMatrices(a, b *Matrix) bool {
	//assumes the matrices are the same dimensions
	result := true
	for i := 0; i < len(*a); i++ {
		for j := 0; j < len((*a)[i]); j++ {
			if (*a)[i][j] != (*b)[i][j] {
				result = false
			}
		}
	}
	return result
}

func transposeMat(m Matrix) Matrix {
	transposedM := makeEmptyMatrix(len(m[0]), len(m))
	for i, rows := range m {
		for j:= range rows {
			transposedM[j][i] = m[i][j]
		}
	}
	return transposedM
}

func makeEmptyMatrix(rows, cols int) Matrix {
	m := make([][]float64, rows)
	for i := range m {
		m[i] = make([]float64, cols)
		for j := range m[i] {
			m[i][j] = 0
		}
	}
	return m
}