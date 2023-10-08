local gfx <const> = playdate.graphics
local geom <const> = playdate.geometry

local wallTop = geom.lineSegment.new(0, 0, 400, 0)
local wallBottom = geom.lineSegment.new(0, 240, 400, 240)
local wallLeft = geom.lineSegment.new(0, 0, 0, 240)
local wallRight = geom.lineSegment.new(400, 0, 400, 240)

local function loadLevel(levelName)
  local level = playdate.datastore.read(levelName)
  if level == nil then
    level = {
      walls = {}
    }
  end

  local walls = {}
  table.insert(walls, wallTop)
  table.insert(walls, wallBottom)
  table.insert(walls, wallLeft)
  table.insert(walls, wallRight)
  for _, wall in ipairs(level.walls) do
    table.insert(walls, geom.lineSegment.new(wall.x1, wall.y1, wall.x2, wall.y2))
  end

  return {
    walls = walls,
  }
end

local function saveLevel(levelName, level)
  local walls = {}
  for _, wall in ipairs(level.walls) do
    if wall ~= wallTop and wall ~= wallBottom and wall ~= wallLeft and wall ~= wallRight then
      table.insert(walls, {x1=wall.x1, y1=wall.y1, x2=wall.x2, y2=wall.y2})
    end
  end

  playdate.datastore.write({
    walls = walls
  }, levelName)
end

local currentLevel = "level1"
local level = loadLevel(currentLevel)

local foodSprites = gfx.imagetable.new("images/food")
assert(foodSprites)

local function canSee(pos, viewRange, angle, range, player)
  if pos:distanceToPoint(player) > viewRange then
    return false
  end

  local playerAngle = geom.vector2D.new(0, -1):angleBetween(geom.lineSegment.new(pos.x, pos.y, player.x, player.y):segmentVector())
  if playerAngle < 0 then
    playerAngle += 360
  end

  return angle-range/2 < playerAngle and playerAngle < angle+range/2
end

local player = geom.point.new(200, 120)
local caught = false

local function lineBetween(x1, y1, x2, y2)
  local ls = geom.lineSegment.new(x1, y1, x2, y2)
  local dx, dy = ls:segmentVector():unpack()
  return geom.lineSegment.new(x1, y1, x1+dx*1000000, y1+dy*1000000)
end

local function intersect(ray, lines)
  local minDist = math.huge
  local hit = nil
  local hitLine = nil
  for _, line in ipairs(lines) do
    local doesIntersect, point = ray:intersectsLineSegment(line)
    if doesIntersect then
      local dist = geom.lineSegment.new(ray.x1, ray.y1, point.x, point.y):length()
      if dist < minDist then
        hit = point
        hitLine = line

        minDist = dist
      end
    end
  end
  return hit, hitLine
end

local function haveSeen(t, val)
  local seen = t[val]
  if not seen then
    t[val] = true
  end
  return seen
end

local wallStart = nil

local function newConeTimer(dur, angle, range, easing)
  local t = playdate.timer.new(dur, angle-range/2, angle+range/2, easing)
  t.repeats = true
  t.reverses = true
  return t
end

local fruitDir = 90
local coneTimer = newConeTimer(2000, fruitDir, 60, playdate.easingFunctions.inOutSine)
local fruitTimer = playdate.timer.new(5000, 20, 100)
fruitTimer.repeats = true
fruitTimer.discardOnCompletion = false
fruitTimer.timerEndedCallback = function(timer)
  -- reverse manually
  local startValue = timer.startValue
  local endValue = timer.endValue
  timer.startValue = endValue
  timer.endValue = startValue
  timer:start()

  if fruitDir == 90 then
    fruitDir = 360-fruitDir
  else
    fruitDir = 90
  end
  coneTimer:remove()
  coneTimer = newConeTimer(2000, fruitDir, 60, playdate.easingFunctions.inOutSine)
end

-- coneTimer:pause()
-- fruitTimer:pause()

