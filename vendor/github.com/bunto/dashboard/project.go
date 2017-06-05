package dashboard

import (
	"log"
	"sync"
	"time"
)

var (
	defaultProjectMap map[string]*Project
	defaultProjects   = []*Project{
		newProject("bunto", "bunto/bunto", "master", "bunto"),
		newProject("jemoji", "bunto/bemoji", "master", "jemoji"),
		newProject("mercenary", "bunto/mercenary", "master", "mercenary"),
		newProject("bunto-import", "bunto/bunto-import", "master", "bunto-import"),
		newProject("bunto-feed", "bunto/bunto-feed", "master", "bunto-feed"),
		newProject("bunto-seo-tag", "bunto/bunto-seo-tag", "master", "bunto-seo-tag"),
		newProject("bunto-sitemap", "bunto/bunto-sitemap", "master", "bunto-sitemap"),
		newProject("bunto-mentions", "bunto/bunto-mentions", "master", "bunto-mentions"),
		newProject("bunto-watch", "bunto/bunto-watch", "master", "bunto-watch"),
		newProject("bunto-compose", "bunto/bunto-compose", "master", "bunto-compose"),
		newProject("bunto-paginate", "bunto/bunto-paginate", "master", "bunto-paginate"),
		newProject("bunto-gist", "bunto/bunto-gist", "master", "bunto-gist"),
		newProject("bunto-coffeescript", "bunto/bunto-coffeescript", "master", "bunto-coffeescript"),
		newProject("bunto-opal", "bunto/bunto-opal", "master", "bunto-opal"),
		newProject("classifier-reborn", "bunto/classifier-reborn", "master", "classifier-reborn"),
		newProject("bunto-sass-converter", "bunto/bunto-sass-converter", "master", "bunto-sass-converter"),
		newProject("bunto-textile-converter", "bunto/bunto-textile-converter", "master", "bunto-textile-converter"),
		newProject("bunto-redirect-from", "bunto/bunto-redirect-from", "master", "bunto-redirect-from"),
		newProject("github-metadata", "bunto/github-metadata", "master", "bunto-github-metadata"),
		newProject("plugins website", "bunto/plugins", "gh-pages", ""),
		newProject("bunto docker", "bunto/docker", "master", ""),
	}
)

func init() {
	go resetProjectsPeriodically()
}

func resetProjectsPeriodically() {
	for range time.Tick(time.Hour / 2) {
		log.Println("resetting projects' cache")
		resetProjects()
	}
}

func resetProjects() {
	for _, p := range defaultProjects {
		p.reset()
	}
}

type Project struct {
	Name    string `json:"name"`
	Nwo     string `json:"nwo"`
	Branch  string `json:"branch"`
	GemName string `json:"gem_name"`

	Gem     *RubyGem      `json:"gem"`
	Travis  *TravisReport `json:"travis"`
	GitHub  *GitHub       `json:"github"`
	fetched bool
}

func (p *Project) fetch() {
	rubyGemChan := rubygem(p.GemName)
	travisChan := travis(p.Nwo, p.Branch)
	githubChan := github(p.Nwo)

	if p.Gem == nil {
		p.Gem = <-rubyGemChan
	}

	if p.Travis == nil {
		p.Travis = <-travisChan
	}

	if p.GitHub == nil {
		p.GitHub = <-githubChan
	}

	p.fetched = true
}

func (p *Project) reset() {
	p.fetched = false
	p.Gem = nil
	p.Travis = nil
	p.GitHub = nil
}

func buildProjectMap() {
	defaultProjectMap = map[string]*Project{}
	for _, p := range defaultProjects {
		defaultProjectMap[p.Name] = p
	}
}

func newProject(name, nwo, branch, rubygem string) *Project {
	return &Project{
		Name:    name,
		Nwo:     nwo,
		Branch:  branch,
		GemName: rubygem,
	}
}

func getProject(name string) *Project {
	if defaultProjectMap == nil {
		buildProjectMap()
	}

	if p, ok := defaultProjectMap[name]; ok {
		if !p.fetched {
			p.fetch()
		}
		return p
	}

	return nil
}

func getAllProjects() []*Project {
	var wg sync.WaitGroup
	for _, p := range defaultProjects {
		wg.Add(1)
		go func(project *Project) {
			project.fetch()
			wg.Done()
		}(p)
	}
	wg.Wait()
	return defaultProjects
}

func getProjects() []*Project {
	return defaultProjects
}
