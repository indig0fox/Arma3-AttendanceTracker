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

// [_name, _author, _workshopID];

[
	["worldName", _name],
	["author", _author],
	["worldSize", worldSize],
	["workshopID", _workshopID]
];