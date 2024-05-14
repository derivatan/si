# SI
SI (Secret Ingredient), is a database tool to handle query generation and relation handling. It can also be used to save models in the database.
The idea is to simplify the relation handling and to work only with model structs, instead of queries.

Its syntax is heavily based on the laravel eloquent, but very simplified and typed.


## Configuration

### Defining a simple model
```go
type artist struct {
    // Note: `si.Model` must be the first declared thing on the struct.
    si.Model
    
    Name string `si:"DB_COLUMN_NAME"`
}

func (a artist) GetModel() si.Model {
    return a.Model
}

func (a artist) GetTable() string {
    return "artists"
}
```
A model must implement the [`Modeler`](https://github.com/derivatan/si/blob/095f3ca8e974635a8ac20e8b2e327af27556c781/common.go#L20C2-L20C2) interface.

It must also embed the `si.Model` as the first declared field on the model.

All exported fields on the model should have a matching column in the database with the naming as `snake_case(FieldName)`.
This can be overwritten with the si-tag (`DB_COLUMN_NAME` in the example above). The tag can be excluded.


### Define Relationships

```go
type Album struct {
    si.Model

    Name string
    Year int
    ArtistID uuid.UUID

    // Note: must be unexported.
    artist si.RelationData[Artist] `si:FIELD_NAME`
}

func (a Album) GetModel() si.Model {
    return a.Model
}

func (a Album) GetTable() string {
    return "albums"
}

func (a Album) Artist() si.Relation[Album, Artist] {
    return si.BelongsTo[Album, Artist](a, "artist", func(a *Album) *si.RelationData[Artist] {
        return &a.artist
    })
}
```

A relationship is defined by two things.
* An **unexported** field with the type `si.RelationData[To]`
* An exported function that returns a `si.Relation[From, To]`

The field is only for _si_:s internal use and should not be used or modified in any way. To get a relation you must use the function as a query builder.

To ignore a field on a model, the tag ``si:"-"`` can be used or make it unexported.
```go
type Artist struct {
  si.Model
  
  Name string
  IgnoredField `si:"-"` // Ignored because of the tag.
  ignoredField // Ignored because of not exported.
}
```


### Database setup
In order to be completely agnostic about the database, _si_ uses these [interfaces](https://github.com/derivatan/si/blob/main/db.go) for database communication.
This is based on the `sql.DB`, but can easily be implemented with whatever you want to use. A simple example of such an implementaiton can be found in [`sql_wrap.go`](https://github.com/derivatan/si/blob/main/sql_wrap.go)


## Usage

* `si.Query[T]()` is the main entry point for retrieving data from the database.

  Examples
```go
// Get alla albums that start with the letter 'a'.
albums, err := si.Query[Album]().Where("name", "ILIKE", "a%").OrderBy("name", true).Get(db)
// Get the Artist from an Album.
artist, err := albums[0].Artist().Find(db)

// Get all Artist, with all their albums. This will only execute two queries.
artists, err := si.Query[Artist]().With(func(a Artist, r []Artist) error {
    return a.Albums().Execute(db, r)
}).Get(db)
// Get the albums from an Artist, since irs already fetched from the database, it does not require a `db`, and there can be no error.
albums := artists[0].Albums().MustFind(nil)
```

* `si.Save(model)` is used to create or update a model, with the values upon the model.
  To save relations, you must update the ID column, just as a normal column. This will **not** change what's stored in relation field if it is already loaded. 

* If you need to debug the generated queries, or get some silent errors, you can use `si.SetLogger(...)`.
  This logger will be called with all the queries and their arguments that _si_ generates, and might in some cases give some debugging messages. 

  This example will print all queries.
```go
si.SetLogger(func(a ...any) {
    fmt.Println(a...)
})
```


## Example and tests

There are integration tests for all major functionalities in a [separate repo](http://github.com/derivatan/si_test)
The tests are put there, in its own repository, because I don't want the library itself to import packages that are only needed for the testing.

These are also a good example for how to use the library.




### Comments

 * There is no mapping for a many-to-many relation yet. In order to achieve this anyway, you can make a model for the pivot-table and use the one-to-many relation in both directions to the other models.
