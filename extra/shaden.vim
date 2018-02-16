" To install:
"
"   # mkdir -p $HOME/.vim/plugin
"   # cp extra/shaden.vim $HOME/.vim/plugin/shaden.vim
"
" Add these lines to your vimrc:
"
"   " Send a visual block of code to Shaden
"   vnoremap <leader>p :<C-U>ShadenPatchSelection<CR>
"
"   " Send a line of code to Shaden
"   nnoremap <leader>p :<C-U>ShadenPatchLine<CR>
"
"   " Clear the patch and resend the entire file to Shaden
"   nnoremap <leader>r :<C-U>ShadenRepatch<CR>

if (exists("g:loaded_shaden"))
    finish
endif
let g:loaded_shaden = 1

let s:error_message = "unable to communicate with shaden"

if (!exists("g:shaden_http_addr"))
    let g:shaden_http_addr = '127.0.0.1:5000'
endif

function! ShadenRepatch()
    let result = system(s:command("(clear)"))
    if v:shell_error
        echo s:error_message
        return
    endif
    call s:patch(s:command(join(getline(1,'$'), "\n")))
endfunction

function! ShadenPatchSelection()
    call s:patch(s:command(s:get_visual_selection()))
endfunction

function! ShadenPatchLine()
    call s:patch(s:command(getline('.')))
endfunction

function! s:patch(cmd)
    let result = systemlist(a:cmd)
    if v:shell_error
        echo s:error_message
        return
    endif
    for LINE in result
        echo LINE
    endfor
endfunction

function! s:escape(str)
    return substitute(a:str, '"', '\\\"', 'g')
endfunction

function! s:command(content)
    return printf('curl -sSfL http://%s/eval -d "%s"', g:shaden_http_addr, s:escape(a:content))
endfunction

function! s:get_visual_selection()
    let [line_start, column_start] = getpos("'<")[1:2]
    let [line_end, column_end] = getpos("'>")[1:2]
    let lines = getline(line_start, line_end)
    if len(lines) == 0
        return ''
    endif
    let lines[-1] = lines[-1][: column_end - (&selection == 'inclusive' ? 1 : 2)]
    let lines[0] = lines[0][column_start - 1:]
    return join(lines, "\n")
endfunction

command! ShadenPatchSelection call ShadenPatchSelection()
command! ShadenPatchLine call ShadenPatchLine()
command! ShadenRepatch call ShadenRepatch()
