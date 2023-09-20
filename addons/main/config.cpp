#include "script_mod.hpp"

class CfgPatches {
	class AttendanceTracker {
		units[] = {};
		weapons[] = {};
		requiredVersion = 2.10;
		requiredAddons[] = {
			"cba_main",
			"cba_xeh",
			"cba_settings"
		};
		VERSION_CONFIG;
		author[] = {"IndigoFox"};
		authorUrl = "https://github.com/indig0fox";
	};
};

class CfgFunctions {
	class attendanceTracker {
		class functions {
			file = "x\addons\attendancetracker\main\functions";
			class postInit {postInit = 1;};
			class callbackHandler {postInit = 1;};
			class getMissionHash {};
			class getMissionInfo {};
			class getSettings {};
			class getWorldInfo {};
			class log {};
			class missionLoaded {};
			class onPlayerConnected {};
			class timestamp {};
			class writePlayer {};
		};
	};
};