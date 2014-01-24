using HTTPClient
using JSON
using DataFrames
using Gadfly

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
