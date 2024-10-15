// open with `scide`, ctrl-enter to run, ctrl-. to stop, ctrl-d to look up things.  (and maybe ctrl-b to start the sound server.)

// why?  because it's in norns, the rather mysterious sound box thing.

(
{ var ampOsc, freq;
	ampOsc = SinOsc.kr(100, 0, 0.2, 0.2);
	freq = [[660, 880], [440, 660], [1320, 880]].choose;
	SinOsc.ar(freq, 0, ampOsc);
}.play;
)

{ Pan2.ar(PinkNoise.ar(0.2), -1) }.play;

{ PinkNoise.ar(0.1) + SinOsc.ar(100, 0, 0.2) + Saw.ar(100, 0.2) }.play;

{ Mix.new([SinOsc.ar(440, 0, 0.2), Saw.ar(660, 0.2)]).postln }.play;

(
{
	var a, b;
	a = [SinOsc.ar(440, 0, 0.2), Saw.ar(662, 0.2)];
	b = [SinOsc.ar(442, 0, 0.2), Saw.ar(660, 0.2)];
	Mix([a, b]).postln;
}.play;
)

// also, you can put .scope and .plot anywhere to inspect things (also .postln)

// bunch of sines, maybe half-way to an organ?
(
var n = 1000;
{ Mix.fill(n, { SinOsc.ar(500 + 500.0.rand, 0, 1/n) }) }.play;
)