params [
	["_message", "", [""]],
	["_level", "INFO", [""]],
	"_function"
];

if (isNil "_message") exitWith {false};
if (
	missionNamespace getVariable ["AttendanceTracker_debug", false] &&
	_level == "DEBUG"
) exitWith {};

"AttendanceTracker" callExtension ["log", [_level, _message]];

if (!isNil "_function") then {
	diag_log formatText["[AttendanceTracker] (%1): <%2> %3", _level, _function, _message];
} else {
	diag_log formatText["[AttendanceTracker] (%1): %2", _level, _message];
};

true;