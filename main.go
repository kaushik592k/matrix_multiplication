package main

import (
	"fmt"
	"sync"
)

// Multiply two matrices in parallel
func parallelMatrixMultiply(A, B [][]int) [][]int {
	// Get dimensions
	rowsA := len(A)
	colsA := len(A[0])
	colsB := len(B[0])

	// Initialize result matrix
	C := make([][]int, rowsA)
	for i := range C {
		C[i] = make([]int, colsB)
	}

	// Create a wait group to synchronize goroutines
	var wg sync.WaitGroup

	// Create a channel to collect results
	resultChannel := make(chan [3]int, rowsA*colsB) // [row, col, value]

	// Dispatch goroutines for matrix multiplication
	for i := 0; i < rowsA; i++ {
		wg.Add(1)
		go func(row int) {
			defer wg.Done()
			for j := 0; j < colsB; j++ {
				sum := 0
				for k := 0; k < colsA; k++ {
					sum += A[row][k] * B[k][j]
				}
				resultChannel <- [3]int{row, j, sum} // Send computed value to channel
			}
		}(i)
	}

	// start another goroutine which will Close result channel when all goroutines are done
	go func() {
		wg.Wait()
		close(resultChannel)
	}()

	// Collect results from channel and construct result matrix
	for result := range resultChannel {
		row, col, value := result[0], result[1], result[2]
		C[row][col] = value
	}

	return C
}

func printMatrix(M [][]int) {
	for _, row := range M {
		for _, val := range row {
			fmt.Printf("%4d ", val)
		}
		fmt.Println()
	}
}

func main() {
	A := [][]int{ // declare matrix as a 2d slice
		{1, 2, 3},
		{4, 5, 6},
		{7, 8, 9},
	}

	B := [][]int{
		{9, 8, 7},
		{6, 5, 4},
		{3, 2, 1},
	}

	C := parallelMatrixMultiply(A, B)
	fmt.Println("Matrix A:")
	printMatrix(A)
	fmt.Println("\nMatrix B:")
	printMatrix(B)
	fmt.Println("\nResultant Matrix (C = A * B):")
	printMatrix(C)
}
