; Reverse a list
(define (reverse lst)
  (define (reverse-accum lst accum)
    (if (empty? lst) accum
      (reverse-accum (rest lst) (cons (first lst) accum))))
  (reverse-accum lst nil))

; Shift pitches up a half step
(define (shift _ p) 
  (theory/transpose p (theory/interval :semitone 1)))

(define (interleave x y)
  (if (empty? x) y
    (cons (first x)
          (interleave y (rest x)))))

; Ascending and descending whole-tone scales
(define root (theory/pitch "Eb4"))
(define asc (theory/scale root :whole-tone 1))
(define desc (map shift (reverse asc)))

; Interleave the ascending and descending lists and convert to Hz
(define freqs 
  (map (fn (_ p) (hz p)) 
       (interleave asc desc)))

(define clock (unit/clock))
(define gen (unit/gen))
(define switch (unit/switch (table :size (len freqs))))
(define slope (unit/slope))
(define gate (unit/gate))

(-> clock (table :tempo (hz 5)))
(-> switch (table :trigger (<- clock)) freqs)
(-> gen (table :freq (<- switch)))

(-> slope (table :trigger (<- clock)
                 :rise (ms 10)
                 :fall (ms 300)))

(-> gate (table :in (<- gen :triangle)
                :control (<- slope)))

(emit (<- gate))
