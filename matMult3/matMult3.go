package main

import (
	"fmt"
	"math/rand"
	"os"
	"runtime"
	"time"
)

type Matrix [][]float64

func rowByFullMatrix(row *[]float64, cols, resultMatrix *Matrix, rowNum int) {
	
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

}

func main() {
	a, b := makeMatrix(1024, 900), makeMatrix(900, 1500)
	
	rowsa, rowsb := len(a), len(b)
	colsa, colsb := len((a)[0]), len((b)[0])

	if colsa != rowsb {
		fmt.Println("Matrices cannot be multiplied!")
		os.Exit(3)
	}

	runtime.GOMAXPROCS(1)

	//start
	start := time.Now()

	matrixC := make([][]float64, rowsa)
	for i := range matrixC {
		matrixC[i] = make([]float64, colsb)
	}

	resultMatrix := Matrix(matrixC)

	numRowsInResultMatrix := rowsa

	for i := 0; i < numRowsInResultMatrix; i++ {
		
		go rowByFullMatrix(&a[i], &b, &resultMatrix, i)

	}

	//end
	elapsed := time.Since(start)

	time.Sleep(1* time.Second)

	fmt.Print("\nFinished. Elapsed Time: ", elapsed.Microseconds(), " Microseconds")
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
