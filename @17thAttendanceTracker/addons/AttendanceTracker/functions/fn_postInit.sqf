
AttendanceTracker = false call CBA_fnc_createNamespace;

AttendanceTracker setVariable ["missionContext", createHashMapFromArray [
	["missionName", missionName],
	["briefingName", briefingName],
	["missionNameSource", missionNameSource],
	["onLoadName", getMissionConfigValue ["onLoadName", ""]],
	["author", getMissionConfigValue ["author", ""]],
	["serverName", serverName],
	["serverProfile", profileName],
	["missionStart", call attendanceTracker_fnc_timestamp]
]];

// store all user details in a hash when they connect so we can reference it in disconnect events
AttendanceTracker setVariable ["allUsers", createHashMap];

private _database = "AttendanceTracker" callExtension "connectDB";
systemChat "AttendanceTracker: Connecting to database...";
["Connecting to database...", "INFO"] call attendanceTracker_fnc_log;

{
	if (!isServer) exitWith {};
	_x params ["_ehName", "_code"];

	_handle = (addMissionEventHandler [_ehName, _code]);
    if (isNil "_handle") then {
        [format["Failed to add Mission event handler: %1", _x], "ERROR"] call attendanceTracker_fnc_log;
		false;
    } else {
        missionNamespace setVariable [
            ("AttendanceTracker" + "_MEH_" + _ehName),
            _handle
        ];
        true;
    };
} forEach (call attendanceTracker_fnc_eventHandlers);