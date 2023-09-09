import "CoreLibs/animator"
import "CoreLibs/graphics"
import "CoreLibs/object"
import "CoreLibs/timer"

import "fps"

local gfx <const> = playdate.graphics

local function translateMovement(button, factor)
  local offsetX = 0
  local offsetY = 0

  if button == playdate.kButtonLeft then
    offsetX = -0.5
    offsetY = 0.5
  end

  if button == playdate.kButtonRight then
    offsetX = 0.5
    offsetY = -0.5
  end

  if button == playdate.kButtonUp then
    offsetX = -1
    offsetY = -1
  end

  if button == playdate.kButtonDown then
    offsetX = 1
    offsetY = 1
  end

  return offsetX*factor, offsetY*factor
end

local function make(x, y)
  local len = 10
  return {
    pos = {x = x, y = y, z = 0},
    arc = playdate.geometry.arc.new(0, 0, len, 0, 360),
    dir = 10,
    speed = 0.01,
    len = len,
    shield = false,
    collisionImage = gfx.image.new(400, 240),
    sprite = nil,

    anim = nil,

    draw = function(self)
      local sx, sy = toScreenPos(self.pos)
      local screenPos = {x = sx, y = sy}

      gfx.setDitherPattern(0.3, gfx.image.kDitherTypeDiagonalLine)
      if self.sprite == nil then
        -- gfx.fillRect(self.pos.x-5, self.pos.y-5, 10, 10)
        gfx.fillCircleAtPoint(screenPos.x, screenPos.y, 5)
      else
        self.sprite:drawAnchored(screenPos.x, screenPos.y, 0.25, 0.25)
      end
      gfx.setDitherPattern(0)

      -- if not self.shield then
      --   if self.anim == nil then
      --     local dx, dy = self.arc:pointOnArc(self.dir):unpack()
      --     gfx.drawLine(screenPos.x, screenPos.y, screenPos.x+dx, screenPos.y+dy)
      --   else
      --     local dx, dy = self.arc:pointOnArc(self.anim:currentValue()):unpack()
      --     gfx.drawLine(screenPos.x, screenPos.y, screenPos.x+dx, screenPos.y+dy)

      --     if self.anim:ended() then
      --       self.anim = nil
      --     end
      --   end
      -- else
      --   -- TODO: draw arc "behind" character to get an offset arc for bigger size
      --   local dx, dy = self.arc:pointOnArc(self.dir):unpack()
      --   local angle = (self.dir / self.arc:length()) * 360
      --   gfx.drawArc(screenPos.x+dx*3, screenPos.y+dy*3, 5, angle-30, angle+30)
      -- end
    end,

    updateCollision = function(self)
      self.collisionImage:clear(gfx.kColorClear)
      gfx.pushContext(self.collisionImage)

      local sx, sy = toScreenPos(self.pos)
      local screenPos = {x = sx, y = sy}
      gfx.fillCircleAtPoint(screenPos.x, screenPos.y, 3)

      gfx.popContext(self.collisionImage)
    end,

    move = function(self, button, world)
      if self.anim ~= nil then
        return
      end

      if playdate.buttonIsPressed(playdate.kButtonB) then
        self.speed = 0.125
      else
        self.speed = 0.25
      end

      -- fix moving diagonally (y is 2x the pixel size of x)
      -- local bothDirections = (playdate.buttonIsPressed(playdate.kButtonLeft) or playdate.buttonIsPressed(playdate.kButtonRight)) and (playdate.buttonIsPressed(playdate.kButtonUp) or playdate.buttonIsPressed(playdate.kButtonDown))
      -- local speedChange = 1
      -- if bothDirections then
      --   speedChange = 2
      -- end

      if button == playdate.kButtonLeft then
        self.dir = self.arc:length()*0.75
      end
      if button == playdate.kButtonRight then
        self.dir = self.arc:length()*0.25
      end
      if button == playdate.kButtonUp then
        self.dir = 0
      end
      if button == playdate.kButtonDown then
        self.dir = self.arc:length()*0.5
      end

      local offsetX, offsetY = translateMovement(button, self.speed)
      self.pos.x = self.pos.x + offsetX
      self.pos.y = self.pos.y + offsetY
      if math.abs(offsetX) > 0 or math.abs(offsetY) > 0 then
        self:updateCollision()
      end

      local layer = world.cachedLayers[math.floor(self.pos.z+1)]
      if layer ~= nil then
        if world:getTile({x = math.floor(self.pos.x+offsetX), y = math.floor(self.pos.y+offsetY), z = self.pos.z + 1}) then
        -- if gfx.checkAlphaCollision(layer, 0, 0, gfx.kImageUnflipped, self.collisionImage, 0, 0, gfx.kImageUnflipped) then
          -- colliding with the next level up, stop!

          -- need to undo collision/position update
          self.pos.x = self.pos.x - offsetX
          self.pos.y = self.pos.y - offsetY
          if math.abs(offsetX) > 0 or math.abs(offsetY) > 0 then
            self:updateCollision()
          end

          return
        end
      end
    end,

    attack = function(self)
      if self.anim ~= nil then
        return
      end

      self.anim = gfx.animator.new(700, self.dir, self.dir+self.arc:length()*0.25)
    end,

    jump = function(self)
      self.anim = nil

      self.pos.z = self.pos.z + 4
    end
  }
