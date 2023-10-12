#include "script_component.hpp"

if (!isServer) exitWith {};

if (typeName _this != "ARRAY") exitWith {
    diag_log format ["[%1]: Invalid log params: %2", GVAR(logPrefix), _this];
};

params [
    ["_level", "INFO", [""]],
    ["_text", "", ["", []]]
];

if (_text isEqualType []) then {
    _text = format ["%1", _text];
};

if (
    _level == "DEBUG" && 
    !GVAR(debug)
) exitWith {};

if (_text isEqualTo "") exitWith {};

diag_log formatText [
    "[%1] %2: %3",
    GVAR(logPrefix),
    _level,
    _text
];

