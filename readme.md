This repo is just me messing around with project structure, organization and approaches.

# Overview

 - Go Version 1.21
 - Multiple entry points (cmd):
   - `api` Service
   - `migrations` Service
   - `seed` Service
   - `stub` (to run locally such as a script, reusing code)
 - Shared `internal` modules
 - Dependency Injection ([uber-go/dig](https://github.com/uber-go/dig))
 - Logger ([slog](https://pkg.go.dev/log/slog))
 - Live Reload ([cosmtrek/air](https://github.com/cosmtrek/air)) (TODO: not yet inside container)
 - Makefile
 - Semi-automatic Database Transaction handling
 
For list of dependencies, check [go.mod](https://github.com/rafael1mc/go-structure/blob/main/go.mod).

# Important

I have commited `.env` and `keys` just to simplify the demonstration. They should NOT be commited in your _real_ project.

# How to run

To run the project, you have to do the following:
 1. Run compose
 2. Run migrations
 3. Run seed (optional)


### 1. Run compose:

```
docker compose up --build
```

### 2. Run Migrations
To run migrations, you have 3 options. **Pick one** of them:

#### Option 1: Run `migration` service through vscode:
1. Hit F5 (run on vscode)
2. Select `migration` and press enter

#### Option 2: Run `migration` manually:
```
go run ./cmd/migration
```

#### Option 3: Run migrations through make:
```
make migration-up
```

### 3. Run Seed
To seed the database:

1. Hit F5
2. Select `seed` and press enter

(You can also run it manually `go run ./cmd/seed`)

# Manual Test
Launch [localhost](http://localhost:8080/ping) to see if it works.

Data won't persist between database executions (missing volume, so if the database container dies, the data will be lost).

You can launch an HTTP Client and make the following request to login:
```
curl --request POST \
  --url http://localhost:8080/auth \
  --header 'Content-Type: application/json' \
  --data '{
	"email": "admin@example.com",
	"password": "admin123"
}'
```

# License
Check at [LICENSE](https://github.com/rafael1mc/go-structure/blob/main/LICENSE).
