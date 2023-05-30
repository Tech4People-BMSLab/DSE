# Log in to Docker Registry
Write-Host "Logging in..." -ForegroundColor Green
docker login registry.gitlab.utwente.nl -u registry-manager -p Jwo4JKWpkn12mz9k2Vvz

Set-Location ..

# Print directory
Write-Host "Current directory: " -ForegroundColor Green -NoNewline
Write-Host $PWD -ForegroundColor Green

# Build and push api image
Write-Host "Building api image..." -ForegroundColor Green
docker build --file api/Dockerfile --tag "registry.gitlab.utwente.nl/bmslabexternal/dse:api-devserver" .
Write-Host "Pushing api image..." -ForegroundColor Green
docker push "registry.gitlab.utwente.nl/bmslabexternal/dse:api-devserver"

# Complete
Write-Host "Done!" -ForegroundColor Green

Set-Location .\docker