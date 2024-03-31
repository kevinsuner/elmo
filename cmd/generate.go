package cmd

import (
	"bytes"
	"fmt"
	"html/template"
	"log/slog"
	"os"
	"path/filepath"
	"strings"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/yuin/goldmark"
	"github.com/yuin/goldmark-meta"
	"github.com/yuin/goldmark/extension"
	"github.com/yuin/goldmark/parser"
	"golang.org/x/text/cases"
	"golang.org/x/text/language"
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

    filenames := make([]string, 0, len(files))
    for _, file := range files {
        if !file.IsDir() && filepath.Ext(file.Name()) == ext {
            filenames = append(filenames, fmt.Sprintf("%s/%s", dir, file.Name()))
        }
    }

    return filenames, nil
}

type entry struct {
    kind    string
    name    string
}

func getFilesByExt(dir, ext string) ([]entry, error) {
    files, err := os.ReadDir(dir)
    if err != nil {
        return nil, err
    }

    var kind string = "template"
    if dir == viper.GetString("ContentPostsDir") {
        kind = "content"
    }

    entries := make([]entry, 0, len(files))
    for _, file := range files {
        if !file.IsDir() {
            if filepath.Ext(file.Name()) == ext {
                entries = append(entries, entry{
                    kind: kind,
                    name: file.Name(),
                })
            }
        }
    }

    return entries, nil
}

type content struct {
    kind    string
    meta    map[string]interface{}
    html    template.HTML
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
            kind:   entry.kind,
            meta:   meta.Get(ctx),
            html:   template.HTML(buf.String()),
        }
    }

    return nil
}

type post struct {
    Title   string
    Slug    string
}

func getPosts() ([]post, error) {
    files, err := os.ReadDir(viper.GetString("ContentPostsDir"))
    if err != nil {
        return nil, err
    }

    posts := make([]post, 0, len(files))
    for _, file := range files {
        if !file.IsDir() {
            filename := strings.Split(file.Name(), ".")[0]
            lang, err := language.Parse(viper.GetString("Language"))
            if err != nil {
                return nil, err
            }

            posts = append(posts, post{
                Title: cases.Title(lang).String(
                    strings.ReplaceAll(filename, "-", " "),
                ),
                Slug: filename,
            })
        }
    }

    return posts, nil
}

type tpl struct {
   kind     string
   content  string
}

func parseTemplates(
    dir, ext string,
    templates map[string]tpl,
    contents map[string]content,
) error {
    entries, err := getFilesByExt(dir, ext)
    if err != nil {
        return err
    }

    files, err := getFilepaths(viper.GetString("PartialsDir"), ".tmpl")
    if err != nil {
        return err
    }

    posts, err := getPosts()
    if err != nil {
        return err
    }

    for _, entry := range entries {
        filename := strings.Split(entry.name, ".")[0]
        if _, ok := contents[filename]; ok {
            files = append(files, fmt.Sprintf("%s/%s", dir, entry.name))
            tmpl, err := template.New(filename).ParseFiles(files...)
            if err != nil {
                return err
            }

            var buf bytes.Buffer
            if err := tmpl.Execute(
                &buf, 
                map[string]interface{}{
                    "meta":     contents[filename].meta,
                    "html":     contents[filename].html,
                    "posts":    posts,
                },
            ); err != nil {
                return err
            }

            templates[filename] = tpl{
                kind: contents[filename].kind,
                content: buf.String(),
            }
        }
    }

    return nil
}

func parsePosts(templates map[string]tpl, contents map[string]content) error {
    files, err := getFilepaths(viper.GetString("PartialsDir"), ".tmpl")
    if err != nil {
        return err
    }

    for filename, content := range contents {
        if content.kind == "content" {
            files = append(
                files, 
                fmt.Sprintf(
                    "%s/%s",
                    viper.GetString("ThemeDir"),
                    "post.html.tmpl",
                ),
            )

            tmpl, err := template.New("post").ParseFiles(files...)
            if err != nil {
                return err
            }

            var buf bytes.Buffer
            if err := tmpl.Execute(
                &buf,
                map[string]interface{}{
                    "meta": content.meta,
                    "html": content.html,
                },
            ); err != nil {
                return err
            }

            templates[filename] = tpl{
                kind: content.kind,
                content: buf.String(),
            }
        }
    }

    return nil
}

func createFiles(templates map[string]tpl) error {
    for filename, tpl := range templates {
        var dir string = viper.GetString("PublicDir")
        if tpl.kind == "content" {
             dir = viper.GetString("PublicPostsDir")
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
    logger := viper.Get("Logger").(*slog.Logger)
    if err := createDir(viper.GetString("PublicPostsDir")); err != nil {
        logger.Error("createDir", "error", err.Error())
        os.Exit(1)
    }

    contents := make(map[string]content)
    err := parseMarkdown(viper.GetString("ContentDir"), ".md", contents)
    if err != nil {
        logger.Error("parseMarkdown", "error", err.Error())
        os.Exit(1)
    }

    err = parseMarkdown(viper.GetString("ContentPostsDir"), ".md", contents)
    if err != nil {
        logger.Error("parseMarkdown", "error", err.Error())
        os.Exit(1)
    }

    templates := make(map[string]tpl)
    err = parseTemplates(viper.GetString("ThemeDir"), ".tmpl", templates, contents)
    if err != nil {
        logger.Error("parseTemplates", "error", err.Error())
        os.Exit(1)
    }

    if err = parsePosts(templates, contents); err != nil {
        logger.Error("parsePosts", "error", err.Error())
        os.Exit(1)
    }

    if err = createFiles(templates); err != nil {
        logger.Error("createFiles", "error", err.Error())
        os.Exit(1)
    }
}
