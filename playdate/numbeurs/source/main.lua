import "CoreLibs/graphics"
import "CoreLibs/timer"
import "CoreLibs/ui"

import "gfxext"

local gfx <const> = playdate.graphics
local gfxext <const> = playdate.graphicsext

local magicNumber = 0
local saveState = {["magicNumber"] = 0, timePlayed = 0}

local evilNumber

local piDigits = ".1415926535897932384626433832795028841971693993751058209749445923078164062862089986280348253421170679821480865132823066470938446095505822317253594081284811174502841027019385211055596446229489549303819644288109756659334461284756482337867831652712019091456485669234603486104543266482133936072602491412737245870066063155881748815209209628292540917153643678925903600113305305488204665213841469519415116094330572703657595919530921861173819326117931051185480744623799627495673518857527248912279381830119491298336733624406566430860213949463952247371907021798609437027705392171762931767523846748184676694051320005681271452635608277857713427577896091736371787214684409012249534301465495853710507922796892589235420199561121290219608640344181598136297747713099605187072113499999983729780499510597317328160963185950244594553469083026425223082533446850352619311881710100031378387528865875332083814206171776691473035982534904287554687311595628638823537875937519577818577805321712268066130019278766111959092164201989"
local piTimer = nil
local piPos = 1
local piSynth
local piSeed = nil

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

local function restoreGame()
    local data = playdate.datastore.read()
    if data ~= nil then
        saveState = data

        magicNumber = saveState.magicNumber
    end
end

local function saveGame()
    saveState.magicNumber = magicNumber

    print("saving...", saveState.magicNumber)
    playdate.datastore.write(saveState)
end

local function initGame()
    restoreGame()

    evilNumber = gfx.image.new("images/evil.png")
    assert(evilNumber)

    playdate.ui.crankIndicator:start()

    playdate.getSystemMenu():addMenuItem("save game", saveGame)

    print("HELO")
end

initGame()

function playdate.deviceWillSleep()
    print("zZzz")

    saveGame()
end

function playdate.gameWillTerminate()
    saveState.timePlayed = saveState.timePlayed + math.ceil(playdate.getCurrentTimeMilliseconds() / 1000)
    saveGame()

    print("BYEEE")
end

function playdate.update()
    playdate.display.setInverted(magicNumber < 0)

    gfx.clear()

    playdate.drawFPS(385, 0)

    local timePlayed = playdate.getCurrentTimeMilliseconds()
    if timePlayed < 2000 then
        playdate.ui.crankIndicator:update()
    end
    gfx.drawText(tostring(saveState.timePlayed + math.ceil(timePlayed / 1000)), 5, 220)

    local drawNumber = true
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
                        piPos = piPos + 1
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
        drawNumber = false

        local width, height = evilNumber:getSize()
        local scale = 3
        evilNumber:drawScaled(200-5-(width/2)*scale, 120-(height/2)*scale, scale)
    -- TODO: 3.145...
    -- TODO: prime numbers (?), maybe naive check that brute forces things?
    -- TODO: 583 the best
    -- TODO: 999 HI SCORE (hi!)
    -- TODO: 11 (D.M. number) -> "said" by a character/speech bubble after a while
    -- TODO: 12 also great
    -- TODO: 13: trying too hard
    elseif magicNumber == 42 then
        gfx.drawText("the answer...", 250, 130)
    elseif magicNumber % 2 == 0 then
        gfx.drawText("2", 200, 150)
    end

    local textWidth = 0
    local textHeight = 0
    if drawNumber then
        textWidth, textHeight = gfxext.drawTextScaled("*" .. tostring(magicNumber) .. "*", 200, 120, 4, gfx.getSystemFont())
    end

    local now = playdate.getTime()
    if magicNumber == now.year or magicNumber == now.weekday or magicNumber == now.month or magicNumber == now.day or magicNumber == now.hour or magicNumber == now.minute or magicNumber == now.second then
        gfx.drawText("*!*", 200+textWidth+4, 120+textHeight-2)
    end


    if piTimer ~= nil then
        if piPos > 0 then
            math.randomseed(piSeed)
            if math.random() < 0.9 then
                gfx.drawTextInRect(string.sub(piDigits, 1, piPos), 222, 130, 180, 120)
            else
                for i = 1,piPos do
                    gfxext.drawTextRotated(string.sub(piDigits, i, i), math.random(400), math.random(240), math.random(360), gfx.getSystemFont())
                end
            end
        end
    end

    if playdate.buttonJustPressed( playdate.kButtonA ) then
        magicNumber = magicNumber + 1
    end
    if playdate.buttonJustPressed( playdate.kButtonB ) then
        magicNumber = magicNumber - 1
    end

    -- gfx.sprite.update()
    playdate.timer.updateTimers()

end

function playdate.cranked(_, acceleratedChange)
    magicNumber = magicNumber + math.floor(acceleratedChange / 10)
end
