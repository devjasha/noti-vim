" noti.vim - Vim plugin for Noti notes
" Maintainer: devjasha
" Version: 0.1.0

if exists('g:loaded_noti')
  finish
endif
let g:loaded_noti = 1

" Configuration
let g:noti_notes_dir = get(g:, 'noti_notes_dir', expand('~/notes'))
let g:noti_default_folder = get(g:, 'noti_default_folder', '')
let g:noti_default_tags = get(g:, 'noti_default_tags', [])
let g:noti_git_auto_commit = get(g:, 'noti_git_auto_commit', 0)

" Check if noti CLI is available
function! s:CheckNotiCLI()
  if !executable('noti')
    echoerr 'noti CLI not found. Install it with: go install github.com/devjasha/noti-vim/cmd/noti@latest'
    return 0
  endif
  return 1
endfunction

" Create a new note
function! noti#New(...)
  if !s:CheckNotiCLI()
    return
  endif

  let l:title = a:0 > 0 ? a:1 : input('Note title: ')
  if empty(l:title)
    echo 'Cancelled'
    return
  endif

  let l:cmd = 'noti new "' . l:title . '"'

  if !empty(g:noti_default_folder)
    let l:cmd .= ' --folder ' . g:noti_default_folder
  endif

  if !empty(g:noti_default_tags)
    let l:cmd .= ' --tags ' . join(g:noti_default_tags, ',')
  endif

  let l:output = system(l:cmd . ' --json')
  let l:result = json_decode(l:output)

  if v:shell_error == 0
    execute 'edit ' . l:result.file_path
    echo 'Created note: ' . l:result.title
  else
    echoerr 'Failed to create note: ' . l:output
  endif
endfunction

" List all notes
function! noti#List(...)
  if !s:CheckNotiCLI()
    return
  endif

  let l:cmd = 'noti list --json'

  if a:0 > 0 && !empty(a:1)
    let l:cmd .= ' --folder ' . a:1
  endif

  let l:output = system(l:cmd)
  if v:shell_error != 0
    echoerr 'Failed to list notes: ' . l:output
    return
  endif

  let l:notes = json_decode(l:output)

  " Create a new buffer for the list
  new
  setlocal buftype=nofile
  setlocal bufhidden=wipe
  setlocal noswapfile
  setlocal nowrap
  setlocal cursorline

  " Add header
  call setline(1, 'Noti Notes (' . len(l:notes) . ' total)')
  call setline(2, repeat('=', 80))

  let l:line = 3
  for note in l:notes
    let l:tags = empty(note.tags) ? '' : ' [' . join(note.tags, ', ') . ']'
    call setline(l:line, note.title . l:tags)
    let l:line += 1
    call setline(l:line, '  ' . note.slug)
    let l:line += 1
    call setline(l:line, '')
    let l:line += 1
  endfor

  " Make it readonly
  setlocal nomodifiable
  setlocal readonly

  " Add keybinding to open note
  nnoremap <buffer> <CR> :call <SID>OpenNoteFromList()<CR>
  nnoremap <buffer> q :close<CR>

  " Store notes data for later use
  let b:noti_notes = l:notes
endfunction

" Open note from list buffer
function! s:OpenNoteFromList()
  if !exists('b:noti_notes')
    return
  endif

  let l:line = getline('.')

  " Find the note by slug (lines starting with spaces contain slugs)
  if l:line =~ '^\s\+'
    let l:slug = substitute(l:line, '^\s\+', '', '')

    for note in b:noti_notes
      if note.slug == l:slug
        close
        execute 'edit ' . note.file_path
        return
      endif
    endfor
  endif
endfunction

