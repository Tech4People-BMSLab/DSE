# Make sure to run build_chrome.ps1 first!
try {
    # Create directory 'dist-firefox' if it doesn't exist
    if (!(Test-Path "dist-firefox")) {
        New-Item -ItemType Directory -Path "dist-firefox"
    }

    if (!(Test-Path "dist-firefox/images")) {
        New-Item -ItemType Directory -Path "dist-firefox/images"
    }

    # Change directory to the 'dist-firefox' folder
    Set-Location "dist-firefox"

    # Remove the 'images' folder recursively and forcefully
    Remove-Item -Recurse -Force "images"

    # Delete all CSS, HTML, JS, and ZIP files in the current directory
    Remove-Item -Force "*.css", "*.html", "*.js", "*.zip"

    # Change directory back to the root folder
    Set-Location ".."

    # Install dependencies using pnpm
    pnpm i

    # Build the app for Chrome using pnpm
    cross-env BROWSER='chrome' webpack --config webpack/webpack.dev.js

    # Change directory to the 'dist' folder
    Set-Location "dist"

    # Copy the 'images' folder to the 'dist-firefox' folder
    Copy-Item -Recurse -Force "images" "..\dist-firefox\images\"

    # Copy all CSS, HTML, and JS files to the 'dist-firefox' folder
    Copy-Item -Force "*.css", "*.html", "*.js" "..\dist-firefox\"

    # Copy the 'firefox_manifest.json' file to the 'dist-firefox' folder
    Copy-Item -Force "..\public\firefox_manifest.json" "..\dist-firefox\manifest.json"

    # Change directory to the 'dist-firefox' folder
    Set-Location "..\dist-firefox"

    # Compress the contents of the 'dist-firefox' folder using 7-Zip
    & "C:\Program Files\7-Zip\7z.exe" a -tzip "DSE_Firefox.zip" "*"

    # Change directory to the src folder
    Set-Location "../src"

    # Compress the contents of the 'dist-firefox' folder using 7-Zip
    & "C:\Program Files\7-Zip\7z.exe" a -tzip "../dist-firefox/DSE_Source.zip" "*"
} finally {
    # Change directory back to the root folder
    Set-Location ".."

    # Print a message indicating that the build is complete
    Write-Host "Build complete"
}