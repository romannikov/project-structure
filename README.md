# CLI for a project structure creation

It's a CLI for an initial project structure creation. 

## Examples

There is only one required parameter: it's a project name (`project-name`).
```
project-structure --project-name example
```

By default CLI creates a go-like project structure that looks like:
```
project-name
  api/
  build/
  cmd/
    main.go
  deploy/
  docs/
  internal/
    project-name/
      rest/
      service/
  pkg/
```
Also, it can be customize by a `project-structure` parameter that should contain a path ot the file with the structure.
File example:
```
{
  "items":[
    {
      "name":"api",
      "is_directory":true
    },
    {
      "name":"build",
      "is_directory":true
    },
    {
      "name":"cmd",
      "is_directory":true,
      "items":[
        {
          "name":"main.go",
          "is_directory":false
        }
      ]
    }
  ]
}
```