end

local tileWidthHalf <const> = 16 / 2
local tileHeightHalf <const> = 8 / 2

function toTilePos(pos)
  return math.floor((pos.x / tileWidthHalf + pos.y / tileHeightHalf) / 2),
         math.floor((pos.y / tileHeightHalf - (pos.x / tileWidthHalf)) / 2),
         pos.z
end

function toScreenPos(pos)
  local heightOffsetY = -pos.z * (tileHeightHalf*2)
  -- https://clintbellanger.net/articles/isometric_math/
  return (pos.x - pos.y) * tileWidthHalf, (pos.x + pos.y) * tileHeightHalf + heightOffsetY
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

  cachedLayers = {},
  isCached = {},
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
  if color ~= gfx.kColorBlack then
    return false
  end

  return true
end

function world.setTile(self, pos, tile)
  if self.layers[pos.z+self.offset.z] == nil then
    self.layers[pos.z+self.offset.z] = gfx.image.new(100, 100)
  end
  local layer = self.layers[pos.z+self.offset.z]

  gfx.pushContext(layer)
  local color = gfx.kColorClear
  if tile then
    color = gfx.kColorBlack
  end
  gfx.setColor(color)
  gfx.drawPixel(pos.x+self.offset.x, pos.y+self.offset.y)
  gfx.popContext()

  self.isCached[pos.z+self.offset.z] = false
  return tile
end

function world.draw(self, from, to)
  playdate.resetElapsedTime()
  local wasNotCached = false
  for h = math.max(-10, from),math.min(to, 10),1 do
    if not self.isCached[h] then
      wasNotCached = true

      if self.cachedLayers[h] == nil then
        self.cachedLayers[h] = gfx.image.new(400, 240)
      end

      self.cachedLayers[h]:clear(gfx.kColorClear)

      gfx.pushContext(self.cachedLayers[h])
      -- gfx.drawRect(0, 0, 400, 240)
      local fadedBrick = brick:fadedImage((10-math.abs(h))/10, gfx.image.kDitherTypeFloydSteinberg)
      -- local fadedBrick = brick:blurredImage(3, 1, gfx.image.kDitherTypeFloydSteinberg)
      for x = 0, 52, 1 do
        for y = -24, 39, 1 do
          local tile = world:getTile({x = x, y = y, z = h})
          if tile then
            local sx, sy = toScreenPos({x = x, y = y, z = h})
            fadedBrick:draw(sx, sy)
          end
        end
      end
      gfx.popContext()

      self.isCached[h] = true
    end

    gfx.setImageDrawMode(gfx.kDrawModeCopy)
    self.cachedLayers[h]:draw(0, 0)
  end

  if wasNotCached then
    local tookMs = playdate.getElapsedTime()*1000
    local availableMs = 1 / playdate.display.getRefreshRate() * 1000
    print("tile cache took "..tookMs.."ms ("..(tookMs/availableMs*100).."% of "..availableMs.."ms)")
  end
end

local cursor = {
  x = 0,
  y = 0,
  z = 0,
}

function fix(pos)
  return {x = math.floor(pos.x), y = math.floor(pos.y), z = pos.z}
end

local repeatTimer = nil
function removeRepeatTimer()
  if repeatTimer ~= nil then
    repeatTimer:remove()
  end
end

