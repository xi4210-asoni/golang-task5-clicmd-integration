/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"bufio"
	"fmt"
	"net/url"
	"os"
	"strings"

	"github.com/gorilla/websocket"
	"github.com/spf13/cobra"
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "chat",
	Short: "Chat binary used to run the commands for connecting with the websockets",
	Long: `A longer description that spans multiple lines and likely contains
examples and usage of using your application. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	// Uncomment the following line if your bare application
	// has an action associated with it:
	// Run: func(cmd *cobra.Command, args []string) { },
}

type ReqRespBody struct {
	Event string      `json:"event" mapstructure:"event"`
	Data  interface{} `json:"data" mapstructure:"data"`
}

type ChatMessage struct {
	Message string `json:"message"`
}

var chatBeginCmd = &cobra.Command{
	Use:   "begin",
	Short: "use to start the begin command directly",
	Run: func(cmd *cobra.Command, args []string) {
		u := url.URL{Scheme: "ws", Host: "localhost:8092", Path: "/socket"}

		//? get username from the flag
		username, err := cmd.Flags().GetString("user")
		if err != nil {
			fmt.Println(err)
		}

		//set the connection from dialer
		conn, _, err := websocket.DefaultDialer.Dial(u.String(), nil)
		if err != nil {
			fmt.Println(err)
		}

		// defer conn.Close() // Todo: Change this

		//? tologin to the socket

		loginmsg := fmt.Sprintf(`{"event": "login", "data": {"user": "%s"}}`, username)
		if err := conn.WriteMessage(websocket.TextMessage, []byte(loginmsg)); err != nil {
			fmt.Println(err)
		}

		go func() {
			for {
				_, message, err := conn.ReadMessage()
				if err != nil {
					fmt.Println(err)
					os.Exit(-1)
				}
				fmt.Println("Response: ", string(message))
			}
		}()

		// fmt.Println("Press ctrl+c to exit")
		// select {} //holds for the user input
		//TODO Change to use the input from the user multiple times
		//! Is there any specific input type from user
		// var chatusername string
		// var chatusermessage string

		// go func() {
		for {
			scanner := bufio.NewScanner(os.Stdin)
			fmt.Println("Enter message(user: your message): ")
			scanner.Scan()
			chat := scanner.Text()
			if err != nil {
				fmt.Println("Unable to login:", err)
				os.Exit(-1)
			}
			fmt.Scanf("%s", &chat)

			// split into two parts
			// temp := strings.Split(chat, ":")
			// fmt.Println(temp)
			parts := strings.Split(chat, ":")
			if len(parts) != 2 {
				fmt.Println("Invalid input format")
				return
			}
			chatusername := strings.TrimSpace(parts[0])
			// chatmessage := strings.Trim(parts[1])
			chatmessage := parts[1]
			// fmt.Println("Enter user message: ")
			// fmt.Scanf("%s", &chatusermessage)
			// fmt.Println(chatmessage)
			fmt.Println("Sending:", chatmessage, "to", chatusername)
			// jsonBytes, err := json.Marshal(cha)
			// apiChat := &ReqRespBody{
			// 	Event: "chat",
			// 	Data: &ChatMessage{
			// 		Chat:    chatusername,
			// 		Message: chatmessage,
			// 	},
			// }
			// apichat := &ChatMessage{Chat: chatusername, Message: chatmessage}
			// json.Marshal(apiChat)
			chatmsg := fmt.Sprintf(`{"event": "chat", "data": {"user": "%s","message": "%s"}}`, chatusername, chatmessage)
			// fmt.Printf(chatmsg)
			err = conn.WriteMessage(websocket.TextMessage, []byte(chatmsg))
			// err = conn.
			// err = conn.WriteJSON(apiChat)
			if err != nil {
				fmt.Println(err)
				break
			}
		}
		// }()
	},
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	// Here you will define your flags and configuration settings.
	// Cobra supports persistent flags, which, if defined here,
	// will be global for your application.

	// rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.chat.yaml)")

	// Cobra also supports local flags, which will only run
	// when this action is called directly.
	chatBeginCmd.Flags().StringP("user", "u", "test", "Username for chat session")
	rootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
	rootCmd.AddCommand(chatBeginCmd)
}
