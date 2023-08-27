import "CoreLibs/graphics"
import "CoreLibs/timer"
import "CoreLibs/ui"

local gfx <const> = playdate.graphics

local magicNumber = 0
local saveState = {["magicNumber"] = 0}
local saveTimer = nil

local evilNumber = nil

local piDigits = ".1415926535897932384626433832795028841971693993751058209749445923078164062862089986280348253421170679821480865132823066470938446095505822317253594081284811174502841027019385211055596446229489549303819644288109756659334461284756482337867831652712019091456485669234603486104543266482133936072602491412737245870066063155881748815209209628292540917153643678925903600113305305488204665213841469519415116094330572703657595919530921861173819326117931051185480744623799627495673518857527248912279381830119491298336733624406566430860213949463952247371907021798609437027705392171762931767523846748184676694051320005681271452635608277857713427577896091736371787214684409012249534301465495853710507922796892589235420199561121290219608640344181598136297747713099605187072113499999983729780499510597317328160963185950244594553469083026425223082533446850352619311881710100031378387528865875332083814206171776691473035982534904287554687311595628638823537875937519577818577805321712268066130019278766111959092164201989"
local piTimer = nil
local piPos = 1
local piSynth = nil
local piSeed

function myGameSetUp()
    local data = playdate.datastore.read()
    if data ~= nil then
        saveState = data

        magicNumber = saveState.magicNumber
    end

    saveTimer = playdate.timer.new(1000, function()
        if magicNumber == saveState.magicNumber then
            return
        end

        saveState.magicNumber = magicNumber

        print("saving...", saveState.magicNumber)
        playdate.datastore.write(saveState)
    end)
    saveTimer.repeats = true

    evilNumber = gfx.image.new("images/evil.png")
    assert(evilNumber)

    playdate.ui.crankIndicator:start()

    print("HELO")
end

myGameSetUp()

-- function genPrimes(upto)
--     if n == 2 then
--         return
--     end

--     local table = {}
--     for i = 1, n do
--         table[i] = true
--     end
--     sqrtn = math.tointeger(math.ceil(math.sqrt(n)))

--     -- Starting with 2, for each True (prime) number I in the table, mark all
--     -- its multiples as composite (starting with I*I, since earlier multiples
--     -- should have already been marked as multiples of smaller primes).
--     -- At the end of this process, the remaining True items in the table are
--     -- primes, and the False items are composites.
--     for i in range(2, sqrtn):
--         if table[i]:
--             for j in range(i * i, n, i):
--                 table[j] = False

--     # Yield all the primes in the table.
--     yield 2
--     for i in range(3, n, 2):
--         if table[i]:
--             yield i
-- end

-- from https://devforum.play.date/t/add-a-drawtextscaled-api-see-code-example/7108
function playdate.graphics.drawTextScaled(text, x, y, scale, font)
    local padding = string.upper(text) == text and 6 or 0 -- Weird padding hack?
    local w <const> = font:getTextWidth(text)
    local h <const> = font:getHeight() - padding
    local img <const> = gfx.image.new(w, h, gfx.kColorClear)
    gfx.lockFocus(img)
    gfx.setFont(font)
    gfx.drawTextAligned(text, w / 2, 0, kTextAlignment.center)
    gfx.unlockFocus()
    img:drawScaled(x - (scale * w) / 2, y - (scale * h) / 2, scale)
end

function playdate.graphics.drawTextRotated(text, x, y, angle, font)
    local padding = string.upper(text) == text and 6 or 0 -- Weird padding hack?
    local w <const> = font:getTextWidth(text)
    local h <const> = font:getHeight() - padding
    local img <const> = gfx.image.new(w, h, gfx.kColorClear)
    gfx.lockFocus(img)
    gfx.setFont(font)
    gfx.drawTextAligned(text, w / 2, 0, kTextAlignment.center)
    gfx.unlockFocus()
    img:drawRotated(x, y, angle)
end

function playdate.update()
    playdate.ui.crankIndicator:update()

    playdate.display.setInverted(magicNumber < 0)

    gfx.clear()

    if magicNumber == 3 then
        if piTimer == nil then
            piPos = 0
            piSynth = playdate.sound.synth.new(playdate.sound.kWaveNoise)
            piSeed = playdate.getSecondsSinceEpoch()
            piTimer = playdate.timer.new(5000, function()
                math.randomseed(piSeed)
                if math.random() > 0.9 then
                    piSynth:playMIDINote("C3", 0.05)
                end
                piTimer = playdate.timer.new(100+math.random(300), function()
                    if piPos < string.len(piDigits) and piPos < 108 then
                        piPos += 1
                    else
                        piSynth:stop()
                    end
                end)
                piTimer.repeats = true
            end)
        end
    elseif magicNumber ~= 3 and piTimer ~= nil then
        piTimer:remove()
        piTimer = nil

        piSynth:stop()
        piSynth = nil

        piSeed = nil
    elseif magicNumber == 666 or magicNumber == 6 then
        local width, height = evilNumber:getSize()
        local scale = 3
        evilNumber:drawScaled(200-5-(width/2)*scale, 120-(height/2)*scale, scale)
    -- TODO: 3.145...
    -- TODO: prime numbers (?), maybe naive check that brute forces things?
    -- TODO: 583 the best
    -- TODO: 999 HI SCORE (hi!)
    else
        gfx.drawTextScaled("*" .. tostring(magicNumber) .. "*", 200, 120, 4, gfx.getSystemFont())
    end

    if piTimer ~= nil then
        gfx.drawTextScaled("*" .. tostring(magicNumber) .. "*", 200, 120, 4, gfx.getSystemFont())

        if piPos > 0 then
            math.randomseed(piSeed)
            if math.random() < 0.9 then
                gfx.drawTextInRect(string.sub(piDigits, 1, piPos), 222, 130, 180, 120)
            else
                for i = 1,piPos do
                    gfx.drawTextRotated(string.sub(piDigits, i, i), math.random(400), math.random(240), math.random(360), gfx.getSystemFont())
                end
            end
        end
    end

    if playdate.buttonJustPressed( playdate.kButtonA ) then
        magicNumber += 1
    end
    if playdate.buttonJustPressed( playdate.kButtonB ) then
        magicNumber -= 1
    end

    -- gfx.sprite.update()
    playdate.timer.updateTimers()

end

function playdate.cranked(change, acceleratedChange)
    magicNumber += math.floor(acceleratedChange / 10)
end
