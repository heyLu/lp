# pixl!

draw pixls.

![what it looks like](screen.png)

see the [whirlwind intro][] for a quickstart.

it's live at <http://pixl.papill0n.org>.

[whirlwind intro]: https://github.com/heyLu/codegirls/blob/master/2013-12-17-christmas-special.md

# trixl!

`pixl` in 3 dimensions. you can even import the data from pixl.

see it live at <http://pixl.papill0n.org/3>.

    // quick start (type the following in the console):
    trixl.generate.fun()

## todo

* code sharing (especially for worlds created and dynamics, but also for
  libraries/modules that could be shared)
* multiplayer backend as for pixl
* live-editing from a built-in editor (similar to firefox' scratchpad,
  but on the page. should also support saving code to `localStorage`.)
* better performance. noticable lag at 10000 trixls, probably
  batch-drawing would help, but how? first find out what the bottleneck
  is.
* bug: sometimes when moving the mouse to certain positions (edges of
  the screen?) the camera rapidly alternates between two positions
* bug: doesn't work in chrome, as it's missing support for iterators and
  the `keys()` function on maps. i didn't find any convenient polyfills
  (`keys()` should return something iterator-like, e.g. returning raw
  values on `next()` and `StopIteration` when the end is reached.)
* bug: overwriting trixls isn't supported as `Map#set` compares by
  reference and not by value. maybe use chunks instead? (could be a
  performance boost as well.)
