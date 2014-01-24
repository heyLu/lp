using HTTPClient
using JSON
using DataFrames
using Gadfly

# this started out as an experiment, just trying to see what i can do
# with julia and some json.
#
# this will generate a graph of the times repositories of a user were
# created vs the last push times to that repository.
#
# it started out as a silly experiment that simply plotted two
# numeric/ordered values from the data, but it's quite interesting:
# quite a few users have this pattern where they constantly create (or
# fork) repositories and then move on after a while. for most people
# and projects we see a roughly linear increase, suggesting that people
# continue creating repositories on github and work a bit on them and
# then move on.
# still, there are projects that were started early and are still active,
# for example mbostock/d3 or my own heyLu/confidence (my dotfiles)
#
# next up: having fun with commits? (or activity?)

github_get(path) = takebuf_string(get(string("https://api.github.com", path), headers=[("User-Agent", "Julia-Experiment")]).body)

function toDataFrame(dicts, keys = keys(dicts[1]))
    df = DataFrame()
    for k in keys
        df[k] = map(d -> d[k], dicts)
    end
    return df
end

repos_str = github_get("/users/mbostock/repos")
repos = JSON.parse(repos_str)
repos_df = toDataFrame(repos, filter(k -> search(k, "url") == 0:-1, keys(repos[1])))

p = plot(toDataFrame(repos), y = "pushed_at", x = "created_at", label = "name", Geom.point, Geom.label(;hide_overlaps=false))
draw(PNG("created_vs_push.png", 9inch, 9inch/golden), p)
