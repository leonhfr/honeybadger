package uci

import (
	"strconv"
	"strings"
	"time"

	"github.com/notnil/chess"
)

const fenFields = 6

// Parse parses UCI commands and returns a Command object.
func Parse(command []string) Command {
	var index int
top:
	switch command[index] {
	case "uci":
		return CommandUCI{}
	case "debug":
		if len(command) > 1 {
			return CommandDebug{on: command[index+1] == "on"}
		}
	case "isready":
		return CommandIsReady{}
	case "setoption":
		if len(command) > 1 {
			return parseCommandSetOption(command[index+1:])
		}
	case "ucinewgame":
		return CommandUCINewGame{}
	case "position":
		if len(command) > 1 {
			return parseCommandPosition(command[index+1:])
		}
	case "go":
		if len(command) > 1 {
			return parseCommandGo(command[index+1:])
		}
	case "stop":
		return CommandStop{}
	case "quit":
		return CommandQuit{}
	default:
		if len(command) == index+1 {
			break
		}
		// unknown commands should be ignored
		index++
		goto top
	}
	return nil
}

// parseCommandSetOption parses setoption UCI commands
func parseCommandSetOption(command []string) CommandSetOption {
	var c CommandSetOption
	if len(command) >= 4 && command[0] == "name" && command[2] == "value" {
		c.name = command[1]
		c.value = command[3]
	}
	return c
}

// parseCommandPosition parses position UCI commands
func parseCommandPosition(command []string) CommandPosition {
	var notation chess.UCINotation
	var c CommandPosition
	var index int

	if command[0] == "startpos" {
		c.startPos = true
		index = 1
	} else if command[0] == "fen" && len(command) >= fenFields+1 {
		c.fen = strings.Join(command[1:fenFields+1], " ")
		index = fenFields + 1
	}

	if len(command) > index && command[index] == "moves" {
		for index++; index < len(command); index++ {
			move, _ := notation.Decode(nil, command[index])
			c.moves = append(c.moves, move)
		}
	}

	return c
}

// parseCommandGo parses go UCI commands
func parseCommandGo(command []string) CommandGo {
	var notation chess.UCINotation
	var c CommandGo

	for index := 0; index < len(command); index++ {
		switch command[index] {
		case "wtime":
			if len(command) >= index+1 {
				t, _ := strconv.Atoi(command[index+1])
				c.input.WhiteTime = time.Duration(t) * time.Millisecond
				index++
			}
		case "btime":
			if len(command) >= index+1 {
				t, _ := strconv.Atoi(command[index+1])
				c.input.BlackTime = time.Duration(t) * time.Millisecond
				index++
			}
		case "winc":
			if len(command) >= index+1 {
				t, _ := strconv.Atoi(command[index+1])
				c.input.WhiteIncrement = time.Duration(t) * time.Millisecond
				index++
			}
		case "binc":
			if len(command) >= index+1 {
				t, _ := strconv.Atoi(command[index+1])
				c.input.BlackIncrement = time.Duration(t) * time.Millisecond
				index++
			}
		case "movestogo":
			if len(command) >= index+1 {
				n, _ := strconv.Atoi(command[index+1])
				c.input.MovesToGo = n
				index++
			}
		case "searchmoves":
			if len(command) >= index+1 {
				for index++; index < len(command); index++ {
					move, _ := notation.Decode(nil, command[index])
					c.input.SearchMoves = append(c.input.SearchMoves, move)
				}
				return c
			}
		case "depth":
			if len(command) >= index+1 {
				n, _ := strconv.Atoi(command[index+1])
				c.input.Depth = n
				index++
			}
		case "nodes":
			if len(command) >= index+1 {
				n, _ := strconv.Atoi(command[index+1])
				c.input.Nodes = n
				index++
			}
		case "movetime":
			if len(command) >= index+1 {
				t, _ := strconv.Atoi(command[index+1])
				c.input.MoveTime = time.Duration(t) * time.Millisecond
				index++
			}
		case "infinite":
			c.input.Infinite = true
		}
	}

	return c
}
