module github.com/go-scaffold/go-sdk/v2

go 1.23.0

replace github.com/go-scaffold/go-sdk/v2/pkg => ./pkg

require (
	github.com/pasdam/go-template-map-loader v0.0.0-20251027152818-839d0eaea9e2
	github.com/pasdam/go-test-utils v0.0.0-20230710135805-45ec4e440661
	github.com/pasdam/go-utils v0.0.0-20230718144448-c56c396f6c77
	github.com/stretchr/testify v1.11.1
)

require github.com/kr/text v0.2.0 // indirect

require (
	github.com/davecgh/go-spew v1.1.1 // indirect
	github.com/pasdam/files-index v0.0.0-20251027145827-bf1f76a08090
	github.com/pasdam/go-io-utilx v0.0.0-20251027152920-7448902636f4
	github.com/pkg/errors v0.9.1 // indirect
	github.com/pmezard/go-difflib v1.0.0 // indirect
	github.com/stretchr/objx v0.5.2 // indirect
	gopkg.in/yaml.v3 v3.0.1 // indirect
)
