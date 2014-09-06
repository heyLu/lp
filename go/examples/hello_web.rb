require 'webrick'

server = WEBrick::HTTPServer.new :Port => 8080

server.mount_proc '/' do |req, res|
  name = req.path == '/' ? "World" : req.path[1..-1]
  res.body = "Hello, #{name}!"
end

trap 'INT' do server.shutdown end
server.start