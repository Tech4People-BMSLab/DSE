
# Create gitignore if does not exist
if (!(Test-Path .gitignore)) {
    Write-Host "Downloading: .gitignore"
    Invoke-WebRequest -Uri 'https://www.toptal.com/developers/gitignore/api/vs,node,visualstudiocode,audio,video,database' -OutFile '.gitignore'
    Add-Content -Path .gitignore -Value "`n.history" -Force
    Write-Host "Created: .gitignore"
}

# if package.json does not exist, create it
if (!(Test-Path package.json)) {
    $json = @{
        name        = "template-ts"
        version     = "1.0.0"
        description = ""
        main        = "src/main.ts"
        keywords    = @()
        author      = ""
        license     = "ISC"
        type        = "module"
        scripts     = @{
            dev     = "nodemon"
        }
    }

    $config = ConvertTo-Json $json -Depth 100
    New-Item -ItemType File -Path "package.json" -Value $config -Force
}

# If nodemon.json does not exist, create it
if (!(Test-Path nodemon.json)) {
    New-Item -ItemType File -Path nodemon.json

    $json = @{
        watch = @("src")
        ext = ".ts,.js"
        exec = "node --loader ts-node/esm --inspect=0.0.0.0:10000 src/main.ts"
    }

    $config = ConvertTo-Json $json -Depth 100
    New-Item -ItemType File -Path "nodemon.json" -Value $config -Force
}

# If tsconfig.json does not exist, create it
if (!(Test-Path tsconfig.json)) {
    New-Item -ItemType File -Path tsconfig.json

    $json = @{
        compilerOptions = @{
            target           = "ES6"
            module           = "ES6"
            lib              = @("es6")
            allowJs          = $true
            outDir           = "build"
            rootDir          = "src"
            strict           = $true
            esModuleInterop  = $true
            moduleResolution = "node"
            types            = @("node", "jest")
            skipLibCheck     = $true
            sourceMap        = $true
        }
        include = @("src/**/*")
        exclude = @("src/**/*.spec.ts")
    }

    $config = ConvertTo-Json $json -Depth 100
    New-Item -ItemType File -Path "tsconfig.json" -Value $config -Force
}

# If src/main.ts does not exist, create it
if (!(Test-Path src/main.ts)) {
    New-Item -ItemType Directory -Path src
    New-Item -ItemType File -Path src/main.ts
    Add-Content -Path src/main.ts -Value "// @ts-nocheck"
}

$packages = @(
    'ts-node'
    '@types/node'
    'lodash'
    '@types/lodash'
    'pretty-error'
)

foreach ($package in $packages) {
    if (!(pnpm list $package)) {
        pnpm add $package
    }
}

pnpm i
nodemon