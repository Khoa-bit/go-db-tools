# go-db-tools

This is my personal tools for working with DB instances.
Because **I love zero dependency <3**

## build nested model

### Models struct:

```go
type Layer0 struct {
	ID      int64
	Layers1  []Layer1
	Layers12 []Layer2
}

type Layer1 struct {
	ID     int64
	Layers2 []Layer2
}

type Layer2 struct {
	ID int64
}
```

### Rows returned by DB one at a time:

#### Situation 1 - Many results at base level (aka. Level0 - Different IDs):

```sql
SELECT layer0.*, layer1.*, layer12.*, layer2.*
FROM layer0
LEFT JOIN layer1 ON ...
LEFT JOIN layer12 ON ...
LEFT JOIN layer2 ON ...
WHERE ...
LIMIT 1;
```

```go
type Row struct {
	Layer0  Layer0
	Layer1  Layer1
	Layer12 Layer2
	Layer2  Layer2
}

rows := []Row{
    {Layer0: Layer0{ID: 101}, Layer1: Layer1{ID: 201}, Layer12: Layer2{ID: 5301}, Layer2: Layer2{ID: 301}},
    {Layer0: Layer0{ID: 101}, Layer1: Layer1{ID: 201}, Layer12: Layer2{ID: 5302}, Layer2: Layer2{ID: 302}},
    {Layer0: Layer0{ID: 101}, Layer1: Layer1{ID: 201}, Layer12: Layer2{ID: 5303}, Layer2: Layer2{ID: 303}},
    {Layer0: Layer0{ID: 101}, Layer1: Layer1{ID: 204}, Layer12: Layer2{ID: 5304}, Layer2: Layer2{ID: 304}},
    {Layer0: Layer0{ID: 101}, Layer1: Layer1{ID: 205}, Layer12: Layer2{ID: 5305}, Layer2: Layer2{ID: 305}},
    {Layer0: Layer0{ID: 106}, Layer1: Layer1{ID: 206}, Layer12: Layer2{ID: 5306}, Layer2: Layer2{ID: 306}},
    {Layer0: Layer0{ID: 106}, Layer1: Layer1{ID: 206}, Layer12: Layer2{ID: 5307}, Layer2: Layer2{ID: 307}},
    {Layer0: Layer0{ID: 106}, Layer1: Layer1{ID: 206}, Layer12: Layer2{ID: 5308}, Layer2: Layer2{ID: 308}},
    {Layer0: Layer0{ID: 109}, Layer1: Layer1{ID: 209}, Layer12: Layer2{ID: 5000}, Layer2: Layer2{ID: 000}},
    {Layer0: Layer0{ID: 110}, Layer1: Layer1{ID: 000}, Layer12: Layer2{ID: 5000}, Layer2: Layer2{ID: 000}},
}
```

=> Results:

```json
{
  "ID": 101,
  "Layer1": [
    {
      "ID": 201,
      "Layer2": [
        {
          "ID": 301
        },
        {
          "ID": 302
        },
        {
          "ID": 303
        }
      ]
    },
    {
      "ID": 204,
      "Layer2": [
        {
          "ID": 304
        }
      ]
    },
    {
      "ID": 205,
      "Layer2": [
        {
          "ID": 305
        }
      ]
    },
    {
      "ID": 206,
      "Layer2": [
        {
          "ID": 306
        },
        {
          "ID": 307
        },
        {
          "ID": 308
        }
      ]
    },
    {
      "ID": 209,
      "Layer2": null
    }
  ],
  "Layer12": [
    {
      "ID": 5301
    },
    {
      "ID": 5302
    },
    {
      "ID": 5303
    },
    {
      "ID": 5304
    },
    {
      "ID": 5305
    },
    {
      "ID": 5306
    },
    {
      "ID": 5307
    },
    {
      "ID": 5308
    },
    {
      "ID": 5000
    }
  ]
}
```

#### Situation 2 - One result at base level (aka. Level0 - Same IDs):

```sql
SELECT layer0.*, layer1.*, layer12.*, layer2.*
FROM layer0
LEFT JOIN layer1 ON ...
LEFT JOIN layer12 ON ...
LEFT JOIN layer2 ON ...
WHERE ...;
```

