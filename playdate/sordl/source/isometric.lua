import "CoreLibs/graphics"

local gfx <const> = playdate.graphics

Isometric = {
  models = {},

  cachedLayers = {},
  isCached = {},
}

function Isometric.addModel(self, name, model)
  self.models[name] = model
end

local tileWidthHalf <const> = 16 / 2
local tileHeightHalf <const> = 8 / 2

function Isometric.toTilePos(pos)
  local x = pos.x - 200
  local y = pos.y - 120
  return math.floor((x / tileWidthHalf + y / tileHeightHalf) / 2),
         math.floor((y / tileHeightHalf - (x / tileWidthHalf)) / 2),
         pos.z
end

function Isometric.toScreenPos(pos)
  local heightOffsetY = -pos.z * (tileHeightHalf*2)
  -- https://clintbellanger.net/articles/isometric_math/
  return 200 + (pos.x - pos.y) * tileWidthHalf, 120 + (pos.x + pos.y) * tileHeightHalf + heightOffsetY
end

function Isometric.hasChanged(self, pos)
  self.isCached[pos.z] = false
end

function Isometric.draw(self, world, opts)
  playdate.resetElapsedTime()
  local from = opts.from or -10
  local to = opts.to or 10
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
      local fadedBrick = self.models["default"]:fadedImage((10-math.abs(h))/10, gfx.image.kDitherTypeFloydSteinberg)
      -- local fadedBrick = brick:blurredImage(3, 1, gfx.image.kDitherTypeFloydSteinberg)
      for x = 0, 52, 1 do
        for y = -24, 39, 1 do
          local tile = world:getTile({x = x, y = y, z = h})
          if tile then
            local sx, sy = Isometric.toScreenPos({x = x, y = y, z = h})
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