local modeEdit = {
  name = "edit",
  update = function()
    local pos = fix(cursor)
    gfx.drawText(tostring(pos.x).." "..tostring(pos.y).." @ "..tostring(cursor.z).." "..tostring(world:getTile(pos)), 5, 220)

    gfx.pushContext()
    local sx, sy = toScreenPos(pos)
    gfx.setColor(gfx.kColorXOR)
    gfx.drawRect(sx, sy, 16, 16)
    gfx.popContext()

    if playdate.buttonJustPressed(playdate.kButtonA) then
      local tile = world:getTile(pos)
      world:setTile(pos, not tile)
    end

    if playdate.buttonIsPressed(playdate.kButtonB) then
      if playdate.buttonJustPressed(playdate.kButtonUp) then
        cursor.z = math.min(cursor.z + 1, 10)
      elseif playdate.buttonJustPressed(playdate.kButtonDown) then
        cursor.z = math.max(cursor.z - 1, -10)
      end

      return
    end
  end,

  inputHandlers = {
    leftButtonDown = function()
      removeRepeatTimer()
      repeatTimer = playdate.timer.keyRepeatTimer(function()
        cursor.x = cursor.x - 0.5
        cursor.y = cursor.y + 0.5
      end)
    end,
    leftButtonUp = function()
      removeRepeatTimer()
    end,

    rightButtonDown = function()
      removeRepeatTimer()
      repeatTimer = playdate.timer.keyRepeatTimer(function()
        cursor.x = cursor.x + 0.5
        cursor.y = cursor.y - 0.5
      end)
    end,
    rightButtonUp = function()
      removeRepeatTimer()
    end,

    upButtonDown = function()
      if playdate.buttonIsPressed(playdate.kButtonB) then
        return
      end
      removeRepeatTimer()
      repeatTimer = playdate.timer.keyRepeatTimer(function()
        cursor.x = cursor.x - 1
        cursor.y = cursor.y - 1
      end)
    end,
    upButtonUp = function()

      removeRepeatTimer()
    end,

    downButtonDown = function()
      if playdate.buttonIsPressed(playdate.kButtonB) then
        return
      end
      removeRepeatTimer()
      repeatTimer = playdate.timer.keyRepeatTimer(function()
        cursor.x = cursor.x + 1
        cursor.y = cursor.y + 1
      end)
    end,
    downButtonUp = function()
      removeRepeatTimer()
    end,
  },
}

local falling = false
local lastOnPos = {x = 0, y = 0, z = 0}
local modePlay = {
  name = "play",
  update = function()
    gfx.drawText(tostring(player.pos.x).." "..tostring(player.pos.y), 0, 220)

    local sx, sy = toScreenPos(player.pos)
    local tx, ty = toTilePos({x = sx, y = sy})
    gfx.drawText(tostring(tx).." "..tostring(ty), 50, 220)


    if sy > 240 then
      gfx.drawText("*oops*", 200, 120)
      -- playdate.stop()
      return
    end

    local layer = world.cachedLayers[math.floor(player.pos.z)]
    local collide = false
    if layer ~= nil then
      collide = gfx.checkAlphaCollision(layer, 0, 0, gfx.kImageUnflipped, player.collisionImage, 0, 0, gfx.kImageUnflipped)
    end

    local correctedPos = table.shallowcopy(player.pos) --{x = math.floor(player.pos.x), y = math.floor(player.pos.y), z = player.pos.z}
    local onTile = world:getTile(correctedPos)
    if onTile or collide then
      lastOnPos = correctedPos
      gfx.drawText("*on*", 100, 220)
    else --if math.abs(lastOnPos.x-correctedPos.x) > 0.5 or math.abs(lastOnPos.y-correctedPos.y)>0.5 then
      if not falling then -- now falling...
        printTable(correctedPos,lastOnPos)
      end
      gfx.drawText("off", 100, 220)
      player.pos.z = player.pos.z - 0.25
      falling = true
    end

    if playdate.buttonIsPressed(playdate.kButtonLeft) then
      player:move(playdate.kButtonLeft, world)
    end
    if playdate.buttonIsPressed(playdate.kButtonRight) then
      player:move(playdate.kButtonRight, world)
    end
    if playdate.buttonIsPressed(playdate.kButtonUp) then
      player:move(playdate.kButtonUp, world)
    end
    if playdate.buttonIsPressed(playdate.kButtonDown) then
      player:move(playdate.kButtonDown, world)
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
  end,
  inputHandlers = {},
}

local mode = modeEdit

function setupState()
  playdate.inputHandlers.pop()
  mode = modeEdit
  if not state.editMode then
    mode = modePlay
  end
  playdate.inputHandlers.push(mode.inputHandlers)
end

function initGame()
  map = gfx.image.new("map.png")
  assert(map)

  platform = gfx.image.new("platform.png")
  assert(platform)

  brick = gfx.image.new("brick.png")
  assert(brick)

  ghost = gfx.image.new("ghost.png")
  assert(ghost)

  player = make(15, 3)
  player.sprite = ghost

  local savedState = playdate.datastore.read()
  if savedState ~= nil then
    state = savedState
  end

  world:load()

  setupState()
  playdate.getSystemMenu():addCheckmarkMenuItem("edit", state.editMode, function()
    state.editMode = not state.editMode
    setupState()
  end)
end

function playdate.gameWillTerminate()
  print("saving")
  playdate.datastore.write(state)
  world:save()
  print("saved")
end

initGame()

local fps = FPS.new(320, 2, 60, 16)

function playdate.update()
  -- playdate.display.setInverted(true)

  gfx.clear()
  playdate.drawFPS(385, 2)

  fps:draw()
 
  -- TODO: proper draw order would mean drawing front to back or kind of diagonally?
  world:draw(-10, math.floor(player.pos.z)+1)
  player:draw()
  world:draw(math.floor(player.pos.z)+2, 10)

  mode.update()

  playdate.timer.updateTimers()
end
