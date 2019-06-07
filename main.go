package main

import (
  "context"
  "server_id_api/servers"
  "server_id_api/api"
  "server_id_api/model"
  "github.com/go-chi/chi"
  "github.com/go-chi/render"
  "github.com/go-chi/chi/middleware"
  "time"
  "net/http"
  "encoding/json"
)

var sslLabsClient = api.API{&http.Client{}, "https://api.ssllabs.com"}

func main() {
  r := chi.NewRouter()

  r.Use(middleware.RealIP)
  r.Use(middleware.Logger)
  r.Use(middleware.Recoverer)
  r.Use(middleware.Timeout(60 * time.Second))
  r.Use(render.SetContentType(render.ContentTypeJSON))

  r.Get("/", func(w http.ResponseWriter, r *http.Request) {
    w.Write([]byte("Welcome to server_id_api"))
  })

  r.Route("/domains", func(r chi.Router) {
    r.Get("/", ListDomains)
    r.Route("/search", func(r chi.Router) {
      r.Use(DomainCtx)
      r.Get("/", GetDomain)   // GET /domains/search?domainName=google.com
    })
  })

  http.ListenAndServe(":3000", r)
}

func GetDomain(w http.ResponseWriter, r *http.Request) {
  domain := r.Context().Value("domain").(*model.Domain)

  if err:= render.Render(w, r, NewDomainResponse(domain)); err != nil {
    render.Render(w, r, ErrRender(err))
    return
  }
}

func ListDomains(w http.ResponseWriter, r *http.Request) {
  allDomains := make(map[string]*DomainResponse)
  domains, err := servers.GetAllDomains()

  if err != nil {
    render.Render(w, r, ErrRender(err))
    return
  }

  for domainName, domain := range domains {
    allDomains[domainName] = NewDomainResponse(domain)
  }

  js, _ := json.Marshal(allDomains)

  w.Write(js)
}

func DomainCtx(next http.Handler) http.Handler {
  return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
    var domain model.Domain

    if domainName := r.URL.Query().Get("domainName"); domainName != "" {
      domain = servers.GetServerData(&sslLabsClient, domainName)
    } else {
      render.Render(w, r, ErrNotFound)
      return
    }

    ctx := context.WithValue(r.Context(), "domain", &domain)
    next.ServeHTTP(w, r.WithContext(ctx))
  })
}

// Render Functions For Domains

type DomainResponse struct {
  Servers          []ServerResponse `json:"servers"`
  ServersChanged   bool             `json:"servers_changed"`
  SslGrade         string           `json:"ssl_grade"`
  PreviousSslGrade string           `json:"previous_ssl_Grade,omitempty"`
  Logo             string           `json:"logo"`
  Title            string           `json:"title"`
  IsDown           bool             `json:"is_down"`
}

type ServerResponse struct {
  Address   string `json:"ip_address"`
  SslGrade  string `json:"grade"`
  Country   string `json:"country"`
  Owner     string `json:"owner"`
}

func NewDomainResponse(domain *model.Domain) *DomainResponse {
  servers := []ServerResponse{}

  domainResponse:= &DomainResponse{
    ServersChanged: domain.ServersChanged,
    SslGrade: domain.SslGrade,
    PreviousSslGrade: domain.PreviousSslGrade,
    Logo: domain.Logo,
    Title: domain.Title,
    IsDown: domain.IsDown,
  }

  for index, _ := range domain.Servers {
    server := domain.Servers[index]
    servers = append(servers, NewServerResponse(&server))
  }

  domainResponse.Servers = servers

  return domainResponse
}

func NewServerResponse(server *model.Server) ServerResponse {
  serverResponse := ServerResponse{
    Address: server.Address,
    SslGrade: server.SslGrade,
    Country: server.Country,
    Owner: server.Owner,
  }

  return serverResponse
}

func (rd *DomainResponse) Render(w http.ResponseWriter, r *http.Request) error {
  return nil
}

// Render Functions for Errors

func ErrRender(err error) render.Renderer {
  return &ErrResponse{
    Err:            err,
    HTTPStatusCode: 422,
    StatusText:     "Error rendering response.",
    ErrorText:      err.Error(),
  }
}

type ErrResponse struct {
  Err             error `json:"-"` // low-level runtime error
  HTTPStatusCode  int   `json:"-"` // http response status code

  StatusText      string `json:"status"`          // user-level status message
  AppCode         int64  `json:"code,omitempty"`  // application-specific error code
  ErrorText       string `json:"error,omitempty"` // application-level error message, for debugging
}

func (e *ErrResponse) Render(w http.ResponseWriter, r *http.Request) error {
  render.Status(r, e.HTTPStatusCode)
  return nil
}

var ErrNotFound = &ErrResponse{HTTPStatusCode: 404, StatusText: "Resource not found."}
