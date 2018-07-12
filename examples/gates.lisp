; Clock

(define clock (unit/clock))

(-> clock (table :tempo (hz 8)))

; Modulation

(define lfo-bipolar (unit/low-gen))
(define lfo-unipolar-sine (unit/adjust))

(-> lfo-bipolar
    (table :freq (hz 0.5)))

(-> lfo-unipolar-sine
    (table :in (<- lfo-bipolar :sine)
           :mult 0.5
           :add 1))

; Voice

(define sequence (unit/stages))
(define oscillator (unit/gen))
(define pulse-width (unit/adjust))
(define mix (unit/mix (table :size 3)))
(define slope (unit/slope))
(define gate (unit/gate))

(-> sequence
    (table :clock (<- clock))
    (list (table :freq (hz "C3") :mode mode/first :pulses 1)
          (table :freq (hz "C3") :mode mode/first :pulses 1)
          (table :freq (hz "C3") :mode mode/first :pulses 1)
          (table :freq (hz "C3") :mode mode/first :pulses 1)
          (table :freq (hz "C4") :mode mode/first :pulses 1)))

(-> pulse-width
    (table :in (<- lfo-unipolar-sine)
           :mult 0.5))

(-> oscillator
    (table :freq (<- sequence :freq)
           :pulse-width (<- pulse-width)))

(-> mix
    (table :master (db -12))
    (list
      (table :in (<- oscillator :pulse))
      (table :in (<- oscillator :saw) :level (db -6))
      (table :in (<- oscillator :sub-pulse))))

; Envelope

(-> slope
    (table :rise (ms 1)
           :fall (ms 700)
           :trigger (<- sequence :gate)))

(-> gate
    (table :in (<- mix)
           :control (<- slope)
           :cutoff-high (hz 600)))

(emit (<- gate))
