package domain

//criação de pessoa ex: JSON {Id , Name e Age}
type Person struct {
	Id   int    `json:"id"`
	Name string `json:"name"`
	Age  int    `json:"age"`
}

type People struct{
	//Lista de pessoas
	People []Person `json:"people"`
}