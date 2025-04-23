package root

//go:generate mockery
//go:generate wire-config -c internal/config/wire/wire.go -o internal/app/server/wire.go -m github.com/danielmesquitta/flight-api/internal/app/server -e dev,staging,test,prod
