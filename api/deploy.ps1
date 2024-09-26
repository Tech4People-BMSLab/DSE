# Set GOOS and GOARCH for Linux build
$env:CGO_ENABLED = 0
$env:GOGC        = 50
$env:GOOS        = "linux"
$env:GOARCH      = "amd64"

$dateTime = Get-Date -Format "yyyy-MM-dd HH:mm"
Set-Content -Path "build.txt" -Value $dateTime

# Build the web application
Push-Location "../web"
Invoke-Expression "bun run build"
Pop-Location

# Build the application
go build -trimpath -gcflags="all=-l -B" -ldflags="-s -w -extldflags '-static'" -o ./dist/main ./src/main.go

# Reset the environment variables if needed
Remove-Item Env:\GOOS
Remove-Item Env:\GOARCH

docker build -t registry.gitlab.utwente.nl/bmslabexternal/dse:latest .

docker login registry.gitlab.utwente.nl -u portainer -p pYnp8zViCYCcbA6hizfb
docker push registry.gitlab.utwente.nl/bmslabexternal/dse:latest
