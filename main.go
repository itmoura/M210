package main

import (
	"fmt"
)

// Função principal para resolver o PL
func main() {
	//Coeficientes da função objetivo (maximização)
	objectiveFunction := []float64{24, 22, 45}
	for i := 0; i < len(objectiveFunction); i++ {
		objectiveFunction[i] = -objectiveFunction[i]
	}

	// Coeficientes das restrições
	constraints := [][]float64{
		{2, 1, 3},
		{2, 1, 2},
		{1, 1, 1},
	}

	// Lados direitos das restrições
	rightHandSide := []float64{42, 40, 45}

	// Chamando a função simplex para encontrar a solução
	table, err := simplex(objectiveFunction, constraints, rightHandSide)

	if err != nil {
		fmt.Println("Erro ao resolver o PL:", err)
		return
	}

	// RESULTADOS

	var optimalPosition []int
	var priceShadow []float64
	aux := len(table[0]) - len(constraints) - 1
	for i := aux; i < len(table[0])-1; i++ {
		priceShadow = append(priceShadow, table[0][i])
		if table[0][i] == 0 {
			continue
		}

		if len(constraints) < 3 {
			optimalPosition = append(optimalPosition, i-1)
		} else {
			optimalPosition = append(optimalPosition, i)
		}

	}

	// Percorrer as variaveis de lucro LD
	var lastLine []float64
	numberOfColumns := len(table[0]) - 1
	for j := 0; j < len(constraints)+1; j++ {
		lastLine = append(lastLine, table[j][numberOfColumns])
	}

	fmt.Println("\nPonto Ótimo de Operação:")
	fmt.Println("Z = ", lastLine[0])
	for i := 0; i < len(optimalPosition); i++ {
		fmt.Printf("x%d = %.2f\n", optimalPosition[i]-1, lastLine[optimalPosition[i]-1])
	}

	// Preço Sombra
	fmt.Println("\nPreço Sombra:")
	for i := 0; i < len(priceShadow); i++ {
		fmt.Printf("x%d = %.2f\n", i+1, priceShadow[i])
	}

}

func simplex(objectiveFunction []float64, constraints [][]float64, rightHandSide []float64) ([][]float64, error) {
	// montando a tabela inicial
	newTable := append([][]float64{append(objectiveFunction)}, constraints...)

	// adicionando as variáveis de folga
	for i := 0; i < len(constraints); i++ {
		for j := 0; j < len(constraints); j++ {
			if i == j {
				newTable[i+1] = append(newTable[i+1], 1)
			} else {
				newTable[i+1] = append(newTable[i+1], 0)
			}
		}
	}

	// adicionando as variáveis de folga na função objetivo
	for i := 0; i < len(constraints); i++ {
		newTable[0] = append(newTable[0], 0)
	}

	// adicionando LD na tabela
	newTable[0] = append(newTable[0], 0)
	for i := 0; i < len(rightHandSide); i++ {
		newTable[i+1] = append(newTable[i+1], rightHandSide[i])
	}

	// mostrar a tabela inicial
	fmt.Println("\n\nTabela inicial:")
	for i := range newTable {
		fmt.Print("[")
		for _, value := range newTable[i] {
			fmt.Printf("%.2f ", value)
		}
		fmt.Println("]")
	}

	return calcSimplex(newTable, rightHandSide)
}

func calcSimplex(table [][]float64, rightHandSide []float64) ([][]float64, error) {
	// Pegando o mais negativo da função objetivo
	var mostNegative float64
	var mostNegativeIndex int
	for i := 0; i < len(table[0]); i++ {
		if table[0][i] < mostNegative {
			mostNegative = table[0][i]
			mostNegativeIndex = i
		}
	}

	// Verificando o PIVO
	var pivo float64
	var pivoIndex int
	calcPivot := 99999999.99
	for i := 1; i < len(table); i++ {
		if table[i][mostNegativeIndex] == 0 {
			continue
		}
		calcAux := rightHandSide[i-1] / table[i][mostNegativeIndex]
		if calcAux < calcPivot {
			calcPivot = calcAux
			pivo = table[i][mostNegativeIndex]
			pivoIndex = i
		}
	}

	// calcular a nova linha do pivo
	var newPivoLine []float64
	for i := 0; i < len(table[pivoIndex]); i++ {
		if i == len(table[pivoIndex])-1 {
			newPivoLine = append(newPivoLine, calcPivot)
			continue
		}
		newPivoLine = append(newPivoLine, table[pivoIndex][i]/pivo)
	}

	var newRightHandSide []float64
	// calcular as novas linhas
	var newTable [][]float64
	for i := 0; i < len(table); i++ {
		if i != pivoIndex {
			var newLine []float64
			for j := 0; j < len(table[i]); j++ {
				newLine = append(newLine, table[i][j]-(newPivoLine[j]*table[i][mostNegativeIndex]))
				if j == len(table[i])-1 {
					newRightHandSide = append(newRightHandSide, newLine[j])
				}
			}
			newTable = append(newTable, newLine)
		} else {
			newTable = append(newTable, newPivoLine)
		}
	}

	// verificar se é ótimo
	isOptimal := true
	for i := 0; i < len(newTable[0]); i++ {
		if newTable[0][i] < 0 {
			isOptimal = false
		}
	}

	// mostrar a nova tabela
	fmt.Println("\n\nNova tabela:")
	for i := range newTable {
		fmt.Print("[")
		for _, value := range newTable[i] {
			fmt.Printf("%.2f ", value)
		}
		fmt.Println("]")
	}

	if isOptimal {
		return newTable, nil
	} else {
		return calcSimplex(newTable, newRightHandSide)
	}
}
