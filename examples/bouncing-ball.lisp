(define clock (unit/clock))

(-> clock (table :tempo (hz 0.2)))

(define slope1 (unit/slope))
(define fall (unit/mult))
(define slope2 (unit/slope))

; slope1 controls the "gravity"
(-> slope1 
    (table :rise (ms 10)
           :fall (ms 8000)
           :ratio 0.02
           :trigger (<- clock)))

(-> fall (table :x (ms 450) :y (<- slope1)))

; slope2 cycles; the fall for each cycle will get shorter and shorter based on the state of slope1
(-> slope2
    (table :rise (ms 10)
           :fall (<- fall)
           :ratio 0.001
           :trigger (<- clock)
           :cycle 1))


(define source (unit/gen))
(define amp (unit/mult))
(define gate (unit/gate))
(define filter (unit/filter))

(-> source (table :freq (hz "C4")))
(-> filter (table :in (<- source :pulse) 
                  :cutoff (hz 800)
                  :res 5))

(-> amp 
    (table :x (<- slope1) 
           :y (<- slope2)))

(-> gate
    (table :in (<- filter :bp)
           :control (<- amp)
           :cutoff-high (hz 3000)))

(emit (<- gate))
