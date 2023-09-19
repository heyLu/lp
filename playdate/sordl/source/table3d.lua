---@class Pos
---@field x integer
---@field y integer
---@field z integer

Table3D = {
  models = {},
  rotated = {angle=999},
}

function Table3D.addModel(self, name, model)
  self.models[name] = model
end

local tileWidthHalf  <const> = 64 / 2
local tileHeightHalf <const> = 32 / 2

function Table3D.toScreenPos(pos)
  -- https://clintbellanger.net/articles/isometric_math/
    local heightOffsetY = -pos.z * (tileHeightHalf*2 - 3)
    return (pos.x - pos.y) * (tileWidthHalf - 7),
           (pos.x + pos.y) * (tileHeightHalf-3) + heightOffsetY
end

function Table3D.toMapPos(_)
  return -999, -999, -999 -- FIXME: not implemented
end

function Table3D.hasChanged(_, _)
end

local function map(tbl, f)
  local t = {}
  for k,v in pairs(tbl) do
    t[k] = f(v)
  end
  return t
end

function Table3D.draw(self, world, opts)
  local angle = opts.angle or 0
  local scale = opts.scale or 0.5

  if self.rotated.angle ~= angle then
    local cos = math.cos(math.rad(-angle))
    local sin = math.sin(math.rad(-angle))
    local rotate = function(pos)
      -- rotation curtesy of https://gamedev.stackexchange.com/questions/186667/rotation-grid-positions
      local rx = pos.x * cos - pos.y * sin
      local ry = pos.x * sin + pos.y * cos
      return {x=rx, y=ry, z=pos.z, model=pos.model}
    end

    self.rotated = map(world.positions, rotate)
    table.sort(self.rotated,
      ---@param a Pos
      ---@param b Pos
      function(a, b)
        -- https://gamedev.stackexchange.com/questions/103442/how-do-i-determine-the-draw-order-of-isometric-2d-objects-occupying-multiple-til
        return a.x + a.y + a.z < b.x + b.y + b.z
      end
    )

    self.rotated.angle = angle
  end

  for _, pos in pairs(self.rotated) do
    if type(pos) == "table" then
      local sx, sy = self.toScreenPos(pos)
      local model = self.models["default"]
      if pos.model ~= nil then
        model = pos.model
      end
      local frame = 1+math.floor(angle/360*#model) -- #model is supposedly expensive?
      model[frame]:drawScaled(200+sx*scale, 120+sy*scale, scale)
      -- model[frame]:draw(200+sx, 120+sy)
    end
  end
end

