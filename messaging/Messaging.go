package messaging

import (
	"bufio"
	"encoding/binary"
	"encoding/gob"
	"encoding/json"
	"fmt"
	"io"
	"net"
	"os"
	"strconv"
	"thaley/fileServerClient/util"
)

func ReadServerJSONMessage(conn net.Conn) util.JsonMessage {
	var msg []byte
	//Receive message from server
	_ = gob.NewDecoder(conn).Decode(&msg)
	msgObj := ParseServerMessage(msg)
	return msgObj

}

func ParseServerMessage(msgByte []byte) util.JsonMessage {
	var msgObj util.JsonMessage
	json.Unmarshal(msgByte, &msgObj)
	return msgObj
}

func RecvFile(name string, conn net.Conn) error {
	f, err := os.Create(name)
	if err != nil {
		fmt.Printf("Error creating file: %v\n", err)
		return err
	}
	defer f.Close()

	buf := bufio.NewReader(conn)
	var sz int64
	err = binary.Read(buf, binary.BigEndian, &sz)
	if err != nil {
		fmt.Printf("Error reading connection buffer: %v\n", err)
		return err
	}

	_, err = io.CopyN(f, buf, sz)
	if err != nil {
		fmt.Printf("Error copying data into file: %v\n", err)
		return err
	}

	return nil
}

func BuildMessage(message string) []byte {
	msgString := fmt.Sprintf(`{"Status":"good", "Message": "%v"}`, message)
	fmt.Println("msgString:")
	fmt.Println(msgString)
	var msgObj util.JsonMessage

	json.Unmarshal([]byte(msgString), &msgObj)
	fmt.Println("Built JSON message")
	fmt.Println(msgObj)
	mshMsg, _ := json.Marshal(msgObj)
	return mshMsg
}

func BuildRoomConfigMsg(roomConfig util.RoomConfig) []byte {
	var msgObj util.JsonMessage
	msgObj.Status = "good"
	msgObj.Message = ""
	config := map[string]string{"MEDIA_FILE": "0", "ROOM_USER": ""}
	fileString := strconv.Itoa(roomConfig.MEDIA_FILE)
	config["MEDIA_FILE"] = fileString
	config["ROOM_USER"] = roomConfig.ROOM_USER
	msgObj.Config = config

	mshMsh, _ := json.Marshal(msgObj)
	return mshMsh
}

func SendClientMessage(conn net.Conn, msgContent []byte) {

	err := gob.NewEncoder(conn).Encode(msgContent)
	if err != nil {
		fmt.Println(err)
	}
}

func ReceiveRoomData(conn net.Conn) util.WatchRoom {
	msg := ReadServerJSONMessage(conn)
	mshConfig, _ := json.Marshal(msg.Config)
	var watchRoom util.WatchRoom
	json.Unmarshal(mshConfig, &watchRoom)
	fmt.Println("Watch Room info")
	fmt.Println(watchRoom)
	return watchRoom
}
