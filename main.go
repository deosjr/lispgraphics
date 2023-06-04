package main

import (
    "math"

    "github.com/deosjr/whistle/lisp"
	"github.com/faiface/pixel/pixelgl"
)

func main() {
	pixelgl.Run(run)
}

func run() {
    l := lisp.New()
    LoadPixel(l.Env)
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
    // turtle properties
    l.Eval("(define turtle-line-width 5)")
    l.Eval("(define turtle-pen-colour red)")
    l.Eval("(define turtle-pen-down #t)")
    l.Eval("(define turtle-pos (cons (/ win-w 2) (/ win-h 2)))")
    l.Eval("(define turtle-heading 0)")

    l.Eval(`(define mod-fixed (lambda (n m)
        (if (> n 0) (mod n m) (mod (+ n m) m)))) `)

    l.Eval(`(define turn (lambda (radians)
        (set! turtle-heading (mod-fixed (+ turtle-heading radians) (* 2 pi)))))`)
    l.Eval("(define left (lambda () (turn (/ pi 2))))")
    l.Eval("(define right (lambda () (turn (- 0 (/ pi 2)))))")

    l.Eval(`(define forward (lambda (n) (let ((px (car turtle-pos)) (py (cdr turtle-pos)) (vx (cos turtle-heading)) (vy (sin turtle-heading)))
        (let ((newx (+ px (* n vx))) (newy (+ py (* n vy))))
            (let ((newpos (cons newx newy)))
                (if (eqv? turtle-pen-down #t) (draw-line turtle-pos newpos))
                (set! turtle-pos newpos)
                (wraparound)
    )))))`)

    l.Eval(`(define wraparound (lambda () (let ((px (car turtle-pos)) (py (cdr turtle-pos)))
        (set! turtle-pos (cons (mod-fixed px win-w) (mod-fixed py win-h)))
    )))`)

    l.Eval(`(define draw-line (lambda (from to)
        (let ((imd (imdraw)))
            (im-set-color! imd turtle-pen-colour)
            (im-push imd (vec2d (car from) (cdr from)))
            (im-push imd (vec2d (car to) (cdr to)))
            (line imd turtle-line-width)
            (im-draw imd win)
            (update win)
    )))`)

    l.Eval("(define gosper3 (quote (f r f r r f l f l l f f l f r r l f r f f r r f r f l l f l f r r l f r f f r r f r f l l f l f l f r f r r f l f l l f f l f r l l f r f r r f l f l l f f l f r f r f r r f l f l l f f l f r l l f r f f r r f r f l l f l f r r l f r f r r f l f l l f f l f r r l f r f f r r f r f l l f l f l f r f f r r f r f l l f l f r r l f r f f r r f r f l l f l f r f r f r r f l f l l f f l f r l l f r f r r f l f l l f f l f r l l f r f f r r f r f l l f l f r r l f r f r r f l f l l f f l f r r l f r f f r r f r f l l f l f l f r f f r r f r f l l f l f r r l f r f f r r f r f l l f l f r f r f r r f l f l l f f l f r l l f r f r r f l f l l f f l f r l l f r f f r r f r f l l f l f l f r f r r f l f l l f f l f r r l f r f f r r f r f l l f l f r r l f r f f r r f r f l l f l f l f r f r r f l f l l f f l f r l l f r f r r f l f l l f f l f r f r f r r f l f l l f f l f r l l f r f f r r f r f l l f l f r l l f r f r r f l f l l f f l f r r l f r f f r r f r f l l f l f r r l f r f f r r f r f l l f l f l f r f r r f l f l l f f l f r l l f r f r r f l f l l f f l f r f r f r r f l f l l f f l f r l l f r f f r r f r f l l f l f r f r f r r f l f l l f f l f r r l f r f f r r f r f l l f l f r r l f r f f r r f r f l l f l f l f r f r r f l f l l f f l f r l l f r f r r f l f l l f f l f r f r f r r f l f l l f f l f r l l f r f f r r f r f l l f l f r l l f r f r r f l f l l f f l f r r l f r f f r r f r f l l f l f l f r f f r r f r f l l f l f r r l f r f f r r f r f l l f l f r f r f r r f l f l l f f l f r l l f r f r r f l f l l f f l f r l l f r f f r r f r f l l f l f r)))")
    l.Eval("(define program (quote ()))")

    l.Eval(`(define start (lambda ()
        (begin
        (set! program gosper3)
        (set! tick (lambda () (begin
            (if (null? program) (set! 'program (quote (done))))
            (let ((next (car program)) (rem (cdr program)))
                (cond
                    ((eqv? next 'f) (forward 10))
                    ((eqv? next 'l)  (turn (/ pi 3)))
                    ((eqv? next 'r) (turn (- 0 (/ pi 3))))
                    (else #t))
                (set! program rem)
    )))))))`)
}
