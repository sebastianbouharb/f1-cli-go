# f1-cli-go

## Welcome to the F1 Stats Fetcher ! 🏎️ 🏁 🏆


## Usage
Download the file for your system from the Releases page, then follow the [Install](#install-run-f1-from-anywhere) steps below.

```sh
f1                 # standings for the default season (2026)
f1 --season 2025   # standings for a specific season
```
Mac/Linux:
  chmod +x f1-mac-arm64
  ./f1-mac-arm64 --season (specified season year, ex. 2025)

## Install (run `f1` from anywhere)
To use `f1` from any folder, rename the file to `f1` and put it in a folder on your PATH.

### Mac/Linux
Use the file you downloaded for your system: `f1-mac-arm64` (Apple Silicon), `f1-mac-intel` (Intel Mac), or `f1-linux`. The example below uses `f1-mac-arm64`.

1. From the folder containing the download, move it into `~/.local/bin` as `f1`:
   ```sh
   mkdir -p ~/.local/bin
   mv f1-mac-arm64 ~/.local/bin/f1
   chmod +x ~/.local/bin/f1
   ```
2. Mac only: if it says the developer can't be verified, run:
   ```sh
   xattr -d com.apple.quarantine ~/.local/bin/f1
   ```
3. Add the folder to your PATH. This only needs doing once.

   zsh (the default on Mac):
   ```sh
   echo 'export PATH="$HOME/.local/bin:$PATH"' >> ~/.zshrc
   source ~/.zshrc
   ```
   bash (the default on most Linux systems):
   ```sh
   echo 'export PATH="$HOME/.local/bin:$PATH"' >> ~/.bashrc
   source ~/.bashrc
   ```

### Windows (PowerShell)
1. From the folder containing the download, move it into `bin` in your home folder as `f1.exe`:
   ```powershell
   New-Item -ItemType Directory -Force "$HOME\bin"
   Move-Item .\f1-windows.exe "$HOME\bin\f1.exe"
   ```
2. Add the folder to your PATH. This only needs doing once:
   ```powershell
   [Environment]::SetEnvironmentVariable("Path", [Environment]::GetEnvironmentVariable("Path", "User") + ";$HOME\bin", "User")
   ```
3. Close and reopen PowerShell so the change takes effect.

### Check that it works
```sh
f1 --season 2025
```
If you see `command not found` (or `not recognized` on Windows), open a new terminal and check that the folder is on your PATH.
Windows:
  .\f1-windows.exe --season (specified season year, ex. 2025)
