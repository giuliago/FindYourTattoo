package main

import("net/http"
       "github.com/globalsign/mgo/bson"
	   "log"
	"fmt"
	"gopkg.in/go-playground/validator.v9")

	type User struct {
		Name     string `validate:"required"`
		Username string `validate:"required"`
		Email    string `validate:"required,email"`
		Password string `validate:"required,min=8"`
	}

type UserUpdate struct {

	Name      string `json:"name" validate:"required"`
	Birthdate string `json:"birthdate"`
	Gender    string `json:"gender"` 
}  

func RegisterUser(response http.ResponseWriter, request *http.Request) {
	// ---------------------------
	// PADRONIZAÇÃO DA INFORMAÇÃO
	// ---------------------------
	// retirando corpo da requisição 
	name := request.FormValue("name")
	username := request.FormValue("username")
	//address := request.FormValue("address")
	//typePerson := request.FormValue("type")
	//state := request.FormValue("state")
	//city := request.FormValue("city")
	//telefone := request.FormValue("telefone")
	email := request.FormValue("email")
	password := request.FormValue("password")


	// criando estrutura de usuário vazia
	user := User{
		Name: name,
		Username: username,
		Email: email,
		Password: password,
	}

	// decodificando o JSON no estrutura de usuário
	err := DecodeJson(body, &user)

	// validando erro de decodificação
	if err != nil {
		// resposta HTTP de formato inválido da informação
		response.Write([]byte(`{"code": 400, "message": "you know nothing!"}`))
		return
	}*/

	// ---------------------------
	// VALIDAÇÃO DA INFORMAÇÃO
	// ---------------------------

	// validação da informação 
	err = validate.Struct(user)

	// validando erro de validação
	if err != nil {
		// resposta HTTP de invalidação dos campos da requisição
		response.Write([]byte(`{"code": 400, "message": "invalid parameters!"}`+ errormessage))
		return
	}*/

	// -----------------------------------
	// REGRAS DE NEGÓCIO E BANCO DE DADOS
	// -----------------------------------

	// regras de negócio para salvar usuários no banco de dados
	SaveUser(user)

	// resposta de sucesso HTTP

}

func SaveUser(user User) (err error) {

	conn, err := GetConnection()

	collection := conn.DB("tattoo").C("users")

	collection.Insert(user)

	return
}

func SearchUser(response http.ResponseWriter, request *http.Request) {

	conn, _ := GetConnection()

	collection := conn.DB("tattoo").C("users")

	username := request.FormValue("username")
	password := request.FormValue("password")

	login := User {
		Username: username,
		Password: password,
	}

	  err := collection.Find(bson.M{"username": username}).One(&login)
	  err2 := collection.Find(bson.M{"password": password}).One(&login)
  	if(err != nil && err2 != nil) {
		fmt.Printf("Passou")
    	return 
	  }
		fmt.Printf("credenciais erradas.")
}

func UpdateUser(response http.ResponseWriter, request *http.Request) {

	// ---------------------------
	// PADRONIZAÇÃO DA INFORMAÇÃO
	// ---------------------------

	// retirando corpo da requisição 
	body := request.Body

	// criando estrutura de usuário vazia
	user := UserUpdate{}

	// decodificando o JSON no estrutura de usuário
	err := DecodeJson(body, &user)

	// validando erro de decodificação
	if err != nil {
		// resposta HTTP de formato inválido da informação
		response.Write([]byte(`{"code": 400, "message": "you know nothing!"}`))
		return
	}

	// ---------------------------
	// VALIDAÇÃO DA INFORMAÇÃO
	// ---------------------------

	// validação da informação 
	err = validate.Struct(user)

	// validando erro de validação
	if err != nil {
		// resposta HTTP de invalidação dos campos da requisição
		response.Write([]byte(`{"code": 400, "message": "invalid parameters!"}`+ errormessage))
		return
	}

	// -----------------------------------
	// REGRAS DE NEGÓCIO E BANCO DE DADOS
	// -----------------------------------

	// regras de negócio para salvar usuários no banco de dados
	err = EditUser(user)

	// validação na interação com banco de dados
	if err != nil{
		// resposta HTTP de falha interna do servidor
		response.Write([]byte("500"))
		return
	}

	// resposta de sucesso HTTP
	response.Write([]byte("200"))

}

func EditUser(user UserUpdate) (err error) {

	conn, err := GetConnection()

	if err != nil {
		log.Printf("[ERROR] could not get BD connection")
		return
	}

	collection := conn.DB("webtatu").C("users")

	err = collection.Update(bson.M{"name": user.Name}, user)

	if err != nil {
		log.Printf("[ERROR] could not edit user")
		return
	}

	return
}