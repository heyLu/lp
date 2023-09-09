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

    jump = function(self)
      self.pos.y = self.pos.y - 30
    end
  }
end

local tileWidthHalf <const> = 16 / 2
local tileHeightHalf <const> = 8 / 2

function toTilePos(pos)
  return math.floor((pos.x / tileWidthHalf + pos.y / tileHeightHalf) / 2),
         math.floor((pos.y / tileHeightHalf - (pos.x / tileWidthHalf)) / 2)
end

function toScreenPos(pos)
  -- https://clintbellanger.net/articles/isometric_math/
  return (pos.x - pos.y) * tileWidthHalf, (pos.x + pos.y) * tileHeightHalf
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

local state = {
  editMode = true,
}

local world = {
  layers = {},
  offset = {x = 50, y = 50, z = 0},
}

function world.load(self)
  for i = -10, 10, 1 do
    local name = "world-"..tostring(i)
    local img = playdate.datastore.readImage(name)
    -- print("tried to load "..name.." -> "..tostring(img))
    if img ~= nil then
      self.layers[i] = img
    end
  end
end

function world.save(self)
  for i = -10, 10, 1 do
    if self.layers[i] ~= nil then
      local name = "world-"..tostring(i)
      playdate.datastore.writeImage(self.layers[i], name)
      if playdate.isSimulator then
        playdate.datastore.writeImage(self.layers[i], name..".gif")
        -- print("saved "..name)
      end
    end
  end
end

function world.getTile(self, pos)
  -- TODO: check if table is faster than image in memory
  local layer = self.layers[pos.z+self.offset.z]
  if layer == nil then
    return false
  end
  local color = layer:sample(pos.x+self.offset.x, pos.y+self.offset.y)
  if color ~= playdate.graphics.kColorBlack then
    return false
  end

  return true
end

function world.setTile(self, pos, tile)
  if self.layers[pos.z+self.offset.z] == nil then
    self.layers[pos.z+self.offset.z] = playdate.graphics.image.new(100, 100)
  end
  local layer = self.layers[pos.z+self.offset.z]

  playdate.graphics.pushContext(layer)
  local color = playdate.graphics.kColorClear
  if tile then
    color = playdate.graphics.kColorBlack
  end
  playdate.graphics.setColor(color)
  playdate.graphics.drawPixel(pos.x+self.offset.x, pos.y+self.offset.y)
  playdate.graphics.popContext()
  return tile
end

function initGame()
  map = playdate.graphics.image.new("map.png")
  assert(map)

  platform = playdate.graphics.image.new("platform.png")
  assert(platform)

  brick = playdate.graphics.image.new("brick.png")
  assert(brick)

  local sx, sy = toScreenPos({x = 15, y = 3})
  player = make(sx, sy)

  local savedState = playdate.datastore.read()
  if savedState ~= nil then
    state = savedState
  end

  world:load()

  playdate.getSystemMenu():addCheckmarkMenuItem("edit", state.editMode, function()
    state.editMode = not state.editMode
  end)
end

function playdate.gameWillTerminate()
  print("saving")
  playdate.datastore.write(state)
  world:save()
  print("saved")
end

initGame()

local cursor = {
  x = 0,
  y = 0,
  z = 0,
}

function fix(pos)
  return {x = math.floor(pos.x), y = math.floor(pos.y), z = pos.z}
end

local tilesChanged = true
local tileCache = gfx.image.new(400, 240)

local fpsHistory = {}
local fpsOffset = 0

for i=0,61,1 do
  fpsHistory[i] = 30
end

