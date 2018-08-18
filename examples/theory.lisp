; Shift pitches up a half step.
(define (shift-semitone pitches)
  (map (fn (_ p)
           (theory/transpose p (theory/interval :semitone 1)))
       pitches))

; Interleave two lists.
(define (interleave x y)
  (if (empty? x) y
    (cons (first x)
          (interleave y (rest x)))))

; Create the two whole-tone scales for this key; one ascending and the other
; descending.
(define root (theory/pitch "Eb4"))
(define asc (theory/scale root :whole-tone 1))
(define desc (shift-semitone (reverse asc)))

; Interleave the ascending and descending lists and convert to Hz.
(define freqs
  (map (fn (_ p) (hz p))
       (interleave asc desc)))

(define clock (unit/clock))
(define gen (unit/gen))
(define switch (unit/switch (table :size (len freqs))))
(define slope (unit/slope))
(define gate (unit/gate))

(-> clock (table :tempo (bpm 5)))

(-> switch
    freqs
    (table :trigger (<- clock)))

(-> gen (table :freq (<- switch)))

(-> slope
    (table :trigger (<- clock)
           :rise (ms 10)
           :fall (ms 300)))

(-> gate
    (table :in (<- gen :triangle)
           :control (<- slope)))

(emit (<- gate))
