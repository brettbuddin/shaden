(define (make-drum)
  (let ((stretch (unit/slope))
        (stretch-mult (unit/mult))
        (gen (unit/gen))
        (wave (unit/mix))
        (tone-slope (unit/slope))
        (tone-gate (unit/gate))
        (noise-slope (unit/slope))
        (noise-gate (unit/gate))
        (mix (unit/mix (table :size 2)))
        (distort (unit/overload)))

    ; route signal paths
    (-> stretch (table :ratio 0.001))
    (-> stretch-mult (table :x (<- stretch)))
    (-> gen (table :freq-mod (<- stretch-mult)))
    (-> wave (list (table :in (<- gen :sine))
                   (table :in (<- gen :triangle))))
    (-> tone-slope (table :ratio 0.001))
    (-> noise-slope (table :ratio 0.001))
    (-> tone-gate (table :in (<- wave) :control (<- tone-slope)))
    (-> noise-gate (table :in (<- gen :noise) :control (<- noise-slope)))
    (-> mix (list (table :in (<- tone-gate))
                  (table :in (<- noise-gate))))
    (-> distort (table :in (<- mix)))

    ; collect members into a list for unmounting them if we choose to
    (define members 
      (list stretch stretch-mult gen wave tone-slope tone-gate noise-slope noise-gate mix distort))

    ; return functions for obtaining the final output and for sparsly setting configuration inputs
    (table :out (fn () (<- distort))
           :unmount (fn () (map (fn (i v) (unit-unmount v)) members))
           :set (fn (opts)
                    (=> distort (table :gain (opts :gain)))
                    (=> stretch (table :trigger (opts :trigger) :rise (opts :stretch-rise) :fall (opts :stretch-fall)))
                    (=> stretch-mult (table :y (opts :stretch-amount)))
                    (=> gen (table :freq (opts :pitch) :sync (opts :trigger)))
                    (=> tone-slope (table :trigger (opts :trigger) :rise (opts :tone-rise) :fall (opts :tone-fall)))
                    (=> tone-gate (table :cutoff-high (opts :tone-cutoff)))
                    (=> noise-slope (table :trigger (opts :trigger) :rise (opts :noise-rise) :fall (opts :noise-fall)))
                    (=> noise-gate (table :cutoff-high (opts :noise-cutoff-high) :cutoff-low (opts :noise-cutoff-low)))))))

(define clock (unit/clock))
(define kick (make-drum))
(define gain (unit/mult))

((:set kick) 
  (table :gain 1
         :stretch-rise (ms 1)
         :stretch-fall (ms 70)
         :stretch-amount (hz 400)
         :pitch (hz :C3)
         :tone-rise (ms 1)
         :tone-fall (ms 500)
         :tone-cutoff (hz 3000)
         :noise-rise (ms 1)
         :noise-fall (ms 50)
         :noise-cutoff-high (hz 2500)
         :trigger (<- clock)))

(-> gain (table :x ((:out kick)) :y (db -6)))

(emit (<- gain))
