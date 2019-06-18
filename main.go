package main

import (
	//"io/ioutil"
	"bufio"
	"fmt"

	//"log"
	"os"
	"strconv"
	"strings"
)

type PgnLineType int

const (
	Event PgnLineType = iota
	Site
	Date
	Round
	White
	Black
	Result
	FEN
	PlyCount
)

func (pgnlinetype PgnLineType) String() string {
	return [...]string{"Event", "Site", "Date", "Round", "White", "Black", "Result", "FEN", "PlyCount"}[pgnlinetype]
}

type PlyCountItem struct {
	move  int
	equal int
	sideA int
	sideB int
}

//GetNewPlyCountData returns a new object
func GetNewPlyCountItem() *PlyCountItem {
	return &PlyCountItem{
		move:  0,
		equal: 0,
		sideA: 0,
		sideB: 0,
	}
}
func main() {
	var plyCountData map[int]PlyCountItem
	plyCountData = make(map[int]PlyCountItem)
	defaultp := true
	numberofresults := 0
	file, err := os.Open("t40games.pgn")
	if err != nil {
		fmt.Printf("error occured %v", err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	var plyCountItem PlyCountItem
	for scanner.Scan() { // internally, it advances token based on sperator

		line := scanner.Text()
		pgnlinetype := GetPgnLineType(line)
		//fmt.Println(line) // token in unicode-char

		if pgnlinetype == White {
			if strings.Index(line, "889") != -1 {
				defaultp = false
			} else {
				defaultp = true
			}
			plyCountItem = *GetNewPlyCountItem()
		}
		// do not have to check this side...
		//if pgnlinetype == Black {
		if pgnlinetype == Result {
			numberofresults++
			fmt.Println(numberofresults) // token in unicode-char
			switch line {

			case "[Result \"1-0\"]":
				{
					if defaultp {
						plyCountItem.sideA++
					} else {
						plyCountItem.sideB++
					}
				}
			case "[Result \"0-1\"]":
				{
					if defaultp {
						plyCountItem.sideB++
					} else {
						plyCountItem.sideA++
					}
				}
			case "[Result \"1/2-1/2\"]":
				{
					plyCountItem.equal++
				}
			default:
				fmt.Println(line)
			}
		}
		if pgnlinetype == PlyCount {
			// search holder or create new holder
			plyCountItem, plyCountData = AddPlyCountItem(plyCountData, plyCountItem, line)
			//fmt.Println("new or updated item:")
			//fmt.Println(plyCountItem)
		}

		//fmt.Print(".")

	}
	//fmt.Println("Final results:")
	//fmt.Println(plyCountData)
	WriteCsv(plyCountData)
}

func WriteCsv(plyCountData map[int]PlyCountItem) {
	for i := 0; i < 200; i++ {
		item, exists := plyCountData[i]
		if !exists {
			fmt.Println(i, ",0,0,0")
		} else {
			fmt.Printf("%v,%v,%v,%v", item.move, item.equal, item.sideA, item.sideB)
			fmt.Println()
		}
	}

}

// getPlyCountItem should return the existing item based on movecount or return a new one
// we are adding a new item if needed within the function therefore pass address of PlyCountdata
func AddPlyCountItem(plyCountData map[int]PlyCountItem, plyCountItem PlyCountItem, plycountline string) (PlyCountItem, map[int]PlyCountItem) {
	plycountnumber := GetPlyCount(plycountline)
	// search for correct item
	pcItem, exists := plyCountData[plycountnumber]
	if !exists {
		pcItem = *GetNewPlyCountItem()
	}
	pcItem.move = plycountnumber
	pcItem.equal = pcItem.equal + plyCountItem.equal
	pcItem.sideA = pcItem.sideA + plyCountItem.sideA
	pcItem.sideB = pcItem.sideB + plyCountItem.sideB

	plyCountData[plycountnumber] = pcItem
	return pcItem, plyCountData
}

// GetPlyCount returns the number from a string looking like [whatever "88"]
func GetPlyCount(line string) int {
	quoteIndex := strings.Index(line, "\"")
	lastquoteIndex := strings.LastIndex(line, "\"")
	number, err := strconv.Atoi(line[quoteIndex+1 : lastquoteIndex])
	if err != nil {
		fmt.Printf("failed getting plyCount from %v", line)
	}
	return number
}

//GetPgnLineType returns the type for this line
func GetPgnLineType(line string) PgnLineType {
	var parts []string
	parts = strings.Split(line, " ")
	switch parts[0] {
	case "[Event":
		return Event
	case "[Site":
		return Site
	case "[Date":
		return Date
	case "[White":
		return White
	case "[Black":
		return Black
	case "[Result":
		return Result
	case "[PlyCount":
		return PlyCount
	default:
		return FEN

	}
}
