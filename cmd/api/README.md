# GraphQL Server

Generate code with :

```shell
cd cmd/graphqlapi
go run github.com/99designs/gqlgen generate
```

- Update `schema.graphqls`
- Complete `schema.resolver.go` to access service layer and `resolver.go` to inject service.