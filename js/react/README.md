# Playing with React

More opinions/impressions coming later, when I have played with this more.

## Impressions

* ES6 is quite cool, especially if you wrote just plain JS before (using
    [6to5](http://6to5.org))

    Some things: `import React from "react";`, template strings, `() => ...`,
* `webpack` is fun/convenient/surprisingly fast (and has support for
    all the things we use here, e.g. react/jsx, 6to5, css)
* JSX works well, even though I'd prefer doing it in Clojure (`hiccup`
    syntax, e.g. nested vectors & collections. has someone written a library
    for that already?)
* showdown, which was suggested by the tutorial, doesn't work with module
    loaders and it was a pain to figure out how to fix it. (the fix was
    using a different library.)

## How to run this thing

```
$ npm install -g webpack
$ npm install
$ webpack --watch

# visit index.html in your browser

$ vi entry.js      # change something!
```

## Ideas

* follow "the way", at least at first, see how it works.
    - e.g. JSX, 6to5, webpack (JSX is a react thing, 6to5 & webpack are for
        sanity)
* things to try
    - build a tiny thingy ("profile badge")
    - put it in a list
    - sort the list
    - filter the list
    - make the filter dynamic (optional/difficult)
* in the future:
    - try out immutable.js & reimplement the above
    - try out `css-layout` and/or `styles` in js

## Resources

* [React.js Conf 2015](https://www.youtube.com/playlist?list=PLb0IAmt7-GS1cbw4qonlQztYV1TAW0sCr)
    seemed fun & quite interesting. For inspiration, watch [HYPE!][hype]
* [om](https://github.com/omcljs/om), because you actually want to use
    Clojure ;)

[hype]: https://www.youtube.com/watch?v=z5e7kWSHWTg

## Misc

- `npm config set prefix $HOME/.local --global`
