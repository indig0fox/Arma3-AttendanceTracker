#include "..\script_mod.hpp"

params [
	["_message", "", [""]],
	["_level", "INFO", [""]],
	"_function"
];

if (isNil "_message") exitWith {false};
if (
	missionNamespace getVariable ["ATDebug", true] &&
	_level != "WARN" && _level != "ERROR"
) exitWith {};

LOG_SYS(_level, _message);

true;