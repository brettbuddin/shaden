(define clock (unit/clock))
(define sample (unit/sample (table :file "hat.wav")))

(-> sample (table :trigger (<- clock) :begin 0.4 :end 0.5))

(emit (<- sample :a) (<- sample :b))
