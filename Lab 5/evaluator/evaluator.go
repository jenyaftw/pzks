package evaluator

import (
	"fmt"
	"slices"

	"github.com/jenyaftw/lab1/parser"
	"github.com/jenyaftw/lab1/token"
)

type Evaluator struct{}

func NewEvaluator() *Evaluator {
	return &Evaluator{}
}

const (
	NULL = iota
	MULTIPLICATION
	DIVISION
	ADDITION
	SUBTRACTION
	RESERVED
)

const NUM_LAYERS = 6

type GanttChart struct {
	TopTime           int
	TopOperatorsTime  int
	Time              int
	LastIdx           int
	Operators         []int
	FinishedOperators []int
	S                 [][]int
	Pads              int
}

func NewGanttChart() *GanttChart {
	return &GanttChart{
		S: make([][]int, NUM_LAYERS),
	}
}

func (c GanttChart) PadWithZeroes() {
	for i := 0; i < NUM_LAYERS; i++ {
		c.S[i] = append([]int{0}, c.S[i]...)
		c.S[i] = append(c.S[i], 0)
	}
	c.Pads += 2
}

func (c GanttChart) Print(operators []*parser.TreeNode) {
	fmt.Println("|    T    |    S1    |    S2    |    S3    |    S4    |    S5    |    S6    |")
	fmt.Println("|---------|----------|----------|----------|----------|----------|----------|")
	for i := 0; i < len(c.S[len(c.S)-1]); i++ {
		fmt.Print("|")
		if i < 9 {
			fmt.Print(" ")
		}

		fmt.Printf("   %d    |", i+1)
		for j := 0; j < len(c.S); j++ {
			if i < len(c.S[j]) && operators[c.S[j][i]].Token.Type == token.OperatorType {
				// If the operator is single digit we need to add a space to make it look better
				if operators[c.S[j][i]].Id < 10 {
					fmt.Print("  ")
				} else if operators[c.S[j][i]].Id < 100 {
					fmt.Print(" ")
				}
				fmt.Printf(" %d(%s)   |", operators[c.S[j][i]].Id, operators[c.S[j][i]].Token.Text)
			} else {
				fmt.Print("          |")
			}
		}
		fmt.Println()
	}
	fmt.Println("|----------|----------|----------|----------|----------|----------|----------|")
}

func (c *GanttChart) AddOperation(op int, operators []*parser.TreeNode) {
	c.Operators = append(c.Operators, op)
	c.Time = getOperationTime(op, operators)
	if c.Time > c.TopTime {
		c.TopTime = c.Time
	}

	for i := 0; i < c.Time; i++ {
		c.S[0] = append(c.S[0], op)
	}
	for i := 0; i < c.TopTime-c.Time; i++ {
		c.S[0] = append(c.S[0], 0)
	}

	for i := 1; i < NUM_LAYERS; i++ {
		baseLayer := make([]int, 0)

		for j, v := range c.Operators {
			rangeOps := j + i
			if rangeOps >= len(c.Operators) {
				rangeOps = len(c.Operators) - 1
			}

			topTime := 0
			for k := 0; k <= rangeOps; k++ {
				if getOperationTime(c.Operators[k], operators) > topTime {
					topTime = getOperationTime(c.Operators[k], operators)
				}
			}

			vTime := getOperationTime(v, operators)
			for k := 0; k < vTime; k++ {
				baseLayer = append(baseLayer, v)
			}

			if j < len(c.Operators)-1 {
				for k := 0; k < topTime-vTime; k++ {
					baseLayer = append(baseLayer, 0)
				}
			}
		}

		c.S[i] = baseLayer
	}
}

func (c *GanttChart) PrependZeros(operators []*parser.TreeNode) {
	opIdx := 0
	prependAmt := getOperationTime(c.Operators[opIdx], operators)
	for i := 1; i < NUM_LAYERS; i++ {
		for j := 0; j < prependAmt; j++ {
			c.S[i] = append([]int{0}, c.S[i]...)
		}

		prependAmt += getOperationTime(c.Operators[opIdx], operators)

		if opIdx+1 >= len(c.Operators) {
			opIdx = len(c.Operators) - 1
		}
	}
}

