package main

import (
	"errors"
	"github.com/gin-gonic/gin"
	"strconv"
)

// Un peu comme une interface en TypeScript
// En gros je donne un type à ma data, et après avec les strings literal, je peux préciser
// comment c'est censé être interprété en JSON
type todo struct {
	ID        int    `json:"id"`
	Item      string `json:"item"`
	Completed bool   `json:"completed"`
}

// Je déclare un truc, avec son type, ici un array de "todos"
var todos = []todo{
	{
		ID:        1,
		Item:      "Clean room",
		Completed: false,
	},
	{
		ID:        2,
		Item:      "Pet dog",
		Completed: false,
	},
	{
		ID:        3,
		Item:      "Murder neighbor",
		Completed: false,
	},
}

// Intended JSON renvoie un JSON Pretty Print
// JSON renvoie un JSON normal (mieux pour la prod)
func getTodos(context *gin.Context) {
	context.IndentedJSON(200, todos)
}

func getTodo(ctx *gin.Context) {
	// Ca peut retourner une erreur mais ballek
	// je convertis ma string en int, osef
	id, _ := strconv.Atoi(ctx.Param("id"))

	todo, err := getTodoById(id)
	if err != nil {
		ctx.IndentedJSON(400, gin.H{
			"message": "todo not found",
		})
		return
	}

	ctx.IndentedJSON(200, todo)
}

// Un helper
func getTodoById(id int) (*todo, error) {
	for i, todo := range todos {
		if todo.ID == id {
			// Je retourne une réference, je pourrai modifier le truc directement
			return &todos[i], nil
		}
	}

	return nil, errors.New("todo not found")
}

func addTodo(ctx *gin.Context) {
	var newTodo todo

	// Une syntaxe plus complexe de gestion d'erreur, je déclare ma variable
	// direct dans le if
	// Je fais un passage par référence pour que ça soit la variable (avec son type)
	// que j'ai déclaré avant qu'il reçoive la valeur
	if err := ctx.BindJSON(&newTodo); err != nil {
		return
	}

	todos = append(todos, newTodo)
	ctx.IndentedJSON(200, newTodo)
}

func toggleTodoStatus(ctx *gin.Context) {
	id, _ := strconv.Atoi(ctx.Param("id"))

	todo, err := getTodoById(id)
	if err != nil {
		ctx.IndentedJSON(400, gin.H{
			"message": "todo not found",
		})
		return
	}

	// Je change le status mec !
	todo.Completed = !todo.Completed
	ctx.IndentedJSON(200, todo)
}

// Et puis là, rien de particulier, c'est vraiment comme en NodeJS
func main() {
	router := gin.Default()
	router.GET("/todos", getTodos)
	router.GET("/todos/:id", getTodo) // Une route dynamique
	router.POST("/todos", addTodo)
	router.PATCH("/todos/:id", toggleTodoStatus)

	// Gestion d'erreurs, c'est assez étrange parce-que les fonctions ont des retours multiples
	// Donc ils peuvent envoyer soit rien soit un erreur, si l'erreur n'est pas renseignée, c'est que
	// rien n'a planté, y'a pas de système de throw comme en PHP
	err := router.Run("localhost:9090")
	if err != nil {
		return
	}
}
