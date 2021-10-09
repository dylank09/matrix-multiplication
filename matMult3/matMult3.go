package main

import (
	"fmt"
	"os"
	"sync"
	"time"
)

type Matrix [][]float64

func algorithm(row, col *[]float64, wg *sync.WaitGroup, elements chan<- float64) {
	sum := 0.0
	for i := range *row {
		sum += (*row)[i] * (*col)[i]
	}
	elements <- sum
	wg.Done()
}

func main() {
	a := Matrix{{7, 8, 2}, {1, 9, 21}, {34, 14, 8}, {1, 4, 11}, {21, 4, 2}}
    b := Matrix{{2, 11, 17, 21}, {3, 6, 8, 91}, {3, 4, 5, 2}}

	fmt.Println("Matrix A")
	printMatrix(&a)
	
	fmt.Println("Matrix B")
	printMatrix(&b)

	//start
	start := time.Now()

	rowsa, rowsb := len(a), len(b)
	colsa, colsb := len((a)[0]), len((b)[0])

	if colsa != rowsb {
		fmt.Println("Matrices cannot be multiplied!")
		os.Exit(3)
	}

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
		go algorithm(&a[currentRow], &currentColData, &wg, elements)
		wg.Wait()

		currentCol += 1

	}

	fmt.Println("Multiply Matrices A and B to get Matrix C:")

	c := make([][]float64, rowsa)
	for i := 0; i < rowsa; i++ {
		row := make([]float64, colsb)
		for j := 0; j < colsb; j++ {
			row[j] = <-elements
		}
		c[i] = row
	}
	result := Matrix(c)
	printMatrix(&result)

	//end
	elapsed := time.Since(start)
	
	//expected answer is: Matrix{{44, 	133, 	193, 	879}, 
	//							 {92, 	149, 	194, 	882}, 
	//							 {134,  490, 	730, 	2004}, 
	//							 {47, 	79, 	104, 	407}, 
	//							 {60, 	263, 	399, 	809}}

	fmt.Print("\nFinished. Elapsed Time: ", elapsed)
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