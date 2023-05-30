$global_packages = @(
    'cross-env'
    'webpack'
    'webpack-cli'
    'web-ext'
)

foreach ($package in $global_packages) {
    if (!(pnpm list $package)) {
        npm i -g $package
    }
}

try {
    pnpm i
    cross-env BROWSER='chrome' webpack --config webpack/webpack.dev.js --watch &

    Start-Sleep -Seconds 10

    New-Item -ItemType Directory -Force -Path ".\dist"
    Copy-Item -Force ".\public\firefox_manifest.json" ".\dist\manifest.json"
    Set-Location "dist"
    web-ext run
}
finally {
    Set-Location ".."
}