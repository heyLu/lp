local gfx <const> = playdate.graphics
local geom <const> = playdate.geometry

local walls = table.pack(
  -- screen "walls"
  geom.lineSegment.new(0, 0, 400, 0), -- top
  geom.lineSegment.new(0, 240, 400, 240), -- bottom
  geom.lineSegment.new(400, 0, 400, 240), -- left
  geom.lineSegment.new(0, 0, 0, 240), -- right

  geom.lineSegment.new(100, 110, 100, 130),
  geom.lineSegment.new(300, 110, 300, 130),
  geom.lineSegment.new(150, 50, 200, 25)
)

local player = geom.point.new(200, 120)

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

local function update()
  gfx.clear()
  gfx.fillRect(0, 0, 400, 240)

  local uniquePoints = {}
  local seenPoint = {}
  for _, line in ipairs(walls) do
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
    local hit, _ = intersect(ray, walls) -- iterates over lines, makes this O^2
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

  if wallStart then
    gfx.pushContext()
    gfx.setColor(gfx.kColorXOR)
    gfx.drawLine(wallStart.x, wallStart.y, player.x, player.y)
    gfx.popContext()
  end

  if playdate.buttonJustPressed(playdate.kButtonA) then
    if wallStart then
      local wall = geom.lineSegment.new(wallStart.x, wallStart.y, player.x, player.y)
      table.insert(walls, wall)
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
end

return update
