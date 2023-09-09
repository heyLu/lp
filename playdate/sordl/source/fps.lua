import "CoreLibs/graphics"

local gfx <const> = playdate.graphics

FPS = {}

FPS.new = function(x, y, width, height)
  local history = {}

  for i=0,width,1 do
    history[i] = playdate.display.getRefreshRate()
  end

  return {
    offset = 0,
    history = history,
    line = playdate.geometry.polygon.new(width),

    draw = function(self)
      local targetRate = playdate.display.getRefreshRate()

      self.history[self.offset] = playdate.getFPS()
      if self.history[self.offset] < 0 then
        -- playdate wrongly shows negative fps at the start?
        self.history[self.offset] = targetRate
      end

      gfx.drawLine(x, y+height, x, y) -- vertical axis
      gfx.drawLine(x, y+height, x+width, y+height) -- horizontal axis

      local p = width
      for i=self.offset,1,-1 do
        self.line:setPointAt(p, x+p, y+height-math.max(self.history[i]/targetRate*height, 0))
        p = p-1
      end
      for i=#self.history,self.offset+1,-1 do
        self.line:setPointAt(p, x+p, y+height-math.max(self.history[i]/targetRate*height, 0))
        p = p-1
      end

      gfx.drawPolygon(self.line)

      self.offset = (self.offset + 1) % 60
    end,
  }
end
