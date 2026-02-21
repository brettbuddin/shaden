package builtin

import (
	"bytes"
	"testing"

	"github.com/brettbuddin/shaden/lisp"
)

func newBenchEnv() *lisp.Environment {
	env := lisp.NewEnvironment()
	Load(env)
	return env
}

func mustParse(code string) any {
	node, err := lisp.Parse(bytes.NewBufferString(code))
	if err != nil {
		panic(err)
	}
	return node
}

func BenchmarkEval(b *testing.B) {
	b.Run("arithmetic", func(b *testing.B) {
		env := newBenchEnv()
		node := mustParse(`(+ 1 2 3 4 5)`)
		b.ReportAllocs()
		for b.Loop() {
			env.Eval(node)
		}
	})

	b.Run("arithmetic_float", func(b *testing.B) {
		env := newBenchEnv()
		node := mustParse(`(+ 1.0 2.0 3.0 4.0 5.0)`)
		b.ReportAllocs()
		for b.Loop() {
			env.Eval(node)
		}
	})

	b.Run("nested_arithmetic", func(b *testing.B) {
		env := newBenchEnv()
		node := mustParse(`(+ (* 2 3) (- 10 4) (/ 8 2))`)
		b.ReportAllocs()
		for b.Loop() {
			env.Eval(node)
		}
	})

	b.Run("function_call", func(b *testing.B) {
		env := newBenchEnv()
		env.Eval(mustParse(`(define (add x y) (+ x y))`))
		node := mustParse(`(add 5 10)`)
		b.ReportAllocs()
		for b.Loop() {
			env.Eval(node)
		}
	})

	b.Run("let_binding", func(b *testing.B) {
		env := newBenchEnv()
		node := mustParse(`(let ((x 1) (y 2) (z 3)) (+ x y z))`)
		b.ReportAllocs()
		for b.Loop() {
			env.Eval(node)
		}
	})

	b.Run("conditional", func(b *testing.B) {
		env := newBenchEnv()
		node := mustParse(`(if true 1 2)`)
		b.ReportAllocs()
		for b.Loop() {
			env.Eval(node)
		}
	})

	b.Run("list_operations", func(b *testing.B) {
		env := newBenchEnv()
		node := mustParse(`(reverse (append (list 1 2 3) 4 5))`)
		b.ReportAllocs()
		for b.Loop() {
			env.Eval(node)
		}
	})

	b.Run("map", func(b *testing.B) {
		env := newBenchEnv()
		node := mustParse(`(map (fn (i v) (+ v 1)) (list 1 2 3 4 5))`)
		b.ReportAllocs()
		for b.Loop() {
			env.Eval(node)
		}
	})

	b.Run("tail_call_recursion", func(b *testing.B) {
		env := newBenchEnv()
		env.Eval(mustParse(`(define (countdown n) (if (= n 0) 0 (countdown (- n 1))))`))
		node := mustParse(`(countdown 1000)`)
		b.ReportAllocs()
		for b.Loop() {
			env.Eval(node)
		}
	})

	b.Run("symbol_lookup", func(b *testing.B) {
		env := newBenchEnv()
		node := mustParse(`(let ((a 1)) (let ((b 2)) (let ((c 3)) (let ((d 4)) (let ((e 5)) (+ a b c d e))))))`)
		b.ReportAllocs()
		for b.Loop() {
			env.Eval(node)
		}
	})
}

func BenchmarkParse(b *testing.B) {
	b.Run("simple", func(b *testing.B) {
		code := []byte(`(+ 1 2)`)
		b.ReportAllocs()
		for b.Loop() {
			lisp.Parse(bytes.NewBuffer(code))
		}
	})

	b.Run("complex", func(b *testing.B) {
		code := []byte(`
(define (factorial n acc)
  (if (= n 0)
    acc
    (factorial (- n 1) (* acc n))))
(define (sum-list lst)
  (reduce (fn (r i v) (+ r v)) 0 lst))
(let ((x (factorial 10 1))
      (y (sum-list (list 1 2 3 4 5))))
  (+ x y))
`)
		b.ReportAllocs()
		for b.Loop() {
			lisp.Parse(bytes.NewBuffer(code))
		}
	})
}
