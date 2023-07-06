package main

/*
#include <stdlib.h>
#include <stdio.h>
#include <string.h>
#include "extensionCallback.h"
*/
import "C" // This is required to import the C code

import (
	"crypto/md5"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"path"
	"runtime"
	"strings"
	"time"
	"unsafe"

	_ "github.com/go-sql-driver/mysql"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var EXTENSION_VERSION string = "0.0.1"
var extensionCallbackFnc C.extensionCallback

// file paths
var ADDON_FOLDER string = getDir() + "\\@AttendanceTracker"
var LOG_FILE string = ADDON_FOLDER + "\\attendanceTracker.log"
var CONFIG_FILE string = ADDON_FOLDER + "\\config.json"

var ATTENDANCE_TABLE string = "attendance"
var MISSIONS_TABLE string = "missions"
var WORLDS_TABLE string = "worlds"

// ! TODO make a hash to save key:netId from A3 value:rowId from join event

var ATConfig AttendanceTrackerConfig

type AttendanceTrackerConfig struct {
	MySQLHost     string `json:"mysqlHost"`
	MySQLPort     int    `json:"mysqlPort"`
	MySQLUser     string `json:"mysqlUser"`
	MySQLPassword string `json:"mysqlPassword"`
	MySQLDatabase string `json:"mysqlDatabase"`
}

// database connection
var db *gorm.DB

// configure log output
func init() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)
	// log to file
	f, err := os.OpenFile(LOG_FILE, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		log.Fatalf("error opening file: %v", err)
	}
	// log to console as well
	// log.SetOutput(io.MultiWriter(f, os.Stdout))
	// log only to file
	log.SetOutput(f)
}

func version() {
	functionName := "version"
	writeLog(functionName, fmt.Sprintf(`["AttendanceTracker Extension Version:%s", "INFO"]`, EXTENSION_VERSION))
}

func getDir() string {
	dir, err := os.Getwd()
	if err != nil {
		writeLog("getDir", fmt.Sprintf(`["Error getting working directory: %v", "ERROR"]`, err))
		return ""
	}
	return dir
}

func loadConfig() {
	// load config from file as JSON
	functionName := "loadConfig"

	// get location of this dll
	// dllPath, err := filepath.Abs(os.Args[0])
	// if err != nil {
	// 	writeLog(functionName, fmt.Sprintf(`["Error getting DLL path: %v", "ERROR"]`, err))
	// 	return
	// }

	// set the addon directory to the parent directory of the dll
	// ADDON_FOLDER = filepath.Dir(dllPath)
	// LOG_FILE = ADDON_FOLDER + "\\attendanceTracker.log"
	// CONFIG_FILE = ADDON_FOLDER + "\\config.json"

	file, err := os.OpenFile(CONFIG_FILE, os.O_RDONLY|os.O_CREATE, 0666)
	if err != nil {
		writeLog(functionName, fmt.Sprintf(`["%s", "ERROR"]`, err))
		return
	}
	defer file.Close()

	decoder := json.NewDecoder(file)
	err = decoder.Decode(&ATConfig)
	if err != nil {
		writeLog(functionName, fmt.Sprintf(`["%s", "ERROR"]`, err))
		return
	}
	writeLog(functionName, `["Config loaded", "INFO"]`)
}

func getMissionHash() string {
	functionName := "getMissionHash"
	// get md5 hash of string
	// https://stackoverflow.com/questions/2377881/how-to-get-a-md5-hash-from-a-string-in-golang
	hash := md5.Sum([]byte(time.Now().UTC().Format("2006-01-02 15:04:05")))

	// convert to string
	hashString := fmt.Sprintf(`%x`, hash)
	writeLog(functionName, fmt.Sprintf(`["Mission hash: %s", "INFO"]`, hashString))
	return hashString
}

