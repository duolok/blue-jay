package interfaces

type Engine interface {
    Search()
    LoadScrapers()
    LoadLastSearch()
    UpdateConfig()
}
