
AttendanceTracker = false call CBA_fnc_createNamespace;

AttendanceTracker_missionStartTimestamp = call attendanceTracker_fnc_timestamp;
AttendanceTracker_missionHash = "AttendanceTracker" callExtension ["getMissionHash", AttendanceTracker_missionStartTimestamp];

AttendanceTracker setVariable ["missionContext", createHashMapFromArray [
	["missionName", missionName],
	["briefingName", briefingName],
	["missionNameSource", missionNameSource],
	["onLoadName", getMissionConfigValue ["onLoadName", ""]],
	["author", getMissionConfigValue ["author", ""]],
	["serverName", serverName],
	["serverProfile", profileName],
	["missionStart", AttendanceTracker_missionStartTimestamp],
	["missionHash", AttendanceTracker_missionHash]
]];



// store all user details in a hash when they connect so we can reference it in disconnect events
AttendanceTracker setVariable ["allUsers", createHashMap];
missionNamespace setVariable ["AttendanceTracker_debug", false];

call attendanceTracker_fnc_connectDB;

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