package main

import (
	"database/sql"
	"errors"
	"log"
	"net/http"
	"os"

	"github.com/graphql-go/graphql"
	"github.com/graphql-go/handler"
	_ "github.com/lib/pq"
)

type Celebrity struct {
	Id           int    `json:"id"`
	Fullname     string `json:"fullname"`
	Nationality  string `json:"nationality"`
	ReqPhotoPath string `json:"reqPhotoPath"`
}

var db *sql.DB

var celebrityType = graphql.NewObject(graphql.ObjectConfig{
	Name: "Celebrity",
	Fields: graphql.Fields{
		"id":           &graphql.Field{Type: graphql.Int},
		"fullname":     &graphql.Field{Type: graphql.String},
		"nationality":  &graphql.Field{Type: graphql.String},
		"reqPhotoPath": &graphql.Field{Type: graphql.String},
	},
})

var celebrityInputType = graphql.NewInputObject(graphql.InputObjectConfig{
	Name: "CelebrityInput",
	Fields: graphql.InputObjectConfigFieldMap{
		"id":           &graphql.InputObjectFieldConfig{Type: graphql.NewNonNull(graphql.Int)},
		"fullname":     &graphql.InputObjectFieldConfig{Type: graphql.NewNonNull(graphql.String)},
		"nationality":  &graphql.InputObjectFieldConfig{Type: graphql.NewNonNull(graphql.String)},
		"reqPhotoPath": &graphql.InputObjectFieldConfig{Type: graphql.NewNonNull(graphql.String)},
	},
})

var celebrityUpdateInputType = graphql.NewInputObject(graphql.InputObjectConfig{
	Name: "CelebrityUpdateInput",
	Fields: graphql.InputObjectConfigFieldMap{
		"fullname":     &graphql.InputObjectFieldConfig{Type: graphql.NewNonNull(graphql.String)},
		"nationality":  &graphql.InputObjectFieldConfig{Type: graphql.NewNonNull(graphql.String)},
		"reqPhotoPath": &graphql.InputObjectFieldConfig{Type: graphql.NewNonNull(graphql.String)},
	},
})

func initDB() {
	var err error

	connStr := os.Getenv("DATABASE_URL")
	if connStr == "" {
		connStr = "postgres://postgres:12345678@localhost:5432/Celebrities?sslmode=disable" // используем postgresql
		log.Println("Используется строка подключения по умолчанию. Установите DATABASE_URL для production среды")
	}

	db, err = sql.Open("postgres", connStr) // используем postgresql
	if err != nil {
		log.Fatal("Не удалось подключиться к базе данных:", err)
	}

	err = db.Ping()
	if err != nil {
		log.Fatal("Не удалось выполнить ping базы данных:", err)
	}

	// используем postgresql
	createTableSQL := ` 
	CREATE TABLE IF NOT EXISTS Celebrities (
		id INTEGER PRIMARY KEY,
		fullname TEXT NOT NULL,
		nationality TEXT NOT NULL,
		reqPhotoPath TEXT NOT NULL
	);`

	_, err = db.Exec(createTableSQL)
	if err != nil {
		log.Fatal("Не удалось создать таблицу:", err)
	}

	log.Println("База данных PostgreSQL успешно инициализирована") // используем postgresql
}

func scanCelebrity(rows *sql.Rows) (Celebrity, error) {
	var c Celebrity
	err := rows.Scan(&c.Id, &c.Fullname, &c.Nationality, &c.ReqPhotoPath)
	return c, err
}

func celebrityToMap(c Celebrity) map[string]interface{} {
	return map[string]interface{}{
		"id":           c.Id,
		"fullname":     c.Fullname,
		"nationality":  c.Nationality,
		"reqPhotoPath": c.ReqPhotoPath,
	}
}

func getAllCelebrities() ([]map[string]interface{}, error) {
	rows, err := db.Query("SELECT id, fullname, nationality, reqPhotoPath FROM Celebrities ORDER BY id")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var result []map[string]interface{}
	for rows.Next() {
		c, err := scanCelebrity(rows)
		if err != nil {
			return nil, err
		}
		result = append(result, celebrityToMap(c))
	}
	return result, rows.Err()
}

func getCelebrityByID(id int) (map[string]interface{}, error) {
	var c Celebrity
	err := db.QueryRow("SELECT id, fullname, nationality, reqPhotoPath FROM Celebrities WHERE id = $1", id).
		Scan(&c.Id, &c.Fullname, &c.Nationality, &c.ReqPhotoPath)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return celebrityToMap(c), nil
}

