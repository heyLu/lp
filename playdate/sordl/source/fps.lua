import "CoreLibs/graphics"

local gfx <const> = playdate.graphics

FPS = {}

FPS.new = function(x, y, width, height)
  local history = {}

  for i=0,width,1 do
    history[i] = 30
  end

  return {
    offset = 0,
    history = history,
    line = playdate.geometry.polygon.new(width),

    draw = function(self)
      self.history[self.offset] = playdate.getFPS()
      if self.history[self.offset] < 0 then
        -- playdate wrongly shows negative fps at the start?
        self.history[self.offset] = 30
      end

      gfx.drawLine(x, y+height, x, y) -- vertical axis
      gfx.drawLine(x, y+height, x+width, y+height) -- horizontal axis

      local p = 60
      for i=self.offset,1,-1 do
        self.line:setPointAt(p, x+p, y+height-math.max(self.history[i]/2, 0))
        p = p-1
      end
      for i=#self.history,self.offset+1,-1 do
        self.line:setPointAt(p, x+p, y+height-math.max(self.history[i]/2, 0))
        p = p-1
      end

      gfx.drawPolygon(self.line)

      self.offset = (self.offset + 1) % 60
    end,
  }
end
