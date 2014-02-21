# Introduction to Datomic

Datomic is a new database with a *flexible schema* that keeps a *history*
of all changes, allows you to assign arbitrary *attributes* to *entities*
and considers the database to be a *value*.

More buzzwords: ACID transactions, client-side query, elastic scaling,
pluggable storage, datalog for queries.

It's similar to git, in that it doesn't forget, but that's more or less where the
similarities end: Datomic has a linear history and cares about entities with
attributes, not about changes/arbitrary contents.

## Examples

* [Day of Datomic](https://github.com/Datomic/day-of-datomic): Code from a workshop
    from the creators of Datomic.
* [Music Brainz example](https://github.com/Datomic/mbrainz-sample)

## Resources

* [Learn Datalog Today!](http://learndatalogtoday.org)
* [The Architecture of Datomic](http://www.infoq.com/articles/Architecture-Datomic)
    describes the overall architecture of Datomic, with some details about the design
    and benefits of Datomic.
* [Videos about Datomic](http://www.datomic.com/videos.html), the "Tutorial Screencasts"
    section is a good place to start with Datomic basics.
* [Datomic Documentation](http://docs.datomic.com)
* [Data modelling in Datomic](http://stackoverflow.com/questions/10357778/data-modeling-in-datomic):
    How does Datomic differ from relational databases and what advantages
    does it offer. (You can change the schema and permanently delete old
    values now, though.)