func buildSchema() graphql.Schema {
	rootQuery := graphql.NewObject(graphql.ObjectConfig{
		Name: "Query",
		Fields: graphql.Fields{
			"celebrities": &graphql.Field{
				Type: graphql.NewList(celebrityType),
				Resolve: func(p graphql.ResolveParams) (interface{}, error) {
					list, err := getAllCelebrities()
					if err != nil {
						log.Printf("GraphQL celebrities: %v", err)
						return nil, err
					}
					log.Printf("Query celebrities - возвращено %d знаменитостей", len(list))
					return list, nil
				},
			},
			"celebrity": &graphql.Field{
				Type: celebrityType,
				Args: graphql.FieldConfigArgument{
					"id": &graphql.ArgumentConfig{Type: graphql.NewNonNull(graphql.Int)},
				},
				Resolve: func(p graphql.ResolveParams) (interface{}, error) {
					id := p.Args["id"].(int)
					item, err := getCelebrityByID(id)
					if err != nil {
						log.Printf("GraphQL celebrity(%d): %v", id, err)
						return nil, err
					}
					if item == nil {
						log.Printf("Query celebrity(%d) - не найдено", id)
						return nil, nil
					}
					log.Printf("Query celebrity(%d) - найдено", id)
					return item, nil
				},
			},
		},
	})

	rootMutation := graphql.NewObject(graphql.ObjectConfig{
		Name: "Mutation",
		Fields: graphql.Fields{
			"createCelebrity": &graphql.Field{
				Type: celebrityType,
				Args: graphql.FieldConfigArgument{
					"input": &graphql.ArgumentConfig{Type: graphql.NewNonNull(celebrityInputType)},
				},
				Resolve: func(p graphql.ResolveParams) (interface{}, error) {
					input := p.Args["input"].(map[string]interface{})
					id := input["id"].(int)
					fullname := input["fullname"].(string)
					nationality := input["nationality"].(string)
					reqPhotoPath := input["reqPhotoPath"].(string)

					var exists bool
					err := db.QueryRow("SELECT EXISTS(SELECT 1 FROM Celebrities WHERE id = $1)", id).Scan(&exists)
					if err != nil {
						log.Printf("Mutation createCelebrity: %v", err)
						return nil, err
					}
					if exists {
						return nil, errors.New("Знаменитость с таким ID уже существует")
					}

					_, err = db.Exec(
						"INSERT INTO Celebrities (id, fullname, nationality, reqPhotoPath) VALUES ($1, $2, $3, $4)",
						id, fullname, nationality, reqPhotoPath,
					)
					if err != nil {
						log.Printf("Mutation createCelebrity insert: %v", err)
						return nil, err
					}

					log.Printf("Mutation createCelebrity - добавлена знаменитость с ID %d", id)
					return celebrityToMap(Celebrity{Id: id, Fullname: fullname, Nationality: nationality, ReqPhotoPath: reqPhotoPath}), nil
				},
			},
			"updateCelebrity": &graphql.Field{
				Type: celebrityType,
				Args: graphql.FieldConfigArgument{
					"id":    &graphql.ArgumentConfig{Type: graphql.NewNonNull(graphql.Int)},
					"input": &graphql.ArgumentConfig{Type: graphql.NewNonNull(celebrityUpdateInputType)},
				},
				Resolve: func(p graphql.ResolveParams) (interface{}, error) {
					id := p.Args["id"].(int)
					input := p.Args["input"].(map[string]interface{})
					fullname := input["fullname"].(string)
					nationality := input["nationality"].(string)
					reqPhotoPath := input["reqPhotoPath"].(string)

					var exists bool
					err := db.QueryRow("SELECT EXISTS(SELECT 1 FROM Celebrities WHERE id = $1)", id).Scan(&exists)
					if err != nil {
						log.Printf("Mutation updateCelebrity: %v", err)
						return nil, err
					}
					if !exists {
						return nil, errors.New("Знаменитость не найдена")
					}

					_, err = db.Exec(
						"UPDATE Celebrities SET fullname = $1, nationality = $2, reqPhotoPath = $3 WHERE id = $4",
						fullname, nationality, reqPhotoPath, id,
					)
					if err != nil {
						log.Printf("Mutation updateCelebrity update: %v", err)
						return nil, err
					}

					log.Printf("Mutation updateCelebrity(%d) - обновлено", id)
					return celebrityToMap(Celebrity{Id: id, Fullname: fullname, Nationality: nationality, ReqPhotoPath: reqPhotoPath}), nil
				},
			},
			"deleteCelebrity": &graphql.Field{
				Type: graphql.String,
				Args: graphql.FieldConfigArgument{
					"id": &graphql.ArgumentConfig{Type: graphql.NewNonNull(graphql.Int)},
				},
				Resolve: func(p graphql.ResolveParams) (interface{}, error) {
					id := p.Args["id"].(int)

					var exists bool
					err := db.QueryRow("SELECT EXISTS(SELECT 1 FROM Celebrities WHERE id = $1)", id).Scan(&exists)
					if err != nil {
						log.Printf("Mutation deleteCelebrity: %v", err)
						return nil, err
					}
					if !exists {
						return nil, errors.New("Знаменитость не найдена")
					}

					_, err = db.Exec("DELETE FROM Celebrities WHERE id = $1", id)
					if err != nil {
						log.Printf("Mutation deleteCelebrity delete: %v", err)
						return nil, err
					}

					log.Printf("Mutation deleteCelebrity(%d) - удалено", id)
					return "Знаменитость успешно удалена", nil
				},
			},
		},
	})

	schema, err := graphql.NewSchema(graphql.SchemaConfig{
		Query:    rootQuery,
		Mutation: rootMutation,
	})
	if err != nil {
		log.Fatal("Не удалось создать GraphQL схему:", err)
	}
	return schema
}

func main() {
	initDB()
	defer db.Close()

	schema := buildSchema()

	h := handler.New(&handler.Config{
		Schema: &schema,
	})

	http.Handle("/graphql", h)

	log.Println("GraphQL сервер запущен на порту 3000")
	log.Println("Endpoint: POST http://localhost:3000/graphql")
	log.Fatal(http.ListenAndServe(":3000", nil))
}
