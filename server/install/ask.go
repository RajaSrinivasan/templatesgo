package install

import (
	"bufio"
	"fmt"
	"os"
)

func Ask(prompt string, def string) string {
	scanner := bufio.NewScanner(os.Stdin)
	fmt.Print(prompt)
	if len(def) > 0 {
		fmt.Printf("[%s]", def)
	}
	fmt.Print(" : ")
	scanner.Scan()
	text := scanner.Text()
	if len(text) == 0 {
		return def
	}
	return text
}
