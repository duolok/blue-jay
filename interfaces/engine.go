package interfaces

type Engine interface {
    search()
    loadScrapers()
    loadLastSearch()
    updateConfig()
}
