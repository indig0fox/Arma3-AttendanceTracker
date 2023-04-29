package main

/*
#include <stdlib.h>
#include <stdio.h>
#include <string.h>
#include "extensionCallback.h"
*/
import "C" // This is required to import the C code

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
	"time"
	"unsafe"

	_ "github.com/go-sql-driver/mysql"
)

var EXTENSION_VERSION string = "0.0.1"
var extensionCallbackFnc C.extensionCallback

// file paths
var ADDON_FOLDER string = getDir() + "\\@17thAttendanceTracker"
var LOG_FILE string = ADDON_FOLDER + "\\attendanceTracker.log"
var CONFIG_FILE string = ADDON_FOLDER + "\\config.json"

var ATConfig AttendanceTrackerConfig

type AttendanceTrackerConfig struct {
	MySQLHost     string `json:"mysqlHost"`
	MySQLPort     int    `json:"mysqlPort"`
	MySQLUser     string `json:"mysqlUser"`
	MySQLPassword string `json:"mysqlPassword"`
	MySQLDatabase string `json:"mysqlDatabase"`
}

// database connection
var db *sql.DB

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
		log.Fatal(err)
	}
	return dir
}

func loadConfig() {
	// load config from file as JSON
	functionName := "loadConfig"
	file, err := os.Open(CONFIG_FILE)
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

func connectDB() string {
	functionName := "connectDB"
	var err error

	// load config
	loadConfig()

	// connect to database
	connectionString := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s", ATConfig.MySQLUser, ATConfig.MySQLPassword, ATConfig.MySQLHost, ATConfig.MySQLPort, ATConfig.MySQLDatabase)

	db, err = sql.Open("mysql", connectionString)
	if err != nil {
		writeLog(functionName, fmt.Sprintf(`["%s", "ERROR"]`, err))
		return "ERROR"
	}
	if db == nil {
		writeLog(functionName, fmt.Sprintf(`["%s", "ERROR"]`, "db is nil"))
		return "ERROR"
	}
	// defer db.Close()

	db.SetConnMaxLifetime(time.Minute * 3)
	db.SetMaxOpenConns(10)
	db.SetMaxIdleConns(10)

	pingErr := db.Ping()
	if pingErr != nil {
		writeLog(functionName, fmt.Sprintf(`["%s", "ERROR"]`, pingErr))
		return "ERROR"
	}

	// Connect and check the server version
	var version string
	err = db.QueryRow("SELECT VERSION()").Scan(&version)
	if err != nil {
		writeLog(functionName, fmt.Sprintf(`["%s", "ERROR"]`, err))
		return "ERROR"
	}
	writeLog(functionName, fmt.Sprintf(`["Connected to MySQL/MariaDB version %s", "INFO"]`, version))
	return version
}

type AttendanceLogItem struct {
	MissionName       string `json:"missionName"`
	BriefingName      string `json:"briefingName"`
	MissionNameSource string `json:"missionNameSource"`
	OnLoadName        string `json:"onLoadName"`
	Author            string `json:"author"`
	ServerName        string `json:"serverName"`
	ServerProfile     string `json:"serverProfile"`
	MissionStart      string `json:"missionStart"`
	// situational
	EventType       string `json:"eventType"`
	PlayerId        string `json:"playerId"`
	PlayerUID       string `json:"playerUID"`
	ProfileName     string `json:"profileName"`
	SteamName       string `json:"steamName"`
	IsJIP           bool   `json:"isJIP"`
	RoleDescription string `json:"roleDescription"`
}

func writeAttendance(data string) {
	functionName := "writeAttendance"
	// data is json, parse it
	stringjson := fixEscapeQuotes(trimQuotes(data))
	var event AttendanceLogItem
	err := json.Unmarshal([]byte(stringjson), &event)
	if err != nil {
		writeLog(functionName, fmt.Sprintf(`["%s", "ERROR"]`, err))
		return
	}

	// get MySQL friendly datetime
	// first, convert string to int
	missionStartTime, err := strconv.ParseInt(event.MissionStart, 10, 64)
	if err != nil {
		writeLog(functionName, fmt.Sprintf(`["%s", "ERROR"]`, err))
		return
	}
	t := time.Unix(0, missionStartTime).Format("2006-01-02 15:04:05")
	// get MySQL friendly NOW
	now := time.Now().Format("2006-01-02 15:04:05")

	// prevent crash
	if db == nil {
		writeLog(functionName, `["db is nil", "ERROR"]`)
		return
	}

	// send to DB

	result, err := db.ExecContext(context.Background(), `INSERT INTO AttendanceLog (timestamp, mission_name, briefing_name, mission_name_source, on_load_name, author, server_name, server_profile, mission_start, event_type, player_id, player_uid, profile_name, steam_name, is_jip, role_description) VALUES (?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?)`,
		now,
		event.MissionName,
		event.BriefingName,
		event.MissionNameSource,
		event.OnLoadName,
		event.Author,
		event.ServerName,
		event.ServerProfile,
		t,
		event.EventType,
		event.PlayerId,
		event.PlayerUID,
		event.ProfileName,
		event.SteamName,
		event.IsJIP,
		event.RoleDescription,
	)

	if err != nil {
		writeLog(functionName, fmt.Sprintf(`["%s", "ERROR"]`, err))
		return
	}

	id, err := result.LastInsertId()
	if err != nil {
		writeLog(functionName, fmt.Sprintf(`["%s", "ERROR"]`, err))
		return
	}

	writeLog(functionName, fmt.Sprintf(`["Saved attendance for %s to row id %d", "INFO"]`, event.ProfileName, id))

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
	case "logAttendance": // callExtension ["serverEvent", [_hash] call CBA_fnc_encodeJSON];
		if argc == 1 {
			go writeAttendance(out[0])
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

func getTimestamp() int64 {
	// get the current unix timestamp in nanoseconds
	return time.Now().UnixNano()
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
		time := getTimestamp()
		temp = fmt.Sprintf(`["%s"]`, strconv.FormatInt(time, 10))
	case "connectDB":
		go connectDB()
		temp = fmt.Sprintf(`["%s"]`, "Connecting to DB")

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
