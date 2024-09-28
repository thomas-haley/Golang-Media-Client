package menu

import (
	"encoding/gob"
	"fmt"
	"net"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"thaley/fileServerClient/messaging"
	"thaley/fileServerClient/util"
	"time"
)

func MediaRoomMainMenu(conn net.Conn) {
	for {
		drawMediaRoomMainMenu()

		//Get Input
		fmt.Print("Enter command: ")
		var cmd string
		fmt.Scanln(&cmd)

		switch cmd {
		case "1":
			gob.NewEncoder(conn).Encode("1")
			activeRooms := messaging.ReadServerJSONMessage(conn)
			fmt.Print("\n\n=========================\n")
			fmt.Println("Active Rooms")
			fmt.Println("=========================\n")
			fmt.Println(activeRooms.Message)
			fmt.Println("")
			break
		case "2":
			gob.NewEncoder(conn).Encode("2")
			handleCreateRoom(conn)
			break
		case "3":
			gob.NewEncoder(conn).Encode("3")
			handleJoinRoom(conn)
			break
		case "0":
			gob.NewEncoder(conn).Encode("0")
			return
		default:
			gob.NewEncoder(conn).Encode("0")
			return
		}
	}

}

func drawMediaRoomMainMenu() {
	fmt.Print("\n\n===========================\n")
	fmt.Println("Media Rooms")
	fmt.Println("===========================")
	fmt.Println("[1] - List Active Rooms")
	fmt.Println("[2] - Create Room")
	fmt.Println("[3] - Join Room")
	fmt.Println("[0] - Exit")
}

func handleCreateRoom(conn net.Conn) {
	var selectedMedia string
	var mediaInt int
	var userName string
	availableMedia := messaging.ReadServerJSONMessage(conn)
	var roomConfig util.RoomConfig
	validInput := false
	for !validInput {
		//Get media list from server and display
		fmt.Println("==========================")
		fmt.Println("Available Media")
		fmt.Println("==========================")
		fmt.Println(availableMedia.Message)

		fmt.Print("Enter the file you would like to play: ")
		fmt.Scanln(&selectedMedia)
		fmt.Print("Enter your username: ")
		fmt.Scanln(&userName)
		var err error
		mediaInt, err = strconv.Atoi(selectedMedia)
		if err != nil {
			fmt.Println("Not a valid media file")
		}
		validInput = true
	}
	roomConfig.MEDIA_FILE = mediaInt
	roomConfig.ROOM_USER = userName
	roomConfigJSON := messaging.BuildRoomConfigMsg(roomConfig)
	// msg := messaging.BuildMessage(roomConfigJSON)
	messaging.SendClientMessage(conn, roomConfigJSON)
	ensureMediaFileDownloaded(conn)
	fmt.Println("Media File Downloaded")
	// ensureMediaFileDownloaded(conn)
	RoomHostMenu(conn)
}

func handleJoinRoom(conn net.Conn) {
	fmt.Print("\n\n===========================\n")
	fmt.Println("Media Rooms")
	fmt.Println("===========================")

	activeRooms := messaging.ReadServerJSONMessage(conn)
	fmt.Println(activeRooms.Message)

	var userOp string
	fmt.Print("Enter the room # you wish to join: ")
	fmt.Scanln(&userOp)

	var userName string
	fmt.Print("Enter your username: ")
	fmt.Scanln(&userName)
	msg := messaging.BuildMessage(fmt.Sprintf("%v|%v", userOp, userName))
	messaging.SendClientMessage(conn, msg)

	roomInfo := messaging.ReceiveRoomData(conn)
	ensureMediaFileDownloaded(conn)

	for {

		if roomInfo.IS_RUNNING == "false" {
			fmt.Println("Room joined, awaiting start command")
			awaitStartMessage(conn, roomInfo)
		}

		fmt.Println("======================")
		fmt.Println("Room Menu")
		fmt.Println("======================")
		fmt.Println("[1] - Open Media")
		fmt.Println("[0] - Leave room")
		fmt.Println("")
		fmt.Print("Enter a command: ")
		var menuOp string
		fmt.Scanln(&menuOp)
		switch menuOp {
		case "1":
			gob.NewEncoder(conn).Encode("1")
			roomInfo = messaging.ReceiveRoomData(conn)
			fmt.Println("Starting media again for guests")
			fmt.Println(roomInfo)
			startMedia(roomInfo)
			break
		case "0":
		default:
			gob.NewEncoder(conn).Encode("0")
			messaging.SendClientMessage(conn, msg)
			return
		}
	}

}

