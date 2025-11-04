package main

import (
	"bytes"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"text/template"

	fileType "github.com/karabomatabane/bash-flutter-cli/constants"
)

const version = "1.0.0"

// Template functions
var funcMap = template.FuncMap{
	"toSnake":  toSnakeCase,
	"toCamel":  toCamelCase,
	"toPascal": toPascalCase,
	"toLower":  strings.ToLower,
	"toUpper":  strings.ToUpper,
}

type TemplateData struct {
	Name string
}

func main() {
	if len(os.Args) < 2 {
		printUsage()
		os.Exit(1)
	}

	command := os.Args[1]

	switch command {
	case "g", "generate":
		if len(os.Args) < 4 {
			fmt.Println("Usage: bf g <type> <path>")
			fmt.Println("Example: bf g p pages/home")
			os.Exit(1)
		} else if len(os.Args) == 4 {
			generateBoilerplate(os.Args[2], os.Args[3])
			return
		}
		generateBoilerplate(os.Args[2], os.Args[3], os.Args[4])
	case "list", "ls":
		listTemplates()
	case "version", "-v", "--version":
		fmt.Printf("bf v%s\n", version)
	case "help", "-h", "--help":
		printUsage()
	default:
		fmt.Printf("Unknown command: %s\n", command)
		printUsage()
		os.Exit(1)
	}
}

func generateBoilerplate(args ...string) {
	if len(args) == 0 {
		fmt.Println("Error: Expected at least two arguments, but received none.")
		fmt.Println("\nRun bf -h for help")
		os.Exit(1)
	}
	templateType := args[0]
	path := args[1]
	var skipBloc bool
	if len(args) > 2 {
		skipBloc = args[2] == "--skip-bloc"
	}

	// Map short codes to template names
	typeMap := map[string]string{
		"p": "page",
		"b": "bloc",
		"e": "event",
		"s": "state",
	}

	// Get full template name
	fullType := templateType
	if mapped, ok := typeMap[templateType]; ok {
		fullType = mapped
	}

	// If input is bf g p <path> with --skip-bloc
	// Generate page with bloc
	if !skipBloc && fullType == fileType.Page {
		generateFile(fileType.BloCPage, path)
		generateFile(fileType.BLoC, path)
		generateFile(fileType.Event, path)
		generateFile(fileType.State, path)
	} else {
		generateFile(fullType, path)
	}
}

