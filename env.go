package dotenv

import (
	"bufio"
	"fmt"
	"maps"
	"os"
	"strings"
	"sync"
	"unicode"
)

var (
	buffer    = make(map[string]string)
	errGlobal error
	mu        sync.RWMutex
)

// init carrega automaticamente o .env
func init() {
	errGlobal = loadFile(".env")
}

// Load permite carregar manualmente outro arquivo
func Load(path string) error {
	mu.Lock()
	defer mu.Unlock()

	buffer = make(map[string]string)
	errGlobal = loadFile(path)
	return errGlobal
}

// Reload recarrega o .env padrão
func Reload() error {
	return Load(".env")
}

// função interna de leitura
func loadFile(path string) error {
	file, err := os.Open(path)
	if err != nil {
		return err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	lineNum := 0

	for scanner.Scan() {
		lineNum++
		line := strings.TrimSpace(scanner.Text())

		// ignora vazios e comentários
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}

		key, value, err := parseLine(line, lineNum)
		if err != nil {
			return err
		}

		buffer[key] = value
	}

	return scanner.Err()
}

// valida chave
func isValidKey(s string) bool {
	if len(s) == 0 {
		return false
	}

	for i, r := range s {
		// primeiro caractere não pode ser número
		if i == 0 && unicode.IsDigit(r) {
			return false
		}

		if !unicode.IsLetter(r) && !unicode.IsDigit(r) && r != '_' {
			return false
		}
	}

	return true
}

// parse de linha
func parseLine(line string, lineNum int) (string, string, error) {
	parts := strings.SplitN(line, "=", 2)
	if len(parts) != 2 {
		return "", "", fmt.Errorf("invalid syntax at line %d: %s", lineNum, line)
	}

	key := strings.TrimSpace(parts[0])
	value := strings.TrimSpace(parts[1])

	if !isValidKey(key) {
		return "", "", fmt.Errorf("invalid key at line %d: %s", lineNum, key)
	}

	return key, value, nil
}

// retorna cópia segura
func GetEnvs() (map[string]string, error) {
	mu.RLock()
	defer mu.RUnlock()

	if errGlobal != nil {
		return nil, errGlobal
	}

	return maps.Clone(buffer), nil
}

// retorna valor específico
func GetEnv(key string) (string, error) {
	mu.RLock()
	defer mu.RUnlock()

	if errGlobal != nil {
		return "", errGlobal
	}

	key = strings.TrimSpace(key)

	if !isValidKey(key) {
		return "", fmt.Errorf("invalid key: %s", key)
	}

	value, ok := buffer[key]
	if !ok {
		return "", fmt.Errorf("env not found: %s", key)
	}

	return value, nil
}
