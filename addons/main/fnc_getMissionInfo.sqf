#include "script_component.hpp"
[
	["missionName", missionName],
	["missionStart", GVAR(missionStart)],
	["missionHash", GVAR(missionHash)],
	["briefingName", briefingName],
	["missionNameSource", missionNameSource],
	["onLoadName", getMissionConfigValue ["onLoadName", "Unknown"]],
	["author", getMissionConfigValue ["author", "Unknown"]],
	["serverName", serverName],
	["serverProfile", profileName],
	["worldName", toLower worldName]
];