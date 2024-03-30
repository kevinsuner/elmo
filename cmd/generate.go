package cmd

import (
	"bytes"
	"errors"
	"fmt"
	"html/template"
	"os"
	"path/filepath"
	"strings"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/yuin/goldmark"
	"github.com/yuin/goldmark-meta"
	"github.com/yuin/goldmark/extension"
	"github.com/yuin/goldmark/parser"
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

func createDir(path string) error {
    _, err := os.Stat(path)
    switch {
    case os.IsExist(err):
        return os.RemoveAll(path)
    case os.IsNotExist(err):
        return os.MkdirAll(path, os.ModePerm)
    }

    return err
}

type fileInfo struct {
    kind string
    name string
}

func getFiles(dir, ext string) ([]fileInfo, error) {
    items, err := os.ReadDir(dir)
    if err != nil {
        return nil, err
    }

    var kind string
    var files []fileInfo
    for _, item := range items {
        if item.IsDir() {
            subitems, err := os.ReadDir(fmt.Sprintf("%s/%s", dir, item.Name()))
            if err != nil {
                return nil, err
            }

            switch {
            case item.Name() == "partials":
                kind = item.Name()
            case item.Name() == "posts":
                kind = item.Name()
            }

            for _, subitem := range subitems {
                if filepath.Ext(subitem.Name()) == ext {
                    files = append(files, fileInfo{
                        kind: kind,
                        name: subitem.Name(),
                    })
                }
            }
        } else {
            if filepath.Ext(item.Name()) == ext {
                files = append(files, fileInfo{
                    name: item.Name(),
                })
            }
        }
    }

    return files, nil
}

type content struct {
    meta map[string]interface{}
    html template.HTML
}

func parseContent() (map[string]content, error) {
    files, err := getFiles(viper.GetString("ContentDir"), ".md")
    if err != nil {
        return nil, err
    }

    contents := make(map[string]content, len(files))
    for _, file := range files {
        filename := strings.Split(file.name, ".")[0]
        if _, ok := contents[filename]; !ok {
            var out []byte
            switch file.kind {
            case "partials":
                continue
            case "posts":
                out, err = os.ReadFile(fmt.Sprintf("%s/%s/%s",
                    viper.GetString("ContentDir"),
                    viper.GetString("PostDir"),
                    file.name,
                ))
            default:
                out, err = os.ReadFile(fmt.Sprintf("%s/%s",
                    viper.GetString("ContentDir"), file.name,
                ))
            }

            if err != nil {
                return nil, err
            }

            markdown := goldmark.New(
                goldmark.WithExtensions(
                    extension.GFM,
                    meta.Meta,
                ),
                goldmark.WithParserOptions(parser.WithAutoHeadingID()),
            )

            var buf bytes.Buffer
            ctx := parser.NewContext()
            err = markdown.Convert(out, &buf, parser.WithContext(ctx))
            if err != nil {
                return nil, err
            }

            contents[filename] = content{
                meta: meta.Get(ctx),
                html: template.HTML(buf.String()),
            }
        }
    }

    return contents, nil
}

func getPartials(files []fileInfo) []string {
    partials := make([]string, 0, len(files))
    for _, file := range files {
        if file.kind == "partials" {
            partials = append(partials, fmt.Sprintf("%s/%s/%s",
                viper.GetString("LayoutDir"),
                viper.GetString("PartialDir"),
                file.name,
            ))
        }
    }

    return partials
}

func parseLayout(contents map[string]content) (map[string]string, error) {
    files, err := getFiles(viper.GetString("LayoutDir"), ".tmpl")
    if err != nil {
        return nil, err
    }

    partials := getPartials(files)
    layouts := make(map[string]string, len(files))
    for _, file := range files {
        if file.kind == "partials" || file.kind == "posts" { continue }

        filename := strings.Split(file.name, ".")[0]
        if len(filename) == 0 {
            return nil, errors.New("invalid file")
        }

        if _, ok := contents[filename]; ok {
            partials = append(partials, fmt.Sprintf("%s/%s",
                viper.GetString("LayoutDir"),
                file.name,
            ))

            tmpl, err := template.New(filename).ParseFiles(partials...)
            if err != nil {
                return nil, err
            }

            var buf bytes.Buffer
            if err := tmpl.Execute(
                &buf, 
                map[string]interface{}{
                    "meta": contents[filename].meta,
                    "html": contents[filename].html,
                },
            ); err != nil {
                return nil, err
            }

            layouts[filename] = buf.String()
        }
    }

    return layouts, nil
}

func parsePosts(contents map[string]content) (map[string]string, error) {
    files, err := getFiles(viper.GetString("LayoutDir"), ".tmpl")
    if err != nil {
        return nil, err
    }
    partials := getPartials(files)

    files, err = getFiles(viper.GetString("ContentDir"), ".md")
    if err != nil {
        return nil, err
    }

    posts := make(map[string]string, len(files))
    for _, file := range files {
        if len(file.kind) == 0 || file.kind == "partials" { continue }

        filename := strings.Split(file.name, ".")[0]
        if len(filename) == 0 {
            return nil, errors.New("invalid file")
        }

        if _, ok := contents[filename]; ok {
            partials = append(partials, fmt.Sprintf("%s/%s",
                viper.GetString("LayoutDir"),
                "post.html.tmpl",
            ))

            tmpl, err := template.New("post").ParseFiles(partials...)
            if err != nil {
                return nil, err
            }

            var buf bytes.Buffer
            if err := tmpl.Execute(
                &buf,
                map[string]interface{}{
                    "meta": contents[filename].meta,
                    "html": contents[filename].html,
                },
            ); err != nil {
                return nil, err
            }

            posts[filename] = buf.String()
        }
    }

    return posts, nil
}

func generatePublic(layouts map[string]string) error {
    for name, html := range layouts {
        file, err := os.Create(
            fmt.Sprintf("%s/%s.html", viper.GetString("PublicDir"), name),
        )
        if err != nil {
            return err
        }

        if _, err := file.WriteString(html); err != nil {
            return err
        }
    }

    return nil
}

func generatePosts(posts map[string]string) error {
    for name, html := range posts {
        file, err := os.Create(
            fmt.Sprintf(
                "%s/%s/%s.html",
                viper.GetString("PublicDir"),
                viper.GetString("PostDir"),
                name,
            ),
        )
        if err != nil {
            return err
        }

        if _, err := file.WriteString(html); err != nil {
            return err
        }
    }

    return nil
}

func generate(cmd *cobra.Command, args []string) {
    if err := createDir(fmt.Sprintf("%s/%s",
        viper.GetString("PublicDir"),
        viper.GetString("PostDir"),
    )); err != nil {
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

    posts, err := parsePosts(contents)
    if err != nil {
        fmt.Println(err.Error())
        os.Exit(1)
    }

    if len(layouts) > 0 {
        if err := generatePublic(layouts); err != nil {
            fmt.Println(err.Error())
            os.Exit(1)
        }
    }

    if len(posts) > 0 {
        if err := generatePosts(posts); err != nil {
            fmt.Println(err.Error())
            os.Exit(1)
        }
    }
}
