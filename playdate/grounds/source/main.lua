import 'CoreLibs/graphics'
import 'CoreLibs/sprites'
import 'CoreLibs/timer'
import 'CoreLibs/crank'

import 'plugin_manager'

PluginManager:load("grid_view")
PluginManager:load("dithering")
PluginManager:load("shady")

local pluginNames = {}
local i = 1
for name, _ in pairs(PluginManager.plugins) do
	pluginNames[i] = name
	i = i+1
end
playdate.getSystemMenu():addOptionsMenuItem("demo", pluginNames, function(pluginName)
	PluginManager:use(pluginName)
end)

function playdate.update()
	playdate.graphics.clear()

	PluginManager:update()
end
