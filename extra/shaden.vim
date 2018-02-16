" To install:
"
"   # mkdir -p $HOME/.vim/plugin
"   # cp extra/shaden.vim $HOME/.vim/plugin/shaden.vim
"
" Add these lines to your vimrc:
"
"   " Send a visual block of code to Shaden for evaluation
"   vnoremap <C-S-P> :<C-U>ShadenPatchSelection<CR>
"
"   " Send a line of code to Shaden for evaluation
"   nnoremap <C-S-P> :<C-U>ShadenPatchLine<CR>

if (exists("g:loaded_shaden"))
    finish
endif
let g:loaded_shaden = 1

if (!exists("g:shaden_http_addr"))
    let g:shaden_http_addr = '127.0.0.1:5000'
endif

function! ShadenRepatch()
    let content = s:escape("(clear)\n" . join(getline(1,'$'), "\n"))
    for LINE in systemlist(s:command(content))
        echo LINE
    endfor
endfunction

function! ShadenPatchSelection()
    let content = s:escape(s:get_visual_selection())
    for LINE in systemlist(s:command(content))
        echo LINE
    endfor
endfunction

function! ShadenPatchLine()
    let content = s:escape(getline('.'))
    for LINE in systemlist(s:command(content))
        echo LINE
    endfor
endfunction

function! s:escape(str)
    return substitute(a:str, '"', '\\\"', 'g')
endfunction

function! s:command(content)
    return printf('curl -sfL http://%s/eval -d "%s"', g:shaden_http_addr, a:content)
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
