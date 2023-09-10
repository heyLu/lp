local gfx <const> = playdate.graphics

local img = gfx.image.new(16, 16)
local width, height = img:getSize()
gfx.pushContext(img)
gfx.fillRect(0, 0, width, height)
gfx.popContext(img)

local function update()
	local offset = width + 2
	local numSteps = (240-30) / offset

	local ditherTypes = {
		{dither=gfx.image.kDitherTypeNone, name="n"},
		{dither=gfx.image.kDitherTypeDiagonalLine, name = "d"},
		{dither=gfx.image.kDitherTypeHorizontalLine, name="h"},
		{dither=gfx.image.kDitherTypeScreen, name="s"},
		{dither=gfx.image.kDitherTypeBayer2x2, name="b2"},
		{dither=gfx.image.kDitherTypeBayer4x4, name="b4"},
		{dither=gfx.image.kDitherTypeBayer8x8, name="b8"},
		{dither=gfx.image.kDitherTypeFloydSteinberg, name="fs"},
		{dither=gfx.image.kDitherTypeBurkes, name="bu"},
		{dither=gfx.image.kDitherTypeAtkinson, name="at"},
	}

	for t, d in ipairs(ditherTypes) do
		for i = 1, numSteps, 1 do
			gfx.drawText(d.name, 50+(offset+5)*t, 10)

			local alpha = math.min(math.floor((numSteps-(i-1))/(numSteps)*1000) / 1000, 1.0)
			gfx.drawText(tostring(alpha), 5, 30+offset*(i-1))
			img:fadedImage(alpha, d.dither):draw(50+(offset+5)*t, 30 + offset*(i-1))
		end
	end
end

return update