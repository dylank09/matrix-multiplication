package main

import (
	"fmt"
	"os"
	"sync"
	"time"
)

type Matrix [][]float64

func algorithm(row *[]float64, cols *Matrix, wg *sync.WaitGroup, rows chan<- []float64) {
	
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

	//start
	start := time.Now()

	rowsa, rowsb := len(a), len(b)
	colsa, _ := len((a)[0]), len((b)[0])

	if colsa != rowsb {
		fmt.Println("Matrices cannot be multiplied!")
		os.Exit(3)
	}

	var wg sync.WaitGroup

	numRowsInResultMatrix := rowsa
	rows := make(chan []float64, numRowsInResultMatrix)
	for i := 0; i < numRowsInResultMatrix; i++ {
		
		wg.Add(1)
		go algorithm(&a[i], &b, &wg, rows)
		wg.Wait()

	}

	fmt.Println("Multiply Matrices A and B to get Matrix C:")

	c := make([][]float64, rowsa)
	for i := 0; i < numRowsInResultMatrix; i++ {
		
		c[i] = <- rows

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