" Search notes
function! noti#Search(query)
  if !s:CheckNotiCLI()
    return
  endif

  if empty(a:query)
    let l:query = input('Search query: ')
    if empty(l:query)
      echo 'Cancelled'
      return
    endif
  else
    let l:query = a:query
  endif

  let l:output = system('noti search "' . l:query . '" --json')
  if v:shell_error != 0
    echoerr 'Search failed: ' . l:output
    return
  endif

  let l:results = json_decode(l:output)

  if empty(l:results)
    echo 'No matches found for: ' . l:query
    return
  endif

  " Create a new buffer for search results
  new
  setlocal buftype=nofile
  setlocal bufhidden=wipe
  setlocal noswapfile
  setlocal nowrap
  setlocal cursorline

  " Add header
  call setline(1, 'Search Results for: "' . l:query . '" (' . len(l:results) . ' matches)')
  call setline(2, repeat('=', 80))

  let l:line = 3
  for result in l:results
    call setline(l:line, result.note.title)
    let l:line += 1
    call setline(l:line, '  ' . result.note.slug)
    let l:line += 1

    for match in result.matches
      if match.context == 'title'
        call setline(l:line, '    • in title: ' . match.line)
      elseif match.context == 'tag'
        call setline(l:line, '    • in tag: ' . match.line)
      else
        call setline(l:line, '    • line ' . match.line_number . ': ' . match.line)
      endif
      let l:line += 1
    endfor

    call setline(l:line, '')
    let l:line += 1
  endfor

  " Make it readonly
  setlocal nomodifiable
  setlocal readonly

  " Add keybinding to open note
  nnoremap <buffer> <CR> :call <SID>OpenNoteFromSearch()<CR>
  nnoremap <buffer> q :close<CR>

  " Store results for later use
  let b:noti_search_results = l:results
endfunction

" Open note from search results
function! s:OpenNoteFromSearch()
  if !exists('b:noti_search_results')
    return
  endif

  let l:line = getline('.')

  " Find note by slug
  if l:line =~ '^\s\+[a-z0-9/-]\+$'
    let l:slug = substitute(l:line, '^\s\+', '', '')

    for result in b:noti_search_results
      if result.note.slug == l:slug
        close
        execute 'edit ' . result.note.file_path
        return
      endif
    endfor
  endif
endfunction

" List tags
function! noti#Tags()
  if !s:CheckNotiCLI()
    return
  endif

  let l:output = system('noti tags --json')
  if v:shell_error != 0
    echoerr 'Failed to list tags: ' . l:output
    return
  endif

  let l:tags = json_decode(l:output)

  " Create a new buffer for tags
  new
  setlocal buftype=nofile
  setlocal bufhidden=wipe
  setlocal noswapfile
  setlocal nowrap
  setlocal cursorline

  call setline(1, 'Tags (' . len(l:tags) . ' total)')
  call setline(2, repeat('=', 80))

  let l:line = 3
  for tag in l:tags
    call setline(l:line, printf('%-30s (%d notes)', tag.tag, tag.count))
    let l:line += 1
  endfor

  setlocal nomodifiable
  setlocal readonly

  nnoremap <buffer> <CR> :call <SID>FilterByTag()<CR>
  nnoremap <buffer> q :close<CR>

  let b:noti_tags = l:tags
endfunction

" Filter notes by tag
function! s:FilterByTag()
  if !exists('b:noti_tags')
    return
  endif

  let l:line = getline('.')
  let l:tag = matchstr(l:line, '^[a-zA-Z0-9_-]\+')

  if !empty(l:tag)
    close
    call noti#ListByTag(l:tag)
  endif
endfunction

" List notes by tag
function! noti#ListByTag(tag)
  if !s:CheckNotiCLI()
    return
  endif

  let l:output = system('noti list --tag ' . a:tag . ' --json')
  if v:shell_error != 0
    echoerr 'Failed to list notes: ' . l:output
    return
  endif

  let l:notes = json_decode(l:output)

  " Create buffer and display notes
  new
  setlocal buftype=nofile
  setlocal bufhidden=wipe
  setlocal noswapfile
  setlocal nowrap
  setlocal cursorline

  call setline(1, 'Notes tagged with "' . a:tag . '" (' . len(l:notes) . ' total)')
  call setline(2, repeat('=', 80))

  let l:line = 3
  for note in l:notes
    call setline(l:line, note.title)
    let l:line += 1
    call setline(l:line, '  ' . note.slug)
    let l:line += 1
    call setline(l:line, '')
    let l:line += 1
  endfor

  setlocal nomodifiable
  setlocal readonly

  nnoremap <buffer> <CR> :call <SID>OpenNoteFromList()<CR>
  nnoremap <buffer> q :close<CR>

  let b:noti_notes = l:notes
endfunction

