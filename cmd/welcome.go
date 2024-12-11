/*
Copyright © 2024 the0xsec the0xsec@gmail.com
*/

package cmd

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"strings"

	"github.com/spf13/cobra"
	"golang.org/x/term"
)

var welcomeCmd = &cobra.Command{
	Use:   "welcome",
	Short: "A brief welcome page that directs the user.",
	Long:  `This is a entry point for interactions that comes with a cookies browser.`,
	Run: func(cmd *cobra.Command, args []string) {
		displayWelcome()
	},
}

func init() {
	rootCmd.AddCommand(welcomeCmd)
}

func displayWelcome() {
	width, _, err := term.GetSize(int(os.Stdout.Fd()))
	if err != nil {
		width = 80
	}

	message := "Welcome to the gaia. Pray to the computer god."
	description := "The tool helps manage and understand what is going on in your environment."
	footer := "Press enter to continue or to backspace to exit the program"

	boxWidth := width - 4
	padding := (boxWidth - len(message)) / 2

	fmt.Println(strings.Repeat("=", boxWidth+4))
	fmt.Printf("||%s%s%s||\n", strings.Repeat(" ", padding), message, strings.Repeat(" ", boxWidth-len(message)-padding))
	fmt.Printf("||%s%s%s||\n", strings.Repeat(" ", padding), description, strings.Repeat(" ", boxWidth-len(description)-padding))
	fmt.Printf("||%s%s%s||\n", strings.Repeat(" ", padding), footer, strings.Repeat(" ", boxWidth-len(footer)-padding))
	fmt.Println(strings.Repeat("=", boxWidth+4))

	waitForInput()
}

func waitForInput() {
	oldState, err := term.MakeRaw(int(os.Stdin.Fd()))
	if err != nil {
		fmt.Println("Could not Set Terminal to RAW:", err)
	}
	defer term.Restore(int(os.Stdin.Fd()), oldState)
	reader := bufio.NewReader(os.Stdin)
	fmt.Print("Press Enter to continue or Backspace to return to the previous page.")

	for {
		char, err := reader.ReadByte()
		if err != nil {
			fmt.Println("Error reading the input:", err)
			return
		}

		if char == 127 {
			displayPreviousPage()
			return
		} else if char == 13 {
			fmt.Println("\n\n\n")
			promptForPassword()
			displayNextPage()
			return
		}
	}
}

func promptForPassword() {
	fmt.Println("Please Enter Your Password: ")

	password, err := term.ReadPassword(int(os.Stdin.Fd()))
	if err != nil {
		fmt.Println("Error reading password:", err)
		return
	}

	credential, err := getKeychainCredential("doppler")
	if err != nil {
		fmt.Println(err)
		return
	}

	if credential == "" {
		fmt.Println("No credential found. Would you like to create one? (y/n): ")
		response, _ := bufio.NewReader(os.Stdin).ReadString('\n')
		if strings.TrimSpace(response) == "y" {
			createKeychainCredential("doopler", string(password))
		} else {
			fmt.Println("Existing....")
		}
	} else {
		fmt.Println("Credential Found: ", credential)
	}
}

func createKeychainCredential(name, password string) {
	cmd := exec.Command("security", "add-generic-password", "-s", name, "-w", password)
	err := cmd.Run()
	if err != nil {
		fmt.Println("Error Creating Credential:", err)
		return
	}
	fmt.Println("Credential Created Successfully. The Computer God blesses you.")
}

func getKeychainCredential(name string) (string, error) {
	cmd := exec.Command("security", "find-generic-password", "-s", name, "-w")
	output, err := cmd.Output()
	if err != nil {
		if exitErr, ok := err.(*exec.ExitError); ok && exitErr.ExitCode() == 1 {
			return "", nil
		}
		return "", fmt.Errorf("error retrieving credential: %v", err)
	}
	return strings.TrimSpace(string(output)), nil
}

func displayNextPage() {
	fmt.Println("Navigating to the next page")
}

func displayPreviousPage() {
	fmt.Println("\n\nNavigating back to the previous page...")
}
