-- based on https://devforum.play.date/t/splitting-a-game-into-several-functional-binaries-nics-plugin-manager/1387
-- with some adjustments to make the linter happy

PluginManager = {
	plugins = {},

	current_plugin = "",
}

function PluginManager.load(self, name)
	local path = name..'/'..name
	local filepath = path..'.pdz'

	if not playdate.file.exists( filepath ) then
		print( 'plugin '..name..' does not exist')
		return
	end

	self.plugins[ name ] = {
		name = name,
		path = path,
		filepath = filepath,
		modtime = playdate.file.modtime( filepath ),
		update_fn = playdate.file.run( path )
	}
	self.current_plugin = name
end

function PluginManager.use(self, name)
	self.current_plugin = name
end

function PluginManager.update(self)
	for name, plugin in pairs(self.plugins) do
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
	end

	local plugin = self.plugins[self.current_plugin]
	if plugin == nil then
		playdate.graphics.drawText("no plugin "..self.current_plugin, 0, 0)
		return
	end
	plugin.update_fn()
end