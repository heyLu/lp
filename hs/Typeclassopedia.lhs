# Typeclassopedia: A journey

I don't think Haskell is scary, but there still is a lot to learn after
you finish LYAH. I kind of stopped there, because I was still scared and
because it felt so good already, why continue?

Well, because Haskell is so much more and I think it's one of the things
that make Haskell special: Working with high-level abstractions that
make things possible that were previously unthinkable.

I think I've missed a lot so far, because if you "just use" Haskell you
aren't really in control, one might even say that you aren't really
writing Haskell, you've stopped right after the doorstep.

So, let's find out what's behind the doorstep.

## Functors

> class Functor f where
>   fmap :: (a -> b) -> f a -> f b

Instances for `Either e` and `((->) e)`.

Mhh, `Either e` seems easier, although it also looks a bit weird:
`Either` is usually about two things, the `l` and the `r` part.

Ok, so why are we declaring `Either e` the `Functor` instance, why not
just `Either`? (Excuse my ignorance here, I've really been avoiding this
stuff, so I'll have to be very explicit about everything for a while.)

Well, let's have a look at the `fmap` again, then. It's type is `(a ->
b) -> f a -> f b`, so for `Either e` it is `(a -> b) -> Either e a ->
Either e b`. That's a hint why it is this way: the `l` (here `e`) part
is usually used for errors and `fmap`ing over that is probably not very
interesting most of the time. So we "change" the `r` part.

But still, *can* we define a `Functor` instance for `Either`? The type
would be `(a -> b) -> Either a -> Either b`. I think that would not
work, because we haven't "applied" `Either` fully, so likely Haskell
will not like it.

But what does `Either l` mean? It's incomplete, but we if we "apply"
some other type to it then we have a type we can use. Is this what kinds
are? Some kind of marker for such incomplete types? Is there a name for
those things?

Playing with this shows that those things really aren't allowed:

```
λ :kind Either
* -> * -> Either
λ :kind Either Int
* -> *
λ Either Int -> Either Int
<interactive>:1:1:
    Expecting one more argument to `Either Int'
    In a type in a GHCi command: Either Int -> Either Int
```

First, `Either` is a *type constructor*, e.g. you can't use it right
away, but you build types with it. You can, however, also apply it
partially and only later decide on the other parameters.

For example, we can have `WithError`, which is just `Either String`.
But to use it we still have to "apply" it to some other type.

> type WithError = Either String -- "incomplete"
>
> zeroIsSpecial :: Int -> WithError Int
> zeroIsSpecial 0 = Left "you can't have my precious zero, sorry"
> zeroIsSpecial n = Right n

The [page on Kinds][kinds] on the Haskell wiki tells me that those
complete types, e.g. types with kind `*` are called monomorphic (data)
types. Are there polymorphic data types? Well, there surely are, but we
can't compute with them. (Surely? Yes, `Either`, for example. But what
about ordinary constructors with two arguments? Mhh.. ok, haskell says
they aren't type constructors, so they have no kind.)

[kinds]: http://www.haskell.org/haskellwiki/Kind

Ok, let's come back to this later.

So, let's define a `Functor` instance for `Either e`. But first, what
will it "mean"? If we make something a Functor, then we look at it as a
whole, so `Either e` as an instance of `Functor` is not two parts, but
one. It's a container for things. And to the value in that container we
want to apply a function. Let's do that first.

> instance Functor (Either e) where
>   fmap f (Right r) = Right $ f r
>   fmap f l = l

That was... easy. Am I missing something?

What about `((->) e)` then? That seems trickier. First, `(->)`s kind is
`* -> * -> *`, just as `Either`s was, but I find it a little bit more
difficult to see what that means.

For `Either` it was that we had to supply the types for the `Left` and
`Right` parts. And for `(->)` it means which arguments we want the
function type to have. So, `((->) e)` means something with an `e` as an
argument.

The Typeclassopedia mentions that the `e` is read-only and that this
type is also known as `Reader`, so let's look at what it does.

> instance Functor ((->) e) where
>   fmap :: (a -> b) -> (((->) e) a) -> (((->) e) b)
>   -- aka  (a -> b) -> (e -> a) -> (e -> b)
>   fmap f g = 
