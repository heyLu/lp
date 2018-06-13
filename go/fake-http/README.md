# fake-http

`fake-http` pretends to be other HTTP APIs.

## TODO

- replay from log instantly (i.e. operate like a request cache)
    - ensure every request is cached max once?  at least in the response
      defn?
    - make it possible to download the cache
        - /_cache, analogous to /_log?
        - or reuse /_log, which wouldn't need new code, but would make
          it impossible to download the log.  damn...
            - how about /_log?cache=true?
            - or make /_log redirect to /_cache if `-cache` is enabled
- /_clear endpoint
- some kind of templating (request info + random functions?)
- lua integration? ;)

## Example: Record and replay kubernetes traffic

This needs a running Minikube instance, which can be started with
`minikube start`.

```
# shortcut for the cmdline below: fake-http -proxy-minikube
$ fake-http -proxy-url=https://$(minikube ip):8443 -proxy-client-cert ~/.minikube/client.crt -proxy-client-key ~/.minikube/client.key
2018/06/09 13:21:19 Listening on http://localhost:8080
2018/06/09 13:21:19 See http://localhost:8080/_help
...

# switch to another terminal

# set up "fake-minikube" context
$ kubectl set-cluster fake-minikube --server=http://localhost:8080
$ kubectl set-context fake-minikube --cluster=fake-minikube --user minikube

$ kubectl --context=fake-minikube get pods
NAME                         READY     STATUS    RESTARTS   AGE
hellogo-2387138299-p83qw     1/1       Running   2          266d
hellonode-1839943766-t2hsv   1/1       Running   2          266d
```

Now have a look at the requests kubectl made: <http://localhost:8080>.

Or, save the log and replay it later:

```
$ curl -H 'Accept: application/yaml' http://localhost:8080/_log > minikube.yaml

# replay later
$ fake-http minikube.yaml
```