func ensureMediaFileDownloaded(conn net.Conn) {
	requiredMediaMessage := messaging.ReadServerJSONMessage(conn)
	requiredMedia := requiredMediaMessage.Message
	filePath := fmt.Sprintf("./download/%v", requiredMedia)
	//Check if we have the exact file already
	fmt.Printf("Saving file to path: %v\n", filePath)
	_, err := os.Stat(filePath)
	//If error, file doesnt exist, fetch from server
	if err != nil {
		fmt.Println("Client does not have file, sending info to server")
		msg := messaging.BuildMessage("false")
		messaging.SendClientMessage(conn, msg)
		fmt.Println("Downloading file from server")
		messaging.RecvFile(filePath, conn)
		fmt.Println("File Downloaded")
		return
	}
	fmt.Println("Client has file")
	msg := messaging.BuildMessage("true")
	messaging.SendClientMessage(conn, msg)

}

func RoomGuestMenu(conn net.Conn) {

}

func RoomHostMenu(conn net.Conn) {
	//Get watch room config and load into object

	for {
		fmt.Println("1:")
		roomConfig := messaging.ReceiveRoomData(conn)
		fmt.Println("2:")
		drawRoomHostMenu(roomConfig)
		fmt.Print("Enter command: ")
		var userCom string
		fmt.Scanln(&userCom)

		switch userCom {
		case "1":
			gob.NewEncoder(conn).Encode("1")
			userList := messaging.ReadServerJSONMessage(conn)
			fmt.Println("====================")
			fmt.Println("Connected Users")
			fmt.Println("====================")
			fmt.Println(userList.Message)
			break
		case "2":
			gob.NewEncoder(conn).Encode("2")
			if roomConfig.IS_RUNNING == "false" {
				fmt.Println("Awaiting start message")
				awaitStartMessage(conn, roomConfig)
			} else {
				startMedia(roomConfig)
			}
			break
		case "0":
			gob.NewEncoder(conn).Encode("0")
			return
		}
	}
}

func drawRoomHostMenu(watchRoom util.WatchRoom) {
	fmt.Println("====================")
	fmt.Println("Host Menu")
	fmt.Printf("Media: %v | Running: %t\n", watchRoom.ROOM_MEDIA, watchRoom.IS_RUNNING)
	fmt.Println("====================")
	fmt.Println("[1] - List Room Members")
	if watchRoom.IS_RUNNING == "true" {
		fmt.Println("[2] - Launch Media")

	} else {
		fmt.Println("[2] - Start Watch Countdown")
	}
	fmt.Println("[0] - Close Room")
}

func awaitStartMessage(conn net.Conn, room util.WatchRoom) {

	cdMsg := messaging.ReadServerJSONMessage(conn)

	if cdMsg.Message == "begin_countdown" {
		ch := make(chan bool)
		go func() {
			time.Sleep(3 * time.Second)
			ch <- true
		}()
		fmt.Println("Beginning watch countdown")
		wd, _ := os.Getwd()
		// strings.Replace(wd, "\\", "/", -1)
		path := fmt.Sprintf(`%v\download\%v`, wd, room.ROOM_MEDIA)
		strings.Replace(path, `\\`, `\`, -1)
		cmd := exec.Command(`vlc`, path)
		fmt.Println(cmd.String())
		_ = <-ch
		go cmd.Run()
	}

}

func startMedia(room util.WatchRoom) {
	layout := "2006-01-02 15:04:05 -0700 MST"
	serverStartTime, err := time.Parse(layout, room.MEDIA_START_TIME)
	if err != nil {
		fmt.Println("Time parse error")
		fmt.Println(err)
	}
	tNow := time.Now().UTC()

	dif := tNow.Sub(serverStartTime).Abs()

	difSec := dif.Seconds()
	wd, _ := os.Getwd()
	// strings.Replace(wd, "\\", "/", -1)
	path := fmt.Sprintf(`%v\download\%v`, wd, room.ROOM_MEDIA)
	strings.Replace(path, `\\`, `\`, -1)
	cmd := exec.Command(`vlc`, path, fmt.Sprintf("--start-time=%f", difSec))
	go cmd.Run()
}
