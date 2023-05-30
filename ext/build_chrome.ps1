# Install dependencies using pnpm
pnpm i

# Build the app for Chrome using pnpm
cross-env BROWSER='chrome' webpack --config webpack/webpack.dev.js

# Change directory to the 'dist' folder
Set-Location "dist"

# Compress the contents of the 'dist' folder using 7-Zip
& "C:\Program Files\7-Zip\7z.exe" a -tzip "DSE_Chrome.zip" "*"

# Change directory back to the root folder
Set-Location ".."

# Print a message indicating that the build is complete
Write-Host "Build complete"