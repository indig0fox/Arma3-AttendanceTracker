params [
	["_eventType", ""],
	["_playerUID", ""],
	["_profileName", ""],
	["_steamName", ""]
];


private _hash = + (AttendanceTracker getVariable ["missionContext", createHashMap]);
_hash set ["eventType", _eventType];
_hash set ["playerUID", _playerUID];
_hash set ["profileName", _profileName];
_hash set ["steamName", _steamName];
_hash set ["isJIP", false];
_hash set ["roleDescription", ""];

"AttendanceTracker" callExtension ["logAttendance", [[_hash] call CBA_fnc_encodeJSON]];

true;