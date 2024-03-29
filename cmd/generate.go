package cmd

import (
	"bytes"
	"errors"
	"fmt"
	"html/template"
	"io/fs"
	"os"
	"path/filepath"
	"strings"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/yuin/goldmark"
)

func init() {
    rootCmd.AddCommand(generateCmd)
}

var generateCmd = &cobra.Command{
    Use: "generate",
    Short: "placeholder",
    Long: "placeholder",
    Run: generate,
}

func createPublicDir() error {
    _, err := os.Stat(viper.GetString("PublicDir"))
    switch {
    case os.IsNotExist(err):
        return os.Mkdir(viper.GetString("PublicDir"), os.ModePerm)
    case os.IsExist(err):
        return os.Remove(viper.GetString("PublicDir"))
    default:
        return err
    }
}

func getFiles(dir, ext string) ([]string, error) {
    var files []string
    fsys := os.DirFS(dir)
    if err := fs.WalkDir(
        fsys,
        ".",
        func(path string, d fs.DirEntry, err error) error {
            if !strings.Contains(path, "/") {
                if filepath.Ext(path) == ext {
                    files = append(files, d.Name())
                }
            }

            return nil
        },
    ); err != nil {
        return nil, err
    }

    return files, nil
}

func parseContent() (map[string]string, error) {
    files, err := getFiles(viper.GetString("ContentDir"), ".md")
    if err != nil {
        return nil, err
    }

    contents := make(map[string]string, len(files))
    for _, file := range files {
        filename := strings.Split(file, ".")[0]
        if _, ok := contents[filename]; !ok {
            content, err := os.ReadFile(
                fmt.Sprintf("%s/%s", viper.GetString("ContentDir"), file),
            )
            if err != nil {
                return nil, err
            }

            var buf bytes.Buffer
            if err := goldmark.Convert(content, &buf); err != nil {
                return nil, err
            }

            contents[filename] = buf.String()
        }
    }

    return contents, nil
}

func parseLayout(contents map[string]string) (map[string]string, error) {
    files, err := getFiles(viper.GetString("LayoutDir"), ".tmpl")
    if err != nil {
        return nil, err
    }

    if len(files) != len(contents) {
        return nil, errors.New("mismatched files and contents length")
    }

    layouts := make(map[string]string, len(files))
    for _, file := range files {
        filename := strings.Split(file, ".")[0]
        if len(filename) == 0 {
            return nil, errors.New("invalid file")
        }

        if _, ok := contents[filename]; ok {
            tmpl, err := template.New(filename).ParseFiles(
                fmt.Sprintf("%s/%s", viper.GetString("LayoutDir"), file),
            )
            if err != nil {
                return nil, err
            }

            var buf bytes.Buffer
            if err := tmpl.Execute(
                &buf, 
                map[string]template.HTML{
                    "content": template.HTML(contents[filename]),
                },
            ); err != nil {
                return nil, err
            }

            layouts[filename] = buf.String()
        }
    }

    return layouts, nil
} 

func generatePublic(layouts map[string]string) error {
    for name, layout := range layouts {
        file, err := os.Create(
            fmt.Sprintf("%s/%s.html", viper.GetString("PublicDir"), name),
        )
        if err != nil {
            return err
        }

        if _, err := file.WriteString(layout); err != nil {
            return err
        }
    }

    return nil
}

func generate(cmd *cobra.Command, args []string) {
    if err := createPublicDir(); err != nil {
        fmt.Println(err.Error())
        os.Exit(1)
    }

    contents, err := parseContent()
    if err != nil {
        fmt.Println(err.Error())
        os.Exit(1)
    }

    layouts, err := parseLayout(contents)
    if err != nil {
        fmt.Println(err.Error())
        os.Exit(1)
    }

    if err := generatePublic(layouts); err != nil {
        fmt.Println(err.Error())
        os.Exit(1)
    }
}
