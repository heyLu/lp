# stars

`stars` fetches your GitHub stars and updates them as necessary.  It is
not necessarily limited to just GitHub, as it will accept any git url
that's publicely accessible.

## Usage

`stars` accepts the input from various GitHub API endpoints that return
lists of repositories.

```
# Fetch starred repositories
$ curl https://api.github.com/users/heyLu/starred | ./stars
...

# Fetch public repos
$ curl https://api.github.com/users/heyLu/repos | ./stars
...

# Fetch all repos (public, private, org member)
$ token=<personal access token>
$ curl -u heyLu:$token https://api.github.com/user/repos | ./stars

# Fetch repos from a file
$ ./stars < repos.json
...
```

However, note that GitHub's API uses pagination, so ensure you set the
`per_page` query parameter appropriately, or use a tool such as
[unpaginate][] which collects the results from the pages into one JSON
document:

```
$ unpaginate https://api.github.com/users/heyLu/starred | ./stars
...
```

[unpaginate]: https://github.com/heyLu/lp/tree/master/go/unpaginate

## Input format

`stars` reads an array of JSON objects from stdin.  It must have
`full_name` and `git_url` fields, should have a `pushed_at` field and
also an optional `description` field.

For example:

```json
[
  {
    "full_name": "ProseMirror/prosemirror",
    "description": "The ProseMirror WYSIWYM editor",
    "pushed_at": "2015-08-06T21:37:16Z",
    "git_url": "git://github.com/ProseMirror/prosemirror.git",
  },
  ...
]
```
