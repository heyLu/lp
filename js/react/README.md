# Playing with React

Opinions/impressions coming later, when I have played with this more.

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

## Resources

* [React.js Conf 2015](https://www.youtube.com/playlist?list=PLb0IAmt7-GS1cbw4qonlQztYV1TAW0sCr)
    seemed fun & quite interesting. For inspiration, watch [HYPE!][hype]
* [om](https://github.com/omcljs/om), because you actually want to use
    Clojure ;)

[hype]: https://www.youtube.com/watch?v=z5e7kWSHWTg

## Misc

- `npm config set prefix $HOME/.local --global`
