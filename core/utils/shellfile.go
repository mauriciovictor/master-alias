package utils

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"master-alias.com/core"
)

func EnsureShellSource() error {
	home, _ := os.UserHomeDir()

	var config, _ = core.LoadConfig()

	var shellConfigFile = ".bashrc"

	switch config.Shell {
	case "zsh":
		shellConfigFile = ".zshrc"
	}

	rcPath := filepath.Join(home, shellConfigFile) // ou ".zshrc"

	content, _ := os.ReadFile(rcPath)
	filePath := filepath.Join(home, ".master-alias", "master_aliases.sh")

	line := fmt.Sprintf("source %s", filePath)

	exec.Command("bash", "-c", line)

	if !strings.Contains(string(content), line) {
		f, err := os.OpenFile(rcPath, os.O_APPEND|os.O_WRONLY, 0644)
		if err != nil {
			return err
		}
		defer f.Close()

		if _, err := f.WriteString("\n" + line); err != nil {
			return err
		}
	}

	return nil
}

func CreateShellAliasFile() error {
	home, _ := os.UserHomeDir()
	filePath := filepath.Join(home, ".master-alias", "master_aliases.sh")

	f, err := os.Create(filePath)

	if err != nil {
		return err
	}
	defer f.Close()

	return nil
}

func WriteAliasToFile(name string) error {
	var command = "master-alias run " + name

	home, err := os.UserHomeDir()
	if err != nil {
		return err
	}

	aliasFile := filepath.Join(home, ".master-alias", "master_aliases.sh")

	f, err := os.OpenFile(aliasFile, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)

	if err != nil {
		return err
	}
	defer f.Close()

	if strings.Contains(command, "$") || strings.Contains(command, ">") || strings.Contains(command, "|") {
		_, err = fmt.Fprintf(f, "%s() { %s; }\n", name, command)
	} else {
		escaped := strings.ReplaceAll(command, "'", `'\''`)
		_, err = fmt.Fprintf(f, "alias %s='%s'\n", name, escaped)
	}
	return err
}

func RemoveAliasFromFile(name string) error {
	home, err := os.UserHomeDir()
	if err != nil {
		return err
	}

	aliasFile := filepath.Join(home, ".master_aliases.sh")

	input, err := os.Open(aliasFile)
	if err != nil {
		// Se o arquivo não existe, apenas retorna
		if os.IsNotExist(err) {
			return nil
		}
		return err
	}
	defer input.Close()

	var lines []string
	scanner := bufio.NewScanner(input)

	for scanner.Scan() {
		line := scanner.Text()

		// Remove a linha se contiver o alias (pode ser `alias name=` ou `name() {`)
		if strings.HasPrefix(line, "alias "+name+"=") ||
			strings.HasPrefix(line, name+"()") {
			continue // pula essa linha
		}

		lines = append(lines, line)
	}

	if err := scanner.Err(); err != nil {
		return err
	}

	// Reescreve o arquivo sem o alias removido
	return os.WriteFile(aliasFile, []byte(strings.Join(lines, "\n")+"\n"), 0644)
}
