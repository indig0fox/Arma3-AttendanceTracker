createHashMapFromArray [
	["missionName", missionName],
	["missionStart", ATNamespace getVariable "missionStartTime"],
	["missionHash", ATNamespace getVariable "missionHash"],
	["briefingName", briefingName],
	["missionNameSource", missionNameSource],
	["onLoadName", getMissionConfigValue ["onLoadName", ""]],
	["author", getMissionConfigValue ["author", ""]],
	["serverName", serverName],
	["serverProfile", profileName],
	["worldName", toLower worldName]
];