package cmd

import (
	"bufio"
	"fmt"
	"log"
	"os"

	"github.com/dumbogo/chess/config"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(signUpCmd)
}

var urlSignUp string

var signUpCmd = &cobra.Command{
	Use:   "signup",
	Short: "Sign up on chess",
	Long:  "Sign up on chess platform to start playing!",
	Run: func(cmd *cobra.Command, args []string) {
		c, err := config.LoadClientConfiguration()
		if err != nil {
			panic(err)
		}
		reader := bufio.NewReader(os.Stdin)
		fmt.Printf("Visit http://%s to signin and follow steps, click enter when you finish\n HIT ENTER", c.HTTPServerURL)
		reader.ReadLine()

		fmt.Print("Type your token here: ")
		token, _, err := reader.ReadLine()
		if err != nil {
			log.Fatalf("Error: %v\n", err)
		}
		// TODO: add some steps to recognize which provider was used, at the moment we are going to leave it to github only
		c.SetAuthToken(string(token))
	},
}
