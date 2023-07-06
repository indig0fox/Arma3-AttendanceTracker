params [
	["_eventType", ""],
	["_playerId", ""],
	["_playerUID", ""],
	["_profileName", ""],
	["_steamName", ""],
	["_isJIP", false, [true, false]],
	["_roleDescription", ""]
];

private _hash = + (AttendanceTracker getVariable ["missionContext", createHashMap]);

_hash set ["eventType", _eventType];
_hash set ["playerId", _playerId];
_hash set ["playerUID", _playerUID];
_hash set ["profileName", _profileName];
_hash set ["steamName", _steamName];
_hash set ["isJIP", _isJIP];
_hash set ["roleDescription", _roleDescription];

[
	{missionNamespace getVariable ["AttendanceTracker_DBConnected", false]},
	{"AttendanceTracker" callExtension ["writeAttendance", [[_this] call CBA_fnc_encodeJSON]]},
	_hash, // args
	30 // timeout in seconds. if DB never connects, we don't want these building up
] call CBA_fnc_waitUntilAndExecute;

true;