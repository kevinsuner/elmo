package cmd

import (
	"bytes"
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

func getFilepaths(dir, ext string) ([]string, error) {
    files, err := os.ReadDir(dir)
    if err != nil {
        return nil, err
    }

    filepaths := make([]string, 0, len(files))
    for _, file := range files {
        if !file.IsDir() {
            if filepath.Ext(file.Name()) == ext {
                filepaths = append(filepaths, fmt.Sprintf("%s/%s", dir, file.Name()))
            }
        }
    }

    return filepaths, nil
}

type entry struct {
    name    string
}

func getFilesByExt(dir, ext string) ([]entry, error) {
    files, err := os.ReadDir(dir)
    if err != nil {
        return nil, err
    }

    entries := make([]entry, 0, len(files))
    for _, file := range files {
        if !file.IsDir() {
            if filepath.Ext(file.Name()) == ext {
                entries = append(entries, entry{
                    name: file.Name(),
                })
            }
        }
    }

    return entries, nil
}

type content struct {
    meta map[string]interface{}
    html template.HTML
}

func parseMarkdown(dir, ext string, contents map[string]content) error {
    entries, err := getFilesByExt(dir, ext)
    if err != nil {
        return err
    }

    for _, entry := range entries {
        out, err := os.ReadFile(fmt.Sprintf("%s/%s", dir, entry.name))
        if err != nil {
            return err
        }
        
        markdown := goldmark.New(
            goldmark.WithExtensions(extension.GFM, meta.Meta),
            goldmark.WithParserOptions(parser.WithAutoHeadingID()),
        )

        var buf bytes.Buffer
        ctx := parser.NewContext()
        err = markdown.Convert(out, &buf, parser.WithContext(ctx))
        if err != nil {
            return err
        }

        filename := strings.Split(entry.name, ".")[0]
        contents[filename] = content{
            meta: meta.Get(ctx),
            html: template.HTML(buf.String()),
        }
    }

    return nil
}

type tpl struct {
   tag string
   content string
}

func parseTemplates(
    dir, ext string,
    post bool,
    templates map[string]tpl,
    contents map[string]content,
) error {
    entries, err := getFilesByExt(dir, ext)
    if err != nil {
        return err
    }

    files, err := getFilepaths(fmt.Sprintf("%s/%s",
        viper.GetString("LayoutDir"),
        viper.GetString("PartialDir"),
    ), ".tmpl")

    for _, entry := range entries {
        filename := strings.Split(entry.name, ".")[0]
        var tmpl *template.Template
        var tag string
        if post {
            files = append(files, 
                fmt.Sprintf(
                    "%s/%s", dir, "post.html.tmpl",
                ),
            )

            tag = "post"
            tmpl, err = template.New("post").ParseFiles(files...)
            if err != nil {
                return err
            }
        } else {
            tag = "page"
            files = append(files, fmt.Sprintf("%s/%s", dir, entry.name))
            tmpl, err = template.New(filename).ParseFiles(files...)
            if err != nil {
                return err
            }
        }
        
        var buf bytes.Buffer
        if err := tmpl.Execute(
            &buf, 
            map[string]interface{}{
                "meta": contents[filename].meta,
                "html": contents[filename].html,
            },
        ); err != nil {
            return err
        }

        templates[filename] = tpl{
            tag: tag,
            content: buf.String(),
        }
    }

    return nil
}

func createFiles(templates map[string]tpl) error {
    for filename, tpl := range templates {
        dir := viper.GetString("PublicDir") 
        if tpl.tag == "post" {
             dir = fmt.Sprintf(
                 "%s/%s", 
                 viper.GetString("PublicDir"), 
                 viper.GetString("PostDir"),
             ) 
        }

        file, err := os.Create(fmt.Sprintf("%s/%s.html", dir, filename))
        if err != nil {
            return err
        }

        if _, err := file.WriteString(tpl.content); err != nil {
            return err
        }
    }

    return nil
}

func generate(cmd *cobra.Command, args []string) {
    err := createDir(fmt.Sprintf("%s/%s",
        viper.GetString("PublicDir"),
        viper.GetString("PostDir"),
    )) 
    if err != nil {
        fmt.Println(err.Error())
        os.Exit(1)
    }

    contents := make(map[string]content)
    err = parseMarkdown(viper.GetString("ContentDir"), ".md", contents)
    if err != nil {
        fmt.Println(err.Error())
        os.Exit(1)
    }

    if err = parseMarkdown(
        fmt.Sprintf(
            "%s/%s",
            viper.GetString("ContentDir"),
            viper.GetString("PostDir"),
        ),
        ".md", contents,
    ); err != nil {
        fmt.Println(err.Error())
        os.Exit(1)
    }

    // TODO: Fix this, only parsing index and post

    templates := make(map[string]tpl)
    if err = parseTemplates(
        viper.GetString("LayoutDir"), 
        ".tmpl", 
        false, 
        templates,
        contents,
    ); err != nil {
        fmt.Println(err.Error())
        os.Exit(1)
    }

    if err = parseTemplates(
        viper.GetString("LayoutDir"), 
        ".tmpl", 
        true, 
        templates,
        contents,
    ); err != nil {
        fmt.Println(err.Error())
        os.Exit(1)
    }

    if err = createFiles(templates); err != nil {
        fmt.Println(err.Error())
        os.Exit(1)
    }
}
