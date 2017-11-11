(define oscillator (unit/gen))
(define mix (unit/mix (table :size 3)))

; Mix a single oscillator's pulse, saw and sub-pulse outputs. 
; The saw output is attenuated by -12dB and the sub-pulse output 
; is attenuated by -3dB. All outputs are attenuated by -3dB via
; the master level input.

(-> oscillator (table :freq (hz 300)))

(-> mix 
    (table :master (db -3))
    (list
      (table :in (<- oscillator :pulse))
      (table :in (<- oscillator :saw) :level (db -12))
      (table :in (<- oscillator :sub-pulse) :level (db -3))))

(emit (<- mix))
