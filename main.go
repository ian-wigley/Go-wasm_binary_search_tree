package main

import (
	"fmt"
	"strings"
	"syscall/js"
)

// Struct to hold each word from the .PTX file
type Node struct {
	value  int
	word   string
	pLeft  *Node
	pRight *Node
}

// Struct to hold the word frequency information
type Statistics struct {
	totalWords        int
	distinctWords     int
	uniqueWords       int
	multiplyUsedWords int
}

var count = 0

func main() {
	js.Global().Set("searchTree", binarySearchTree())
	<-make(chan bool)
}

func binarySearchTree() js.Func {
	jsonfunc := js.FuncOf(func(this js.Value, args []js.Value) interface{} {
		if len(args) != 1 {
			result := map[string]interface{}{
				"error": "Invalid no of arguments passed",
			}
			return result
		}
		jsDoc := js.Global().Get("document")
		if !jsDoc.Truthy() {
			result := map[string]interface{}{
				"error": "Unable to get document object",
			}
			return result
		}
		ouputTextArea := jsDoc.Call("getElementById", "textoutput")
		if !ouputTextArea.Truthy() {
			result := map[string]interface{}{
				"error": "Unable to get output text area",
			}
			return result
		}
		inputText := args[0].String()

		var freqOrder *Node
		freqOrder = nil

		stats := Statistics{
			totalWords:        0,
			distinctWords:     0,
			uniqueWords:       0,
			multiplyUsedWords: 0,
		}

		tree := splitInputText(inputText, nil)

		if tree == nil {
			fmt.Println("Error, can't open file")
		} else {
			temp := ""
			display(tree, &temp)
			freqOrder = copyTree(tree, freqOrder)
			calculateStats(freqOrder, &stats)
			displayStats(&stats, ouputTextArea, temp)
		}
		return nil
	})
	return jsonfunc
}

// Function to split each line of text & populate the tree
func splitInputText(inputText string, tree *Node) *Node {
	temp := strings.Split(inputText, "\n")
	for i := 0; i < len(temp); i++ {
		tree = add(tree, temp[i])
	}
	return tree
}

// Function to create a new Node
func add(tree *Node, word string) *Node {
	tempPtr := Node{
		value:  1,
		word:   word,
		pLeft:  nil,
		pRight: nil,
	}
	count += 1
	return addNode(tree, &tempPtr)
}

// Function to add the words and frequencies to the newly created Node
func addNode(tree *Node, toAdd *Node) *Node {
	if tree == nil {
		return toAdd
	}
	if toAdd.word != "" {
		// Insert recursively into Tree
		if toAdd.word < tree.word && toAdd.word != tree.word {
			tree.pLeft = addNode(tree.pLeft, toAdd)
			return tree
		}

		if toAdd.word > tree.word && toAdd.word != tree.word {
			tree.pRight = addNode(tree.pRight, toAdd)
			return tree
		}

		// If the word already exists then don't add it, just increment the frequency counter
		if toAdd.word == tree.word {
			tree.value++
			return tree
		}
	}
	return tree
}

// Function to display the Words tree to the screen
func display(tree *Node, temp *string) {
	if tree != nil {
		display(tree.pLeft, temp)
		*temp += fmt.Sprintf("(%v)%v\n", tree.value, tree.word)
		display(tree.pRight, temp)
	}
}

// Function to copy the alphabetically ordered tree into a frequency ordered tree
func copyTree(tree *Node, freqOrder *Node) *Node {
	if tree != nil {
		freqOrder = addFreq(freqOrder, tree)
		copyTree(tree.pLeft, freqOrder)
		copyTree(tree.pRight, freqOrder)
	}
	return freqOrder
}

// Function to create a new Node
func addFreq(freqOrder *Node, tree *Node) *Node {
	tempPtr1 := Node{
		value:  tree.value,
		word:   tree.word,
		pLeft:  nil,
		pRight: nil,
	}
	return addFreqNode(freqOrder, &tempPtr1)
}

// Function to add the words and frequencies to the newly created Node
func addFreqNode(freqOrder *Node, toAdd *Node) *Node {
	if freqOrder == nil {
		return toAdd
	} else {
		// Insert recursively into frequency tree
		if toAdd.value < freqOrder.value {
			freqOrder.pLeft = addFreqNode(freqOrder.pLeft, toAdd)
			return freqOrder
		} else if toAdd.value > freqOrder.value {
			freqOrder.pRight = addFreqNode(freqOrder.pRight, toAdd)
			return freqOrder
		}

		// If the frequency already exists then
		if toAdd.value == freqOrder.value {
			// Insert recursively
			if toAdd.word < freqOrder.word {
				freqOrder.pLeft = addFreqNode(freqOrder.pLeft, toAdd)
				return freqOrder
			} else if toAdd.word > freqOrder.word {
				freqOrder.pRight = addFreqNode(freqOrder.pRight, toAdd)
				return freqOrder
			}
		}
	}
	return freqOrder
}

// Function to step through the Frequency ordered tree & gather up the required values
func calculateStats(freqOrder *Node, stats *Statistics) {
	if freqOrder != nil {
		calculateStats(freqOrder.pLeft, stats)
		stats.distinctWords++
		stats.totalWords += freqOrder.value
		if freqOrder.value == 1 {
			stats.uniqueWords++
		} else if freqOrder.value > 1 {
			stats.multiplyUsedWords++
		}
		calculateStats(freqOrder.pRight, stats)
	}
}

// Function to calculate and display the statistical values
func displayStats(stats *Statistics, jsonOuputTextArea js.Value, temp string) {
	var distinctPercent float32 = 0
	var uniquePercent float32 = 0
	var uniquePercentDistinct float32 = 0
	distinctPercent = (float32)(stats.distinctWords*100) / (float32)(stats.totalWords)
	uniquePercent = (float32)(stats.uniqueWords*100) / (float32)(stats.totalWords)
	uniquePercentDistinct = (float32)(stats.uniqueWords*100) / (float32)(stats.distinctWords)

	temp += fmt.Sprintf("Total Words found = %v\n", stats.totalWords)
	temp += fmt.Sprintf("Distinct Words found = %v\n", stats.distinctWords)
	temp += fmt.Sprintf("Multiply-Used Words found = %v\n", stats.multiplyUsedWords)
	temp += fmt.Sprintf("Unique Words found = %v\n", stats.uniqueWords)
	temp += fmt.Sprintf("Distinct Words as a percentage of words = %v\n", distinctPercent)
	temp += fmt.Sprintf("Unique Words as a percentage of words = %v\n", uniquePercent)
	temp += fmt.Sprintf("Unique Words as a percentage of distinct words = %v\n", uniquePercentDistinct)
	jsonOuputTextArea.Set("value", temp)
}
