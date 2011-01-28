package main

import (
    "log"
    "strings"
    "strconv"
    "http"
)

type ViewCtx struct {
    left, right interface{}
}

// Render main page
func show(wr http.ResponseWriter, art_num string) {
    id, _ := strconv.Atoi(art_num)
    main_view.Exec(wr, ViewCtx{getArticleList(), getArticle(id)})
}

// Render edit page
func edit(wr http.ResponseWriter, art_num string) {
    id, _ := strconv.Atoi(art_num)
    edit_view.Exec(wr, ViewCtx{getArticleList(), getArticle(id)})
}

// Update database and render main page
func update(wr http.ResponseWriter, req *http.Request, art_num string) {
    if req.FormValue("submit") == "Save" {
        id, _ := strconv.Atoi(art_num) // id == 0 means new article
        id = updateArticle(
            id, req.FormValue("title"), req.FormValue("body"),
        )
        // If we insert new article, we change art_num to its id. This allows
        // show the article immediately after its creation.
        art_num = strconv.Itoa(id)
    }
    // Show modified/created article
    show(wr, art_num)
}

// Decide which handler to use basis on the request method and URL path.
func router(wr http.ResponseWriter, req *http.Request) {
    root_path  := "/"
    edit_path  := "/edit/"

    switch req.Method {
    case "GET":
        switch {
        case req.URL.Path == "/style.css":
            http.ServeFile(wr, req, "static" + req.URL.Path)

        case strings.HasPrefix(req.URL.Path, edit_path):
            edit(wr, req.URL.Path[len(edit_path):])

        case strings.HasPrefix(req.URL.Path, root_path):
            show(wr, req.URL.Path[len(root_path):])
        }

    case "POST":
        switch {
        case strings.HasPrefix(req.URL.Path, root_path):
            update(wr, req, req.URL.Path[len(root_path):])
        }
    }
}

func main() {
    viewInit()
    mysqlInit()

    err := http.ListenAndServe(":1111", http.HandlerFunc(router))
    if err != nil {
        log.Exitln("ListenAndServe:", err)
    }
}