func connectDB() error {
	var err error
	dsn := fmt.Sprintf(
		"%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		ATConfig.MySQLUser,
		ATConfig.MySQLPassword,
		ATConfig.MySQLHost,
		ATConfig.MySQLPort,
		ATConfig.MySQLDatabase,
	)

	// log dsn and pause
	writeLog("connectDB", fmt.Sprintf(`["DSN: %s", "INFO"]`, dsn))
	var input string
	fmt.Scanln(&input)

	db, err = gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		writeLog("connectDB", fmt.Sprintf(`["%s", "ERROR"]`, err))
		return err
	}

	// Migrate the schema
	err = db.AutoMigrate(&World{}, &Mission{}, &AttendanceItem{})
	if err != nil {
		writeLog("connectDB", fmt.Sprintf(`["%s", "ERROR"]`, err))
		return err
	}

	writeLog("connectDB", `["Database connected", "INFO"]`)
	return nil
}

type World struct {
	gorm.Model
	Author            string  `json:"author"`
	WorkshopID        string  `json:"workshopID"`
	DisplayName       string  `json:"displayName"`
	WorldName         string  `json:"worldName"`
	WorldNameOriginal string  `json:"worldNameOriginal"`
	WorldSize         int     `json:"worldSize"`
	Latitude          float32 `json:"latitude"`
	Longitude         float32 `json:"longitude"`
	Missions          []Mission
}

func writeWorldInfo(worldInfo string) {
	functionName := "writeWorldInfo"
	// writeLog(functionName, fmt.Sprintf(`["%s", "DEBUG"]`, worldInfo))
	// worldInfo is json, parse it
	var wi World
	fixedString := fixEscapeQuotes(trimQuotes(worldInfo))
	err := json.Unmarshal([]byte(fixedString), &wi)
	if err != nil {
		writeLog(functionName, fmt.Sprintf(`["%s", "ERROR"]`, err))
		return
	}

	// prevent crash
	if db == nil {
		connectDB()
	}

	// write world if not exist
	var world World
	db.Where("world_name = ?", wi.WorldName).First(&world)
	if world.ID == 0 {
		writeLog(functionName, `["World not found, writing new world", "INFO"]`)
		result := db.Create(&wi)
		if result.Error != nil {
			writeLog(functionName, fmt.Sprintf(`["%s", "ERROR"]`, result.Error))
			return
		}
		writeLog(functionName, fmt.Sprintf(`["World written with ID %d", "INFO"]`, wi.ID))
	} else {
		// return ID
		writeLog(functionName, fmt.Sprintf(`["World exists with ID %d", "INFO"]`, world.ID))
	}
}

type Mission struct {
	gorm.Model
	MissionName       string `json:"missionName"`
	BriefingName      string `json:"briefingName"`
	MissionNameSource string `json:"missionNameSource"`
	OnLoadName        string `json:"onLoadName"`
	Author            string `json:"author"`
	ServerName        string `json:"serverName"`
	ServerProfile     string `json:"serverProfile"`
	MissionStart      string `json:"missionStart"`
	MissionHash       string `json:"missionHash"`
	WorldName         string `json:"worldName" gorm:"-"`
	WorldID           int
	World             World `gorm:"foreignkey:WorldID"`
	Attendees         []AttendanceItem
}

func writeMission(missionJSON string) {
	functionName := "writeMission"
	var err error
	// writeLog(functionName, fmt.Sprintf(`["%s", "DEBUG"]`, Mission))
	// Mission is json, parse it
	var mi Mission
	fixedString := fixEscapeQuotes(trimQuotes(missionJSON))
	err = json.Unmarshal([]byte(fixedString), &mi)
	if err != nil {
		writeLog(functionName, fmt.Sprintf(`["%s", "ERROR"]`, err))
		return
	}

	// write mission to database
	db.Create(&mi)
	if db.Error != nil {
		writeLog(functionName, fmt.Sprintf(`["%s", "ERROR"]`, db.Error))
		return
	}
	writeLog(functionName, fmt.Sprintf(`["Mission written with ID %d", "INFO"]`, mi.ID))
}

type AttendanceItem struct {
	gorm.Model
	MissionHash     string `json:"missionHash"`
	EventType       string `json:"eventType"`
	PlayerId        string `json:"playerId"`
	PlayerUID       string `json:"playerUID"`
	JoinTime        time.Time
	DisconnectTime  time.Time
	ProfileName     string `json:"profileName"`
	SteamName       string `json:"steamName"`
	IsJIP           bool   `json:"isJIP"`
	RoleDescription string `json:"roleDescription"`
	MissionID       int
}

