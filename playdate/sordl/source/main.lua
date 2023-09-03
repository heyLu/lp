import "CoreLibs/animator"
import "CoreLibs/graphics"
import "CoreLibs/object"

local gfx <const> = playdate.graphics

function make(x, y)
  local len = 10
  return {
    pos = {x = x, y = y},
    arc = playdate.geometry.arc.new(0, 0, len, 0, 360),
    dir = 10,
    speed = 2,
    len = len,
    shield = false,

    anim = nil,

    draw = function(self)
      gfx.setDitherPattern(0.3, playdate.graphics.image.kDitherTypeDiagonalLine)
      -- gfx.fillRect(self.pos.x-5, self.pos.y-5, 10, 10)
      gfx.fillCircleAtPoint(self.pos.x, self.pos.y, 5)
      gfx.setDitherPattern(0)

      if not self.shield then
        if self.anim == nil then
          local dx, dy = self.arc:pointOnArc(self.dir):unpack()
          gfx.drawLine(self.pos.x, self.pos.y, self.pos.x+dx, self.pos.y+dy)
        else
          local dx, dy = self.arc:pointOnArc(self.anim:currentValue()):unpack()
          gfx.drawLine(self.pos.x, self.pos.y, self.pos.x+dx, self.pos.y+dy)

          if self.anim:ended() then
            self.anim = nil
          end
        end
      else
        -- TODO: draw arc "behind" character to get an offset arc for bigger size
        local dx, dy = self.arc:pointOnArc(self.dir):unpack()
        local angle = (self.dir / self.arc:length()) * 360
        gfx.drawArc(self.pos.x+dx*3, self.pos.y+dy*3, 5, angle-30, angle+30)
      end
    end,

    move = function(self, button)
      if self.anim ~= nil then
        return
      end

      if playdate.buttonIsPressed(playdate.kButtonB) then
        self.speed = 1
      else
        self.speed = 2
      end

      if button == playdate.kButtonLeft then
        self.pos.x = self.pos.x - self.speed
        self.dir = self.arc:length()*0.75
      end
      if button == playdate.kButtonRight then
        self.pos.x = self.pos.x + self.speed
        self.dir = self.arc:length()*0.25
      end
      if button == playdate.kButtonUp then
        self.pos.y = self.pos.y - self.speed
        self.dir = 0
      end
      if button == playdate.kButtonDown then
        self.pos.y = self.pos.y + self.speed
        self.dir = self.arc:length()*0.5
      end
    end,

    attack = function(self)
      if self.anim ~= nil then
        return
      end

      self.anim = playdate.graphics.animator.new(700, self.dir, self.dir+self.arc:length()*0.25)
    end,
  }
end

function toTilePos(pos)
    local virtualTileX = pos.x / 8
    local virtualTileY = pos.y / 8

    local isoTileX = virtualTileX - (400 / 8) / 2
    local isoTileY = virtualTileY - (240 / 8) / 2

    return math.floor(isoTileX+0.5), math.floor(isoTileY+0.5)
end

local player
local map = nil
local platform = nil

function initGame()
  map = playdate.graphics.image.new("map.png")
  assert(map)

  platform = playdate.graphics.image.new("platform.png")
  assert(platform)

  player = make(200+4, 120+2)
end

initGame()

function playdate.update()
  gfx.clear()
  playdate.drawFPS(380, 2)

  platform:draw(0, 0)
  local width, height = map:getSize()
  map:draw(400-1-width, 15)

  local x, y = toTilePos(player.pos)
  gfx.drawText(tostring(x).." "..tostring(y), 5, 220)

  gfx.setColor(playdate.graphics.kColorXOR)
  gfx.drawPixel(400-1-width + x + width/2, 15 + y + height/2)
  gfx.setColor(playdate.graphics.kColorBlack)

  player:draw()

  if playdate.buttonIsPressed(playdate.kButtonLeft) then
    player:move(playdate.kButtonLeft)
  end
  if playdate.buttonIsPressed(playdate.kButtonRight) then
    player:move(playdate.kButtonRight)
  end
  if playdate.buttonIsPressed(playdate.kButtonUp) then
    player:move(playdate.kButtonUp)
  end
  if playdate.buttonIsPressed(playdate.kButtonDown) then
    player:move(playdate.kButtonDown)
  end

  if playdate.buttonJustPressed(playdate.kButtonA) then
    player:attack()
  end
  if playdate.buttonIsPressed(playdate.kButtonB) then
    player.shield = true
  else
    player.shield = false
  end
end
