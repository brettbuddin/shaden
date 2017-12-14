(define clock (unit/clock))
(-> clock (table :tempo (hz 1)))

(define voice1-gen (unit/gen))
(define voice1-slope (unit/slope))
(define voice1-gate (unit/gate))

(-> voice1-gen (table :freq (hz "C4")))
(-> voice1-slope (table :rise (ms 10) :fall (ms 1000) :trigger (<- clock)))
(-> voice1-gate 
    (table :control (<- voice1-slope) 
           :in (<- voice1-gen :saw) 
           :cutoff-high (hz 1500)))

(define reverb (unit/reverb))

(-> reverb 
    (table :a (<- voice1-gate)
           :b (<- voice1-gate)
           :mix 0.4
           :decay 0.9
           :size 0.3
           :shift-semitones 2
           :cutoff-pre (hz 2000)
           :cutoff-post (hz 300)))

(define ampl (unit/mult))
(define ampr (unit/mult))

(-> ampl (table :x (<- reverb :a) :y (db -6)))
(-> ampr (table :x (<- reverb :b) :y (db -6)))

(emit (<- ampl) (<- ampr))
