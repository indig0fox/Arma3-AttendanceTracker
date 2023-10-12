#include "script_component.hpp"

if (!isServer) exitWith {};

GVAR(attendanceTracker) = true;
GVAR(debug) = true;
GVAR(logPrefix) = "AttendanceTracker";
GVAR(extensionName) = "AttendanceTracker";
GVAR(missionLogged) = false;



addMissionEventHandler ["ExtensionCallback", {
	params ["_name", "_function", "_data"];
	if !(_name isEqualTo GVAR(extensionName)) exitWith {};
	
	_dataArr = parseSimpleArray _data;
	if (count _dataArr isEqualTo 0) exitWith {};

	switch (_function) do {
		case ":LOG:MISSION:SUCCESS:": {
			GVAR(missionLogged) = true;
		};
		case ":LOG:": {
			diag_log formatText[
				"[%1] %2",
				GVAR(logPrefix),
				_dataArr select 0
			];
		};
		default {
			["DEBUG", format["%1", _dataArr]] call FUNC(log);
		};
	};
}];


// LOAD EXTENSION
GVAR(extensionName) callExtension ":START:";

// GET MISSION START TIMESTAMP AND UNIQUE HASH
private _missionHashData = parseSimpleArray ("AttendanceTracker" callExtension ":MISSION:HASH:");
if (count _missionHashData isEqualTo 0) exitWith {
	["ERROR", "Failed to get mission hash, exiting"] call FUNC(log);
};

_missionHashData params ["_timestamp", "_hash"];
GVAR(missionStart) = _timestamp;
GVAR(missionHash) = _hash;


// PARSE SETTINGS
private _settings = parseSimpleArray (GVAR(extensionName) callExtension ":GET:SETTINGS:");
if (count _settings isEqualTo 0) exitWith {
	["ERROR", "Failed to get settings, exiting"] call FUNC(log);
};

GVAR(settings) = createHashMapFromArray (_settings#0);
GVAR(debug) = GVAR(settings) getOrDefault ["debug", GVAR(debug)];
private _updateInterval = GVAR(settings) getOrDefault ["dbupdateinterval", 90];
// remove duration by removing the last index
_updateInterval = _updateInterval select [0, count _updateInterval - 1];
GVAR(updateInterval) = parseNumber _updateInterval;

// add player connected (to mission) handler
addMissionEventHandler ["PlayerConnected", {
	_this call FUNC(onPlayerConnected);
}];


// we'll wait for the end of init (DB connect included) of the extension
// then we'll log the world and mission
// the response to THAT is handled above in the extension callback
// and will set GVAR(missionLogged) true
addMissionEventHandler ["ExtensionCallback", {
	params ["_name", "_function", "_data"];
	if !(_name isEqualTo GVAR(extensionName)) exitWith {};
	if !(_function isEqualTo ":READY:") exitWith {};

	// LOAD WORLD AND MISSION INFO
	GVAR(worldInfo) = call FUNC(getWorldInfo);
	GVAR(missionInfo) = call FUNC(getMissionInfo);

	["INFO", (GVAR(extensionName) callExtension [
		":LOG:MISSION:",
		[
			GVAR(worldInfo),
			GVAR(missionInfo)
		]
	]) select 0] call FUNC(log);

	// remove the handler
	removeMissionEventHandler ["ExtensionCallback", _thisEventHandler];
}];