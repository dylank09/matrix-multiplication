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
	
	sum := 0.0
	for i := range *currentRow { // multiply row of first matrix by column in second matrix
		sum += (*currentRow)[i] * (*currentCol)[i]
	}
	
	(*resultMatrix)[resultRowIndex][resultColIndex] = sum //fills in the element in the resulting matrix 
	
	wg.Done()
}

func rowByFullMatrixAlgo2(row *[]float64, matB, resultMatrix *Matrix, rowNum int, wg *sync.WaitGroup) {
	
	defer wg.Done()

	numColsInB := len(*matB)

	sum := 0.0
	for i := 0; i < numColsInB; i++ {

		currentColData := (*matB)[i]

		//multiply the row by the current column to get an element in resulting row
		sum = 0.0
		for index, val := range *row {
			sum += val * currentColData[index]
		}
		
		(*resultMatrix)[rowNum][i] = sum //fills in the elements of row in the resulting matrix
	}

}

func colByFullMatrixAlgo3(col *[]float64, matA, resultMatrix *Matrix, colNum int, wg *sync.WaitGroup) {

	numRowsInA := len(*matA)

	// result := make([] float64, numRowsInA) //new col will be same size as row in mat A
	sum := 0.0
	for i := 0; i < numRowsInA; i++ {

		currentRowData := (*matA)[i]

		//multiply the row by the current column to get an element
		sum = 0.0
		for index, c := range *col {
			sum += c * currentRowData[index]
		}

		(*resultMatrix)[i][colNum] = sum //put the result straight into the result matrix column
		
	}

	wg.Done()
}

func rowByColGetMatrixAlgo3b(resultMatrix *Matrix, colA, rowB *[]float64, wg *sync.WaitGroup) {
	defer wg.Done()

	for i := range *rowB {
		for j := range *colA {
			(*resultMatrix)[j][i] += (*rowB)[i] * (*colA)[j] //add to current value
		}
	}
}

func main() {

	rand.Seed(time.Now().UnixNano())

	a, b := makeMatrix(1000, 1024), makeMatrix(1024, 900)

	rowsa, rowsb := len(a), len(b)
	colsa, colsb := len((a)[0]), len((b)[0])

	if colsa != rowsb {
		fmt.Println("Matrices cannot be multiplied!")
		os.Exit(3)
	}

	//start sequential

	fmt.Println("Sequential Algorithm")
	
	start0 := time.Now()

	sequentialMat := makeEmptyMatrix(rowsa, colsb)

	sum := 0.0
	for i := 0; i < rowsa; i++ {
		for j := 0; j < colsb; j++ {
			sum = 0.0
			for k := 0; k < rowsb; k++ {
				sum += a[i][k]*b[k][j]
			}
			sequentialMat[i][j] = sum
		}
	}
	
	elapsed0 := time.Since(start0)
	
	fmt.Println("Finished sequential. Elapsed Time: ", elapsed0)

	//end sequential --------------------------------------------------------

	//start algorithm 1

	fmt.Println("\nFirst Algorithm")
	
	start1 := time.Now()

	transposedB := transposeMatrix(b) //transpose matrix b to make columns into rows

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

	//end algorithm 1 --------------------------------------------------------

	//start algorithm 2
	
	fmt.Println("\n\nSecond Algorithm: ")

	start2 := time.Now()

	transposedB = transposeMatrix(b) //transpose matrix b to make columns into rows

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

	//end algorithm 2 --------------------------------------------------------

	//start algorithm 3

	fmt.Println("\nThird Algorithm: ")

	start3 := time.Now()

	resultMatrix3 := makeEmptyMatrix(rowsa, colsb)

	transposedB = transposeMatrix(b) //transpose matrix b to make columns into rows
	//ran again so that there isn't a bias in the timing
	
	var wg3 sync.WaitGroup

	numColsInResultMatrix := colsb

	wg3.Add(numColsInResultMatrix)
	for i := 0; i < numColsInResultMatrix; i++ {
		
		go colByFullMatrixAlgo3(&transposedB[i], &a, &resultMatrix3, i, &wg3)

	}
	wg3.Wait()
	
	elapsed3 := time.Since(start3)
	fmt.Println("Finished algorithm 3. Elapsed Time: ", elapsed3)

	//end algorithm 3 --------------------------------------------------------

	//start algorithm 3b

	fmt.Println("\nExtra Algorithm: ")

	start3b := time.Now()

	transposedA := transposeMatrix(a) 

	var wg3b sync.WaitGroup

	resultMatrix3b := makeEmptyMatrix(rowsa, colsb)

	wg3b.Add(colsa)

	for i := 0; i < colsa; i++ {
		
		colA := transposedA[i]
		rowB := b[i]
		go rowByColGetMatrixAlgo3b(&resultMatrix3b, &colA, &rowB, &wg3b)

	}

	wg3b.Wait()

	elapsed3b := time.Since(start3b)
	fmt.Println("Finished algorithm 3. Elapsed Time: ", elapsed3b)

	//end algorithm 3b --------------------------------------------------------

	//test equality

	fmt.Println("Check result 1: ", compareMatrices(&sequentialMat, &resultMatrix1))
	fmt.Println("Check result 2: ", compareMatrices(&sequentialMat, &resultMatrix2))
	fmt.Println("Check result 3: ", compareMatrices(&sequentialMat, &resultMatrix3))
	fmt.Println("Check result 3b: ", compareMatrices(&sequentialMat, &resultMatrix3b))
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

//make a matix with specific dimensions filled with random numbers
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

//compare two matrices element by element and if any one is different, return false
func compareMatrices(a, b *Matrix) bool {
	//this function assumes the matrices are the same dimensions
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

//turn a matrix on it's side --> matrix columns become it's rows
func transposeMatrix(m Matrix) Matrix {
	transposedM := makeEmptyMatrix(len(m[0]), len(m))
	for i, rows := range m {
		for j:= range rows {
			transposedM[j][i] = m[i][j]
		}
	}
	return transposedM
}

//fill a matrix with zeros instead of random numbers
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