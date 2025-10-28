# Noti Vim - CLI and Vim Plugin for Noti Notes

A fast, standalone CLI tool and Vim/Neovim plugin for managing markdown notes compatible with the [Noti desktop app](https://github.com/devjasha/Noti).

## Features

- ✅ **Fast Go CLI** - Single binary, no runtime dependencies
- ✅ **Vim/Neovim integration** - Edit notes without leaving your editor
- ✅ **Git integration** - Built-in version control
- ✅ **Tag and folder organization** - Organize notes your way
- ✅ **Full-text search** - Find notes quickly
- ✅ **Compatible with Noti app** - Use both tools seamlessly

## Installation

### CLI Tool

#### From Source
```bash
go install github.com/devjasha/noti-vim/cmd/noti@latest

# Make sure ~/go/bin is in your PATH
export PATH="$PATH:$HOME/go/bin"
```

**Note**: The binary will be installed to `~/go/bin/noti`. Add `~/go/bin` to your PATH if not already present.

#### From Binary
Download the latest release for your platform from [GitHub Releases](https://github.com/devjasha/noti-vim/releases).

```bash
# Linux/macOS
curl -L https://github.com/devjasha/noti-vim/releases/latest/download/noti-$(uname -s)-$(uname -m) -o noti
chmod +x noti
sudo mv noti /usr/local/bin/
```

### Vim Plugin

#### Using vim-plug
```vim
Plug 'devjasha/noti-vim'
```

#### Using packer.nvim
```lua
use 'devjasha/noti-vim'
```

#### Using lazy.nvim
```lua
{
  'devjasha/noti-vim',
  dependencies = {
    'nvim-telescope/telescope.nvim',  -- Optional: for fuzzy finding
  },
  config = function()
    require('noti').setup()
  end
}
```

## Quick Start

### 1. Initialize Notes Directory
```bash
noti init ~/notes
```

### 2. Create Your First Note
```bash
noti new "My First Note"
```

### 3. List Notes
```bash
noti list
```

### 4. Open in Vim
```bash
vim $(noti list --json | jq -r '.[0].file_path')
```

## CLI Usage

### Note Operations

```bash
# Create a new note
noti new "Meeting Notes" --folder meetings --tags work,important

# List all notes
noti list

# List notes in a folder
noti list --folder projects

# List notes with a tag
noti list --tag work

# Show a note
noti show meetings/meeting-notes

# Delete a note
noti delete meetings/meeting-notes
```

### Search

```bash
# Search note content
noti search "project deadline"

# Search with JSON output
noti search "TODO" --json
```

### Organization

```bash
# List all tags
noti tags

# List all folders
noti folders
```

### Git Operations

```bash
# Commit changes
noti commit "Updated project notes"

# Sync with remote
noti sync

# Check status
noti status
```

### Configuration

```bash
# View config
cat ~/.config/noti/config.yaml

# Set notes directory
noti init /path/to/notes
```

## Vim Plugin Usage

### Commands

| Command | Description |
|---------|-------------|
| `:NotiNew [name]` | Create a new note |
| `:NotiFind` | Fuzzy find notes (Telescope) |
| `:NotiList` | List all notes |
| `:NotiSearch <query>` | Search note content |
| `:NotiTags` | Browse by tags |
| `:NotiFolders` | Browse folders |
| `:NotiCommit [msg]` | Commit changes |
| `:NotiSync` | Sync with remote |
| `:NotiStatus` | Git status |

### Default Keybindings

| Key | Action |
|-----|--------|
| `<leader>nn` | New note |
| `<leader>nf` | Find notes |
| `<leader>ns` | Search content |
| `<leader>nt` | Browse tags |
| `<leader>nd` | Browse folders |
| `<leader>nc` | Commit |
| `<leader>np` | Sync |

### Configuration

```vim
" In your vimrc
let g:noti_notes_dir = '~/notes'
let g:noti_default_folder = ''
let g:noti_default_tags = []
let g:noti_git_auto_commit = 0

" Custom keybindings
nmap <leader>n <Plug>NotiNew
nmap <leader>f <Plug>NotiFindFiles
```

For Neovim with Lua:
```lua
require('noti').setup({
  notes_dir = vim.fn.expand('~/notes'),
  default_folder = '',
  default_tags = {},
  git_auto_commit = false,

  -- Telescope integration
  telescope = {
    enabled = true,
  },

  -- Keybindings
  mappings = {
    new = '<leader>nn',
    find = '<leader>nf',
    search = '<leader>ns',
    tags = '<leader>nt',
    folders = '<leader>nd',
    commit = '<leader>nc',
    sync = '<leader>np',
  }
})
```

## File Format

Notes are stored as markdown files with YAML frontmatter:

```markdown
---
title: My Note Title
tags: [work, important]
created: 2025-10-28T12:00:00Z
---

# My Note Title

This is the content of the note.
```

## Development

### Building

```bash
# Build for current platform
make build

# Build for all platforms
make build-all

# Install locally
make install
```

### Testing

```bash
# Run tests
go test ./...

# Run with coverage
go test -cover ./...
```

## Compatibility

- **Noti Desktop App**: Fully compatible - use both tools on the same notes directory
- **File Format**: Same markdown + YAML frontmatter
- **Git Repository**: Shared git history
- **Sync**: File watchers handle automatic updates

## License

ISC

## Contributing

Contributions welcome! Please open an issue or pull request.
