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

      -- fix moving diagonally (y is 2x the pixel size of x)
      local bothDirections = (playdate.buttonIsPressed(playdate.kButtonLeft) or playdate.buttonIsPressed(playdate.kButtonRight)) and (playdate.buttonIsPressed(playdate.kButtonUp) or playdate.buttonIsPressed(playdate.kButtonDown))
      print(bothDirections)
      local speedChange = 1
      if bothDirections then
        speedChange = 2
      end

      if button == playdate.kButtonLeft then
        self.pos.x = self.pos.x - self.speed
        -- self.pos.x = self.pos.x - 1
        -- self.pos.y = self.pos.y - 0.5
        self.dir = self.arc:length()*0.75
      end
      if button == playdate.kButtonRight then
        self.pos.x = self.pos.x + self.speed
        -- self.pos.x = self.pos.x + 1
        -- self.pos.y = self.pos.y + 0.5
        self.dir = self.arc:length()*0.25
      end
      if button == playdate.kButtonUp then
        self.pos.y = self.pos.y - self.speed/speedChange
        -- self.pos.x = self.pos.x + 1
        -- self.pos.y = self.pos.y - 0.5
        self.dir = 0
      end
      if button == playdate.kButtonDown then
        self.pos.y = self.pos.y + self.speed/speedChange
        -- self.pos.x = self.pos.x - 1
        -- self.pos.y = self.pos.y + 0.5
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

local numberOfTilesInX <const> = 16
local numberOfTilesInY <const> = 8

function toTilePos(pos)
    local virtualTileX = pos.x / numberOfTilesInX
    local virtualTileY = pos.y / numberOfTilesInY

    local isoTileX = virtualTileX - (400 / numberOfTilesInX) / 2
    local isoTileY = virtualTileY - (240 / numberOfTilesInY) / 2

    return math.floor(isoTileX+0.5), math.floor(isoTileY+0.5)
end

function toScreenPos(pos)
  -- local screenTileX = pos.x + (400 / numberOfTilesInX) / 2
  -- local screenTileY = pos.y + (240 / numberOfTilesInY) / 2

  -- return math.floor(screenTileX * numberOfTilesInX), math.floor(screenTileY * numberOfTilesInY)

  -- https://clintbellanger.net/articles/isometric_math/
  return (pos.x - pos.y) * numberOfTilesInX / 2, (pos.x + pos.y) * numberOfTilesInY / 2
end

function toMapPos(pos, width, height)
    local screenTileX = pos.x + (width / 1) / 2
    local screenTileY = pos.y + (height / 1) / 2

    return math.floor(screenTileX + 0.5), math.floor(screenTileY + 0.5)
end

local player
local map = nil
local platform = nil
local brick = nil

local world = {}

function initGame()
  map = playdate.graphics.image.new("map.png")
  assert(map)

  platform = playdate.graphics.image.new("platform.png")
  assert(platform)

  brick = playdate.graphics.image.new("brick.png")
  assert(brick)

  player = make(200+2, 120)

  local worldData = playdate.datastore.read("world")
  if worldData ~= nil then
    world = worldData
  end
end

function playdate.gameWillTerminate()
  print("saving")
  print("saved")
end

initGame()

local cursor = {
  x = 0,
  y = 0,
}

function playdate.update()
  gfx.clear()
  playdate.drawFPS(385, 2)

  for x = -12, 12, 0.5 do
    for y = -14, 14, 0.5 do
      if world[x] ~= nil and world[x][y] then
        local sx, sy = toScreenPos({x = x, y = y})
        brick:draw(sx, sy)
      end
    end
  end

  local sx, sy = toScreenPos(cursor)
  gfx.setColor(gfx.kColorXOR)
  gfx.drawPixel(sx, sy)
  gfx.drawRect(sx, sy, 16, 16)

  gfx.drawText(tostring(cursor.x).." "..tostring(cursor.y), 5, 220)

  if playdate.buttonJustPressed(playdate.kButtonA) then
    if world[cursor.x] == nil then
      world[cursor.x] = {}
    end

    world[cursor.x][cursor.y] = not world[cursor.x][cursor.y]
    playdate.datastore.write(world)
  end

  local change = 1
  local bothDirections = (playdate.buttonIsPressed(playdate.kButtonLeft) or playdate.buttonIsPressed(playdate.kButtonRight)) and (playdate.buttonIsPressed(playdate.kButtonUp) or playdate.buttonIsPressed(playdate.kButtonDown))
  if bothDirections then
    change = 0.5
  end

  if playdate.buttonJustPressed(playdate.kButtonLeft) then
    cursor.x = cursor.x - change
  end
  if playdate.buttonJustPressed(playdate.kButtonRight) then
    cursor.x = cursor.x + change
  end
  if playdate.buttonJustPressed(playdate.kButtonUp) then
    cursor.y = cursor.y - change
  end
  if playdate.buttonJustPressed(playdate.kButtonDown) then
    cursor.y = cursor.y + change
  end
end

function play()
  gfx.clear()
  playdate.drawFPS(385, 2)

  platform:draw(0, 0)

  local x, y = toTilePos(player.pos)

  local xOffset = 10
  local yOffset = 7

  local mapCopy = map:rotatedImage(45)
  local width, height = mapCopy:getSize()
  playdate.graphics.lockFocus(mapCopy)
  -- local sx, sy = toMapPos({x = x, y = y}, width, height)
  -- local sx = player.pos.x/400*width
  -- local sy = player.pos.y/240*height
  local sx = math.floor(x/width * width) + width/2 --- xOffset
  local sy = math.floor(y/height * height) + height/2 --- yOffset
  local color = mapCopy:sample(sx, sy)
  if color == playdate.graphics.kColorBlack then
    color = playdate.graphics.kColorWhite
  else
    color = playdate.graphics.kColorBlack
  end
  gfx.setColor(color)
  gfx.drawPixel(sx, sy)
  playdate.graphics.unlockFocus()

  mapCopy:draw(400-50, 15)

  gfx.drawText(tostring(x).." "..tostring(y).." / "..tostring(sx).." "..tostring(sy).." / "..tostring(player.pos.x).." "..tostring(player.pos.y), 5, 220)

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