func generateFile(fullType, path string) {
	// Look for template file
	templatePath := filepath.Join("templates", fullType+".tmpl")
	if _, err := os.Stat(templatePath); os.IsNotExist(err) {
		fmt.Printf("Error: Template not found: %s\n", templatePath)
		fmt.Println("\nEnsure installation completed successfully.")
		os.Exit(1)
	}

	// Read template
	templateContent, err := os.ReadFile(templatePath)
	if err != nil {
		fmt.Printf("Error reading template: %v\n", err)
		os.Exit(1)
	}

	// Parse path to get directory and name
	dir := filepath.Dir(path)
	baseName := filepath.Base(path)

	// Create directory if it doesn't exist
	if dir != "." && dir != "" {
		if err := os.MkdirAll(dir, 0755); err != nil {
			fmt.Printf("Error creating directory: %v\n", err)
			os.Exit(1)
		}
	}

	// Generate filename based on template type
	var filename string
	switch fullType {
	case fileType.Page:
		filename = fmt.Sprintf("%s_page.dart", toSnakeCase(baseName))
	case fileType.BloCPage:
		filename = fmt.Sprintf("%s_page.dart", toSnakeCase(baseName))
	case fileType.BLoC:
		filename = fmt.Sprintf("bloc/%s_bloc.dart", toSnakeCase(baseName))
	case fileType.Event:
		filename = fmt.Sprintf("bloc/%s_event.dart", toSnakeCase(baseName))
	case fileType.State:
		filename = fmt.Sprintf("bloc/%s_state.dart", toSnakeCase(baseName))
	default:
		filename = fmt.Sprintf("%s.dart", toSnakeCase(baseName))
	}

	// Check if bloc dir is required and create it if so
	if fullType != fileType.Page {
		blocPath := filepath.Join(dir, "bloc")

		if info, err := os.Stat(blocPath); os.IsNotExist(err) || !info.IsDir() {
			fmt.Println("Creating bloc directory")
			if err := os.MkdirAll(blocPath, 0755); err != nil {
				fmt.Printf("Error creating bloc directory: %v\n", err)
				os.Exit(1)
			}
		}
	}

	outputPath := filepath.Join(dir, filename)

	// Check if file exists
	if _, err := os.Stat(outputPath); err == nil {
		fmt.Printf("Error: File already exists: %s\n", outputPath)
		os.Exit(1)
	}

	// Parse and execute template
	tmpl, err := template.New("dart").Funcs(funcMap).Parse(string(templateContent))
	if err != nil {
		fmt.Printf("Error parsing template: %v\n", err)
		os.Exit(1)
	}

	var buf bytes.Buffer
	data := TemplateData{Name: baseName}
	if err := tmpl.Execute(&buf, data); err != nil {
		fmt.Printf("Error executing template: %v\n", err)
		os.Exit(1)
	}

	// Write file
	if err := os.WriteFile(outputPath, buf.Bytes(), 0644); err != nil {
		fmt.Printf("Error writing file: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("✓ Generated %s\n", outputPath)
}

func listTemplates() {
	templateDir := "templates"
	entries, err := os.ReadDir(templateDir)
	if err != nil {
		fmt.Println("No templates found. Installation may be broken. Please reinstall and try again.")
		return
	}

	fmt.Println("Available templates:")
	for _, entry := range entries {
		if !entry.IsDir() && strings.HasSuffix(entry.Name(), ".tmpl") {
			name := strings.TrimSuffix(entry.Name(), ".tmpl")
			fmt.Printf("  %s\n", name)
		}
	}
}

func printUsage() {
	fmt.Println("Bash Flutter CLI")
	fmt.Println("\nUsage:")
	fmt.Println("  bf g <type> <path>    Generate a file from template")
	fmt.Println("  bf list               List available templates")
	fmt.Println("  bf version            Show version")
	fmt.Println("\nExamples:")
	fmt.Println("  bf g p pages/home           # Generate pages/home_page.dart with relevate bloc")
	fmt.Println("  bf g page pages/home        # Same as above")
	fmt.Println("  bf g b blocs/counter        # Generate blocs/counter_bloc.dart")
	fmt.Println("  bf g w widgets/custom_btn   # Generate widgets/custom_btn_widget.dart")
	fmt.Println("\nOptional flags:")
	fmt.Println("  --skip-bloc                 # Generates page without bloc")
	fmt.Println("\nShort codes:")
	fmt.Println("  p = page")
	fmt.Println("  b = bloc")
	// fmt.Println("  w = widget")
	// fmt.Println("  m = model")
	fmt.Println("  e = event")
	fmt.Println("  s = state")
}

// toSnakeCase converts string to snake_case
func toSnakeCase(s string) string {
	// Handle acronyms and split on capitals
	re := regexp.MustCompile("([a-z0-9])([A-Z])")
	snake := re.ReplaceAllString(s, "${1}_${2}")

	// Replace any non-alphanumeric with underscore
	re = regexp.MustCompile("[^a-zA-Z0-9]+")
	snake = re.ReplaceAllString(snake, "_")

	return strings.ToLower(strings.Trim(snake, "_"))
}

// toCamelCase converts string to camelCase
func toCamelCase(s string) string {
	s = toPascalCase(s)
	if len(s) == 0 {
		return s
	}
	return strings.ToLower(s[:1]) + s[1:]
}

// toPascalCase converts string to PascalCase
func toPascalCase(s string) string {
	// Split by common delimiters
	words := regexp.MustCompile(`[_\-\s]+`).Split(s, -1)
	for i, word := range words {
		if len(word) > 0 {
			words[i] = strings.ToUpper(word[:1]) + strings.ToLower(word[1:])
		}
	}
	return strings.Join(words, "")
}
