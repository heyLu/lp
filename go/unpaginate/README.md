# unpaginate

`unpaginate` unpaginates JSON resources.

The requested resource is assumed to return an array of JSON
objects.  `unpaginate` prints a new array containing the JSON
objects from all pages on stdout.

Pagination is assumed to be in the format that the GitHub v3
API uses:

```http
HTTP/1.1 200 OK
...
Link: <https://api.github.com/user/527119/repos?per_page=42&page=2>; rel="next", <https://api.github.com/user/527119/repos?per_page=42&page=2>; rel="last"
```

## Usage

```
$ ./unpaginate https://api.github.com/users/heyLu/repos
...

$ ./unpaginate -user heyLu:$token https://api.github.com/user/repos
...
```

Note that the resources must return *arrays* of JSON objects.  Any other
kind of input will break `unpaginate`, sometimes in unexpected ways.

`unpaginate` was written to be used with [stars][], but is hopefully
also useful in other contexts.

[stars]: https://github.com/heyLu/lp/tree/master/go/stars