local function update()
  gfx.clear()
  gfx.fillRect(0, 0, 400, 240)

  local uniquePoints = {}
  local seenPoint = {}
  for _, line in ipairs(level.walls) do
    if not haveSeen(seenPoint, line.x1*1000000+line.y1) then
      table.insert(uniquePoints, geom.point.new(line.x1, line.y1))
    end
    if not haveSeen(seenPoint, line.x2*1000000+line.y2) then
      table.insert(uniquePoints, geom.point.new(line.x2, line.y2))
    end
  end

  local uniqueAngles = {}
  local seenAngle = {}
  for _, point in ipairs(uniquePoints) do
    local angle = math.atan(point.y-player.y, point.x-player.x)
    if not haveSeen(seenAngle, angle) then
      -- add angle of point AND slightly offset to hit points _behind_ walls (VERY important!)
      table.insert(uniqueAngles, angle-0.0001)
      table.insert(uniqueAngles, angle)
      table.insert(uniqueAngles, angle+0.0001)
    end
  end

  -- intersect with angles of all wall endpoints
  local intersects = {}
  for _, angle in ipairs(uniqueAngles) do
    local dx = math.cos(angle)
    local dy = math.sin(angle)

    local ray = lineBetween(player.x, player.y, player.x+dx, player.y+dy)
    local hit, _ = intersect(ray, level.walls) -- iterates over lines, makes this O^2
    if hit ~= nil then
      table.insert(intersects, {point=hit, angle=angle})
    end
  end

  table.sort(intersects, function(a, b)
    return a.angle < b.angle
  end)

  local points = {}
  for _, intersection in ipairs(intersects) do
    table.insert(points, intersection.point)
  end

  gfx.pushContext()
  gfx.setColor(gfx.kColorXOR)
  local light = geom.polygon.new(table.unpack(points))
  light:close()
  gfx.fillPolygon(light)
  gfx.popContext()

  -- debug intersections
  for _, point in ipairs(points) do
    gfx.drawRect(point.x-3, point.y-3, 6, 6)
    -- gfx.drawLine(player.x, player.y, point.x, point.y)
  end

  -- for _, line in ipairs(walls) do
  --   gfx.drawLine(player.x, player.y, line.x1, line.y1)
  --   gfx.drawLine(player.x, player.y, line.x2, line.y2)
  -- end

  gfx.drawCircleAtPoint(player.x, player.y, 3)

  local coneSize = 200
  local pos = geom.point.new(fruitTimer.value, 50)
  local found = canSee(geom.point.new(pos.x+5, pos.y+5), coneSize/2, coneTimer.value, 60, player)

  if not caught and found then
    local noise = playdate.sound.synth.new(playdate.sound.kWaveNoise)
    local crushEffect = playdate.sound.bitcrusher.new()
    crushEffect:setMix(0.5)
    crushEffect:setAmount(0.5)
    playdate.sound.addEffect(crushEffect)
    noise:playNote("C", 0.5, 0.5)

    caught = true
  end

  if not found then
    caught = false
  end

  gfx.pushContext()
  gfx.setImageDrawMode(gfx.kDrawModeWhiteTransparent)
  gfx.setColor(gfx.kColorBlack)

  foodSprites:getImage(1, 2):invertedImage():draw(pos.x, pos.y)

  if found then
    gfx.setStencilPattern(0.7)
  else
    gfx.setStencilPattern(0.3)
  end
  gfx.fillEllipseInRect(pos.x-coneSize/2+5, pos.y-coneSize/2+5, coneSize, coneSize, coneTimer.value-30, coneTimer.value+30)

  gfx.popContext()

  if wallStart then
    gfx.pushContext()
    gfx.setColor(gfx.kColorXOR)
    gfx.drawLine(wallStart.x, wallStart.y, player.x, player.y)
    gfx.popContext()
  end

  if playdate.buttonJustPressed(playdate.kButtonA) then
    if wallStart then
      local wall = geom.lineSegment.new(wallStart.x, wallStart.y, player.x, player.y)
      table.insert(level.walls, wall)
      wallStart = nil
    else
      wallStart = geom.point.new(player.x, player.y)
    end
  end

  if playdate.buttonIsPressed(playdate.kButtonUp) then
    player.y = player.y - 2
  end
  if playdate.buttonIsPressed(playdate.kButtonDown) then
    player.y = player.y + 2
  end
  if playdate.buttonIsPressed(playdate.kButtonLeft) then
    player.x = player.x - 2
  end
  if playdate.buttonIsPressed(playdate.kButtonRight) then
    player.x = player.x + 2
  end

  playdate.timer.updateTimers()
end

return {
  update = update,
  gameWillTerminate = function()
    saveLevel(currentLevel, level)
  end,
}
