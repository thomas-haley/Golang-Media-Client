package util

type WatchRoom struct {
	ROOM_OWNER_NAME  string `json:"ROOM_OWNER_NAME"`
	ROOM_MEDIA       string `json:"ROOM_MEDIA"`
	IS_RUNNING       string `json:"IS_RUNNING"`
	MEDIA_START_TIME string `json:"MEDIA_START_TIME"`
}

// func BuildWatchRoomFromMessage(config map[string]string) WatchRoom{
// 	var newRoom WatchRoom
// 	newRoom.ROOM_OWNER_NAME =
// }
