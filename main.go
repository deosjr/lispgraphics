package main

import (
    "github.com/deosjr/lispadventures/lisp"
	"github.com/faiface/pixel/pixelgl"
)

func main() {
	pixelgl.Run(run)
}

func run() {
    p, env := lisp.New()
    LoadPixel(env)
	env.AddBuiltin("/", func(args []lisp.SExpression) (lisp.SExpression, error) {
		return lisp.NewPrimitive(args[0].AsNumber() / args[1].AsNumber()), nil
	})

    // TODO: this feels too hacky/leaky. lisp.New should abstract over p,env?
    p.Eval(env, "(define env (environment))")
    envexp, _ := p.Eval(env, "env")
    drawenv := envexp.AsPrimitive().(*lisp.Env)

    // TODO hack! this isnt behaviour of LISP set!, this is overriding on toplevel always!
    // Should either implement set! properly or pass turtle _everywhere_
    // usage: (set! 'top 3)
    p.Eval(drawenv, "(define set! (lambda (s v) (eval (list 'define s (quasiquote (quote ,v))) env) ))")
    p.Eval(drawenv, "(define not (lambda (t) (if (eqv? t #t) #f #t)))")

    // these can be defined in toplevel env, but this makes more clear whats needed
    p.Eval(drawenv, "(define refresh 100)")

    p.Eval(env, "(define win-w 1024)")
    p.Eval(env, "(define win-h 768)")
    p.Eval(env, "(define win (window))")
    
    // TODO: lisp.process not exported
    //loadTurtleGraphics(p, drawenv)
    // START TURTLE GRAPHICS
    // turtle properties
    p.Eval(drawenv, "(define turtle-line-width 5)")
    p.Eval(drawenv, "(define turtle-pen-colour red)")
    p.Eval(drawenv, "(define turtle-pen-down #t)")
    p.Eval(drawenv, "(define turtle-pos (cons (/ win-w 2) (/ win-h 2)))")
    p.Eval(drawenv, "(define turtle-heading (cons 0 1))")

    p.Eval(drawenv, `(define left (lambda () (let ((vx (car turtle-heading)) (vy (cdr turtle-heading)))
        (set! 'turtle-heading (cons (* vy -1) vx))
    )))`)
    p.Eval(drawenv, `(define right (lambda () (let ((vx (car turtle-heading)) (vy (cdr turtle-heading)))
        (set! 'turtle-heading (cons vy (- 0 vx)))
    )))`)
    p.Eval(drawenv, `(define forward (lambda (n) (let ((px (car turtle-pos)) (py (cdr turtle-pos)) (vx (car turtle-heading)) (vy (cdr turtle-heading)))
        (let ((newx (+ px (* n vx))) (newy (+ py (* n vy))))
            (let ((newpos (cons newx newy)))
                (if (eqv? turtle-pen-down #t) (draw-line turtle-pos newpos))
                (set! 'turtle-pos newpos)
                (wraparound)
    )))))`)
    p.Eval(drawenv, `(define wraparound (lambda () (let ((px (car turtle-pos)) (py (cdr turtle-pos)))
        (let ((newx (cond ((< px 0) (+ px win-w)) ((> px win-w) (- px win-w)) (else px)))
              (newy (cond ((< py 0) (+ py win-h)) ((> py win-h) (- py win-h)) (else py))))
        (set! 'turtle-pos (cons newx newy))
    ))))`)

    p.Eval(drawenv, `(define draw-line (lambda (from to)
        (let ((imd (imdraw)))
            (im-set-color! imd turtle-pen-colour)
            (im-push imd (vec2d (car from) (cdr from)))
            (im-push imd (vec2d (car to) (cdr to)))
            (line imd turtle-line-width)
            (im-draw imd win)
            (update win)
    )))`)

    p.Eval(drawenv, `(define start (lambda ()
        (set! 'tick (lambda () (begin
            (forward 10)
            (set! 'turtle-pen-down (not turtle-pen-down)
    ))))))`)
    // END TURTLE GRAPHICS

    p.Eval(drawenv, `(define tick (lambda () #t))`)

	p.Eval(drawenv, `(define drawrec (lambda ()
        (receive
            ((in) (quasiquote (input ,in)) ->
                (eval in env) (eval (quote (drawrec)) env))
            (after refresh -> (if (closed? win) #t
                (begin (tick) (drawrec)))))))`)


    // TODO: if drawpid dies and repl isnt killed, things are broken
    // if repl is killed, things are still broken (old repl is still reading)
    // fix: 'read' shouldnt be blocking somehow? -> havent managed to figure out how yet
    // other fix: REPL sends to named pid!
	p.Eval(env, `(define REPL (lambda (drawpid)
        (begin (display "> ")
            (send drawpid (quasiquote (input (unquote (read)))))
            (REPL drawpid))))`)

	p.Eval(env, `(define restarter (lambda ()
        (begin (process_flag 'trap_exit #t)
               (let ((drawpid (spawn_link (lambda () (eval (quote (drawrec)) env)) (quote ()))))
               (let ((repl (spawn_link REPL (quote (drawpid)))))
                    (receive
                        ((reason) (quasiquote (EXIT ,repl ,reason)) -> (display "REPL: ") (display reason) (display newline))
                        ((reason) (quasiquote (EXIT ,drawpid ,reason)) ->
                            (if (eqv? reason "normal") #t
                            (begin (display "** exception error: ") (display reason) (display newline) (restarter))))))))))`)

	p.Eval(env, "(restarter)")
}

