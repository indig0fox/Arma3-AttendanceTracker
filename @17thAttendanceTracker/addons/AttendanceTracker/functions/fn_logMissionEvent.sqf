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

"AttendanceTracker" callExtension ["logAttendance", [[_hash] call CBA_fnc_encodeJSON]];

true;