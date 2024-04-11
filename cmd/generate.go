package cmd

import (
	"bytes"
	"flag"
	"fmt"
	"html/template"
	"os"
	"path/filepath"
	"strings"

	"github.com/yuin/goldmark"
	"github.com/yuin/goldmark-meta"
	"github.com/yuin/goldmark/extension"
	"github.com/yuin/goldmark/parser"
	"golang.org/x/text/cases"
	"golang.org/x/text/language"
)

const (
    ContentDir      string = "content"
    PublicDir       string = "public"
    ThemesDir       string = "themes"
    DefaultTheme    string = "https://github.com/kevinsuner/elmo-thumbalina.git"
)

var (
    contentPostsDir string = fmt.Sprintf("%s/posts", ContentDir)
    publicPostsDir  string = fmt.Sprintf("%s/posts", PublicDir)
    themeDir        string
    partialsDir     string
)

func init() {
    commander["generate"] = Command{
        Flags: flag.NewFlagSet("generate", flag.ExitOnError),
        Use: `Transforms the Markdown files found in the content directory
to HTML, and embeds the resulting output into a Golang
html/template.

It uses the files provided by the theme (which is set in the
configuration file <elmo.toml>), in the html/template parsing
and execution phase.`,
        Run: generate,
    }
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
    if dir == contentPostsDir {
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
    files, err := os.ReadDir(contentPostsDir)
    if err != nil {
        return nil, err
    }

    posts := make([]post, 0, len(files))
    for _, file := range files {
        if !file.IsDir() {
            filename := strings.Split(file.Name(), ".")[0]
            lang, err := language.Parse("en")
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

    files, err := getFilepaths(partialsDir, ".tmpl")
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
                kind:       contents[filename].kind,
                content:    buf.String(),
            }
        }
    }

    return nil
}

func parsePosts(templates map[string]tpl, contents map[string]content) error {
    files, err := getFilepaths(partialsDir, ".tmpl")
    if err != nil {
        return err
    }

    for filename, content := range contents {
        if content.kind == "content" {
            files = append(files, fmt.Sprintf("%s/%s", themeDir, "post.html.tmpl"))
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
                kind:       content.kind,
                content:    buf.String(),
            }
        }
    }

    return nil
}

func createFiles(dir string, templates map[string]tpl) error {
    for filename, tpl := range templates {
        if tpl.kind == "content" {
             dir = publicPostsDir
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

func generate(params map[string]string) {
    themeDir = fmt.Sprintf("%s/%s", ThemesDir, params["theme"])
    partialsDir = fmt.Sprintf("%s/partials", themeDir)

    if err := createDir(publicPostsDir); err != nil {
        logger.Error("createDir", "error", err.Error())
        os.Exit(1)
    }

    contents := make(map[string]content)
    err := parseMarkdown(ContentDir, ".md", contents)
    if err != nil {
        logger.Error("parseMarkdown", "error", err.Error())
        os.Exit(1)
    }

    err = parseMarkdown(contentPostsDir, ".md", contents)
    if err != nil {
        logger.Error("parseMarkdown", "error", err.Error())
        os.Exit(1)
    }

    templates := make(map[string]tpl)
    err = parseTemplates(themeDir, ".tmpl", templates, contents)
    if err != nil {
        logger.Error("parseTemplates", "error", err.Error())
        os.Exit(1)
    }

    if err = parsePosts(templates, contents); err != nil {
        logger.Error("parsePosts", "error", err.Error())
        os.Exit(1)
    }

    if err = createFiles(PublicDir, templates); err != nil {
        logger.Error("createFiles", "error", err.Error())
        os.Exit(1)
    }
}
