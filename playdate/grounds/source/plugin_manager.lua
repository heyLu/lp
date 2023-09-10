-- based on https://devforum.play.date/t/splitting-a-game-into-several-functional-binaries-nics-plugin-manager/1387
-- with some adjustments to make the linter happy

PluginManager = {}

-- private member
local _plugins = {}

function PluginManager.load( name )
	local path = name..'/'..name
	local filepath = path..'.pdz'

	if not playdate.file.exists( filepath ) then
		print( 'plugin '..name..' does not exist')
		return
	end

	_plugins[ name ] = {
		path = path,
		filepath = filepath,
		modtime = playdate.file.modtime( filepath ),
		update_fn = playdate.file.run( path )
	}
end

function PluginManager.update()
	for name, plugin in pairs(_plugins) do
		local modtime = playdate.file.modtime( plugin.filepath )

		-- check if we need to reload
		if not (
			modtime.year==plugin.modtime.year and
			modtime.month==plugin.modtime.month and
			modtime.day==plugin.modtime.day and
			modtime.hour==plugin.modtime.hour and
			modtime.minute==plugin.modtime.minute and
			modtime.second==plugin.modtime.second
			) then
			plugin.modtime = modtime
			print( 'Plugin reload: '..name)
			plugin.update_fn = playdate.file.run( plugin.path )
		end

		-- run plugin update function
		if type(plugin.update_fn)=="function" then
			plugin.update_fn()
		end
	end
end