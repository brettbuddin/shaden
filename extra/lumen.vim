" To install:
"
"   # mkdir -p $HOME/.vim/plugin
"   # cp extra/lumen.vim $HOME/.vim/plugin/lumen.vim
"
" Add these lines to your vimrc:
"
"   " Send a visual block of code to Lumen for evaluation
"   vnoremap <C-S-P> :<C-U>LumenPatchSelection<CR>
"
"   " Send a line of code to Lumen for evaluation
"   nnoremap <C-S-P> :<C-U>LumenPatchLine<CR>

if (exists("g:loaded_lumen"))
    finish
endif
let g:loaded_lumen = 1

if (!exists("g:lumen_http_addr"))
    let g:lumen_http_addr = '127.0.0.1:5000'
endif

function! LumenPatchSelection()
    let content = s:escape(s:get_visual_selection())
	echom "lumen: " . system(s:command(content))
endfunction

function! LumenPatchLine()
    let content = s:escape(getline('.'))
	echom "lumen: " . system(s:command(content))
endfunction

function! s:escape(str)
    return substitute(a:str, '"', '\\\"', 'g')
endfunction

function! s:command(content)
    return printf('curl -sfL http://%s/eval -d "%s"', g:lumen_http_addr, a:content)
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

command! LumenPatchSelection call LumenPatchSelection()
command! LumenPatchLine call LumenPatchLine()
