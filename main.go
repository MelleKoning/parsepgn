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

type PlyCountData struct {
	move  int
	equal int
	sideA int
	sideB int
}

func main() {
	var plyCountData []PlyCountData
	defaultp := true

	file, err := os.Open("t40games.pgn")
	if err != nil {
		fmt.Printf("error occured %v", err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	var plyCountItem PlyCountData
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
		}
		// do not have to check this side...
		//if pgnlinetype == Black {
		if pgnlinetype == PlyCount {
			// search holder or create new holder
			plyCountItem = GetPlyCountItem(&plyCountData, plyCountItem, line)
		}
		if pgnlinetype == Result {
			fmt.Println(line) // token in unicode-char
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
			}
		}
		fmt.Println("new or updated item:")
		fmt.Println(plyCountItem)
		fmt.Println("...")
		fmt.Println(plyCountData)
		fmt.Println("---")

	}
	fmt.Println("Final results:")
	fmt.Println(plyCountData)
}

// getPlyCountItem should return the existing item based on movecount or return a new one
// we are adding a new item if needed within the function therefore pass address of PlyCountdata
func GetPlyCountItem(plyCountData *[]PlyCountData, plyCountItem PlyCountData, plycountline string) PlyCountData {
	plycountfound := GetPlyCount(plycountline)
	// search for correct item
	for _, pc := range *plyCountData {
		if pc.move == plycountfound {
			// update data from given plyCountItem
			pc.equal = plyCountItem.equal
			pc.sideA = plyCountItem.sideA
			pc.sideB = plyCountItem.sideB
			return pc
		}
	}
	newPc := &PlyCountData{
		move:  plycountfound,
		equal: plyCountItem.equal,
		sideA: plyCountItem.sideA,
		sideB: plyCountItem.sideB,
	}
	*plyCountData = append(*plyCountData, *newPc)
	return *newPc
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
