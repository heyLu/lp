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

local function intersect(origin, ray, lines)
  local minDist = math.huge
  local hit = nil
  local hitLine = nil
  for _, line in ipairs(lines) do
    local doesIntersect, point = ray:intersectsLineSegment(line)
    if doesIntersect then
      local isOrigin = false --(line.x1 == origin.dx and line.y1 == origin.dy) or (line.x2 == origin.dx and line.y2 == origin.dy)
      local dist = geom.lineSegment.new(ray.x1, ray.y1, point.x, point.y):length()
      if not isOrigin and dist < minDist then
        hit = point
        hitLine = line
      end
    end
  end
  return hit, hitLine
end

local function update()
  gfx.clear()
  gfx.fillRect(0, 0, 400, 240)

  local endpoints = {}
  local seen = {}
  for _, line in ipairs(walls) do
    gfx.drawLine(line)

    local a = geom.vector2D.new(line.x1, line.y1)
    if not seen[tostring(a)] then
      table.insert(endpoints, a)
      seen[tostring(a)] = true
    end
    local b = geom.vector2D.new(line.x2, line.y2)
    if not seen[tostring(b)] then
      table.insert(endpoints, b)
      seen[tostring(b)] = true
    end
  end

  local zeroAngle = geom.vector2D.newPolar(200, 0)
  table.sort(endpoints, function(a, b)
    local angleA = zeroAngle:angleBetween(geom.lineSegment.new(200, 120, a.dx, a.dy):segmentVector())
    local angleB = zeroAngle:angleBetween(geom.lineSegment.new(200, 120, b.dx, b.dy):segmentVector())
    return angleA < angleB
  end)
  printTable(endpoints)

  gfx.pushContext()
  gfx.setColor(gfx.kColorWhite)

  print("---")
  local n = 1
  local points = {}
  local currentWall = nil
  for _, endpoint in ipairs(endpoints) do
    local newPoint, newWall = intersect(endpoint, lineBetween(player.x, player.y, endpoint.dx, endpoint.dy), walls)
    assert(newPoint)
    assert(newWall)

    if #points == 0 then
      print("init!")
      -- table.insert(points, {x=newWall.x1, y=newWall.y1})
      -- table.insert(points, {x=newWall.x2, y=newWall.y2})
      currentWall = newWall
    end

    print("endpoint intersect", endpoint, newPoint)
    if newWall ~= currentWall then
      print("new wall!", #points, currentWall)
      if #points < 2 then
        print("supplement", newPoint)
        points[2] = newPoint
        -- points[2] = geom.point.new(currentWall.x2, currentWall.y2)
      end
      print("intersect", points[1], points[2])
      local start, _ = intersect(nil, lineBetween(player.x, player.y, points[1].x, points[1].y), {[1]=currentWall})
      local _end, _ = intersect(nil, lineBetween(player.x, player.y, points[2].x, points[2].y), {[1]=currentWall})
      if _end == nil then
        print("retry intersect")
        _end, _ = intersect(nil, lineBetween(player.x, player.y, currentWall.x1, currentWall.y1), {[1]=currentWall})
      end
      if start == nil or _end == nil then
        print("no intersection")
        break
      end
      print("tri", player.x, player.y, start.x, start.y, _end.x, _end.y)
      gfx.fillTriangle(player.x, player.y, start.x, start.y, _end.x, _end.y)
      -- TODO: remove (newest?) point
      table.remove(points, 1)
      if n > 5 then
        print("break")
        break
      end
      n += 1

      currentWall = newWall
    else --if #points < 2 then
      print("add", newPoint)
      table.insert(points, newPoint)
      -- if #points > 2 then
      --   table.remove(points, 1)
      -- end
    end
  end

  gfx.setImageDrawMode(gfx.kDrawModeNXOR)
  gfx.setColor(gfx.kColorXOR)

  for _, point in ipairs(endpoints) do
    gfx.drawRect(point.dx-2, point.dy-2, 5, 5)
    -- gfx.drawText(i, point.dx+5, point.dy-5)
  end

  gfx.popContext()

  gfx.pushContext()
  gfx.setColor(gfx.kColorXOR)
  gfx.setLineWidth(2)
  for _, endpoint in ipairs(endpoints) do
    gfx.drawLine(player.x, player.y, endpoint.dx, endpoint.dy)
  end
  gfx.popContext()

  gfx.drawCircleAtPoint(player.x, player.y, 3)

  if playdate.buttonIsPressed(playdate.kButtonUp) then
    player.y = player.y - 1
  end
  if playdate.buttonIsPressed(playdate.kButtonDown) then
    player.y = player.y + 1
  end
  if playdate.buttonIsPressed(playdate.kButtonLeft) then
    player.x = player.x - 1
  end
  if playdate.buttonIsPressed(playdate.kButtonRight) then
    player.x = player.x + 1
  end
end

return update
