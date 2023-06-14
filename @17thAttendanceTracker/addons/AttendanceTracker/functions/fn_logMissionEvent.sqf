params [
	["_eventType", ""],
	["_playerId", ""],
	["_playerUID", ""],
	["_profileName", ""],
	["_steamName", ""],
	["_isJIP", false, [true, false]],
	["_roleDescription", ""],
	["_rowID", nil]
];

private _hash = + (AttendanceTracker getVariable ["missionContext", createHashMap]);
_hash set ["eventType", _eventType];
_hash set ["playerId", _playerId];
_hash set ["playerUID", _playerUID];
_hash set ["profileName", _profileName];
_hash set ["steamName", _steamName];
_hash set ["isJIP", _isJIP];
_hash set ["roleDescription", _roleDescription];
_hash set ["missionHash", missionNamespace getVariable ["AttendanceTracker_missionHash", ""]];

if (!isNil "_rowID") then {
	_hash set ["rowID", _rowID];
	"AttendanceTracker" callExtension ["writeDisconnectEvent", [[_hash] call CBA_fnc_encodeJSON]];
} else {
	"AttendanceTracker" callExtension ["writeAttendance", [[_hash] call CBA_fnc_encodeJSON]];
};

true;