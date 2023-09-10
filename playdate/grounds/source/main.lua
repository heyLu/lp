import 'CoreLibs/graphics'
import 'CoreLibs/sprites'
import 'CoreLibs/timer'
import 'CoreLibs/crank'

import 'plugin_manager'

PluginManager.load("dithering")

function playdate.update()
	playdate.graphics.clear()

	PluginManager.update()
end