func (c *GanttChart) BalanceZeros() {
	lenLast := len(c.S[len(c.S)-1])
	for i := 0; i < len(c.S)-1; i++ {
		neededZeros := lenLast - len(c.S[i])
		if neededZeros > 0 {
			tempArray := make([]int, neededZeros)
			c.S[i] = append(c.S[i], tempArray...)
		} else {
			for j := 0; j < -neededZeros; j++ {
				c.S[i] = c.S[i][:len(c.S[i])-1]
			}
		}
	}
}

func (c *GanttChart) Append(chart *GanttChart) {
	c.Operators = chart.Operators
	for i := 0; i < len(chart.S); i++ {
		c.S[i] = append(c.S[i], chart.S[i]...)
	}
	c.Pads += chart.Pads
}

func (c *GanttChart) PrintStats(operators []*parser.TreeNode, count int) {
	totalTime := len(c.S[len(c.S)-1])
	fmt.Println("Execution time: ", totalTime)

	sumOfOperators := 0
	for i := range operators {
		sumOfOperators += getOperationTime(i, operators)
	}
	numLayers := len(c.S)
	unpadddedTime := totalTime - count*2
	accel := float64(sumOfOperators*numLayers) / float64(unpadddedTime)
	efficiency := accel / float64(numLayers)

	fmt.Println("Acceleration: ", accel)
	fmt.Println("Efficiency: ", efficiency)
}

func getOperationTime(operation int, operators []*parser.TreeNode) int {
	switch operators[operation].Token.Text {
	case "*":
		return 4
	case "/":
		return 8
	case "+":
		return 1
	case "-":
		return 3
	default:
		return 0
	}
}

func operationToString(op int) string {
	switch op {
	case MULTIPLICATION:
		return "*"
	case DIVISION:
		return "/"
	case ADDITION:
		return "+"
	case SUBTRACTION:
		return "-"
	default:
		return " "
	}
}

func stringToOperation(op string) int {
	switch op {
	case "*":
		return MULTIPLICATION
	case "/":
		return DIVISION
	case "+":
		return ADDITION
	case "-":
		return SUBTRACTION
	default:
		return NULL
	}
}

func findRequired(op int, operators []*parser.TreeNode) []int {
	required := make([]int, 0)
	firstOperator := operators[op]
	left, right := firstOperator.Left, firstOperator.Right
	for i, v := range operators {
		if (v == left || v == right) && v.Token.Type == token.OperatorType {
			required = append(required, i)
		}
	}

	return required
}

func (e *Evaluator) GenerateGanttChart(node parser.TreeNode) (*GanttChart, []*parser.TreeNode, int) {
	chart := NewGanttChart()
	allCharts := make([]*GanttChart, 0)
	reversed := node.ReverseLevelOrderTraversal()
	ops := 0

	for i, v := range reversed {
		if v.Token.Type == token.OperatorType {
			ops += 1
			required := findRequired(i, reversed)

			simpleAdd := func() {
				chart.AddOperation(i, reversed)
				chart.PrependZeros(reversed)
				chart.FinishedOperators = append(chart.FinishedOperators, i)
			}

			if len(required) == 0 {
				simpleAdd()
				continue
			}

			finishedContains1 := slices.Contains(chart.FinishedOperators, required[0])
			finishedContains2 := false

			if len(required) > 1 {
				slices.Contains(chart.FinishedOperators, required[1])
			}

			if finishedContains1 || finishedContains2 {
				operatorsContains1 := slices.Contains(chart.Operators, required[0])
				operatorsContains2 := false
				if len(required) > 1 {
					slices.Contains(chart.Operators, required[1])
				}

				if operatorsContains1 || operatorsContains2 {
					oldChart := chart
					chart.BalanceZeros()
					allCharts = append(allCharts, chart)

					chart = NewGanttChart()
					chart.AddOperation(i, reversed)
					chart.PrependZeros(reversed)
					chart.FinishedOperators = append(chart.FinishedOperators, oldChart.FinishedOperators...)
					chart.FinishedOperators = append(chart.FinishedOperators, i)
					chart.Pads = oldChart.Pads
				} else {
					simpleAdd()
				}
			} else {
				simpleAdd()
			}
		}
	}

	if len(chart.Operators) > 0 {
		allCharts = append(allCharts, chart)
	}

	totalChart := NewGanttChart()
	for _, v := range allCharts {
		v.PadWithZeroes()
		totalChart.Append(v)
	}

	return totalChart, reversed, len(allCharts)
}
