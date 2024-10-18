package main

import (
	"fmt"
	"log"
	"os/exec"
	"sort"

	"go.i3wm.org/i3/v4"
)

func main() {
	//get workspaces
	workspaces, err := i3.GetWorkspaces()
	if err != nil {
		log.Fatal(err)
	}

	//extract and sort workspace numbers
	numbers := make([]int, len(workspaces))
	for i, ws := range workspaces {
		numbers[i] = int(ws.Num)
	}
	sort.Ints(numbers)

	//find the next available workspace number
	newNum := 1
	for _, num := range numbers {
		if num > newNum {
			break
		}
		if num == newNum {
			newNum++
		}
	}

	//move the focused window to the new workspace
	moveCmd := fmt.Sprintf("move container to workspace number %d", newNum)
	if err := exec.Command("i3-msg", moveCmd).Run(); err != nil {
		log.Fatal(err)
	}

	//switch to the new workspace
	switchCmd := fmt.Sprintf("workspace number %d", newNum)
	if err := exec.Command("i3-msg", switchCmd).Run(); err != nil {
		log.Fatal(err)
	}
}