func writeAttendance(data string) {
	functionName := "writeAttendance"
	var err error
	// data is json, parse it
	stringjson := fixEscapeQuotes(trimQuotes(data))
	var event AttendanceItem
	err = json.Unmarshal([]byte(stringjson), &event)
	if err != nil {
		writeLog(functionName, fmt.Sprintf(`["%s", "ERROR"]`, err))
		return
	}

	// prevent crash
	if db == nil {
		connectDB()
	}

	var playerUid string
	var rowId uint
	if event.EventType == "Server" {
		// check for most recent existing attendance row
		var attendance AttendanceItem
		db.Where("player_uid = ? AND event_type = ?", event.PlayerId, "Server").Last(&attendance)
		if attendance.ID != 0 {
			// update disconnect time
			row := db.Model(&attendance).Update("disconnect_time", time.Now().UTC().Format("2006-01-02 15:04:05"))
			if row.Error != nil {
				writeLog(functionName, fmt.Sprintf(`["%s", "ERROR"]`, row.Error))
				return
			}
			rowId, playerUid = attendance.ID, attendance.PlayerUID

		} else {
			// insert new row
			row := db.Create(&event)
			if row.Error != nil {
				writeLog(functionName, fmt.Sprintf(`["%s", "ERROR"]`, row.Error))
				return
			}
			rowId, playerUid = event.ID, event.PlayerUID
		}
	} else if event.EventType == "Mission" {
		// check for most recent join_time for this player within 6 hours without a disconnect_time
		var attendance AttendanceItem
		db.Where("player_uid = ? AND join_time > ? AND disconnect_time IS NULL", event.PlayerUID, time.Now().UTC().Add(-6*time.Hour).Format("2006-01-02 15:04:05")).Last(&attendance)
		if attendance.ID != 0 {
			// update disconnect time
			row := db.Model(&attendance).Update("disconnect_time", time.Now().UTC().Format("2006-01-02 15:04:05"))
			if row.Error != nil {
				writeLog(functionName, fmt.Sprintf(`["%s", "ERROR"]`, row.Error))
				return
			}
			rowId, playerUid = attendance.ID, attendance.PlayerUID
		} else {
			// insert new row
			row := db.Create(&event)
			if row.Error != nil {
				writeLog(functionName, fmt.Sprintf(`["%s", "ERROR"]`, row.Error))
				return
			}
			rowId, playerUid = event.ID, event.PlayerUID
		}
	}

	writeLog(functionName, fmt.Sprintf(`["Saved attendance for %s to row id %d", "INFO"]`, playerUid, rowId))
	if event.EventType == "Server" {
		writeLog(functionName, fmt.Sprintf(`["ATT_LOG", ["SERVER", "%s", "%d"]]`, playerUid, rowId))
	} else if event.EventType == "Mission" {
		writeLog(functionName, fmt.Sprintf(`["ATT_LOG", ["MISSION", "%s", "%d"]]`, playerUid, rowId))
	}

}

func runExtensionCallback(name *C.char, function *C.char, data *C.char) C.int {
	return C.runExtensionCallback(extensionCallbackFnc, name, function, data)
}

//export goRVExtensionVersion
func goRVExtensionVersion(output *C.char, outputsize C.size_t) {
	result := C.CString(EXTENSION_VERSION)
	defer C.free(unsafe.Pointer(result))
	var size = C.strlen(result) + 1
	if size > outputsize {
		size = outputsize
	}
	C.memmove(unsafe.Pointer(output), unsafe.Pointer(result), size)
}

