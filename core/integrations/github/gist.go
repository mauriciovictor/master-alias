package github

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"

	"master-alias.com/core"
	"master-alias.com/core/structs"
)

func ExportGist() {
	defaultPath := GetDefaultPath()

	aliasJson := filepath.Join(defaultPath, "alias.json")

	content, err := os.ReadFile(aliasJson)
	if err != nil {
		panic(err)
	}

	gist := map[string]interface{}{
		"description": "Alias do Master CLI",
		"public":      false,
		"files": map[string]map[string]string{
			"alias.json": {
				"content": string(content),
			},
		},
	}

	data, _ := json.Marshal(gist)
	var config, _ = core.LoadConfig()
	url := "https://api.github.com/gists"
	method := "POST"
	var gistJsonExist = FileExists()

	if gistJsonExist {
		gistConfig, _ := ReadGithubGistJson()
		url = "https://api.github.com/gists/" + gistConfig.Id
		method = "PATCH"
	}

	req, _ := http.NewRequest(method, url, bytes.NewBuffer(data))
	req.Header.Set("Authorization", "token "+config.GithubToken)
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		panic(err)
	}

	// Estrutura simples para extrair o ID do Gist
	var result struct {
		ID  string `json:"id"`
		URL string `json:"html_url"`
	}

	if err := json.Unmarshal(body, &result); err != nil {
		panic(err)
	}
	if !gistJsonExist {
		CreateGithubGistJson(structs.GithubGist{
			Id:  result.ID,
			Url: result.URL,
		})

		fmt.Println("✅ Gist criado com sucesso!")
		fmt.Println("🆔 ID:", result.ID)
		fmt.Println("🌐 URL:", result.URL)

		return
	}

	fmt.Println("✅ Gist atualizado com sucesso!")
}

func ImportGist(gist_id string) {
	url := fmt.Sprintf("https://api.github.com/gists/%s", gist_id)
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		fmt.Println("❌ Erro ao criar requisição:", err)
		return
	}
	var config, _ = core.LoadConfig()
	// Cabeçalho de autenticação com o token
	req.Header.Set("Authorization", "token "+config.GithubToken)
	req.Header.Set("Accept", "application/vnd.github+json")
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("❌ Erro ao enviar requisição:", err)
		return
	}

	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		fmt.Println("❌ Erro: Gist não encontrado ou acesso negado")
		return
	}

	var gistData struct {
		Files map[string]struct {
			Content string `json:"content"`
		} `json:"files"`
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("❌ Erro ao ler resposta:", err)
		return
	}

	if err := json.Unmarshal(body, &gistData); err != nil {
		fmt.Println("❌ Erro ao decodificar JSON:", err)
		return
	}

	// Pegamos o primeiro arquivo do Gist
	for filename, file := range gistData.Files {
		home, _ := os.UserHomeDir()
		dir := filepath.Join(home, ".master-alias")
		os.MkdirAll(dir, 0755)

		dest := filepath.Join(dir, filename)
		err := os.WriteFile(dest, []byte(file.Content), 0644)
		if err != nil {
			fmt.Println("❌ Erro ao salvar arquivo:", err)
			return
		}

		CreateGithubGistJson(structs.GithubGist{
			Id:  gist_id,
			Url: "https://gist.github.com/" + gist_id,
		})

		fmt.Printf("✅ Arquivo %s criado com sucesso em %s\n", filename, dest)
		return
	}

	fmt.Println("⚠️ Nenhum arquivo encontrado no Gist.")
}

func GetDefaultPath() string {
	home, _ := os.UserHomeDir()
	defaultPath := filepath.Join(home, ".master-alias")

	return defaultPath
}

func CreateGithubGistJson(gist structs.GithubGist) (structs.GithubGist, error) {
	defaultPath := GetDefaultPath()
	gistFilePath := filepath.Join(defaultPath, "github_gist.json")

	f, err := os.Create(gistFilePath)

	if err != nil {
		return structs.GithubGist{}, err
	}

	defer f.Close()

	// Codifica o objeto em JSON formatado
	encoder := json.NewEncoder(f)
	encoder.SetIndent("", "  ")

	if err := encoder.Encode(gist); err != nil {
		return structs.GithubGist{}, err
	}

	return gist, nil
}

func ReadGithubGistJson() (structs.GithubGist, error) {
	defaultPath := GetDefaultPath()
	gistFilePath := filepath.Join(defaultPath, "github_gist.json")

	file, err := os.Open(gistFilePath)

	if err != nil {
		return structs.GithubGist{}, err
	}

	defer file.Close()

	var gist structs.GithubGist
	decoder := json.NewDecoder(file)
	err = decoder.Decode(&gist)

	return gist, err
}

func FileExists() bool {
	defaultPath := GetDefaultPath()
	gistFilePath := filepath.Join(defaultPath, "github_gist.json")

	if _, err := os.Stat(gistFilePath); err == nil {
		return true
	} else if os.IsNotExist(err) {
		return false
	}

	return false
}
