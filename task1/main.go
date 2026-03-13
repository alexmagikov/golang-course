package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

func main() {
	reader := bufio.NewReader(os.Stdin)
	fmt.Print("Enter owner name: ")
	name, _ := reader.ReadString('\n')
	name = strings.TrimSpace(name)

	fmt.Print("Enter repo name: ")
	repo, _ := reader.ReadString('\n')
	repo = strings.TrimSpace(repo)

	url := CreatePath(name, repo)

	info, err := GetRepoInfo(url)
	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Printf("\n====Repo info=====\n")

	fmt.Printf("Repo Name: %s\n", info.Name)
	fmt.Printf("Repo description: %s\n", info.Description)
	fmt.Printf("Repo forks: %d\n", info.Forks)
	fmt.Printf("Repo stars: %d\n", info.Stars)
	fmt.Printf("Repo created date: %s\n", info.CreatedAt)
}
