package main

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"sort"
	"strconv"

	"go.i3wm.org/i3/v4"
)

// getWorkspaceNumbers retrieves and sorts workspace numbers from i3.
func getWorkspaceNumbers(workspaces []i3.Workspace) ([]int, error) {
	numbers := make([]int, len(workspaces))
	for i, ws := range workspaces {
		numbers[i] = int(ws.Num)
	}
	sort.Ints(numbers)
	return numbers, nil
}

// findNextWorkspaceNumber finds the next available workspace number.
func findNextWorkspaceNumber(numbers []int) int {
	newNum := 1
	for _, num := range numbers {
		if num > newNum {
			break
		}
		if num == newNum {
			newNum++
		}
	}
	return newNum
}

// moveContainerToWorkspace moves the focused window to the specified workspace.
func moveContainerToWorkspace(workspaceNum int) error {
	moveCmd := fmt.Sprintf("move container to workspace number %d", workspaceNum)
	if err := exec.Command("i3-msg", moveCmd).Run(); err != nil {
		return fmt.Errorf("failed to move container to workspace %d: %w", workspaceNum, err)
	}
	return nil
}

// switchToWorkspace switches to the specified workspace.
func switchToWorkspace(workspaceNum int) error {
	switchCmd := fmt.Sprintf("workspace number %d", workspaceNum)
	if err := exec.Command("i3-msg", switchCmd).Run(); err != nil {
		return fmt.Errorf("failed to switch to workspace %d: %w", workspaceNum, err)
	}
	return nil
}

// traverseAndMoveWindows recursively traverses the i3 tree and moves windows to new workspace numbers.
func traverseAndMoveWindows(node *i3.Node, wsNum int, wsMap map[int]int) error {
	if node.Type == "workspace" {
		// Only parse numeric workspace names
		if num, err := strconv.Atoi(node.Name); err == nil {
			wsNum = num
		} else {
			// Skip non-numeric workspaces like __i3_scratch
			wsNum = 0
		}
	}
	if node.Window != 0 && wsNum > 0 && wsMap[wsNum] != wsNum {
		// Move window to new workspace
		moveCmd := fmt.Sprintf("[con_id=%d] move window to workspace number %d", node.ID, wsMap[wsNum])
		if err := exec.Command("i3-msg", moveCmd).Run(); err != nil {
			return fmt.Errorf("failed to move window %d to workspace %d: %w", node.ID, wsMap[wsNum], err)
		}
	}
	// Traverse child nodes
	for _, child := range node.Nodes {
		if err := traverseAndMoveWindows(child, wsNum, wsMap); err != nil {
			return err
		}
	}
	// Traverse floating nodes
	for _, child := range node.FloatingNodes {
		if err := traverseAndMoveWindows(child, wsNum, wsMap); err != nil {
			return err
		}
	}
	return nil
}

// renumberWorkspaces reassigns workspace numbers to be consecutive, moving windows and restoring focus.
func renumberWorkspaces() error {
	// Get workspaces
	workspaces, err := i3.GetWorkspaces()
	if err != nil {
		return fmt.Errorf("failed to get workspaces: %w", err)
	}

	// Find the currently focused workspace
	focusedWorkspaceNum := 0
	for _, ws := range workspaces {
		if ws.Focused {
			focusedWorkspaceNum = int(ws.Num)
			break
		}
	}

	// Create mapping of old to new workspace numbers
	wsNums, err := getWorkspaceNumbers(workspaces)
	if err != nil {
		return err
	}
	wsMap := make(map[int]int)
	for newNum, oldNum := range wsNums {
		wsMap[oldNum] = newNum + 1
	}

	// Get the i3 tree
	tree, err := i3.GetTree()
	if err != nil {
		return fmt.Errorf("failed to get i3 tree: %w", err)
	}

	// Move windows to new workspace numbers
	if err := traverseAndMoveWindows(tree.Root, 0, wsMap); err != nil {
		return err
	}

	// Rename workspaces to new numbers
	for _, ws := range workspaces {
		if newNum, ok := wsMap[int(ws.Num)]; ok && int(ws.Num) != newNum {
			renameCmd := fmt.Sprintf("rename workspace number %d to %d", ws.Num, newNum)
			if err := exec.Command("i3-msg", renameCmd).Run(); err != nil {
				return fmt.Errorf("failed to rename workspace %d to %d: %w", ws.Num, newNum, err)
			}
		}
	}

	// Switch to the renumbered focused workspace
	if focusedWorkspaceNum > 0 {
		if newNum, ok := wsMap[focusedWorkspaceNum]; ok {
			return switchToWorkspace(newNum)
		}
	}

	return nil
}

// help displays usage information for the script.
func help() {
	fmt.Println("Usage: nextwrk [--switch] [--renumber]")
	fmt.Println("	[no args]		Move focused container to the next free workspace")
	fmt.Println("	--switch		Move container to next free workspace and switch to it")
	fmt.Println("	--renumber		Renumber all workspaces to remove gaps and restore focus")
	fmt.Println("	[any other args]	Show this help message")
	fmt.Println("------------------------------")
	fmt.Println("Built with Go for i3WM. MIT License. github.com/arrayindex-dev/nextwrk")
}

func main() {
	// Parse command-line flags
	doSwitch := false
	renumber := false
	for _, arg := range os.Args[1:] {
		switch arg {
		case "--switch":
			doSwitch = true
		case "--renumber":
			renumber = true
		default:
			help()
			return
		}
	}

	if renumber {
		// Renumber workspaces to remove gaps and restore focus
		if err := renumberWorkspaces(); err != nil {
			log.Fatal(err)
		}
		return
	}

	// Original functionality: move container and optionally switch workspace
	workspaces, err := i3.GetWorkspaces()
	if err != nil {
		log.Fatal(err)
	}
	numbers, err := getWorkspaceNumbers(workspaces)
	if err != nil {
		log.Fatal(err)
	}

	newNum := findNextWorkspaceNumber(numbers)

	if err := moveContainerToWorkspace(newNum); err != nil {
		log.Fatal(err)
	}

	if doSwitch {
		if err := switchToWorkspace(newNum); err != nil {
			log.Fatal(err)
		}
	}
}