(define clock (unit/clock))

(-> clock (table :tempo (hz 9)))

(define euclid1 (unit/euclid))
(define euclid2 (unit/euclid))

(define logic1 (unit/logic))
(define logic2 (unit/logic))
(define logic3 (unit/logic))

(-> euclid1 (table :span 8 :fill 4 :offset 0 :clock (<- clock)))
(-> euclid2 (table :span 8 :fill 3 :offset 4 :clock (<- clock)))

(-> logic1 (table :x (<- euclid1) :y (<- euclid2) :mode logic/and))
(-> logic2 (table :x (<- euclid1) :y (<- euclid2) :mode logic/or))
(-> logic3 (table :x (<- logic1) :y (<- logic2) :mode logic/xor))

(define voice1-gen (unit/gen))
(define voice1-slope (unit/slope))
(define voice1-gate (unit/gate))

(define voice2-gen (unit/gen))
(define voice2-slope (unit/slope))
(define voice2-gate (unit/gate))

(define voice3-gen (unit/gen))
(define voice3-slope (unit/slope))
(define voice3-gate (unit/gate))

(-> voice1-gen (table :freq (hz "C2")))
(-> voice1-slope (table :rise (ms 10) :fall (ms 1000) :trigger (<- logic1)))
(-> voice1-gate 
    (table :control (<- voice1-slope) 
           :in (<- voice1-gen :saw) 
           :cutoff-high (hz 700)))

(-> voice2-gen (table :freq (hz "F3")))
(-> voice2-slope (table :rise (ms 1) :fall (ms 1000) :trigger (<- logic2)))
(-> voice2-gate 
    (table :control (<- voice2-slope) 
           :in (<- voice2-gen :pulse) 
           :cutoff-high (hz 700)))

(-> voice3-gen (table :freq (hz "Eb3")))
(-> voice3-slope (table :rise (ms 1) :fall (ms 1000) :trigger (<- logic3)))
(-> voice3-gate 
    (table :control (<- voice3-slope) 
           :in (<- voice3-gen :saw) 
           :cutoff-high (hz 500)))

(define mix (unit/mix))

(-> mix
    (table :master (db -24))
    (list 
      (table :in (<- voice1-gate))
      (table :in (<- voice2-gate))
      (table :in (<- voice3-gate))))

(define delay (unit/delay))
(define delay-filter (unit/filter))
(define delay-lfo (unit/low-gen))

(-> delay-lfo (table :freq (hz 0.1) :amp (ms 10) :offset (ms 115)))
(-> delay 
    (table :in (<- mix) 
           :time (<- delay-lfo :sine) 
           :fb-return (<- delay-filter :bp) 
           :fb-gain 0.9 
           :mix 0.3))
(-> delay-filter (table :in (<- delay :fb-send) :cutoff (hz 500)))

(emit (<- delay))
