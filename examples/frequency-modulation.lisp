(define modulator (unit/gen))
(define carrier (unit/gen))

; Modulate the carrier's frequency by 40Hz at a rate of 5Hz using
; the modulating oscillator's sine wave output.

(-> modulator 
    (table :freq (hz 5) 
           :amp (hz 40)))

(-> carrier 
    (table :freq (hz 300) 
           :freq-mod (<- modulator :sine)))

(emit (<- carrier :sine))
