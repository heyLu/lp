-- inspired by javidx9's tutorial: https://www.youtube.com/watch?v=ih20l3pJoeU

local gfx <const> = playdate.graphics

local function vector3(x, y, z)
  return {x=x, y=y, z=z}
end

local function triangle(a, b, c)
  return {a=a, b=b, c=c}
end

local function mat4x4(a11, a12, a13, a14,
                      a21, a22, a23, a24,
                      a31, a32, a33, a34,
                      a41, a42, a43, a44)
  return table.pack(
    table.pack(a11, a12, a13, a14),
    table.pack(a21, a22, a23, a24),
    table.pack(a31, a32, a33, a34),
    table.pack(a41, a42, a43, a44)
  )
end

local function multiplyMatrix4x4Vector3(m, v)
  local res = vector3(
    v.x * m[1][1] + v.y * m[2][1] + v.z * m[3][1] + m[4][1],
    v.x * m[1][2] + v.y * m[2][2] + v.z * m[3][2] + m[4][2],
    v.x * m[1][3] + v.y * m[2][3] + v.z * m[3][3] + m[4][3]
  )
  local w = v.x * m[1][4] + v.y * m[2][4] + v.z * m[3][4] + m[4][4]
  if w ~= 0 then
    res.x /= w
    res.y /= w
    res.z /= w
  end
  return res
end

local cube = {
  tris = table.pack(
    -- south
    triangle(vector3(0, 0, 0), vector3(0, 1, 0), vector3(1, 1, 0)),
    triangle(vector3(0, 0, 0), vector3(1, 1, 0), vector3(1, 0, 0)),

    -- east
    triangle(vector3(1, 0, 0), vector3(1, 1, 0), vector3(1, 1, 1)),
    triangle(vector3(1, 0, 0), vector3(1, 1, 1), vector3(1, 0, 1)),

    -- north
    triangle(vector3(1, 0, 1), vector3(1, 1, 1), vector3(0, 1, 1)),
    triangle(vector3(1, 0, 1), vector3(0, 1, 1), vector3(0, 0, 1)),

    -- west
    triangle(vector3(0, 0, 1), vector3(0, 1, 1), vector3(0, 1, 0)),
    triangle(vector3(0, 0, 1), vector3(0, 1, 0), vector3(0, 0, 0)),

    -- top
    triangle(vector3(0, 1, 0), vector3(0, 1, 1), vector3(1, 1, 1)),
    triangle(vector3(0, 1, 0), vector3(1, 1, 1), vector3(1, 1, 0)),

    -- bottom
    triangle(vector3(1, 0, 1), vector3(0, 0, 1), vector3(0, 0, 0)),
    triangle(vector3(1, 0, 1), vector3(0, 0, 0), vector3(1, 0, 0))
  ),
}

local screenWidth = 400
local screenHeight = 240
local aspectRatio = screenHeight/screenWidth
local zNear = 0.1
local zFar = 1000
local fov = math.rad(90)
local factor = 1 / math.tan(fov/2)
local proj = mat4x4(
  aspectRatio*factor, 0, 0, 0,
  0, factor, 0, 0,
  0, 0, zFar / (zFar - zNear), 1,
  0, 0, (-zFar*zNear) / (zFar - zNear), 0
)

local function update()
  gfx.clear()
  playdate.drawFPS(0, 0)

  playdate.display.setRefreshRate(100)

  local speed = 1
  local theta = playdate.getCurrentTimeMilliseconds() / 1000 * speed
  local rotX = mat4x4(
    math.cos(theta), math.sin(theta), 0, 0,
    -math.sin(theta), math.cos(theta), 0, 0,
    0, 0, 1, 0,
    0, 0, 0, 1
  )

  -- print("---")
  for _, tri in ipairs(cube.tris) do
    local transTri = table.deepcopy(tri)

    transTri.a = multiplyMatrix4x4Vector3(rotX, transTri.a)
    transTri.b = multiplyMatrix4x4Vector3(rotX, transTri.b)
    transTri.c = multiplyMatrix4x4Vector3(rotX, transTri.c)

    -- slightly offset to "look into" things
    transTri.a.x -= 0.5
    transTri.b.x -= 0.5
    transTri.c.x -= 0.5

    -- move back
    transTri.a.z += 3
    transTri.b.z += 3
    transTri.c.z += 3

    local projTri = {}
    projTri.a = multiplyMatrix4x4Vector3(proj, transTri.a)
    projTri.b = multiplyMatrix4x4Vector3(proj, transTri.b)
    projTri.c = multiplyMatrix4x4Vector3(proj, transTri.c)

    -- scale into view
    projTri.a.x += 1
    projTri.a.y += 1
    projTri.b.x += 1
    projTri.b.y += 1
    projTri.c.x += 1
    projTri.c.y += 1

    projTri.a.x *= 0.5*screenWidth
    projTri.a.y *= 0.5*screenHeight
    projTri.b.x *= 0.5*screenWidth
    projTri.b.y *= 0.5*screenHeight
    projTri.c.x *= 0.5*screenWidth
    projTri.c.y *= 0.5*screenHeight

    gfx.drawTriangle(
      projTri.a.x, projTri.a.y,
      projTri.b.x, projTri.b.y,
      projTri.c.x, projTri.c.y
    )
    -- print(
    --   projTri.a.x, projTri.a.y,
    --   projTri.b.x, projTri.b.y,
    --   projTri.c.x, projTri.c.y
    -- )
  end
end

return update
