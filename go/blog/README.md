# A very simple (link)blog

Want a very simple blog?  Just run `blog blog.yaml` and be done.

See <https://papill0n.org/blog.html> for an example blog.

## Usage

    $ go get github.com/heyLu/lp/go/blog
    $ blog <yaml-file> > blog.html
    $ firefox blog.html
    # upload it somewhere!

## Features

Not many.

- post types (see [blog.yaml](./blog.yaml) for usage examples)
    - `shell`: Write about an interesting shell command
    - `link`: Quickly post a link
    - `image`: Post an image
    - `song`: Post a song (using `<audio>`)
    - `text`: Write an actual blog post (or just a heading)
    - `video`: Post a video (from YouTube)
- a simple tagging system
