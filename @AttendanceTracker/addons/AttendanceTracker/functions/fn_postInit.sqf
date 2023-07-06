
AttendanceTracker = false call CBA_fnc_createNamespace;

AttendanceTracker_missionStartTimestamp = call attendanceTracker_fnc_timestamp;
diag_log format ["AttendanceTracker: Mission started at %1", AttendanceTracker_missionStartTimestamp];
AttendanceTracker_missionHash = call attendanceTracker_fnc_getMissionHash;
diag_log format ["AttendanceTracker: Mission hash is %1", AttendanceTracker_missionHash];

_settings = call attendanceTracker_fnc_getSettings;
if (count _settings > 0) then {
	for "_i" from 0 to (count _settings) - 1 do {
		_setting = _settings select _i;
		_key = _setting select 0;
		_value = _setting select 1;
		missionNamespace setVariable ["AttendanceTracker_" + _key, _value];
	};
} else {
	[format["Failed to parse settings: %1", _settings], "ERROR"] call attendanceTracker_fnc_log;
};
call attendanceTracker_fnc_connectDB;

AttendanceTracker setVariable ["missionContext", createHashMapFromArray [
	["missionHash", AttendanceTracker_missionHash],
	["missionStart", AttendanceTracker_missionStartTimestamp],
	["missionName", missionName],
	["briefingName", briefingName],
	["missionNameSource", missionNameSource],
	["onLoadName", getMissionConfigValue ["onLoadName", ""]],
	["author", getMissionConfigValue ["author", ""]],
	["serverName", serverName],
	["serverProfile", profileName],
	["missionStart", AttendanceTracker_missionStartTimestamp],
	["missionHash", AttendanceTracker_missionHash],
	["worldName", toLower worldName]
]];



// store all user details in a hash when they connect so we can reference it in disconnect events
AttendanceTracker setVariable ["allUsers", createHashMap];
AttendanceTracker setVariable ["rowIds", createHashMap];

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