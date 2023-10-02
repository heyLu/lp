local gfx <const> = playdate.graphics
local geom <const> = playdate.geometry

local lines = table.pack(
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

local function update()
  gfx.clear()
  gfx.fillRect(0, 0, 400, 240)

  local rays = {}
  for i, line in ipairs(lines) do
    gfx.drawLine(line)

    table.insert(rays, {ray=lineBetween(player.x, player.y, line.x1, line.y1), seg=line})
    table.insert(rays, {ray=lineBetween(player.x, player.y, line.x2, line.y2), seg=line})
  end

  local zeroAngle = geom.vector2D.new(-1, 0)

  local intersects = {}
  for _, ray in ipairs(rays) do
    local minDist = 1000000
    local intersect = nil
    for _, line in ipairs(lines) do
      local doesIntersect, point = ray.ray:intersectsLineSegment(line)
      if doesIntersect then
        local dist = geom.lineSegment.new(player.x, player.y, point.x, point.y):length()
        -- print(doesIntersect, point, dist, ray, line, ray:segmentVector():normalize(), line:segmentVector():normalize())
        if ray.seg == line then -- origin is ourselves, just add it
          table.insert(intersects, {point=point, ray=ray.ray, dist=dist})
        elseif dist < minDist then -- closest intersection
          minDist = dist
          intersect = point
        end
      end
    end

    if intersect ~= nil then
      table.insert(intersects, {point = intersect, ray = ray.ray, dist=minDist})
    end
  end

  -- sort intersections to paint shade polygon in the right order
  table.sort(intersects, function(a, b)
    local angleA = zeroAngle:angleBetween(a.ray:segmentVector())
    local angleB = zeroAngle:angleBetween(b.ray:segmentVector())
    -- if angleA == angleB then
    --   return a.dist < b.dist
    -- end
    return angleA < angleB
  end)

  gfx.pushContext()

  gfx.setColor(gfx.kColorWhite)
  local points = {}
  for _, pr in ipairs(intersects) do
    if pr.point.x ~= 0 and pr.point.y ~= 0 then
      gfx.drawRect(pr.point.x-2, pr.point.y-2, 5, 5)
      table.insert(points, pr.point)
    end
  end

  gfx.setColor(gfx.kColorXOR)
  local shade = geom.polygon.new(table.unpack(points))
  shade:close()

  gfx.setColor(gfx.kColorWhite)
  gfx.setPolygonFillRule(gfx.kPolygonFillNonZero)
  gfx.fillPolygon(shade)

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
