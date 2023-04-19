package person

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"

	"github.com/Computerick/api-crud-persistencia-arquivo/domain"
)

type Service struct {
	// Recebe o caminho completo do arquivo (parâmetro)
	dbFilePath string
	people     domain.People
}

func NewService(dbFilePath string) (Service, error) {
	//Verificar os cenários do arquivo recebido
	_, err := os.Stat((dbFilePath))
	if err != nil {
		if os.IsNotExist(err) {
			//cria um arquivo vazio
			err = createEmptyFile(dbFilePath)
			if err != nil {
				return Service{}, err
			}
			return Service{
				dbFilePath: dbFilePath,
				people:     domain.People{},
			}, nil

		}
	}
	// Se existir, leio o arquivo atualizo a variavel people do serviço com as pessoas do arquivo
	// verifica se existe o arquivo
	jsonFile, err := os.Open(dbFilePath)
	if err != nil {
		return Service{}, fmt.Errorf("Error trying to open file that contains al people: %s", err.Error())
	}
	//Ler tudo por essa variável
	jsonFileContentByte, err := ioutil.ReadAll(jsonFile)
	if err != nil {
		return Service{}, fmt.Errorf("Error trying to read file: %s", err.Error())
	}

	var allPeople domain.People
	json.Unmarshal(jsonFileContentByte, &allPeople)

	return Service{
		dbFilePath: dbFilePath,
		people:     allPeople,
	}, nil
}

func createEmptyFile(dbFilePath string) error {
	var people domain.People = domain.People{
		People: []domain.Person{},
	}
	peopleJSON, err := json.Marshal(people)
	if err != nil {
		return fmt.Errorf("Error trying to encode people as JSON?: %s", err.Error())
	}
	// Passo entre parenteses o caminho do arquivo
	err = ioutil.WriteFile(dbFilePath, peopleJSON, 0755)
	if err != nil {
		return fmt.Errorf("Error trying to write to file. Error: %s", err.Error())
	}
	return nil
}

func (s *Service) Create(person domain.Person) error {
	//verificar se pessoa existe, se ja existe retorna erro
	if s.exists(person) {
		return fmt.Errorf("Error trying to create person. There is a person with this Id already registered")
	}

	// adicionar a pessoa ao slice(lista) de pessoas
	s.people.People = append(s.people.People, person)

	//salvo o arquivo
	err := s.saveFile()
	if err != nil {
		return fmt.Errorf("Error trying save file in method created. Error: %s", err.Error())
	}
	return nil
}

// Metodo verifica se pessoa existe
func (s Service) exists(person domain.Person) bool {
	for _, currentPerson := range s.people.People {
		if currentPerson.Id == person.Id {
			return true
		}
	}
	return false
}

// Metodo para salvar registro
func (s Service) saveFile() error {
	allPeopleJSON, err := json.Marshal(s.people)
	if err != nil {
		return fmt.Errorf("Error trying to encode people as json: %s", err.Error())
	}
	return ioutil.WriteFile(s.dbFilePath, allPeopleJSON, 0755)
}

//Metodo para consultar registros
func (s Service ) List() domain.People{
	return s.people
}

//Metodo para consultar registros por Id
func (s Service) GetbyId(personId int) (domain.Person, error){
	for _, currentPerson := range s.people.People {
		if currentPerson.Id == personId {
			return currentPerson, nil
		}
	}
	return domain.Person{}, fmt.Errorf("Person not found")
}

//Metodo para atualizar o registro
func (s *Service) Update(person domain.Person) error {
	// Procurando o Id na lista
	var indexToUpdate int = -1
	for index, currentPerson := range s.people.People {
		if currentPerson.Id == person.Id {
			indexToUpdate = index
			break
		}
	}
	if indexToUpdate < 0{
		return fmt.Errorf("There is no person with the given Id in our database")
	}
	//atualizar e salva novo arquivo
	s.people.People[indexToUpdate] = person
	return s.saveFile()
}

//Metodo para deletar registro por Id (procurar o registro se existente, deleter e salvar mudança no banco)
// * O Uso do Ponteiro é necessário para Atualizar e Deletar sem defasar da lista principal!
func (s *Service) DeleteById(personId int)error{
	// Procurando o Id na lista
	var indexToDelete int = -1
	for index, currentPerson := range s.people.People {
		if currentPerson.Id == personId {
			indexToDelete = index
			break
		}
	}
	if indexToDelete < 0{
		return fmt.Errorf("There is no person with the given Id in our database")
	}
	//[1 2 X3X 4] == [1 2 4]
	//cria a nova slice e salva novo arquivo
	s.people.People = append(s.people.People[:indexToDelete], s.people.People[indexToDelete+1:]...)
	return s.saveFile()

}