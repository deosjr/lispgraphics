package main

import (
	_ "embed"
	"math"

	"github.com/deosjr/lispgraphics/pixel"
	"github.com/deosjr/whistle/erlang"
	"github.com/deosjr/whistle/kanren"
	"github.com/deosjr/whistle/lisp"
	"github.com/faiface/pixel/pixelgl"
)

//go:embed turtle.lisp
var turtle string

func main() {
	pixelgl.Run(run)
}

func run() {
	l := lisp.New()

	// needed for the repl to work
	kanren.Load(l)
	erlang.Load(l)

	pixel.Load(l)
	l.Env.AddBuiltin("sin", func(args []lisp.SExpression) (lisp.SExpression, error) {
		return lisp.NewPrimitive(math.Sin(args[0].AsNumber())), nil
	})
	l.Env.AddBuiltin("cos", func(args []lisp.SExpression) (lisp.SExpression, error) {
		return lisp.NewPrimitive(math.Cos(args[0].AsNumber())), nil
	})

	l.Eval("(define win-w 1024)")
	l.Eval("(define win-h 768)")
	l.Eval("(define win (window))")
	l.Eval("(define refresh 100)")

	// TODO hack! this isnt behaviour of LISP set!, this is overriding on toplevel always!
	// Should either implement set! properly or pass turtle _everywhere_
	// usage: (set! 'top 3)
	//l.Eval("(define set! (lambda (s v) (eval (list 'define s (quasiquote (quote ,v))) env) ))")
	//l.Eval("(define not (lambda (t) (if (eqv? t #t) #f #t)))")

	// this is still pretty hacky...
	l.Eval("(define env (environment))")
	envexp, _ := l.Eval("env")
	drawenv := envexp.AsPrimitive().(*lisp.Env)
	oldenv := l.Env
	l.Env = drawenv
	loadTurtleGraphics(l)
	l.Env = oldenv

	l.Eval(`(define tick (lambda () #t))`)

	l.Eval(`(define drawrec (lambda ()
        (receive
            ((in) (quasiquote (input ,in)) ->
                (eval in env) (eval (quote (drawrec)) env))
            (after refresh -> (if (closed? win) #t
                (begin (tick) (drawrec)))))))`)

	// TODO: if drawpid dies and repl isnt killed, things are broken
	// if repl is killed, things are still broken (old repl is still reading)
	// fix: 'read' shouldnt be blocking somehow? -> havent managed to figure out how yet
	// other fix: REPL sends to named pid!
	l.Eval(`(define REPL (lambda (drawpid)
        (begin (display "> ")
            (send drawpid (quasiquote (input (unquote (read)))))
            (REPL drawpid))))`)

	l.Eval(`(define restarter (lambda ()
        (begin (process_flag 'trap_exit #t)
               (let ((drawpid (spawn_link (lambda () (eval (quote (drawrec)) env)) (quote ()))))
               (let ((repl (spawn_link REPL (quote (drawpid)))))
                    (receive
                        ((reason) (quasiquote (EXIT ,repl ,reason)) -> (display "REPL: ") (display reason) (display newline))
                        ((reason) (quasiquote (EXIT ,drawpid ,reason)) ->
                            (if (eqv? reason "normal") #t
                            (begin (display "** exception error: ") (display reason) (display newline) (restarter))))))))))`)

	l.Eval("(restarter)")
}

func loadTurtleGraphics(l lisp.Lisp) {
	err := l.Load(turtle)
	if err != nil {
		panic(err)
	}
}
