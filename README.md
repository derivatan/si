# SI
SI (Secret Ingredient), is a database tool to handle query generation and relation handling. It can also be used to save models in the database.
The idea is to simplify the relation handling and to work only with model structs, instead of queries.

Its syntax is heavily based on the laravel eloquent, but very simplified and typed.


## Configuration

### Define models
```go
type artist struct {
    si.Model
    
    Name string `si:"DB_COLUMN_NAME"`
    Year int
}

func (a artist) GetModel() si.Model {
    return a.Model
}
```
All exported fields should have a matching column in the database with the naming as `snake_case(FieldName)`.
This can be overwritten with the si-tag (`DB_COLUMN_NAME` in the example above). The tag can be excluded.
The `si.Model` declared as the first field on the model, Similarly, the database must have the fields defined in this model as columns.


### Define Relationships

```go
type Album struct {
	database.Model

	Name string
	ArtistID uuid.UUID

	artist si.RelationData[Artist] `si:FIELD_NAME`
}

func (a Album) GetModel() si.Model {
    return a.Model
}

func (a Album) Artist() si.Relation[album, artist] {
	return si.BelongsTo[album, artist](a, "artist", func(a *album) *si.RelationData[artist] {
		return &a.artist
	})
}
```

A relationship is defined by a field and a function on the struct that defined the relation.

The field must have type si.RelationData with the generic type of the related struct.


### Init setup
SI needs a database connection and this is configured with `si.InitSecretIngredient(...)`.
This takes something that implements the `si.DB` interface, which can be a simple wrapper of some db connection. (`sqlwrapper.go`-file contians a wrapper for `sql.db`) 

All models and their 'database table name' need to be added with `si.AddModel[STRUCT](TABLE_NAME)`.

Example: 
```go
si.InitSecretIngredient(si.NewSQLDB(db))
si.AddModel[contact]("contacts")
si.AddModel[artist]("artists")
si.AddModel[album]("albums")
```


## Usage

* `si.Query[...]().Get()` is the main entry point for retrieving data from the database.

* `si.Save(MODEL)` is used to create or update a model, with the values upon the model.
  To save relations, you must update the ID column, just as a normal column. This will **not** change what's stored in relation field if it is already loaded. 

* If you want to debug the generated queries, or do whatever with them, you can use `si.SetLogger(...)`.
  It takes an anonymous function as argument, and this function will be called with the query and its parameters before it is executed.

  This example will print all queries.
    ```go
    si.SetLogger(func(a ...any) {
        fmt.Println(a...)
    })
    ```


## Example

There is a small example in the `examples` folder that can be run with:
```
go run examples/*.go
```

In order to run this You'll need to
* Uncomment the postgres driver importer in `examples/example.go`. (its only used for the example, and not the library itself.)
* Update the db-connection details in `examples/example.go`. 
* Run the `examples/db.sql` script.


### Comments

 * There is no mapping for a many-to-many relation yet. In order to achieve this anyway, you can make a model for the pivot-table and use the one-to-many relation in both directions to the other models.
