package utils

import (
	"encoding/json"
	"os"
	"path/filepath"

	"master-alias.com/core/structs"
)

func GetConfigPath(filename string) string {
	home, _ := os.UserHomeDir()
	return filepath.Join(home, ".master-alias", filename)
}

func CreateAliasJsonFile() error {
	home, _ := os.UserHomeDir()
	filePath := filepath.Join(home, ".master-alias", "alias.json")

	// Verifica se o arquivo já existe
	if _, err := os.Stat(filePath); err == nil {
		// Arquivo existe, não faz nada
		return nil
	} else if !os.IsNotExist(err) {
		// Outro erro ao tentar verificar
		return err
	}

	f, err := os.Create(filepath.Join(home, ".master-alias", "alias.json"))

	if err != nil {
		return err
	}
	defer f.Close()

	_, err = f.Write([]byte("[]"))
	if err != nil {
		return err
	}

	return nil
}

func WriteJSON(filename string, alias []structs.Alias) error {
	CreateAliasJsonFile()
	file, err := os.Create(GetConfigPath(filename))
	if err != nil {
		return err
	}
	defer file.Close()

	encoder := json.NewEncoder(file)
	encoder.SetIndent("", "  ")
	return encoder.Encode(alias)
}

func ReadJSON(filename string) ([]structs.Alias, error) {
	CreateAliasJsonFile()
	file, err := os.Open(GetConfigPath(filename))
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var alias []structs.Alias
	decoder := json.NewDecoder(file)
	err = decoder.Decode(&alias)
	return alias, err
}

func RemoveItem(filename, id string) {
	alias := FindById(filename, id)
	if alias.Id == "" {
		return
	}

	aliases, err := ReadJSON(filename)

	if err != nil {
		panic(err)
	}

	newAliases := []structs.Alias{}

	for _, a := range aliases {
		if a.Id != alias.Id {
			newAliases = append(newAliases, a)
		}
	}

	err = WriteJSON(filename, newAliases)
}

func FindById(filename, id string) structs.Alias {
	CreateAliasJsonFile()
	// Lê o arquivo JSON
	data, err := os.ReadFile(GetConfigPath(filename))
	if err != nil {
		panic(err)
	}

	// Converte para slice de structs
	var aliases []structs.Alias
	if err := json.Unmarshal(data, &aliases); err != nil {
		panic(err)
	}

	// Procura o alias pelo nome
	for _, a := range aliases {
		if a.Id == id {
			return structs.Alias{Id: id, Name: a.Name, Command: a.Command}
		}
	}

	return structs.Alias{Id: "", Name: "", Command: ""}
}

func FindByName(filename, name string) structs.Alias {
	CreateAliasJsonFile()

	// Lê o arquivo JSON
	data, err := os.ReadFile(GetConfigPath(filename))
	if err != nil {
		panic(err)
	}

	// Converte para slice de structs
	var aliases []structs.Alias
	if err := json.Unmarshal(data, &aliases); err != nil {
		panic(err)
	}

	// Procura o alias pelo nome
	for _, a := range aliases {
		if a.Name == name {
			return structs.Alias{Name: a.Name, Command: a.Command}
		}
	}

	return structs.Alias{Name: "", Command: ""}
}