//export goRVExtensionArgs
func goRVExtensionArgs(output *C.char, outputsize C.size_t, input *C.char, argv **C.char, argc C.int) {
	var offset = unsafe.Sizeof(uintptr(0))
	var out []string
	for index := C.int(0); index < argc; index++ {
		out = append(out, C.GoString(*argv))
		argv = (**C.char)(unsafe.Pointer(uintptr(unsafe.Pointer(argv)) + offset))
	}

	// temp := fmt.Sprintf("Function: %s nb params: %d params: %s!", C.GoString(input), argc, out)
	temp := fmt.Sprintf("Function: %s nb params: %d", C.GoString(input), argc)

	switch C.GoString(input) {
	case "fillLastMissionNull":
		{
			// go fillLastMissionNull()
		}
	case "writeAttendance":
		{ // callExtension ["logAttendance", [_hash] call CBA_fnc_encodeJSON]];
			if argc == 1 {
				go writeAttendance(out[0])
			}
		}
	case "writeDisconnectEvent":
		{ // callExtension ["writeDisconnectEvent", [[_hash] call CBA_fnc_encodeJSON]];

			if argc == 1 {
				// go writeDisconnectEvent(out[0])
				go writeAttendance(out[0])
			}
		}
	case "logMission":
		if argc == 1 {
			go writeMission(out[0])
		}
	case "logWorld":
		if argc == 1 {
			go writeWorldInfo(out[0])
		}
	}

	// Return a result to Arma
	result := C.CString(temp)
	defer C.free(unsafe.Pointer(result))
	var size = C.strlen(result) + 1
	if size > outputsize {
		size = outputsize
	}

	C.memmove(unsafe.Pointer(output), unsafe.Pointer(result), size)
}

func callBackExample() {
	name := C.CString("arma")
	defer C.free(unsafe.Pointer(name))
	function := C.CString("funcToExecute")
	defer C.free(unsafe.Pointer(function))
	// Make a callback to Arma
	for i := 0; i < 3; i++ {
		time.Sleep(2 * time.Second)
		param := C.CString(fmt.Sprintf("Loop: %d", i))
		defer C.free(unsafe.Pointer(param))
		runExtensionCallback(name, function, param)
	}
}

func getTimestamp() string {
	// get the current unix timestamp in nanoseconds
	// return time.Now().Local().Unix()
	return time.Now().UTC().Format("2006-01-02 15:04:05")
}

func trimQuotes(s string) string {
	// trim the start and end quotes from a string
	return strings.Trim(s, `"`)
}

func fixEscapeQuotes(s string) string {
	// fix the escape quotes in a string
	return strings.Replace(s, `""`, `"`, -1)
}

func writeLog(functionName string, data string) {
	statusName := C.CString("AttendanceTracker")
	defer C.free(unsafe.Pointer(statusName))
	statusFunction := C.CString(functionName)
	defer C.free(unsafe.Pointer(statusFunction))
	statusParam := C.CString(data)
	defer C.free(unsafe.Pointer(statusParam))
	runExtensionCallback(statusName, statusFunction, statusParam)

	// get calling function & line
	_, file, line, _ := runtime.Caller(1)
	log.Printf(`%s:%d: %s`, path.Base(file), line, data)
	log.Printf(`%s: %s`, functionName, data)
}

//export goRVExtension
func goRVExtension(output *C.char, outputsize C.size_t, input *C.char) {

	var temp string

	// logLine("goRVExtension", fmt.Sprintf(`["Input: %s",  "DEBUG"]`, C.GoString(input)), true)

	switch C.GoString(input) {
	case "version":
		temp = EXTENSION_VERSION
	case "getDir":
		temp = getDir()
	case "getTimestamp":
		temp = fmt.Sprintf(`["%s"]`, getTimestamp())
	case "connectDB":
		go connectDB()
		temp = fmt.Sprintf(`["%s"]`, "Connecting to DB")
	case "getMissionHash":
		temp = fmt.Sprintf(`["%s"]`, getMissionHash())
	default:
		temp = fmt.Sprintf(`["%s"]`, "Unknown Function")
	}

	result := C.CString(temp)
	defer C.free(unsafe.Pointer(result))
	var size = C.strlen(result) + 1
	if size > outputsize {
		size = outputsize
	}

	C.memmove(unsafe.Pointer(output), unsafe.Pointer(result), size)
	// return
}

//export goRVExtensionRegisterCallback
func goRVExtensionRegisterCallback(fnc C.extensionCallback) {
	extensionCallbackFnc = fnc
}

func main() {}
