# proto-lib

Generated protobuf/gRPC Go stubs shared by all services.

- Module path: `github.com/rusgainew/tunduck-app-mk/proto-lib`
- Regenerate stubs from the repo root:
  - `cd api/proto && make`
- Generated packages:
  - `common` for shared messages
  - `auth` for auth messages and `AuthService`
  - `company` for company messages and `CompanyService`
  - `document` for document messages and `DocumentService`

If you use a Go workspace locally, `go.work` already includes `./proto-lib` so services import the local copy during development.
