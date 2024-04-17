package main

import (
	"bytes"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"strings"
	"text/template"

	"github.com/sirupsen/logrus"
	"github.com/urfave/cli/v2"
	"gopkg.in/yaml.v2"
)

const (
	CollectionsDir = "resources/collections"
	TemplatesDir   = "resources/templates"
	Template       = "chemfig.tex"
)

type Collection struct {
	Name       string       `yaml:"name"`
	Structures []*Structure `yaml:"structures"`
}

type Structure struct {
	Name         string `yaml:"name"`
	CompoundName string `yaml:"compound_name"`
	Chemfig      string `yaml:"chemfig"`
}

func main() {
	app := &cli.App{
		Name: "chemfig",
		Commands: []*cli.Command{
			{
				Name:   "generate",
				Usage:  "generate LaTeX files",
				Action: Generate,
			},
			{
				Name:   "targets",
				Usage:  "return a list of files that will be generated",
				Action: Targets,
			},
		},
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:     "output-dir",
				Aliases:  []string{"o"},
				Usage:    "output directory",
				Required: true,
			},
		},
	}

	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}

func Generate(c *cli.Context) error {
	outputDir := c.String("output-dir")

	template, err := readTemplate(Template)
	if err != nil {
		return err
	}

	collections, err := readCollections()
	if err != nil {
		return err
	}

	for _, c := range collections {
		for _, s := range c.Structures {
			buf := bytes.NewBuffer(nil)
			if err := template.Execute(buf, s); err != nil {
				return err
			}

			err := replaceFileIfChanged(getOutputFileName(outputDir, c.Name, s.Name), buf.Bytes())
			if err != nil {
				return err
			}
		}
	}

	return nil
}

func Targets(c *cli.Context) error {
	outputDir := c.String("output-dir")

	collections, err := readCollections()
	if err != nil {
		return err
	}

	for _, c := range collections {
		for _, s := range c.Structures {
			fmt.Println(getOutputFileName(outputDir, c.Name, s.Name))
		}
	}

	return nil
}

func readCollections() ([]*Collection, error) {
	files := findCollectionFiles()

	collections := make([]*Collection, 0)

	for _, file := range files {
		c, err := readCollection(file)
		if err != nil {
			return nil, err
		}

		collections = append(collections, c)
	}

	return collections, nil
}

func findCollectionFiles() []string {
	collections := make([]string, 0)

	for _, ext := range []string{".yaml", ".yml"} {
		files, err := filepath.Glob(filepath.Join(CollectionsDir, "*"+ext))
		if err != nil {
			logrus.WithError(err).Error("failed to find collection files")

			continue
		}

		collections = append(collections, files...)
	}

	return collections
}

func readCollection(file string) (*Collection, error) {
	f, err := os.Open(file)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	data, err := io.ReadAll(f)
	if err != nil {
		return nil, err
	}

	collection := new(Collection)

	err = yaml.Unmarshal(data, collection)
	if err != nil {
		return nil, err
	}

	collection.Name = strings.TrimSuffix(filepath.Base(file), filepath.Ext(file))

	return collection, nil
}

func readTemplate(name string) (*template.Template, error) {
	path := filepath.Join("resources/templates", name)

	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	data, err := io.ReadAll(f)
	if err != nil {
		return nil, err
	}

	return template.New(name).Funcs(template.FuncMap{
		"trim": strings.TrimSpace,
		"indentBlock": func(s string, spaces int) string {
			s = strings.TrimSpace(s)

			parts := strings.Split(s, "\n")
			for i, p := range parts {
				parts[i] = strings.TrimSpace(p)
			}

			return strings.Join(strings.Split(s, "\n"), "\n"+strings.Repeat(" ", spaces))
		},
	}).Parse(string(data))
}

func replaceFileIfChanged(path string, data []byte) error {
	if err := os.MkdirAll(filepath.Dir(path), os.ModePerm); err != nil {
		return err
	}

	_, err := os.Stat(path)
	if os.IsNotExist(err) {
		logrus.Infof("Creating file %s", path)

		return os.WriteFile(path, data, 0644)
	}

	d, err := os.ReadFile(path)
	if err != nil {
		return err
	}

	if string(d) == string(data) {
		return nil
	}

	logrus.Infof("Updating file %s", path)
	return os.WriteFile(path, data, 0644)
}

func getOutputFileName(outputDir string, collectionName string, structureName string) string {
	return filepath.Join(outputDir, sanitizeName(collectionName), sanitizeName(structureName)+".tex")
}

func sanitizeName(name string) string {
	s := strings.ReplaceAll(name, " ", "_")
	s = strings.TrimSpace(s)

	return s
}
