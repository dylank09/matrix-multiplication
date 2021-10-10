package main

import (
	"fmt"
	"os"
	"sync"
	"time"
)

type Matrix [][]float64

func algorithm1(row, col *[]float64, wg *sync.WaitGroup, elements chan<- float64) {
	sum := 0.0
	for i := range *row {
		sum += (*row)[i] * (*col)[i]
	}
	elements <- sum
	wg.Done()
}


func algorithm2(row *[]float64, cols *Matrix, wg *sync.WaitGroup, rows chan<- []float64) {
	
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
	rows <- result
	wg.Done()

}

func main() {
	a := Matrix{{7, 8, 2}, {1, 9, 21}, {34, 14, 8}, {1, 4, 11}, {21, 4, 2}}
    b := Matrix{{2, 11, 17, 21}, {3, 6, 8, 91}, {3, 4, 5, 2}}

	fmt.Println("Matrix A")
	printMatrix(&a)
	
	fmt.Println("Matrix B")
	printMatrix(&b)

	
	rowsa, rowsb := len(a), len(b)
	colsa, colsb := len((a)[0]), len((b)[0])

	if colsa != rowsb {
		fmt.Println("Matrices cannot be multiplied!")
		os.Exit(3)
	}

	fmt.Println("\nFirst Algorithm")

	//start
	start1 := time.Now()


	var wg1 sync.WaitGroup

	numElementsInResultMatrix := rowsa * colsb
	elements := make(chan float64, numElementsInResultMatrix)

	currentRow := 0
	currentCol := 0
	for i := 0; i < numElementsInResultMatrix; i++ {
		if i != 0 && i % colsb == 0 { //*
			currentRow += 1
		} 

		currentColData := make([]float64, rowsb)
		if currentCol >= colsb {
			currentCol = 0
		}
		
		for i := 0; i < rowsb; i++ {
			currentColData[i] = b[i][currentCol]
		}
		wg1.Add(1)
		go algorithm1(&a[currentRow], &currentColData, &wg1, elements)
		wg1.Wait()

		currentCol += 1
	}
	
	fmt.Println("Multiply Matrices A and B to get Matrix C:")

	c1 := make([][]float64, rowsa)
	for i := 0; i < rowsa; i++ {
		row := make([]float64, colsb)
		for j := 0; j < colsb; j++ {
			row[j] = <-elements
		}
		c1[i] = row
	}
	result1 := Matrix(c1)
	printMatrix(&result1)

	elapsed1 := time.Since(start1)
	
	//expected answer is: Matrix{{44, 	133, 	193, 	879}, 
	//							 {92, 	149, 	194, 	882}, 
	//							 {134,  490, 	730, 	2004}, 
	//							 {47, 	79, 	104, 	407}, 
	//							 {60, 	263, 	399, 	809}}

	fmt.Println("Finished algorithm 1. Elapsed Time: ", elapsed1)

	//end algorithm 1

	//start algorithm 2
	
	fmt.Println("\n\nSecond Algorithm")

	start2 := time.Now()

	var wg2 sync.WaitGroup

	numRowsInResultMatrix := rowsa
	rows := make(chan []float64, numRowsInResultMatrix)
	for i := 0; i < numRowsInResultMatrix; i++ {
		
		wg2.Add(1)
		go algorithm2(&a[i], &b, &wg2, rows)
		wg2.Wait()

	}

	fmt.Println("Multiply Matrices A and B to get Matrix C:")

	c2 := make([][]float64, rowsa)
	for i := 0; i < numRowsInResultMatrix; i++ {
		
		c2[i] = <- rows

	}
	result2 := Matrix(c2)
	printMatrix(&result2)

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