```go
type Row struct {
	Layer0  Layer0
	Layer1  Layer1
	Layer12 Layer2
	Layer2  Layer2
}

rows := []Row{
    {Layer0: Layer0{ID: 101}, Layer1: Layer1{ID: 201}, Layer12: Layer2{ID: 5301}, Layer2: Layer2{ID: 301}},
    {Layer0: Layer0{ID: 101}, Layer1: Layer1{ID: 201}, Layer12: Layer2{ID: 5302}, Layer2: Layer2{ID: 302}},
    {Layer0: Layer0{ID: 101}, Layer1: Layer1{ID: 201}, Layer12: Layer2{ID: 5303}, Layer2: Layer2{ID: 303}},
    {Layer0: Layer0{ID: 101}, Layer1: Layer1{ID: 204}, Layer12: Layer2{ID: 5304}, Layer2: Layer2{ID: 304}},
    {Layer0: Layer0{ID: 101}, Layer1: Layer1{ID: 205}, Layer12: Layer2{ID: 5305}, Layer2: Layer2{ID: 305}},
    {Layer0: Layer0{ID: 101}, Layer1: Layer1{ID: 206}, Layer12: Layer2{ID: 5306}, Layer2: Layer2{ID: 306}},
    {Layer0: Layer0{ID: 101}, Layer1: Layer1{ID: 206}, Layer12: Layer2{ID: 5307}, Layer2: Layer2{ID: 307}},
    {Layer0: Layer0{ID: 101}, Layer1: Layer1{ID: 206}, Layer12: Layer2{ID: 5308}, Layer2: Layer2{ID: 308}},
    {Layer0: Layer0{ID: 101}, Layer1: Layer1{ID: 209}, Layer12: Layer2{ID: 5000}, Layer2: Layer2{ID: 000}},
    {Layer0: Layer0{ID: 101}, Layer1: Layer1{ID: 000}, Layer12: Layer2{ID: 5000}, Layer2: Layer2{ID: 000}},
}
```

=> Result:

```json
{
  "ID": 101,
  "Layer1": [
    {
      "ID": 201,
      "Layer2": [
        {
          "ID": 301
        },
        {
          "ID": 302
        },
        {
          "ID": 303
        }
      ]
    },
    {
      "ID": 204,
      "Layer2": [
        {
          "ID": 304
        }
      ]
    },
    {
      "ID": 205,
      "Layer2": [
        {
          "ID": 305
        }
      ]
    },
    {
      "ID": 206,
      "Layer2": [
        {
          "ID": 306
        },
        {
          "ID": 307
        },
        {
          "ID": 308
        }
      ]
    },
    {
      "ID": 209,
      "Layer2": null
    }
  ],
  "Layer12": [
    {
      "ID": 5301
    },
    {
      "ID": 5302
    },
    {
      "ID": 5303
    },
    {
      "ID": 5304
    },
    {
      "ID": 5305
    },
    {
      "ID": 5306
    },
    {
      "ID": 5307
    },
    {
      "ID": 5308
    },
    {
      "ID": 5000
    }
  ]
}
```

### Notes

- No matter what the DB driver or library we use, there should be an easy way to convert `query columns` into `a model struct`.
- Here is my example of mapping `stmt` from [zombiezen.com/go/sqlite](https://github.com/zombiezen/go-sqlite) into `a model struct`:

```go
import (
	"context"
	"zombiezen.com/go/sqlite"
	"zombiezen.com/go/sqlite/sqlitex"
)

type Model struct {
	ID        int64
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt time.Time
	DeletedBy string
}

func stmtModel(stmt *sqlite.Stmt, start int) Model {
	return Model{
		ID:        stmt.ColumnInt64(start),
		CreatedAt: stmtColumnTime(stmt, start+1),
		UpdatedAt: stmtColumnTime(stmt, start+2),
		DeletedAt: stmtColumnTime(stmt, start+3),
		DeletedBy: stmt.ColumnText(start + 4),
	}
}

// ---------------------

type Classroom struct {
	Model
	Field1 int64
	Field2 string
	Field3 string
}

func stmtClassroom(stmt *sqlite.Stmt, start int) (end int, classroom Classroom) {
	var buffer []byte
	stmt.ColumnBytes(1, buffer)

	return start + 9, Classroom{
		Model:    stmtModel(stmt, start),
		Field1:   stmt.ColumnInt64(start + 5),
		Field2:   stmt.ColumnText(start + 6),
		Field3:   stmt.ColumnText(start + 7),
	}
}

// ---------------------

type ClassroomWithStudents struct {
	Classroom
	Students []Student
}

// ---------------------

builder := NestedModelBuilder{}
err := sqlitex.Execute(conn, getClassroomWithStudents, &sqlitex.ExecOptions{
    Args: []any{userID, id},
    ResultFunc: func(stmt *sqlite.Stmt) error {
        next, classroom := stmtClassroom(stmt, 0)
        _, student := stmtStudent(stmt, next)
        builder.Build(&ClassroomWithStudents{Classroom: classroom}, &student)
        return nil
    },
})
results := GetAll[*Layer0](builder)
```

## build IN query:

```sql
SELECT *
FROM classroom
WHERE classroom.number IN (?1)
```

=> If there are 6 inputs, it will generate 10 slots, which a configuration step number to reduce the amount the prepared query

```sql
SELECT *
FROM classroom
WHERE classroom.number IN (?1, ?2, ?3, ?4, ?5, ?6, ?7, ?8, ?9, ?10)
```
