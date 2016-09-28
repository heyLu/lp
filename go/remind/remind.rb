require 'chronic'
require 'json'

REMINDERS_PATH = 'remind.json'

reminders = []
if File.exist? REMINDERS_PATH
  reminders = JSON.parse File.read(REMINDERS_PATH)
end

now = Time.now
if ARGV.length == 0
  reminders.each do |r|
    time = Chronic.parse(r['time'])
    if time > now
      print "#{time} - #{r['description']}"
    end
  end
else
  time = Chronic.parse(ARGV[0])
  description = ARGV[1]

  reminders << { time: time, description: description }
end

File.write(REMINDERS_PATH, JSON.dump(reminders))
