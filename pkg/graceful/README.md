# graceful

Utilitaire Go pour gérer l'arrêt propre (graceful shutdown) d'une application via les signaux système (ex: `SIGINT`, `SIGTERM`) avec `context.Context`.

## ✨ Fonctionnalités

- Intercepte les signaux système (`os.Interrupt`, `syscall.SIGTERM`)
- Annule un `context.Context` partagé à travers l'application
- Permet une fermeture propre des goroutines, workers, clients, etc.

## Example 

```go
func main() {
    ctx, cancel := graceful.WaitForSignalContext(context.Background())

    fmt.Println("App running. Press Ctrl+C to exit.")

    <-ctx.Done()
    fmt.Println("Context cancelled. Shutting down gracefully.")
}
```

## 🧠 Bonnes pratiques
- Partagez le ctx retourné dans tous vos composants pour qu'ils puissent s'arrêter proprement.
- Utilisez `<-ctx.Done()` dans vos workers, serveurs, ou boucles pour détecter un arrêt.
- Combinez avec `context.WithTimeout()` si vous avez besoin d'un délai d’arrêt limité.
