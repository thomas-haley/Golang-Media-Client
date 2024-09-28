package main

import (
	"encoding/gob"
	"encoding/json"
	"fmt"
	"net"
	"os"
	"os/exec"
	"thaley/fileServerClient/menu"
	"thaley/fileServerClient/messaging"
	"thaley/fileServerClient/util"
)

var ServerIP string = "73.178.239.141"
var ServerPort string = "9999"

func main() {
	os.Mkdir("./download/", os.ModePerm)

	// //Get IP from user
	// fmt.Println("Enter the IP of the server you want to connect to: ")
	// fmt.Scanln(&ServerIP)

	//Connect to the server
	c, err := net.Dial("tcp", ServerIP+":"+ServerPort)

	if err != nil {
		fmt.Println(err)
		fmt.Scanln()
		return
	}

	var msg []byte
	//Receive message from server
	err = gob.NewDecoder(c).Decode(&msg)
	// fmt.Println(msg)
	var serverConfig util.ServerConfig

	json.Unmarshal(msg, &serverConfig)

	if serverConfig.ALLOW_ACCESS {
		fmt.Println("Access Allowed")
		MainMenu(c, serverConfig)
	} else {
		fmt.Println("Access denied")
	}
}

func MainMenu(conn net.Conn, config util.ServerConfig) {
	for {
		fmt.Println("\n\n============================")
		fmt.Println("Main Menu")
		fmt.Println("============================")
		if config.ALLOW_FILE_READ {
			fmt.Println("[1] Read files from server")
		}
		if config.ALLOW_FILE_WRITE {
			fmt.Println("[2] Upload files to server")
		}
		if config.ALLOW_FILE_DOWNLOAD {
			fmt.Println("[3] Download files from server")
		}
		if config.ALLOW_FILE_TRANSCODE {
			fmt.Println("[4] Transcode files via server")
		}
		if config.ALLOW_WATCH_PARTY {
			fmt.Println("[5] View Media Rooms")
		}

		fmt.Println("[0] Exit program")
		var userOp string
		fmt.Scanln(&userOp)

		switch userOp {
		case "1":
			//Handle file read option
			// msg, _ := json.Marshal("1")
			gob.NewEncoder(conn).Encode("1")
			fileRead(conn)
			// conn.Write(msg)
			break
		case "2":
			//Handle file select/curl to server option, while also opening way for server to read
			gob.NewEncoder(conn).Encode("2")
			break
		case "3":
			//Handle file read/download from available
			gob.NewEncoder(conn).Encode("3")
			fileDownload(conn)
			break
		case "4":
			//Handle file upload/send to ffmpeg/and download again
			gob.NewEncoder(conn).Encode("4")
			break
		case "5":
			//Handle file upload/send to ffmpeg/and download again
			gob.NewEncoder(conn).Encode("5")
			menu.MediaRoomMainMenu(conn)
			break
		case "0":
			//Exit program
			gob.NewEncoder(conn).Encode("0")
			fmt.Println("Goodbye")
			return
		default:
			clearTerminal()
			fmt.Println("That is not a valid option.\n\n")
		}
	}

}

func fileRead(conn net.Conn) {
	var msg []byte
	//Receive message from server
	_ = gob.NewDecoder(conn).Decode(&msg)
	serverContents := messaging.ParseServerMessage(msg)
	clearTerminal()
	fmt.Println("==========================")
	fmt.Println("Server Contents")
	fmt.Println("==========================")
	fmt.Println(serverContents.Message)
	return
}

func fileDownload(conn net.Conn) {
	var userSaveName string
	fmt.Print("Enter the name of the file to be saved:")
	fmt.Scanln(&userSaveName)
	var userSaveExt string
	fmt.Print("Enter the name of the file to be saved:")
	fmt.Scanln(&userSaveExt)
	var userOp string
	fmt.Print("Enter the number of file to download:")
	fmt.Scanln(&userOp)
	gob.NewEncoder(conn).Encode(userOp)
	filePath := fmt.Sprintf("./download/%v.%v", userSaveName, userSaveExt)
	fmt.Printf("Creating file at: %v\n", filePath)
	fmt.Println("Downloading file...")
	messaging.RecvFile(filePath, conn)
	// file, err := os.Create(filePath)
	// if err != nil {
	// 	fmt.Println("Error creating file")
	// 	return
	// }

	// defer file.Close()

	// _, err = io.Copy(file, conn)
	// if err != nil {
	// 	fmt.Println("Error receiving file:", err)
	// 	return
	// }
	fmt.Println("File downloaded")
	return
}

func clearTerminal() {
	cmd := exec.Command("clear")
	cmd.Stdout = os.Stdout
	cmd.Run()
}
