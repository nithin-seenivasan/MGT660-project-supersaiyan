# MGT 660 - Advanced Management of Software Development
2nd year elective at Yale School of Management, as part of the MBA program
## Final Project - Coding a basic Eventbrite Clone in Golang
A fully functioning but basic Eventbrite clone written in Go and hosted in Heroku. Sprint reports are located here https://github.com/nithin-seenivasan/project-report-supersaiyan 


## What is here

| File                      | Role                                                                                                                      |
| ------------------------- | ------------------------------------------------------------------------------------------------------------------------- |
| ./README.md               | This file!                                                                                                                |
| ./server.go               | Entrypoint for the code; contains the `main` function                                                                     |
| ./routes.go               | Maps your URLs to the controllers/handlers for each URL                                                                   |
| ./event_models.go         | Defines your data structure and logic. I put in a few default events.                                                     |
| ./index_controllers.go    | Controllers related to the index (home) page                                                                              |
| ./templates.go            | Sets up your templates from which you generate HTML                                                                       |
| ./templates               | Directory with your templates. You'll need more templates ;P                                                              |
| ./templates/layout.gohtml | The "base" layout for your HTML frontend                                                                                  |
| ./templates/index.gohtml  | The template for the index (home) page                                                                                    |
| ./static.go               | Sets up the static file server (see next entry)                                                                           |
| ./staticfiles             | Directory with our "static" assets like images, css, js                                                                   |
| ./go.mod                  | [Go modules file](https://www.kablamo.com.au/blog/2018/12/10/just-tell-me-how-to-use-go-modules). Lists our dependencies. |
| ./go.sum                  | A "checksum" file that says precisely what versions of our dependencies need to be installed.                             |
| ./vendor                  | A directory containing our dependencies                                                                                   |
| ./init-schema.sql         | An SQL file, which is imported by `event_models.go`, that sets up the database when the app/server is started             |