function playdate.update()
  if tilesChanged then
    playdate.resetElapsedTime()
    gfx.pushContext(tileCache)
    gfx.clear()
    for h = -10,10,1 do
      local offsetY = -h * (tileHeightHalf*2)
      for x = 0, 52, 1 do
        for y = -24, 39, 1 do
          local tile = world:getTile({x = x, y = y, z = h})
          if tile then
            local sx, sy = toScreenPos({x = x, y = y})
            brick:draw(sx, sy+offsetY)
          end
        end
      end
    end
    gfx.popContext()

    local cacheTime = playdate.getElapsedTime()
    print("tile cache took "..tostring(cacheTime*1000).."ms")
    tilesChanged = false
  end

  gfx.clear()
  tileCache:draw(0, 0)
  playdate.drawFPS(385, 2)

  fpsHistory[fpsOffset] = math.max(playdate.getFPS(), 0)
  local p = 60
  for i=fpsOffset,1,-1 do
    gfx.drawLine(320+p, 16, 320+p, 16-math.max(fpsHistory[i]/2, 0))
    p = p-1
  end
  for i=#fpsHistory,fpsOffset+1,-1 do
    gfx.drawLine(320+p, 16, 320+p, 16-math.max(fpsHistory[i]/2, 0))
    p = p-1
  end
  fpsOffset = (fpsOffset + 1) % 60

  player:draw()

  if state.editMode then
    edit()
  else
    play()
  end
end

function edit()
  local pos = fix(cursor)
  gfx.drawText(tostring(pos.x).." "..tostring(pos.y).." @ "..tostring(cursor.z).." "..tostring(world:getTile(pos)), 5, 220)

  local offsetY = -cursor.z * (tileHeightHalf*2)
  local sx, sy = toScreenPos({x = pos.x, y = pos.y})
  gfx.setColor(gfx.kColorXOR)
  gfx.drawRect(sx, sy+offsetY, 16, 16)

  if playdate.buttonJustPressed(playdate.kButtonA) then
    local tile = world:getTile(pos)
    world:setTile(pos, not tile)
    tilesChanged = true
  end

  if playdate.buttonIsPressed(playdate.kButtonB) then
    if playdate.buttonJustPressed(playdate.kButtonUp) then
      cursor.z = math.min(cursor.z + 1, 10)
    elseif playdate.buttonJustPressed(playdate.kButtonDown) then
      cursor.z = math.max(cursor.z - 1, -10)
    end

    return
  end

  if playdate.buttonJustPressed(playdate.kButtonLeft) then
    cursor.x = cursor.x - 0.5
    cursor.y = cursor.y + 0.5
  end
  if playdate.buttonJustPressed(playdate.kButtonRight) then
    -- cursor.x = cursor.x + change
    cursor.x = cursor.x + 0.5
    cursor.y = cursor.y - 0.5
  end
  if playdate.buttonJustPressed(playdate.kButtonUp) then
    -- cursor.y = cursor.y - change
    cursor.x = cursor.x - 1
    cursor.y = cursor.y - 1
  end
  if playdate.buttonJustPressed(playdate.kButtonDown) then
    -- cursor.y = cursor.y + change
    cursor.x = cursor.x + 1
    cursor.y = cursor.y + 1
  end
end

function play()
  local x, y = toTilePos(player.pos)
  gfx.drawText(tostring(x).." "..tostring(y), 0, 220)

  local sx, sy = toScreenPos({x = x, y = y})
  local tx, ty = toTilePos({x = sx, y = sy})
  gfx.drawText(tostring(tx).." "..tostring(ty), 50, 220)


  if player.pos.y > 240 then
    gfx.drawText("*oops*", 200, 120)
    -- playdate.stop()
    return
  end

  if world[x] ~= nil and world[x][y] then
    gfx.drawText("*on*", 100, 220)
  else
    gfx.drawText("off", 100, 220)
    player.pos.y = player.pos.y + 5
  end

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

  if playdate.buttonIsPressed(playdate.kButtonA) and playdate.buttonJustPressed(playdate.kButtonUp) then
    player:jump()
  elseif playdate.buttonJustPressed(playdate.kButtonA) then
    player:attack()
  end
  if playdate.buttonIsPressed(playdate.kButtonB) then
    player.shield = true
  else
    player.shield = false
  end
end
