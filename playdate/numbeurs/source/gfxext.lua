import "CoreLibs/graphics"

local gfx <const> = playdate.graphics

playdate.graphicsext = {
  -- from https://devforum.play.date/t/add-a-drawtextscaled-api-see-code-example/7108
  drawTextScaled = function(text, x, y, scale, font)
      local padding = string.upper(text) == text and 6 or 0 -- Weird padding hack?
      local w <const> = font:getTextWidth(text)
      local h <const> = font:getHeight() - padding
      local img <const> = gfx.image.new(w, h, gfx.kColorClear)
      gfx.lockFocus(img)
      gfx.setFont(font)
      gfx.drawTextAligned(text, w / 2, 0, kTextAlignment.center)
      gfx.unlockFocus()
      img:drawScaled(x - (scale * w) / 2, y - (scale * h) / 2, scale)
      return w, h
  end,

  drawTextRotated = function(text, x, y, angle, font)
      local padding = string.upper(text) == text and 6 or 0 -- Weird padding hack?
      local w <const> = font:getTextWidth(text)
      local h <const> = font:getHeight() - padding
      local img <const> = gfx.image.new(w, h, gfx.kColorClear)
      gfx.lockFocus(img)
      gfx.setFont(font)
      gfx.drawTextAligned(text, w / 2, 0, kTextAlignment.center)
      gfx.unlockFocus()
      img:drawRotated(x, y, angle)
      return w, h
  end
}
