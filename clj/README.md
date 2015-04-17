# Learning (aka. playing with) Clojure

- first heard about it early 2012 (I think)
- dipped my feet in a few times since then
- came across (i.e. read or saw things they did) a few interesting
  people using/writing Clojure in interesting way (fogus, Rich Hickey,
  Chris Granger)
- now reading 'The Joy of Clojure' (and enjoying it)

## Tools

- mostly LightTable
- previously Vim with the VimClojure plugin, but not right now
  (LightTable 'feels' more interactive to me and what I miss most from
  Vim (keybindings) will be in LT soonish, I also think that something
  like LT has more potential to be extensible and has a much saner
  extension language)
- the `clj` repl

## Little thoughts

- a lot of interesting people with interesting ideas use Clojure, so
    maybe it encourages thinking about problems first? you can simply
    start coding something before having any idea what you're going to
    do (which is both good and bad).

## Giggles & quibbles

+ the giggles
    * homoiconic
    * has macros
    * good feature inheritance (supposedly, not as in OO)
    * lots of cool projects (datomic, matchure, ring, lighttable,
      typed-clojure)
- the quibbles
    * the dynamic type-system continues to bite me
        could be my fault, but I often jump right in without reading the
        whole documentation and have spend too much time tracking weird
        errors down that were all fixed by changing one place in the
        code
    * learning new libraries without sufficient docs/tutorials (e.g.
        what are the arguments to this function, what does this map mean;
        maybe clojure wants me to think differently)
    * the jvm overhead is bad for my little computer (too slow, too much
        memory consumed)

        it also makes it difficult to quickly explore some ideas without
        creating a new project (in most cases, mostly to get libraries),
        starting a jvm... etc.

        i'd really love a version of clojure without this overhead, even
        if it is a lot slower. maybe [clojure-scheme](https://github.com/takeoutweight/clojure-scheme)
        could be that.

## Hickups

Some examples for the quibbles.

* `(if (= something something-undefined) true-expr false-expr)` not
    throwing an exeption for `something-undefined` being not defined

## How To Learn

- write something that you wish existed
- play with the repl (Clojure excels at manipulating data, so take some data and play with it)
    - `doc` and `ns-publics` are your friends
- koans:
    - [4clojure](http://www.4clojure.com/)
    - ["Clojure Katas inspired by Alice in Wonderland"](https://github.com/gigasquid/wonderland-clojure-katas)
