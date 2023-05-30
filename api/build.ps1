try {
    Set-Location ..
    docker build --file api/Dockerfile --tag "registry.gitlab.utwente.nl/bmslabexternal/dse:api-devserver" .
}
finally {
    Set-Location api
}