" List folders
function! noti#Folders()
  if !s:CheckNotiCLI()
    return
  endif

  let l:output = system('noti folders --json')
  if v:shell_error != 0
    echoerr 'Failed to list folders: ' . l:output
    return
  endif

  let l:folders = json_decode(l:output)

  " Create a new buffer for folders
  new
  setlocal buftype=nofile
  setlocal bufhidden=wipe
  setlocal noswapfile
  setlocal nowrap
  setlocal cursorline

  call setline(1, 'Folders (' . len(l:folders) . ' total)')
  call setline(2, repeat('=', 80))

  let l:line = 3
  for folder in l:folders
    call setline(l:line, printf('%-40s (%d notes)', folder.path, folder.count))
    let l:line += 1
  endfor

  setlocal nomodifiable
  setlocal readonly

  nnoremap <buffer> <CR> :call <SID>FilterByFolder()<CR>
  nnoremap <buffer> q :close<CR>

  let b:noti_folders = l:folders
endfunction

" Filter notes by folder
function! s:FilterByFolder()
  if !exists('b:noti_folders')
    return
  endif

  let l:line = getline('.')
  let l:folder = matchstr(l:line, '^[a-zA-Z0-9/_-]\+')

  if !empty(l:folder)
    close
    call noti#ListByFolder(l:folder)
  endif
endfunction

" List notes by folder
function! noti#ListByFolder(folder)
  if !s:CheckNotiCLI()
    return
  endif

  let l:output = system('noti list --folder ' . a:folder . ' --json')
  if v:shell_error != 0
    echoerr 'Failed to list notes: ' . l:output
    return
  endif

  let l:notes = json_decode(l:output)

  " Create buffer and display notes
  new
  setlocal buftype=nofile
  setlocal bufhidden=wipe
  setlocal noswapfile
  setlocal nowrap
  setlocal cursorline

  call setline(1, 'Notes in folder "' . a:folder . '" (' . len(l:notes) . ' total)')
  call setline(2, repeat('=', 80))

  let l:line = 3
  for note in l:notes
    call setline(l:line, note.title)
    let l:line += 1
    call setline(l:line, '  ' . note.slug)
    let l:line += 1
    call setline(l:line, '')
    let l:line += 1
  endfor

  setlocal nomodifiable
  setlocal readonly

  nnoremap <buffer> <CR> :call <SID>OpenNoteFromList()<CR>
  nnoremap <buffer> q :close<CR>

  let b:noti_notes = l:notes
endfunction

" Git operations
function! noti#GitStatus()
  if !s:CheckNotiCLI()
    return
  endif

  let l:output = system('noti git status')
  echo l:output
endfunction

function! noti#GitCommit(...)
  if !s:CheckNotiCLI()
    return
  endif

  let l:message = a:0 > 0 ? a:1 : input('Commit message: ')
  if empty(l:message)
    echo 'Cancelled'
    return
  endif

  let l:output = system('noti git commit "' . l:message . '"')
  echo l:output
endfunction

function! noti#GitSync(...)
  if !s:CheckNotiCLI()
    return
  endif

  let l:message = a:0 > 0 ? a:1 : ''
  let l:cmd = 'noti git sync'

  if !empty(l:message)
    let l:cmd .= ' --message "' . l:message . '"'
  endif

  echo 'Syncing...'
  let l:output = system(l:cmd)
  echo l:output
endfunction

" Commands
command! -nargs=? NotiNew call noti#New(<f-args>)
command! -nargs=? NotiList call noti#List(<f-args>)
command! -nargs=? NotiSearch call noti#Search(<q-args>)
command! NotiTags call noti#Tags()
command! NotiFolders call noti#Folders()
command! NotiGitStatus call noti#GitStatus()
command! -nargs=? NotiGitCommit call noti#GitCommit(<f-args>)
command! -nargs=? NotiGitSync call noti#GitSync(<f-args>)

" Default keymappings (can be disabled by setting g:noti_no_default_mappings = 1)
if !get(g:, 'noti_no_default_mappings', 0)
  nnoremap <leader>nn :NotiNew<CR>
  nnoremap <leader>nl :NotiList<CR>
  nnoremap <leader>ns :NotiSearch<Space>
  nnoremap <leader>nt :NotiTags<CR>
  nnoremap <leader>nf :NotiFolders<CR>
  nnoremap <leader>ng :NotiGitStatus<CR>
  nnoremap <leader>nc :NotiGitCommit<Space>
  nnoremap <leader>np :NotiGitSync<CR>
endif
