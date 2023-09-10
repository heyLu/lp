import "CoreLibs/graphics"
import "CoreLibs/ui"

local gfx <const> = playdate.graphics

local columnNote = 1
local columnVelocity = 2

local notes = {
  "C", "Db", "D", "Eb", "E", "F", "Gb", "Bb", "B", "Ab", "A", "C",
}

function toNoteString(pitch)
  local note = math.floor(pitch%12)
  local octave = math.floor(pitch//12)
  return notes[note+1]..tostring(octave-1)
end

local sequence = playdate.sound.sequence.new()

local fuerElise = playdate.sound.sequence.new("fuer-elise.mid")

local globalEffect = playdate.sound.twopolefilter.new(playdate.sound.kFilterBandPass)
local filterFreq = 1000
globalEffect:setMix(1)
-- globalEffect:setResonance(0.3)
globalEffect:setFrequency(filterFreq)
playdate.sound.addEffect(globalEffect)

local bitmore = gfx.font.new("bitmore")
assert(bitmore)
bitmore:setTracking(1)
gfx.setFont(bitmore)

function makeTrack(waveform)
  local track1 = playdate.sound.track.new()
  for step = 1, 64, 1 do
    if waveform == playdate.sound.kWaveNoise then
      if (step-1) % 4 == 0 then
        local note = "C7"
        if (step-1)%16 == 0 then
          note = "C4"
        end
        track1:addNote(step, note, 1, 0.1)
      end
    else
      local r = math.random()
      if r > 0.95 then
        track1:addNote(step, 60+math.random(12), 4, math.random(127)/255)
      else
        track1:addNote(step, "C4", 4, 0.0)
      end
    end
  end
  track1:setInstrument(playdate.sound.synth.new(waveform))

  local track1View = playdate.ui.gridview.new(25, 11)
  track1View:setNumberOfColumns(2)
  track1View:setNumberOfRows(track1:getLength())
  track1View:setScrollDuration(100)

  local active = false
  local setActive = function(a)
    active = a
  end

  track1View.drawCell = function(self, section, row, column, selected, x, y, width, height)
    local note = track1:getNotes(row)[1]

    local noteString = "---"
    local velocityString = "--"
    if note ~= nil and note.velocity > 0 then
      noteString = toNoteString(note.note)
      velocityString = tostring(math.floor(note.velocity*255))
    end

    selected = selected and active

    if column == columnNote then
      if selected then
        -- gfx.drawRect(x, y, 25*2-2, 11)
        gfx.fillRect(x, y, width, height)
        gfx.setImageDrawMode(gfx.kDrawModeFillWhite)
      else
        gfx.setImageDrawMode(gfx.kDrawModeCopy)
      end
      gfx.drawTextInRect(noteString, x+1, y+1, width, height, nil, nil, kTextAlignment.left)
    elseif column == columnVelocity then
      if selected then
        gfx.fillRect(x, y, width, height)
        gfx.setImageDrawMode(gfx.kDrawModeFillWhite)
      else
        gfx.setImageDrawMode(gfx.kDrawModeCopy)
      end
      gfx.drawTextInRect(velocityString, x+1, y+1, width, height, nil, nil, kTextAlignment.left)  end
  end

  return {
    track = track1,
    view = track1View,
    setActive = setActive,
  }
end

local tracks = {
  makeTrack(playdate.sound.kWaveSine),
  -- makeTrack(playdate.sound.kWaveSquare),
  -- makeTrack(playdate.sound.kWaveSawtooth),
  -- makeTrack(playdate.sound.kWaveTriangle),
  makeTrack(playdate.sound.kWaveNoise),
  makeTrack(playdate.sound.kWavePOPhase),
  makeTrack(playdate.sound.kWavePODigital),
  -- makeTrack(playdate.sound.kWavePOVosim),
}

-- local t = 1
-- for i=1,fuerElise:getTrackCount(),1 do
--   if tracks[i] == nil then
--     break
--   end

--   local track = fuerElise:getTrackAtIndex(i)
--   if track:getLength() > 0 then
--     tracks[t].track:setNotes(track:getNotes())
--     t = t + 1
--   end
-- end

local selectedTrack = 1

function setActive()
  for i, track in ipairs(tracks) do
    track.setActive(i == selectedTrack)
  end
end

for _, track in ipairs(tracks) do
  sequence:addTrack(track.track)
end

local bpm = 120
sequence:setTempo(bpm/7.5) -- FIXME: not accurate, not sure how to exact bpm with integer steps?
sequence:setLoops(0)
sequence:play()

local playhead = playdate.ui.gridview.new(10, 11)
playhead:setNumberOfColumns(1)
playhead:setNumberOfRows(tracks[1].track:getLength())
playhead:setScrollDuration(100)
playhead.drawCell = function(self, section, row, column, selected, x, y, width, height)
  if selected then
    gfx.drawTextInRect(">", x, y, width, height, nil, nil, kTextAlignment.left)
  end
end

local selectHandlers = {}
local editHandlers = {}

local function moveToPreviousColumn()
  local _, _, oldColumn = tracks[selectedTrack].view:getSelection()
  tracks[selectedTrack].view:selectPreviousColumn()
  local section, row, column = tracks[selectedTrack].view:getSelection()

  if column == oldColumn then
    selectedTrack = selectedTrack - 1
    if selectedTrack < 1 then
      selectedTrack = #tracks
    elseif selectedTrack > #tracks then
      selectedTrack = 1
    end
    tracks[selectedTrack].view:setSelection(section, row, columnVelocity)
  else
    tracks[selectedTrack].view:scrollToCell(section, row, column)
  end
end

local function moveToNextColumn()
  local _, _, oldColumn = tracks[selectedTrack].view:getSelection()
  tracks[selectedTrack].view:selectNextColumn()
  local section, row, column = tracks[selectedTrack].view:getSelection()
  tracks[selectedTrack].view:scrollToCell(section, row, column)

  if column == oldColumn then
    selectedTrack = selectedTrack + 1
    if selectedTrack < 1 then
      selectedTrack = #tracks
    elseif selectedTrack > #tracks then
      selectedTrack = 1
    end
    tracks[selectedTrack].view:setSelection(section, row, columnNote)
  else
    tracks[selectedTrack].view:scrollToCell(section, row, column)
  end
end

selectHandlers.AButtonDown = function()
  playdate.inputHandlers.pop()
  playdate.inputHandlers.push(editHandlers)
end
selectHandlers.BButtonDown = function()
  if sequence:isPlaying() then
    sequence:stop()
  else
    sequence:play()
  end
end
selectHandlers.leftButtonUp = function()
  moveToPreviousColumn()
end
selectHandlers.rightButtonUp = function()
  moveToNextColumn()
end
selectHandlers.upButtonUp = function()
  tracks[selectedTrack].view:selectPreviousRow(true)
  local section, row, column = tracks[selectedTrack].view:getSelection()
  for i=1,#tracks,1 do
    tracks[i].view:scrollToCell(section, row, column)
  end
end
selectHandlers.downButtonUp = function()
  tracks[selectedTrack].view:selectNextRow(true)
  local section, row, column = tracks[selectedTrack].view:getSelection()
  for i=1,#tracks,1 do
    tracks[i].view:scrollToCell(section, row, column)
  end
end
selectHandlers.cranked = function(change, acceleratedChange)
  filterFreq = filterFreq + change
  globalEffect:setFrequency(filterFreq)
end

editHandlers.BButtonUp = function()
  playdate.inputHandlers.pop()
  playdate.inputHandlers.push(selectHandlers)
end
editHandlers.leftButtonUp = function()
  moveToPreviousColumn()
end
editHandlers.rightButtonUp = function()
  moveToNextColumn()
end
editHandlers.upButtonUp = function()
  local track = tracks[selectedTrack].track
  local _, row, column = tracks[selectedTrack].view:getSelection()
  local note = track:getNotes(row)[1]
  track:removeNote(row, note.note)
  if note.velocity == 0 then
    note.velocity = 1.0
  else
    if column == columnNote then
      note.note = note.note + 1 % 127
    elseif column == columnVelocity then
      note.velocity = ((note.velocity*255 + 1)%255) / 255
    end
  end
  track:addNote(row, note.note, note.length, note.velocity)
  if not sequence:isPlaying() then
    track:getInstrument():allNotesOff()
    track:getInstrument():playMIDINote(note.note, note.velocity, note.length)
  end
end
editHandlers.downButtonUp = function()
  local track = tracks[selectedTrack].track
  local _, row, column = tracks[selectedTrack].view:getSelection()
  local note = track:getNotes(row)[1]
  track:removeNote(row, note.note)
  if note.velocity == 0 then
    note.velocity = 1.0
  else
    if column == columnNote then
      note.note = note.note - 1 % 127
    elseif column == columnVelocity then
      note.velocity = ((note.velocity*255 - 1)%255) / 255
    end
  end
  track:addNote(row, note.note, note.length, note.velocity)
  if not sequence:isPlaying() then
    track:getInstrument():allNotesOff()
    track:getInstrument():playMIDINote(note.note, note.velocity, note.length)
  end
end

playdate.inputHandlers.push(selectHandlers)

math.randomseed(playdate.getSecondsSinceEpoch())

function playdate.update()
  playdate.display.setInverted(true)

  setActive()

  local needsDisplay = playdate.getButtonState() ~= 0

  playdate.drawFPS(382, 2)

  gfx.pushContext()
  gfx.setColor(gfx.kColorWhite)
  gfx.fillRect(2, 2, 20, 10)
  gfx.popContext()
  if sequence:isPlaying() then
    gfx.fillTriangle(2, 2, 2, 9, 10, 6)
  else
    gfx.fillRect(2, 2, 8, 8)
  end
  gfx.setImageDrawMode(gfx.kDrawModeCopy)
  local info = tostring(bpm).."bpm".." "..tostring(sequence:getTempo()).."st/s  "..sequence:getTrackCount().." tracks  "..tracks[1].track:getLength().."steps"
  gfx.drawTextInRect(info, 16, 2, string.len(info)*10, 11)

  local currentStep = (sequence:getCurrentStep()%sequence:getLength())
  if playhead:getSelectedRow() ~= currentStep then
    playhead:setSelectedRow(currentStep)
  end
  if needsDisplay or playhead.needsDisplay then
    gfx.pushContext()
    gfx.setColor(gfx.kColorWhite)
    gfx.fillRect(0, 12, 20, 240-12)
    gfx.popContext()
    playhead:drawInRect(0, 20, 20, 240-12)
  end

  for i, track in ipairs(tracks) do
    local width = 80
    if needsDisplay or track.view.needsDisplay then
      gfx.pushContext()
      gfx.setColor(gfx.kColorWhite)
      gfx.fillRect(20+(i-1)*width, 12, width, 240-12)
      gfx.popContext()
      track.view:drawInRect(20+(i-1)*width, 12, width, 240-12)
    end
  end

  playdate.timer:updateTimers()
end