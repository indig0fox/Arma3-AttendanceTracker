#include "script_component.hpp"

_world = ( configfile >> "CfgWorlds" >> worldName );
_author = getText( _world >> "author" );
_name = getText ( _world >> "description" );

_source = configSourceMod ( _world );

_workshopID = '';
{
	if ( ( _x#1 ) == _source ) then	{
		_workshopID = _x#7;
		break;
	};
} foreach getLoadedModsInfo;

if (_workshopID isEqualTo "") then {
	_workshopID = "0";
};

// [_name, _author, _workshopID];
_return = [
	["author", _author],
	["workshopID", _workshopID],
	["displayName", _name],
	["worldName", toLower worldName],
	["worldNameOriginal", _name],
	["worldSize", worldSize],
	["latitude", -1 * getNumber( _world >> "latitude" )],
	["longitude", getNumber( _world >> "longitude" )]
];
["DEBUG", format["WorldInfo is: %1", _return]] call FUNC(log);
_return
